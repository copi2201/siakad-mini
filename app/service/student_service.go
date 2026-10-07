package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"siakad-mini/app/model"
	"siakad-mini/app/repository"
	"siakad-mini/helper"
)

type StudentService struct {
	DB          *pgxpool.Pool
	UserRepo    *repository.UserRepository
	StudentRepo *repository.StudentRepository
}

func NewStudentService(
	db *pgxpool.Pool,
	userRepo *repository.UserRepository,
	studentRepo *repository.StudentRepository,
) *StudentService {
	return &StudentService{
		DB:          db,
		UserRepo:    userRepo,
		StudentRepo: studentRepo,
	}
}

func (s *StudentService) Create(
	ctx context.Context,
	req model.CreateStudentRequest,
) (*model.Student, error) {

	if err := ValidateNIM(req.NIM); err != nil {
		return nil, err
	}

	if err := ValidateAngkatan(req.Angkatan); err != nil {
		return nil, err
	}

	if err := ValidateIPK(req.IPK); err != nil {
		return nil, err
	}

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	hashedPassword, err := helper.HashPassword(req.NIM)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "mahasiswa",
	}

	if err := s.UserRepo.Create(ctx, tx, user); err != nil {
		return nil, err
	}

	student := &model.Student{
		UserID:      user.ID,
		NIM:         req.NIM,
		Nama:        req.Nama,
		Prodi:       req.Prodi,
		Angkatan:    req.Angkatan,
		IPKTerakhir: req.IPK,
	}

	if err := s.StudentRepo.Create(ctx, tx, student); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return student, nil
}

func (s *StudentService) Update(
	ctx context.Context,
	id int64,
	req model.UpdateStudentRequest,
) (*model.Student, error) {

	if err := ValidateAngkatan(req.Angkatan); err != nil {
		return nil, err
	}

	if err := ValidateIPK(req.IPK); err != nil {
		return nil, err
	}

	student, err := s.StudentRepo.Update(
		ctx,
		id,
		req.Nama,
		req.Prodi,
		req.Angkatan,
		req.IPK,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("data mahasiswa tidak ditemukan")
		}

		return nil, err
	}

	return student, nil
}
func (s *StudentService) GetAll(
	ctx context.Context,
	page int,
	perPage int,
	prodi string,
	angkatan int,
	search string,
	sort string,
) ([]model.Student, int, error) {

	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 10
	}

	if perPage > 50 {
		perPage = 50
	}

	students, total, err := s.StudentRepo.FindAllPaginated(
		ctx,
		page,
		perPage,
		prodi,
		angkatan,
		search,
		sort,
	)

	if err != nil {
		return nil, 0, err
	}

	return students, total, nil
}
func (s *StudentService) GetDetail(
	ctx context.Context,
	id int64,
) (*model.StudentDetailResponse, error) {

	student, err := s.StudentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("data mahasiswa tidak ditemukan")
		}

		return nil, err
	}

	enrollmentRepo := repository.NewEnrollmentRepository(s.DB)

	courses, err := enrollmentRepo.FindCoursesByStudentID(
		ctx,
		student.ID,
	)
	if err != nil {
		return nil, err
	}

	totalSKS := 0

	for _, course := range courses {
		totalSKS += course.SKS
	}

	batasSKS := GetSKSLimit(student.IPKTerakhir)

	return &model.StudentDetailResponse{
		Student:  *student,
		Courses:  courses,
		TotalSKS: totalSKS,
		BatasSKS: batasSKS,
	}, nil
}
func (s *StudentService) Delete(
	ctx context.Context,
	id int64,
) error {

	err := s.StudentRepo.SoftDelete(ctx, id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("data mahasiswa tidak ditemukan")
		}

		return err
	}

	return nil
}
