package services

import (
	"errors"

	"github.com/TKSpectro/go-todo-api/app/models"
	"github.com/TKSpectro/go-todo-api/app/types"
	"github.com/TKSpectro/go-todo-api/core"
	"github.com/TKSpectro/go-todo-api/utils"
	"github.com/TKSpectro/go-todo-api/utils/middleware/locals"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetTodos   godoc
// @Summary    List todos
// @Tags       todos
// @Accept     json
// @Param	   meta query pagination.QueryParams false "Pagination Query Parameters"
// @Produce    json
// @Success    200  {object}  types.GetTodosResponse
// @Router     /todos [get]
func GetTodos(c *fiber.Ctx) error {
	var meta = locals.Meta(c)

	var todos = &[]models.Todo{}
	if err := models.FindTodosByAccount(todos, meta, locals.JwtPayload(c).AccountID).Error; err != nil {
		return &core.INTERNAL_SERVER_ERROR
	}

	return c.JSON(&types.GetTodosResponse{
		Todos: *todos,
		Meta:  *meta,
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

	if err := models.FindTodoByID(todo, remoteId, locals.JwtPayload(c).AccountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &core.NOT_FOUND
		}
		return &core.INTERNAL_SERVER_ERROR
	}

	return c.JSON(&types.GetTodoResponse{
		Todo: *todo,
	})
}

func CreateRandomTodo(c *fiber.Ctx) error {
	return models.CreateRandomTodo(locals.JwtPayload(c).AccountID).Error
}

// CreateTodo    godoc
// @Summary      Create todo
// @Description  Create todo
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        todo body 		types.CreateTodoRequest true "Todo"
// @Success      200  {object}  types.CreateTodoResponse
// @Router       /todos [post]
func CreateTodo(c *fiber.Ctx) error {
	remoteData := &types.CreateTodoRequest{}

	if err := utils.ParseBodyAndValidate(c, remoteData); err != nil {
		return err
	}

	var todo = &models.Todo{}
	todo.WriteRemote(remoteData.Todo)
	todo.AccountID = locals.JwtPayload(c).AccountID

	if err := models.CreateTodo(todo).Error; err != nil {
		return core.RequestErrorFrom(&core.INTERNAL_SERVER_ERROR, err.Error())
	}

	return c.JSON(&types.CreateTodoResponse{
		Todo: *todo,
	})
}

// UpdateTodo    godoc
// @Summary      Update todo
// @Description  Update todo
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Todo ID"
// @Param        todo body      types.UpdateTodoRequest true "Todo"
// @Success      200  {object}  types.UpdateTodoResponse
// @Router       /todos/{id} [put]
func UpdateTodo(c *fiber.Ctx) error {
	var todo = &models.Todo{}

	remoteId := c.Params("id")
	if remoteId == "" {
		return &core.BAD_REQUEST
	}

	var remoteData = &types.UpdateTodoRequest{}
	if err := utils.ParseBodyAndValidate(c, remoteData); err != nil {
		return err
	}

	if err := models.FindTodoByID(todo, remoteId, locals.JwtPayload(c).AccountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &core.NOT_FOUND
		}
		return &core.INTERNAL_SERVER_ERROR
	}

	todo.WriteRemote(remoteData.Todo)

	if err := models.UpdateTodo(todo).Error; err != nil {
		return core.RequestErrorFrom(&core.INTERNAL_SERVER_ERROR, err.Error())
	}

	return c.JSON(&types.UpdateTodoResponse{
		Todo: *todo,
	})
}

// DeleteTodo    godoc
// @Summary      Delete todo
// @Description  Delete todo
// @Tags         todos
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Todo ID"
// @Success      204  {object}  nil  "No Content"
// @Router       /todos/{id} [delete]
func DeleteTodo(c *fiber.Ctx) error {
	var todo = &models.Todo{}

	remoteId := c.Params("id")
	if remoteId == "" {
		return &core.BAD_REQUEST
	}

	if err := models.FindTodoByID(todo, remoteId, locals.JwtPayload(c).AccountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &core.NOT_FOUND
		}
		return &core.INTERNAL_SERVER_ERROR
	}

	if err := models.DeleteTodoByID(remoteId).Error; err != nil {
		return &core.INTERNAL_SERVER_ERROR
	}

	return c.SendStatus(fiber.StatusNoContent)
}
