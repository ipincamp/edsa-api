package domain

// PaginationMetaDTO adalah DTO untuk metadata paginasi
type PaginationMetaDTO struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalPage int64 `json:"total_page"`
	TotalData int64 `json:"total_data"`
}

// PaginatedDTO adalah DTO generik untuk respons paginasi
type PaginatedDTO struct {
	List interface{}       `json:"list"`
	Meta PaginationMetaDTO `json:"meta"`
}
