package main

import (
	"errors"
	"tkspectro/vefeast/core"
	"tkspectro/vefeast/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: ErrorHandler,
	})

	app.Use(logger.New())

	router.SetupRoutes(app)

	core.SetupDatabase()

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
