package gorm

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userGameScoreRepositoryGORM struct {
	db *gorm.DB
}

func NewUserGameScoreRepository(db *gorm.DB) usecase.UserGameScoreRepository {
	return &userGameScoreRepositoryGORM{db: db}
}

func (r *userGameScoreRepositoryGORM) FindByUserAndGame(ctx context.Context, userID uuid.UUID, gameID uint) (*domain.UserGameScore, error) {
	var gormScore UserGameScoreGORM
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND game_id = ?", userID, gameID).
		First(&gormScore)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Tidak ditemukan
		}
		return nil, result.Error
	}
	return gormScore.ToDomain(), nil
}

func (r *userGameScoreRepositoryGORM) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]domain.UserGameScore, error) {
	var gormScores []UserGameScoreGORM
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&gormScores).Error; err != nil {
		return nil, err
	}

	var domainScores []domain.UserGameScore
	for _, s := range gormScores {
		domainScores = append(domainScores, *s.ToDomain())
	}
	return domainScores, nil
}

func (r *userGameScoreRepositoryGORM) Upsert(ctx context.Context, score *domain.UserGameScore) error {
	gormScore := UserGameScoreFromDomain(score)

	// Gunakan OnConflict untuk Create atau Update
	// Jika (user_id, game_id) sudah ada, update 'highest_score'
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "game_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"highest_score", "updated_at"}),
		}).
		Create(gormScore).Error
}
