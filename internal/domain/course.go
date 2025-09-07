package domain

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID           string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name         string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Description  string `gorm:"type:text"`
	CourseGroups []CourseGroup
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
