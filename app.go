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
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Return status code with error message
	return c.Status(code).JSON(fiber.Map{

		"error":      "Internal Server Error",
		"statusCode": code,
	})
}
