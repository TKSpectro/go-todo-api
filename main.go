package main

import (
	"errors"

	_ "github.com/TKSpectro/go-todo-api/api"
	"github.com/TKSpectro/go-todo-api/app/models"
	"github.com/TKSpectro/go-todo-api/app/routes"
	"github.com/TKSpectro/go-todo-api/config/database"
	"github.com/TKSpectro/go-todo-api/core"
	"github.com/TKSpectro/go-todo-api/utils"
	"github.com/TKSpectro/go-todo-api/utils/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// @title           fiber-api
// @version         1.0
// @BasePath  /api
func main() {
	app := Setup()

	app.Listen(":3000")
}

func Setup() *fiber.App {
	database.Connect()
	database.Migrate(&models.Account{}, &models.Todo{})

	jwt.Init()

	utils.RegisterCustomValidators()

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Use(logger.New())

	// Recover from panics anywhere in the chain and handle the control to the centralized ErrorHandler
	app.Use(recover.New())

	routes.Setup(app)

	return app
}

var ErrorHandler = func(c *fiber.Ctx, err error) error {
	// Status statusCode defaults to 500
	statusCode := fiber.StatusInternalServerError
	code := 0
	message := ""
	detail := ""

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		statusCode = e.Code
	}

	// Check if its our custom request error
	var requestError *core.RequestError
	if errors.As(err, &requestError) {
		statusCode = requestError.StatusCode
		code = requestError.Code
		message = requestError.Message
		detail = requestError.Detail
	}

	// Return status code with error message
	return c.Status(statusCode).JSON(fiber.Map{
		"statusCode": statusCode,
		"code":       code,
		"error":      message,
		"detail":     detail,
	})
}
