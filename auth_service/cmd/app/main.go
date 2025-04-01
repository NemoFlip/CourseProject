package main

import (
	"CourseProject/auth_service/config"
	"CourseProject/auth_service/internal/app"
	customLogger "CourseProject/auth_service/pkg/log"
	_ "CourseProject/docs"
)

func main() {
	logger := customLogger.InitLogger()

	cfg := config.MustLoad()

	app.StartServer(logger, cfg)
}
