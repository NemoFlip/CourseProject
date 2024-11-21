package main

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/handlers"
	"CourseProject/auth_service/pkg/auth"
	_ "CourseProject/docs"
	"CourseProject/pkg"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
)

// @title Auth Service
// @description This is the auth services of course project
// @host localhost:8080
// @BasePath /
func main() {
	router := gin.Default()
	db, err := pkg.PostgresConnect("usersdb")
	if err != nil {
		log.Fatalf(err.Error())
	}

	jwtSecret, exists := os.LookupEnv("JWT_SECRET_KEY")
	if !exists {
		log.Printf("unable to parse jwt secret key from environment")
		return
	}
	tokenManager := auth.NewTokenManager(jwtSecret)
	userStorage := database.NewUserStorage(db)
	userServer := handlers.NewUserServer(*userStorage, *tokenManager)

	router.POST("/registration", userServer.RegisterUser)
	router.POST("/login", userServer.LoginUser)
	router.POST("/logout", userServer.LogoutUser)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("unable to run server on port (:8080): %s", err)
	}
}
