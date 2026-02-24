package handlers

import (
	"core-backend/internal/database"
	"core-backend/internal/models"

	"github.com/gofiber/fiber/v2"
)

func GetPublicKey(c *fiber.Ctx) error {
	targetID := c.Params("id")

	var userKey models.UserKey
	result := database.DB.Where("user_id = ?", targetID).First(&userKey)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Açık anahtar bulunamadı",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"public_key": userKey.PublicKey,
	})
}
