package main

import (
	"db-service/internal/database"
	"db-service/internal/server"

	"db-service/internal/config"
)

func main() {
	cfg := config.Load()

	database.ConnectDB(cfg)
	defer database.CloseDB()

	db := database.GetDB()
	server.StartGRPCServer(cfg, db)
}
