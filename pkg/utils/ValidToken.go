package utils

import (
	"api/pkg/model"
	"errors"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

func ValidToken(authHeader string) error {
	// Check if the header is missing or doesn't start with "Bearer "
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return errors.New("token is missing or malformed")
	}

	// Extract the token from the Authorization header
	tokenVal := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &model.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenVal, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_TOKEN")), nil
	})

	if err != nil {
		return errors.New("token is missing or malformed")
	}
	if !tkn.Valid {
		return errors.New("token is missing or malformed")
	}
	return nil
}
