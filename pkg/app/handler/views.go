package handler

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/TKSpectro/go-todo-api/config"
	"github.com/TKSpectro/go-todo-api/pkg/app/model"
	"github.com/TKSpectro/go-todo-api/pkg/app/types"
	"github.com/TKSpectro/go-todo-api/pkg/app/types/pagination"
	"github.com/TKSpectro/go-todo-api/pkg/jwt"
	"github.com/TKSpectro/go-todo-api/pkg/middleware/locals"
	"github.com/TKSpectro/go-todo-api/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func defaultMap(c *fiber.Ctx, h *Handler, m *fiber.Map) fiber.Map {
	if m == nil {
		m = &fiber.Map{}
	}

	(*m)["IsAuthenticated"] = locals.JwtPayload(c).Valid
	(*m)["AccountID"] = locals.JwtPayload(c).AccountID

	if locals.JwtPayload(c).Valid {
		account := &model.Account{}
		h.accountService.FindAccountByID(account, locals.JwtPayload(c).AccountID)

		(*m)["Account"] = account
	}

	return *m
}

// renderPartial
// When rendering partials we need to use this instead of c.Render because it would include the whole template with the layout etc.
// As we use htmx this makes no sense and we only want to render the partial itself
//
// The templateName and name CAN be the same if the whole partial is in its own file
// but sometimes a partial will be the templateName ...-list and the name will be
// ...-item (...-list.html containing a block/define for ...-item.html)
func renderPartial(c *fiber.Ctx, templateName string, name string, data interface{}) error {
	tmpl := template.Must(template.ParseFiles(config.ROOT_PATH + "/pkg/view/partials/" + templateName + ".html"))
	err := tmpl.ExecuteTemplate(c.Response().BodyWriter(), name, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) VIndex(c *fiber.Ctx) error {
	return c.Render("index", defaultMap(c, h, nil))
}

func (h *Handler) VTodosIndex(c *fiber.Ctx) error {
	var meta = locals.Meta(c)

	meta.Order = append(meta.Order, pagination.OrderEntry{
		Key:       "created_at",
		Direction: "desc",
	})

	var todos = &[]model.Todo{}
	if err := h.FindWithMeta(todos, &model.Todo{}, meta, h.db.Where("account_id = ?", locals.JwtPayload(c).AccountID)).Error; err != nil {
		return &utils.INTERNAL_SERVER_ERROR
	}

	return c.Render("todos", defaultMap(c, h, &fiber.Map{
		"Todos": todos,
	}))
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

	if err := renderPartial(c, "todo-list", "todo-item", todo); err != nil {
		return &utils.INTERNAL_SERVER_ERROR
	}

	return nil
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

	if err := renderPartial(c, "todo-list", "todo-complete-toggle", todo); err != nil {
		return &utils.INTERNAL_SERVER_ERROR
	}

	return nil
}

func (h *Handler) VTodosDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	h.todoService.DeleteTodoByID(id)

	return c.Status(http.StatusOK).SendString("")
}

func (h *Handler) VLogin(c *fiber.Ctx) error {
	return c.Render("login", defaultMap(c, h, nil))
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
