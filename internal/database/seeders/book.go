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

// Helper function to seed a single book cover
func seedBookCover(db *gorm.DB, bookID uint, sourceFilename string) (string, error) {
	cfg := config.AppConfig.Storage
	sourceAssetDir := filepath.Join("internal", "assets", "books")
	// Buat subdirektori "covers" di dalam direktori upload
	coversUploadDir := filepath.Join(cfg.StoragePath, cfg.StorageUploadDir, "covers")
	uploadDirRelative := filepath.Join(cfg.StorageUploadDir, "covers") // "cdn/covers"

	// Pastikan direktori tujuan ada
	if err := os.MkdirAll(coversUploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create covers upload directory '%s': %w", coversUploadDir, err)
	}

	// --- File Handling & Path Generation ---
	fileUUID := uuid.New()
	fileExt := filepath.Ext(sourceFilename)
	newFilenameUUID := fileUUID.String() + fileExt // Nama file fisik (UUID)

	sourcePath := filepath.Join(sourceAssetDir, sourceFilename)     // Path file sumber asli
	destPathUUID := filepath.Join(coversUploadDir, newFilenameUUID) // Path file tujuan fisik (UUID)

	// Path DB relatif terhadap STORAGE_PATH (e.g., "cdn/covers/uuid.png")
	dbFilePathUUID := path.Join(uploadDirRelative, newFilenameUUID)
	// URL Publik (e.g., "http://localhost:8000/cdn/covers/uuid.png")
	publicURLUUID := cfg.StoragePublicBaseURL + path.Join(cfg.StoragePublicURL, dbFilePathUUID)

	// --- Cek & Salin File Fisik ---
	if _, err := os.Stat(destPathUUID); os.IsNotExist(err) {
		log.Printf("   Copying book cover '%s' to '%s'...", sourcePath, destPathUUID)
		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			log.Printf("   WARNING: Could not open source cover '%s': %v. Skipping.", sourcePath, err)
			return "", err // Kembalikan error jika file sumber tidak ada
		}
		defer sourceFile.Close()

		destFile, err := os.Create(destPathUUID)
		if err != nil {
			log.Printf("   WARNING: Could not create destination cover '%s': %v. Skipping.", destPathUUID, err)
			return "", err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, sourceFile)
		if err != nil {
			log.Printf("   WARNING: Could not copy cover from '%s' to '%s': %v. Skipping.", sourcePath, destPathUUID, err)
			os.Remove(destPathUUID) // Coba hapus file yang gagal disalin
			return "", err
		}
		log.Printf("   Successfully copied book cover.")
	} else if err == nil {
		log.Printf("   Physical book cover file '%s' already exists, skipping copy.", newFilenameUUID)
	} else {
		log.Printf("   WARNING: Error checking destination cover file '%s': %v. Skipping.", destPathUUID, err)
		return "", err
	}

	// --- Baca File Size & MIME Type dari file SUMBER ---
	var fileSize int64
	var mimeType string
	fileInfo, err := os.Stat(sourcePath)
	if err != nil {
		log.Printf("   WARNING: Could not stat source cover file '%s' to get size: %v. Using size 0.", sourcePath, err)
		fileSize = 0
	} else {
		fileSize = fileInfo.Size()
	}
	mime, err := mimetype.DetectFile(sourcePath)
	if err != nil {
		log.Printf("   WARNING: Could not detect MIME type for source cover file '%s': %v. Using default 'image/png'.", sourcePath, err)
		mimeType = "image/png" // Fallback
	} else {
		mimeType = mime.String()
	}

	// --- Database Seeding ---
	assetID := uuid.New()
	asset := repo.MediaAssetGORM{
		ID:               assetID,
		FileName:         sourceFilename, // Nama file asli (e.g., cover1.png)
		FilePath:         dbFilePathUUID, // Path relatif disimpan (e.g., cdn/covers/uuid.png)
		PublicURL:        publicURLUUID,  // URL publik lengkap
		MimeType:         mimeType,
		FileSize:         fileSize,
		OwnerType:        domain.OwnerTypeBookCover,              // Tipe pemilik
		OwnerID:          strconv.FormatUint(uint64(bookID), 10), // ID buku sebagai string
		UploadedByUserID: nil,                                    // Di-seed oleh sistem
	}

	// Gunakan FirstOrCreate berdasarkan OwnerType dan OwnerID untuk buku ini
	// Ini mencegah duplikasi aset cover untuk buku yang sama jika seeder dijalankan lagi
	result := db.Where(repo.MediaAssetGORM{
		OwnerType: domain.OwnerTypeBookCover,
		OwnerID:   strconv.FormatUint(uint64(bookID), 10),
	}).Attrs(asset).FirstOrCreate(&asset) // Attrs mengisi nilai jika record baru dibuat

	if result.Error != nil {
		return "", fmt.Errorf("failed to seed book cover asset DB record for book ID %d: %w", bookID, result.Error)
	}

	if result.RowsAffected > 0 {
		log.Printf("   Seeded cover asset DB record: OriginalName='%s', StoredAs='%s' (AssetID: %s)",
			asset.FileName, newFilenameUUID, asset.ID)
	} else {
		log.Printf("   Cover asset DB record for book ID %d already exists (AssetID: %s). Using existing URL.", bookID, asset.ID)
	}

	// Kembalikan URL publik dari asset yang ada atau baru dibuat
	return asset.PublicURL, nil
}

func BookSeeder(db *gorm.DB) error {
	log.Println("Seeding books and their covers...")

	// Map dari tema ke nama file cover
	bookCoverMap := map[string]string{
		"Alphabet":   "cover1.png",
		"Numbers":    "cover2.png",
		"Body Parts": "cover3.png",
		"Family":     "cover4.png",
		"Colors":     "cover1.png", // Contoh: gunakan cover1 lagi karena belum ada
		"Animals":    "cover2.png", // Contoh: gunakan cover2 lagi
		"Time":       "cover3.png", // Contoh: gunakan cover3 lagi
		"Verbs":      "cover4.png", // Contoh: gunakan cover4 lagi
	}

	// Definisikan data buku (tanpa CoverImageURL awal)
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
		book := bookData // Salin data untuk iterasi ini
		// Buat atau cari buku berdasarkan Judul
		result := db.FirstOrCreate(&book, repo.BookGORM{Title: book.Title})
		if result.Error != nil {
			return fmt.Errorf("failed to seed book '%s': %w", book.Title, result.Error)
		}

		created := result.RowsAffected > 0
		if created {
			log.Printf("Seeded book record: %s (ID: %d)", book.Title, book.ID)
		} else {
			log.Printf("Book record '%s' (ID: %d) already exists.", book.Title, book.ID)
		}

		// Ambil nama file cover dari map berdasarkan tema
		coverFilename, ok := bookCoverMap[book.Theme]
		if !ok {
			log.Printf("   WARNING: No cover image defined for theme '%s'. Skipping cover seeding for book '%s'.", book.Theme, book.Title)
			continue // Lanjut ke buku berikutnya jika tidak ada definisi cover
		}

		// Panggil helper untuk seed cover DAN dapatkan URL publiknya
		publicURL, err := seedBookCover(db, book.ID, coverFilename)
		if err != nil {
			// Jika seeding cover gagal, log error tapi lanjutkan (buku tetap ada tanpa cover)
			log.Printf("   ERROR seeding cover for book '%s': %v", book.Title, err)
			continue
		}

		// Perbarui CoverImageURL buku HANYA jika URL baru berbeda ATAU jika buku baru dibuat
		if created || book.CoverImageURL != publicURL {
			book.CoverImageURL = publicURL
			if err := db.Save(&book).Error; err != nil {
				// Jika gagal menyimpan URL cover, log error tapi jangan gagalkan seeder
				log.Printf("   ERROR updating CoverImageURL for book '%s': %v", book.Title, err)
			} else {
				log.Printf("   Updated CoverImageURL for book '%s' to: %s", book.Title, publicURL)
			}
		} else {
			log.Printf("   CoverImageURL for book '%s' is already up-to-date.", book.Title)
		}
	}
	log.Println("Finished seeding books and covers.")
	return nil
}
