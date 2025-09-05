package util

import (
	"encoding/json"
	"log"
	"sort"
	"sync"

	"github.com/ipincamp/go-edsa-api/domain"
	"github.com/ipincamp/go-edsa-api/internal/constant"
	"gorm.io/gorm"
)

var (
	rolesByID   = make(map[string]domain.Role)
	rolesByName = make(map[string]domain.Role)
	once        sync.Once
)

func LoadRolesAndPermissions(db *gorm.DB) {
	once.Do(func() {
		var roles []domain.Role
		if err := db.Find(&roles).Error; err != nil {
			log.Fatalf("Failed to load roles for cache: %v", err)
		}

		rolesByID = make(map[string]domain.Role)
		rolesByName = make(map[string]domain.Role)

		for _, role := range roles {
			rolesByID[role.ID] = role
			rolesByName[role.Name] = role
		}
	})
}

func GetAllRoles() []domain.Role {
	allRoles := make([]domain.Role, 0, len(rolesByID))
	for _, role := range rolesByID {
		allRoles = append(allRoles, role)
	}

	return allRoles
}

func GetRoleByName(name string) (domain.Role, bool) {
	role, found := rolesByName[name]
	return role, found
}

func GetRoleByID(id string) (domain.Role, bool) {
	role, found := rolesByID[id]
	return role, found
}

func PrintPermissionTree(roles []domain.Role) {
	log.Printf(
		"├── %s%sCaching roles and permissions...%s",
		constant.Color("bold"),
		constant.Color("yellow"),
		constant.Color("reset"),
	)

	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Name < roles[j].Name
	})

	basePrefix := "│   "
	for i, role := range roles {
		isLastRole := i == len(roles)-1
		roleConnector := "├──"
		permParentPrefix := basePrefix + "│   "
		if isLastRole {
			roleConnector = "└──"
			permParentPrefix = basePrefix + "    "
		}

		log.Printf(
			"%s%s%s %s%s%s%s",
			basePrefix, roleConnector, constant.Color("reset"),
			constant.Color("bold"), constant.Color("yellow"), role.Name, constant.Color("reset"),
		)

		var permissions map[string]bool
		if err := json.Unmarshal(role.Permissions, &permissions); err != nil {
			continue
		}

		permKeys := make([]string, 0, len(permissions))
		for pKey := range permissions {
			permKeys = append(permKeys, pKey)
		}
		sort.Strings(permKeys)

		for j, pName := range permKeys {
			isLastPerm := j == len(permKeys)-1
			permConnector := "├──"
			if isLastPerm {
				permConnector = "└──"
			}
			if permissions[pName] {
				log.Printf(
					"%s%s%s %s%s%s%s",
					permParentPrefix, constant.Color("gray"), permConnector, constant.Color("reset"),
					constant.Color("green"), pName, constant.Color("reset"),
				)
			}
		}
	}

	log.Printf(
		"├── %s%sCached permissions for %d roles.%s",
		constant.Color("bold"), constant.Color("green"), len(roles), constant.Color("reset"),
	)
}
