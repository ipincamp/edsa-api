package seeders

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

// Permissions adalah tipe untuk mapping permission pada role.
type Permissions map[constant.Permission]bool

var roleSeedData = map[string]Permissions{
	constant.RoleAdmin.String(): {
		constant.UsersCreate:      true,
		constant.UsersListAll:     true,
		constant.UsersViewOther:   true,
		constant.UsersUpdateOther: true,
		constant.UsersDelete:      true,
		constant.RolesManage:      true,
	},
	constant.RoleTeacher.String(): {
		constant.CoursesCreate: true,
		constant.CoursesUpdate: true,
		constant.GradesManage:  true,
	},
	constant.RoleStudent.String(): {
		constant.CoursesRead:      true,
		constant.AssignmentSubmit: true,
		constant.UsersViewSelf:    true,
		constant.UsersUpdateSelf:  true,
	},
	constant.RoleGuest.String(): {
		// Guest tidak memiliki permission
	},
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
	stringPermissions := make(map[string]bool)
	for p, v := range permissions {
		stringPermissions[p.String()] = v
	}

	permissionsJSON, err := json.Marshal(stringPermissions)
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
