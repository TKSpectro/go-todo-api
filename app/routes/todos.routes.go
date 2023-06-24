package routes

import (
	"tkspectro/vefeast/app/services"
	"tkspectro/vefeast/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesTodos(router fiber.Router) {
	router.Get("/", middleware.Protected, middleware.Pagination, services.GetTodos)
	router.Get("/:id", middleware.Protected, services.GetTodo)
}
