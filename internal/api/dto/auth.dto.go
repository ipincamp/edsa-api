package dto

import "github.com/go-playground/validator/v10"

type RegisterRequest struct {
	Name                 string `json:"name" validate:"required,min=3,max=100"`
	Email                string `json:"email" validate:"required,email,max=255"`
	Password             string `json:"password" validate:"required,min=8,max=32"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func ValidateStruct(s interface{}) error {
	validate := validator.New()
	return validate.Struct(s)
}
