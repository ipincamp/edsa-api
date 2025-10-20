package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User adalah entitas domain inti.
type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	RoleID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// --- Data Transfer Objects (DTOs) ---
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UserResponse adalah DTO untuk data pengguna yang aman dikirim ke klien.
type UserResponse struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	RoleID uint      `json:"role_id"`
}

// TokenResponse adalah DTO untuk token autentikasi.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// AuthResponse adalah DTO untuk respons login/register.
type AuthResponse struct {
	User  UserResponse  `json:"user"`
	Token TokenResponse `json:"token"`
}
