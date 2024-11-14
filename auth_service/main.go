package main

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/handlers"
	_ "CourseProject/docs"
	"CourseProject/pkg"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
)

// @title Auth Service
// @description This is the auth service of course project
// @host localhost:8080
// @BasePath /
func main() {
	router := gin.Default()
	db, err := pkg.PostgresConnect("usersdb")
	if err != nil {
		log.Fatalf(err.Error())
	}
	userStorage := database.NewUserStorage(db)
	userServer := handlers.NewUserServer(*userStorage)

	router.POST("/register", userServer.RegisterHandler)
	router.POST("/login", userServer.LoginUser)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("unable to run server on port (:8080): %s", err)
	}
}
