package handler

import (
	"github.com/TKSpectro/go-todo-api/app/service"
)

type Handler struct {
	accountService service.IAccountService
	todoService    service.ITodoService
	validator      *Validator
}

func NewHandler(as service.IAccountService, ts service.ITodoService) *Handler {
	v := NewValidator()

	return &Handler{
		accountService: as,
		todoService:    ts,
		validator:      v,
	}
}
