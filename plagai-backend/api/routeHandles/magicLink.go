package routeHandles

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	mailserver "github.com/plagai/plagai-backend/mail-server"
	"github.com/plagai/plagai-backend/middleware"
)

type MagicRequest struct {
	Email string `json:"email"`
}

type MagicTokenInfo struct {
	Email         string
	Authenticated bool
	MagicToken    string
}

var (
	userSessions = make(map[string]MagicTokenInfo) // key: magicToken
	mutex        sync.Mutex
)

func GenerateToken(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func MagicRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req MagicRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, "Email required", http.StatusBadRequest)
		return
	}

	magicToken, err := GenerateToken(20)
	if err != nil {
		http.Error(w, "Failed to generate magic token", http.StatusInternalServerError)
		return
	}

	mutex.Lock()
	userSessions[magicToken] = MagicTokenInfo{
		Email:         req.Email,
		Authenticated: false,
		MagicToken:    magicToken,
	}
	mutex.Unlock()

	fmt.Printf("Magic link for %s: http://plaggy.xyz/api/v1/auth/magic-consume?magic=%s\n", req.Email, magicToken)
	loginUrl := fmt.Sprintf("http://plaggy.xyz/api/v1/auth/magic-consume?magic=%s", magicToken)
	mailcontent := fmt.Sprintf("Login by going to this URL: %s", loginUrl)
	mailserver.SendMailWithPlaintextContent(req.Email, "Login From PlagAI", "Magic Link Login", mailcontent)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Magic link sent",
		"magicId": magicToken,
	})
}

func MagicStatusHandler(w http.ResponseWriter, r *http.Request) {
	magic := r.URL.Query().Get("magic")
	if magic == "" {
		http.Error(w, "Missing magic token", http.StatusBadRequest)
		return
	}

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" {
		http.Error(w, "Missing or invalid email in body", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	info, exists := userSessions[magic]
	mutex.Unlock()

	if !exists {
		http.Error(w, "Magic token not found", http.StatusNotFound)
		return
	}

	if info.Email != req.Email {
		http.Error(w, "Email does not match magic token", http.StatusUnauthorized)
		return
	}

	if !info.Authenticated {
		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}

	expirationTime := time.Now().Add(72 * time.Hour)
	claims := &middleware.Claims{
		Email: req.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	userToken, err := jwtToken.SignedString(middleware.JWTKey)
	if err != nil {
		http.Error(w, "Failed to create user token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"token":         userToken,
	})
}

func MagicConsumeHandler(w http.ResponseWriter, r *http.Request) {
	magic := r.URL.Query().Get("magic")
	fmt.Println("Attempting to consume magic token:", magic)

	if magic == "" {
		http.Error(w, "Missing magic token", http.StatusBadRequest)
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	info, exists := userSessions[magic]
	if !exists {
		http.Error(w, "Invalid or expired magic token", http.StatusNotFound)
		return
	}

	info.Authenticated = true
	userSessions[magic] = info
	fmt.Fprint(w, "You are now logged in! You can close this browser.")
}
