package main

import (
	"log"
	"net/http"

	"api-gateway/internal/config"
	"api-gateway/internal/grpc"
	"api-gateway/internal/handlers"
)

func main() {
	cfg := config.Load()

	grpcClient, err := grpc.NewClient(cfg.UserServiceAddress)
	if err != nil {
		log.Fatalf("failed to setup gRPC client: %v", err)
	}
	defer grpcClient.Close()

	userHandler := handlers.NewUserHandler(grpcClient)

	http.HandleFunc("/users/", userHandler.GetUser)

	log.Printf("API Gateway started on port %s", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, nil); err != nil {
		log.Fatal(err)
	}
}
