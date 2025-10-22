package domain

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// Interaction adalah entitas domain inti untuk interaksi halaman
type Interaction struct {
	ID        uint
	PageID    uint
	Type      string
	Config    json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	Page Page // Relasi
}

// --- Data Transfer Objects (DTOs) ---

// InteractionResponse adalah DTO untuk respons
type InteractionResponse struct {
	ID     uint            `json:"id"`
	PageID uint            `json:"page_id"`
	Type   string          `json:"type"`
	Config json.RawMessage `json:"config"`
	Page   PageResponse    `json:"page,omitempty"`
}

// CreateInteractionRequest adalah DTO untuk membuat interaksi baru
type CreateInteractionRequest struct {
	Type   string          `json:"type" validate:"required,max=50"`
	Config json.RawMessage `json:"config" validate:"required"`
}

// UpdateInteractionRequest adalah DTO untuk memperbarui interaksi
type UpdateInteractionRequest struct {
	Type   string          `json:"type" validate:"required,max=50"`
	Config json.RawMessage `json:"config" validate:"required"`
}
