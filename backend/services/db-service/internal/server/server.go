package server

import (
	"context"
	"db-service/internal/config"
	"db-service/internal/repository"
	"log"
	"net"
	"time"

	pb "db-service/proto"

	"google.golang.org/grpc"
)

type DBServer struct {
	pb.UnimplementedDBServiceServer
}

type Video struct {
	ObjectKey string
	Bucket    string
	CreatedAt time.Time
	ID        int64
}

func (s *DBServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("Checking user: %s", req.Auth0Id)

	user, err := repository.GetUserByAuth0ID(req.Auth0Id)
	if err != nil {
		return &pb.GetUserResponse{
			Success: false,
			Message: "Database error: " + err.Error(),
		}, err
	}

	if user == nil {
		log.Println("User not found, creating new user...")
		err := repository.CreateUser(req.Auth0Id, req.Email)
		if err != nil {
			return &pb.GetUserResponse{
				Success: false,
				Message: "Failed to create user: " + err.Error(),
			}, err
		}
		log.Println("User created successfully")
	}

	return &pb.GetUserResponse{
		Success: true,
		Message: "User exists or was created successfully",
	}, nil
}

func (s *DBServer) SaveUploadedVideo(ctx context.Context, req *pb.UploadVideoRequest) (*pb.UploadVideoResponse, error) {
	log.Printf("Saving video for user: %s", req.Auth0Id)

	videoID, err := repository.SaveUploadedVideo(req.Auth0Id, req.Bucket, req.ObjectKey, req.ExerciseName)
	if err != nil {
		return &pb.UploadVideoResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &pb.UploadVideoResponse{
		Success: true,
		Message: "Video saved successfully",
		VideoId: videoID,
	}, nil
}

func StartGRPCServer(cfg *config.Config) {
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDBServiceServer(grpcServer, &DBServer{})

	log.Printf("DB Service running on port :%s", cfg.GRPCPort)

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *DBServer) GetUserVideosByExercise(ctx context.Context, req *pb.GetUserVideosByExerciseRequest) (*pb.GetUserVideosByExerciseResponse, error) {
	log.Printf("Getting videos for user: %s and exercise: %s", req.Auth0Id, req.ExerciseName)

	videos, err := repository.GetUserVideosByExercise(req.Auth0Id, req.ExerciseName)
	if err != nil {
		return &pb.GetUserVideosByExerciseResponse{
			Success: false,
			Message: "Database error: " + err.Error(),
		}, err
	}

	var videoInfos []*pb.VideoInfo
	for _, v := range videos {
		videoInfos = append(videoInfos, &pb.VideoInfo{
			Id:           v.ID,
			ObjectKey:    v.ObjectKey,
			Bucket:       v.Bucket,
			ExerciseName: req.ExerciseName,
			Auth0Id:      req.Auth0Id,
			CreatedAt:    v.CreatedAt.Format(time.RFC3339),
		})
	}
	return &pb.GetUserVideosByExerciseResponse{
		Success: true,
		Message: "Videos retrieved successfully",
		Videos:  videoInfos,
	}, nil
}

func (s *DBServer) SaveAnalysis(ctx context.Context, req *pb.VideoAnalysisRequest) (*pb.SaveAnalysisResponse, error) {
	log.Printf("Saving analysis for video ID: %d", req.VideoId)

	_, err := repository.SaveAnalysis(req.VideoId, req.Type, req.Bucket, req.ObjectKey)
	if err != nil {
		return &pb.SaveAnalysisResponse{
			Success: false,
			Message: "Database error: " + err.Error(),
		}, err
	}

	return &pb.SaveAnalysisResponse{
		Success: true,
		Message: "Analysis saved successfully",
	}, nil
}
