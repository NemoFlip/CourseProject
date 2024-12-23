package router

import (
	"CourseProject/auth_service/internal/delivery/handlers"
	customLogger "CourseProject/auth_service/pkg/log"
	"CourseProject/auth_service/pkg/managers"
	"github.com/gin-gonic/gin"
)

func InitRouting(r *gin.Engine, logger *customLogger.Logger, userServer *handlers.UserServer, tokenManager *managers.TokenManager, passRecoveryServer *handlers.PassRecoveryServer) {
	RegisterAuthRouters(r, logger, userServer, tokenManager, passRecoveryServer)
}
