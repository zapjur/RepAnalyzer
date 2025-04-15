package main

import (
	"orchestrator/internal/config"
	"orchestrator/internal/server"
)

func main() {

	cfg := config.Load()

	server.StartGRPCServer(cfg)
}
