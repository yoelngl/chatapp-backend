package cmd

import (
	"github.com/chatapp/backend/database"
	"github.com/chatapp/backend/models"
	"github.com/gofiber/fiber/v2"
)

func User(c *fiber.Ctx) error {
	users := []models.User{}
	database.DB.Db.Find(&users)

	return c.Status(200).JSON(users)
}
