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

func UserSeeder(cnf *config.Config) func(db *gorm.DB) error {
	return func(db *gorm.DB) error {
		var adminRole domain.Role
		if err := db.Where("name = ?", constant.RoleAdmin).First(&adminRole).Error; err != nil {
			return fmt.Errorf("admin role not found for user seeder: %w", err)
		}

		hashParams := &hash.Argon2Params{
			Memory:      64 * 1024,
			Iterations:  3,
			Parallelism: 2,
			SaltLength:  16,
			KeyLength:   32,
		}

		hashedPassword, err := hash.CreateHash(cnf.Seeder.Admin.Password, hashParams)
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
