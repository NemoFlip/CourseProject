package main

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/handlers"
	"CourseProject/pkg"
	"github.com/gin-gonic/gin"
	"log"
)

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

	err = router.Run(":8080")
	if err != nil {
		log.Fatalf("unable to run server on port (:8080): %s", err)
	}
}
