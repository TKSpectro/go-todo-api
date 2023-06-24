package middleware

import (
	"tkspectro/vefeast/app/types/pagination"

	"github.com/gofiber/fiber/v2"
)

// PaginationMiddleware is a middleware that parses the query string for pagination parameters and sets them in the fiber.Ctx.Locals object
func Pagination(c *fiber.Ctx) error {
	page := c.QueryInt("page")
	if page <= 0 {
		page = 1
	}

	limit := c.QueryInt("limit")
	if limit <= 0 {
		limit = 10
	}

	c.Locals("meta", pagination.Meta{
		Page:   page,
		Limit:  limit,
		Skip:   (page - 1) * limit,
		Offset: (page - 1) * limit,
	})

	return c.Next()
}
