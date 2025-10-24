package gorm

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type mediaAssetRepositoryGORM struct {
	db *gorm.DB
}

func NewMediaAssetRepository(db *gorm.DB) usecase.MediaAssetRepository {
	return &mediaAssetRepositoryGORM{db: db}
}

func (r *mediaAssetRepositoryGORM) Create(ctx context.Context, asset *domain.MediaAsset) error {
	gormAsset := MediaAssetFromDomain(asset)
	result := r.db.WithContext(ctx).Create(gormAsset)
	if result.Error != nil {
		return result.Error
	}
	asset.ID = gormAsset.ID // Kembalikan ID yang baru dibuat
	return nil
}

func (r *mediaAssetRepositoryGORM) FindByID(ctx context.Context, id uuid.UUID) (*domain.MediaAsset, error) {
	var gormAsset MediaAssetGORM
	result := r.db.WithContext(ctx).First(&gormAsset, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormAsset.ToDomain(), nil
}

func (r *mediaAssetRepositoryGORM) SoftDelete(ctx context.Context, assetID uuid.UUID, deleterID *uuid.UUID) error {
	// Update manual agar bisa mengisi DeletedByUserID
	result := r.db.WithContext(ctx).
		Model(&MediaAssetGORM{}).
		Where("id = ?", assetID).
		Updates(map[string]interface{}{
			"deleted_at":         time.Now(),
			"deleted_by_user_id": deleterID,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("media asset not found or already deleted")
	}
	return nil
}
