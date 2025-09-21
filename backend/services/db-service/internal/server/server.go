package server

import (
	"context"
	"db-service/internal/config"
	"db-service/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"time"

	pb "db-service/proto"

	"google.golang.org/grpc"
)

type DBServer struct {
	pb.UnimplementedDBServiceServer
	repo *repository.Repository
}

func (s *DBServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("Checking user: %s", req.Auth0Id)

	user, err := s.repo.GetUserByAuth0ID(req.Auth0Id)
	if err != nil {
		return &pb.GetUserResponse{
			Success: false,
			Message: "Database error: " + err.Error(),
		}, err
	}

	if user == nil {
		log.Println("User not found, creating new user...")
		err := s.repo.CreateUser(req.Auth0Id, req.Email)
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

	videoID, err := s.repo.SaveUploadedVideo(req.Auth0Id, req.Bucket, req.ObjectKey, req.ExerciseName)
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

func StartGRPCServer(cfg *config.Config, db *pgxpool.Pool) {
	repo := repository.NewRepository(db)
	server := &DBServer{repo: repo}

	grpcServer := grpc.NewServer()
	pb.RegisterDBServiceServer(grpcServer, server)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("DB Service running on port :%s", cfg.GRPCPort)

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *DBServer) GetUserVideosByExercise(ctx context.Context, req *pb.GetUserVideosByExerciseRequest) (*pb.GetUserVideosByExerciseResponse, error) {
	log.Printf("Getting videos for user: %s and exercise: %s", req.Auth0Id, req.ExerciseName)

	videos, err := s.repo.GetUserVideosByExercise(req.Auth0Id, req.ExerciseName)
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

	_, err := s.repo.SaveAnalysis(req.VideoId, req.Type, req.Bucket, req.ObjectKey)
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

func (s *DBServer) CheckOwnership(ctx context.Context, req *pb.CheckOwnershipRequest) (*pb.CheckOwnershipResponse, error) {
	if req == nil {
		log.Println("CheckOwnership called with nil request")
		return nil, status.Error(codes.InvalidArgument, "nil request")
	}
	if s.repo == nil {
		log.Println("CheckOwnership called but repository is not initialized")
		return nil, status.Error(codes.Internal, "repository not initialized")
	}

	video, err := s.repo.GetVideoByID(req.VideoId)
	if err != nil {
		log.Printf("Failed to fetch video with ID %d: %v", req.VideoId, err)
		return nil, status.Errorf(codes.Internal, "failed to fetch video: %v", err)
	}
	if video == nil {
		log.Printf("Video with ID %d not found", req.VideoId)
		return &pb.CheckOwnershipResponse{
			Owned:     false,
			Message:   "video not found",
			ObjectKey: "",
			Bucket:    "",
		}, status.Error(codes.NotFound, "video not found")
	}

	user, err := s.repo.GetUserByAuth0ID(req.Auth0Id)
	if err != nil {
		log.Printf("Failed to fetch user with Auth0 ID %s: %v", req.Auth0Id, err)
		return &pb.CheckOwnershipResponse{
			Owned:     false,
			Message:   "database error: " + err.Error(),
			ObjectKey: "",
			Bucket:    "",
		}, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
	}
	if user == nil {
		log.Printf("User with Auth0 ID %s not found", req.Auth0Id)
		return &pb.CheckOwnershipResponse{
			Owned:     false,
			Message:   "user not found",
			ObjectKey: "",
			Bucket:    "",
		}, status.Error(codes.NotFound, "user not found")
	}

	owned := video.UserID == user.ID

	log.Printf("User %s ownership check for video ID %d: %t", req.Auth0Id, req.VideoId, owned)
	return &pb.CheckOwnershipResponse{
		Owned:     owned,
		Message:   "success",
		ObjectKey: video.ObjectKey,
		Bucket:    video.Bucket,
	}, nil
}

func (s *DBServer) GetVideoAnalysis(ctx context.Context, req *pb.GetVideoAnalysisRequest) (*pb.GetVideoAnalysisResponse, error) {
	log.Printf("Getting analysis for video ID: %d", req.VideoId)

	analyses, err := s.repo.GetVideoAnalysis(req.VideoId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch video analysis: %v", err)
	}

	var analysisInfos []*pb.VideoAnalysis
	for _, a := range analyses {
		analysisInfos = append(analysisInfos, &pb.VideoAnalysis{
			Id:        a.ID,
			Type:      a.Type,
			Bucket:    a.Bucket,
			ObjectKey: a.ObjectKey,
			VideoId:   a.VideoID,
		})
	}

	return &pb.GetVideoAnalysisResponse{
		Success:  true,
		Message:  "Analysis retrieved successfully",
		Analyses: analysisInfos,
	}, nil
}

func (s *DBServer) SaveAnalysisJSON(ctx context.Context, req *pb.SaveAnalysisJSONRequest) (*pb.SaveAnalysisJSONResponse, error) {
	log.Printf("SaveAnalysisJSON: video_id=%d payload_len=%d", req.VideoId, len(req.PayloadJson))

	id, err := s.repo.SaveAnalysisJSON(req.VideoId, req.PayloadJson)
	if err != nil {
		return &pb.SaveAnalysisJSONResponse{
			Success:    false,
			Message:    "Database error: " + err.Error(),
			AnalysisId: 0,
		}, err
	}

	return &pb.SaveAnalysisJSONResponse{
		Success:    true,
		Message:    "Analysis JSON saved successfully",
		AnalysisId: id,
	}, nil
}
