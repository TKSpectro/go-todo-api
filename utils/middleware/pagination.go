package middleware

import (
	"strings"

	"github.com/TKSpectro/go-todo-api/app/types/pagination"

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

	// TODO: Sorting needs to be completed. This is just a placeholder/simple implementation
	order := c.Query("order")
	if order == "" {
		order = "id asc"
	}

	c.Locals("meta", pagination.Meta{
		Page:   page,
		Limit:  limit,
		Skip:   (page - 1) * limit,
		Offset: (page - 1) * limit,

		Order:   order,
		Search:  parseSearch(c.Query("search")),
		Filters: parseFilters(c),
	})

	return c.Next()
}

func parseFilters(c *fiber.Ctx) map[string]string {
	queryArgs := c.Context().QueryArgs()

	filters := make(map[string]string)
	// for queryArgs.Len
	queryArgs.VisitAll(func(key, value []byte) {
		// if key starts with filter
		if strings.HasPrefix(string(key[:6]), "filter") {
			modelKey := strings.Split(string(key[7:]), "]")[0]
			filters[modelKey] = string(value)
		}
	})

	return filters
}

// ParseSearch parses the search string and returns a string that can be used in a gorm query
func parseSearch(search string) string {
	if search == "" {
		return "%"
	}

	return "%" + search + "%"
}
