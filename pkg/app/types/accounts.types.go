package types

import (
	"github.com/TKSpectro/go-todo-api/pkg/app/model"
	"github.com/TKSpectro/go-todo-api/pkg/app/types/pagination"
)

type GetAccountsResponse struct {
	Accounts []model.Account `json:"accounts"`
	Meta     pagination.Meta `json:"_meta"`
}

type GetAccountResponse struct {
	Account model.Account `json:"account"`
}
