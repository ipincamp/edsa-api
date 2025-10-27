package domain

import (
	"time"
)

// Konstanta Tipe Token
const (
	TokenTypeAccess            = "access"         // Token akses
	TokenTypeRefresh           = "refresh"        // Token refresh
	TokenTypeEmailVerification = "email_verify"   // Token verifikasi email
	TokenTypePasswordReset     = "password_reset" // Token reset password
	TokenTypeAccountDeletion   = "account_delete" // Token konfirmasi hapus akun
)

// PasetoPayload adalah struktur data yang disimpan dalam token
type PasetoPayload struct {
	UserID                string    `json:"user_id"`
	Email                 string    `json:"email"`
	SessionID             string    `json:"session_id"`
	RoleName              string    `json:"role_name"`
	TokenType             string    `json:"token_type"`
	VerificationAttemptID string    `json:"verification_attempt_id,omitempty"`
	IssuedAt              time.Time `json:"iat"`
	ExpiresAt             time.Time `json:"exp"`
}

// TokenResponse adalah DTO untuk token autentikasi.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
