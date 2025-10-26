package factories

import (
	"log"
	"sync"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"github.com/ipincamp/go-edsa-api/internal/service/argon2id"
	"gorm.io/gorm"
)

// avatarIdCache stores fetched avatar IDs to prevent repeated DB queries.
var (
	avatarIdCache = make(map[string]*uuid.UUID)
	cacheMutex    sync.RWMutex // Mutex to protect concurrent access to the cache
	// Use a special value to indicate "not found" has been cached
	notFoundMarker = &uuid.UUID{} // A non-nil pointer distinct from actual UUIDs
)

// getOrFetchAvatarIDFromDB retrieves avatar ID from cache or DB.
func getOrFetchAvatarIDFromDB(db *gorm.DB, avatarFileName string) *uuid.UUID {
	cacheMutex.RLock()
	cachedID, found := avatarIdCache[avatarFileName]
	cacheMutex.RUnlock()

	if found {
		if cachedID == notFoundMarker {
			return nil // We cached the fact that it wasn't found
		}
		return cachedID // Return from cache
	}

	// Not in cache, query the database
	var avatarAsset repo.MediaAssetGORM
	var fetchedID *uuid.UUID
	var cacheValueToStore *uuid.UUID // What to store in the cache

	result := db.Where("file_name = ?", avatarFileName).First(&avatarAsset)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("WARNING: Default avatar '%s' not found in database. Caching this result.", avatarFileName)
			fetchedID = nil
			cacheValueToStore = notFoundMarker // Mark as not found in cache
		} else {
			// Log other DB errors but don't cache the error state, maybe it's temporary
			log.Printf("ERROR: Failed to query avatar '%s' from database: %v. Returning nil without caching.", avatarFileName, result.Error)
			fetchedID = nil
			cacheValueToStore = nil // Don't cache DB errors
		}
	} else {
		// Found in DB
		fetchedID = &avatarAsset.ID
		cacheValueToStore = fetchedID // Cache the actual ID pointer
		log.Printf("INFO: Fetched avatar ID for '%s' from DB and cached it.", avatarFileName)
	}

	// Store the result (or notFoundMarker) in cache if we determined a cache value
	if cacheValueToStore != nil {
		cacheMutex.Lock()
		avatarIdCache[avatarFileName] = cacheValueToStore
		cacheMutex.Unlock()
	}

	return fetchedID
}

// UserFactory now uses the caching mechanism.
func UserFactory(db *gorm.DB, roleID uint, defaultAvatarFileName string) *repo.UserGORM {
	passSvc := argon2id.NewPasswordService()
	hashedPassword, err := passSvc.Hash("password")
	if err != nil {
		log.Fatalf("Failed to hash password for factory: %v", err)
	}

	// Get avatar ID using the caching helper function
	profilePicIDPtr := getOrFetchAvatarIDFromDB(db, defaultAvatarFileName)

	return &repo.UserGORM{
		Name:             faker.Name(),
		Email:            faker.Email(),
		Password:         hashedPassword,
		RoleID:           roleID,
		ProfilePictureID: profilePicIDPtr,
	}
}
