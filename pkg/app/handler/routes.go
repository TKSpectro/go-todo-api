package handler

import (
	"github.com/TKSpectro/go-todo-api/config"
	"github.com/TKSpectro/go-todo-api/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// New registers all routes for the application
func (h *Handler) RegisterRoutes(app *fiber.App) {
	h.RegisterHyperMediaRoutes(app)
	h.RegisterApiRoutes(app)
}

// RegisterHyperMediaRoutes registers all routes for the application that are used for rendering views
//
// Good Read for Hypermedia-Driven Applications: https://hypermedia.systems/json-data-apis/
func (h *Handler) RegisterHyperMediaRoutes(app *fiber.App) {
	app.Static("/js", config.ROOT_PATH+"/pkg/view/js")

	app.Get("/", middleware.LoadAuth, h.VIndex)

	app.Get("/login", h.VLogin)
	app.Post("/login", h.VLoginPost)
	app.Post("/logout", h.VLogout)

	app.Get("/todos", middleware.Pagination, middleware.Protected, h.VTodosIndex)
	app.Post("/todos", middleware.Protected, h.VTodosCreate)
	app.Put("/todos/:id/complete", middleware.Protected, h.VTodosUpdate)
	app.Delete("/todos/:id", middleware.Protected, h.VTodosDelete)
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
