package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"core-backend/internal/config"
	"core-backend/internal/database"
	"core-backend/internal/handlers"
	"core-backend/pkg/logger"
)

func main() {
	logger.InitLogger()
	defer logger.Log.Sync()

	config.LoadConfig()

	database.Connect()
	defer database.Close()

	database.Migrate()

	app := fiber.New()

	app.Get("/keys/:id", handlers.GetPublicKey)

	logger.Log.Info("CoreGuard API 3000 portunda başlatılıyor...")
	log.Fatal(app.Listen(":3000"))
}
