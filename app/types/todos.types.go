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

type CreateTodoRequest struct {
	Todo models.Todo `json:"todo"`
}

type CreateTodoResponse struct {
	Todo models.Todo `json:"todo"`
}

type UpdateTodoRequest struct {
	Todo models.Todo `json:"todo"`
}

type UpdateTodoResponse struct {
	Todo models.Todo `json:"todo"`
}
