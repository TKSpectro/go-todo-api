package types

import (
	"github.com/TKSpectro/go-todo-api/app/models"
	"github.com/TKSpectro/go-todo-api/app/types/pagination"
)

type GetTodosResponse struct {
	Todos []models.Todo   `json:"todos"`
	Meta  pagination.Meta `json:"_meta"`
}

type GetTodoResponse struct {
	Todo models.Todo `json:"todo"`
}
