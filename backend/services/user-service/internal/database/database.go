package database

import (
	"context"
	"log"
	"time"
	"user-service/internal/config"

	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func ConnectDB(cfg *config.Config) {
	const (
		maxRetries    = 5
		retryInterval = 3 * time.Second
	)

	var conn *pgx.Conn
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		conn, err = pgx.Connect(context.Background(), cfg.DatabaseURL)
		if err == nil {
			log.Println("Connected to PostgreSQL")
			db = conn
			return
		}

		log.Printf("Attempt %d/%d: Failed to connect to PostgreSQL: %v. Retrying in %v...\n", attempt, maxRetries, err, retryInterval)
		time.Sleep(retryInterval)
	}

	log.Fatalf("Could not connect to PostgreSQL after %d attempts: %v", maxRetries, err)
}

func GetDB() *pgx.Conn {
	if db == nil {
		log.Fatal("Database connection is not initialized!")
	}
	return db
}

func CloseDB() {
	if db != nil {
		db.Close(context.Background())
		log.Println("Database connection closed.")
	}
}
