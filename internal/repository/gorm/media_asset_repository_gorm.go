package gorm

import (
	"context"

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
