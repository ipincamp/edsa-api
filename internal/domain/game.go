package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Game adalah entitas domain inti untuk game
type Game struct {
	ID               uint
	Name             string
	Type             string
	RelatedBookTheme string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

// UserGameScore adalah entitas domain inti untuk skor game user
type UserGameScore struct {
	ID           uint
	UserID       uuid.UUID
	GameID       uint
	HighestScore int
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt

	User User
	Game Game
}

// --- Data Transfer Objects (DTOs) ---

// GameResponse adalah DTO untuk 'GET /games'
type GameResponse struct {
	ID               uint   `json:"id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	RelatedBookTheme string `json:"related_book_theme"`
	HighestScore     int    `json:"highest_score"` // Skor tertinggi user ini
}

// SubmitGameScoreRequest adalah DTO untuk 'POST /games/:gameId/score'
type SubmitGameScoreRequest struct {
	Score int `json:"score" validate:"required,number,min=0,max=100"`
}
