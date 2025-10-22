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
func (s *localStorageService) Upload(file *multipart.FileHeader, fileID uuid.UUID) (string, error) {
	// 1. Buat nama file unik
	ext := filepath.Ext(file.Filename)
	uniqueFilename := fileID.String() + ext

	// 2. Tentukan path tujuan
	// Cth: ./public/uploads
	uploadDir := filepath.Join(s.cfg.Storage.StoragePath, s.cfg.Storage.StorageUploadDir)
	// Cth: ./public/uploads/xxxxxxxx-xxxx.png
	destPath := filepath.Join(uploadDir, uniqueFilename)

	// 3. Buat direktori jika belum ada
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
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
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	// 6. Salin file
	if _, err = io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to copy file to destination: %w", err)
	}

	// 7. Tentukan path relatif untuk DB dan URL publik
	// Cth: uploads/xxxxxxxx-xxxx.png
	relativePath := filepath.ToSlash(filepath.Join(s.cfg.Storage.StorageUploadDir, uniqueFilename))

	return relativePath, nil
}

func (s *localStorageService) Delete(filePath string) error {
	// filePath adalah path relatif, cth: "uploads/uuid.png"

	// 1. Tentukan path absolut
	// Cth: ./public + uploads/uuid.png
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
