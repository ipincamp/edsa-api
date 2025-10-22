package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/pkg/validator"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type MediaHandler struct {
	mediaService usecase.MediaService
	validate     *validator.GoPlaygroundValidator
}

func NewMediaHandler(
	ms usecase.MediaService,
	v *validator.GoPlaygroundValidator,
) *MediaHandler {
	return &MediaHandler{
		mediaService: ms,
		validate:     v,
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

	// 3. Panggil usecase
	asset, err := h.mediaService.UploadFile(c.Context(), file, ownerID, ownerType)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, "File uploaded successfully", asset)
}
