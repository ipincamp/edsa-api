package http

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/pkg/validator"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type MediaHandler struct {
	mediaService  usecase.MediaService
	validate      *validator.GoPlaygroundValidator
	loggerService usecase.ActivityLoggerService
}

func NewMediaHandler(
	ms usecase.MediaService,
	v *validator.GoPlaygroundValidator,
	logger usecase.ActivityLoggerService,
) *MediaHandler {
	return &MediaHandler{
		mediaService:  ms,
		validate:      v,
		loggerService: logger,
	}
}

// UploadFile menangani 'POST /media/upload'
func (h *MediaHandler) UploadFile(c *fiber.Ctx) error {
	// 1. Ambil file dari form
	file, err := c.FormFile("file")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Missing file", err.Error())
	}

	// 2. Ambil data polimorfik
	ownerID := c.FormValue("owner_id")
	ownerType := c.FormValue("owner_type") // Cth: "book_cover", "interaction_audio"

	if ownerID == "" {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Missing owner_id", "Field 'owner_id' is required")
	}

	if ownerType == "" {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Missing owner_type", "Field 'owner_type' is required")
	}

	uploaderID, ok := c.Locals("userID").(uuid.UUID)
	var uploaderIDPtr *uuid.UUID
	if ok {
		uploaderIDPtr = &uploaderID
	}

	sessionID, _ := c.Locals("sessionID").(uuid.UUID)

	// 3. Panggil usecase
	asset, err := h.mediaService.UploadFile(c.Context(), file, ownerID, ownerType, uploaderIDPtr)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	details, _ := json.Marshal(map[string]interface{}{
		"asset_id":   asset.ID,
		"file_name":  asset.FileName,
		"owner_id":   ownerID,
		"owner_type": ownerType,
	})
	h.loggerService.Log(c.Context(), domain.ActivityLog{
		UserID:    uploaderID, // Gunakan UUID asli, bukan pointer
		SessionID: sessionID,
		Action:    domain.ActionMediaUpload,
		Details:   details,
	})

	return utils.SendSuccess(c, fiber.StatusCreated, "File uploaded successfully", asset)
}

// DeleteFile menangani 'DELETE /media/:id'
func (h *MediaHandler) DeleteFile(c *fiber.Ctx) error {
	// 1. Ambil asset ID dari URL
	assetIdStr := c.Params("id")
	assetID, err := uuid.Parse(assetIdStr)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid asset ID format", err.Error())
	}

	// 2. Ambil ID penghapus (deleter) dan sesi dari token
	deleterID, ok := c.Locals("userID").(uuid.UUID)
	var deleterIDPtr *uuid.UUID
	if ok {
		deleterIDPtr = &deleterID
	}

	sessionID, _ := c.Locals("sessionID").(uuid.UUID)

	// 3. Panggil usecase
	if err := h.mediaService.DeleteFile(c.Context(), assetID, deleterIDPtr); err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}

	// 4. Log aktivitas (Solusi 2)
	details, _ := json.Marshal(map[string]interface{}{
		"asset_id": assetID,
	})

	h.loggerService.Log(c.Context(), domain.ActivityLog{
		UserID:    deleterID, // Gunakan UUID asli, bukan pointer
		SessionID: sessionID,
		Action:    domain.ActionMediaDelete,
		Details:   details,
	})

	return utils.SendSuccess(c, fiber.StatusOK, "File deleted successfully", nil)
}
