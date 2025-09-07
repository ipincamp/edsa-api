package repository

import (
	"context"

	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

// UserRepository adalah kontrak untuk akses data user
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindAllEmails(ctx context.Context) ([]string, error)
	List(ctx context.Context, limit, offset int, roleName string) ([]domain.User, int64, error)
	Update(ctx context.Context, user *domain.User) error
}

// userRepository adalah implementasi UserRepository menggunakan GORM
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository membuat instance baru userRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create menambah user baru ke database
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByEmail mencari user berdasarkan email
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Preload("Role").Where("email = ?", email).Take(&user).Error
	return &user, err
}

// FindByID mencari user berdasarkan ID
func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&user).Error
	return &user, err
}

// FindAllEmails mengambil semua email user
func (r *userRepository) FindAllEmails(ctx context.Context) ([]string, error) {
	var emails []string
	err := r.db.WithContext(ctx).Model(&domain.User{}).Pluck("email", &emails).Error
	return emails, err
}

// List mengambil daftar user dengan filter role dan paginasi
func (r *userRepository) List(ctx context.Context, limit, offset int, roleName string) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	// Build base query
	baseQuery := r.db.WithContext(ctx).
		Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.name != ?", constant.RoleAdmin)

	if roleName != "" {
		baseQuery = baseQuery.Where("roles.name = ?", roleName)
	}

	// Get total count
	err := baseQuery.Model(&domain.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get users with pagination - only select id, name, created_at
	err = baseQuery.
		Select("users.id, users.name, users.created_at").
		Limit(limit).
		Offset(offset).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Update mengubah data user di database
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
