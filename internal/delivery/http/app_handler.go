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
	appService       usecase.AppService
	dashboardService usecase.DashboardService
	validate         *validator.GoPlaygroundValidator
}

func NewAppHandler(as usecase.AppService, ds usecase.DashboardService, v *validator.GoPlaygroundValidator) *AppHandler {
	return &AppHandler{
		appService:       as,
		dashboardService: ds,
		validate:         v,
	}
}

// --- Helper ---

// getUserIDFromLocals adalah helper untuk mengambil ID dari middleware
func getUserIDFromLocals(c *fiber.Ctx) (uuid.UUID, error) {
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		// Mengembalikan error yang akan dikirim oleh global error handler
		return uuid.Nil, utils.SendSimpleError(c, fiber.StatusUnauthorized, "Invalid token", "Invalid user ID in token")
	}
	return userID, nil
}

// getSessionIDFromLocals adalah helper untuk mengambil session ID dari middleware
func getSessionIDFromLocals(c *fiber.Ctx) uuid.UUID {
	sessionID, ok := c.Locals("sessionID").(uuid.UUID)
	if !ok {
		sessionID = uuid.Nil // Fallback
	}
	return sessionID
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
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Books retrieved successfully", books)
}

// GetProgressToRestore menangani 'GET /app/books/:bookId/restore'
func (h *AppHandler) GetProgressToRestore(c *fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return err
	}
	sessionID := getSessionIDFromLocals(c)

	bookIdStr := c.Params("bookId")
	bookID, err := strconv.Atoi(bookIdStr)
	if err != nil || bookID <= 0 {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid book ID", err.Error())
	}

	restoreData, err := h.appService.GetProgressToRestore(c.Context(), userID, sessionID, uint(bookID))
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Progress restored successfully", restoreData)
}

// UpdatePageProgress menangani 'POST /app/progress/update'
func (h *AppHandler) UpdatePageProgress(c *fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return err
	}
	sessionID := getSessionIDFromLocals(c)

	var req domain.UpdateProgressRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	if err := h.appService.UpdatePageProgress(c.Context(), userID, sessionID, &req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	// Mengirim data 'nil' karena pesannya sudah ada di root response
	return utils.SendSuccess(c, fiber.StatusOK, "Progress updated", nil)
}

// CompleteBookProgress menangani 'POST /app/progress/complete'
func (h *AppHandler) CompleteBookProgress(c *fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return err
	}
	sessionID := getSessionIDFromLocals(c)

	var req domain.CompleteProgressRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	if err := h.appService.CompleteBookProgress(c.Context(), userID, sessionID, &req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Book completed and progress saved", nil)
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
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Games retrieved successfully", games)
}

// SubmitGameScore menangani 'POST /app/games/:gameId/score'
func (h *AppHandler) SubmitGameScore(c *fiber.Ctx) error {
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return err
	}
	sessionID := getSessionIDFromLocals(c)

	gameIdStr := c.Params("gameId")
	gameID, err := strconv.Atoi(gameIdStr)
	if err != nil || gameID <= 0 {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid game ID", err.Error())
	}

	var req domain.SubmitGameScoreRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	if err := h.appService.SubmitGameScore(c.Context(), userID, sessionID, uint(gameID), &req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Score updated successfully", nil)
}

// --- Activity Handler ---

// GetMyActivity menangani 'GET /app/activity'
func (h *AppHandler) GetMyActivity(c *fiber.Ctx) error {
	// 1. Ambil 'userID' dari token
	userID, err := getUserIDFromLocals(c)
	if err != nil {
		return err // Error sudah dikirim oleh helper
	}

	// 2. Ambil filter paginasi dari query params
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	filters := &domain.ActivityLogQuery{
		Page:      page,
		Limit:     limit,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// 3. Panggil Usecase (dari dashboardService)
	paginatedData, err := h.dashboardService.GetStudentActivity(c.Context(), userID, filters)
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

	return utils.SendPagination(c, "My activity retrieved successfully", paginatedData.List, meta)
}
