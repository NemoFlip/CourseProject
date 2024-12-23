package main

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/delivery"
	"CourseProject/auth_service/internal/delivery/handlers"
	"CourseProject/auth_service/pkg/auth"
	customLogger "CourseProject/auth_service/pkg/log"
	_ "CourseProject/docs"
	"CourseProject/pkg"
)

func main() {
	logger := customLogger.InitLogger()

	db, err := pkg.PostgresConnect("usersdb")
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("usersdb in initialized")

	tokenManager := auth.NewTokenManager()
	emailManager := auth.NewEmailManager()
	refreshStorage := database.NewRefreshStorage()
	verifyCodeStorage := database.NewVerifyCodeStorage()

	if tokenManager == nil || refreshStorage == nil {
		logger.Fatal("unable to connect to token_manager or to refresh_storage")
	}
	userStorage := database.NewUserStorage(db)
	userServer := handlers.NewUserServer(*userStorage, *tokenManager, *refreshStorage, emailManager, verifyCodeStorage)
	passRecoveryServer := handlers.NewPassRecoveryServer(verifyCodeStorage, *userStorage)

	delivery.StartServer(logger, userServer, tokenManager, passRecoveryServer)

}
