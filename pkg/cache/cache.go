package cache

import (
	"log"
	"sync"

	"github.com/ipincamp/go-edsa-api/internal/constant"
	"github.com/ipincamp/go-edsa-api/internal/domain"
	"gorm.io/gorm"
)

// rolesByID adalah cache role berdasarkan ID
var rolesByID = make(map[string]domain.Role)

// rolesByName adalah cache role berdasarkan nama
var rolesByName = make(map[string]domain.Role)

// usersByID adalah cache user berdasarkan ID
var usersByID = make(map[string]domain.User)

// cacheMutex digunakan untuk concurrency cache
var cacheMutex = &sync.RWMutex{}

// once memastikan cache hanya di-load sekali
var once sync.Once

func LoadCache(db *gorm.DB) {
	log.Printf(
		"├── %s%sCaching roles and permissions...%s",
		constant.Color("bold"),
		constant.Color("yellow"),
		constant.Color("reset"),
	)

	once.Do(func() {
		cacheMutex.Lock()
		defer cacheMutex.Unlock()

		var users []domain.User
		if err := db.Preload("Role").Find(&users).Error; err != nil {
			log.Fatalf("Failed to load users and roles for cache: %v", err)
		}

		roleSet := make(map[string]domain.Role)
		for _, user := range users {
			usersByID[user.ID] = user
			if user.Role.ID != "" {
				roleSet[user.Role.ID] = user.Role
			}
		}

		// Populate role caches
		for _, role := range roleSet {
			rolesByID[role.ID] = role
			rolesByName[role.Name] = role
		}

		log.Printf(
			"│   ├── %s%sCached %d roles.%s",
			constant.Color("bold"), constant.Color("green"), len(roleSet), constant.Color("reset"),
		)
		log.Printf(
			"│   └── %s%sCached %d users.%s",
			constant.Color("bold"), constant.Color("green"), len(users), constant.Color("reset"),
		)
	})
}

func AddUserToCache(user domain.User) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	usersByID[user.ID] = user
	log.Printf("User %s added/updated in cache", user.ID)
}

func GetUserFromCacheByID(id string) (domain.User, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	user, found := usersByID[id]
	return user, found
}

func GetAllRoles() []domain.Role {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	allRoles := make([]domain.Role, 0, len(rolesByID))
	for _, role := range rolesByID {
		allRoles = append(allRoles, role)
	}

	return allRoles
}

func GetRoleByName(name string) (domain.Role, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	role, found := rolesByName[name]
	return role, found
}

func GetRoleByID(id string) (domain.Role, bool) {
	cacheMutex.RLock()
	defer cacheMutex.RUnlock()
	role, found := rolesByID[id]
	return role, found
}
