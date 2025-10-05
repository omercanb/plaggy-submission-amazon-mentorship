package routeHandles

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/plagai/plagai-backend/middleware"
	"github.com/plagai/plagai-backend/models/database"
	"github.com/plagai/plagai-backend/security"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials LoginRequest

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var userFromDB database.Instructor
	if err := h.DB.
		Where("email LIKE ?", credentials.Email).
		First(&userFromDB).Error; err != nil {
		log.Println("user query failed, probably no user with email:", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	// We recieve a SHA256 hashed password from the client and compare that to the SHA256 + argon2 encypted
	// password we have in the db. So we never actually send plain text over-the-air. Also if we forget the 
	// passwords of the account we create, if we ever have a signup functionality, we are fucked. ~brtcrt
	ok, err := security.VerifyFromClientDigestHex(credentials.Password, userFromDB.Password);
	if !ok || err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middleware.Claims{
		Email: credentials.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JWTKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		Token: tokenString,
	}
	w.WriteHeader(http.StatusOK)
	// In fact, returning like this seems kinda dumb tbh. Maybe we should switch to cookies which would make
	// it a lot easier to manage on the frontend. Although I really don't give a shit at this point.
	w.Header().Set("token", tokenString) // This probably doesn't work bc we haven't allowed this header in main.go
	json.NewEncoder(w).Encode(response)
}
