package dto

type RegisterRequest struct {
	Name                 string `json:"name" validate:"required,min=3"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,gte=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string   `json:"token"`
	User  UserData `json:"user"`
}
