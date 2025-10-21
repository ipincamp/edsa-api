package utils

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ipincamp/go-edsa-api/internal/pkg/validator"
)

// ErrorDetail struct for the "error" array
type ErrorDetail struct {
	Field *string `json:"field"`
	Tag   *string `json:"tag,omitempty"`
	Value *string `json:"value"` // 'value' adalah detail error utama
}

// SuccessResponse struct for successful responses
type SuccessResponse struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse struct for error responses
type ErrorResponse struct {
	Status  bool          `json:"status"`
	Message string        `json:"message"`
	Error   []ErrorDetail `json:"error,omitempty"`
}

// PaginationMeta struct for pagination
type PaginationMeta struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalPage int64 `json:"total_page"`
	TotalData int64 `json:"total_data"`
}

// PaginationData struct for paginated data
type PaginationData struct {
	List interface{}    `json:"list"`
	Meta PaginationMeta `json:"meta"`
}

// StringPtr adalah helper untuk membuat pointer ke string
func StringPtr(s string) *string {
	if s == "" {
		// Jangan tampilkan field jika string kosong
		return nil
	}
	return &s
}

// SendSuccess mengirim respons sukses terstruktur
// 'message' bisa dikosongkan untuk menggunakan "Success"
func SendSuccess(c *fiber.Ctx, code int, message string, data interface{}) error {
	if message == "" {
		message = "Success"
	}
	return c.Status(code).JSON(SuccessResponse{
		Status:  true,
		Message: message,
		Data:    data,
	})
}

// SendPagination mengirim respons sukses untuk data paginasi
func SendPagination(c *fiber.Ctx, message string, listData interface{}, meta PaginationMeta) error {
	if message == "" {
		message = "Success"
	}
	return c.Status(fiber.StatusOK).JSON(SuccessResponse{
		Status:  true,
		Message: message,
		Data: PaginationData{
			List: listData,
			Meta: meta,
		},
	})
}

// SendError mengirim respons error terstruktur
func SendError(c *fiber.Ctx, code int, message string, details ...ErrorDetail) error {
	if message == "" {
		message = "Failed"
	}
	return c.Status(code).JSON(ErrorResponse{
		Status:  false,
		Message: message,
		Error:   details,
	})
}

// SendSimpleError adalah helper untuk error tunggal yang tidak spesifik (non-validasi)
// 'errorValue' adalah detail yang akan dimasukkan ke field 'value'
func SendSimpleError(c *fiber.Ctx, code int, message string, errorValue string) error {
	if message == "" {
		message = "Failed"
	}
	return SendError(c, code, message, ErrorDetail{
		Field: nil, // Tidak ada field spesifik
		Tag:   nil,
		Value: StringPtr(errorValue),
	})
}

// SendValidationErrors mengonversi error validasi ke format ErrorDetail yang baru
func SendValidationErrors(c *fiber.Ctx, validationErrors []*validator.ValidationError) error {
	var details []ErrorDetail
	for _, err := range validationErrors {
		details = append(details, ErrorDetail{
			Field: StringPtr(err.Field),
			Tag:   StringPtr(err.Tag),
			Value: StringPtr(err.Value),
		})
	}
	return c.Status(fiber.StatusUnprocessableEntity).JSON(ErrorResponse{
		Status:  false,
		Message: "Validation failed",
		Error:   details,
	})
}
