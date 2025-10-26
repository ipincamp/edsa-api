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

func DefaultAvatarsSeeder(db *gorm.DB) error {
	log.Println("Seeding default avatar media assets (Original Filename, UUID Storage)...")

	cfg := config.AppConfig.Storage
	originalAvatarMap := map[string]string{
		"avatar1.png": "avatar1.png",
		"avatar2.png": "avatar2.png",
		"avatar3.png": "avatar3.png",
		"avatar4.png": "avatar4.png",
	}

	sourceAssetDir := filepath.Join("internal", "assets", "avatars")
	uploadDir := filepath.Join(cfg.StoragePath, cfg.StorageUploadDir)

	// Pastikan direktori tujuan ada
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create public upload directory '%s': %w", uploadDir, err)
	}

	for logicalName, sourceFilename := range originalAvatarMap {
		// --- File Handling & Path Generation ---
		fileUUID := uuid.New()
		fileExt := filepath.Ext(sourceFilename)
		newFilenameUUID := fileUUID.String() + fileExt // Nama file fisik (UUID)

		sourcePath := filepath.Join(sourceAssetDir, sourceFilename) // Path file sumber asli
		destPathUUID := filepath.Join(uploadDir, newFilenameUUID)   // Path file tujuan fisik (UUID)

		// UBAH CARA PEMBUATAN PATH DB DAN URL PUBLIK
		// Path DB sekarang hanya nama direktori upload + nama file UUID
		// e.g., "cdn/uuid.png"
		dbFilePathUUID := path.Join(cfg.StorageUploadDir, newFilenameUUID)
		// URL Publik sekarang base URL + StoragePublicURL (/) + path DB
		// e.g., "http://localhost:8000" + "/" + "cdn/uuid.png" -> "http://localhost:8000/cdn/uuid.png"
		publicURLUUID := cfg.StoragePublicBaseURL + path.Join(cfg.StoragePublicURL, dbFilePathUUID)

		// --- Cek & Salin File Fisik ---
		if _, err := os.Stat(destPathUUID); os.IsNotExist(err) {
			log.Printf("   Copying '%s' to '%s'...", sourcePath, destPathUUID)
			sourceFile, err := os.Open(sourcePath)
			if err != nil {
				log.Printf("   WARNING: Could not open source avatar '%s': %v. Skipping.", sourcePath, err)
				continue
			}

			destFile, err := os.Create(destPathUUID)
			if err != nil {
				log.Printf("   WARNING: Could not create destination avatar '%s': %v. Skipping.", destPathUUID, err)
				sourceFile.Close() // Pastikan source ditutup
				continue
			}

			_, err = io.Copy(destFile, sourceFile)
			sourceFile.Close() // Tutup source setelah copy
			destFile.Close()   // Tutup dest setelah copy
			if err != nil {
				log.Printf("   WARNING: Could not copy avatar from '%s' to '%s': %v. Skipping.", sourcePath, destPathUUID, err)
				os.Remove(destPathUUID)
				continue
			}
			log.Printf("   Successfully copied.")
		} else if err == nil {
			log.Printf("   Physical file '%s' already exists, skipping copy.", newFilenameUUID)
		} else {
			// Error lain saat cek file UUID
			log.Printf("   WARNING: Error checking destination file '%s': %v. Skipping.", destPathUUID, err)
			continue
		}

		// --- Baca File Size & MIME Type dari file SUMBER ---
		var fileSize int64
		var mimeType string
		fileInfo, err := os.Stat(sourcePath) // Baca dari path sumber baru
		if err != nil {
			log.Printf("   WARNING: Could not stat source file '%s' to get size: %v. Using size 0.", sourcePath, err)
			fileSize = 0
		} else {
			fileSize = fileInfo.Size()
		}
		mime, err := mimetype.DetectFile(sourcePath) // Deteksi dari path sumber baru
		if err != nil {
			log.Printf("   WARNING: Could not detect MIME type for source file '%s': %v. Using default 'image/png'.", sourcePath, err)
			mimeType = "image/png"
		} else {
			mimeType = mime.String()
		}

		// --- Database Seeding ---
		assetID := uuid.New()
		asset := repo.MediaAssetGORM{
			ID:               assetID,
			FileName:         sourceFilename, // Nama Asli
			FilePath:         dbFilePathUUID, // Path UUID (e.g., "cdn/uuid.png")
			PublicURL:        publicURLUUID,  // URL UUID (e.g., "http://localhost:8000/cdn/uuid.png")
			MimeType:         mimeType,
			FileSize:         fileSize,
			OwnerType:        domain.OwnerTypeUserAvatar,
			OwnerID:          fmt.Sprintf("default-%s", logicalName),
			UploadedByUserID: nil,
		}
		result := db.Where(repo.MediaAssetGORM{FilePath: dbFilePathUUID}).FirstOrCreate(&asset)
		if result.Error != nil {
			// TODO: Jika DB gagal setelah file disalin, perlu rollback manual file
			return fmt.Errorf("failed to seed default avatar asset DB record for '%s': %w", newFilenameUUID, result.Error)
		}

		if result.RowsAffected > 0 {
			log.Printf("Seeded asset DB record: OriginalName='%s', StoredAs='%s' (AssetID: %s, Size: %d, Type: %s)",
				asset.FileName, newFilenameUUID, asset.ID, asset.FileSize, asset.MimeType)
		} else {
			// Jika record DB sudah ada (berdasarkan FilePath UUID), pastikan map shared tetap terisi ID yang benar
			log.Printf("Asset DB record for stored file '%s' already exists (AssetID: %s). Ensuring map uses AssetID.",
				newFilenameUUID, asset.ID)
		}
	}
	log.Println("Finished seeding default avatar assets.")
	return nil
}
