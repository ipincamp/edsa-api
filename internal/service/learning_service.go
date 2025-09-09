package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ipincamp/go-edsa-api/internal/delivery/http/dto"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/repository"
	"gorm.io/gorm"
)

// Definisikan error baru
var (
	ErrInteractionNotFound  = errors.New("interaction not found")
	ErrBookAlreadyCompleted = errors.New("book is already completed")
	ErrInvalidAnswer        = errors.New("invalid answer format")
)

// LearningService adalah kontrak untuk logika bisnis terkait pembelajaran
type LearningService interface {
	GetAvailableBooks(ctx context.Context, userID string) ([]dto.BookResponse, error)
	GetBookDetail(ctx context.Context, userID string, bookID string) (*dto.BookDetailResponse, error)
	SubmitInteraction(ctx context.Context, userID string, req dto.SubmitInteractionRequest) error
}

type learningService struct {
	db               *gorm.DB
	bookRepo         repository.BookRepository
	userProgressRepo repository.UserProgressRepository
	interactionRepo  repository.InteractionRepository
	pageRepo         repository.PageRepository
}

// NewLearningService membuat instance baru LearningService
func NewLearningService(
	db *gorm.DB,
	bookRepo repository.BookRepository,
	userProgressRepo repository.UserProgressRepository,
	interactionRepo repository.InteractionRepository,
	pageRepo repository.PageRepository,
) LearningService {
	return &learningService{
		db:               db,
		bookRepo:         bookRepo,
		userProgressRepo: userProgressRepo,
		interactionRepo:  interactionRepo,
		pageRepo:         pageRepo,
	}
}

// GetAvailableBooks mengambil semua buku dan menentukan statusnya untuk user tertentu
func (s *learningService) GetAvailableBooks(ctx context.Context, userID string) ([]dto.BookResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	type result struct {
		resp []dto.BookResponse
		err  error
	}
	resultChan := make(chan result, 1)

	go func() {
		allBooks, err := s.bookRepo.FindAllOrderedByLevel(ctx)
		if err != nil {
			resultChan <- result{nil, err}
			return
		}

		userProgresses, err := s.userProgressRepo.GetProgressByUser(ctx, userID)
		if err != nil {
			resultChan <- result{nil, err}
			return
		}

		// Buat map untuk lookup progres user dengan cepat
		progressMap := make(map[string]domain.UserProgress)
		for _, p := range userProgresses {
			progressMap[p.BookID] = p
		}

		var response []dto.BookResponse
		for _, book := range allBooks {
			progress, hasProgress := progressMap[book.ID]

			// Tentukan status buku (bisa dibuka atau tidak)
			var bookStatus dto.BookStatus
			isUnlocked := book.UnlockDependencyID == nil || (progressMap[*book.UnlockDependencyID].Status == "completed")

			if !isUnlocked {
				bookStatus = dto.StatusLocked
			} else if hasProgress && progress.Status == "completed" {
				bookStatus = dto.StatusCompleted
			} else if hasProgress && progress.Status == "in_progress" {
				bookStatus = dto.StatusInProgress
			} else {
				bookStatus = dto.StatusAvailable
			}

			response = append(response, dto.BookResponse{
				ID:            book.ID,
				Title:         book.Title,
				Description:   book.Description,
				CoverImageURL: book.CoverImageURL,
				Level:         book.Level,
				Status:        bookStatus,
			})
		}
		resultChan <- result{response, nil}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

// GetBookDetail mengambil detail sebuah buku beserta halaman dan progres user
func (s *learningService) GetBookDetail(ctx context.Context, userID string, bookID string) (*dto.BookDetailResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	type result struct {
		resp *dto.BookDetailResponse
		err  error
	}
	resultChan := make(chan result, 1)

	go func() {
		book, err := s.bookRepo.FindByID(ctx, bookID)
		if err != nil {
			resultChan <- result{nil, err}
			return
		}

		// Cari atau buat progres user untuk buku ini
		progress := &domain.UserProgress{
			UserID: userID,
			BookID: bookID,
			Status: "in_progress", // Jika baru dibuat, statusnya langsung 'in_progress'
		}
		if err := s.userProgressRepo.FindOrCreate(ctx, progress); err != nil {
			resultChan <- result{nil, err}
			return
		}

		response := &dto.BookDetailResponse{
			ID:                  book.ID,
			Title:               book.Title,
			UserProgressStatus:  progress.Status,
			LastCompletedPageID: progress.LastCompletedPageID,
			Pages:               dto.ToPageListResponse(book.Pages),
		}
		resultChan <- result{response, nil}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

// SubmitInteraction memproses jawaban user, menghitung skor, dan memperbarui progres
func (s *learningService) SubmitInteraction(ctx context.Context, userID string, req dto.SubmitInteractionRequest) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	type result struct {
		err error
	}
	resultChan := make(chan result, 1)

	go func() {
		err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			// Buat instance repo dengan transaksi
			txProgressRepo := repository.NewUserProgressRepository(tx)
			txInteractionRepo := repository.NewInteractionRepository(tx)
			txPageRepo := repository.NewPageRepository(tx)

			// 1. Ambil progres user saat ini
			progress, err := txProgressRepo.FindByUserAndBook(ctx, userID, req.BookID)
			if err != nil {
				return err
			}
			if progress.Status == "completed" {
				return ErrBookAlreadyCompleted
			}

			// 2. Ambil data interaksi untuk validasi dan scoring
			interaction, err := txInteractionRepo.FindByID(ctx, req.InteractionID)
			if err != nil {
				return ErrInteractionNotFound
			}

			// 3. Validasi & Hitung Skor
			isCorrect, err := checkAnswer(req.Answer, interaction.ConfigData)
			if err != nil {
				return err
			}

			scoreToAdd := 0
			if isCorrect {
				scoreToAdd = interaction.PointsAwarded
			}

			// 4. Update Progres
			progress.BookScore += scoreToAdd
			progress.LastCompletedPageID = &req.PageID

			// 5. Cek apakah buku selesai
			isLast, err := txPageRepo.IsLastPage(ctx, req.PageID, req.BookID)
			if err != nil {
				return err
			}
			if isLast && isCorrect { // Buku dianggap selesai jika interaksi terakhir benar
				now := time.Now()
				progress.Status = "completed"
				progress.CompletedAt = &now
			}

			// 6. Simpan perubahan progres
			return txProgressRepo.Update(ctx, progress)
		})
		resultChan <- result{err: err}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case res := <-resultChan:
		return res.err
	}
}

// checkAnswer adalah helper untuk memvalidasi jawaban user.
// Logika di sini bisa sangat kompleks tergantung tipe interaksi.
func checkAnswer(userAnswer string, configData json.RawMessage) (bool, error) {
	// Contoh sederhana jika jawaban disimpan dalam format {"correct_answer": "the_answer"}
	var config struct {
		CorrectAnswer string `json:"correct_answer"`
	}
	if err := json.Unmarshal(configData, &config); err != nil {
		return false, ErrInvalidAnswer
	}

	return userAnswer == config.CorrectAnswer, nil
}
