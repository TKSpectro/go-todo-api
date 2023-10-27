package handler

import (
	"errors"
	"net/http"

	"github.com/TKSpectro/go-todo-api/pkg/app/model"
	"github.com/TKSpectro/go-todo-api/pkg/app/types"
	"github.com/TKSpectro/go-todo-api/pkg/app/types/pagination"
	"github.com/TKSpectro/go-todo-api/pkg/jwt"
	"github.com/TKSpectro/go-todo-api/pkg/middleware/locals"
	"github.com/TKSpectro/go-todo-api/pkg/view"
	"github.com/TKSpectro/go-todo-api/utils"
	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"gorm.io/gorm"
)

func (h *Handler) GetBaseData(c *fiber.Ctx) view.BaseData {
	account := &model.Account{}
	if locals.JwtPayload(c).Valid {
		h.accountService.FindAccountByID(account, locals.JwtPayload(c).AccountID)
	}

	return view.BaseData{
		IsAuthenticated: locals.JwtPayload(c).Valid,
		Account:         account,
	}
}

func (h *Handler) VIndex(c *fiber.Ctx) error {
	return adaptor.HTTPHandler(templ.Handler(view.IndexPage(h.GetBaseData(c))))(c)
}

func (h *Handler) VTodosIndex(c *fiber.Ctx) error {
	var meta = locals.Meta(c)

	meta.Order = append(meta.Order, pagination.OrderEntry{
		Key:       "created_at",
		Direction: "desc",
	})

	var todos = []model.Todo{}
	if err := h.FindWithMeta(&todos, &model.Todo{}, meta, h.db.Where("account_id = ?", locals.JwtPayload(c).AccountID)).Error; err != nil {
		return &utils.INTERNAL_SERVER_ERROR
	}

	return adaptor.HTTPHandler(templ.Handler(view.TodosIndexPage(h.GetBaseData(c), todos)))(c)
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

	return adaptor.HTTPHandler(templ.Handler(view.TodoItem(*todo)))(c)
}

func (h *Handler) VTodosComplete(c *fiber.Ctx) error {
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

	return adaptor.HTTPHandler(templ.Handler(view.TodoCompleteToggle(*todo)))(c)
}

func (h *Handler) VTodosDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	h.todoService.DeleteTodoByID(id)

	return c.Status(http.StatusOK).SendString("")
}

func (h *Handler) VLogin(c *fiber.Ctx) error {
	return adaptor.HTTPHandler(templ.Handler(view.LoginPage(h.GetBaseData(c))))(c)
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

func (h *Handler) VLogout(c *fiber.Ctx) error {
	c.ClearCookie("go-todo-api_auth")
	c.ClearCookie("go-todo-api_refresh")

	c.Response().Header.Set("HX-Redirect", "/")

	return c.Status(http.StatusOK).SendString("")
}
