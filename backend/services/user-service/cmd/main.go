package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"user-service/internal/config"
	"user-service/internal/server"
)

func main() {
	cfg := config.Load()

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server.Register(grpcServer)

	log.Printf("User Service running on port :%s", cfg.GRPCPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
