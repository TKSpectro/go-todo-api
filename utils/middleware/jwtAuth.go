package middleware

import (
	"strings"

	"github.com/TKSpectro/go-todo-api/core"
	"github.com/TKSpectro/go-todo-api/utils/jwt"
	"github.com/TKSpectro/go-todo-api/utils/middleware/locals"

	"github.com/gofiber/fiber/v2"
)

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
		// TODO: Only send the error message in development mode
		return core.RequestErrorFrom(&core.UNAUTHORIZED, err.Error())
	}

	c.Locals(locals.KEY_PAYLOAD, payload)

	return c.Next()
}
