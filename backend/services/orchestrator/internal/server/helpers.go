package server

import (
	"fmt"
	"log"
	anPb "orchestrator/proto/analysis"
)

func errorResponse(msg string, err error) (*anPb.VideoToAnalyzeResponse, error) {
	log.Printf("%s: %v", msg, err)
	return &anPb.VideoToAnalyzeResponse{
		Success: false,
		Message: msg,
	}, fmt.Errorf("%s: %w", msg, err)
}
