package media

import (
	"context"
	"mime/multipart"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/domain"
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

func (s *mediaService) UploadFile(ctx context.Context, file *multipart.FileHeader, ownerID, ownerType string) (*domain.MediaAssetResponse, error) {
	// 1. Buat UUID di sini
	assetID := uuid.New()

	// 2. Simpan file fisik, dapatkan path relatif (cth: "uploads/uuid.png")
	filePath, err := s.storageSvc.Upload(file, assetID)
	if err != nil {
		return nil, err
	}

	// 3. Buat URL statis lengkap
	// Cth: /public + uploads/uuid.png -> /public/uploads/uuid.png
	baseURL := s.cfg.Storage.StoragePublicBaseURL
	publicURL := baseURL + path.Join(s.cfg.Storage.StoragePublicURL, filePath)

	// 4. Slugify nama file asli untuk disimpan di DB
	cleanFileName := slugifyFilename(file.Filename)

	// 5. Buat entitas domain
	asset := &domain.MediaAsset{
		ID:        assetID,
		FileName:  cleanFileName,
		FilePath:  filePath,  // "uploads/uuid.png"
		PublicURL: publicURL, // "http://localhost:8080/public/uploads/uuid.png"
		MimeType:  file.Header.Get("Content-Type"),
		FileSize:  file.Size,
		OwnerID:   ownerID,
		OwnerType: ownerType,
	}

	// 6. Simpan metadata ke database
	if err := s.mediaRepo.Create(ctx, asset); err != nil {
		// TODO: Implementasikan rollback (hapus file fisik jika gagal simpan DB)
		return nil, err
	}

	// 7. Kembalikan respons DTO
	return toMediaAssetResponse(asset), nil
}
