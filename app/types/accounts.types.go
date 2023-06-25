package types

import (
	"github.com/TKSpectro/go-todo-api/app/models"
	"github.com/TKSpectro/go-todo-api/app/types/pagination"
)

type GetAccountsResponse struct {
	Accounts []models.Account `json:"accounts"`
	Meta     pagination.Meta  `json:"_meta"`
}

type GetAccountResponse struct {
	Account models.Account `json:"account"`
}
