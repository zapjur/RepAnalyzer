package server

import (
	"context"

	"google.golang.org/grpc"
	pb "user-service/proto"
)

type grpcServer struct {
	pb.UnimplementedUserServiceServer
}

func Register(s *grpc.Server) {
	pb.RegisterUserServiceServer(s, &grpcServer{})
}

func (s *grpcServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	return &pb.GetUserResponse{
		Auth0Id: req.Auth0Id,
		Email:   "example@example.com",
		Name:    "Jan Kowalski",
		Exists:  true,
	}, nil
}
