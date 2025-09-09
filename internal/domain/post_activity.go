package domain

import (
	"encoding/json"
	"time"
)

type PostActivity struct {
	ID            string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	BookID        string          `gorm:"type:uuid;not null;index"`
	Type          string          `gorm:"type:varchar(100);not null"`
	ConfigData    json.RawMessage `gorm:"type:jsonb"`
	PointsAwarded int             `gorm:"not null"`
	Book          Book            `gorm:"foreignKey:BookID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
