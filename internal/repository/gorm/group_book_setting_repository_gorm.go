package gorm

import (
	"context"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type groupBookSettingRepositoryGORM struct {
	db *gorm.DB
}

// NewGroupBookSettingRepository membuat instance baru dari repository pengaturan buku grup
func NewGroupBookSettingRepository(db *gorm.DB) usecase.GroupBookSettingRepository {
	return &groupBookSettingRepositoryGORM{db: db}
}

// Upsert membuat atau memperbarui pengaturan buku (OnConflict)
func (r *groupBookSettingRepositoryGORM) Upsert(ctx context.Context, setting *domain.GroupBookSetting) error {
	gormSetting := GroupBookSettingFromDomain(setting)

	// Gunakan OnConflict untuk Create atau Update
	// Jika (group_id, book_id) sudah ada, update 'is_unlocked' dan 'updated_at'
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "group_id"}, {Name: "book_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"is_unlocked", "updated_at"}),
		}).
		Create(gormSetting).Error
}
