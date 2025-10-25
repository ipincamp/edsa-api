package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func UpdatePagesTableForL10nAndAudio() *gormigrate.Migration {
	type Page struct {
		ID uint `gorm:"primarykey"`

		// Kolom baru untuk L10n
		Narration_ID string `gorm:"type:text"`
		Narration_EN string `gorm:"type:text"`

		// Kolom baru untuk Audio
		AudioNarrationURL string `gorm:"type:varchar(255)"` // Nullable by default

		// Kolom baru untuk Post Activity
		IsPostActivity bool `gorm:"default:false;not null;index"`
	}

	// Definisikan struct untuk rollback (menambahkan kembali kolom lama)
	type OldPage struct {
		ID            uint   `gorm:"primarykey"`
		NarrativeText string `gorm:"type:text"`
	}

	return &gormigrate.Migration{
		ID: "20251025193204",
		Migrate: func(tx *gorm.DB) error {
			// 1. Tambahkan kolom-kolom baru
			if err := tx.AutoMigrate(&Page{}); err != nil {
				return err
			}

			// 2. Hapus kolom 'narrative_text' yang lama
			return tx.Migrator().DropColumn(&Page{}, "narrative_text")
		},
		Rollback: func(tx *gorm.DB) error {
			// 1. Tambahkan kembali kolom 'narrative_text'
			if err := tx.Migrator().AddColumn(&OldPage{}, "NarrativeText"); err != nil {
				return err
			}

			// 2. Hapus kolom-kolom baru
			if err := tx.Migrator().DropColumn(&Page{}, "narration_id"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&Page{}, "narration_en"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&Page{}, "audio_narration_url"); err != nil {
				return err
			}
			if err := tx.Migrator().DropColumn(&Page{}, "is_post_activity"); err != nil {
				return err
			}
			return nil
		},
	}
}
