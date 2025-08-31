package main

import (
	"log"
	"orchestrator/internal/client"
	"orchestrator/internal/config"
	"orchestrator/internal/rabbitmq"
	"orchestrator/internal/redis"
	"orchestrator/internal/server"
)

func main() {

	cfg := config.Load()

	rabbitClient, err := rabbitmq.ConnectToRabbitMQ(cfg.RabbitMQURI)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitClient.Channel.Close()
	defer rabbitClient.Connection.Close()

	queues := []string{"bar_path_queue", "bar_path_results_queue", "pose_queue", "pose_results_queue", "analysis_queue", "analysis_results_queue"}
	err = rabbitClient.DeclareQueues(queues)
	if err != nil {
		log.Fatalf("Failed to declare RabbitMQ queues: %v", err)
	}

	redisClient, err := redis.ConnectToRedis(cfg.RedisAddr)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	redisManager := &redis.RedisManager{RedisClient: redisClient}

	grpcClient, err := client.NewClient(cfg.DBServiceAddress)
	if err != nil {
		log.Fatalf("failed to setup gRPC DB client: %v", err)
	}
	defer grpcClient.Close()

	rabbitClient.StartConsumers(redisManager, grpcClient)

	server.StartGRPCServer(cfg, redisManager, rabbitClient)
}
