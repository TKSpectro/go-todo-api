package handler

import (
	"github.com/TKSpectro/go-todo-api/core"

	"github.com/gofiber/fiber/v2"
)

func ParseBody(c *fiber.Ctx, body interface{}) *core.RequestError {
	if err := c.BodyParser(body); err != nil {
		return core.RequestErrorFrom(&core.BAD_REQUEST, err.Error())
	}

	return nil
}

func ParseBodyAndValidate(c *fiber.Ctx, body interface{}, v Validator) *core.RequestError {
	if err := ParseBody(c, body); err != nil {
		return err
	}

	return v.Validate(body)
}
