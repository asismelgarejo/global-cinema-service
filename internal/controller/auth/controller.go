package auth

import (
	"api/pkg/model"
	"context"
	"errors"
	"os"
	"strings"
	"time"

	repoUser "api/internal/repository/mongodb/user"

	"github.com/dgrijalva/jwt-go"
)

type ControllerAuth struct {
	repoUser *repoUser.RepoUser
}

func New(repoUser *repoUser.RepoUser) *ControllerAuth {
	return &ControllerAuth{repoUser: repoUser}
}

func (a *ControllerAuth) SignIn(ctx context.Context, userData model.User) (*model.JWTOutput, error) {
	_, err := a.repoUser.GetUser(ctx, userData)
	if err != nil {
		return nil, err
	}
	expirationTime := time.Now().Add(time.Minute * 10)
	claims := &model.Claims{
		Username:       userData.Username,
		StandardClaims: jwt.StandardClaims{ExpiresAt: expirationTime.Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_TOKEN")))
	if err != nil {
		return nil, err
	}
	return &model.JWTOutput{Token: tokenString, Expires: expirationTime}, nil
}

func ValidToken(authHeader string) (*model.Claims, error) {
	// Check if the header is missing or doesn't start with "Bearer "
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, errors.New("token is missing or malformed")
	}

	// Extract the token from the Authorization header
	tokenVal := strings.TrimPrefix(authHeader, "Bearer ")

	claims := &model.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenVal, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_TOKEN")), nil
	})

	if err != nil {
		return nil, errors.New("token is missing or malformed")
	}
	if !tkn.Valid {
		return nil, errors.New("token is missing or malformed")
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		return nil, errors.New("Token is not expired yet")
	}
	return claims, nil
}

func (a *ControllerAuth) RefreshToken(ctx context.Context, authHeader string) (*model.JWTOutput, error) {
	claims, err := ValidToken(authHeader)
	if err != nil {
		return nil, err
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET"))
	if err != nil {
		return nil, err
	}
	return &model.JWTOutput{
		Token:   tokenString,
		Expires: expirationTime,
	}, nil
}
