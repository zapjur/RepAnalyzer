package database

import (
	"context"
	"log"
	"user-service/internal/config"

	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn

func ConnectDB(cfg *config.Config) {
	conn, err := pgx.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL")
	db = conn
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
