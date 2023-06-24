package routes

import "github.com/gofiber/fiber/v2"

func RegisterRoutesTodos(router fiber.Router) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("/todos")
	})

	router.Get("/:id", func(c *fiber.Ctx) error {
		return c.SendString("/todos/:id")
	})
}
