package seeders

import (
	"log"

	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"gorm.io/gorm"
)

func BookSeeder(db *gorm.DB) {
	log.Println("Seeding books...")

	// 8 Buku berdasarkan 8 tema di BRD
	books := []repo.BookGORM{
		{
			Title:         "The Alphabet in the Land of Dewi Sri",
			Description:   "An exciting journey through the ABCs.",
			CoverImageURL: "https://example.com/images/covers/alphabet.png",
			Theme:         "Alphabet",
			BookOrder:     1,
		},
		{
			Title:         "The Numbers in the Village of Ten Hills",
			Description:   "Learn numbers with friendly creatures.",
			CoverImageURL: "https://example.com/images/covers/numbers.png",
			Theme:         "Numbers",
			BookOrder:     2,
		},
		{
			Title:         "Bawang Putih and the Kind Body Parts",
			Description:   "Discover all the parts of your body.",
			CoverImageURL: "https://example.com/images/covers/body.png",
			Theme:         "Body Parts",
			BookOrder:     3,
		},
		{
			Title:         "The Big Mango Tree and the Family of Five",
			Description:   "A wonderful day out with the whole family.",
			CoverImageURL: "https://example.com/images/covers/family.png",
			Theme:         "Family",
			BookOrder:     4,
		},
		{
			Title:         "The Magic Paintbrush",
			Description:   "Mixing colors to create new ones.",
			CoverImageURL: "https://example.com/images/covers/colors.png",
			Theme:         "Colors",
			BookOrder:     5,
		},
		{
			Title:         "Sounds of the Wild",
			Description:   "Listen to the sounds animals make.",
			CoverImageURL: "https://example.com/images/covers/animals.png",
			Theme:         "Animals",
			BookOrder:     6,
		},
		{
			Title:         "What Time Is It?",
			Description:   "Learn to tell time with Timmy the clock.",
			CoverImageURL: "https://example.com/images/covers/time.png",
			Theme:         "Time",
			BookOrder:     7,
		},
		{
			Title:         "Action Day!",
			Description:   "Jumping, running, and playing all day.",
			CoverImageURL: "https://example.com/images/covers/verbs.png",
			Theme:         "Verbs",
			BookOrder:     8,
		},
	}

	for _, book := range books {
		// Buat atau cari berdasarkan Judul
		result := db.FirstOrCreate(&book, repo.BookGORM{Title: book.Title})
		if result.Error != nil {
			log.Printf("Failed to seed book '%s': %v\n", book.Title, result.Error)
		}
		if result.RowsAffected > 0 {
			log.Printf("Seeded book: %s\n", book.Title)
		}
	}
}
