package factories

import (
	"log"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	"github.com/ipincamp/go-edsa-api/internal/database/shared"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"github.com/ipincamp/go-edsa-api/internal/service/argon2id"
	"gorm.io/gorm"
)

func UserFactory(db *gorm.DB, roleID uint, avatarFileName string) *repo.UserGORM {
	passSvc := argon2id.NewPasswordService()
	hashedPassword, err := passSvc.Hash("password")
	if err != nil {
		log.Fatalf("Failed to hash password for factory: %v", err)
	}

	// Ambil ID avatar dari map global di paket shared
	var profilePicIDPtr *uuid.UUID
	if avatarID, ok := shared.DefaultAvatarAssets[avatarFileName]; ok {
		idCopy := avatarID
		profilePicIDPtr = &idCopy
	} else {
		log.Printf("WARNING: Default avatar '%s' not found in seeded assets. User will have no avatar.", avatarFileName)
		profilePicIDPtr = nil
	}

	return &repo.UserGORM{
		Name:             faker.Name(),
		Email:            faker.Email(),
		Password:         hashedPassword,
		RoleID:           roleID,
		ProfilePictureID: profilePicIDPtr,
	}
}
