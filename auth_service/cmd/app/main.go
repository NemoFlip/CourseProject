package main

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/delivery"
	"CourseProject/auth_service/internal/delivery/handlers"
	customLogger "CourseProject/auth_service/pkg/log"
	"CourseProject/auth_service/pkg/managers"
	_ "CourseProject/docs"
	"CourseProject/pkg"
)

func main() {
	logger := customLogger.InitLogger()

	db, err := pkg.PostgresConnect("usersdb")
	if err != nil {
		logger.ErrorLogger.Fatal().Msg(err.Error())
	}
	logger.InfoLogger.Info().Msg("usersdb in initialized")

	tokenManager := managers.NewTokenManager(logger)
	emailManager := managers.NewEmailManager(logger)
	refreshStorage := database.NewRefreshStorage(logger)
	verifyCodeStorage := database.NewVerifyCodeStorage(logger)

	if tokenManager == nil || refreshStorage == nil {
		logger.ErrorLogger.Fatal().Msg("unable to connect to token_manager or to refresh_storage")
	}
	userStorage := database.NewUserStorage(db)
	userServer := handlers.NewUserServer(
		*userStorage,
		*tokenManager,
		*refreshStorage,
		emailManager,
		verifyCodeStorage,
		logger)
	passRecoveryServer := handlers.NewPassRecoveryServer(verifyCodeStorage, *userStorage, logger)

	delivery.StartServer(logger, userServer, tokenManager, passRecoveryServer)
}
