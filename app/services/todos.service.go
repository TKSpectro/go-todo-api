package services

import (
	"errors"

	"github.com/TKSpectro/go-todo-api/app/models"
	"github.com/TKSpectro/go-todo-api/app/types"
	"github.com/TKSpectro/go-todo-api/app/types/pagination"
	"github.com/TKSpectro/go-todo-api/core"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetTodos   godoc
// @Summary    List todos
// @Tags       todos
// @Accept     json
// @Produce    json
// @Success    200  {object}  types.GetTodosResponse
// @Router     /todos [get]
func GetTodos(c *fiber.Ctx) error {
	var meta = c.Locals("meta").(pagination.Meta)

	var todos = &[]models.Todo{}
	err := models.FindTodosByAccount(todos, &meta, c.Locals("AccountId").(uint)).Error
	if err != nil {
		return &core.INTERNAL_SERVER_ERROR
	}

	return c.JSON(&types.GetTodosResponse{
		Todos: *todos,
		Meta:  meta,
	})
}

// GetTodo    godoc
// @Summary    Get todo
// @Description  get string by ID
// @Tags       todos
// @Accept     json
// @Produce    json
// @Param      id   path      int  true  "Todo ID"
// @Success    200  {object}  types.GetTodoResponse
// @Router     /todos/{id} [get]
func GetTodo(c *fiber.Ctx) error {
	var todo = &models.Todo{}

	remoteId := c.Params("id")
	if remoteId == "" {
		return &core.BAD_REQUEST
	}

	err := models.FindTodoByID(todo, remoteId, c.Locals("AccountId").(uint)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &core.NOT_FOUND
		}
		return &core.INTERNAL_SERVER_ERROR
	}

	return c.JSON(&types.GetTodoResponse{
		Todo: *todo,
	})
}
