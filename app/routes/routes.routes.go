package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func New(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from root")
	})

	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from base api path")
	})

	api.Get("/docs/*", swagger.HandlerDefault)

	RegisterRoutesAuth(api.Group("/auth"))

	RegisterRoutesAccounts(api.Group("/accounts"))
	RegisterRoutesTodos(api.Group("/todos"))
}
