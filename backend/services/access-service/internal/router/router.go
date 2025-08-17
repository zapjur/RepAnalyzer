package router

import (
	"access-service/internal/auth"
	"access-service/internal/grpc"
	"access-service/internal/handlers"
	"access-service/internal/minio"
	"access-service/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
)

func Setup(minioClient *minio.Client, grpcClient *grpc.Client) http.Handler {
	r := chi.NewRouter()

	auth.SetAuth0Domain(os.Getenv("AUTH0_DOMAIN"))
	r.Use(auth.JwtMiddleware)

	svc := service.NewAccessService(minioClient, grpcClient)
	h := handler.NewAccessHandler(svc)

	r.Get("/access/video/{videoId}", h.GetPresignedURL)
	r.Get("/access/video-analysis/{videoId}", h.GetVideoAnalysis)
	return r
}
