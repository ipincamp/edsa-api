package media

import (
	"context"
	"mime/multipart"
	"path"

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

	// 4. Buat entitas domain
	asset := &domain.MediaAsset{
		ID:        assetID,
		FileName:  file.Filename,
		FilePath:  filePath,  // "uploads/uuid.png"
		PublicURL: publicURL, // "http://localhost:8080/public/uploads/uuid.png"
		MimeType:  file.Header.Get("Content-Type"),
		FileSize:  file.Size,
		OwnerID:   ownerID,
		OwnerType: ownerType,
	}

	// 5. Simpan metadata ke database
	if err := s.mediaRepo.Create(ctx, asset); err != nil {
		// TODO: Implementasikan rollback (hapus file fisik jika gagal simpan DB)
		return nil, err
	}

	// 6. Kembalikan respons DTO
	return toMediaAssetResponse(asset), nil
}
