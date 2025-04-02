package minio

import (
	"context"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	Minio     *minio.Client
	endpoint  string
	accessKey string
	secretKey string
	useSSL    bool
}

func NewMinioClient() *Client {
	client := &Client{
		endpoint:  "minio:9000",
		accessKey: "admin",
		secretKey: "admin123",
		useSSL:    false,
	}

	client.connectWithRetry()

	return client
}

func (c *Client) connectWithRetry() {
	var err error
	for {
		c.Minio, err = minio.New(c.endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(c.accessKey, c.secretKey, ""),
			Secure: c.useSSL,
		})

		if err == nil {
			log.Println("Successfully connected to MinIO")
			break
		}

		log.Printf("Failed to connect to MinIO (%v). Retrying in 5 seconds...\n", err)
		time.Sleep(5 * time.Second)
	}
}

func (c *Client) EnsureBucketExists(bucket string) {
	ctx := context.Background()
	var err error
	var exists bool

	const maxRetries = 5
	const retryInterval = 3 * time.Second

	for i := 1; i <= maxRetries; i++ {
		exists, err = c.Minio.BucketExists(ctx, bucket)
		if err == nil {
			break
		}
		log.Printf("Bucket check attempt %d/%d failed: %v. Retrying in %v...\n", i, maxRetries, err, retryInterval)
		time.Sleep(retryInterval)
		c.connectWithRetry()
	}

	if err != nil {
		log.Fatalf("Failed to check bucket after %d attempts: %v", maxRetries, err)
	}

	if !exists {
		err = c.Minio.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}
	}
}
