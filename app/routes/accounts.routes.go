package routes

import (
	"github.com/TKSpectro/go-todo-api/app/services"
	"github.com/TKSpectro/go-todo-api/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesAccounts(router fiber.Router) {
	router.Get("/", middleware.Protected, middleware.Pagination, services.GetAccounts)
	router.Get("/:id", middleware.Protected, services.GetAccount)
}
