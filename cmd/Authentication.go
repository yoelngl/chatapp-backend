package cmd

import (
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/chatapp/backend/config"
	"github.com/chatapp/backend/database"
	"github.com/chatapp/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

func emailExistsAndUsername(email string, username string) error {
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

func getUserData(email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Db.Where(&models.User{Email: email}).Find(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func CheckPasswordHash(password string, hash string) bool {
	passwordPlain := []byte(password)
	passwordHashed := []byte(hash)
	err := bcrypt.CompareHashAndPassword(passwordHashed, passwordPlain)
	fmt.Println(err)
	return err == nil
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

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

	if !validEmail(user.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please use an email format!",
		})
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"message": "Failed to hash password. " + err.Error(),
		})
	}

	user.Password = string(passwordHashed)

	if checkExists := emailExistsAndUsername(user.Email, user.Username); checkExists != nil {
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

func Login(c *fiber.Ctx) error {
	if c.Method() != "POST" {
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.ErrMethodNotAllowed)
	}

	input := new(Authentication)
	var userData UserData

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Error on login request!",
			"data":    err,
		})
	}

	email := input.Email
	password := input.Password
	userModel, _ := new(models.User), *new(error)

	if validEmail(email) {
		userModel, _ = getUserData(email)
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"message": "Please use an email format!",
		})
	}

	if userModel.ID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
			"message": "User not found!",
		})
	} else {
		userData = UserData{
			Username: userModel.Username,
			Email:    userModel.Email,
			FullName: userModel.FullName,
			Password: userModel.Password,
		}
	}

	if !CheckPasswordHash(password, userData.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid password", "data": nil})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["full_name"] = userData.FullName
	claims["username"] = userData.Username
	claims["email"] = userData.Email
	claims["expired"] = time.Now().Add(time.Hour * 8).Unix()

	t, err := token.SignedString([]byte(config.Config("JWT_SECRET_KEY")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.ErrInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully Login!",
		"token":   t,
	})
}
