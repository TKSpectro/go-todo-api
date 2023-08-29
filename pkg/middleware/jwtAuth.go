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
	authCookie := c.Cookies("go-todo-api_auth")

	token := ""

	if authHeader != "" {
		chunks := strings.Split(authHeader, " ")

		if len(chunks) < 2 {
			return &utils.UNAUTHORIZED
		}

		token = chunks[1]
	} else if authCookie != "" {
		token = authCookie
	} else {
		return &utils.UNAUTHORIZED
	}

	payload, err := jwt.Verify(token)
	if err != nil {
		return utils.RequestErrorFrom(&utils.UNAUTHORIZED, err.Error())
	}

	c.Locals(locals.KEY_PAYLOAD, payload)

	return c.Next()
}

func LoadAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	authCookie := c.Cookies("go-todo-api_auth")

	token := ""

	if authHeader != "" {
		chunks := strings.Split(authHeader, " ")

		if len(chunks) < 2 {
			return &utils.UNAUTHORIZED
		}

		token = chunks[1]
	} else if authCookie != "" {
		token = authCookie
	} else {
		return c.Next()
	}

	payload, err := jwt.Verify(token)

	if err == nil {
		c.Locals(locals.KEY_PAYLOAD, payload)
	}

	return c.Next()
}
