package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User adalah entitas domain inti untuk pengguna
type User struct {
	ID               uuid.UUID
	Name             string
	Email            string
	Password         string
	RoleID           uint
	Role             Role
	ProfilePictureID *uuid.UUID
	ProfilePicture   MediaAsset
	EmailVerifiedAt  *time.Time
	IsActive         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

// --- Data Transfer Objects (DTOs) ---

// RegisterRequest adalah DTO untuk pendaftaran pengguna baru.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest adalah DTO untuk autentikasi pengguna.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// RefreshTokenRequest adalah DTO for token refresh.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ChangePasswordRequest adalah DTO untuk mengubah password.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

// UpdateDetailsRequest adalah DTO untuk mengubah nama (memerlukan konfirmasi password).
type UpdateDetailsRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Password string `json:"password" validate:"required"` // Untuk konfirmasi
}

// DeleteAccountRequest adalah DTO untuk penghapusan akun.
// DEPRECATED: Gunakan DeletionRequest dan ConfirmDeletionRequest
type DeleteAccountRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	DeletionReason  string `json:"deletion_reason" validate:"omitempty,max=255"`
}

// DeletionRequest adalah DTO untuk handler hapus akun (langkah 1)
type DeletionRequest struct {
	ConfirmationToken string `json:"confirmation_token,omitempty"`
	DeletionReason    string `json:"deletion_reason,omitempty"`
	CurrentPassword   string `json:"current_password,omitempty"`
}

// ConfirmDeletionRequest adalah DTO untuk validasi hapus akun (langkah 2)
type ConfirmDeletionRequest struct {
	ConfirmationToken string `json:"confirmation_token" validate:"required"`
	DeletionReason    string `json:"deletion_reason" validate:"required,min=10,max=500"`
	CurrentPassword   string `json:"current_password" validate:"required"`
}

// VerifyEmailRequest adalah DTO untuk 'POST /auth/verify-email'
type VerifyEmailRequest struct {
	Token string `json:"token" validate:"required"`
}

// ResendVerificationRequest adalah DTO untuk 'POST /auth/resend-verification'
type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ForgotPasswordRequest adalah DTO untuk 'POST /auth/forgot-password'
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// ResetPasswordRequest adalah DTO untuk 'POST /auth/reset-password'
type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// UpdateAvatarRequest adalah DTO untuk 'PATCH /users/me/avatar'
type UpdateAvatarRequest struct {
	MediaID uuid.UUID `json:"media_id" validate:"required,uuid"`
}

// UserResponse adalah DTO untuk data pengguna yang aman dikirim ke klien.
type UserResponse struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	RoleName          string    `json:"role"`
	JoinedAt          time.Time `json:"joined_at"`
	ProfilePictureURL string    `json:"profile_picture_url,omitempty"`
	OverallScore      int       `json:"overall_score,omitempty"`
	IsActive          bool      `json:"is_active,omitempty"`
	EmailVerified     bool      `json:"email_verified"`
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
