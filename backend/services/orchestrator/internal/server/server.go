package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"log"
	"net"
	"orchestrator/internal/config"
	"orchestrator/internal/redis"
	anPb "orchestrator/proto/analysis"
	"strings"
)

type OrchestratorServer struct {
	anPb.UnimplementedOrchestratorServer
	Redis         *redis.RedisManager
	RabbitChannel *amqp.Channel
}

type VideoToAnalyze struct {
	Bucket       string
	ObjectKey    string
	ExerciseName string
	Auth0Id      string
	VideoId      int64
}

type TaskMessage struct {
	Bucket       string `json:"bucket"`
	ObjectKey    string `json:"object_key"`
	ExerciseName string `json:"exercise_name"`
	VideoID      string `json:"video_id"`
	Auth0Id      string `json:"auth0_id"`
	ReplyQueue   string `json:"reply_queue"`
}

func StartGRPCServer(cfg *config.Config, redisManager *redis.RedisManager, rabbitChannel *amqp.Channel) {
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	anPb.RegisterOrchestratorServer(grpcServer, &OrchestratorServer{
		Redis:         redisManager,
		RabbitChannel: rabbitChannel,
	})

	log.Printf("Orchestrator Service running on port :%s", cfg.GRPCPort)

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *OrchestratorServer) AnalyzeVideo(ctx context.Context, req *anPb.VideoToAnalyzeRequest) (*anPb.VideoToAnalyzeResponse, error) {

	// adding task to redis
	videoIDStr := fmt.Sprintf("%d", req.VideoId)
	err := s.Redis.SetTaskStatus(videoIDStr, "bar_path", "pending")
	if err != nil {
		log.Printf("Failed to set task status in Redis: %v", err)
		return errorResponse("Failed to set task status in Redis", err)

	}

	auth0IDEdited := strings.ReplaceAll(req.Auth0Id, "|", "_")

	// sending video data to appropriate queue
	msg := TaskMessage{
		Bucket:       req.Bucket,
		ObjectKey:    req.ObjectKey,
		ExerciseName: req.ExerciseName,
		VideoID:      videoIDStr,
		Auth0Id:      auth0IDEdited,
		ReplyQueue:   "bar_path_results_queue",
	}
	taskBody, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal task message: %v", err)
		return errorResponse("Failed to marshal task message", err)
	}
	err = s.RabbitChannel.Publish(
		"",
		"bar_path_queue",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        taskBody,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message to RabbitMQ: %v", err)
		return errorResponse("Failed to publish message to RabbitMQ", err)

	}

	// send back response
	response := &anPb.VideoToAnalyzeResponse{
		Success: true,
		Message: "Video analysis started successfully",
	}

	return response, nil
}
