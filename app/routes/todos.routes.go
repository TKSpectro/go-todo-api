package routes

import (
	"github.com/TKSpectro/go-todo-api/app/handler"
	"github.com/TKSpectro/go-todo-api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesTodos(router fiber.Router) {
	router.Get("/", middleware.Protected, middleware.Pagination, handler.GetTodos)
	router.Get("/:id", middleware.Protected, handler.GetTodo)
	router.Post("/", middleware.Protected, handler.CreateTodo)
	router.Put("/:id", middleware.Protected, handler.UpdateTodo)
	router.Delete("/:id", middleware.Protected, handler.DeleteTodo)

	router.Post("/random", middleware.Protected, handler.CreateRandomTodo)
}
