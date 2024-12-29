package managers

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/entity"
	customLogger "CourseProject/auth_service/pkg/log"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

type TokenManager struct {
	signingKey string
}

func NewTokenManager(logger *customLogger.Logger) *TokenManager {
	signingKey, exists := os.LookupEnv("JWT_SECRET_KEY")
	if !exists {
		logger.ErrorLogger.Error().Msg("unable to parse jwt secret key from environment")
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

func (tm *TokenManager) GenerateRefreshToken() (string, error) {
	tokenSlice := make([]byte, 32)
	_, err := rand.Read(tokenSlice)
	if err != nil {
		return "", fmt.Errorf("unable to generate bytes: %s", err)
	}
	tokenString := base64.URLEncoding.EncodeToString(tokenSlice)
	return tokenString, nil
}

func (tm *TokenManager) GetHashedRefreshToken(refreshToken string) (string, error) {
	newHashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("unable to generate hashed refresh token: %s", err)
	}

	return string(newHashedToken), nil
}

func (tm *TokenManager) GenerateBothTokens(userID string) (newAccessToken string, newRefreshToken string, err error) {
	newAccessToken, err = tm.GenerateAccessToken(userID, jwt.SigningMethodHS512)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err = tm.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}
	return newAccessToken, newRefreshToken, nil
}

func (tm *TokenManager) GetUserID(ctx *gin.Context, logger *customLogger.Logger) (string, bool) {
	userID, exists := ctx.Get("userID")
	if !exists {
		logger.ErrorLogger.Error().Msg("invalid token credentials: userID is absent")
		ctx.Writer.WriteHeader(http.StatusUnauthorized)
		return "", false
	}
	userIDStr, ok := userID.(string)
	if !ok {
		logger.ErrorLogger.Error().Msg("userID is not a string")
		ctx.Writer.WriteHeader(http.StatusUnauthorized)
		return "", false
	}
	return userIDStr, true
}

func (tm *TokenManager) CompareRefreshTokens(refreshStorage database.RefreshStorage, inputToken, userID string) error {
	savedToken, err := refreshStorage.Get(userID)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	if err = bcrypt.CompareHashAndPassword([]byte(savedToken), []byte(inputToken)); err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func (tm *TokenManager) PostHashedRefreshToken(refreshStorage database.RefreshStorage, refreshToken, userID string) error {
	hashedRefreshToken, err := tm.GetHashedRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	expTime := time.Now().Add(time.Minute * 43200).UTC() // 30 days refresh_token is valid
	refreshTokenEntity := entity.RefreshToken{
		UserID:       userID,
		RefreshToken: hashedRefreshToken,
		ExpiresAt:    expTime,
	}
	if err = refreshStorage.Post(refreshTokenEntity); err != nil {
		return err
	}
	return nil
}
