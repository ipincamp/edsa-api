package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	studentIdStr := c.Params("studentId")
	studentID, err := uuid.Parse(studentIdStr)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid student ID format")
	}

	// Panggil Usecase
	activities, err := h.dashboardService.GetStudentActivity(c.Context(), studentID)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, activities)
}
