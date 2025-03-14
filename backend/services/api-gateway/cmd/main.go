package main

import (
	"log"
	"net/http"

	"api-gateway/internal/config"
	"api-gateway/internal/grpc"
	"api-gateway/internal/handlers"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	cfg := config.Load()

	grpcClient, err := grpc.NewClient(cfg.UserServiceAddress)
	if err != nil {
		log.Fatalf("failed to setup gRPC client: %v", err)
	}
	defer grpcClient.Close()

	userHandler := handlers.NewUserHandler(grpcClient)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/", userHandler.GetUser)

	handlerWithCORS := corsMiddleware(mux)

	log.Printf("API Gateway started on port %s", cfg.HTTPPort)
	if err := http.ListenAndServe(":"+cfg.HTTPPort, handlerWithCORS); err != nil {
		log.Fatal(err)
	}
}
