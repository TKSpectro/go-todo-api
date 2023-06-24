package types

import (
	"tkspectro/vefeast/app/models"
	"tkspectro/vefeast/app/types/pagination"
)

type GetTodosResponse struct {
	Todos []models.Todo   `json:"todos"`
	Meta  pagination.Meta `json:"_meta"`
}

type GetTodoResponse struct {
	Todo models.Todo `json:"todo"`
}
