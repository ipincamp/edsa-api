package seeders

import (
	"fmt"
	"log"

	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"gorm.io/gorm"
)

func RolePermissionSeeder(db *gorm.DB) error {
	var roles []domain.Role
	if err := db.Find(&roles).Error; err != nil {
		return fmt.Errorf("failed to get roles: %w", err)
	}
	rolesMap := make(map[string]domain.Role)
	for _, r := range roles {
		rolesMap[r.Name] = r
	}

	var permissions []domain.Permission
	if err := db.Find(&permissions).Error; err != nil {
		return fmt.Errorf("failed to get permissions: %w", err)
	}
	permissionsMap := make(map[string]domain.Permission)
	for _, p := range permissions {
		permissionsMap[p.Name] = p
	}

	rolePermissions := map[string][]string{
		constant.RoleAdmin.String(): {
			constant.PermissionRead.String(),
			constant.PermissionWrite.String(),
			constant.PermissionDelete.String(),
		},
		constant.RoleTeacher.String(): {
			constant.PermissionRead.String(),
			constant.PermissionWrite.String(),
		},
		constant.RoleStudent.String(): {
			constant.PermissionRead.String(),
		},
	}

	for roleName, permissionNames := range rolePermissions {
		role, ok := rolesMap[roleName]
		if !ok {
			log.Printf("Role '%s' not found, skipping.", roleName)
			continue
		}

		var permsToAssign []domain.Permission
		for _, pName := range permissionNames {
			perm, ok := permissionsMap[pName]
			if ok {
				permsToAssign = append(permsToAssign, perm)
			}
		}

		if err := db.Model(&role).Association("Permissions").Replace(permsToAssign); err != nil {
			return fmt.Errorf("failed to assign permissions to role %s: %w", roleName, err)
		}
	}

	log.Println("Role-Permission seeder ran successfully")
	return nil
}
