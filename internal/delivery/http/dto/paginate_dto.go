package dto

type Pagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalData int64 `json:"total_data"`
	TotalPage int   `json:"total_page"`
}

type PaginatedResponse struct {
	Data interface{} `json:"data"`
	Meta *Pagination `json:"meta"`
}
