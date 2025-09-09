package domain

import (
	"encoding/json"
	"time"
)

type Interaction struct {
	ID            string          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	PageID        string          `gorm:"type:uuid;not null;index"`
	Type          string          `gorm:"type:varchar(100);not null"`
	ConfigData    json.RawMessage `gorm:"type:jsonb"`
	PointsAwarded int             `gorm:"not null"`
	Page          Page            `gorm:"foreignKey:PageID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
