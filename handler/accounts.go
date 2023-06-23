package handler

import "github.com/gofiber/fiber/v2"

func GetAccounts(c *fiber.Ctx) error {
	return c.SendString("/accounts")
}

func GetAccount(c *fiber.Ctx) error {
	return c.SendString("/accounts/:id")
}
