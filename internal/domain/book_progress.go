package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Status progres
const (
	BookProgressStatusLocked    = "locked"
	BookProgressStatusUnlocked  = "unlocked"
	BookProgressStatusCompleted = "completed"
)

// UserBookProgress adalah entitas domain inti untuk melacak progres buku pengguna
type UserBookProgress struct {
	ID                   uint
	UserID               uuid.UUID
	BookID               uint
	Status               string
	HighestScore         float64
	LastPageID           uint
	CurrentSessionPoints int
	Rating               int
	CreatedAt            time.Time
	UpdatedAt            time.Time
	DeletedAt            gorm.DeletedAt

	User User
	Book Book
}

// --- Data Transfer Objects (DTOs) ---

// RestoreProgressResponse adalah DTO untuk 'GET /books/:bookId/restore'
type RestoreProgressResponse struct {
	LastPageID           uint `json:"last_page_id"`
	CurrentSessionPoints int  `json:"current_session_points"`
}

// UpdateProgressRequest adalah DTO untuk 'POST /progress/update'
type UpdateProgressRequest struct {
	BookID             uint `json:"book_id" validate:"required"`
	PageID             uint `json:"page_id" validate:"required"`
	PointsEarnedOnPage int  `json:"points_earned_on_page"` // Poin yang didapat di halaman ini
}

// CompleteProgressRequest adalah DTO untuk 'POST /progress/complete'
type CompleteProgressRequest struct {
	BookID     uint    `json:"book_id" validate:"required"`
	FinalScore float64 `json:"final_score" validate:"number,min=0,max=100"`
	DurationMs int     `json:"duration_ms"` // Durasi dalam milidetik
}
