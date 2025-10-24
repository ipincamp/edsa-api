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
	ActionLogout       = "LOGOUT"
	ActionRegister     = "REGISTER"
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

// ActivityLogQuery adalah DTO untuk filter paginasi log
type ActivityLogQuery struct {
	Page      int
	Limit     int
	StartDate string // ISO 8601 Format (YYYY-MM-DD)
	EndDate   string // ISO 8601 Format (YYYY-MM-DD)
}

// PaginatedActivityLogs adalah struct internal untuk membawa hasil dari repo
type PaginatedActivityLogs struct {
	Logs      []ActivityLog
	TotalData int64
}

// PaginatedActivityLogDTO adalah DTO untuk respons paginasi
// Strukturnya mencerminkan utils.PaginationData
type PaginatedActivityLogDTO struct {
	List interface{}       `json:"list"`
	Meta PaginationMetaDTO `json:"meta"`
}

// PaginationMetaDTO adalah DTO untuk metadata paginasi
// Strukturnya mencerminkan utils.PaginationMeta
type PaginationMetaDTO struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalPage int64 `json:"total_page"`
	TotalData int64 `json:"total_data"`
}
