package utils

import (
	"github.com/TKSpectro/go-todo-api/core"

	"github.com/gofiber/fiber/v2"
)

func ParseBody(c *fiber.Ctx, body interface{}) *core.RequestError {
	if err := c.BodyParser(body); err != nil {
		return &core.BAD_REQUEST
	}

	return nil
}

func ParseBodyAndValidate(c *fiber.Ctx, body interface{}) *core.RequestError {
	if err := ParseBody(c, body); err != nil {
		return err
	}

	return Validate(body)
}

func GetCurrentAccountId(c *fiber.Ctx) *uint {
	id, _ := c.Locals("AccountId").(uint)
	return &id
}
