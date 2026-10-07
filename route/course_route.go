package route

import (
	"github.com/gofiber/fiber/v2"

	"siakad-mini/app/service"
	"siakad-mini/helper"
	"siakad-mini/middleware"
)

func SetupCourseRoutes(
	api fiber.Router,
	courseService *service.CourseService,
) {
	courses := api.Group(
		"/courses",
		middleware.RequireAuth(),
	)

	// GET /api/v1/courses
	courses.Get("/", func(c *fiber.Ctx) error {
		coursesData, err := courseService.GetAll(
			c.Context(),
		)

		if err != nil {
			return helper.Error(
				c,
				fiber.StatusInternalServerError,
				"Gagal mengambil data mata kuliah",
				nil,
			)
		}

		return helper.Success(
			c,
			fiber.StatusOK,
			"Data mata kuliah berhasil diambil",
			coursesData,
		)
	})
}