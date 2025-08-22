package main

import (
	"log"

	"github.com/ipincamp/edsa/internal/config"
	"github.com/ipincamp/edsa/internal/database"
)

func main() {
	log.Println("Starting seeder application...")

	env := config.LoadEnv()
	database.ConnectDB(env)

	database.RunSeeder(database.DB, env)
}
