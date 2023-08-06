package middleware

import (
	"strings"

	"github.com/TKSpectro/go-todo-api/pkg/jwt"
	"github.com/TKSpectro/go-todo-api/pkg/middleware/locals"
	"github.com/TKSpectro/go-todo-api/utils"

	"github.com/gofiber/fiber/v2"
)

func Protected(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return &utils.UNAUTHORIZED
	}

	chunks := strings.Split(authHeader, " ")

	if len(chunks) < 2 {
		return &utils.UNAUTHORIZED
	}

	payload, err := jwt.Verify(chunks[1])
	if err != nil {
		return utils.RequestErrorFrom(&utils.UNAUTHORIZED, err.Error())
	}

	c.Locals(locals.KEY_PAYLOAD, payload)

	return c.Next()
}
