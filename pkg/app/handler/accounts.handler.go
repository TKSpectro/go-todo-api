package handler

import (
	"errors"
	"strconv"

	"github.com/TKSpectro/go-todo-api/pkg/app/model"
	"github.com/TKSpectro/go-todo-api/pkg/app/types"
	"github.com/TKSpectro/go-todo-api/pkg/middleware/locals"
	"github.com/TKSpectro/go-todo-api/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetAccounts   godoc
// @Summary      List accounts
// @Tags         accounts
// @Accept       json
// @Param		 meta query pagination.QueryParams false "Pagination Query Parameters"
// @Produce      json
// @Success      200  {array}  model.Account
// @Router       /accounts [get]
func (h *Handler) GetAccounts(c *fiber.Ctx) error {
	var accounts = &[]model.Account{}
	var meta = locals.Meta(c)

	err := h.accountService.FindAccounts(accounts, meta).Error
	if err != nil {
		return &utils.INTERNAL_SERVER_ERROR
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
// @Success      200  {object}  model.Account
// @Router       /accounts/{id} [get]
func (h *Handler) GetAccount(c *fiber.Ctx) error {
	var account = &model.Account{}

	remoteIdString := c.Params("id")
	if remoteIdString == "" {
		return &utils.BAD_REQUEST
	}

	remoteId, err := strconv.ParseUint(remoteIdString, 10, 32)
	if err != nil {
		return &utils.BAD_REQUEST
	}

	err = h.accountService.FindAccountByID(account, uint(remoteId)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utils.NOT_FOUND
		}
		return &utils.INTERNAL_SERVER_ERROR
	}

	return c.JSON(&types.GetAccountResponse{
		Account: *account,
	})
}
