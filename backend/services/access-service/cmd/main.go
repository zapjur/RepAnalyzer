package main

import (
	"access-service/internal/auth"
	"access-service/internal/grpc"
	"access-service/internal/minio"
	"access-service/internal/router"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	auth.SetAuth0Domain(os.Getenv("AUTH0_DOMAIN"))

	minioClient := minio.NewClient()

	dbAddr := "db-service:50051"
	grpcClient, err := grpc.NewClient(dbAddr)
	if err != nil {
		log.Fatal("failed to setup gRPC client:", err)
	}

	r := router.Setup(minioClient, grpcClient)

	srv := &http.Server{Addr: ":8082", Handler: r}
	go func() {
		log.Println("Access Service running on :8082")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
	_ = grpcClient.Close()
	log.Println("Access Service stopped")
}
