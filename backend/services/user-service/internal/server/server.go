package server

import (
	"context"
	"log"
	"net"
	"time"
	"user-service/internal/config"
	"user-service/internal/repository"

	pb "user-service/proto"

	"google.golang.org/grpc"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
}

type Video struct {
	URL       string
	CreatedAt time.Time
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
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

func (s *UserServer) SaveUploadedVideo(ctx context.Context, req *pb.UploadVideoRequest) (*pb.UploadVideoResponse, error) {
	log.Printf("Saving video for user: %s", req.Auth0Id)

	err := repository.SaveUploadedVideo(req.Auth0Id, req.Url, req.ExerciseName)
	if err != nil {
		return &pb.UploadVideoResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &pb.UploadVideoResponse{
		Success: true,
		Message: "Video saved successfully",
	}, nil
}

func StartGRPCServer(cfg *config.Config) {
	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &UserServer{})

	log.Printf("User Service running on port :%s", cfg.GRPCPort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *UserServer) GetUserVideosByExercise(ctx context.Context, req *pb.GetUserVideosByExerciseRequest) (*pb.GetUserVideosByExerciseResponse, error) {
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
			Url:          v.URL,
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
