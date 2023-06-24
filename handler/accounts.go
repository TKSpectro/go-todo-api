package handler

import "github.com/gofiber/fiber/v2"

// GetAccounts   godoc
// @Summary      List accounts
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Success      200  {array}  model.Account
// @Router       /accounts [get]
func GetAccounts(c *fiber.Ctx) error {
	return c.SendString("/accounts")
}

// GetAccount    godoc
// @Summary      Get account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  model.Account
// @Router       /accounts/{id} [get]
func GetAccount(c *fiber.Ctx) error {
	return c.SendString("/accounts/:id")
}
