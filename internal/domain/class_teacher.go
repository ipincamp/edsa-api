package domain

import "time"

type ClassTeacher struct {
	ID               string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	TeacherID        string `gorm:"type:uuid;not null"`
	CourseGroupID    string `gorm:"type:uuid;not null"`
	TeachingCapacity int
	Teacher          User        `gorm:"foreignKey:TeacherID"`
	CourseGroup      CourseGroup `gorm:"foreignKey:CourseGroupID"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
