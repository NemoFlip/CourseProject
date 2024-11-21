package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type TokenManager struct {
	signingKey string
}

func NewTokenManager(signingKey string) *TokenManager {
	return &TokenManager{signingKey: signingKey}
}

func (tm *TokenManager) ValidateToken(tokenString string) error {
	jwtToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return tm.signingKey, nil
	})
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("unable to get claims from token: %w", err)
	}
	_, ok = claims["sub"].(string)
	if !ok {
		return fmt.Errorf("unable to get `sub` claim from token: %w", err)
	}
	exp, ok := claims["exp"].(int64)
	if !ok {
		return fmt.Errorf("unable to get `exp` claim from token: %w", err)
	}
	currentTime := time.Now().Unix()
	if currentTime > exp {
		return fmt.Errorf("token is expired")
	}
	return nil
}

func (tm *TokenManager) SignToken(token *jwt.Token) (string, error) {
	jwtTokenString, err := token.SignedString([]byte(tm.signingKey))
	if err != nil {
		return "", fmt.Errorf("unable to sign jwt token: %s", err)
	}
	return jwtTokenString, nil

}
