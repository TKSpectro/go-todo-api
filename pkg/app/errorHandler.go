package app

import (
	"errors"

	_ "github.com/TKSpectro/go-todo-api/api"
	"github.com/TKSpectro/go-todo-api/core"

	"github.com/gofiber/fiber/v2"
)

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
