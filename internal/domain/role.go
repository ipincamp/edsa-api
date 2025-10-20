package domain

import (
	"time"

	"gorm.io/gorm"
)

// Definisi nama-nama role yang tersedia
const (
	RoleNameAdmin   = "admin"
	RoleNameTeacher = "teacher"
	RoleNameStudent = "student"
	RoleNamePublic  = "public"
)

// Role adalah entitas domain untuk role pengguna
type Role struct {
	ID        uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
