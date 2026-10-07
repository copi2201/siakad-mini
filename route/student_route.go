package route

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"siakad-mini/app/model"
	"siakad-mini/app/service"
	"siakad-mini/helper"
	"siakad-mini/middleware"
)

func SetupStudentRoutes(
	api fiber.Router,
	studentService *service.StudentService,
) {
	// Route mahasiswa yang hanya dapat diakses admin
	adminStudents := api.Group(
		"/students",
		middleware.RequireAuth(),
		middleware.RequireRole("admin"),
	)

	// GET /api/v1/students
	adminStudents.Get("/", func(c *fiber.Ctx) error {
		page := c.QueryInt("page", 1)
		perPage := c.QueryInt("per_page", 10)

		if page < 1 {
			page = 1
		}

		if perPage < 1 {
			perPage = 10
		}

		if perPage > 50 {
			perPage = 50
		}

		angkatan := c.QueryInt("angkatan", 0)
		prodi := c.Query("prodi")
		search := c.Query("search")
		sort := c.Query("sort", "nama")

		if sort != "nama" && sort != "-ipk_terakhir" {
			return helper.Error(
				c,
				fiber.StatusUnprocessableEntity,
				"Parameter sort tidak valid",
				nil,
			)
		}

		studentsData, total, err := studentService.GetAll(
			c.Context(),
			page,
			perPage,
			prodi,
			angkatan,
			search,
			sort,
		)

		if err != nil {
			return helper.Error(
				c,
				fiber.StatusInternalServerError,
				"Gagal mengambil data mahasiswa",
				nil,
			)
		}

		lastPage := 0

		if total > 0 {
			lastPage = int(
				math.Ceil(
					float64(total) / float64(perPage),
				),
			)
		}

		meta := model.PaginationMeta{
			CurrentPage: page,
			PerPage:     perPage,
			Total:       total,
			LastPage:    lastPage,
		}

		return helper.SuccessWithMeta(
			c,
			fiber.StatusOK,
			"Data mahasiswa berhasil diambil",
			studentsData,
			meta,
		)
	})
	// GET /api/v1/students/{id}
	studentsDetail := api.Group(
		"/students",
		middleware.RequireAuth(),
	)

	studentsDetail.Get("/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(
			c.Params("id"),
			10,
			64,
		)

		if err != nil || id <= 0 {
			return helper.Error(
				c,
				fiber.StatusUnprocessableEntity,
				"ID mahasiswa tidak valid",
				nil,
			)
		}

		student, err := studentService.GetDetail(
			c.Context(),
			id,
		)

		if err != nil {
			return helper.Error(
				c,
				fiber.StatusNotFound,
				err.Error(),
				nil,
			)
		}

		role := c.Locals("role")
		userID := c.Locals("user_id").(int64)

		// Admin dapat melihat semua mahasiswa.
		if role == "admin" {
			return helper.Success(
				c,
				fiber.StatusOK,
				"Data mahasiswa berhasil diambil",
				student,
			)
		}

		// Mahasiswa hanya dapat melihat data dirinya sendiri.
		if role == "mahasiswa" {
			if student.Student.UserID != userID {
				return helper.Error(
					c,
					fiber.StatusForbidden,
					"Akses ditolak",
					nil,
				)
			}

			return helper.Success(
				c,
				fiber.StatusOK,
				"Data mahasiswa berhasil diambil",
				student,
			)
		}

		return helper.Error(
			c,
			fiber.StatusForbidden,
			"Akses ditolak",
			nil,
		)
	})

	// POST /api/v1/students
	adminStudents.Post("/", func(c *fiber.Ctx) error {
		var req model.CreateStudentRequest

		if err := c.BodyParser(&req); err != nil {
			return helper.Error(
				c,
				fiber.StatusUnprocessableEntity,
				"Format request tidak valid",
				nil,
			)
		}

		student, err := studentService.Create(
			c.Context(),
			req,
		)

		if err != nil {
			return helper.Error(
				c,
				fiber.StatusUnprocessableEntity,
				err.Error(),
				nil,
			)
		}

		return helper.Success(
			c,
			fiber.StatusCreated,
			"Mahasiswa berhasil dibuat",
			student,
		)
	})

	// PUT /api/v1/students/{id}
	adminStudents.Put("/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(
			c.Params("id"),
			10,
			64,
		)

		if err != nil || id <= 0 {
			return helper.Error(
				c,
				fiber.StatusUnprocessableEntity,
				"ID mahasiswa tidak valid",
				nil,
			)
		}

		var req model.UpdateStudentRequest

		if err := c.BodyParser(&req); err != nil {
			return helper.Error(
				c,
				fiber.StatusUnprocessableEntity,
				"Format request tidak valid",
				nil,
			)
		}

		student, err := studentService.Update(
			c.Context(),
			id,
			req,
		)

		if err != nil {
			return helper.Error(
				c,
				fiber.StatusUnprocessableEntity,
				err.Error(),
				nil,
			)
		}

		return helper.Success(
			c,
			fiber.StatusOK,
			"Data mahasiswa berhasil diperbarui",
			student,
		)
	})
	// DELETE /api/v1/students/{id}
	adminStudents.Delete("/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(
			c.Params("id"),
			10,
			64,
		)

		if err != nil || id <= 0 {
			return helper.Error(
				c,
				fiber.StatusUnprocessableEntity,
				"ID mahasiswa tidak valid",
				nil,
			)
		}

		if err := studentService.Delete(
			c.Context(),
			id,
		); err != nil {
			if err.Error() == "data mahasiswa tidak ditemukan" {
				return helper.Error(
					c,
					fiber.StatusNotFound,
					err.Error(),
					nil,
				)
			}

			return helper.Error(
				c,
				fiber.StatusInternalServerError,
				"Gagal menghapus mahasiswa",
				nil,
			)
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}
