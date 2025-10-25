package domain

import (
	"time"

	"gorm.io/gorm"
)

// Page adalah entitas domain inti untuk halaman buku
type Page struct {
	ID                uint
	BookID            uint
	PageNumber        int
	InstructionText   string
	Narration_ID      string
	Narration_EN      string
	AudioNarrationURL string
	IsPostActivity    bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt

	Book Book // Relasi
}

// --- Data Transfer Objects (DTOs) ---

// PageResponse adalah DTO untuk respons
type PageResponse struct {
	ID                uint         `json:"id"`
	BookID            uint         `json:"book_id"`
	PageNumber        int          `json:"page_number"`
	InstructionText   string       `json:"instruction_text"`
	Narration_ID      string       `json:"narration_id"`
	Narration_EN      string       `json:"narration_en"`
	AudioNarrationURL string       `json:"audio_narration_url"`
	IsPostActivity    bool         `json:"is_post_activity"`
	Book              BookResponse `json:"book,omitempty"`
}

// CreatePageRequest adalah DTO untuk membuat halaman baru
type CreatePageRequest struct {
	PageNumber        int    `json:"page_number" validate:"required,number,min=1"`
	InstructionText   string `json:"instruction_text"`
	Narration_ID      string `json:"narration_id"`
	Narration_EN      string `json:"narration_en"`
	AudioNarrationURL string `json:"audio_narration_url,omitempty" validate:"omitempty,url"`
	IsPostActivity    bool   `json:"is_post_activity,omitempty"`
}

// UpdatePageRequest adalah DTO untuk memperbarui halaman
type UpdatePageRequest struct {
	PageNumber        int    `json:"page_number" validate:"required,number,min=1"`
	InstructionText   string `json:"instruction_text"`
	Narration_ID      string `json:"narration_id"`
	Narration_EN      string `json:"narration_en"`
	AudioNarrationURL string `json:"audio_narration_url,omitempty" validate:"omitempty,url"`
	IsPostActivity    bool   `json:"is_post_activity,omitempty"`
}
