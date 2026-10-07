package route

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"siakad-mini/app/model"
	"siakad-mini/app/service"
	"siakad-mini/helper"
	"siakad-mini/middleware"
)

func SetupEnrollmentRoutes(
	api fiber.Router,
	enrollmentService *service.EnrollmentService,
) {
	enrollments := api.Group(
		"/enrollments",
		middleware.RequireAuth(),
		middleware.RequireRole("mahasiswa"),
	)

	// POST /api/v1/enrollments
	enrollments.Post("/", func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(int64)
		if !ok {
			return helper.Error(c, 401, "User tidak valid", nil)
		}

		var req model.CreateEnrollmentRequest

		if err := c.BodyParser(&req); err != nil {
			return helper.Error(c, 422, "Request tidak valid", nil)
		}

		if req.CourseID <= 0 {
			return helper.Error(c, 422, "course_id wajib diisi", nil)
		}

		if req.TahunAkademik == "" {
			return helper.Error(c, 422, "tahun_akademik wajib diisi", nil)
		}

		student, err := enrollmentService.StudentRepo.FindByUserID(
			c.Context(),
			userID,
		)
		if err != nil {
			return helper.Error(
				c,
				404,
				"Data mahasiswa tidak ditemukan",
				nil,
			)
		}

		enrollment, err := enrollmentService.Create(
			c.Context(),
			student.ID,
			req.CourseID,
			req.TahunAkademik,
		)
		if err != nil {
			return helper.Error(c, 422, err.Error(), nil)
		}

		return helper.Success(
			c,
			201,
			"KRS berhasil ditambahkan",
			enrollment,
		)
	})

	// DELETE /api/v1/enrollments/:id
	enrollments.Delete("/:id", func(c *fiber.Ctx) error {
		enrollmentID, err := strconv.ParseInt(
			c.Params("id"),
			10,
			64,
		)
		if err != nil {
			return helper.Error(
				c,
				422,
				"ID KRS tidak valid",
				nil,
			)
		}

		userID, ok := c.Locals("user_id").(int64)
		if !ok {
			return helper.Error(c, 401, "User tidak valid", nil)
		}

		student, err := enrollmentService.StudentRepo.FindByUserID(
			c.Context(),
			userID,
		)
		if err != nil {
			return helper.Error(
				c,
				404,
				"Data mahasiswa tidak ditemukan",
				nil,
			)
		}

		err = enrollmentService.Delete(
			c.Context(),
			enrollmentID,
			student.ID,
		)
		if err != nil {
			return helper.Error(
				c,
				404,
				"KRS tidak ditemukan",
				nil,
			)
		}

		return c.SendStatus(204)
	})
}