package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func UpdateUserBookProgressTable() *gormigrate.Migration {
	type UserBookProgress struct {
		ID uint `gorm:"primarykey"`

		// 1. Ubah tipe data kolom ini
		HighestScore float64 `gorm:"type:decimal(5,2);default:0"`

		// 2. Tambahkan kolom baru
		Rating int `gorm:"default:0;index"`
	}

	return &gormigrate.Migration{
		ID: "20251025193919",
		Migrate: func(tx *gorm.DB) error {
			return tx.AutoMigrate(&UserBookProgress{})
		},
		Rollback: func(tx *gorm.DB) error {
			// 1. Hapus kolom 'rating'
			if err := tx.Migrator().DropColumn(&UserBookProgress{}, "rating"); err != nil {
				return err
			}

			// 2. Ubah tipe 'highest_score' kembali ke INT
			// Kita gunakan TRUNC() untuk mengonversi decimal (cth: 85.5) kembali ke integer (85)
			return tx.Exec("ALTER TABLE user_book_progresses ALTER COLUMN highest_score TYPE INT USING (TRUNC(highest_score)::integer), ALTER COLUMN highest_score SET DEFAULT 0").Error
		},
	}
}
