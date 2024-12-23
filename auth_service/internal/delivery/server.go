package delivery

import (
	"CourseProject/auth_service/internal/delivery/handlers"
	"CourseProject/auth_service/internal/delivery/router"
	"CourseProject/auth_service/pkg/auth"
	"CourseProject/auth_service/pkg/log"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Auth Service
// @description This is the auth services of course project
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /
func StartServer(logger *log.Logger, userServer *handlers.UserServer, tokenManager *auth.TokenManager, passRecoveryServer *handlers.PassRecoveryServer) {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.InitRouting(r, userServer, tokenManager, passRecoveryServer)

	if err := r.Run(":8080"); err != nil {
		logger.Fatal(fmt.Sprintf("unable to run server on port (:8080): %s", err))
	}
}
