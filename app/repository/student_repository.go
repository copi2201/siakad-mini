package repository

import (
	"fmt"

	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"siakad-mini/app/model"
)

type StudentRepository struct {
	DB *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) *StudentRepository {
	return &StudentRepository{DB: db}
}

func (r *StudentRepository) FindByID(ctx context.Context, id int64) (*model.Student, error) {
	student := &model.Student{}

	err := r.DB.QueryRow(
		ctx,
		`SELECT id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at
		 FROM students
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(
		&student.ID,
		&student.UserID,
		&student.NIM,
		&student.Nama,
		&student.Prodi,
		&student.Angkatan,
		&student.IPKTerakhir,
		&student.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return student, nil
}

func (r *StudentRepository) FindByUserID(ctx context.Context, userID int64) (*model.Student, error) {
	student := &model.Student{}

	err := r.DB.QueryRow(
		ctx,
		`SELECT id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at
		 FROM students
		 WHERE user_id = $1 AND deleted_at IS NULL`,
		userID,
	).Scan(
		&student.ID,
		&student.UserID,
		&student.NIM,
		&student.Nama,
		&student.Prodi,
		&student.Angkatan,
		&student.IPKTerakhir,
		&student.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return student, nil
}

func (r *StudentRepository) FindAll(ctx context.Context) ([]model.Student, error) {
	rows, err := r.DB.Query(
		ctx,
		`SELECT id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at
		 FROM students
		 WHERE deleted_at IS NULL
		 ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []model.Student

	for rows.Next() {
		var student model.Student

		err := rows.Scan(
			&student.ID,
			&student.UserID,
			&student.NIM,
			&student.Nama,
			&student.Prodi,
			&student.Angkatan,
			&student.IPKTerakhir,
			&student.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

func (r *StudentRepository) Create(ctx context.Context, tx pgx.Tx, student *model.Student) error {
	return tx.QueryRow(
		ctx,
		`INSERT INTO students
			(user_id, nim, nama, prodi, angkatan, ipk_terakhir)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id`,
		student.UserID,
		student.NIM,
		student.Nama,
		student.Prodi,
		student.Angkatan,
		student.IPKTerakhir,
	).Scan(&student.ID)
}

func (r *StudentRepository) SoftDelete(ctx context.Context, id int64) error {
	_, err := r.DB.Exec(
		ctx,
		`UPDATE students
		 SET deleted_at = CURRENT_TIMESTAMP
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)

	return err
}
func (r *StudentRepository) Update(
	ctx context.Context,
	id int64,
	nama string,
	prodi string,
	angkatan int,
	ipk float64,
) (*model.Student, error) {
	student := &model.Student{}

	err := r.DB.QueryRow(
		ctx,
		`UPDATE students
		 SET nama = $1,
		     prodi = $2,
		     angkatan = $3,
		     ipk_terakhir = $4
		 WHERE id = $5
		   AND deleted_at IS NULL
		 RETURNING id, user_id, nim, nama, prodi, angkatan, ipk_terakhir, deleted_at`,
		nama,
		prodi,
		angkatan,
		ipk,
		id,
	).Scan(
		&student.ID,
		&student.UserID,
		&student.NIM,
		&student.Nama,
		&student.Prodi,
		&student.Angkatan,
		&student.IPKTerakhir,
		&student.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return student, nil
}
func (r *StudentRepository) FindAllPaginated(
	ctx context.Context,
	page int,
	perPage int,
	prodi string,
	angkatan int,
	search string,
	sort string,
) ([]model.Student, int, error) {

	offset := (page - 1) * perPage

	baseQuery := `
		FROM students
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argIndex := 1

	if prodi != "" {
		baseQuery += fmt.Sprintf(" AND prodi = $%d", argIndex)
		args = append(args, prodi)
		argIndex++
	}

	if angkatan > 0 {
		baseQuery += fmt.Sprintf(" AND angkatan = $%d", argIndex)
		args = append(args, angkatan)
		argIndex++
	}

	if search != "" {
		baseQuery += fmt.Sprintf(
			" AND (nim ILIKE $%d OR nama ILIKE $%d)",
			argIndex,
			argIndex,
		)
		args = append(args, "%"+search+"%")
		argIndex++
	}

	var total int

	countQuery := "SELECT COUNT(*) " + baseQuery

	err := r.DB.QueryRow(
		ctx,
		countQuery,
		args...,
	).Scan(&total)

	if err != nil {
		return nil, 0, err
	}

	orderBy := "nama ASC"

	switch sort {
	case "nama":
		orderBy = "nama ASC"
	case "-ipk_terakhir":
		orderBy = "ipk_terakhir DESC"
	}

	dataQuery := `
		SELECT
			id,
			user_id,
			nim,
			nama,
			prodi,
			angkatan,
			ipk_terakhir,
			deleted_at
		` + baseQuery + `
		ORDER BY ` + orderBy + `
		LIMIT $` + fmt.Sprint(argIndex) + `
		OFFSET $` + fmt.Sprint(argIndex+1)

	args = append(args, perPage, offset)

	rows, err := r.DB.Query(
		ctx,
		dataQuery,
		args...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	students := []model.Student{}

	for rows.Next() {
		var student model.Student

		err := rows.Scan(
			&student.ID,
			&student.UserID,
			&student.NIM,
			&student.Nama,
			&student.Prodi,
			&student.Angkatan,
			&student.IPKTerakhir,
			&student.DeletedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return students, total, nil
}
