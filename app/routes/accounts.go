package routes

import (
	"tkspectro/vefeast/app/handler"
	"tkspectro/vefeast/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesAccounts(router fiber.Router) {
	router.Get("/", middleware.Protected, handler.GetAccounts)
	router.Get("/:id", middleware.Protected, handler.GetAccount)
}
