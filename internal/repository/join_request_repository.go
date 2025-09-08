package repository

import (
	"context"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

// JoinRequestRepository adalah kontrak untuk akses data permintaan join grup
type JoinRequestRepository interface {
	Create(ctx context.Context, request *domain.JoinGroupRequest) error
	UpdateStatus(ctx context.Context, requestID string, status string) error
	GetPendingByGroupID(ctx context.Context, groupID string) ([]domain.JoinGroupRequest, error)
	FindByID(ctx context.Context, id string) (*domain.JoinGroupRequest, error)
}

type joinRequestRepository struct {
	db *gorm.DB
}

// NewJoinRequestRepository membuat instance baru
func NewJoinRequestRepository(db *gorm.DB) JoinRequestRepository {
	return &joinRequestRepository{db: db}
}

func (r *joinRequestRepository) Create(ctx context.Context, request *domain.JoinGroupRequest) error {
	return r.db.WithContext(ctx).Create(request).Error
}

func (r *joinRequestRepository) UpdateStatus(ctx context.Context, requestID string, status string) error {
	return r.db.WithContext(ctx).Model(&domain.JoinGroupRequest{}).Where("id = ?", requestID).Update("status", status).Error
}

// GetPendingByGroupID mendapatkan semua request yang masih pending untuk sebuah kelas
func (r *joinRequestRepository) GetPendingByGroupID(ctx context.Context, groupID string) ([]domain.JoinGroupRequest, error) {
	var requests []domain.JoinGroupRequest
	err := r.db.WithContext(ctx).
		Preload("Applicant"). // Memuat data user yang mengajukan
		Where("target_group_id = ? AND status = ?", groupID, "pending").
		Find(&requests).Error
	return requests, err
}

func (r *joinRequestRepository) FindByID(ctx context.Context, id string) (*domain.JoinGroupRequest, error) {
	var request domain.JoinGroupRequest
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&request).Error
	return &request, err
}
