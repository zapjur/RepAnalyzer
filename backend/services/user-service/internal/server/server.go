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
		return nil, err
	}

	if user == nil {
		log.Println("User not found, creating new user...")
		err := repository.CreateUser(req.Auth0Id, req.Email)
		if err != nil {
			return nil, err
		}
		user, _ = repository.GetUserByAuth0ID(req.Auth0Id)
	}

	return &pb.GetUserResponse{
		Auth0Id: user.Auth0ID,
		Email:   user.Email,
		Exists:  true,
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
