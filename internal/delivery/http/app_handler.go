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

type AppHandler struct {
	appService usecase.AppService
	validate   *validator.GoPlaygroundValidator
}

func NewAppHandler(as usecase.AppService, v *validator.GoPlaygroundValidator) *AppHandler {
	return &AppHandler{
		appService: as,
		validate:   v,
	}
}

// --- Helper ---

// getUserIDFromLocals adalah helper untuk mengambil ID dari middleware
func getUserIDFromLocals(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		return uuid.Nil, utils.SendError(c, fiber.StatusUnauthorized, "Invalid token")
	}
	return userID, nil
}

// --- Book Handler ---

// GetBooksWithProgress menangani 'GET /app/books'
func (h *AppHandler) GetBooksWithProgress(c *fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return err // Error sudah dikirim oleh helper
	}

	books, err := h.appService.GetBooksWithProgress(c.Context(), userID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, books)
}

// GetProgressToRestore menangani 'GET /app/books/:bookId/restore'
func (h *AppHandler) GetProgressToRestore(c *fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return err
	}

	bookIdStr := c.Params("bookId")
	bookID, err := strconv.Atoi(bookIdStr)
	if err != nil || bookID <= 0 {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid book ID")
	}

	restoreData, err := h.appService.GetProgressToRestore(c.Context(), userID, uint(bookID))
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, restoreData)
}

// UpdatePageProgress menangani 'POST /app/progress/update'
func (h *AppHandler) UpdatePageProgress(c *fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return err
	}

	var req domain.UpdateProgressRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	if err := h.appService.UpdatePageProgress(c.Context(), userID, &req); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, fiber.Map{"message": "Progress updated"})
}

// CompleteBookProgress menangani 'POST /app/progress/complete'
func (h *AppHandler) CompleteBookProgress(c *fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return err
	}

	var req domain.CompleteProgressRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	if err := h.appService.CompleteBookProgress(c.Context(), userID, &req); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, fiber.Map{"message": "Book completed and progress saved"})
}

// --- Game Handler ---

// GetAllGames menangani 'GET /app/games'
func (h *AppHandler) GetAllGames(c *fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return err
	}

	games, err := h.appService.GetAllGames(c.Context(), userID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, games)
}

// SubmitGameScore menangani 'POST /app/games/:gameId/score'
func (h *AppHandler) SubmitGameScore(c *fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return err
	}

	gameIdStr := c.Params("gameId")
	gameID, err := strconv.Atoi(gameIdStr)
	if err != nil || gameID <= 0 {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid game ID")
	}

	var req domain.SubmitGameScoreRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	if err := h.appService.SubmitGameScore(c.Context(), userID, uint(gameID), &req); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, fiber.Map{"message": "Score updated successfully"})
}
