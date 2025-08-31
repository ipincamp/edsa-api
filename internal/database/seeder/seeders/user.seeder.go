package seeders

import (
	"fmt"
	"log"

	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/util"
	"gorm.io/gorm"
)

func UserSeeder(cnf *config.Config) func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		var adminRole domain.Role
		if err := db.Where("name = ?", constant.RoleAdmin).First(&adminRole).Error; err != nil {
			return fmt.Errorf("admin role not found for user seeder: %w", err)
		}

		hashedPassword, err := util.HashPassword(cnf.Seeder.Admin.Password)
		if err != nil {
			return fmt.Errorf("failed to hash password for admin seeder: %w", err)
		}

		adminUser := domain.User{
			Name:     cnf.Seeder.Admin.Name,
			Email:    cnf.Seeder.Admin.Email,
			Password: hashedPassword,
			RoleID:   adminRole.ID,
		}

		err = db.FirstOrCreate(&adminUser, "email = ?", adminUser.Email).Error
		if err != nil {
			return fmt.Errorf("failed to seed admin user: %w", err)
		}
		log.Println("User seeder ran successfully")

		return nil
	}
}
