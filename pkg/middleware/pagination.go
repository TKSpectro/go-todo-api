package middleware

import (
	"strings"

	"github.com/TKSpectro/go-todo-api/pkg/app/types/pagination"
	"github.com/TKSpectro/go-todo-api/pkg/middleware/locals"

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

	c.Locals(locals.KEY_META, &pagination.Meta{
		Page:   page,
		Limit:  limit,
		Skip:   (page - 1) * limit,
		Offset: (page - 1) * limit,

		Order:   parseOrders(c),
		Search:  parseSearch(c.Query("search")),
		Filters: parseFilters(c),
	})

	return c.Next()
}

func parseFilters(c *fiber.Ctx) []pagination.FilterEntry {
	queryArgs := c.Context().QueryArgs()
	filters := make([]pagination.FilterEntry, 0)

	queryArgs.VisitAll(func(key, value []byte) {
		if strings.HasPrefix(string(key[:6]), "filter") {
			// Split up the filter string by ] and already skip the first "["
			entries := strings.Split(string(key[7:]), "]")

			// Default
			operator := "eq"

			// If there are more than 2 entries, the operator is specified in the filter
			if len(entries) > 2 {
				// Split up the filter string by "[" so we can get the actual operator
				// The last "]" was already removed in the previous split
				operator = strings.Split(entries[1], "[")[1]
			}

			filters = append(filters, pagination.FilterEntry{
				Key:      entries[0],
				Operator: operator,
				Value:    string(value),
			})
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

// parseOrders parses the order string and returns a slice of OrderEntry
// The allowed format is ?order[id]=asc or ?order=id or ?order[id]=desc
// If the direction is not specified, asc is used
func parseOrders(c *fiber.Ctx) []pagination.OrderEntry {
	orders := []pagination.OrderEntry{}

	queryArgs := c.Context().QueryArgs()

	queryArgs.VisitAll(func(key, value []byte) {
		if strings.HasPrefix(string(key), "order") {
			entries := strings.Split(string(key[6:]), "]")

			direction := "asc"

			if value != nil && string(value) == "desc" {
				direction = string(value)
			}

			orders = append(orders, pagination.OrderEntry{
				Key:       entries[0],
				Direction: direction,
			})
		}
	})

	return orders
}
