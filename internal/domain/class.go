package domain

import (
	"time"

	"gorm.io/gorm"
)

// Class adalah entitas domain inti untuk kelas
type Class struct {
	ID        uint
	Name      string
	SubjectID uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Relasi (opsional di domain, tapi bisa berguna)
	Subject Subject
}

// --- Data Transfer Objects (DTOs) ---

// ClassResponse adalah DTO untuk respons
type ClassResponse struct {
	ID        uint            `json:"id"`
	Name      string          `json:"name"`
	SubjectID uint            `json:"subject_id"`
	Subject   SubjectResponse `json:"subject,omitempty"` // Tampilkan detail subject
}

// CreateClassRequest adalah DTO untuk membuat class baru
type CreateClassRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=100"`
	SubjectID uint   `json:"subject_id" validate:"required,number"`
}

// UpdateClassRequest adalah DTO untuk memperbarui class
type UpdateClassRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=100"`
	SubjectID uint   `json:"subject_id" validate:"required,number"`
}
