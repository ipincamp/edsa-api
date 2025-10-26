package seeders

import (
	"fmt"
	"log"

	"github.com/ipincamp/go-edsa-api/internal/config"
	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"gorm.io/gorm"
)

func BookSeeder(db *gorm.DB) error {
	log.Println("Seeding books...")

	// Ambil base URL dari config untuk contoh URL
	baseURL := config.AppConfig.Storage.StoragePublicBaseURL // e.g., http://localhost:8000

	// 8 Buku berdasarkan 8 tema di BRD
	books := []repo.BookGORM{
		{
			Title:         "The Alphabet in the Land of Dewi Sri",
			Description:   "An exciting journey through the ABCs.",
			CoverImageURL: fmt.Sprintf("%s/cdn/covers/alphabet.png", baseURL), // <-- Ubah URL
			Theme:         "Alphabet",
			BookOrder:     1,
		},
		{
			Title:         "The Numbers in the Village of Ten Hills",
			Description:   "Learn numbers with friendly creatures.",
			CoverImageURL: fmt.Sprintf("%s/cdn/covers/numbers.png", baseURL), // <-- Ubah URL
			Theme:         "Numbers",
			BookOrder:     2,
		},
		{
			Title:         "Bawang Putih and the Kind Body Parts",
			Description:   "Discover all the parts of your body.",
			CoverImageURL: fmt.Sprintf("%s/cdn/covers/body.png", baseURL), // <-- Ubah URL
			Theme:         "Body Parts",
			BookOrder:     3,
		},
		{
			Title:         "The Big Mango Tree and the Family of Five",
			Description:   "A wonderful day out with the whole family.",
			CoverImageURL: fmt.Sprintf("%s/cdn/covers/family.png", baseURL), // <-- Ubah URL
			Theme:         "Family",
			BookOrder:     4,
		},
		{
			Title:         "The Magic Paintbrush",
			Description:   "Mixing colors to create new ones.",
			CoverImageURL: fmt.Sprintf("%s/cdn/covers/colors.png", baseURL), // <-- Ubah URL
			Theme:         "Colors",
			BookOrder:     5,
		},
		{
			Title:         "Sounds of the Wild",
			Description:   "Listen to the sounds animals make.",
			CoverImageURL: fmt.Sprintf("%s/cdn/covers/animals.png", baseURL), // <-- Ubah URL
			Theme:         "Animals",
			BookOrder:     6,
		},
		{
			Title:         "What Time Is It?",
			Description:   "Learn to tell time with Timmy the clock.",
			CoverImageURL: fmt.Sprintf("%s/cdn/covers/time.png", baseURL), // <-- Ubah URL
			Theme:         "Time",
			BookOrder:     7,
		},
		{
			Title:         "Action Day!",
			Description:   "Jumping, running, and playing all day.",
			CoverImageURL: fmt.Sprintf("%s/cdn/covers/verbs.png", baseURL), // <-- Ubah URL
			Theme:         "Verbs",
			BookOrder:     8,
		},
	}

	for _, book := range books {
		// Buat atau cari berdasarkan Judul
		result := db.FirstOrCreate(&book, repo.BookGORM{Title: book.Title})
		if result.Error != nil {
			return fmt.Errorf("failed to seed book '%s': %w", book.Title, result.Error)
		}
		if result.RowsAffected > 0 {
			log.Printf("Seeded book: %s\n", book.Title)
		}
	}
	return nil
}
