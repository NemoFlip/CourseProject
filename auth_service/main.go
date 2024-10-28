package main

import (
	"CourseProject/auth_service/internal/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/register", handlers.RegisterUser)
	router.POST("/login", handlers.LoginUser)
}
