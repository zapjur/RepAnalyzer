package minio

import (
	"analysis-service/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"time"
)

type Client struct {
	Minio     *minio.Client
	endpoint  string
	accessKey string
	secretKey string
	useSSL    bool
}

func NewMinioClient(cfg *config.Config) *Client {
	c := &Client{
		endpoint:  cfg.MinioEndpoint,
		accessKey: cfg.MinioAccessKey,
		secretKey: cfg.MinioSecretKey,
		useSSL:    cfg.MinieUseSSL,
	}
	c.connectWithRetry()
	log.Printf("MinIO presign client ready (endpoint=%s, ssl=%v)\n", c.endpoint, c.useSSL)
	return c
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
