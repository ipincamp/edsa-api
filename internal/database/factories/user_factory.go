package factories

import (
	"log"
	"path"

	"github.com/go-faker/faker/v4"
	"github.com/ipincamp/go-edsa-api/internal/config"
	"github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"github.com/ipincamp/go-edsa-api/internal/service/argon2id"
)

func UserFactory(roleID uint, avatarFileName string) *gorm.UserGORM {
	passSvc := argon2id.NewPasswordService()
	hashedPassword, err := passSvc.Hash("password")
	if err != nil {
		log.Fatalf("Failed to hash password for factory: %v", err)
	}

	baseURL := config.AppConfig.Storage.StoragePublicBaseURL
	defaultAvatarURL := baseURL + path.Join("/public/uploads", avatarFileName)

	return &gorm.UserGORM{
		Name:              faker.Name(),
		Email:             faker.Email(),
		Password:          hashedPassword,
		RoleID:            roleID,
		ProfilePictureURL: defaultAvatarURL,
	}
}

// --- AKHIR PERUBAHAN ---
