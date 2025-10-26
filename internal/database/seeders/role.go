package seeders

import (
	"fmt"
	"log"

	"github.com/ipincamp/go-edsa-api/internal/domain"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"gorm.io/gorm"
)

func RoleSeeder(db *gorm.DB, logger *log.Logger) error {
	logger.Println("Seeding roles...")
	roles := []repo.RoleGORM{
		{Name: domain.RoleNameAdmin},
		{Name: domain.RoleNameTeacher},
		{Name: domain.RoleNameStudent},
		{Name: domain.RoleNamePublic},
	}

	for _, role := range roles {
		result := db.FirstOrCreate(&role, repo.RoleGORM{Name: role.Name})
		if result.Error != nil {
			return fmt.Errorf("failed to seed role '%s': %w", role.Name, result.Error)
		}
		if result.RowsAffected > 0 {
			logger.Printf("Seeded role: %s\n", role.Name)
		}
	}
	return nil
}
