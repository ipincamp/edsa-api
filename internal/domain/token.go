package domain

import (
	"time"
)

// PasetoPayload adalah struktur data yang disimpan dalam token
type PasetoPayload struct {
	UserID                string    `json:"user_id"`
	Email                 string    `json:"email"`
	SessionID             string    `json:"session_id"`
	RoleName              string    `json:"role_name"`
	VerificationAttemptID string    `json:"verification_attempt_id,omitempty"` // ID unik untuk link verifikasi
	IssuedAt              time.Time `json:"iat"`
	ExpiresAt             time.Time `json:"exp"`
}

// TokenResponse adalah DTO untuk token autentikasi.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
