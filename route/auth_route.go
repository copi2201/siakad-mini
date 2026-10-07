package route

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	"siakad-mini/app/model"
	"siakad-mini/app/service"
	"siakad-mini/helper"
	"siakad-mini/middleware"
)

func SetupAuthRoutes(
	api fiber.Router,
	authService *service.AuthService,
) {
	auth := api.Group("/auth")

	// POST /api/v1/auth/login
	auth.Post(
		"/login",
		limiter.New(limiter.Config{
			Max:        5,
			Expiration: 1 * time.Minute,
			LimitReached: func(c *fiber.Ctx) error {
				return helper.Error(
					c,
					fiber.StatusTooManyRequests,
					"Terlalu banyak percobaan login. Silakan coba lagi nanti.",
					nil,
				)
			},
		}),
		func(c *fiber.Ctx) error {
			var req model.LoginRequest

			if err := c.BodyParser(&req); err != nil {
				return helper.Error(
					c,
					fiber.StatusUnprocessableEntity,
					"Format request tidak valid",
					nil,
				)
			}

			result, err := authService.Login(
				c.Context(),
				req,
			)
			if err != nil {
				return helper.Error(
					c,
					fiber.StatusUnauthorized,
					err.Error(),
					nil,
				)
			}

			return helper.Success(
				c,
				fiber.StatusOK,
				"Login berhasil",
				result,
			)
		},
	)

	// GET /api/v1/auth/me
	auth.Get(
		"/me",
		middleware.RequireAuth(),
		func(c *fiber.Ctx) error {
			userID := c.Locals("user_id").(int64)

			result, err := authService.GetMe(
				c.Context(),
				userID,
			)
			if err != nil {
				return helper.Error(
					c,
					fiber.StatusNotFound,
					err.Error(),
					nil,
				)
			}

			return helper.Success(
				c,
				fiber.StatusOK,
				"Data user berhasil diambil",
				result,
			)
		},
	)
}