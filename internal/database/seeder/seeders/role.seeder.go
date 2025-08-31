package seeders

import (
	"fmt"
	"log"

	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"gorm.io/gorm"
)

func RoleSeeder(db *gorm.DB) error {
	roles := []domain.Role{
		{Name: constant.RoleAdmin.String()},
		{Name: constant.RoleTeacher.String()},
		{Name: constant.RoleStudent.String()},
	}

	for _, role := range roles {
		err := db.FirstOrCreate(&role, "name = ?", role.Name).Error
		if err != nil {
			return fmt.Errorf("failed to seed role %s: %w", role.Name, err)
		}
	}
	log.Println("Role seeder ran successfully")

	return nil
}
