package server

import (
	"fmt"
	"log"
	pb "orchestrator/proto"
)

func errorResponse(msg string, err error) (*pb.VideoToAnalyzeResponse, error) {
	log.Printf("%s: %v", msg, err)
	return &pb.VideoToAnalyzeResponse{
		Success: false,
		Message: msg,
	}, fmt.Errorf("%s: %w", msg, err)
}
