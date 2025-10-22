package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type GoPlaygroundValidator struct {
	Validator *validator.Validate
}

func NewValidator() *GoPlaygroundValidator {
	return &GoPlaygroundValidator{
		Validator: validator.New(),
	}
}

// Struct kustom untuk error validasi
type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func (v *GoPlaygroundValidator) ValidateStruct(s interface{}) []*ValidationError {
	var errors []*ValidationError
	err := v.Validator.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var elem ValidationError
			elem.Field = err.Field()
			elem.Tag = err.Tag()
			elem.Value = err.Param()
			errors = append(errors, &elem)
		}
	}
	return errors
}

// Helper untuk parsing dan validasi dalam satu langkah
func ParseAndValidate(c *fiber.Ctx, s interface{}, v *GoPlaygroundValidator) error {
	if err := c.BodyParser(s); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if errs := v.ValidateStruct(s); len(errs) > 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Validation failed")
	}
	return nil
}
