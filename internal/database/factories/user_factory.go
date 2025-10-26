package factories

import (
	"log"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"

	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"github.com/ipincamp/go-edsa-api/internal/service/argon2id"
	"gorm.io/gorm"
)

// Fungsi getAvatarIDFromDB adalah helper untuk mencari ID asset avatar di DB
func getAvatarIDFromDB(db *gorm.DB, avatarFileName string) *uuid.UUID {
	var avatarAsset repo.MediaAssetGORM
	// Cari asset berdasarkan nama file asli
	result := db.Where("file_name = ?", avatarFileName).First(&avatarAsset)
	if result.Error != nil {
		// Jika tidak ditemukan atau ada error lain, log peringatan dan kembalikan nil
		log.Printf("WARNING: Default avatar '%s' not found in database. User will have no avatar. Error: %v", avatarFileName, result.Error)
		return nil
	}
	// Jika ditemukan, kembalikan pointer ke ID-nya
	return &avatarAsset.ID
}

// Ubah UserFactory: Hapus parameter avatarFileName, cari ID di dalam fungsi
func UserFactory(db *gorm.DB, roleID uint, defaultAvatarFileName string) *repo.UserGORM {
	passSvc := argon2id.NewPasswordService()
	hashedPassword, err := passSvc.Hash("password")
	if err != nil {
		log.Fatalf("Failed to hash password for factory: %v", err)
	}

	// Cari ID avatar dari DB menggunakan helper
	profilePicIDPtr := getAvatarIDFromDB(db, defaultAvatarFileName)

	return &repo.UserGORM{
		Name:             faker.Name(),
		Email:            faker.Email(),
		Password:         hashedPassword,
		RoleID:           roleID,
		ProfilePictureID: profilePicIDPtr, // Gunakan ID yang didapat dari DB
	}
}
