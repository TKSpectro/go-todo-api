package types

import "github.com/TKSpectro/go-todo-api/app/models"

type GetAccountsResponse struct {
	Accounts []models.Account `json:"accounts"`
}

type GetAccountResponse struct {
	Account models.Account `json:"account"`
}
