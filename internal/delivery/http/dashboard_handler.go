package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/pkg/validator"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type DashboardHandler struct {
	dashboardService usecase.DashboardService
	validate         *validator.GoPlaygroundValidator
}

func NewDashboardHandler(ds usecase.DashboardService, v *validator.GoPlaygroundValidator) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: ds,
		validate:         v,
	}
}

// GetStudentActivity menangani 'GET /dashboard/students/:studentId/activity'
func (h *DashboardHandler) GetStudentActivity(c *fiber.Ctx) error {
	// 1. Ambil 'id' (studentId) dari URL
	studentIdStr := c.Params("studentId")
	studentID, err := uuid.Parse(studentIdStr)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid student ID format", err.Error())
	}

	// 2. Ambil filter paginasi dari query params
	// 'id' sudah didapat dari URL param (studentID)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	startDate := c.Query("start_date") // cth: 2025-10-20
	endDate := c.Query("end_date")     // cth: 2025-10-25

	filters := &domain.ActivityLogQuery{
		Page:      page,
		Limit:     limit,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// 3. Panggil Usecase
	paginatedData, err := h.dashboardService.GetStudentActivity(c.Context(), studentID, filters)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}

	// 4. Kembalikan menggunakan format SendPagination
	meta := utils.PaginationMeta{
		Page:      paginatedData.Meta.Page,
		Limit:     paginatedData.Meta.Limit,
		TotalPage: paginatedData.Meta.TotalPage,
		TotalData: paginatedData.Meta.TotalData,
	}

	return utils.SendPagination(c, "Student activity retrieved successfully", paginatedData.List, meta)
}
