package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"siakad-mini/app/model"
)

type UserRepository struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	err := r.DB.QueryRow(
		ctx,
		`SELECT id, email, password, role
		 FROM users
		 WHERE email = $1`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}

	err := r.DB.QueryRow(
		ctx,
		`SELECT id, email, password, role
		 FROM users
		 WHERE id = $1`,
		id,
	).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
func (r *UserRepository) Create(
	ctx context.Context,
	tx pgx.Tx,
	user *model.User,
) error {
	return tx.QueryRow(
		ctx,
		`INSERT INTO users
			(email, password, role)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		user.Email,
		user.Password,
		user.Role,
	).Scan(&user.ID)
}
