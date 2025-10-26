package factories

import (
	"log"

	"github.com/go-faker/faker/v4"
	"github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"github.com/ipincamp/go-edsa-api/internal/service/argon2id"
)

func UserFactory(roleID uint) *gorm.UserGORM {
	passSvc := argon2id.NewPasswordService()
	hashedPassword, err := passSvc.Hash("password")
	if err != nil {
		log.Fatalf("Failed to hash password for factory: %v", err)
	}

	return &gorm.UserGORM{
		Name:             faker.Name(),
		Email:            faker.Email(),
		Password:         hashedPassword,
		RoleID:           roleID,
		ProfilePictureID: nil,
	}
}

// --- AKHIR PERUBAHAN ---
