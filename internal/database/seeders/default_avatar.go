package seeders

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/database/shared"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"gorm.io/gorm"
)

func DefaultAvatarsSeeder(db *gorm.DB) error {
	log.Println("Seeding default avatar media assets...")

	cfg := config.AppConfig.Storage
	avatarFilenames := []string{"avatar1.png", "avatar2.png", "avatar3.png", "avatar4.png"}

	for _, filename := range avatarFilenames {
		assetID := uuid.New() // Generate UUID baru untuk setiap aset
		filePath := filepath.ToSlash(filepath.Join(cfg.StorageUploadDir, filename))
		publicURL := cfg.StoragePublicBaseURL + filepath.ToSlash(filepath.Join(cfg.StoragePublicURL, filePath))

		asset := repo.MediaAssetGORM{
			ID:               assetID,
			FileName:         filename,
			FilePath:         filePath,
			PublicURL:        publicURL,
			MimeType:         "image/png", // Asumsi semua adalah PNG
			FileSize:         0,           // Ukuran tidak terlalu penting untuk seeder
			OwnerType:        domain.OwnerTypeUserAvatar,
			OwnerID:          "default", // Tandai sebagai default
			UploadedByUserID: nil,       // Tidak diupload oleh user spesifik
		}

		// Gunakan FirstOrCreate berdasarkan FilePath agar idempotensi terjaga
		result := db.Where(repo.MediaAssetGORM{FilePath: filePath}).FirstOrCreate(&asset)
		if result.Error != nil {
			return fmt.Errorf("failed to seed default avatar asset '%s': %w", filename, result.Error)
		}

		// Simpan ID yang (mungkin baru) dibuat ke map global di paket shared
		shared.DefaultAvatarAssets[filename] = asset.ID // Gunakan asset.ID karena bisa jadi sudah ada

		if result.RowsAffected > 0 {
			log.Printf("Seeded default avatar asset: %s (ID: %s)", filename, asset.ID)
		} else {
			// Jika sudah ada, log ID yang ada
			log.Printf("Default avatar asset '%s' already exists (ID: %s), skipping creation.", filename, asset.ID)
		}
	}
	log.Println("Finished seeding default avatar assets.")
	return nil
}
