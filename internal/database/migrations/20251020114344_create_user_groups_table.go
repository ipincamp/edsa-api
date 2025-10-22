package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateUserGroupsTable() *gormigrate.Migration {
	type User struct {
		ID uuid.UUID `gorm:"type:uuid;primarykey"`
	}
	type Group struct {
		ID uint `gorm:"primarykey"`
	}

	type UserGroup struct {
		UserID  uuid.UUID `gorm:"primaryKey"`
		User    User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
		GroupID uint      `gorm:"primaryKey"`
		Group   Group     `gorm:"foreignKey:GroupID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	}

	return &gormigrate.Migration{
		ID: "20251020114344",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&UserGroup{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("user_groups")
		},
	}
}
