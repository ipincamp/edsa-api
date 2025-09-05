package seeders

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"gorm.io/gorm"
)

type Permissions map[string]bool

func RoleSeeder(db *gorm.DB) error {
	rolePermissions := map[string]Permissions{
		constant.RoleAdmin.String(): {
			"users.create": true,
			"users.read":   true,
			"users.update": true,
			"users.delete": true,
			"roles.manage": true,
		},
		constant.RoleTeacher.String(): {
			"courses.create": true,
			"courses.update": true,
			"grades.manage":  true,
		},
		constant.RoleStudent.String(): {
			"courses.read":      true,
			"assignment.submit": true,
		},
		constant.RoleGuest.String(): {},
	}

	for roleName, permissions := range rolePermissions {
		permissionsJSON, err := json.Marshal(permissions)
		if err != nil {
			return fmt.Errorf("failed to marshal permissions for role %s: %w", roleName, err)
		}

		role := domain.Role{
			Name:        roleName,
			Permissions: permissionsJSON,
		}

		if err := db.FirstOrCreate(&role, "name = ?", role.Name).Error; err != nil {
			return fmt.Errorf("failed to seed role %s: %w", role.Name, err)
		}
	}
	log.Println("Role seeder ran successfully")

	return nil
}
