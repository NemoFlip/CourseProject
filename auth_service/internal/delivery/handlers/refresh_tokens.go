package handlers

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/entity"
	customLogger "CourseProject/auth_service/pkg/log"
	"CourseProject/auth_service/pkg/managers"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type TokenServer struct {
	tokenManager   managers.TokenManager
	refreshStorage database.RefreshStorage
	logger         *customLogger.Logger
}

func NewTokenServer(tokenManager managers.TokenManager, refreshStorage database.RefreshStorage, logger *customLogger.Logger) *TokenServer {
	return &TokenServer{
		tokenManager:   tokenManager,
		refreshStorage: refreshStorage,
		logger:         logger,
	}
}

type refreshInput struct {
	UserID       string
	RefreshToken string
}

// @Summary Refresh tokens
// @Description get access and refresh tokens via user_id
// @Tags tokens
// @Param token body refreshInput true "Данные для регистрации пользователя"
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200
// @Failure 500
// @Router /token/refresh [post]
func (ts *TokenServer) RefreshTokens(ctx *gin.Context) {
	// Get token from request body for refreshing
	var inputToken refreshInput
	if err := ctx.BindJSON(&inputToken); err != nil {
		ts.logger.ErrorLogger.Error().Msgf("invalid input for refreshing tokens: %s", err)
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Compare it with already saved token
	if err := ts.tokenManager.CompareRefreshTokens(ts.refreshStorage, inputToken.RefreshToken, inputToken.UserID); err != nil {
		ts.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	newAccessToken, newRefreshToken, err := ts.tokenManager.GenerateBothTokens(inputToken.UserID)
	if err != nil {
		ts.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = ts.refreshStorage.Delete(inputToken.UserID); err != nil {
		ts.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	hashedRefreshToken, err := ts.tokenManager.GetHashedRefreshToken(newRefreshToken)
	if err != nil {
		ts.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	expTime := time.Now().Add(time.Minute * 43200).UTC() // 30 days refresh_token is valid
	refreshToken := entity.RefreshToken{
		UserID:       inputToken.UserID,
		RefreshToken: hashedRefreshToken,
		ExpiresAt:    expTime,
	}
	if err = ts.refreshStorage.Post(refreshToken); err != nil {
		ts.logger.ErrorLogger.Error().Msg(err.Error())
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
