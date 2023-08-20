package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/chatapp/backend/cmd"
	"github.com/chatapp/backend/database"
	"github.com/chatapp/backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	app := fiber.New()

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	database.DBConnect()

	app.Post("/v1/api/register", cmd.SignUp)

	tests := []struct {
		name           string
		method         string
		payload        *models.User
		expectedStatus int
		message        string
	}{
		{
			name:           "Valid Email",
			method:         "POST",
			payload:        &models.User{Email: "yoelngl321@gmail.com", Username: "yoelngl", Password: "12345678", FullName: "Yoel Gabriel Nainggolan"},
			expectedStatus: fiber.StatusCreated,
			message:        "Account successfully Registered!",
		},
		{
			name:           "Invalid Email",
			method:         "POST",
			payload:        &models.User{Email: "invalidemail.com", Username: "yoelngl", Password: "12345678", FullName: "Yoel Gabriel Nainggolan"},
			expectedStatus: fiber.StatusBadRequest,
			message:        "Please use an email format!",
		},
		{
			name:           "Email already Exists",
			method:         "POST",
			payload:        &models.User{Email: "yoelngl321@gmail.com", Username: "yoelngl", Password: "12345678", FullName: "Yoel Gabriel Nainggolan"},
			expectedStatus: fiber.StatusConflict,
			message:        "Email address already exists, please use another email!",
		},
		{
			name:           "Username already Exists",
			method:         "POST",
			payload:        &models.User{Email: "yoelmarket123@gmail.com", Username: "yoelngl", Password: "12345678", FullName: "Yoel Gabriel Nainggolan"},
			expectedStatus: fiber.StatusConflict,
			message:        "Username already exists, please use another username!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)

			req := httptest.NewRequest(tt.method, "/v1/api/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var result map[string]string
			json.NewDecoder(resp.Body).Decode(&result)
			assert.Equal(t, tt.message, result["message"])
		})
	}
}

func TestBodyInputSignUp(t *testing.T) {
	app := fiber.New()

	app.Post("/v1/api/register", cmd.SignUp)

	req := httptest.NewRequest("POST", "/v1/api/register", bytes.NewBufferString("invalid Body Request"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.Nil(t, err)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}
