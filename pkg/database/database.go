package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ipincamp/go-edsa-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect membuat koneksi ke database PostgreSQL
func Connect(cfg config.Database) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		cfg.Host, cfg.User, cfg.Pass, cfg.Name, cfg.Port, cfg.Tz,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // Output ke stdout
			logger.Config{
				SlowThreshold: time.Second, // Query lambat > 1 detik
				LogLevel:      logger.Info, // Level log: Info (bisa Debug untuk lebih detail)
				Colorful:      true,        // Output berwarna
			},
		),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err = sqlDB.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
