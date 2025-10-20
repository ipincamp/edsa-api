package domain

import (
	"time"

	"gorm.io/gorm"
)

// Page adalah entitas domain inti untuk halaman buku
type Page struct {
	ID              uint
	BookID          uint
	PageNumber      int
	NarrativeText   string
	InstructionText string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt

	Book Book // Relasi
}

// --- Data Transfer Objects (DTOs) ---

// PageResponse adalah DTO untuk respons
type PageResponse struct {
	ID              uint         `json:"id"`
	BookID          uint         `json:"book_id"`
	PageNumber      int          `json:"page_number"`
	NarrativeText   string       `json:"narrative_text"`
	InstructionText string       `json:"instruction_text"`
	Book            BookResponse `json:"book,omitempty"`
}

// CreatePageRequest adalah DTO untuk membuat halaman baru
type CreatePageRequest struct {
	PageNumber      int    `json:"page_number" validate:"required,number,min=1"`
	NarrativeText   string `json:"narrative_text"`
	InstructionText string `json:"instruction_text"`
}

// UpdatePageRequest adalah DTO untuk memperbarui halaman
type UpdatePageRequest struct {
	PageNumber      int    `json:"page_number" validate:"required,number,min=1"`
	NarrativeText   string `json:"narrative_text"`
	InstructionText string `json:"instruction_text"`
}
