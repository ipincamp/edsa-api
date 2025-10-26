package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUsersTable() *gormigrate.Migration {
	type Role struct {
		ID uint `gorm:"primarykey"`
	}

	type User struct {
		ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primarykey"`
		Name      string    `gorm:"type:varchar(255)"`
		Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
		Password  string    `gorm:"type:varchar(255);not null"`
		RoleID    uint      `gorm:"not null"`
		Role      Role      `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251019121635",
		Migrate: func(tx *gorm.DB) error {
			if err := tx.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\"").Error; err != nil {
				return err
			}
			return tx.AutoMigrate(&User{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable(&User{})
		},
	}
}
