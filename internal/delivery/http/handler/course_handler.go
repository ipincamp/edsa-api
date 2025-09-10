package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/dto"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/service"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"github.com/ipincamp/go-edsa-api/pkg/validator"
	"gorm.io/gorm"
)

type CourseHandler struct {
	courseService service.CourseService
	validator     *validator.CustomValidator
}

func NewCourseHandler(courseService service.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
		validator:     validator.NewValidator(),
	}
}

// ApplyToJoinGroup - Handler for guests to apply to a course group
func (h *CourseHandler) ApplyToJoinGroup(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)

	var req dto.ApplyToGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if errs := h.validator.Validate(req); errs != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", errs)
	}

	err := h.courseService.ApplyToJoinGroup(c.Context(), user.ID, req.GroupCode)
	if err != nil {
		if errors.Is(err, service.ErrGroupNotFound) {
			return util.SendError(c, fiber.StatusNotFound, err.Error())
		}
		// Anda bisa menambahkan error mapping lain di sini
		return util.SendError(c, fiber.StatusInternalServerError, "failed to apply to group")
	}

	return util.SendSuccess(c, fiber.StatusOK, "application submitted successfully", nil)
}

// HandleJoinRequest - Handler for teachers to approve/reject a request
func (h *CourseHandler) HandleJoinRequest(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)

	requestID := c.Params("requestID")
	var req dto.HandleJoinRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if errs := h.validator.Validate(req); errs != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", errs)
	}

	err := h.courseService.HandleJoinRequest(c.Context(), user.ID, requestID, req.Approved)
	if err != nil {
		// Tambahkan error mapping yang lebih spesifik jika perlu
		return util.SendError(c, fiber.StatusInternalServerError, "failed to handle join request")
	}

	message := "request approved successfully"
	if !req.Approved {
		message = "request rejected successfully"
	}

	return util.SendSuccess(c, fiber.StatusOK, message, nil)
}

// GetMyClasses - Handler for a teacher to get their classes
func (h *CourseHandler) GetMyClasses(c *fiber.Ctx) error {
	user := c.Locals("user").(*domain.User)

	classes, err := h.courseService.GetMyClasses(c.Context(), user.ID)
	if err != nil {
		return util.SendError(c, fiber.StatusInternalServerError, "failed to retrieve classes")
	}

	return util.SendSuccess(c, fiber.StatusOK, "classes retrieved successfully", classes)
}

// GetStudentsByClass - Handler to get students in a class
func (h *CourseHandler) GetStudentsByClass(c *fiber.Ctx) error {
	groupID := c.Params("groupID")

	students, err := h.courseService.GetStudentsByClass(c.Context(), groupID)
	if err != nil {
		return util.SendError(c, fiber.StatusInternalServerError, "failed to retrieve students")
	}

	return util.SendSuccess(c, fiber.StatusOK, "students retrieved successfully", students)
}

// GetTeachersByClass - Handler to get teachers in a class
func (h *CourseHandler) GetTeachersByClass(c *fiber.Ctx) error {
	groupID := c.Params("groupID")

	teachers, err := h.courseService.GetTeachersByClass(c.Context(), groupID)
	if err != nil {
		return util.SendError(c, fiber.StatusInternalServerError, "failed to retrieve teachers")
	}

	return util.SendSuccess(c, fiber.StatusOK, "teachers retrieved successfully", teachers)
}

// Admin

// CreateCourse - Handler to create a new course
func (h *CourseHandler) CreateCourse(c *fiber.Ctx) error {
	var req dto.CreateCourseRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if errs := h.validator.Validate(req); errs != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", errs)
	}

	course, err := h.courseService.CreateCourse(c.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrCourseNameExists) {
			return util.SendError(c, fiber.StatusConflict, err.Error())
		}
		return util.SendError(c, fiber.StatusInternalServerError, "failed to create course")
	}

	return util.SendSuccess(c, fiber.StatusCreated, "course created successfully", course)
}

// GetCourseByID - Handler to get a course by its ID
func (h *CourseHandler) GetCourseByID(c *fiber.Ctx) error {
	var courseIDReq dto.CourseIDRequest
	if err := c.ParamsParser(&courseIDReq); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid course ID", nil)
	}
	if validationErrors := h.validator.Validate(courseIDReq); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	course, err := h.courseService.GetCourseByID(c.Context(), courseIDReq.CourseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.SendError(c, fiber.StatusNotFound, "course not found")
		}
		return util.SendError(c, fiber.StatusInternalServerError, "failed to get course")
	}

	return util.SendSuccess(c, fiber.StatusOK, "course retrieved successfully", course)
}

// GetAllCourses - Handler to get all courses with pagination
func (h *CourseHandler) GetAllCourses(c *fiber.Ctx) error {
	var req dto.CourseFilterRequest
	if err := c.QueryParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid query parameters", nil)
	}
	if validationErrors := h.validator.Validate(req); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	courses, err := h.courseService.GetAllCourses(c.Context(), req.Page, req.Limit)
	if err != nil {
		return util.SendError(c, fiber.StatusInternalServerError, "failed to get courses")
	}
	return util.SendSuccess(c, fiber.StatusOK, "courses retrieved successfully", courses)
}

// UpdateCourse - Handler to update a course by its ID
func (h *CourseHandler) UpdateCourse(c *fiber.Ctx) error {
	var courseIDReq dto.CourseIDRequest
	if err := c.ParamsParser(&courseIDReq); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid course ID", nil)
	}
	if validationErrors := h.validator.Validate(courseIDReq); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	var req dto.UpdateCourseRequest
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if errs := h.validator.Validate(req); errs != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", errs)
	}

	course, err := h.courseService.UpdateCourse(c.Context(), courseIDReq.CourseId, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.SendError(c, fiber.StatusNotFound, "course not found")
		}
		return util.SendError(c, fiber.StatusInternalServerError, "failed to update course")
	}

	return util.SendSuccess(c, fiber.StatusOK, "course updated successfully", course)
}

// DeleteCourse - Handler to delete a course by its ID
func (h *CourseHandler) DeleteCourse(c *fiber.Ctx) error {
	var courseIDReq dto.CourseIDRequest
	if err := c.ParamsParser(&courseIDReq); err != nil {
		return util.SendError(c, fiber.StatusBadRequest, "invalid course ID", nil)
	}
	if validationErrors := h.validator.Validate(courseIDReq); validationErrors != nil {
		return util.SendError(c, fiber.StatusBadRequest, "validation failed", validationErrors)
	}

	err := h.courseService.DeleteCourse(c.Context(), courseIDReq.CourseId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.SendError(c, fiber.StatusNotFound, "course not found")
		}
		if errors.Is(err, service.ErrCourseAlreadyDeleted) {
			return util.SendError(c, fiber.StatusBadRequest, err.Error())
		}
		return util.SendError(c, fiber.StatusInternalServerError, "failed to delete course")
	}

	return util.SendSuccess(c, fiber.StatusOK, "course deleted successfully", nil)
}
