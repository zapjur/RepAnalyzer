package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"log"

	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	Minio *minio.Client
}

func NewMinioClient() *Client {
	minioClient, err := minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("admin", "admin123", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("failed to initialize MinIO client: %v", err)
	}

	return &Client{Minio: minioClient}
}

func (c *Client) EnsureBucketExists(bucket string) {
	ctx := context.Background()

	exists, err := c.Minio.BucketExists(ctx, bucket)
	if err != nil {
		log.Fatalf("failed to check bucket: %v", err)
	}
	if !exists {
		err = c.Minio.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("failed to create bucket: %v", err)
		}
	}
}
