package factories

import (
	"log"
	"sync"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"

	"gorm.io/gorm"
)

// avatarIdCache stores fetched avatar IDs to prevent repeated DB queries.
var (
	avatarIdCache  = make(map[string]*uuid.UUID)
	cacheMutex     sync.RWMutex
	notFoundMarker = &uuid.UUID{}
)

// getOrFetchAvatarIDFromDB (tetap sama seperti sebelumnya)
func getOrFetchAvatarIDFromDB(db *gorm.DB, avatarFileName string) *uuid.UUID {
	cacheMutex.RLock()
	cachedID, found := avatarIdCache[avatarFileName]
	cacheMutex.RUnlock()

	if found {
		if cachedID == notFoundMarker {
			return nil
		}
		return cachedID
	}

	var avatarAsset repo.MediaAssetGORM
	var fetchedID *uuid.UUID
	var cacheValueToStore *uuid.UUID

	result := db.Where("file_name = ?", avatarFileName).First(&avatarAsset)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("WARNING: Default avatar '%s' not found in database. Caching this result.", avatarFileName)
			fetchedID = nil
			cacheValueToStore = notFoundMarker
		} else {
			log.Printf("ERROR: Failed to query avatar '%s' from database: %v. Returning nil without caching.", avatarFileName, result.Error)
			fetchedID = nil
			cacheValueToStore = nil
		}
	} else {
		fetchedID = &avatarAsset.ID
		cacheValueToStore = fetchedID
		log.Printf("INFO: Fetched avatar ID for '%s' from DB and cached it.", avatarFileName)
	}

	if cacheValueToStore != nil {
		cacheMutex.Lock()
		avatarIdCache[avatarFileName] = cacheValueToStore
		cacheMutex.Unlock()
	}

	return fetchedID
}

// UserFactory sekarang menerima hashedPassword sebagai parameter
func UserFactory(db *gorm.DB, roleID uint, defaultAvatarFileName string, hashedPassword string) *repo.UserGORM {
	// Hapus pembuatan instance passSvc dan hashing di sini

	// Dapatkan avatar ID menggunakan cache
	profilePicIDPtr := getOrFetchAvatarIDFromDB(db, defaultAvatarFileName)

	return &repo.UserGORM{
		Name:             faker.Name(),
		Email:            faker.Email(),
		Password:         hashedPassword, // Gunakan hash dari parameter
		RoleID:           roleID,
		ProfilePictureID: profilePicIDPtr,
	}
}
