package domain

import (
	"time"

	"gorm.io/gorm"
)

// Subject adalah entitas domain inti untuk mata pelajaran
type Subject struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

// --- Data Transfer Objects (DTOs) ---

// SubjectResponse adalah DTO untuk respons
type SubjectResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// CreateSubjectRequest adalah DTO untuk membuat subject baru
type CreateSubjectRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

// UpdateSubjectRequest adalah DTO untuk memperbarui subject
type UpdateSubjectRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}
