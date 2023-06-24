package routes

import (
	"tkspectro/vefeast/app/services"
	"tkspectro/vefeast/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesAccounts(router fiber.Router) {
	router.Get("/", middleware.Protected, services.GetAccounts)
	router.Get("/:id", middleware.Protected, services.GetAccount)
}
