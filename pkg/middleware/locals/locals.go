package locals

import (
	"github.com/TKSpectro/go-todo-api/pkg/app/types/pagination"
	"github.com/TKSpectro/go-todo-api/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

const (
	KEY_PAYLOAD = "LOCALS_PAYLOAD"
	KEY_META    = "LOCALS_META"
)

func JwtPayload(c *fiber.Ctx) (payload *jwt.TokenPayload) {
	if c.Locals(KEY_PAYLOAD) == nil {
		return &jwt.TokenPayload{
			Valid: false,
		}
	}

	return c.Locals(KEY_PAYLOAD).(*jwt.TokenPayload)
}

func Meta(c *fiber.Ctx) *pagination.Meta {
	return c.Locals(KEY_META).(*pagination.Meta)
}

func Can(c *fiber.Ctx, needs uint64) bool {
	return (JwtPayload(c).Permission & needs) > 0
}
