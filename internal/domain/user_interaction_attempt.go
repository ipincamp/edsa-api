package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserInteractionAttempt adalah entitas domain inti
type UserInteractionAttempt struct {
	ID              uint
	UserID          uuid.UUID
	InteractionID   uint
	Timestamp       time.Time
	UserAnswer      json.RawMessage
	IsCorrect       bool
	ScoreAwarded    float64
	DurationSeconds int
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt

	// Relasi
	User        User
	Interaction Interaction
}

// --- Data Transfer Objects (DTOs) ---
