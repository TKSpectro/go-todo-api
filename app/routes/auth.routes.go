package routes

import (
	"github.com/TKSpectro/go-todo-api/app/handler"
	"github.com/TKSpectro/go-todo-api/pkg/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesAuth(router fiber.Router) {
	router.Put("/login", handler.Login)
	router.Post("/register", handler.Register)
	router.Put("/refresh", middleware.Protected, handler.Refresh)
	router.Get("/me", middleware.Protected, handler.Me)

	router.Put("/jwk-rotate", middleware.AllowedIps, handler.RotateJWK)
}
