package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"siakad-mini/app/model"
)

type EnrollmentRepository struct {
	DB *pgxpool.Pool
}

func NewEnrollmentRepository(db *pgxpool.Pool) *EnrollmentRepository {
	return &EnrollmentRepository{DB: db}
}

func (r *EnrollmentRepository) FindByID(ctx context.Context, id int64) (*model.Enrollment, error) {
	enrollment := &model.Enrollment{}

	err := r.DB.QueryRow(
		ctx,
		`SELECT id, student_id, course_id, tahun_akademik, created_at
		 FROM enrollments
		 WHERE id = $1`,
		id,
	).Scan(
		&enrollment.ID,
		&enrollment.StudentID,
		&enrollment.CourseID,
		&enrollment.TahunAkademik,
		&enrollment.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return enrollment, nil
}

func (r *EnrollmentRepository) Exists(
	ctx context.Context,
	studentID int64,
	courseID int64,
	tahunAkademik string,
) (bool, error) {
	var exists bool

	err := r.DB.QueryRow(
		ctx,
		`SELECT EXISTS(
			SELECT 1
			FROM enrollments
			WHERE student_id = $1
			AND course_id = $2
			AND tahun_akademik = $3
		)`,
		studentID,
		courseID,
		tahunAkademik,
	).Scan(&exists)

	return exists, err
}

func (r *EnrollmentRepository) Create(
	ctx context.Context,
	tx pgx.Tx,
	enrollment *model.Enrollment,
) error {
	return tx.QueryRow(
		ctx,
		`INSERT INTO enrollments
			(student_id, course_id, tahun_akademik)
		 VALUES ($1, $2, $3)
		 RETURNING id, created_at`,
		enrollment.StudentID,
		enrollment.CourseID,
		enrollment.TahunAkademik,
	).Scan(
		&enrollment.ID,
		&enrollment.CreatedAt,
	)
}

func (r *EnrollmentRepository) Delete(
	ctx context.Context,
	id int64,
	studentID int64,
) error {
	result, err := r.DB.Exec(
		ctx,
		`DELETE FROM enrollments
		 WHERE id = $1 AND student_id = $2`,
		id,
		studentID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *EnrollmentRepository) GetTotalSKS(
	ctx context.Context,
	studentID int64,
	tahunAkademik string,
) (int, error) {
	var totalSKS int

	err := r.DB.QueryRow(
		ctx,
		`SELECT COALESCE(SUM(c.sks), 0)
		 FROM enrollments e
		 JOIN courses c ON c.id = e.course_id
		 WHERE e.student_id = $1
		 AND e.tahun_akademik = $2`,
		studentID,
		tahunAkademik,
	).Scan(&totalSKS)

	return totalSKS, err
}

func (r *EnrollmentRepository) FindByStudentID(
	ctx context.Context,
	studentID int64,
) ([]model.Enrollment, error) {
	rows, err := r.DB.Query(
		ctx,
		`SELECT id, student_id, course_id, tahun_akademik, created_at
		 FROM enrollments
		 WHERE student_id = $1
		 ORDER BY created_at DESC`,
		studentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enrollments []model.Enrollment

	for rows.Next() {
		var enrollment model.Enrollment

		err := rows.Scan(
			&enrollment.ID,
			&enrollment.StudentID,
			&enrollment.CourseID,
			&enrollment.TahunAkademik,
			&enrollment.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		enrollments = append(enrollments, enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return enrollments, nil
}
func (r *EnrollmentRepository) FindCoursesByStudentID(
	ctx context.Context,
	studentID int64,
) ([]model.Course, error) {
	rows, err := r.DB.Query(
		ctx,
		`SELECT
			c.id,
			c.kode_mk,
			c.nama_mk,
			c.sks,
			c.semester,
			c.kuota,
			(
				SELECT COUNT(*)
				FROM enrollments e2
				WHERE e2.course_id = c.id
			) AS terisi,
			(
				c.kuota - (
					SELECT COUNT(*)
					FROM enrollments e3
					WHERE e3.course_id = c.id
				)
			) AS sisa_kuota
		FROM courses c
		INNER JOIN enrollments e
			ON e.course_id = c.id
		WHERE e.student_id = $1
		ORDER BY c.semester, c.kode_mk`,
		studentID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses := []model.Course{}

	for rows.Next() {
		var course model.Course

		if err := rows.Scan(
			&course.ID,
			&course.KodeMK,
			&course.NamaMK,
			&course.SKS,
			&course.Semester,
			&course.Kuota,
			&course.Terisi,
			&course.SisaKuota,
		); err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}
