package util

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type GoValidator struct {
	validate *validator.Validate
}

func NewValidator() *GoValidator {
	return &GoValidator{
		validate: validator.New(),
	}
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v *GoValidator) Validate(payload interface{}) []ValidationError {
	err := v.validate.Struct(payload)
	if err == nil {
		return nil
	}

	var validationErrors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		field := toSnakeCase(err.Field())
		validationErrors = append(validationErrors, ValidationError{
			Field:   field,
			Message: generateErrorMessage(err),
		})
	}
	return validationErrors
}

func generateErrorMessage(err validator.FieldError) string {
	field := toSnakeCase(err.Field())
	tag := err.Tag()
	param := err.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("Field '%s' is required", field)
	case "email":
		return fmt.Sprintf("Field '%s' must be a valid email address", field)
	case "min":
		return fmt.Sprintf("Field '%s' must be at least %s characters long", field, param)
	case "gte":
		return fmt.Sprintf("Field '%s' must be at least %s characters long", field, param)
	case "eqfield":
		return fmt.Sprintf("Field '%s' must be the same as field '%s'", field, toSnakeCase(param))
	default:
		return fmt.Sprintf("Field '%s' is not valid", field)
	}
}

func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
