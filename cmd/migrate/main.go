package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/lib/pq"
	"github.com/masatrio/bookstore-api/config"
)

func main() {
	cfg := config.LoadConfig()
	dbURL := cfg.Database.URL

	_, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	m, err := migrate.New(
		"file://migrations",
		dbURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to migrate: %v", err)
	}

	log.Println("Migrations applied successfully")
}
