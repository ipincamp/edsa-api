package http

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/pkg/utils"
	"github.com/ipincamp/go-edsa-api/internal/pkg/validator"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type AdminHandler struct {
	adminService usecase.AdminService
	validate     *validator.GoPlaygroundValidator
}

func NewAdminHandler(as usecase.AdminService, v *validator.GoPlaygroundValidator) *AdminHandler {
	return &AdminHandler{
		adminService: as,
		validate:     v,
	}
}

// --- Helper ---
func (h *AdminHandler) getIDParam(c *fiber.Ctx) (uint, error) {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id parameter")
	}
	return uint(id), nil
}

func (h *AdminHandler) getUintIDParam(c *fiber.Ctx, paramName string) (uint, error) {
	idStr := c.Params(paramName)
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id parameter: " + paramName)
	}
	return uint(id), nil
}

// --- Subject Handlers ---

func (h *AdminHandler) CreateSubject(c *fiber.Ctx) error {
	var req domain.CreateSubjectRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	subject, err := h.adminService.CreateSubject(c.Context(), &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, subject)
}

func (h *AdminHandler) GetAllSubjects(c *fiber.Ctx) error {
	subjects, err := h.adminService.GetAllSubjects(c.Context())
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, subjects)
}

func (h *AdminHandler) GetSubjectByID(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	subject, err := h.adminService.GetSubjectByID(c.Context(), id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, subject)
}

func (h *AdminHandler) UpdateSubject(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req domain.UpdateSubjectRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	subject, err := h.adminService.UpdateSubject(c.Context(), id, &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, subject)
}

func (h *AdminHandler) DeleteSubject(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.adminService.DeleteSubject(c.Context(), id); err != nil {
		// Bisa jadi error 404 (not found) atau 409 (conflict jika ada relasi)
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// --- Class Handlers ---

func (h *AdminHandler) CreateClass(c *fiber.Ctx) error {
	var req domain.CreateClassRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	class, err := h.adminService.CreateClass(c.Context(), &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error()) // 400 jika subject_id tidak ada
	}
	return utils.SendSuccess(c, fiber.StatusCreated, class)
}

func (h *AdminHandler) GetAllClasses(c *fiber.Ctx) error {
	classes, err := h.adminService.GetAllClasses(c.Context())
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, classes)
}

func (h *AdminHandler) GetClassByID(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	class, err := h.adminService.GetClassByID(c.Context(), id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, class)
}

func (h *AdminHandler) UpdateClass(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req domain.UpdateClassRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	class, err := h.adminService.UpdateClass(c.Context(), id, &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error()) // 404 jika class/subject tidak ada
	}
	return utils.SendSuccess(c, fiber.StatusOK, class)
}

func (h *AdminHandler) DeleteClass(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.adminService.DeleteClass(c.Context(), id); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// --- Group Handlers ---

func (h *AdminHandler) CreateGroup(c *fiber.Ctx) error {
	var req domain.CreateGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	group, err := h.adminService.CreateGroup(c.Context(), &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error()) // 400 jika class_id tidak ada
	}
	return utils.SendSuccess(c, fiber.StatusCreated, group)
}

func (h *AdminHandler) GetAllGroups(c *fiber.Ctx) error {
	groups, err := h.adminService.GetAllGroups(c.Context())
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, groups)
}

func (h *AdminHandler) GetGroupByID(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	group, err := h.adminService.GetGroupByID(c.Context(), id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, group)
}

func (h *AdminHandler) UpdateGroup(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req domain.UpdateGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	group, err := h.adminService.UpdateGroup(c.Context(), id, &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error()) // 404 jika group/class tidak ada
	}
	return utils.SendSuccess(c, fiber.StatusOK, group)
}

func (h *AdminHandler) DeleteGroup(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.adminService.DeleteGroup(c.Context(), id); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// --- Book Handlers ---

func (h *AdminHandler) CreateBook(c *fiber.Ctx) error {
	var req domain.CreateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	book, err := h.adminService.CreateBook(c.Context(), &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, book)
}

func (h *AdminHandler) GetAllBooks(c *fiber.Ctx) error {
	books, err := h.adminService.GetAllBooks(c.Context())
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, books)
}

func (h *AdminHandler) GetBookByID(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "bookId")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	book, err := h.adminService.GetBookByID(c.Context(), id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, book)
}

func (h *AdminHandler) UpdateBook(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "bookId")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req domain.UpdateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	book, err := h.adminService.UpdateBook(c.Context(), id, &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, book)
}

func (h *AdminHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "bookId")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.adminService.DeleteBook(c.Context(), id); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// --- Page Handlers ---

func (h *AdminHandler) CreatePage(c *fiber.Ctx) error {
	bookID, err := h.getUintIDParam(c, "bookId")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req domain.CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	page, err := h.adminService.CreatePage(c.Context(), bookID, &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error()) // 400 jika bookId tidak ada
	}
	return utils.SendSuccess(c, fiber.StatusCreated, page)
}

func (h *AdminHandler) GetAllPagesForBook(c *fiber.Ctx) error {
	bookID, err := h.getUintIDParam(c, "bookId")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	pages, err := h.adminService.GetAllPagesForBook(c.Context(), bookID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, pages)
}

func (h *AdminHandler) GetPageByID(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "pageId")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	page, err := h.adminService.GetPageByID(c.Context(), id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, page)
}

func (h *AdminHandler) UpdatePage(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "pageId")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req domain.UpdatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	page, err := h.adminService.UpdatePage(c.Context(), id, &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, page)
}

func (h *AdminHandler) DeletePage(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "pageId")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.adminService.DeletePage(c.Context(), id); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// --- Interaction Handlers ---

func (h *AdminHandler) CreateInteraction(c *fiber.Ctx) error {
	pageID, err := h.getUintIDParam(c, "pageId")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req domain.CreateInteractionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	interaction, err := h.adminService.CreateInteraction(c.Context(), pageID, &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error()) // 400 jika pageId tidak ada
	}
	return utils.SendSuccess(c, fiber.StatusCreated, interaction)
}

func (h *AdminHandler) GetAllInteractionsForPage(c *fiber.Ctx) error {
	pageID, err := h.getUintIDParam(c, "pageId")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	interactions, err := h.adminService.GetAllInteractionsForPage(c.Context(), pageID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, interactions)
}

func (h *AdminHandler) GetInteractionByID(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "interactionId")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	interaction, err := h.adminService.GetInteractionByID(c.Context(), id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, interaction)
}

func (h *AdminHandler) UpdateInteraction(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "interactionId")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	var req domain.UpdateInteractionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	interaction, err := h.adminService.UpdateInteraction(c.Context(), id, &req)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, interaction)
}

func (h *AdminHandler) DeleteInteraction(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "interactionId")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	if err := h.adminService.DeleteInteraction(c.Context(), id); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(fiber.StatusNoContent)
}
