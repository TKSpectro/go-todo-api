package handler

import (
	"errors"
	"net/http"

	"github.com/TKSpectro/go-todo-api/pkg/app/model"
	"github.com/TKSpectro/go-todo-api/pkg/app/types"
	"github.com/TKSpectro/go-todo-api/pkg/jwt"
	"github.com/TKSpectro/go-todo-api/pkg/middleware/locals"
	"github.com/TKSpectro/go-todo-api/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func (h *Handler) VIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello, World!",
	})
}

func (h *Handler) VTodosIndex(c *fiber.Ctx) error {
	var meta = locals.Meta(c)

	meta.Order = "created_at desc"

	var todos = &[]model.Todo{}
	if err := h.todoService.FindTodosByAccount(todos, meta, locals.JwtPayload(c).AccountID).Error; err != nil {
		return &utils.INTERNAL_SERVER_ERROR
	}

	return c.Render("todos", fiber.Map{
		"Title": "Todos",
		"Todos": todos,
	})
}

func (h *Handler) VTodosCreate(c *fiber.Ctx) error {
	todo := &model.Todo{}

	if err := ParseBodyAndValidate(c, todo, *h.validator); err != nil {
		return err
	}

	todo.AccountID = locals.JwtPayload(c).AccountID

	if err := h.todoService.CreateTodo(todo).Error; err != nil {
		return &utils.INTERNAL_SERVER_ERROR
	}

	return c.Render("todo-item", todo)
}

func (h *Handler) VTodosUpdate(c *fiber.Ctx) error {
	id := c.Params("id")

	todo := &model.Todo{}
	if err := h.todoService.FindTodoByID(todo, id, locals.JwtPayload(c).AccountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utils.NOT_FOUND
		}
		return &utils.INTERNAL_SERVER_ERROR
	}

	todo.Completed = !todo.Completed

	if err := h.todoService.UpdateTodo(todo).Error; err != nil {
		return &utils.INTERNAL_SERVER_ERROR
	}

	return c.Render("todo-complete-toggle", todo)
}

func (h *Handler) VTodosDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	h.todoService.DeleteTodoByID(id)

	return c.Status(http.StatusOK).SendString("")
}

func (h *Handler) VLogin(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title": "Login",
	})
}

func (h *Handler) VLoginPost(c *fiber.Ctx) error {
	remoteData := &types.LoginDTOBody{}

	if err := ParseBodyAndValidate(c, remoteData, *h.validator); err != nil {
		return err
	}

	account := &model.Account{}
	if err := h.accountService.FindAccountByEmail(account, remoteData.Email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utils.NOT_FOUND
		}
		return err
	}

	if !model.CheckPasswordHash(remoteData.Password, account.Password) {
		return &utils.AUTH_LOGIN_WRONG_PASSWORD
	}

	auth, err := jwt.Generate(account)
	if err != nil {
		return err
	}

	//Set cookie and return 200
	c.Cookie(&fiber.Cookie{
		Name:     "go-todo-api_auth",
		Value:    auth.Token,
		HTTPOnly: true,
	})

	c.Cookie(&fiber.Cookie{
		Name:     "go-todo-api_refresh",
		Value:    auth.RefreshToken,
		HTTPOnly: true,
	})

	c.Response().Header.Set("HX-Redirect", "/")

	return c.Status(http.StatusOK).SendString("")
}
