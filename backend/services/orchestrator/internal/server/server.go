package server

import pb "orchestrator/proto"

type OrchestratorServer struct {
	pb.UnimplementedOrchestratorServer
}

type VideoToAnalyze struct {
	URL          string
	ExerciseName string
	Auth0Id      string
	VideoId      string
}
