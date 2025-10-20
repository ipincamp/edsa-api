package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Konstanta untuk Aksi Log
const (
	ActionLogin        = "LOGIN"
	ActionStartBook    = "START_BOOK"
	ActionCompleteBook = "COMPLETE_BOOK"
	ActionCompleteGame = "COMPLETE_GAME"
)

// ActivityLog adalah entitas domain inti untuk log aktivitas pengguna
type ActivityLog struct {
	ID             uint
	UserID         uuid.UUID
	SessionID      uuid.UUID
	Action         string
	TimestampStart time.Time
	DurationMs     *int
	Details        json.RawMessage
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt

	User User
}

// --- Data Transfer Objects (DTOs) ---

// ActivityLogResponse adalah DTO untuk 'GET /students/:studentId/activity'
type ActivityLogResponse struct {
	ID             uint            `json:"id"`
	Action         string          `json:"action"`
	TimestampStart time.Time       `json:"timestamp_start"`
	DurationMs     *int            `json:"duration_ms,omitempty"`
	Details        json.RawMessage `json:"details,omitempty"`
}
