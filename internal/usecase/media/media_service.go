package media

import (
	"context"
	"mime/multipart"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type mediaService struct {
	storageSvc usecase.FileStorageService
	mediaRepo  usecase.MediaAssetRepository
}

func NewMediaService(
	storageSvc usecase.FileStorageService,
	mediaRepo usecase.MediaAssetRepository,
) usecase.MediaService {
	return &mediaService{
		storageSvc: storageSvc,
		mediaRepo:  mediaRepo,
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
	// 1. Simpan file fisik menggunakan storage service
	publicURL, filePath, err := s.storageSvc.Upload(file)
	if err != nil {
		return nil, err
	}

	// 2. Buat entitas domain
	asset := &domain.MediaAsset{
		FileName:  file.Filename,
		FilePath:  filePath,
		PublicURL: publicURL,
		MimeType:  file.Header.Get("Content-Type"),
		FileSize:  file.Size,
		OwnerID:   ownerID,
		OwnerType: ownerType,
	}

	// 3. Simpan metadata ke database
	if err := s.mediaRepo.Create(ctx, asset); err != nil {
		// TODO: Implementasikan rollback (hapus file fisik jika gagal simpan DB)
		return nil, err
	}

	// 4. Kembalikan respons DTO
	return toMediaAssetResponse(asset), nil
}
