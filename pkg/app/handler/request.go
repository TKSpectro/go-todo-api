package handler

import (
	"github.com/TKSpectro/go-todo-api/utils"
	"github.com/gofiber/fiber/v2"
)

func ParseBody(c *fiber.Ctx, body interface{}) *utils.RequestError {
	if err := c.BodyParser(body); err != nil {
		return utils.RequestErrorFrom(&utils.BAD_REQUEST, err.Error())
	}

	return nil
}

func ParseBodyAndValidate(c *fiber.Ctx, body interface{}, v Validator) *utils.RequestError {
	if err := ParseBody(c, body); err != nil {
		return err
	}

	return v.Validate(body)
}
