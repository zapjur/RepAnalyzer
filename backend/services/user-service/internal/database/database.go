package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
	"user-service/internal/config"
)

var db *pgxpool.Pool

func ConnectDB(cfg *config.Config) {
	const (
		maxRetries    = 5
		retryInterval = 3 * time.Second
	)

	var pool *pgxpool.Pool
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		pool, err = pgxpool.New(context.Background(), cfg.DatabaseURL)
		if err == nil {
			log.Println("Connected to PostgreSQL")
			db = pool
			return
		}

		log.Printf("Attempt %d/%d: Failed to connect to PostgreSQL: %v. Retrying in %v...\n", attempt, maxRetries, err, retryInterval)
		time.Sleep(retryInterval)
	}

	log.Fatalf("Could not connect to PostgreSQL after %d attempts: %v", maxRetries, err)
}

func GetDB() *pgxpool.Pool {
	if db == nil {
		log.Fatal("Database connection is not initialized!")
	}
	return db
}

func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("Database connection closed.")
	}
}
