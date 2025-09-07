package dto

// RegisterRequest adalah DTO untuk request registrasi user baru
type RegisterRequest struct {
	Name                 string `json:"name" validate:"required,min=3"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,gte=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,eqfield=Password"`
}

// LoginRequest adalah DTO untuk request login user
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest adalah DTO untuk request refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// AuthResponse adalah DTO untuk response autentikasi
type AuthResponse struct {
	Token Token        `json:"token"`
	User  UserResponse `json:"user"`
}

// Token adalah DTO untuk response token
type Token struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token"`
}
