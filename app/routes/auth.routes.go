package routes

import (
	"github.com/TKSpectro/go-todo-api/app/services"
	"github.com/TKSpectro/go-todo-api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesAuth(router fiber.Router) {
	router.Put("/login", services.Login)
	router.Post("/register", services.Register)
	router.Put("/refresh", middleware.Protected, services.Refresh)
	router.Get("/me", middleware.Protected, services.Me)

	router.Put("/jwk-rotate", middleware.AllowedIps, services.RotateJWK)
}
