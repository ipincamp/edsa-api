package seeders

import (
	"fmt"
	"log"

	repo "github.com/ipincamp/go-edsa-api/internal/repository/gorm"
	"gorm.io/gorm"
)

func GameSeeder(db *gorm.DB, logger *log.Logger) error {
	logger.Println("Seeding games...")

	games := []repo.GameGORM{
		{Name: "Alphabet Match", Type: "Matching", RelatedBookTheme: "Alphabet"},
		{Name: "Number Pop", Type: "Tapping", RelatedBookTheme: "Numbers"},
		{Name: "Body Parts Label", Type: "DragDrop", RelatedBookTheme: "Body Parts"},
		{Name: "Family Tree", Type: "DragDrop", RelatedBookTheme: "Family"},
		{Name: "Color Mixing", Type: "Quiz", RelatedBookTheme: "Colors"},
		{Name: "Animal Sounds", Type: "Listening", RelatedBookTheme: "Animals"},
		{Name: "Time Quiz", Type: "Quiz", RelatedBookTheme: "Time"},
		{Name: "Verb Action", Type: "Matching", RelatedBookTheme: "Verbs"},
	}

	for _, game := range games {
		result := db.FirstOrCreate(&game, repo.GameGORM{Name: game.Name})
		if result.Error != nil {
			return fmt.Errorf("failed to seed game '%s': %w", game.Name, result.Error)
		}
		if result.RowsAffected > 0 {
			logger.Printf("Seeded game: %s\n", game.Name)
		}
	}
	return nil
}
