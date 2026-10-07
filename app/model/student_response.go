package model

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	PerPage     int `json:"per_page"`
	Total       int `json:"total"`
	LastPage    int `json:"last_page"`
}

type StudentListResponse struct {
	Data []Student      `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type StudentDetailResponse struct {
	Student  Student  `json:"student"`
	Courses  []Course `json:"courses"`
	TotalSKS int      `json:"total_sks"`
	BatasSKS int      `json:"batas_sks"`
}
