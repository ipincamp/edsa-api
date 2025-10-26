package seeders

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"gorm.io/gorm"
)

// Update function signature to accept logger
func DefaultAvatarsSeeder(db *gorm.DB, logger *log.Logger) error {
	logger.Println("Seeding default avatar media assets (Original Filename, UUID Storage)...")

	cfg := config.AppConfig.Storage
	originalAvatarMap := map[string]string{
		"avatar1.png": "avatar1.png",
		"avatar2.png": "avatar2.png",
		"avatar3.png": "avatar3.png",
		"avatar4.png": "avatar4.png",
	}

	sourceAssetDir := filepath.Join("internal", "assets", "avatars")
	// Direktori fisik upload (e.g., ./public/cdn/avatars)
	avatarSubDir := "avatars"
	uploadDir := filepath.Join(cfg.StoragePath, cfg.StorageUploadDir, avatarSubDir)

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create public upload directory '%s': %w", uploadDir, err)
	}

	for logicalName, sourceFilename := range originalAvatarMap {
		logger.Printf("   Processing avatar '%s' (logical name: %s)...", sourceFilename, logicalName)
		sourcePath := filepath.Join(sourceAssetDir, sourceFilename)

		// Cek file sumber dulu
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			logger.Printf("   -> WARNING: Source avatar file '%s' not found. Skipping.", sourcePath)
			continue // Lanjut ke avatar berikutnya jika file sumber tidak ada
		} else if err != nil {
			logger.Printf("   -> ERROR: Could not check source avatar file '%s': %v. Skipping.", sourcePath, err)
			continue // Error lain saat cek file sumber, skip
		}

		// File sumber ada, lanjutkan
		fileUUID := uuid.New()
		fileExt := filepath.Ext(sourceFilename)
		newFilenameUUID := fileUUID.String() + fileExt            // Nama file fisik (UUID)
		destPathUUID := filepath.Join(uploadDir, newFilenameUUID) // Path file tujuan fisik (UUID)

		// Path DB (e.g., "cdn/avatars/uuid.png")
		dbFilePathUUID := path.Join(cfg.StorageUploadDir, avatarSubDir, newFilenameUUID)
		// URL Publik (e.g., "http://localhost:8000/cdn/avatars/uuid.png")
		publicURLUUID := cfg.StoragePublicBaseURL + path.Join(cfg.StoragePublicURL, dbFilePathUUID)

		// --- Cek & Salin File Fisik ---
		if _, err := os.Stat(destPathUUID); os.IsNotExist(err) {
			logger.Printf("   -> Copying '%s' to '%s'...", sourcePath, destPathUUID)
			sourceFile, err := os.Open(sourcePath)
			if err != nil {
				logger.Printf("   -> WARNING: Could not open source avatar '%s' after stat: %v. Skipping.", sourcePath, err)
				continue // Gagal buka sumber
			}
			// Use anonymous function for scoped defer
			func() {
				defer sourceFile.Close()
				destFile, err := os.Create(destPathUUID)
				if err != nil {
					logger.Printf("   -> ERROR: Could not create destination avatar '%s': %v. Skipping.", destPathUUID, err)
					return // Return from anonymous function
				}
				defer destFile.Close()

				_, err = io.Copy(destFile, sourceFile)
				if err != nil {
					logger.Printf("   -> ERROR: Could not copy avatar from '%s' to '%s': %v. Skipping.", sourcePath, destPathUUID, err)
					os.Remove(destPathUUID) // Coba hapus file gagal
					return                  // Return from anonymous function
				}
				logger.Printf("   -> Successfully copied avatar.")
			}() // Call the anonymous function
			// Check if error occurred inside anonymous function (alternative: pass error out)
			// This simplified version relies on logging inside
		} else if err == nil {
			logger.Printf("   -> Physical avatar file '%s' already exists, skipping copy.", newFilenameUUID)
		} else {
			logger.Printf("   -> ERROR: Error checking destination avatar file '%s': %v. Skipping.", destPathUUID, err)
			continue // Gagal cek tujuan
		}

		// --- Baca File Size & MIME Type ---
		var fileSize int64
		var mimeType string
		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			logger.Printf("   -> WARNING: Could not stat source avatar file '%s' to get size: %v. Using size 0.", sourcePath, err)
			fileSize = 0
		} else {
			fileSize = fileInfo.Size()
		}
		mime, err := mimetype.DetectFile(sourcePath)
		if err != nil {
			logger.Printf("   -> WARNING: Could not detect MIME type for source avatar file '%s': %v. Using default 'image/png'.", sourcePath, err)
			mimeType = "image/png"
		} else {
			mimeType = mime.String()
		}

		// --- Database Seeding MediaAsset ---
		assetID := uuid.New()
		asset := repo.MediaAssetGORM{
			ID:               assetID,
			FileName:         sourceFilename,
			FilePath:         dbFilePathUUID,
			PublicURL:        publicURLUUID,
			MimeType:         mimeType,
			FileSize:         fileSize,
			OwnerType:        domain.OwnerTypeUserAvatar,
			OwnerID:          fmt.Sprintf("default-%s", logicalName),
			UploadedByUserID: nil,
		}

		// Gunakan FirstOrCreate berdasarkan OwnerType dan OwnerID unik ini
		result := db.Where(repo.MediaAssetGORM{
			OwnerType: domain.OwnerTypeUserAvatar,
			OwnerID:   fmt.Sprintf("default-%s", logicalName),
		}).Attrs(asset).FirstOrCreate(&asset) // Attrs mengisi nilai jika record baru dibuat

		if result.Error != nil {
			// Gagal seed asset -> Log error tapi lanjutkan ke avatar berikutnya
			logger.Printf("   -> ERROR seeding default avatar asset DB record for '%s': %v. Skipping.", sourceFilename, result.Error)
			continue
		}

		// Log detail asset yang di-seed atau ditemukan
		logMsgFormat := "   -> Seeded asset DB record: OriginalName='%s', StoredAs='%s' (AssetID: %s, Size: %d, Type: %s)"
		if result.RowsAffected == 0 { // Jika record sudah ada
			logMsgFormat = "   -> Found existing asset DB record: OriginalName='%s', StoredAs='%s' (AssetID: %s, Size: %d, Type: %s)"
		}
		logger.Printf(logMsgFormat, asset.FileName, newFilenameUUID, asset.ID, asset.FileSize, asset.MimeType)

	} // End loop for originalAvatarMap

	logger.Println("Finished seeding default avatar assets.")
	return nil
}
