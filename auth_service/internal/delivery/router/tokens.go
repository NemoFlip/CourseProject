package router

import (
	"CourseProject/auth_service/internal/delivery/handlers"
	"CourseProject/auth_service/internal/delivery/middleware"
	customLogger "CourseProject/auth_service/pkg/log"
	"CourseProject/auth_service/pkg/managers"
	"github.com/gin-gonic/gin"
)

func RegisterTokenRouters(r *gin.Engine, logger *customLogger.Logger, tokenManager *managers.TokenManager, tokenServer *handlers.TokenServer) {
	checkAuth := middleware.CheckAuthorization(tokenManager, logger)
	secureGroup := r.Group("/", checkAuth)
	{
		secureGroup.POST("/token/refresh", tokenServer.RefreshTokens)
	}
}
