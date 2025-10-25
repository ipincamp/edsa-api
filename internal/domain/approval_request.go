package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Konstanta untuk Tipe & Status
const (
	ApprovalRequestTypeDeleteAccount = "DELETE_ACCOUNT"
	ApprovalRequestTypeRegisterUser  = "REGISTER_USER"

	ApprovalStatusPending  = "pending"
	ApprovalStatusApproved = "approved"
	ApprovalStatusRejected = "rejected"
)

// ApprovalRequest adalah entitas domain inti
type ApprovalRequest struct {
	ID              uint
	UserID          uuid.UUID
	RequestType     string
	Status          string
	Reason          string
	ReviewerID      *uuid.UUID // Nullable
	ReviewTimestamp *time.Time // Nullable
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt

	// Relasi
	User     User
	Reviewer User
}

// --- Data Transfer Objects (DTOs) ---
