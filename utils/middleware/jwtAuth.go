package middleware

import (
	"strings"

	"github.com/TKSpectro/go-todo-api/core"
	"github.com/TKSpectro/go-todo-api/utils/jwt"
	"github.com/TKSpectro/go-todo-api/utils/middleware/locals"

	"github.com/gofiber/fiber/v2"
)

// TODO: Switch to RS256 with signing files https://github.com/gofiber/contrib/tree/main/jwt#rs256-example
func Protected(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return &core.UNAUTHORIZED
	}

	chunks := strings.Split(authHeader, " ")

	if len(chunks) < 2 {
		return &core.UNAUTHORIZED
	}

	payload, err := jwt.Verify(chunks[1])
	if err != nil {
		return &core.UNAUTHORIZED
	}

	c.Locals(locals.KEY_PAYLOAD, payload)

	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
