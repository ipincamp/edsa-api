package dashboard

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
)

type dashboardService struct {
	activityLogRepo usecase.ActivityLogRepository
	userRepo        usecase.UserRepository
}

func NewDashboardService(
	activityLogRepo usecase.ActivityLogRepository,
	userRepo usecase.UserRepository,
) usecase.DashboardService {
	return &dashboardService{
		activityLogRepo: activityLogRepo,
		userRepo:        userRepo,
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

func (s *dashboardService) GetStudentActivity(ctx context.Context, studentID uuid.UUID) ([]domain.ActivityLogResponse, error) {
	// TODO: Otorisasi - Cek apakah guru yang meminta berhak melihat studentID ini.

	// 1. Validasi apakah studentID ada
	student, err := s.userRepo.FindByID(ctx, studentID)
	if err != nil {
		return nil, errors.New("database error checking student")
	}
	if student == nil {
		return nil, errors.New("student not found")
	}

	// 2. Ambil semua log untuk user tersebut
	logs, err := s.activityLogRepo.FindAllByUserID(ctx, studentID)
	if err != nil {
		return nil, errors.New("database error fetching logs")
	}

	// 3. Map ke DTO Response
	var responses []domain.ActivityLogResponse
	for _, log := range logs {
		responses = append(responses, *toActivityLogResponse(&log))
	}

	return responses, nil
}
