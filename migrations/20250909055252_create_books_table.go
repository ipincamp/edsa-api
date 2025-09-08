package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func CreateBooksTable() *gormigrate.Migration {
	type Book struct {
		ID            string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
		Title         string `gorm:"type:varchar(100);not null;uniqueIndex"`
		Description   string `gorm:"type:text"`
		CoverImageURL string `gorm:"type:varchar(255)"`
		Level         int    `gorm:"type:smallint;not null"`
		// Menggunakan pointer string (*string) karena kolom ini bisa NULL (untuk buku pertama).
		UnlockDependencyID *string `gorm:"type:uuid"`
		// Menambahkan constraint foreign key secara eksplisit
		UnlockDependency *Book `gorm:"foreignKey:UnlockDependencyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		CreatedAt        time.Time
		UpdatedAt        time.Time
	}

	return &gormigrate.Migration{
		ID: "20250909055252",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&Book{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("books")
		},
	}
}
