package model

type CreateStudentRequest struct {
	NIM      string  `json:"nim" validate:"required"`
	Nama     string  `json:"nama" validate:"required"`
	Email    string  `json:"email" validate:"required,email"`
	Prodi    string  `json:"prodi" validate:"required"`
	Angkatan int     `json:"angkatan" validate:"required"`
	IPK      float64 `json:"ipk_terakhir" validate:"omitempty,min=0,max=4"`
}

type UpdateStudentRequest struct {
	Nama     string  `json:"nama" validate:"required"`
	Prodi    string  `json:"prodi" validate:"required"`
	Angkatan int     `json:"angkatan" validate:"required"`
	IPK      float64 `json:"ipk_terakhir" validate:"omitempty,min=0,max=4"`
}
