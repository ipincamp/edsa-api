package repository

import (
	"context"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

type InteractionRepository interface {
	FindByID(ctx context.Context, id string) (*domain.Interaction, error)
}

type interactionRepository struct {
	db *gorm.DB
}

func NewInteractionRepository(db *gorm.DB) InteractionRepository {
	return &interactionRepository{db: db}
}

func (r *interactionRepository) FindByID(ctx context.Context, id string) (*domain.Interaction, error) {
	var interaction domain.Interaction
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&interaction).Error
	return &interaction, err
}
