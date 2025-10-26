package media

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/pkg/applogger"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type mediaService struct {
	storageSvc usecase.FileStorageService
	mediaRepo  usecase.MediaAssetRepository
	cfg        *config.Config
}

func NewMediaService(
	storageSvc usecase.FileStorageService,
	mediaRepo usecase.MediaAssetRepository,
	cfg *config.Config,
) usecase.MediaService {
	return &mediaService{
		storageSvc: storageSvc,
		mediaRepo:  mediaRepo,
		cfg:        cfg,
	}
}

// --- Mapper DTO ---
// (Helper untuk konversi domain ke response DTO)

func toMediaAssetResponse(a *domain.MediaAsset) *domain.MediaAssetResponse {
	return &domain.MediaAssetResponse{
		ID:        a.ID,
		FileName:  a.FileName,
		PublicURL: a.PublicURL,
		MimeType:  a.MimeType,
		FileSize:  a.FileSize,
		OwnerType: a.OwnerType,
	}
}

// --- Helper Slugify ---

var (
	// Regex untuk menemukan karakter APAPUN yang BUKAN alphanumeric (a-z, A-Z, 0-9)
	slugInvalidChars = regexp.MustCompile(`[^a-zA-Z0-9]+`)

	// Whitelist tipe MIME yang diizinkan untuk di-upload
	allowedMimeTypes = map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/gif":       true,
		"image/svg+xml":   true, // Untuk favicon/logo
		"audio/mpeg":      true, // Untuk MP3
		"audio/wav":       true,
		"video/mp4":       true,
		"application/pdf": true, // Untuk dokumen PDF
	}
)

// slugifyFilename membersihkan nama file dari karakter spesial
// Cth: "bg=$asdfklasdf asdf asdf zoom_11.jpg" -> "bg-asdfklasdf-asdf-asdf-zoom-11.jpg"
func slugifyFilename(originalFilename string) string {
	// 1. Pisahkan nama file dan ekstensi
	ext := filepath.Ext(originalFilename)             // Cth: ".jpg"
	name := strings.TrimSuffix(originalFilename, ext) // Cth: "bg=$asdfklasdf asdf asdf zoom_11"

	// 2. Ganti semua karakter non-alphanumeric (termasuk spasi, _, $, =, dll) dengan 1 strip
	sluggedName := slugInvalidChars.ReplaceAllString(name, "-") // Cth: "bg-asdfklasdf-asdf-asdf-zoom-11"

	// 3. Hapus strip di awal atau akhir (jika ada)
	sluggedName = strings.Trim(sluggedName, "-")

	// 4. Tangani kasus jika nama file menjadi kosong (cth: "---.jpg")
	if sluggedName == "" {
		sluggedName = "file"
	}

	// 5. Gabungkan kembali dengan ekstensi asli
	return sluggedName + ext
}

// --- Methods ---

func (s *mediaService) UploadFile(ctx context.Context, file *multipart.FileHeader, ownerID, ownerType string, uploaderID *uuid.UUID) (*domain.MediaAssetResponse, error) {
	// Validasi tipe MIME file yang di-upload
	src, err := file.Open()
	if err != nil {
		return nil, errors.New("failed to open uploaded file")
	}
	defer src.Close()

	// Deteksi tipe MIME dari konten file, BUKAN dari header
	mime, err := mimetype.DetectReader(src)
	if err != nil {
		return nil, errors.New("failed to detect file type")
	}

	// Cek apakah tipe MIME ada di daftar putih
	if !allowedMimeTypes[mime.String()] {
		applogger.ErrorLogger.Printf("UploadFile: BLOCKED! User tried to upload disallowed file type: %s", mime.String())
		return nil, errors.New("file type not allowed: " + mime.String())
	}

	// Reset file reader kembali ke awal agar bisa dibaca lagi oleh storageSvc
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, errors.New("failed to reset file reader")
	}

	// 1. Buat UUID di sini
	assetID := uuid.New()

	// Tentukan SubDirectory
	var subDirectory string
	switch ownerType {
	case domain.OwnerTypeUserAvatar:
		subDirectory = "avatars"
	case domain.OwnerTypeBookCover:
		subDirectory = filepath.Join("books", "covers")
	// Tambahkan case lain jika perlu
	// case "interaction_audio":
	// 	subDirectory = filepath.Join("interactions", "audio")
	default:
		// Fallback ke direktori 'others' atau root cdn jika tipe tidak dikenal
		subDirectory = "others"
		applogger.ErrorLogger.Printf("UploadFile: Unknown ownerType '%s', saving to '%s' directory", ownerType, subDirectory)
	}

	// 2. Simpan file fisik, dapatkan path relatif (cth: "cdn/avatars/uuid.png")
	// Teruskan subDirectory ke storageSvc.Upload
	filePath, err := s.storageSvc.Upload(file, assetID, subDirectory) // <-- Ubah di sini
	if err != nil {
		applogger.ErrorLogger.Printf("UploadFile: Failed to upload to storage (subdir: %s): %v", subDirectory, err)
		return nil, err
	}

	// 3. Buat URL statis lengkap (Logika ini tetap sama karena filePath sudah benar)
	baseURL := s.cfg.Storage.StoragePublicBaseURL
	publicPath := path.Join(s.cfg.Storage.StoragePublicURL, filePath) // filePath sudah termasuk subdir, misal "/cdn/avatars/uuid.png"
	publicURL := baseURL + publicPath                                 // misal "http://localhost:8000/cdn/avatars/uuid.png"

	// 4. Slugify nama file asli untuk disimpan di DB
	cleanFileName := slugifyFilename(file.Filename)

	// 5. Buat entitas domain
	asset := &domain.MediaAsset{
		ID:               assetID,
		FileName:         cleanFileName,
		FilePath:         filePath,  // "cdn/avatars/uuid.png" atau "cdn/books/covers/uuid.png"
		PublicURL:        publicURL, // "http://localhost:8000/cdn/avatars/uuid.png"
		MimeType:         mime.String(),
		FileSize:         file.Size,
		OwnerID:          ownerID,
		OwnerType:        ownerType,
		UploadedByUserID: uploaderID,
	}

	// 6. Simpan metadata ke database (Logika ini tetap sama)
	if err := s.mediaRepo.Create(ctx, asset); err != nil {
		applogger.ErrorLogger.Printf("UploadFile: Failed to create media asset in DB: %v", err)
		// Jika simpan DB gagal, hapus file fisik
		if delErr := s.storageSvc.Delete(filePath); delErr != nil {
			applogger.ErrorLogger.Printf("UploadFile: CRITICAL! DB insert failed AND physical file delete failed for %s: %v", filePath, delErr)
		}
		return nil, err
	}

	// 7. Kembalikan respons DTO
	return toMediaAssetResponse(asset), nil
}

func (s *mediaService) DeleteFile(ctx context.Context, assetID uuid.UUID, deleterID *uuid.UUID) error {
	// 1. Ambil data aset untuk mendapatkan FilePath
	asset, err := s.mediaRepo.FindByID(ctx, assetID)
	if err != nil {
		applogger.ErrorLogger.Printf("DeleteFile: Failed to find media asset %s: %v", assetID, err)
		return err
	}
	if asset == nil {
		return errors.New("media asset not found")
	}

	// 2. Hapus file fisik dari storage
	if err := s.storageSvc.Delete(asset.FilePath); err != nil {
		// Log error ini tapi jangan hentikan proses.
		// Kita tetap ingin menandai di DB bahwa file ini "dihapus"
		// meskipun file fisiknya gagal dihapus.
		applogger.ErrorLogger.Printf("DeleteFile: CRITICAL! Failed to delete physical file %s: %v", asset.FilePath, err)
	}

	// 3. Soft delete data di database
	if err := s.mediaRepo.SoftDelete(ctx, assetID, deleterID); err != nil {
		applogger.ErrorLogger.Printf("DeleteFile: Failed to soft delete media asset %s in DB: %v", assetID, err)
		return err
	}

	return nil
}
