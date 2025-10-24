package gorm

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type groupRepositoryGORM struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) usecase.GroupRepository {
	return &groupRepositoryGORM{db: db}
}

func (r *groupRepositoryGORM) Create(ctx context.Context, group *domain.Group) error {
	gormGroup := GroupFromDomain(group)
	result := r.db.WithContext(ctx).Create(gormGroup)
	if result.Error != nil {
		return result.Error
	}
	group.ID = gormGroup.ID
	return nil
}

func (r *groupRepositoryGORM) FindAll(ctx context.Context) ([]domain.Group, error) {
	var gormGroups []GroupGORM
	if err := r.db.WithContext(ctx).Preload("Class.Subject").Find(&gormGroups).Error; err != nil {
		return nil, err
	}

	var domainGroups []domain.Group
	for _, g := range gormGroups {
		domainGroups = append(domainGroups, *g.ToDomain())
	}
	return domainGroups, nil
}

func (r *groupRepositoryGORM) FindByID(ctx context.Context, id uint) (*domain.Group, error) {
	var gormGroup GroupGORM
	result := r.db.WithContext(ctx).Preload("Class.Subject").First(&gormGroup, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormGroup.ToDomain(), nil
}

func (r *groupRepositoryGORM) Update(ctx context.Context, group *domain.Group) error {
	gormGroup := GroupFromDomain(group)
	return r.db.WithContext(ctx).Save(gormGroup).Error
}

func (r *groupRepositoryGORM) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&GroupGORM{}, id).Error
}

func (r *groupRepositoryGORM) FindGroupsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Group, error) {
	var gormGroups []GroupGORM
	// Cari grup yang memiliki user dengan ID yang cocok di tabel relasi user_groups
	if err := r.db.WithContext(ctx).
		Joins("JOIN user_groups on user_groups.group_id = groups.id").
		Where("user_groups.user_id = ?", userID).
		Find(&gormGroups).Error; err != nil {
		return nil, err
	}

	var domainGroups []domain.Group
	for _, g := range gormGroups {
		domainGroups = append(domainGroups, *g.ToDomain())
	}
	return domainGroups, nil
}
