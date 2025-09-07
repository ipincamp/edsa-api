package seeders

import (
	"fmt"
	"log"

	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"github.com/ipincamp/go-edsa-api/pkg/hash"
	"gorm.io/gorm"
)

// UserSeeder melakukan seed data user admin ke database
func UserSeeder(cnf *config.Config) func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		adminRole, err := findAdminRole(db)
		if err != nil {
			return err
		}
		adminUser, err := buildAdminUser(cnf, adminRole.ID)
		if err != nil {
			return err
		}
		if err := db.FirstOrCreate(&adminUser, "email = ?", adminUser.Email).Error; err != nil {
			return fmt.Errorf("failed to seed admin user: %w", err)
		}
		log.Println("User seeder ran successfully")
		return nil
	}
}

// findAdminRole mencari role admin di database
func findAdminRole(db *gorm.DB) (*domain.Role, error) {
	var adminRole domain.Role
	if err := db.Where("name = ?", constant.RoleAdmin).First(&adminRole).Error; err != nil {
		return nil, fmt.Errorf("admin role not found for user seeder: %w", err)
	}
	return &adminRole, nil
}

// buildAdminUser membangun data user admin dari config dan roleID
func buildAdminUser(cnf *config.Config, roleID string) (domain.User, error) {
	hashedPassword, err := hash.CreateHash(cnf.Seeder.Admin.Password, &hash.DefaultArgon2Params)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to hash password for admin seeder: %w", err)
	}
	return domain.User{
		Name:     cnf.Seeder.Admin.Name,
		Email:    cnf.Seeder.Admin.Email,
		Password: hashedPassword,
		RoleID:   roleID,
	}, nil
}
