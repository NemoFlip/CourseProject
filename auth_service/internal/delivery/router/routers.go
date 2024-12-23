package router

import (
	"CourseProject/auth_service/internal/delivery/handlers"
	"CourseProject/auth_service/pkg/auth"
	"github.com/gin-gonic/gin"
)

func InitRouting(r *gin.Engine, userServer *handlers.UserServer, tokenManager *auth.TokenManager, passRecoveryServer *handlers.PassRecoveryServer) {
	RegisterAuthRouters(r, userServer, tokenManager, passRecoveryServer)
}
