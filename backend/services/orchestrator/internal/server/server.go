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
	pb "orchestrator/proto"
	"strings"
)

type OrchestratorServer struct {
	pb.UnimplementedOrchestratorServer
	Redis         *redis.RedisManager
	RabbitChannel *amqp.Channel
}

type VideoToAnalyze struct {
	URL          string
	ExerciseName string
	Auth0Id      string
	VideoId      int64
}

type TaskMessage struct {
	URL          string `json:"url"`
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
	pb.RegisterOrchestratorServer(grpcServer, &OrchestratorServer{
		Redis:         redisManager,
		RabbitChannel: rabbitChannel,
	})

	log.Printf("Orchestrator Service running on port :%s", cfg.GRPCPort)

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *OrchestratorServer) AnalyzeVideo(ctx context.Context, req *pb.VideoToAnalyzeRequest) (*pb.VideoToAnalyzeResponse, error) {

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
		URL:          req.Url,
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
	response := &pb.VideoToAnalyzeResponse{
		Success: true,
		Message: "Video analysis started successfully",
	}

	return response, nil
}
