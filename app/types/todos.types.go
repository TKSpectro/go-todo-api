package types

import (
	"github.com/TKSpectro/go-todo-api/app/model"
	"github.com/TKSpectro/go-todo-api/app/types/pagination"
)

type GetTodosResponse struct {
	Todos []model.Todo    `json:"todos"`
	Meta  pagination.Meta `json:"_meta"`
}

type GetTodoResponse struct {
	Todo model.Todo `json:"todo"`
}

type CreateTodoRequest struct {
	Todo model.Todo `json:"todo"`
}

type CreateTodoResponse struct {
	Todo model.Todo `json:"todo"`
}

type UpdateTodoRequest struct {
	Todo model.Todo `json:"todo"`
}

type UpdateTodoResponse struct {
	Todo model.Todo `json:"todo"`
}
