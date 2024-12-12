package auth

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/entity"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

type TokenManager struct {
	signingKey string
}

func NewTokenManager() *TokenManager {
	signingKey, exists := os.LookupEnv("JWT_SECRET_KEY")
	if !exists {
		log.Printf("unable to parse jwt secret key from environment")
		return nil
	}
	return &TokenManager{signingKey: signingKey}
}

func (tm *TokenManager) ValidateAccessToken(tokenString string) (*jwt.Token, error) {
	jwtToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(tm.signingKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	return jwtToken, nil
}

func (tm *TokenManager) SignToken(token *jwt.Token) (string, error) {
	jwtTokenString, err := token.SignedString([]byte(tm.signingKey))
	if err != nil {
		return "", fmt.Errorf("unable to sign jwt token: %s", err)
	}
	return jwtTokenString, nil
}

func (tm *TokenManager) GenerateAccessToken(userID string, signingMethod jwt.SigningMethod) (string, error) {
	payload := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Minute * 15).Unix(),
	}
	accessToken := jwt.NewWithClaims(signingMethod, payload)
	signedAccessToken, err := tm.SignToken(accessToken)
	if err != nil {
		return "", err
	}
	return signedAccessToken, nil
}

func (tm *TokenManager) CreateRefreshToken() (string, error) {
	tokenSlice := make([]byte, 32)
	_, err := rand.Read(tokenSlice)
	if err != nil {
		return "", fmt.Errorf("unable to generate bytes: %s", err)
	}
	tokenString := base64.URLEncoding.EncodeToString(tokenSlice)
	return tokenString, nil
}

func (tm *TokenManager) PostRefreshToken(refreshStorage database.RefreshStorage, userID string) (string, error) {
	refreshTokenString, err := tm.CreateRefreshToken()
	if err != nil {
		return "", err
	}
	newHashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshTokenString), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("unable to generate hashed refresh token: %s", err)
	}
	expTime := time.Now().Add(time.Minute * 43200).UTC() // 30 days refresh token is valid
	refreshToken := entity.RefreshToken{
		UserID:       userID,
		RefreshToken: string(newHashedToken),
		ExpiresAt:    expTime,
	}

	err = refreshStorage.Post(refreshToken)
	if err != nil {
		return "", err
	}
	return refreshTokenString, nil
}
func (tm *TokenManager) GenerateBothTokens(refreshStorage database.RefreshStorage, userID string) (string, string, error) {
	newAccessToken, err := tm.GenerateAccessToken(userID, jwt.SigningMethodHS512)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := tm.PostRefreshToken(refreshStorage, userID)
	if err != nil {
		return "", "", err
	}
	return newAccessToken, newRefreshToken, nil
}

func (tm *TokenManager) GetUserID(ctx *gin.Context) (string, bool) {
	userID, exists := ctx.Get("userID")
	if !exists {
		log.Println("invalid token credentials: userID is absent")
		ctx.Writer.WriteHeader(http.StatusUnauthorized)
		return "", false
	}
	userIDStr, ok := userID.(string)
	if !ok {
		log.Println("userID is not a string")
		ctx.Writer.WriteHeader(http.StatusUnauthorized)
		return "", false
	}
	return userIDStr, true
}
