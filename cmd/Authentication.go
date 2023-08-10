package cmd

import (
	"errors"

	"github.com/chatapp/backend/database"
	"github.com/chatapp/backend/models"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	if c.Method() != "POST" {
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.ErrMethodNotAllowed)
	}

	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": err.Error(),
		})
	}

	passwordHashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(passwordHashed)

	if checkExists := emailExistsAndUsername(c, user.Email, user.Username); checkExists != nil {
		return c.Status(fiber.StatusConflict).JSON(&fiber.Map{
			"message": checkExists.Error(),
		})
	}

	result := database.DB.Db.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Failed to create user!, " + result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"message": "Account successfully Registered!",
	})
}

func emailExistsAndUsername(c *fiber.Ctx, email string, username string) error {
	var user models.User

	emailUniqueness := database.DB.Db.Where("email = ?", email).First(&user)
	if emailUniqueness.RowsAffected > 0 {
		return errors.New("Email address already exists, please use another email!")
	}

	usernameUniqueness := database.DB.Db.Where("username = ?", username).First(&user)
	if usernameUniqueness.RowsAffected > 0 {
		return errors.New("Username already exists, please use another username!")
	}

	return nil
}
