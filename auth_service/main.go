package main

import (
	"CourseProject/auth_service/internal/database"
	"CourseProject/auth_service/internal/handlers"
	"CourseProject/pkg"
	"github.com/gin-gonic/gin"
	"log"
)

var dbname = "usersdb"

func main() {
	router := gin.Default()
	db, err := pkg.ConnectToDB(dbname)
	if err != nil {
		log.Fatalf(err.Error())
	}
	userStorage := database.NewUserStorage(db)
	userServer := handlers.NewUserServer(*userStorage)

	router.POST("/register", userServer.RegisterHandler)
	router.POST("/login", userServer.LoginUser)
}
