package main

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/transport/handlers"
	"CourseProject/auth_service/internal/transport/middleware"
	"CourseProject/auth_service/pkg/auth"
	customLogger "CourseProject/auth_service/pkg/log"
	_ "CourseProject/docs"
	"CourseProject/pkg"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title Auth Service
// @description This is the auth services of course project
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:8080
// @BasePath /
func main() {
	router := gin.Default()
	db, err := pkg.PostgresConnect("usersdb")
	if err != nil {
		log.Fatalf(err.Error())
	}

	logger := customLogger.InitLogger()

	logger.Logger.Info().Msg("Logger initialized")

	tokenManager := auth.NewTokenManager()
	emailManager := auth.NewEmailManager()
	refreshStorage := database.NewRefreshStorage()
	verifyCodeStorage := database.NewVerifyCodeStorage()

	if tokenManager == nil || refreshStorage == nil {
		log.Fatalf("unable to connect to token_manager or to refresh_storage")
	}
	userStorage := database.NewUserStorage(db)
	userServer := handlers.NewUserServer(*userStorage, *tokenManager, *refreshStorage, emailManager, verifyCodeStorage)
	passRecoveryServer := handlers.NewPassRecoveryServer(verifyCodeStorage, *userStorage)

	router.POST("/registration", userServer.RegisterUser)
	router.POST("/login", userServer.LoginUser)

	checkAuth := middleware.CheckAuthorization(tokenManager)
	secureGroup := router.Group("/", checkAuth)
	{
		secureGroup.GET("/refresh")
		secureGroup.POST("/logout", userServer.LogoutUser)
	}

	g2 := router.Group("/password")
	{
		g2.POST("/recovery", userServer.PasswordRecovery)
		g2.POST("/verify", passRecoveryServer.VerifyCode)
		g2.POST("/update", passRecoveryServer.UpdatePassword) // как разрешать /update только при успешном /verify?
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("unable to run server on port (:8080): %s", err)
	}
}
