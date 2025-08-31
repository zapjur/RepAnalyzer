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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	rabbitClient, err := rabbitmq.ConnectToRabbitMQ(cfg.RabbitMQURI, ctx)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitClient.Channel.Close()
	defer rabbitClient.Connection.Close()

	if err := rabbitClient.StartConsumers(); err != nil {
		log.Fatalf("Failed to start consumers: %v", err)
	}

	<-ctx.Done()
	log.Println("Shutting down gracefully...")
}
