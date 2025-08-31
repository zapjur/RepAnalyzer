package main

import (
	"analysis-service/internal/config"
	"analysis-service/internal/rabbitmq"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	cfg := config.Load()

	rabbitConn, rabbitChannel, err := rabbitmq.ConnectToRabbitMQ(cfg.RabbitMQURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitConn.Close()
	defer rabbitChannel.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := rabbitmq.StartConsumers(ctx, rabbitChannel); err != nil {
		log.Fatalf("Failed to start consumers: %v", err)
	}

	<-ctx.Done()
	log.Println("Shutting down gracefully...")
}
