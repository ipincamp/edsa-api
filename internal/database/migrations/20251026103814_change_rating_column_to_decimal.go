package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func ChangeRatingColumnToDecimal() *gormigrate.Migration {
	type UserBookProgress struct {
		ID     uint    `gorm:"primarykey"`
		Rating float64 `gorm:"type:decimal(2,1);default:0.0"`
	}

	return &gormigrate.Migration{
		ID: "20251026103814",
		Migrate: func(tx *gorm.DB) error {
			// Mengubah tipe kolom 'rating' dari INT ke DECIMAL(2,1)
			// `USING (rating::decimal(2,1))` mengkonversi data yang ada (cth: 4 -> 4.0)
			return tx.Exec("ALTER TABLE user_book_progresses ALTER COLUMN rating TYPE DECIMAL(2,1) USING (rating::decimal(2,1)), ALTER COLUMN rating SET DEFAULT 0.0").Error
		},
		Rollback: func(tx *gorm.DB) error {
			// Mengembalikan tipe kolom 'rating' ke INT
			// `USING (TRUNC(rating)::integer)` membulatkan 4.5 kembali ke 4
			return tx.Exec("ALTER TABLE user_book_progresses ALTER COLUMN rating TYPE INT USING (TRUNC(rating)::integer), ALTER COLUMN rating SET DEFAULT 0").Error
		},
	}
}
