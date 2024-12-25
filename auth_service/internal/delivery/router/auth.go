package router

import (
	"CourseProject/auth_service/internal/delivery/handlers"
	"CourseProject/auth_service/internal/delivery/middleware"
	customLogger "CourseProject/auth_service/pkg/log"
	"CourseProject/auth_service/pkg/managers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRouters(r *gin.Engine, logger *customLogger.Logger, userServer *handlers.UserServer, tokenManager *managers.TokenManager, passRecoveryServer *handlers.PassRecoveryServer) {
	r.POST("/auth/register", userServer.RegisterUser)
	r.POST("/auth/login", userServer.LoginUser)

	checkAuth := middleware.CheckAuthorization(tokenManager, logger)
	secureGroup := r.Group("/", checkAuth)
	{
		secureGroup.GET("/token/refresh")
		secureGroup.POST("/auth/logout", userServer.LogoutUser)
	}
	g2 := r.Group("/password")
	{
		g2.POST("/request-reset", userServer.PasswordRecovery)
		g2.POST("/validate-code", passRecoveryServer.ValidateCode)
		g2.POST("/reset", passRecoveryServer.ResetPassword)
	}
}
