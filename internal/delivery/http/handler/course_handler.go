package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/delivery/http/dto"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/service"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"github.com/ipincamp/go-edsa-api/pkg/validator"
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
