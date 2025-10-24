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

// getIDParam adalah helper untuk mengambil parameter ID dari URL dan mengonversinya ke uint
func (h *AdminHandler) getIDParam(c *fiber.Ctx) (uint, error) {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id parameter")
	}
	return uint(id), nil
}

// getUintIDParam adalah helper untuk mengambil parameter ID dari URL dan mengonversinya ke uint
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
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	subject, err := h.adminService.CreateSubject(c.Context(), &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, "Subject created successfully", subject)
}

func (h *AdminHandler) GetAllSubjects(c *fiber.Ctx) error {
	subjects, err := h.adminService.GetAllSubjects(c.Context())
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Subjects retrieved successfully", subjects)
}

func (h *AdminHandler) GetSubjectByID(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	subject, err := h.adminService.GetSubjectByID(c.Context(), id)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Subject retrieved successfully", subject)
}

func (h *AdminHandler) UpdateSubject(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	var req domain.UpdateSubjectRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	subject, err := h.adminService.UpdateSubject(c.Context(), id, &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Subject updated successfully", subject)
}

func (h *AdminHandler) DeleteSubject(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	if err := h.adminService.DeleteSubject(c.Context(), id); err != nil {
		// Bisa jadi error 404 (not found) atau 409 (conflict jika ada relasi)
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Subject deleted successfully", nil)
}

// --- Class Handlers ---

func (h *AdminHandler) CreateClass(c *fiber.Ctx) error {
	var req domain.CreateClassRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	class, err := h.adminService.CreateClass(c.Context(), &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error()) // 400 jika subject_id tidak ada
	}
	return utils.SendSuccess(c, fiber.StatusCreated, "Class created successfully", class)
}

func (h *AdminHandler) GetAllClasses(c *fiber.Ctx) error {
	classes, err := h.adminService.GetAllClasses(c.Context())
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Classes retrieved successfully", classes)
}

func (h *AdminHandler) GetClassByID(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	class, err := h.adminService.GetClassByID(c.Context(), id)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Class retrieved successfully", class)
}

func (h *AdminHandler) UpdateClass(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	var req domain.UpdateClassRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	class, err := h.adminService.UpdateClass(c.Context(), id, &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error()) // 404 jika class/subject tidak ada
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Class updated successfully", class)
}

func (h *AdminHandler) DeleteClass(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	if err := h.adminService.DeleteClass(c.Context(), id); err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Class deleted successfully", nil)
}

// --- Group Handlers ---

func (h *AdminHandler) CreateGroup(c *fiber.Ctx) error {
	var req domain.CreateGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	group, err := h.adminService.CreateGroup(c.Context(), &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error()) // 400 jika class_id tidak ada
	}
	return utils.SendSuccess(c, fiber.StatusCreated, "Group created successfully", group)
}

func (h *AdminHandler) GetAllGroups(c *fiber.Ctx) error {
	groups, err := h.adminService.GetAllGroups(c.Context())
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Groups retrieved successfully", groups)
}

func (h *AdminHandler) GetGroupByID(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	group, err := h.adminService.GetGroupByID(c.Context(), id)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Group retrieved successfully", group)
}

func (h *AdminHandler) UpdateGroup(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	var req domain.UpdateGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	group, err := h.adminService.UpdateGroup(c.Context(), id, &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error()) // 404 jika group/class tidak ada
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Group updated successfully", group)
}

func (h *AdminHandler) DeleteGroup(c *fiber.Ctx) error {
	id, err := h.getIDParam(c)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	if err := h.adminService.DeleteGroup(c.Context(), id); err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Group deleted successfully", nil)
}

// --- Book Handlers ---

func (h *AdminHandler) CreateBook(c *fiber.Ctx) error {
	var req domain.CreateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	book, err := h.adminService.CreateBook(c.Context(), &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, "Book created successfully", book)
}

func (h *AdminHandler) GetAllBooks(c *fiber.Ctx) error {
	// 1. Ambil filter paginasi dari query params
	// Konversi "0" jika tidak ada, agar bisa diabaikan oleh service/repo
	id, _ := strconv.Atoi(c.Query("id", "0"))
	bookOrder, _ := strconv.Atoi(c.Query("book_order", "0"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date") // Menambahkan end_date

	filters := &domain.BookQuery{
		ID:        uint(id),
		BookOrder: bookOrder,
		Page:      page,
		Limit:     limit,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// 2. Panggil Usecase
	paginatedData, err := h.adminService.GetAllBooks(c.Context(), filters)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}

	// 3. Kembalikan menggunakan format SendPagination
	meta := utils.PaginationMeta{
		Page:      paginatedData.Meta.Page,
		Limit:     paginatedData.Meta.Limit,
		TotalPage: paginatedData.Meta.TotalPage,
		TotalData: paginatedData.Meta.TotalData,
	}

	return utils.SendPagination(c, "Books retrieved successfully", paginatedData.List, meta)
}

func (h *AdminHandler) GetBookByID(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "bookId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	book, err := h.adminService.GetBookByID(c.Context(), id)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Book retrieved successfully", book)
}

func (h *AdminHandler) UpdateBook(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "bookId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	var req domain.UpdateBookRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	book, err := h.adminService.UpdateBook(c.Context(), id, &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Book updated successfully", book)
}

func (h *AdminHandler) DeleteBook(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "bookId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	if err := h.adminService.DeleteBook(c.Context(), id); err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Book deleted successfully", nil)
}

// --- Page Handlers ---

func (h *AdminHandler) CreatePage(c *fiber.Ctx) error {
	bookID, err := h.getUintIDParam(c, "bookId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	var req domain.CreatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	page, err := h.adminService.CreatePage(c.Context(), bookID, &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error()) // 400 jika bookId tidak ada
	}
	return utils.SendSuccess(c, fiber.StatusCreated, "Page created successfully", page)
}

func (h *AdminHandler) GetAllPagesForBook(c *fiber.Ctx) error {
	bookID, err := h.getUintIDParam(c, "bookId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	pages, err := h.adminService.GetAllPagesForBook(c.Context(), bookID)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Pages retrieved successfully", pages)
}

func (h *AdminHandler) GetPageByID(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "pageId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	page, err := h.adminService.GetPageByID(c.Context(), id)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Page retrieved successfully", page)
}

func (h *AdminHandler) UpdatePage(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "pageId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	var req domain.UpdatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	page, err := h.adminService.UpdatePage(c.Context(), id, &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Page updated successfully", page)
}

func (h *AdminHandler) DeletePage(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "pageId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	if err := h.adminService.DeletePage(c.Context(), id); err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Page deleted successfully", nil)
}

// --- Interaction Handlers ---

func (h *AdminHandler) CreateInteraction(c *fiber.Ctx) error {
	pageID, err := h.getUintIDParam(c, "pageId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	var req domain.CreateInteractionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	interaction, err := h.adminService.CreateInteraction(c.Context(), pageID, &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error()) // 400 jika pageId tidak ada
	}
	return utils.SendSuccess(c, fiber.StatusCreated, "Interaction created successfully", interaction)
}

func (h *AdminHandler) GetAllInteractionsForPage(c *fiber.Ctx) error {
	pageID, err := h.getUintIDParam(c, "pageId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	interactions, err := h.adminService.GetAllInteractionsForPage(c.Context(), pageID)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Interactions retrieved successfully", interactions)
}

func (h *AdminHandler) GetInteractionByID(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "interactionId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	interaction, err := h.adminService.GetInteractionByID(c.Context(), id)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Interaction retrieved successfully", interaction)
}

func (h *AdminHandler) UpdateInteraction(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "interactionId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	var req domain.UpdateInteractionRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, "Invalid request body", err.Error())
	}
	if errs := h.validate.ValidateStruct(req); len(errs) > 0 {
		return utils.SendValidationErrors(c, errs)
	}

	interaction, err := h.adminService.UpdateInteraction(c.Context(), id, &req)
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusNotFound, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Interaction updated successfully", interaction)
}

func (h *AdminHandler) DeleteInteraction(c *fiber.Ctx) error {
	id, err := h.getUintIDParam(c, "interactionId")
	if err != nil {
		return utils.SendSimpleError(c, fiber.StatusBadRequest, err.Error(), err.Error())
	}

	if err := h.adminService.DeleteInteraction(c.Context(), id); err != nil {
		return utils.SendSimpleError(c, fiber.StatusInternalServerError, err.Error(), err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, "Interaction deleted successfully", nil)
}
