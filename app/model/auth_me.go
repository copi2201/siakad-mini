package model

type AuthMeResponse struct {
	ID       int64        `json:"id"`
	Email    string       `json:"email"`
	Role     string       `json:"role"`
	Student  *StudentInfo `json:"student,omitempty"`
}

type StudentInfo struct {
	NIM      string `json:"nim"`
	Nama     string `json:"nama"`
	Prodi    string `json:"prodi"`
	Angkatan int    `json:"angkatan"`
}