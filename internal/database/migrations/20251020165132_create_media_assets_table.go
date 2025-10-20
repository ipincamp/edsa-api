package migrations

import (
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateMediaAssetsTable() *gormigrate.Migration {
	type MediaAsset struct {
		ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primarykey"`
		FileName  string    `gorm:"type:varchar(255);not null"` // Nama file asli
		FilePath  string    `gorm:"type:varchar(255);not null"` // Path relatif di disk
		PublicURL string    `gorm:"type:varchar(255);not null"` // URL lengkap
		MimeType  string    `gorm:"type:varchar(100)"`
		FileSize  int64     `gorm:"not null"`

		// Relasi Polimorfik
		// OwnerID bisa UUID (user) atau UINT (book), jadi kita gunakan string
		OwnerID string `gorm:"type:varchar(255);index"` // Cth: "123e4567-e89b-12d3-a456-426614174000" atau "42"
		// OwnerType untuk menentukan tipe pemilik
		OwnerType string `gorm:"type:varchar(100);index"` // Cth: "book_cover", "user_avatar"

		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt gorm.DeletedAt `gorm:"index"`
	}

	return &gormigrate.Migration{
		ID: "20251020165132",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&MediaAsset{})
		},
		Rollback: func(tx *gorm.DB) error {
			return tx.Migrator().DropTable("media_assets")
		},
	}
}
