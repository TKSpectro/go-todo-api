package routes

import (
	"tkspectro/vefeast/app/services"
	"tkspectro/vefeast/utils/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesAuth(router fiber.Router) {
	router.Put("/login", services.Login)
	router.Post("/register", services.Register)
	router.Put("/refresh", middleware.Protected, services.Refresh)
}