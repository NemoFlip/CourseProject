package router

import (
	"CourseProject/auth_service/internal/delivery/handlers"
	"CourseProject/auth_service/internal/delivery/middleware"
	"CourseProject/auth_service/pkg/managers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRouters(r *gin.Engine, userServer *handlers.UserServer, tokenManager *managers.TokenManager, passRecoveryServer *handlers.PassRecoveryServer) {
	r.POST("/registration", userServer.RegisterUser)
	r.POST("/login", userServer.LoginUser)

	checkAuth := middleware.CheckAuthorization(tokenManager)
	secureGroup := r.Group("/", checkAuth)
	{
		secureGroup.GET("/refresh")
		secureGroup.POST("/logout", userServer.LogoutUser)
	}

	g2 := r.Group("/password")
	{
		g2.POST("/recovery", userServer.PasswordRecovery)
		g2.POST("/verify", passRecoveryServer.VerifyCode)
		g2.POST("/update", passRecoveryServer.UpdatePassword) // как разрешать /update только при успешном /verify?
	}
}
