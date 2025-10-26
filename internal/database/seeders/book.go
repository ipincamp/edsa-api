package seeders

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"gorm.io/gorm"
)

const placeholderCoverURL = "https://placehold.co/800x400/EEE/31343C/png?text=No+Cover+Yet"

// Helper function to seed a single book cover
func seedBookCover(db *gorm.DB, bookID uint, bookTitle, sourceFilename string) (string, error) { // Tambahkan bookTitle untuk logging
	cfg := config.AppConfig.Storage
	sourceAssetDir := filepath.Join("internal", "assets", "books")
	coversUploadDir := filepath.Join(cfg.StoragePath, cfg.StorageUploadDir, "covers")
	uploadDirRelative := filepath.Join(cfg.StorageUploadDir, "covers") // "cdn/covers"
	sourcePath := filepath.Join(sourceAssetDir, sourceFilename)        // Path file sumber asli

	// Log awal untuk cover buku ini
	log.Printf("   Processing cover '%s' for book '%s' (ID: %d)...", sourceFilename, bookTitle, bookID)

	// *** PENGECEKAN KEBERADAAN FILE SUMBER ***
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		log.Printf("   -> WARNING: Source cover file '%s' not found. Using placeholder URL.", sourcePath)
		return placeholderCoverURL, nil // File tidak ada -> gunakan placeholder
	} else if err != nil {
		log.Printf("   -> ERROR: Could not check source cover file '%s': %v. Using placeholder URL.", sourcePath, err)
		return placeholderCoverURL, err // Error lain -> gunakan placeholder, kembalikan error
	}

	// *** FILE SUMBER ADA, LANJUTKAN PROSES ***
	if err := os.MkdirAll(coversUploadDir, os.ModePerm); err != nil {
		// Gagal buat direktori tujuan -> error fatal untuk cover ini
		return "", fmt.Errorf("   -> ERROR: Failed to create covers upload directory '%s': %w", coversUploadDir, err)
	}

	// --- File Handling & Path Generation ---
	fileUUID := uuid.New()
	fileExt := filepath.Ext(sourceFilename)
	newFilenameUUID := fileUUID.String() + fileExt
	destPathUUID := filepath.Join(coversUploadDir, newFilenameUUID)
	dbFilePathUUID := path.Join(uploadDirRelative, newFilenameUUID)
	publicURLUUID := cfg.StoragePublicBaseURL + path.Join(cfg.StoragePublicURL, dbFilePathUUID)

	// --- Cek & Salin File Fisik ---
	if _, err := os.Stat(destPathUUID); os.IsNotExist(err) {
		log.Printf("   -> Copying '%s' to '%s'...", sourcePath, destPathUUID) // Log copy
		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			log.Printf("   -> WARNING: Could not open source cover '%s' after stat: %v. Using placeholder.", sourcePath, err)
			return placeholderCoverURL, nil // Gagal buka sumber -> gunakan placeholder
		}
		defer sourceFile.Close()

		destFile, err := os.Create(destPathUUID)
		if err != nil {
			log.Printf("   -> ERROR: Could not create destination cover '%s': %v. Using placeholder.", destPathUUID, err)
			return placeholderCoverURL, err // Gagal buat tujuan -> gunakan placeholder, kembalikan error
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, sourceFile)
		if err != nil {
			log.Printf("   -> ERROR: Could not copy cover from '%s' to '%s': %v. Using placeholder.", sourcePath, destPathUUID, err)
			os.Remove(destPathUUID)
			return placeholderCoverURL, err // Gagal copy -> gunakan placeholder, kembalikan error
		}
		log.Printf("   -> Successfully copied book cover.") // Log sukses copy
	} else if err == nil {
		log.Printf("   -> Physical book cover file '%s' already exists, skipping copy.", newFilenameUUID)
	} else {
		log.Printf("   -> ERROR: Error checking destination cover file '%s': %v. Using placeholder.", destPathUUID, err)
		return placeholderCoverURL, err // Gagal cek tujuan -> gunakan placeholder, kembalikan error
	}

	// --- Baca File Size & MIME Type ---
	var fileSize int64
	var mimeType string
	fileInfo, err := os.Stat(sourcePath)
	if err != nil {
		log.Printf("   -> WARNING: Could not stat source cover file '%s' to get size: %v. Using size 0.", sourcePath, err)
		fileSize = 0 // Tetap lanjutkan dengan size 0
	} else {
		fileSize = fileInfo.Size()
	}
	mime, err := mimetype.DetectFile(sourcePath)
	if err != nil {
		log.Printf("   -> WARNING: Could not detect MIME type for source cover file '%s': %v. Using default 'image/png'.", sourcePath, err)
		mimeType = "image/png" // Fallback
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
		OwnerType:        domain.OwnerTypeBookCover,
		OwnerID:          strconv.FormatUint(uint64(bookID), 10),
		UploadedByUserID: nil,
	}

	// FirstOrCreate berdasarkan OwnerType dan OwnerID
	result := db.Where(repo.MediaAssetGORM{
		OwnerType: domain.OwnerTypeBookCover,
		OwnerID:   strconv.FormatUint(uint64(bookID), 10),
	}).Attrs(asset).FirstOrCreate(&asset)

	if result.Error != nil {
		// Gagal seed asset -> Log error tapi kembalikan placeholder
		log.Printf("   -> ERROR seeding book cover asset DB record for book ID %d: %v. Using placeholder.", bookID, result.Error)
		return placeholderCoverURL, nil
	}

	// Log detail asset yang di-seed atau ditemukan (mirip AvatarSeeder)
	logMsgFormat := "   -> Seeded asset DB record: OriginalName='%s', StoredAs='%s' (AssetID: %s, Size: %d, Type: %s)"
	if result.RowsAffected == 0 { // Jika record sudah ada
		logMsgFormat = "   -> Found existing asset DB record: OriginalName='%s', StoredAs='%s' (AssetID: %s, Size: %d, Type: %s)"
	}
	log.Printf(logMsgFormat, asset.FileName, newFilenameUUID, asset.ID, asset.FileSize, asset.MimeType)

	return asset.PublicURL, nil // Kembalikan URL asli yang di-seed/ditemukan
}

func BookSeeder(db *gorm.DB, logger *log.Logger) error {
	logger.Println("Seeding books and their covers...")

	bookCoverMap := map[string]string{
		"Alphabet":   "cover1.png",
		"Numbers":    "cover2.png",
		"Body Parts": "cover3.png",
		"Family":     "cover4.png",
		"Colors":     "cover5.png",
		"Animals":    "cover6.png",
		"Time":       "cover7.png",
		"Verbs":      "cover8.png",
	}

	booksData := []repo.BookGORM{
		{Title: "The Alphabet in the Land of Dewi Sri", Description: "An exciting journey through the ABCs.", Theme: "Alphabet", BookOrder: 1},
		{Title: "The Numbers in the Village of Ten Hills", Description: "Learn numbers with friendly creatures.", Theme: "Numbers", BookOrder: 2},
		{Title: "Bawang Putih and the Kind Body Parts", Description: "Discover all the parts of your body.", Theme: "Body Parts", BookOrder: 3},
		{Title: "The Big Mango Tree and the Family of Five", Description: "A wonderful day out with the whole family.", Theme: "Family", BookOrder: 4},
		{Title: "The Magic Paintbrush", Description: "Mixing colors to create new ones.", Theme: "Colors", BookOrder: 5},
		{Title: "Sounds of the Wild", Description: "Listen to the sounds animals make.", Theme: "Animals", BookOrder: 6},
		{Title: "What Time Is It?", Description: "Learn to tell time with Timmy the clock.", Theme: "Time", BookOrder: 7},
		{Title: "Action Day!", Description: "Jumping, running, and playing all day.", Theme: "Verbs", BookOrder: 8},
	}

	for _, bookData := range booksData {
		book := bookData
		result := db.FirstOrCreate(&book, repo.BookGORM{Title: book.Title})
		if result.Error != nil {
			return fmt.Errorf("failed to seed book '%s': %w", book.Title, result.Error)
		}

		created := result.RowsAffected > 0
		if created {
			logger.Printf("Seeded book record: '%s' (ID: %d)", book.Title, book.ID)
		} else {
			logger.Printf("Book record '%s' (ID: %d) already exists.", book.Title, book.ID)
		}

		coverFilename, ok := bookCoverMap[book.Theme]
		var targetURL string
		var coverErr error

		if !ok {
			logger.Printf("   -> WARNING: No cover image defined for theme '%s'. Setting placeholder URL for book '%s'.", book.Theme, book.Title)
			targetURL = placeholderCoverURL
			coverErr = nil // Tidak ada error jika hanya tidak ada definisi
		} else {
			// Panggil helper seedBookCover dengan book.Title untuk logging
			targetURL, coverErr = seedBookCover(db, book.ID, book.Title, coverFilename)
			if coverErr != nil {
				// Log error dari helper jika ada (helper sudah log detailnya)
				logger.Printf("   -> ERROR processing cover for book '%s'. URL set to placeholder.", book.Title)
				targetURL = placeholderCoverURL // Pastikan targetURL adalah placeholder jika ada error
			}
		}

		// Perbarui CoverImageURL buku jika perlu (jika baru dibuat atau URL berbeda)
		if created || book.CoverImageURL != targetURL {
			oldURL := book.CoverImageURL // Simpan URL lama untuk logging
			book.CoverImageURL = targetURL
			if err := db.Save(&book).Error; err != nil {
				logger.Printf("   -> ERROR updating CoverImageURL for book '%s' (ID: %d): %v", book.Title, book.ID, err)
			} else {
				if oldURL == "" && created { // Kasus buku baru
					logger.Printf("   -> Set CoverImageURL for new book '%s' to: %s", book.Title, targetURL)
				} else { // Kasus update URL
					logger.Printf("   -> Updated CoverImageURL for book '%s' from '%s' to: %s", book.Title, oldURL, targetURL)
				}
			}
		} else {
			// Log jika URL sudah sesuai, baik itu placeholder maupun URL asli
			logger.Printf("   -> CoverImageURL for book '%s' is already set to: %s", book.Title, book.CoverImageURL)
		}
	}
	logger.Println("Finished seeding books and covers.")
	return nil
}
