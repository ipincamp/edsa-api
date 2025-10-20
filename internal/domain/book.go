package domain

import (
	"time"

	"gorm.io/gorm"
)

// Book adalah entitas domain inti untuk buku
type Book struct {
	ID            uint
	Title         string
	Description   string
	CoverImageURL string
	Theme         string
	BookOrder     int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
}

// --- Data Transfer Objects (DTOs) ---

// BookResponse adalah DTO untuk respons
type BookResponse struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	CoverImageURL string `json:"cover_image_url"`
	Theme         string `json:"theme"`
	BookOrder     int    `json:"book_order"`
	Status        string `json:"status,omitempty"`
	HighestScore  int    `json:"highest_score,omitempty"`
}

// CreateBookRequest adalah DTO untuk membuat buku baru
type CreateBookRequest struct {
	Title         string `json:"title" validate:"required,min=3,max=255"`
	Description   string `json:"description"`
	CoverImageURL string `json:"cover_image_url" validate:"omitempty,url"`
	Theme         string `json:"theme" validate:"omitempty,max=100"`
	BookOrder     int    `json:"book_order" validate:"number,min=1"`
}

// UpdateBookRequest adalah DTO untuk memperbarui buku
type UpdateBookRequest struct {
	Title         *string `json:"title,omitempty" validate:"omitempty,min=3,max=255"`
	Description   *string `json:"description,omitempty"`
	CoverImageURL *string `json:"cover_image_url,omitempty" validate:"omitempty,url"`
	Theme         *string `json:"theme,omitempty" validate:"omitempty,max=100"`
	BookOrder     *int    `json:"book_order,omitempty" validate:"omitempty,number,min=1"`
}
