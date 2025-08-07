package service

import (
	"access-service/internal/types"
	"context"
	"fmt"
)

type MinioClient interface {
	GeneratePresignedURL(ctx context.Context, bucket, objectKey string) (string, error)
}

type GrpcClient interface {
	UserOwnsVideo(ctx context.Context, auth0ID string, videoID int64) (types.GrpcDBServiceResponse, error)
}

type AccessService struct {
	minioClient MinioClient
	grpcClient  GrpcClient
}

func NewAccessService(minio MinioClient, grpc GrpcClient) *AccessService {
	return &AccessService{
		minioClient: minio,
		grpcClient:  grpc,
	}
}

func (s *AccessService) GeneratePresignedURL(ctx context.Context, auth0ID string, videoID int64) (string, error) {
	resp, err := s.grpcClient.UserOwnsVideo(ctx, auth0ID, videoID)
	if err != nil {
		return "", fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !resp.Owned {
		return "", fmt.Errorf("user does not own this video")
	}

	return s.minioClient.GeneratePresignedURL(ctx, resp.Bucket, resp.ObjectKey)
}
