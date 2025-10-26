package shared

import "github.com/google/uuid"

// DefaultAvatarAssets menyimpan ID UUID yang sudah dibuat untuk avatar default
// agar bisa diakses oleh seeder dan factory.
var DefaultAvatarAssets map[string]uuid.UUID

func init() {
	DefaultAvatarAssets = make(map[string]uuid.UUID)
}
