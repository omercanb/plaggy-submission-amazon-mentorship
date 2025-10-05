package api

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/plagai/plagai-backend/middleware"
)

var (
	ErrMissingAuthHeader = errors.New("missing Authorization header")
	ErrInvalidToken      = errors.New("invalid or expired token")
)

func GetClaimsFromAuthorization(r *http.Request) (middleware.Claims, error) {
	// Parse JWT from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("Missing Authorization header", http.StatusUnauthorized)
		return middleware.Claims{}, ErrMissingAuthHeader
	}

	log.Printf("Authorization header: %q", authHeader)

	var tokenStr string
	fmt.Sscanf(authHeader, "Bearer %s", &tokenStr)
	claims := &middleware.Claims{}
	tkn, err := middleware.ParseJWT(tokenStr, claims)
	if err != nil || !tkn.Valid {
		log.Println("Invalid or expired token", http.StatusUnauthorized)
		return middleware.Claims{}, ErrInvalidToken
	}

	return *claims, nil
}
