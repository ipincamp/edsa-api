package gorm

import (
	"context"
	"errors"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type interactionRepositoryGORM struct {
	db *gorm.DB
}

func NewInteractionRepository(db *gorm.DB) usecase.InteractionRepository {
	return &interactionRepositoryGORM{db: db}
}

func (r *interactionRepositoryGORM) Create(ctx context.Context, interaction *domain.Interaction) error {
	gormInteraction := InteractionFromDomain(interaction)
	result := r.db.WithContext(ctx).Create(gormInteraction)
	if result.Error != nil {
		return result.Error
	}
	interaction.ID = gormInteraction.ID
	return nil
}

func (r *interactionRepositoryGORM) FindAllByPageID(ctx context.Context, pageID uint) ([]domain.Interaction, error) {
	var gormInteractions []InteractionGORM
	if err := r.db.WithContext(ctx).
		Preload("Page.Book"). // Nested preload
		Where("page_id = ?", pageID).
		Find(&gormInteractions).Error; err != nil {
		return nil, err
	}

	var domainInteractions []domain.Interaction
	for _, i := range gormInteractions {
		domainInteractions = append(domainInteractions, *i.ToDomain())
	}
	return domainInteractions, nil
}

func (r *interactionRepositoryGORM) FindByID(ctx context.Context, id uint) (*domain.Interaction, error) {
	var gormInteraction InteractionGORM
	result := r.db.WithContext(ctx).Preload("Page.Book").First(&gormInteraction, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormInteraction.ToDomain(), nil
}

func (r *interactionRepositoryGORM) Update(ctx context.Context, interaction *domain.Interaction) error {
	gormInteraction := InteractionFromDomain(interaction)
	return r.db.WithContext(ctx).Save(gormInteraction).Error
}

func (r *interactionRepositoryGORM) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&InteractionGORM{}, id).Error
}
