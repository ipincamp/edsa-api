package http

import (
	"errors"
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

// getUintIDParam adalah helper untuk mengambil parameter ID dari URL dan mengonversinya ke uint
// (Helper ini disalin dari admin_handler.go untuk konsistensi)
func (h *DashboardHandler) getUintIDParam(c *fiber.Ctx, paramName string) (uint, error) {
	idStr := c.Params(paramName)
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id parameter: " + paramName)
	}
	return uint(id), nil
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

// UnlockBookForGroup menangani 'POST /dashboard/groups/:groupId/books/:bookId/unlock'
func (h *DashboardHandler) UnlockBookForGroup(c *fiber.Ctx) error {
	// 1. Ambil 'groupId' dari URL
	groupID, err := h.getUintIDParam(c, "groupId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	// 2. Ambil 'bookId' dari URL
	bookID, err := h.getUintIDParam(c, "bookId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	// 3. Panggil Usecase
	if err := h.dashboardService.UnlockBookForGroup(c.Context(), groupID, bookID); err != nil {
		// Cek apakah errornya adalah "not found"
		if err.Error() == "group not found" || err.Error() == "book not found" {
			return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
		}
		// Error lain
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	// 4. Kembalikan sukses
	return utils.SendSuccess(c, fiber.StatusOK, "Book unlocked for group successfully", nil)
}
