package gorm

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RoleGORM adalah representasi tabel 'roles' di database
type RoleGORM struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"type:varchar(50);uniqueIndex;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (RoleGORM) TableName() string {
	return "roles"
}

// UserGORM adalah representasi tabel 'users' di database
type UserGORM struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primarykey"`
	Name      string    `gorm:"type:varchar(255)"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string    `gorm:"type:varchar(255);not null"`
	RoleID    uint      `gorm:"not null"`
	Role      RoleGORM  `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (UserGORM) TableName() string {
	return "users"
}

// SubjectGORM adalah representasi tabel 'subjects' di database
type SubjectGORM struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"type:varchar(255);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (SubjectGORM) TableName() string {
	return "subjects"
}
