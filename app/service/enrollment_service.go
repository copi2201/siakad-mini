package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"siakad-mini/app/model"
	"siakad-mini/app/repository"
)

type EnrollmentService struct {
	DB             *pgxpool.Pool
	EnrollmentRepo *repository.EnrollmentRepository
	StudentRepo    *repository.StudentRepository
	CourseRepo     *repository.CourseRepository
}

func NewEnrollmentService(
	db *pgxpool.Pool,
	enrollmentRepo *repository.EnrollmentRepository,
	studentRepo *repository.StudentRepository,
	courseRepo *repository.CourseRepository,
) *EnrollmentService {
	return &EnrollmentService{
		DB:             db,
		EnrollmentRepo: enrollmentRepo,
		StudentRepo:    studentRepo,
		CourseRepo:     courseRepo,
	}
}
func (s *EnrollmentService) Create(
	ctx context.Context,
	studentID int64,
	courseID int64,
	tahunAkademik string,
) (*model.Enrollment, error) {

	student, err := s.StudentRepo.FindByID(ctx, studentID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("data mahasiswa tidak ditemukan")
		}

		return nil, err
	}

	course, err := s.CourseRepo.FindByID(ctx, courseID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("mata kuliah tidak ditemukan")
		}

		return nil, err
	}

	exists, err := s.EnrollmentRepo.Exists(
		ctx,
		studentID,
		courseID,
		tahunAkademik,
	)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, errors.New(
			"mata kuliah sudah diambil pada tahun akademik tersebut",
		)
	}

	totalSKS, err := s.EnrollmentRepo.GetTotalSKS(
		ctx,
		studentID,
		tahunAkademik,
	)
	if err != nil {
		return nil, err
	}

	batasSKS := GetSKSLimit(student.IPKTerakhir)

	if totalSKS+course.SKS > batasSKS {
		return nil, errors.New(
			"pengambilan mata kuliah melebihi batas SKS",
		)
	}

	if course.Terisi >= course.Kuota {
		return nil, errors.New(
			"kuota mata kuliah sudah penuh",
		)
	}

	enrollment := &model.Enrollment{
		StudentID:     studentID,
		CourseID:      courseID,
		TahunAkademik: tahunAkademik,
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if err := s.EnrollmentRepo.Create(
		ctx,
		tx,
		enrollment,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return enrollment, nil
}
func (s *EnrollmentService) Delete(
	ctx context.Context,
	enrollmentID int64,
	studentID int64,
) error {
	err := s.EnrollmentRepo.Delete(
		ctx,
		enrollmentID,
		studentID,
	)
	if err != nil {
		return err
	}

	return nil
}