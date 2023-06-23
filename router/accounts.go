package router

import (
	"tkspectro/vefeast/handler"
	"tkspectro/vefeast/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutesAccounts(router fiber.Router) {
	router.Get("/", middleware.Protected(), handler.GetAccounts)
	router.Get("/:id", middleware.Protected(), handler.GetAccount)
}
