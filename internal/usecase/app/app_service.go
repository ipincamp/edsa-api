package app

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type appService struct {
	bookRepo     usecase.BookRepository
	progressRepo usecase.UserBookProgressRepository
}

func NewAppService(
	bookRepo usecase.BookRepository,
	progressRepo usecase.UserBookProgressRepository,
) usecase.AppService {
	return &appService{
		bookRepo:     bookRepo,
		progressRepo: progressRepo,
	}
}

// GetBooksWithProgress mengambil semua buku dan menggabungkannya dengan progres user
func (s *appService) GetBooksWithProgress(ctx context.Context, userID uuid.UUID) ([]domain.BookResponse, error) {
	// 1. Ambil semua buku
	books, err := s.bookRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Ambil semua progres user
	progresses, err := s.progressRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. Buat map progres untuk lookup cepat
	progressMap := make(map[uint]*domain.UserBookProgress)
	for i, p := range progresses {
		progressMap[p.BookID] = &progresses[i]
	}

	// 4. Gabungkan data
	var bookResponses []domain.BookResponse
	for _, book := range books {
		resp := domain.BookResponse{
			ID:            book.ID,
			Title:         book.Title,
			Description:   book.Description,
			CoverImageURL: book.CoverImageURL,
			Theme:         book.Theme,
			BookOrder:     book.BookOrder,
		}

		// Cek apakah ada progres untuk buku ini
		if progress, ok := progressMap[book.ID]; ok {
			resp.Status = progress.Status
			resp.HighestScore = progress.HighestScore
		} else {
			// Jika tidak ada progres, atur status default
			// Asumsi: BookOrder 1 adalah buku pertama
			if book.BookOrder == 1 {
				resp.Status = domain.BookProgressStatusUnlocked
			} else {
				resp.Status = domain.BookProgressStatusLocked
			}
		}
		bookResponses = append(bookResponses, resp)
	}
	return bookResponses, nil
}

// GetProgressToRestore mengambil data untuk session restore
func (s *appService) GetProgressToRestore(ctx context.Context, userID uuid.UUID, bookID uint) (*domain.RestoreProgressResponse, error) {
	progress, err := s.progressRepo.FindByUserAndBook(ctx, userID, bookID)
	if err != nil {
		return nil, err
	}
	if progress == nil {
		// Jika belum ada progres, kembalikan nilai default
		return &domain.RestoreProgressResponse{
			LastPageID:           0,
			CurrentSessionPoints: 0,
		}, nil
	}

	return &domain.RestoreProgressResponse{
		LastPageID:           progress.LastPageID,
		CurrentSessionPoints: progress.CurrentSessionPoints,
	}, nil
}

// UpdatePageProgress dipanggil setiap pindah halaman
func (s *appService) UpdatePageProgress(ctx context.Context, userID uuid.UUID, req *domain.UpdateProgressRequest) error {
	progress := &domain.UserBookProgress{
		UserID: userID,
		BookID: req.BookID,
	}

	// 1. Cari atau buat entri progres
	if err := s.progressRepo.FindOrCreate(ctx, progress); err != nil {
		return err
	}

	// 2. Update progres
	progress.LastPageID = req.PageID
	// Logika penambahan poin (bisa juga 'set' total poin, tergantung kesepakatan)
	// Kontrak bilang 'points_earned_on_page', logika bilang 'menambah'
	progress.CurrentSessionPoints += req.PointsEarnedOnPage

	// 3. Simpan
	return s.progressRepo.Update(ctx, progress)
}

// CompleteBookProgress dipanggil saat buku selesai
func (s *appService) CompleteBookProgress(ctx context.Context, userID uuid.UUID, req *domain.CompleteProgressRequest) error {
	progress, err := s.progressRepo.FindByUserAndBook(ctx, userID, req.BookID)
	if err != nil {
		return err
	}
	if progress == nil {
		// Seharusnya tidak terjadi jika UpdatePageProgress sudah dipanggil
		return errors.New("progress not found for this book")
	}

	// 1. Update progres
	progress.Status = domain.BookProgressStatusCompleted
	if req.FinalScore > progress.HighestScore {
		progress.HighestScore = req.FinalScore // Simpan nilai tertinggi
	}
	progress.CurrentSessionPoints = 0 // Reset session

	if err := s.progressRepo.Update(ctx, progress); err != nil {
		return err
	}

	// 2. Buka buku berikutnya
	currentBook, err := s.bookRepo.FindByID(ctx, req.BookID)
	if err != nil {
		log.Printf("Warning: could not find current book (ID: %d) to unlock next: %v", req.BookID, err)
		return nil // Selesaikan progres, tapi gagal buka buku baru
	}
	if currentBook == nil {
		return nil // Buku tidak ada
	}

	nextBook, err := s.bookRepo.FindByOrder(ctx, currentBook.BookOrder+1)
	if err != nil {
		log.Printf("Warning: db error finding next book (Order: %d): %v", currentBook.BookOrder+1, err)
		return nil
	}
	if nextBook == nil {
		log.Println("User has completed the last book.")
		return nil // Ini adalah buku terakhir
	}

	// 3. Buat/Update progres buku berikutnya menjadi 'unlocked'
	nextProgress := &domain.UserBookProgress{
		UserID: userID,
		BookID: nextBook.ID,
	}
	if err := s.progressRepo.FindOrCreate(ctx, nextProgress); err != nil {
		log.Printf("Warning: failed to FindOrCreate progress for next book (ID: %d): %v", nextBook.ID, err)
		return nil
	}

	// Hanya unlock jika masih 'locked'
	if nextProgress.Status == domain.BookProgressStatusLocked {
		nextProgress.Status = domain.BookProgressStatusUnlocked
		if err := s.progressRepo.Update(ctx, nextProgress); err != nil {
			log.Printf("Warning: failed to update status for next book (ID: %d): %v", nextBook.ID, err)
		}
	}

	return nil
}
