package middleware

import (
	"github.com/gofiber/fiber/v2"
	"siakad-mini/helper"
)

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return c.Next()
			}
		}

		return helper.Error(
			c,
			fiber.StatusForbidden,
			"Akses ditolak",
			nil,
		)
	}
}