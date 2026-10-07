package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"siakad-mini/app/model"
	"siakad-mini/app/repository"
	"siakad-mini/helper"
)

type AuthService struct {
	UserRepo    *repository.UserRepository
	StudentRepo *repository.StudentRepository
}

func NewAuthService(
	userRepo *repository.UserRepository,
	studentRepo *repository.StudentRepository,
) *AuthService {
	return &AuthService{
		UserRepo:    userRepo,
		StudentRepo: studentRepo,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	req model.LoginRequest,
) (*model.LoginResponse, error) {

	user, err := s.UserRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("email atau password salah")
		}

		return nil, err
	}

	if !helper.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("email atau password salah")
	}

	// Mahasiswa yang sudah soft delete tidak boleh login.
	if user.Role == "mahasiswa" {
		_, err := s.StudentRepo.FindByUserID(ctx, user.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, errors.New("akun mahasiswa tidak aktif")
			}

			return nil, err
		}
	}

	token, expiresIn, err := helper.GenerateToken(
		user.ID,
		user.Role,
	)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expiresIn,
		User: model.User{
			ID:    user.ID,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}