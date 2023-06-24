package types

import "tkspectro/vefeast/app/models"

type GetAccountsResponse struct {
	Accounts []models.Account `json:"accounts"`
}

type GetAccountResponse struct {
	Account models.Account `json:"account"`
}
