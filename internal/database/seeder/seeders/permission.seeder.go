package seeders

import (
	"fmt"
	"log"

	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"gorm.io/gorm"
)

func PermissionSeeder(db *gorm.DB) error {
	permissions := []domain.Permission{
		{Name: constant.PermissionRead.String()},
		{Name: constant.PermissionWrite.String()},
		{Name: constant.PermissionDelete.String()},
	}

	for _, permission := range permissions {
		err := db.FirstOrCreate(&permission, "name = ?", permission.Name).Error
		if err != nil {
			return fmt.Errorf("failed to seed permission %s: %w", permission.Name, err)
		}
	}
	log.Println("Permission seeder ran successfully")

	return nil
}
