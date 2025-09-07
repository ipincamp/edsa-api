package util

import (
	"math"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/dto"
)

// GeneratePagination membuat metadata pagination untuk response API.
// page: halaman saat ini
// limit: jumlah data per halaman
// totalData: total data yang tersedia
func GeneratePagination(page int, limit int, totalData int64) *dto.Pagination {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(limit)))

	return &dto.Pagination{
		Page:      page,
		Limit:     limit,
		TotalData: totalData,
		TotalPage: totalPage,
	}
}

// GetPaginationParams mengambil parameter page dan limit dari query string request.
// Jika tidak ada, akan menggunakan nilai default.
func GetPaginationParams(c *fiber.Ctx) (int, int) {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	return page, limit
}

// CalculateOffset menghitung offset untuk query database berdasarkan page dan limit.
func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}
