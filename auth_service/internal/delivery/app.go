package delivery

import (
	"CourseProject/auth_service/internal/delivery/handlers"
	"CourseProject/auth_service/internal/delivery/router"
	"CourseProject/auth_service/pkg/log"
	"CourseProject/auth_service/pkg/managers"
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
func StartServer(logger *log.Logger, userServer *handlers.UserServer, tokenManager *managers.TokenManager, passRecoveryServer *handlers.PassRecoveryServer, tokenServer *handlers.TokenServer) {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.InitRouting(r, logger, userServer, tokenManager, passRecoveryServer, tokenServer)

	if err := r.Run(":8080"); err != nil {
		logger.ErrorLogger.Error().Msgf("unable to run server on port (:8080): %s", err)
	}
}
