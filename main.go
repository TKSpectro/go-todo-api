package main

import (
	"errors"
	"tkspectro/vefeast/app/models"
	"tkspectro/vefeast/app/routes"
	"tkspectro/vefeast/config/database"
	"tkspectro/vefeast/core"
	_ "tkspectro/vefeast/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// @title           fiber-api
// @version         1.0

// @BasePath  /api
func main() {
	database.Connect()
	database.Migrate(&models.Account{}, &models.Todo{})

	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Use(logger.New())

	routes.Setup(app)

	app.Listen(":3000")
}

var ErrorHandler = func(c *fiber.Ctx, err error) error {
	// Status statusCode defaults to 500
	statusCode := fiber.StatusInternalServerError
	code := 0
	message := ""

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
	}

	// Return status code with error message
	return c.Status(statusCode).JSON(fiber.Map{
		"statusCode": statusCode,
		"code":       code,
		"error":      message,
	})
}
