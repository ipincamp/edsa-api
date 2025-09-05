package util

import (
	"math"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/domain/dto"
)

func GeneratePagination(page, limit int, total int64) *dto.Pagination {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

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

func CalculateOffset(page, limit int) int {
	return (page - 1) * limit
}
