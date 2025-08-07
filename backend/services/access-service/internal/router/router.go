package router

import (
	"access-service/internal/auth"
	"access-service/internal/grpc"
	"access-service/internal/handlers"
	"access-service/internal/minio"
	"access-service/internal/service"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"os"
)

func Setup() http.Handler {
	r := chi.NewRouter()

	r.Use(auth.JwtMiddleware)
	auth.SetAuth0Domain(os.Getenv("AUTH0_DOMAIN"))

	minioClient := minio.NewClient()
	grpcClient, err := grpc.NewClient("db-service:50051")
	if err != nil {
		log.Fatal("failed to setup gRPC client:", err)
	}
	defer grpcClient.Close()

	svc := service.NewAccessService(minioClient, grpcClient)
	h := handler.NewAccessHandler(svc)

	r.Get("/access/video/{videoId}", h.GetPresignedURL)

	return r
}
