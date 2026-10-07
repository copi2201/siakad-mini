package service

import (
	"context"

	"siakad-mini/app/model"
	"siakad-mini/app/repository"
)

type CourseService struct {
	CourseRepo *repository.CourseRepository
}

func NewCourseService(
	courseRepo *repository.CourseRepository,
) *CourseService {
	return &CourseService{
		CourseRepo: courseRepo,
	}
}

func (s *CourseService) GetAll(
	ctx context.Context,
) ([]model.Course, error) {
	return s.CourseRepo.FindAll(ctx)
}
