package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chatapp/backend/database"
	"github.com/chatapp/backend/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New(fiber.Config{
		IdleTimeout: 5 * time.Second,
	})

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("APP_PORT")

	router.SetupRoutes(app)

	app.Use(cors.New())
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.ErrNotFound)
	})

	go func() {
		if err := app.Listen(":" + port); err != nil {
			log.Panic(err)
		}
	}()

	database.DBConnect()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	_ = <-c
	fmt.Println("Gracefully Shutting down...")
	_ = app.Shutdown()

	fmt.Println("Running cleaning tasks...")

	fmt.Println("Application was successfully shutdown!")
}
