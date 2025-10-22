package logger

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type activityLoggerService struct {
	repo usecase.ActivityLogRepository
}

func NewActivityLoggerService(repo usecase.ActivityLogRepository) usecase.ActivityLoggerService {
	return &activityLoggerService{repo: repo}
}

// --- Methods ---

// Log menjalankan logging di goroutine baru agar tidak memblokir request utama
func (s *activityLoggerService) Log(ctx context.Context, logData domain.ActivityLog) {
	go func() {
		// Buat context baru untuk goroutine
		bgCtx := context.Background()

		// Isi data yang hilang
		if logData.TimestampStart.IsZero() {
			logData.TimestampStart = time.Now()
		}
		if logData.SessionID == uuid.Nil {
			logData.SessionID = uuid.New() // Buat ID unik untuk setiap log
		}

		if err := s.repo.Create(bgCtx, &logData); err != nil {
			// Jangan panic, cukup log error-nya
			log.Printf("ERROR: Failed to create activity log: %v\n", err)
		}
	}()
}
