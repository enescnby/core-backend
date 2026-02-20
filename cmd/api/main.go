package main

import (
	"core-backend/internal/config"
	"core-backend/internal/database"
	"core-backend/pkg/logger"
)

func main() {
	logger.InitLogger()
	defer logger.Log.Sync()

	config.LoadConfig()

	database.Connect()
	defer database.Close()

	database.Migrate()
}
