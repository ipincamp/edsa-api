package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func New() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

func (cv *CustomValidator) Validate(data interface{}) map[string]string {
	err := cv.validator.Struct(data)
	if err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			field := strings.ToLower(err.Field())
			errors[field] = formatErrorMessage(err)
		}
		return errors
	}
	return nil
}

func formatErrorMessage(err validator.FieldError) string {
	field := err.Field()
	field = strings.ToLower(field)
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("The '%s' field is required.", field)
	case "email":
		return fmt.Sprintf("The '%s' must be a valid email address.", field)
	case "min":
		return fmt.Sprintf("The '%s' must be at least %s characters.", field, param)
	case "max":
		return fmt.Sprintf("The '%s' may not be greater than %s characters.", field, param)
	case "eqfield":
		return fmt.Sprintf("The '%s' must match the '%s' field.", field, param)
	default:
		return fmt.Sprintf("The '%s' is invalid.", field)
	}
}
