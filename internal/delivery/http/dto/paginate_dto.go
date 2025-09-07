package dto

// Pagination adalah DTO untuk metadata pagination pada response API
type Pagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalData int64 `json:"total_data"`
	TotalPage int   `json:"total_page"`
}

// PaginatedResponse adalah DTO untuk response data yang dipaginasi
type PaginatedResponse struct {
	Data interface{} `json:"data"`
	Meta *Pagination `json:"meta"`
}
