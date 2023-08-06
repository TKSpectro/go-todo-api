package handler

import (
	"github.com/TKSpectro/go-todo-api/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// TODO: Rename New to Register (Need to rewrite handler building first)

// New registers all routes for the application
func New(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from root")
	})

	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from base api path")
	})

	api.Get("/docs/*", swagger.HandlerDefault)

	auth := api.Group("/auth")
	auth.Put("/login", Login)
	auth.Post("/register", Register)
	auth.Put("/refresh", middleware.Protected, Refresh)
	auth.Get("/me", middleware.Protected, Me)

	auth.Put("/jwk-rotate", middleware.AllowedIps, RotateJWK)

	accounts := api.Group("/accounts")
	accounts.Get("/", middleware.Protected, middleware.Pagination, GetAccounts)
	accounts.Get("/:id", middleware.Protected, GetAccount)

	todos := api.Group("/todos")
	todos.Get("/", middleware.Protected, middleware.Pagination, GetTodos)
	todos.Get("/:id", middleware.Protected, GetTodo)
	todos.Post("/", middleware.Protected, CreateTodo)
	todos.Put("/:id", middleware.Protected, UpdateTodo)
	todos.Delete("/:id", middleware.Protected, DeleteTodo)

	todos.Post("/random", middleware.Protected, CreateRandomTodo)
}
