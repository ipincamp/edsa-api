package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type localStorageService struct {
	cfg *config.Config
}

func NewLocalStorageService(cfg *config.Config) usecase.FileStorageService {
	return &localStorageService{cfg: cfg}
}

// Upload menyimpan file ke disk lokal
func (s *localStorageService) Upload(file *multipart.FileHeader, fileID uuid.UUID, subDirectory string) (string, error) {
	// 1. Buat nama file unik
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fileID.String() + ext

	// 2. Tentukan path tujuan
	// Cth base upload dir: ./public/cdn
	uploadBaseDir := filepath.Join(s.cfg.Storage.StoragePath, s.cfg.Storage.StorageUploadDir)
	// Cth upload dir spesifik: ./public/cdn/avatars
	uploadDir := filepath.Join(uploadBaseDir, subDirectory)
	// Cth dest path: ./public/cdn/avatars/xxxxxxxx-xxxx.png
	destPath := filepath.Join(uploadDir, uniqueFilename)

	// 3. Buat direktori jika belum ada (termasuk subdirektori)
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory %s: %w", uploadDir, err)
	}

	// 4. Buka file sumber
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// 5. Buat file tujuan
	dst, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file %s: %w", destPath, err)
	}
	defer dst.Close()

	// 6. Salin file
	if _, err = io.Copy(dst, src); err != nil {
		// Jika copy gagal, coba hapus file tujuan yang mungkin sebagian terbuat
		os.Remove(destPath)
		return "", fmt.Errorf("failed to copy file to destination %s: %w", destPath, err)
	}

	// 7. Tentukan path relatif untuk DB (HARUS menyertakan subdirektori)
	// Cth: cdn/avatars/xxxxxxxx-xxxx.png
	relativePath := filepath.ToSlash(filepath.Join(s.cfg.Storage.StorageUploadDir, subDirectory, uniqueFilename))

	return relativePath, nil
}

func (s *localStorageService) Delete(filePath string) error {
	// filePath adalah path relatif, cth: "cdn/uuid.png"

	// 1. Tentukan path absolut
	// Cth: ./public + cdn/uuid.png
	fullPath := filepath.Join(s.cfg.Storage.StoragePath, filePath)

	// 2. Hapus file
	if err := os.Remove(fullPath); err != nil {
		// Jika file tidak ada, itu bukan error kritis, anggap sukses
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
