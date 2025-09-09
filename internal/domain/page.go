package domain

import (
	"encoding/json"
	"time"
)

// Page adalah entitas untuk halaman dalam sebuah buku
type Page struct {
	ID             string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	BookID         string          `gorm:"type:uuid;not null;index"`
	PageNumber     int             `gorm:"not null"`
	ContentType    string          `gorm:"type:varchar(255);not null"`
	ContentData    json.RawMessage `gorm:"type:jsonb"` // Fleksibel untuk menyimpan teks, URL gambar, dll.
	HasInteraction bool            `gorm:"not null;default:false"`

	// Relasi
	Book         Book          `gorm:"foreignKey:BookID"`
	Interactions []Interaction `gorm:"foreignKey:PageID"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
