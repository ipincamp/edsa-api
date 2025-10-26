package gorm

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/internal/usecase"
	"gorm.io/gorm"
)

type userRepositoryGORM struct {
	db       *gorm.DB
	roleRepo usecase.RoleRepository
}

// Modifikasi NewUserRepository untuk menerima roleRepo
func NewUserRepository(db *gorm.DB, roleRepo usecase.RoleRepository) usecase.UserRepository {
	return &userRepositoryGORM{
		db:       db,
		roleRepo: roleRepo,
	}
}

func (r *userRepositoryGORM) Create(ctx context.Context, user *domain.User) error {
	gormUser := UserFromDomain(user)

	result := r.db.WithContext(ctx).Create(gormUser)
	if result.Error != nil {
		return result.Error
	}

	user.ID = gormUser.ID
	user.CreatedAt = gormUser.CreatedAt
	user.UpdatedAt = gormUser.UpdatedAt

	return nil
}

func (r *userRepositoryGORM) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var gormUser UserGORM
	result := r.db.WithContext(ctx).
		Preload("ProfilePicture").
		Where("email = ?", email).
		First(&gormUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// Konversi ke domain (tanpa role)
	domainUser := gormUser.ToDomain()

	// Ambil role dari cache
	role, err := r.roleRepo.FindByID(ctx, domainUser.RoleID)
	if err != nil {
		// Ini adalah error cache/DB, bukan error "tidak ditemukan"
		return nil, errors.New("database error: failed to find role for user")
	}
	if role == nil {
		// Ini adalah masalah integritas data
		return nil, errors.New("data integrity error: user role not found")
	}

	domainUser.Role = *role // Set role dari cache
	return domainUser, nil
}

func (r *userRepositoryGORM) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var gormUser UserGORM
	result := r.db.WithContext(ctx).
		Preload("ProfilePicture").
		First(&gormUser, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	// Konversi ke domain (tanpa role)
	domainUser := gormUser.ToDomain()

	// Ambil role dari cache
	role, err := r.roleRepo.FindByID(ctx, domainUser.RoleID)
	if err != nil {
		// Ini adalah error cache/DB, bukan error "tidak ditemukan"
		return nil, errors.New("database error: failed to find role for user")
	}
	if role == nil {
		// Ini adalah masalah integritas data
		return nil, errors.New("data integrity error: user role not found")
	}

	domainUser.Role = *role // Set role dari cache
	return domainUser, nil
}

func (r *userRepositoryGORM) Update(ctx context.Context, user *domain.User) error {
	gormUser := UserFromDomain(user)
	// Save akan memperbarui semua kolom, termasuk password yang diubah
	return r.db.WithContext(ctx).Save(gormUser).Error
}

func (r *userRepositoryGORM) Delete(ctx context.Context, userID uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&UserGORM{}, userID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found or already deleted")
	}
	return nil
}
