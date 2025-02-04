package middleware

import (
	customLogger "CourseProject/auth_service/pkg/log"
	"CourseProject/auth_service/pkg/managers"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

func CheckAuthorization(tm *managers.TokenManager, logger *customLogger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		if len(authHeader) == 0 {
			logger.ErrorLogger.Error().Msg("authorization token is absent")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if bearerToken[0] != "Bearer" || len(bearerToken) != 2 {
			logger.ErrorLogger.Error().Msg("invalid format of auth token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]
		accessToken, err := tm.ValidateAccessToken(tokenString)
		if err != nil {
			logger.ErrorLogger.Error().Msg(err.Error())
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, ok := accessToken.Claims.(jwt.MapClaims)
		if !ok {
			logger.ErrorLogger.Error().Msg("invalid format of auth token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, ok := claims["sub"].(string)
		if !ok {
			logger.ErrorLogger.Error().Msg("unable to get `sub` claim from token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		exp, ok := claims["exp"].(float64)
		if !ok {
			logger.ErrorLogger.Error().Msg("unable to get `ip` claim from token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		currentTime := time.Now().Unix()
		if currentTime > int64(exp) {
			logger.ErrorLogger.Error().Msg("token is expired")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set("userID", userID)

		ctx.Next()
	}
}
