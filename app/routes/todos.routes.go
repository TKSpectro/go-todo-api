package routes

import (
	"github.com/TKSpectro/go-todo-api/app/handlers"
	"github.com/TKSpectro/go-todo-api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesTodos(router fiber.Router) {
	router.Get("/", middleware.Protected, middleware.Pagination, handlers.GetTodos)
	router.Get("/:id", middleware.Protected, handlers.GetTodo)
	router.Post("/", middleware.Protected, handlers.CreateTodo)
	router.Put("/:id", middleware.Protected, handlers.UpdateTodo)
	router.Delete("/:id", middleware.Protected, handlers.DeleteTodo)

	router.Post("/random", middleware.Protected, handlers.CreateRandomTodo)
}
