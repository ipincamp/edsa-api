package connection

import (
	"fmt"
	"log"

	"github.com/ipincamp/go-edsa-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GetDatabase(conf config.Database) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		conf.Host,
		conf.User,
		conf.Pass,
		conf.Name,
		conf.Port,
		conf.Tz,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database:", err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying sql.DB:", err.Error())
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err.Error())
	}

	return db
}
