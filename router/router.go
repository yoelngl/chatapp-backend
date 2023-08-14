package router

import (
	"github.com/chatapp/backend/cmd"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1/")
	api.Get("/", cmd.Home)
	api.Get("/users", cmd.User)

	api.Post("/register", cmd.SignUp)
	api.Post("/login", cmd.Login)
}
