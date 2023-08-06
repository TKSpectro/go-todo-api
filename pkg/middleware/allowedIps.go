package middleware

import (
	"fmt"

	"github.com/TKSpectro/go-todo-api/config"
	"github.com/TKSpectro/go-todo-api/utils"
	"github.com/gofiber/fiber/v2"
)

// AllowedIps is a middleware that checks if the request comes from an allowed IP
func AllowedIps(c *fiber.Ctx) error {
	fmt.Println("AllowedIps middleware", c.IP())
	fmt.Println("AllowedIps middleware", config.ALLOWED_IPS)

	ip := c.IP()

	if !isAllowed(ip) {
		return &utils.FORBIDDEN
	}

	return c.Next()

}

func isAllowed(ip string) bool {
	for _, allowedIp := range config.ALLOWED_IPS {
		if ip == allowedIp {
			return true
		}

		// "*" means all ips are allowed
		if allowedIp == "*" {
			return true
		}
	}

	return false
}
