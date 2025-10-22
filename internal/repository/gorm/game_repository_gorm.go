package gorm

import (
	"context"
	"errors"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type gameRepositoryGORM struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) usecase.GameRepository {
	return &gameRepositoryGORM{db: db}
}

func (r *gameRepositoryGORM) FindAll(ctx context.Context) ([]domain.Game, error) {
	var gormGames []GameGORM
	if err := r.db.WithContext(ctx).Find(&gormGames).Error; err != nil {
		return nil, err
	}

	var domainGames []domain.Game
	for _, g := range gormGames {
		domainGames = append(domainGames, *g.ToDomain())
	}
	return domainGames, nil
}

func (r *gameRepositoryGORM) FindByID(ctx context.Context, id uint) (*domain.Game, error) {
	var gormGame GameGORM
	result := r.db.WithContext(ctx).First(&gormGame, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return gormGame.ToDomain(), nil
}
