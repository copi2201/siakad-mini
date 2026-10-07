package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"siakad-mini/app/model"
)

func (s *AuthService) GetMe(
	ctx context.Context,
	userID int64,
) (*model.AuthMeResponse, error) {

	user, err := s.UserRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user tidak ditemukan")
		}

		return nil, err
	}

	response := &model.AuthMeResponse{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}

	if user.Role == "mahasiswa" {
		student, err := s.StudentRepo.FindByUserID(ctx, user.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errors.New("data mahasiswa tidak ditemukan")
			}

			return nil, err
		}

		response.Student = &model.StudentInfo{
			NIM:      student.NIM,
			Nama:     student.Nama,
			Prodi:    student.Prodi,
			Angkatan: student.Angkatan,
		}
	}

	return response, nil
}