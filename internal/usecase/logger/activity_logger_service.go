package logger

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/pkg/applogger"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

const (
	logChannelBufferSize = 100             // Jumlah log yang bisa ditampung di channel sebelum blocking
	maxBatchSize         = 50              // Jumlah log per batch insert ke DB
	batchTimeout         = 3 * time.Second // Waktu tunggu maks sebelum batch dikirim
	dbInsertTimeout      = 5 * time.Second // Waktu tunggu maks untuk query DB
)

type activityLoggerService struct {
	repo       usecase.ActivityLogRepository
	logChannel chan domain.ActivityLog
	wg         sync.WaitGroup
	done       chan struct{}
}

func NewActivityLoggerService(repo usecase.ActivityLogRepository) usecase.ActivityLoggerService {
	s := &activityLoggerService{
		repo:       repo,
		logChannel: make(chan domain.ActivityLog, logChannelBufferSize),
		done:       make(chan struct{}),
	}

	// Mulai worker di background
	s.wg.Add(1)
	go s.worker()

	return s
}

func (s *activityLoggerService) Log(ctx context.Context, logData domain.ActivityLog) {
	// Isi data yang hilang
	if logData.TimestampStart.IsZero() {
		logData.TimestampStart = time.Now()
	}
	if logData.SessionID == uuid.Nil {
		logData.SessionID = uuid.New() // Buat ID unik jika tidak ada
	}

	// Kirim ke channel (non-blocking)
	select {
	case s.logChannel <- logData:
		// Log berhasil masuk antrian
	default:
		// Channel penuh, log dijatuhkan (dropped)
		applogger.ErrorLogger.Printf("WARNING: Activity log channel is full. Dropping log for action: %s", logData.Action)
	}
}

func (s *activityLoggerService) Shutdown(ctx context.Context) error {
	log.Println("Shutting down activity logger service...")
	// Kirim sinyal berhenti ke worker
	close(s.done)
	// Tunggu worker selesai memproses sisa batch
	s.wg.Wait()
	log.Println("Activity logger service shut down gracefully.")
	return nil
}

func (s *activityLoggerService) worker() {
	defer s.wg.Done()

	batch := make([]domain.ActivityLog, 0, maxBatchSize)
	ticker := time.NewTicker(batchTimeout)

	for {
		select {
		case logEntry := <-s.logChannel:
			batch = append(batch, logEntry)
			if len(batch) >= maxBatchSize {
				// Batch penuh, kirim
				s.commitBatch(batch)
				batch = make([]domain.ActivityLog, 0, maxBatchSize)
				ticker.Reset(batchTimeout) // Reset timer
			}

		case <-ticker.C:
			// Waktu habis, kirim batch yang ada (jika tidak kosong)
			if len(batch) > 0 {
				s.commitBatch(batch)
				batch = make([]domain.ActivityLog, 0, maxBatchSize)
			}

		case <-s.done:
			// Menerima sinyal shutdown
			ticker.Stop()
			// Kirim sisa batch terakhir sebelum keluar
			if len(batch) > 0 {
				s.commitBatch(batch)
			}
			return // Keluar dari loop worker
		}
	}
}

func (s *activityLoggerService) commitBatch(batch []domain.ActivityLog) {
	log.Printf("Committing %d activity logs to database...", len(batch))

	// Buat context dengan timeout untuk operasi DB
	ctx, cancel := context.WithTimeout(context.Background(), dbInsertTimeout)
	defer cancel()

	if err := s.repo.CreateBatch(ctx, batch); err != nil {
		applogger.ErrorLogger.Printf("Failed to create activity log batch: %v", err)
		// Di aplikasi production, simpan log yang gagal
		// ke file untuk diproses ulang
	}
}
