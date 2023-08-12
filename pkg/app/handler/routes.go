package handler

import (
	"errors"
	"net/http"

	"github.com/TKSpectro/go-todo-api/pkg/app/model"
	"github.com/TKSpectro/go-todo-api/pkg/app/types"
	"github.com/TKSpectro/go-todo-api/pkg/jwt"
	"github.com/TKSpectro/go-todo-api/pkg/middleware"
	"github.com/TKSpectro/go-todo-api/pkg/middleware/locals"
	"github.com/TKSpectro/go-todo-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"gorm.io/gorm"
)

// New registers all routes for the application
func (h *Handler) RegisterRoutes(app *fiber.App) {
	h.RegisterViewRoutes(app)
	h.RegisterApiRoutes(app)
}

// RegisterViewRoutes registers all routes for the application that are used for rendering views
func (h *Handler) RegisterViewRoutes(app *fiber.App) {
	// TODO: Cleanup the RegisterViewRoutes into a separate file and only register the API routes here
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Get("/todos", middleware.Pagination, middleware.Protected, func(c *fiber.Ctx) error {
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
	})

	app.Post("/todos", middleware.Protected, func(c *fiber.Ctx) error {
		todo := &model.Todo{}

		if err := ParseBodyAndValidate(c, todo, *h.validator); err != nil {
			return err
		}

		todo.AccountID = locals.JwtPayload(c).AccountID

		if err := h.todoService.CreateTodo(todo).Error; err != nil {
			return &utils.INTERNAL_SERVER_ERROR
		}

		return c.Render("todo-item", todo)
	})

	app.Put("/todos/:id/complete", middleware.Protected, func(c *fiber.Ctx) error {
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
	})

	app.Delete("/todos/:id", middleware.Protected, func(c *fiber.Ctx) error {
		id := c.Params("id")
		h.todoService.DeleteTodoByID(id)

		return c.Status(http.StatusOK).SendString("")
	})

	app.Get("/login", func(c *fiber.Ctx) error {
		return c.Render("login", fiber.Map{
			"Title": "Login",
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
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
	})
}

func (h *Handler) RegisterApiRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from base api path")
	})

	api.Get("/docs/*", swagger.HandlerDefault)

	auth := api.Group("/auth")
	auth.Put("/login", h.Login)
	auth.Post("/register", h.Register)
	auth.Put("/refresh", middleware.Protected, h.Refresh)
	auth.Get("/me", middleware.Protected, h.Me)

	auth.Put("/jwk-rotate", middleware.AllowedIps, h.RotateJWK)

	accounts := api.Group("/accounts")
	accounts.Get("/", middleware.Protected, middleware.Pagination, h.GetAccounts)
	accounts.Get("/:id", middleware.Protected, h.GetAccount)

	todos := api.Group("/todos")
	todos.Get("/", middleware.Protected, middleware.Pagination, h.GetTodos)
	todos.Get("/:id", middleware.Protected, h.GetTodo)
	todos.Post("/", middleware.Protected, h.CreateTodo)
	todos.Put("/:id", middleware.Protected, h.UpdateTodo)
	todos.Delete("/:id", middleware.Protected, h.DeleteTodo)

	todos.Post("/random", middleware.Protected, h.CreateRandomTodo)
}
