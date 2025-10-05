package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var JWTKey []byte

// Claims = data we hold inside our token ~brtcrt
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func init() {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		secret = "böylesisteminamınıyurdunusikeyim"
	}
	JWTKey = []byte(secret)
}

func ParseJWT(tokenStr string, claims *Claims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		/*
			{
				Authorization: "Bearer tokenstringgoeshere"
			}
		*/
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := headerParts[1]

		claims := Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (any, error) {
			return JWTKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			log.Println("JWT error:", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !tkn.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		// log.Println("Authenticated user with email: ", claims.Email)

		next.ServeHTTP(w, r)
	})
}
