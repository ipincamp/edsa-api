package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Konstanta untuk OwnerType
const (
	OwnerTypeUserAvatar = "user_avatar"
	OwnerTypeBookCover  = "book_cover"
	// Tambahkan tipe lain di sini jika perlu
)

// MediaAsset adalah entitas domain inti untuk aset media
type MediaAsset struct {
	ID               uuid.UUID
	FileName         string
	FilePath         string
	PublicURL        string
	MimeType         string
	FileSize         int64
	OwnerID          string
	OwnerType        string
	UploadedByUserID *uuid.UUID
	DeletedByUserID  *uuid.UUID
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt
}

// --- Data Transfer Objects (DTOs) ---

// MediaAssetResponse adalah DTO untuk respons
type MediaAssetResponse struct {
	ID        uuid.UUID `json:"id"`
	FileName  string    `json:"file_name"`
	PublicURL string    `json:"url"`
	MimeType  string    `json:"mime_type"`
	FileSize  int64     `json:"file_size"`
	OwnerType string    `json:"owner_type"`
}
