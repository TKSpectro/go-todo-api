package routes

import (
	"github.com/TKSpectro/go-todo-api/app/services"
	"github.com/TKSpectro/go-todo-api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesTodos(router fiber.Router) {
	router.Get("/", middleware.Protected, middleware.Pagination, services.GetTodos)
	router.Get("/:id", middleware.Protected, services.GetTodo)
	router.Post("/", middleware.Protected, services.CreateTodo)
	router.Put("/:id", middleware.Protected, services.UpdateTodo)
	router.Delete("/:id", middleware.Protected, services.DeleteTodo)

	router.Post("/random", middleware.Protected, services.CreateRandomTodo)
}
