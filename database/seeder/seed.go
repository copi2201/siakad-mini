package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"siakad-mini/helper"
)

type StudentSeed struct {
	NIM      string
	Nama     string
	Email    string
	Prodi    string
	Angkatan int
	IPK      float64
}

type CourseSeed struct {
	KodeMK   string
	NamaMK   string
	SKS      int
	Semester int
	Kuota    int
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Gagal membaca file .env:", err)
	}

	user := url.QueryEscape(os.Getenv("DB_USER"))
	password := url.QueryEscape(os.Getenv("DB_PASSWORD"))
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if user == "" || password == "" || host == "" || port == "" || dbName == "" {
		log.Fatal("Konfigurasi database pada .env belum lengkap")
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		dbName,
	)

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Gagal membuat koneksi database:", err)
	}
	defer db.Close()

	ctx := context.Background()

	if err := db.Ping(ctx); err != nil {
		log.Fatal("Gagal koneksi database:", err)
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		log.Fatal("Gagal memulai transaksi:", err)
	}
	defer tx.Rollback(ctx)

	// =========================
	// ADMIN
	// =========================

	adminPassword, err := helper.HashPassword("admin12345")
	if err != nil {
		log.Fatal(err)
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO users (email, password, role)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (email) DO UPDATE
		 SET password = EXCLUDED.password,
		     role = EXCLUDED.role`,
		"admin@siakad.local",
		adminPassword,
		"admin",
	)

	if err != nil {
		log.Fatal("Gagal membuat admin:", err)
	}

	// =========================
	// MAHASISWA
	// =========================

	students := []StudentSeed{
		{"240411000001", "Mahasiswa 01", "mhs01@siakad.local", "D4 Teknik Informatika", 2024, 3.80},
		{"240411000002", "Mahasiswa 02", "mhs02@siakad.local", "D4 Teknik Informatika", 2024, 3.50},
		{"240411000003", "Mahasiswa 03", "mhs03@siakad.local", "D4 Teknik Informatika", 2024, 3.20},
		{"240411000004", "Mahasiswa 04", "mhs04@siakad.local", "D4 Teknik Informatika", 2024, 3.00},
		{"240411000005", "Mahasiswa 05", "mhs05@siakad.local", "D4 Teknik Informatika", 2024, 2.90},
		{"240411000006", "Mahasiswa 06", "mhs06@siakad.local", "D4 Teknik Informatika", 2024, 2.80},
		{"240411000007", "Mahasiswa 07", "mhs07@siakad.local", "D4 Teknik Informatika", 2024, 2.70},
		{"240411000008", "Mahasiswa 08", "mhs08@siakad.local", "D4 Teknik Informatika", 2024, 2.60},
		{"240411000009", "Mahasiswa 09", "mhs09@siakad.local", "D4 Teknik Informatika", 2024, 2.50},
		{"240411000010", "Mahasiswa 10", "mhs10@siakad.local", "D4 Teknik Informatika", 2024, 2.40},
		{"240411000011", "Mahasiswa 11", "mhs11@siakad.local", "D4 Teknik Informatika", 2024, 2.30},
		{"240411000012", "Mahasiswa 12", "mhs12@siakad.local", "D4 Teknik Informatika", 2024, 3.70},
		{"240411000013", "Mahasiswa 13", "mhs13@siakad.local", "D4 Teknik Informatika", 2024, 3.60},
		{"240411000014", "Mahasiswa 14", "mhs14@siakad.local", "D4 Teknik Informatika", 2024, 3.40},
		{"240411000015", "Mahasiswa 15", "mhs15@siakad.local", "D4 Teknik Informatika", 2024, 3.10},
		{"240411000016", "Mahasiswa 16", "mhs16@siakad.local", "D4 Teknik Informatika", 2024, 2.95},
		{"240411000017", "Mahasiswa 17", "mhs17@siakad.local", "D4 Teknik Informatika", 2024, 2.85},
		{"240411000018", "Mahasiswa 18", "mhs18@siakad.local", "D4 Teknik Informatika", 2024, 2.75},
		{"240411000019", "Mahasiswa 19", "mhs19@siakad.local", "D4 Teknik Informatika", 2024, 2.65},
		{"240411000020", "Mahasiswa 20", "mhs20@siakad.local", "D4 Teknik Informatika", 2024, 2.45},
	}

	for _, student := range students {
		hashedPassword, err := helper.HashPassword(student.NIM)
		if err != nil {
			log.Fatal(err)
		}

		var userID int64

		err = tx.QueryRow(
			ctx,
			`INSERT INTO users (email, password, role)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (email) DO UPDATE
			 SET password = EXCLUDED.password,
			     role = EXCLUDED.role
			 RETURNING id`,
			student.Email,
			hashedPassword,
			"mahasiswa",
		).Scan(&userID)

		if err != nil {
			log.Fatal("Gagal membuat user mahasiswa:", err)
		}

		_, err = tx.Exec(
			ctx,
			`INSERT INTO students
				(user_id, nim, nama, prodi, angkatan, ipk_terakhir)
			 VALUES ($1, $2, $3, $4, $5, $6)
			 ON CONFLICT (nim) DO UPDATE
			 SET nama = EXCLUDED.nama,
			     prodi = EXCLUDED.prodi,
			     angkatan = EXCLUDED.angkatan,
			     ipk_terakhir = EXCLUDED.ipk_terakhir,
			     deleted_at = NULL`,
			userID,
			student.NIM,
			student.Nama,
			student.Prodi,
			student.Angkatan,
			student.IPK,
		)

		if err != nil {
			log.Fatal("Gagal membuat student:", err)
		}
	}

	// =========================
	// MATA KULIAH
	// =========================

	courses := []CourseSeed{
		{"IF101", "Algoritma dan Pemrograman", 3, 1, 30},
		{"IF102", "Basis Data", 3, 2, 30},
		{"IF103", "Pemrograman Web", 3, 3, 30},
		{"IF104", "Pemrograman Backend", 3, 4, 30},
		{"IF105", "Keamanan Siber", 3, 5, 30},
		{"IF106", "Analisis Data", 3, 5, 30},
		{"IF107", "Rekayasa Perangkat Lunak", 3, 4, 30},
		{"IF108", "Jaringan Komputer", 3, 3, 30},
		{"IF109", "Sistem Operasi", 3, 2, 30},
		{"IF110", "Interaksi Manusia dan Komputer", 3, 4, 30},
	}

	for _, course := range courses {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO courses
				(kode_mk, nama_mk, sks, semester, kuota)
			 VALUES ($1, $2, $3, $4, $5)
			 ON CONFLICT (kode_mk) DO UPDATE
			 SET nama_mk = EXCLUDED.nama_mk,
			     sks = EXCLUDED.sks,
			     semester = EXCLUDED.semester,
			     kuota = EXCLUDED.kuota`,
			course.KodeMK,
			course.NamaMK,
			course.SKS,
			course.Semester,
			course.Kuota,
		)

		if err != nil {
			log.Fatal("Gagal membuat course:", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.Fatal("Gagal commit seeder:", err)
	}

	fmt.Println("Seeder berhasil dijalankan.")
	fmt.Println("Admin       : admin@siakad.local")
	fmt.Println("Password    : admin12345")
	fmt.Println("Mahasiswa   : mhs01@siakad.local s.d. mhs20@siakad.local")
	fmt.Println("Password    : NIM masing-masing")
}
