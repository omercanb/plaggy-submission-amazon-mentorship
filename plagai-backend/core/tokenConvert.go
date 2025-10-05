package core

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/plagai/plagai-backend/middleware"
)

func ConvertToken(authHeader string, parsedClaims *middleware.Claims) {
	tokenStr := strings.Split(authHeader, " ")[1]
	jwt.ParseWithClaims(tokenStr, parsedClaims, func(token *jwt.Token) (any, error) {
		return nil, nil
	})
}