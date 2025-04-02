package server

import (
	"context"
	"log"
	"net"
	"user-service/internal/config"
	"user-service/internal/repository"

	pb "user-service/proto"

	"google.golang.org/grpc"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
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
