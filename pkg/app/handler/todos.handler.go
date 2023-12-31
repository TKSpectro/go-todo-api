package handler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/TKSpectro/go-todo-api/pkg/app/model"
	"github.com/TKSpectro/go-todo-api/pkg/app/types"
	"github.com/TKSpectro/go-todo-api/pkg/middleware/locals"
	"github.com/TKSpectro/go-todo-api/utils"

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
func (h *Handler) GetTodos(c *fiber.Ctx) error {
	var meta = locals.Meta(c)

	var todos = &[]model.Todo{}
	if err := h.FindWithMeta(todos, &model.Todo{}, meta, h.db.Where("account_id = ?", locals.JwtPayload(c).AccountID)).Error; err != nil {
		return &utils.INTERNAL_SERVER_ERROR
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
func (h *Handler) GetTodo(c *fiber.Ctx) error {
	var todo = &model.Todo{}

	remoteId := c.Params("id")
	if remoteId == "" {
		return &utils.BAD_REQUEST
	}

	if err := h.todoService.FindTodoByID(todo, remoteId, locals.JwtPayload(c).AccountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utils.NOT_FOUND
		}
		return &utils.INTERNAL_SERVER_ERROR
	}

	return c.JSON(&types.GetTodoResponse{
		Todo: *todo,
	})
}

func (h *Handler) CreateRandomTodo(c *fiber.Ctx) error {
	return h.todoService.CreateRandomTodo(locals.JwtPayload(c).AccountID).Error
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
func (h *Handler) CreateTodo(c *fiber.Ctx) error {
	remoteData := &types.CreateTodoRequest{}

	if err := ParseBodyAndValidate(c, remoteData, *h.validator); err != nil {
		return err
	}

	var todo = &model.Todo{}
	todo.New(remoteData.Todo)
	todo.AccountID = locals.JwtPayload(c).AccountID

	if err := h.todoService.CreateTodo(todo).Error; err != nil {
		return utils.RequestErrorFrom(&utils.INTERNAL_SERVER_ERROR, err.Error())
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
func (h *Handler) UpdateTodo(c *fiber.Ctx) error {
	var todo = &model.Todo{}

	remoteId := c.Params("id")
	if remoteId == "" {
		return &utils.BAD_REQUEST
	}

	var remoteData = &types.UpdateTodoRequest{}
	if err := ParseBodyAndValidate(c, remoteData, *h.validator); err != nil {
		return err
	}

	if err := h.todoService.FindTodoByID(todo, remoteId, locals.JwtPayload(c).AccountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utils.NOT_FOUND
		}
		return &utils.INTERNAL_SERVER_ERROR
	}

	todo.New(remoteData.Todo)

	if err := h.todoService.UpdateTodo(todo).Error; err != nil {
		return utils.RequestErrorFrom(&utils.INTERNAL_SERVER_ERROR, err.Error())
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
func (h *Handler) DeleteTodo(c *fiber.Ctx) error {
	var todo = &model.Todo{}

	remoteId := c.Params("id")
	if remoteId == "" {
		return &utils.BAD_REQUEST
	}

	if err := h.todoService.FindTodoByID(todo, remoteId, locals.JwtPayload(c).AccountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utils.NOT_FOUND
		}
		return &utils.INTERNAL_SERVER_ERROR
	}

	if err := h.todoService.DeleteTodoByID(remoteId).Error; err != nil {
		return &utils.INTERNAL_SERVER_ERROR
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) ExportCSVTodos(c *fiber.Ctx) error {
	var todos = []model.Todo{}

	if err := h.todoService.FindTodos(&todos, locals.JwtPayload(c).AccountID).Error; err != nil {
		return &utils.INTERNAL_SERVER_ERROR
	}

	filePath := path.Join(os.TempDir(), fmt.Sprintf("%d", time.Now().UnixMilli())+"-todos.csv")

	w, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer w.Close()

	cw := csv.NewWriter(w)

	header := make([]string, len(model.TodoColumns))
	for i, col := range model.TodoColumns {
		header[i] = col.Name
	}
	cw.Write(header)

	for _, todo := range todos {
		record, err := todo.MarshalRecord()
		if err != nil {
			return err
		}
		cw.Write(record)
	}
	cw.Flush()

	if err := cw.Error(); err != nil {
		return err
	}

	c.Response().Header.Set("Content-Type", "text/csv")
	c.Response().Header.Set("Content-Disposition", "attachment; filename=todos.csv")
	c.Download(filePath)

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

func (h *Handler) ImportCSVTodos(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	cr := csv.NewReader(src)
	cr.FieldsPerRecord = len(model.TodoColumns)

	cr.Read() // Skip header

	errors := []error{}

	for {
		record, err := cr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			errors = append(errors, err)
			continue
		}

		todo := &model.Todo{}
		if err := todo.UnmarshalRecord(record); err != nil {
			errors = append(errors, err)
			continue
		}

		todo.AccountID = locals.JwtPayload(c).AccountID

		if err := h.todoService.CreateTodo(todo).Error; err != nil {
			errors = append(errors, err)
			continue
		}
	}

	if len(errors) > 0 {
		return c.JSON(&types.ImportCSVTodosResponse{
			Errors: errors,
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
