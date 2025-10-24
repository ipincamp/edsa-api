package gorm

import (
	"context"
	"time"

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

func (r *activityLogRepositoryGORM) CreateBatch(ctx context.Context, logs []domain.ActivityLog) error {
	if len(logs) == 0 {
		return nil
	}

	// Konversi domain slice ke gorm slice
	gormLogs := make([]ActivityLogGORM, len(logs))
	for i, logEntry := range logs {
		// Kita butuh pointer ke logEntry untuk fungsi mapper
		entry := logEntry
		gormLogs[i] = *ActivityLogFromDomain(&entry)
	}

	// GORM's Create() secara otomatis melakukan bulk insert jika diberi slice
	result := r.db.WithContext(ctx).Create(&gormLogs)
	if result.Error != nil {
		return result.Error
	}
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

func (r *activityLogRepositoryGORM) FindPaginatedByUserID(ctx context.Context, userID uuid.UUID, filters *domain.ActivityLogQuery) (*domain.PaginatedActivityLogs, error) {
	var gormLogs []ActivityLogGORM
	var totalData int64

	// 1. Buat query dasar
	query := r.db.WithContext(ctx).Model(&ActivityLogGORM{}).Where("user_id = ?", userID)

	// 2. Terapkan filter tanggal (inklusif)
	if filters.StartDate != "" {
		// Parsing YYYY-MM-DD
		start, err := time.Parse("2006-01-02", filters.StartDate)
		if err == nil {
			// Mengambil data dari awal hari (00:00:00)
			query = query.Where("timestamp_start >= ?", start)
		}
	}
	if filters.EndDate != "" {
		// Parsing YYYY-MM-DD
		end, err := time.Parse("2006-01-02", filters.EndDate)
		if err == nil {
			// Mengambil data sampai akhir hari (23:59:59)
			query = query.Where("timestamp_start <= ?", end.Add(24*time.Hour-time.Nanosecond))
		}
	}

	// 3. Dapatkan total data (sebelum limit/offset)
	if err := query.Count(&totalData).Error; err != nil {
		return nil, err
	}

	if totalData == 0 {
		return &domain.PaginatedActivityLogs{
			Logs:      []domain.ActivityLog{},
			TotalData: 0,
		}, nil
	}

	// 4. Hitung offset dan terapkan limit/offset
	offset := (filters.Page - 1) * filters.Limit
	if err := query.
		Order("timestamp_start desc").
		Limit(filters.Limit).
		Offset(offset).
		Find(&gormLogs).Error; err != nil {
		return nil, err
	}

	// 5. Konversi GORM ke Domain
	domainLogs := make([]domain.ActivityLog, len(gormLogs))
	for i, l := range gormLogs {
		domainLogs[i] = *l.ToDomain()
	}

	return &domain.PaginatedActivityLogs{
		Logs:      domainLogs,
		TotalData: totalData,
	}, nil
}
