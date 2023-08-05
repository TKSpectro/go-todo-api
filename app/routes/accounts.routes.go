package routes

import (
	"github.com/TKSpectro/go-todo-api/app/handlers"
	"github.com/TKSpectro/go-todo-api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesAccounts(router fiber.Router) {
	router.Get("/", middleware.Protected, middleware.Pagination, handlers.GetAccounts)
	router.Get("/:id", middleware.Protected, handlers.GetAccount)
}
