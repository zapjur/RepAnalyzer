package main

import (
	"user-service/internal/database"
	"user-service/internal/server"

	"user-service/internal/config"
)

func main() {
	cfg := config.Load()

	database.ConnectDB(cfg)
	defer database.CloseDB()

	server.StartGRPCServer(cfg)
}
