package app

import (
	"CourseProject/auth_service/config"
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/delivery/handlers"
	"CourseProject/auth_service/internal/delivery/router"
	"CourseProject/auth_service/pkg/db"
	"CourseProject/auth_service/pkg/log"
	"CourseProject/auth_service/pkg/managers"
	"database/sql"
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
func StartServer(logger *log.Logger, cfg *config.Config) {
	userServer, passRecoveryServer, tokenServer, tokenManager := initDependencies(logger, cfg)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.InitRouting(r, logger, userServer, tokenManager, passRecoveryServer, tokenServer)

	address := fmt.Sprintf(":%s", cfg.HTTPServer.Port)
	if err := r.Run(address); err != nil {
		logger.ErrorLogger.Error().Msgf("unable to run server on address (%s): %s", address, err)
	}
}

func initDependencies(logger *log.Logger, cfg *config.Config) (*handlers.UserServer, *handlers.PassRecoveryServer, *handlers.TokenServer, *managers.TokenManager) {
	// Connect to Database
	db := mustConnectDB(logger, cfg)

	// Create deps
	tokenManager := managers.NewTokenManager(logger)
	emailManager := managers.NewEmailManager(logger)
	refreshStorage := database.NewRefreshStorage(logger)
	verifyCodeStorage := database.NewVerifyCodeStorage(logger)
	userStorage := database.NewUserStorage(db)

	// Check deps
	mustCheckDependencies(logger, tokenManager, refreshStorage)

	return handlers.NewUserServer(*userStorage, *tokenManager, *refreshStorage, logger),
		handlers.NewPassRecoveryServer(verifyCodeStorage, *userStorage, logger, emailManager),
		handlers.NewTokenServer(*tokenManager, *refreshStorage, logger),
		tokenManager
}

// mustConnectDB connects to database and fails in other case
func mustConnectDB(logger *log.Logger, cfg *config.Config) *sql.DB {
	db, err := db.PostgresConnect(cfg)
	if err != nil {
		logger.ErrorLogger.Fatal().Msgf("Failed to connect to database: %s", err)
	}
	logger.InfoLogger.Info().Msgf("Connected to database: %s", cfg.UsersDB.Name)
	return db
}

func mustCheckDependencies(logger *log.Logger, tokenManager *managers.TokenManager, refreshStorage *database.RefreshStorage) {
	if tokenManager == nil || refreshStorage == nil {
		logger.ErrorLogger.Fatal().Msg("TokenManager or RefreshStorage is nil")
	}
}
