package gorm

import (
	"context"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type activityLogRepositoryGORM struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) usecase.ActivityLogRepository {
	return &activityLogRepositoryGORM{db: db}
}

func (r *activityLogRepositoryGORM) Create(ctx context.Context, log *domain.ActivityLog) error {
	gormLog := ActivityLogFromDomain(log)
	result := r.db.WithContext(ctx).Create(gormLog)
	if result.Error != nil {
		return result.Error
	}
	log.ID = gormLog.ID
	return nil
}

func (r *activityLogRepositoryGORM) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.ActivityLog, error) {
	var gormLogs []ActivityLogGORM
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("timestamp_start desc"). // Tampilkan yang terbaru dulu
		Find(&gormLogs).Error; err != nil {
		return nil, err
	}

	var domainLogs []domain.ActivityLog
	for _, l := range gormLogs {
		domainLogs = append(domainLogs, *l.ToDomain())
	}
	return domainLogs, nil
}
