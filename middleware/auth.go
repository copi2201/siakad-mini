package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"siakad-mini/helper"
)

func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")

		if authHeader == "" {
			return helper.Error(c, fiber.StatusUnauthorized, "Token tidak ditemukan", nil)
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			return helper.Error(c, fiber.StatusUnauthorized, "Format token tidak valid", nil)
		}

		claims, err := helper.ParseToken(parts[1])
		if err != nil {
			return helper.Error(c, fiber.StatusUnauthorized, "Token tidak valid atau kedaluwarsa", nil)
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}