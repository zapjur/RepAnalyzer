package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"orchestrator/internal/config"
	"orchestrator/internal/rabbitmq"
	"orchestrator/internal/redis"
	anPb "orchestrator/proto/analysis"
	"orchestrator/types"
	"strings"
)

type OrchestratorServer struct {
	anPb.UnimplementedOrchestratorServer
	Redis        *redis.RedisManager
	RabbitClient *rabbitmq.RabbitClient
}

type VideoToAnalyze struct {
	Bucket       string
	ObjectKey    string
	ExerciseName string
	Auth0Id      string
	VideoId      int64
}

func StartGRPCServer(cfg *config.Config, redisManager *redis.RedisManager, rabbitClient *rabbitmq.RabbitClient) {
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	anPb.RegisterOrchestratorServer(grpcServer, &OrchestratorServer{
		Redis:        redisManager,
		RabbitClient: rabbitClient,
	})

	log.Printf("Orchestrator Service running on port :%s", cfg.GRPCPort)

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *OrchestratorServer) AnalyzeVideo(ctx context.Context, req *anPb.VideoToAnalyzeRequest) (*anPb.VideoToAnalyzeResponse, error) {

	videoIDStr := fmt.Sprintf("%d", req.VideoId)
	if err := s.Redis.SetTaskStatus(videoIDStr, "bar_path", "pending"); err != nil {
		return errorResponse("Failed to set task status in Redis", err)
	}
	if err := s.Redis.SetTaskStatus(videoIDStr, "pose", "pending"); err != nil {
		return errorResponse("Failed to set task status in Redis", err)
	}

	auth0IDEdited := strings.ReplaceAll(req.Auth0Id, "|", "_")

	if err := s.RabbitClient.PublishToQueue(
		"bar_path_queue",
		buildTaskMessage(req, videoIDStr, auth0IDEdited, "bar_path_results_queue"),
	); err != nil {
		return errorResponse("Failed to publish bar_path task to RabbitMQ", err)
	}

	if err := s.RabbitClient.PublishToQueue(
		"pose_queue",
		buildTaskMessage(req, videoIDStr, auth0IDEdited, "pose_results_queue"),
	); err != nil {
		return errorResponse("Failed to publish pose task to RabbitMQ", err)
	}

	return &anPb.VideoToAnalyzeResponse{
		Success: true,
		Message: "Video analysis started successfully",
	}, nil
}

func buildTaskMessage(req *anPb.VideoToAnalyzeRequest, videoIDStr, auth0IDEdited, replyQueue string) types.TaskMessage {
	return types.TaskMessage{
		Bucket:       req.Bucket,
		ObjectKey:    req.ObjectKey,
		ExerciseName: req.ExerciseName,
		VideoID:      videoIDStr,
		Auth0Id:      auth0IDEdited,
		ReplyQueue:   replyQueue,
	}
}
