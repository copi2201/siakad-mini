package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"siakad-mini/app/model"
)

type CourseRepository struct {
	DB *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) *CourseRepository {
	return &CourseRepository{DB: db}
}

func (r *CourseRepository) FindByID(ctx context.Context, id int64) (*model.Course, error) {
	course := &model.Course{}

	err := r.DB.QueryRow(
		ctx,
		`SELECT id, kode_mk, nama_mk, sks, semester, kuota
		 FROM courses
		 WHERE id = $1`,
		id,
	).Scan(
		&course.ID,
		&course.KodeMK,
		&course.NamaMK,
		&course.SKS,
		&course.Semester,
		&course.Kuota,
	)

	if err != nil {
		return nil, err
	}

	return course, nil
}

func (r *CourseRepository) FindAll(ctx context.Context) ([]model.Course, error) {
	rows, err := r.DB.Query(
		ctx,
		`SELECT 
			c.id,
			c.kode_mk,
			c.nama_mk,
			c.sks,
			c.semester,
			c.kuota,
			COUNT(e.id) AS terisi,
			c.kuota - COUNT(e.id) AS sisa_kuota
		FROM courses c
		LEFT JOIN enrollments e ON e.course_id = c.id
		GROUP BY c.id
		ORDER BY c.semester, c.kode_mk`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []model.Course

	for rows.Next() {
		var course model.Course

		err := rows.Scan(
			&course.ID,
			&course.KodeMK,
			&course.NamaMK,
			&course.SKS,
			&course.Semester,
			&course.Kuota,
			&course.Terisi,
			&course.SisaKuota,
		)
		if err != nil {
			return nil, err
		}

		courses = append(courses, course)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return courses, nil
}