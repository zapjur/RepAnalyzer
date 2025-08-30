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

	videoIDStr := fmt.Sprintf("%d", req.VideoId)
	if err := s.Redis.SetTaskStatus(videoIDStr, "bar_path", "pending"); err != nil {
		return errorResponse("Failed to set task status in Redis", err)
	}
	if err := s.Redis.SetTaskStatus(videoIDStr, "pose", "pending"); err != nil {
		return errorResponse("Failed to set task status in Redis", err)
	}

	auth0IDEdited := strings.ReplaceAll(req.Auth0Id, "|", "_")

	if err := s.publishToQueue(
		"bar_path_queue",
		buildTaskMessage(req, videoIDStr, auth0IDEdited, "bar_path_results_queue"),
	); err != nil {
		return errorResponse("Failed to publish bar_path task to RabbitMQ", err)
	}

	if err := s.publishToQueue(
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

func buildTaskMessage(req *anPb.VideoToAnalyzeRequest, videoIDStr, auth0IDEdited, replyQueue string) TaskMessage {
	return TaskMessage{
		Bucket:       req.Bucket,
		ObjectKey:    req.ObjectKey,
		ExerciseName: req.ExerciseName,
		VideoID:      videoIDStr,
		Auth0Id:      auth0IDEdited,
		ReplyQueue:   replyQueue,
	}
}

func (s *OrchestratorServer) publishToQueue(queueName string, msg TaskMessage) error {
	taskBody, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal task message: %v", err)
		return fmt.Errorf("marshal error: %w", err)
	}

	if err := s.RabbitChannel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         taskBody,
			DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		log.Printf("Failed to publish message to RabbitMQ: %v", err)
		return fmt.Errorf("publish error: %w", err)
	}

	return nil
}
