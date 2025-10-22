package domain

import (
	"time"

	"gorm.io/gorm"
)

// Group adalah entitas domain inti untuk grup
type Group struct {
	ID        uint
	Name      string
	ClassID   uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt

	// Relasi
	Class Class
}

// --- Data Transfer Objects (DTOs) ---

// GroupResponse adalah DTO untuk respons
type GroupResponse struct {
	ID      uint          `json:"id"`
	Name    string        `json:"name"`
	ClassID uint          `json:"class_id"`
	Class   ClassResponse `json:"class,omitempty"` // Tampilkan detail class
}

// CreateGroupRequest adalah DTO untuk membuat group baru
type CreateGroupRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=100"`
	ClassID uint   `json:"class_id" validate:"required,number"`
}

// UpdateGroupRequest adalah DTO untuk memperbarui group
type UpdateGroupRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=100"`
	ClassID uint   `json:"class_id" validate:"required,number"`
}
