package services

import (
	"errors"

	"github.com/TKSpectro/go-todo-api/app/models"
	"github.com/TKSpectro/go-todo-api/app/types"
	"github.com/TKSpectro/go-todo-api/core"
	"github.com/TKSpectro/go-todo-api/utils/middleware/locals"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetAccounts   godoc
// @Summary      List accounts
// @Tags         accounts
// @Accept       json
// @Param		 meta query pagination.QueryParams false "Pagination Query Parameters"
// @Produce      json
// @Success      200  {array}  models.Account
// @Router       /accounts [get]
func GetAccounts(c *fiber.Ctx) error {
	var accounts = &[]models.Account{}
	var meta = locals.Meta(c)

	err := models.FindAccounts(accounts, meta).Error
	if err != nil {
		return &core.INTERNAL_SERVER_ERROR
	}

	return c.JSON(&types.GetAccountsResponse{
		Accounts: *accounts,
		Meta:     *meta,
	})
}

// GetAccount    godoc
// @Summary      Get account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  models.Account
// @Router       /accounts/{id} [get]
func GetAccount(c *fiber.Ctx) error {
	var account = &models.Account{}

	remoteId := c.Params("id")
	if remoteId == "" {
		return &core.BAD_REQUEST
	}

	err := models.FindAccountByID(account, remoteId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &core.NOT_FOUND
		}
		return &core.INTERNAL_SERVER_ERROR
	}

	return c.JSON(&types.GetAccountResponse{
		Account: *account,
	})
}
