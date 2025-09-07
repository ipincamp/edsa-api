package repository

import (
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
	FindAllEmails() ([]string, error)
	List(limit, offset int, roleName string) ([]domain.User, int64, error)
	Update(user *domain.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *userRepository) FindAllEmails() ([]string, error) {
	var emails []string
	err := r.db.Model(&domain.User{}).Pluck("email", &emails).Error
	return emails, err
}

func (r *userRepository) List(limit, offset int, roleName string) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	query := r.db.Preload("Role")

	if roleName != "" {
		query = query.Joins("JOIN roles ON users.role_id = roles.id").Where("roles.name = ?", roleName)
	}

	if err := query.Model(&domain.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(limit).Offset(offset).Find(&users).Error
	return users, total, err
}

func (r *userRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}
