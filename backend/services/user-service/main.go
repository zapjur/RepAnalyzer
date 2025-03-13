package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "user-service/proto"
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func (s *server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("Request dla użytkownika: %s", req.Auth0Id)

	resp := &pb.GetUserResponse{
		Auth0Id: req.Auth0Id,
		Email:   "example@example.com",
		Name:    "Jan Kowalski",
		Exists:  true,
	}

	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Nie udało się nasłuchiwać: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &server{})

	log.Println("User Service działa na porcie :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Błąd uruchomienia serwera: %v", err)
	}
}
