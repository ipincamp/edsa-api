package dashboard

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/pkg/applogger"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type dashboardService struct {
	activityLogRepo      usecase.ActivityLogRepository
	userRepo             usecase.UserRepository
	groupRepo            usecase.GroupRepository
	bookRepo             usecase.BookRepository
	groupBookSettingRepo usecase.GroupBookSettingRepository
}

func NewDashboardService(
	activityLogRepo usecase.ActivityLogRepository,
	userRepo usecase.UserRepository,
	groupRepo usecase.GroupRepository,
	bookRepo usecase.BookRepository,
	groupBookSettingRepo usecase.GroupBookSettingRepository,
) usecase.DashboardService {
	return &dashboardService{
		activityLogRepo:      activityLogRepo,
		userRepo:             userRepo,
		groupRepo:            groupRepo,
		bookRepo:             bookRepo,
		groupBookSettingRepo: groupBookSettingRepo,
	}
}

// --- Mapper DTO ---
// (Helper untuk konversi domain ke response DTO)

func toActivityLogResponse(l *domain.ActivityLog) *domain.ActivityLogResponse {
	return &domain.ActivityLogResponse{
		ID:             l.ID,
		Action:         l.Action,
		TimestampStart: l.TimestampStart,
		DurationMs:     l.DurationMs,
		Details:        l.Details,
	}
}

// --- Methods ---

func (s *dashboardService) GetStudentActivity(ctx context.Context, studentID uuid.UUID, filters *domain.ActivityLogQuery) (*domain.PaginatedDTO, error) {
	// 1. Validasi apakah studentID ada
	student, err := s.userRepo.FindByID(ctx, studentID)
	if err != nil {
		applogger.ErrorLogger.Printf("GetStudentActivity: DB error checking student %s: %v", studentID, err)
		return nil, errors.New("database error checking student")
	}
	if student == nil {
		return nil, errors.New("student not found")
	}

	// 2. Set default paginasi
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.Limit <= 0 {
		filters.Limit = 10 // Default limit
	}

	// 3. Ambil data paginasi dari repo
	paginatedResult, err := s.activityLogRepo.FindPaginatedByUserID(ctx, studentID, filters)
	if err != nil {
		applogger.ErrorLogger.Printf("GetStudentActivity: DB error fetching paginated logs for student %s: %v", studentID, err)
		return nil, errors.New("database error fetching logs")
	}

	// 4. Map ke DTO Response
	var responses []domain.ActivityLogResponse
	for _, log := range paginatedResult.Logs {
		responses = append(responses, *toActivityLogResponse(&log))
	}

	// 5. Hitung metadata paginasi
	totalPage := int64(math.Ceil(float64(paginatedResult.TotalData) / float64(filters.Limit)))
	if totalPage == 0 && paginatedResult.TotalData > 0 {
		totalPage = 1
	}

	// 6. Buat DTO respons paginasi
	return &domain.PaginatedDTO{
		List: responses,
		Meta: domain.PaginationMetaDTO{
			Page:      filters.Page,
			Limit:     filters.Limit,
			TotalPage: totalPage,
			TotalData: paginatedResult.TotalData,
		},
	}, nil
}

func (s *dashboardService) UnlockBookForGroup(ctx context.Context, groupID uint, bookID uint) error {
	// 1. Validasi apakah groupID ada
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		applogger.ErrorLogger.Printf("UnlockBookForGroup: DB error checking group %d: %v", groupID, err)
		return errors.New("database error checking group")
	}
	if group == nil {
		return errors.New("group not found")
	}

	// 2. Validasi apakah bookID ada
	book, err := s.bookRepo.FindByID(ctx, bookID)
	if err != nil {
		applogger.ErrorLogger.Printf("UnlockBookForGroup: DB error checking book %d: %v", bookID, err)
		return errors.New("database error checking book")
	}
	if book == nil {
		return errors.New("book not found")
	}

	// 3. Buat entitas domain untuk di-upsert
	setting := &domain.GroupBookSetting{
		GroupID:    groupID,
		BookID:     bookID,
		IsUnlocked: true,
		UpdatedAt:  time.Now(), // Mapper akan menggunakannya
	}

	// 4. Panggil repository untuk upsert
	if err := s.groupBookSettingRepo.Upsert(ctx, setting); err != nil {
		applogger.ErrorLogger.Printf("UnlockBookForGroup: Failed to upsert setting for group %d and book %d: %v", groupID, bookID, err)
		return errors.New("failed to update book setting")
	}

	return nil
}
