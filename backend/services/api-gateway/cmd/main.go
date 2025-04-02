package main

import (
	"api-gateway/internal/minio"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"api-gateway/internal/auth"
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
	auth.SetAuth0Domain(cfg.Auth0Domain)

	grpcClient, err := grpc.NewClient(cfg.UserServiceAddress)
	if err != nil {
		log.Fatalf("failed to setup gRPC client: %v", err)
	}
	defer grpcClient.Close()

	minioClient := minio.NewMinioClient()
	minioClient.EnsureBucketExists("videos")

	userHandler := handlers.NewUserHandler(grpcClient)
	videoHandler := handlers.NewVideoHandler(minioClient.Minio, grpcClient)

	r := chi.NewRouter()

	r.Use(corsMiddleware)
	r.Use(auth.JwtMiddleware)

	r.Get("/users/{auth0ID}", userHandler.GetUser)
	r.Post("/upload/{auth0ID}", videoHandler.UploadVideo)

	log.Printf("API Gateway started on port %s", cfg.HTTPPort)
	if err = http.ListenAndServe(":"+cfg.HTTPPort, r); err != nil {
		log.Fatal(err)
	}
}
