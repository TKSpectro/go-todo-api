package routes

import (
	"github.com/TKSpectro/go-todo-api/app/handlers"
	"github.com/TKSpectro/go-todo-api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesAuth(router fiber.Router) {
	router.Put("/login", handlers.Login)
	router.Post("/register", handlers.Register)
	router.Put("/refresh", middleware.Protected, handlers.Refresh)
	router.Get("/me", middleware.Protected, handlers.Me)

	router.Put("/jwk-rotate", middleware.AllowedIps, handlers.RotateJWK)
}
