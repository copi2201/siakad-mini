package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"siakad-mini/app/repository"
	"siakad-mini/app/service"
	"siakad-mini/config"
	"siakad-mini/database"
	"siakad-mini/route"
)

func main() {
	config.LoadEnv()

	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	enrollmentRepo := repository.NewEnrollmentRepository(db)

	authService := service.NewAuthService(userRepo, studentRepo)
	studentService := service.NewStudentService(
		db,
		userRepo,
		studentRepo,
	)
	courseService := service.NewCourseService(courseRepo)
	enrollmentService := service.NewEnrollmentService(
		db,
		enrollmentRepo,
		studentRepo,
		courseRepo,
	)

	app := fiber.New()

	api := app.Group("/api/v1")

	route.SetupAuthRoutes(api, authService)
	route.SetupStudentRoutes(api, studentService)
	route.SetupCourseRoutes(api, courseService)
	route.SetupEnrollmentRoutes(api, enrollmentService)

	log.Println("Server berjalan di http://localhost:3000")

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}