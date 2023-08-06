package handler

import (
	"github.com/TKSpectro/go-todo-api/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// New registers all routes for the application
func (h *Handler) RegisterRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from root")
	})

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
