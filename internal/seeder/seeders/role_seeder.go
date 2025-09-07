package seeders

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

// Permissions adalah tipe untuk mapping permission pada role
type Permissions map[string]bool

// roleSeedData berisi data role dan permission yang akan di-seed
var roleSeedData = map[string]Permissions{
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

// RoleSeeder melakukan seed data role ke database
func RoleSeeder(db *gorm.DB) error {
	for roleName, permissions := range roleSeedData {
		if err := seedRole(db, roleName, permissions); err != nil {
			return err
		}
	}
	log.Println("Role seeder ran successfully")
	return nil
}

// seedRole adalah helper untuk membuat satu role ke database
func seedRole(db *gorm.DB, roleName string, permissions Permissions) error {
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
	return nil
}
