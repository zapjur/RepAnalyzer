package server

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"orchestrator/internal/config"
	pb "orchestrator/proto"
)

type OrchestratorServer struct {
	pb.UnimplementedOrchestratorServer
}

type VideoToAnalyze struct {
	URL          string
	ExerciseName string
	Auth0Id      string
	VideoId      int64
}

func StartGRPCServer(cfg *config.Config) {
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrchestratorServer(grpcServer, &OrchestratorServer{})

	log.Printf("Orchestrator Service running on port :%s", cfg.GRPCPort)

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *OrchestratorServer) AnalyzeVideo(ctx context.Context, req *pb.VideoToAnalyzeRequest) (*pb.VideoToAnalyzeResponse, error) {

	// sending video data to appropriate queue

	// send back response

	return nil, nil
}
