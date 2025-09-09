package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/dto"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/service"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"github.com/ipincamp/go-edsa-api/pkg/validator"
)

type LearningHandler struct {
	learningService service.LearningService
	validator       *validator.CustomValidator
}

func NewLearningHandler(learningService service.LearningService) *LearningHandler {
	return &LearningHandler{
		learningService: learningService,
		validator:       validator.NewValidator(),
	}
}

// GetAvailableBooks - Handler to get all books with user-specific status
func (h *LearningHandler) GetAvailableBooks(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)

	books, err := h.learningService.GetAvailableBooks(c.Context(), user.ID)
	if err != nil {
		return util.SendError(c, fiber.StatusInternalServerError, "failed to retrieve books")
	}

	return util.SendSuccess(c, fiber.StatusOK, "books retrieved successfully", books)
}

// GetBookDetail - Handler to get detail of a specific book
func (h *LearningHandler) GetBookDetail(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)
	bookID := c.Params("bookID")

	bookDetail, err := h.learningService.GetBookDetail(c.Context(), user.ID, bookID)
	if err != nil {
		return util.SendError(c, fiber.StatusInternalServerError, "failed to retrieve book detail")
	}

	return util.SendSuccess(c, fiber.StatusOK, "book detail retrieved successfully", bookDetail)
}

// SubmitInteraction - Handler for submitting interaction results
func (h *LearningHandler) SubmitInteraction(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)

	var req dto.SubmitInteractionRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if errs := h.validator.Validate(req); errs != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", errs)
	}

	err := h.learningService.SubmitInteraction(c.Context(), user.ID, req)
	if err != nil {
		return util.SendError(c, fiber.StatusInternalServerError, "failed to submit interaction")
	}

	return util.SendSuccess(c, fiber.StatusOK, "progress updated successfully", nil)
}
