package routes

import (
	"tkspectro/vefeast/app/handler"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesAuth(router fiber.Router) {
	router.Put("/login", handler.Login)
	router.Post("/register", handler.Register)
}
