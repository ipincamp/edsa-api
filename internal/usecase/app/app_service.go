package app

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/pkg/applogger"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type appService struct {
	bookRepo             usecase.BookRepository
	progressRepo         usecase.UserBookProgressRepository
	gameRepo             usecase.GameRepository
	userGameScoreRepo    usecase.UserGameScoreRepository
	logger               usecase.ActivityLoggerService
	userRepo             usecase.UserRepository
	groupRepo            usecase.GroupRepository
	groupBookSettingRepo usecase.GroupBookSettingRepository
}

func NewAppService(
	bookRepo usecase.BookRepository,
	progressRepo usecase.UserBookProgressRepository,
	gameRepo usecase.GameRepository,
	userGameScoreRepo usecase.UserGameScoreRepository,
	logger usecase.ActivityLoggerService,
	userRepo usecase.UserRepository,
	groupRepo usecase.GroupRepository,
	groupBookSettingRepo usecase.GroupBookSettingRepository,
) usecase.AppService {
	return &appService{
		bookRepo:             bookRepo,
		progressRepo:         progressRepo,
		gameRepo:             gameRepo,
		userGameScoreRepo:    userGameScoreRepo,
		logger:               logger,
		userRepo:             userRepo,
		groupRepo:            groupRepo,
		groupBookSettingRepo: groupBookSettingRepo,
	}
}

// GetBooksWithProgress mengambil semua buku dan menggabungkannya dengan progres user
func (s *appService) GetBooksWithProgress(ctx context.Context, userID uuid.UUID) ([]domain.BookResponse, error) {
	// 1. Ambil semua buku
	books, err := s.bookRepo.FindAll(ctx)
	if err != nil {
		applogger.ErrorLogger.Printf("GetBooksWithProgress: Failed to FindAll books: %v", err)
		return nil, err
	}

	// 2. Ambil semua progres user
	progresses, err := s.progressRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		applogger.ErrorLogger.Printf("GetBooksWithProgress: Failed to FindAllByUserID progress for user %s: %v", userID, err)
		return nil, err
	}

	// 3. Buat map progres untuk lookup cepat
	progressMap := make(map[uint]*domain.UserBookProgress)
	for i, p := range progresses {
		progressMap[p.BookID] = &progresses[i]
	}

	// 4. Ambil data user untuk cek role
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		applogger.ErrorLogger.Printf("GetBooksWithProgress: Failed to FindByID user %s: %v", userID, err)
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	var bookResponses []domain.BookResponse

	// 5. Terapkan logika berdasarkan role
	if user.Role.Name == domain.RoleNameStudent {
		// --- LOGIKA BARU UNTUK STUDENT ---

		// 5a. Ambil grup pengguna
		groups, err := s.groupRepo.FindGroupsByUserID(ctx, userID)
		if err != nil {
			applogger.ErrorLogger.Printf("GetBooksWithProgress: Failed to FindGroupsByUserID for student %s: %v", userID, err)
			return nil, err
		}

		var groupIDs []uint
		for _, g := range groups {
			groupIDs = append(groupIDs, g.ID)
		}

		// 5b. Ambil pengaturan buku untuk grup tsb
		groupUnlockMap := make(map[uint]bool)
		if len(groupIDs) > 0 {
			settings, err := s.groupBookSettingRepo.FindSettingsByGroupIDs(ctx, groupIDs)
			if err != nil {
				applogger.ErrorLogger.Printf("GetBooksWithProgress: Failed to FindSettingsByGroupIDs for groups %v: %v", groupIDs, err)
				return nil, err
			}
			// Buat map status unlock dari guru
			for _, s := range settings {
				if s.IsUnlocked {
					groupUnlockMap[s.BookID] = true
				}
			}
		}

		// 5c. Gabungkan pengaturan grup dengan UserBookProgress
		for _, book := range books {
			resp := domain.BookResponse{
				ID:            book.ID,
				Title:         book.Title,
				Description:   book.Description,
				CoverImageURL: book.CoverImageURL,
				Theme:         book.Theme,
				BookOrder:     book.BookOrder,
			}

			isUnlockedByGroup := groupUnlockMap[book.ID]
			progress, hasProgress := progressMap[book.ID]

			if hasProgress && progress.Status == domain.BookProgressStatusCompleted {
				// Jika sudah selesai, statusnya 'completed'
				resp.Status = domain.BookProgressStatusCompleted
				resp.HighestScore = progress.HighestScore
			} else if isUnlockedByGroup {
				// Jika di-unlock oleh guru (dan belum selesai)
				resp.Status = domain.BookProgressStatusUnlocked
				if hasProgress {
					resp.HighestScore = progress.HighestScore
				}
			} else {
				// Jika tidak di-unlock dan tidak selesai
				resp.Status = domain.BookProgressStatusLocked
			}

			bookResponses = append(bookResponses, resp)
		}

	} else {
		// --- LOGIKA LAMA (UNTUK PUBLIC, ADMIN, TEACHER) ---
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
	}

	return bookResponses, nil
}

// GetProgressToRestore mengambil data untuk session restore
func (s *appService) GetProgressToRestore(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, bookID uint) (*domain.RestoreProgressResponse, error) {
	progress, err := s.progressRepo.FindByUserAndBook(ctx, userID, bookID)
	if err != nil {
		return nil, err
	}

	// Logging aktivitas memulai/melanjutkan buku
	details, _ := json.Marshal(map[string]interface{}{"book_id": bookID})
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    userID,
		Action:    domain.ActionStartBook,
		Details:   details,
		SessionID: sessionID,
	})

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
func (s *appService) UpdatePageProgress(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, req *domain.UpdateProgressRequest) error {
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

	// Catatan: Fungsi ini sendiri tidak me-log, jadi hanya perlu menerima parameter
	// Jika Anda *ingin* me-log setiap pindah halaman, tambahkan panggilannya di sini:
	details, _ := json.Marshal(map[string]interface{}{"book_id": req.BookID, "page_id": req.PageID})
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    userID,
		Action:    domain.ActionPageView,
		SessionID: sessionID,
		Details:   details,
	})

	// 3. Simpan
	return s.progressRepo.Update(ctx, progress)
}

// CompleteBookProgress dipanggil saat buku selesai
func (s *appService) CompleteBookProgress(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, req *domain.CompleteProgressRequest) error {
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

	// Log aktivitas menyelesaikan buku
	details, _ := json.Marshal(map[string]interface{}{
		"book_id": req.BookID,
		"score":   req.FinalScore,
	})
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:     userID,
		Action:     domain.ActionCompleteBook,
		DurationMs: &req.DurationMs,
		Details:    details,
		SessionID:  sessionID,
	})

	// 2. Buka buku berikutnya
	currentBook, err := s.bookRepo.FindByID(ctx, req.BookID)
	if err != nil {
		applogger.ErrorLogger.Printf("CompleteBookProgress: could not find current book (ID: %d) to unlock next: %v", req.BookID, err)
		return nil // Selesaikan progres, tapi gagal buka buku baru
	}
	if currentBook == nil {
		return nil // Buku tidak ada
	}

	// Cek role user sebelum unlock buku berikutnya
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		applogger.ErrorLogger.Printf("CompleteBookProgress: could not find user %s: %v", userID, err)
		return nil
	}

	// Hanya lakukan auto-unlock buku berikutnya untuk role PUBLIC
	if user != nil && user.Role.Name == domain.RoleNamePublic {
		nextBook, err := s.bookRepo.FindByOrder(ctx, currentBook.BookOrder+1)
		if err != nil {
			applogger.ErrorLogger.Printf("CompleteBookProgress: db error finding next book (Order: %d): %v", currentBook.BookOrder+1, err)
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
	}
	// Jika role-nya STUDENT, buku berikutnya TIDAK otomatis unlock.
	// Harus di-unlock oleh guru via dashboard.

	return nil
}

func (s *appService) GetAllGames(ctx context.Context, userID uuid.UUID) ([]domain.GameResponse, error) {
	// 1. Ambil semua master game
	games, err := s.gameRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Ambil semua skor game milik user
	scores, err := s.userGameScoreRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. Buat map skor untuk lookup
	scoreMap := make(map[uint]int)
	for _, score := range scores {
		scoreMap[score.GameID] = score.HighestScore
	}

	// 4. Gabungkan data
	var gameResponses []domain.GameResponse
	for _, game := range games {
		resp := domain.GameResponse{
			ID:               game.ID,
			Name:             game.Name,
			Type:             game.Type,
			RelatedBookTheme: game.RelatedBookTheme,
			HighestScore:     0, // Default 0
		}

		// Jika user punya skor, timpa nilainya
		if score, ok := scoreMap[game.ID]; ok {
			resp.HighestScore = score
		}

		gameResponses = append(gameResponses, resp)
	}
	return gameResponses, nil
}

func (s *appService) SubmitGameScore(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, gameID uint, req *domain.SubmitGameScoreRequest) error {
	// 1. Cek apakah gameID valid
	game, err := s.gameRepo.FindByID(ctx, gameID)
	if err != nil {
		return err
	}
	if game == nil {
		return errors.New("game not found")
	}

	// 2. Dapatkan skor yang ada (jika ada)
	existingScore, err := s.userGameScoreRepo.FindByUserAndGame(ctx, userID, gameID)
	if err != nil {
		return err
	}

	newScore := req.Score

	// 3. Simpan nilai tertinggi
	if existingScore != nil {
		if newScore <= existingScore.HighestScore {
			// Skor baru tidak lebih tinggi, tidak perlu update
			return nil
		}
		// Skor baru lebih tinggi, update skor yang ada
		existingScore.HighestScore = newScore
		existingScore.UpdatedAt = time.Now()
		return s.userGameScoreRepo.Upsert(ctx, existingScore)
	}

	// 4. Jika belum ada skor, buat entri baru
	newScoreEntry := &domain.UserGameScore{
		UserID:       userID,
		GameID:       gameID,
		HighestScore: newScore,
	}

	// Log aktivitas menyelesaikan game
	details, _ := json.Marshal(map[string]interface{}{
		"game_id": gameID,
		"score":   req.Score,
	})
	s.logger.Log(ctx, domain.ActivityLog{
		UserID:    userID,
		Action:    domain.ActionCompleteGame,
		SessionID: sessionID,
		Details:   details,
	})

	return s.userGameScoreRepo.Upsert(ctx, newScoreEntry)
}
