package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Client struct {
	client    *minio.Client
	endpoint  string
	accessKey string
	secretKey string
	useSSL    bool
}

func NewClient() *Client {
	endpoint := getenvDefault("MINIO_ENDPOINT")
	accessKey := getenvDefault("MINIO_ACCESS_KEY")
	secretKey := getenvDefault("MINIO_SECRET_KEY")
	useSSL := getenvBool("MINIO_USE_SSL")

	c := &Client{
		endpoint:  endpoint,
		accessKey: accessKey,
		secretKey: secretKey,
		useSSL:    useSSL,
	}
	c.connectWithRetry()
	return c
}

func (c *Client) connectWithRetry() {
	var err error
	for {
		c.client, err = minio.New(c.endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(c.accessKey, c.secretKey, ""),
			Secure: c.useSSL,
		})
		if err == nil {
			log.Println("Connected to MinIO")
			break
		}
		log.Printf("Failed to connect to MinIO: %v. Retrying in 5s...\n", err)
		time.Sleep(5 * time.Second)
	}
}

func (c *Client) GeneratePresignedURL(ctx context.Context, bucket, objectKey string) (string, error) {
	reqParams := make(url.Values)
	presignedUrl, err := c.client.PresignedGetObject(ctx, bucket, objectKey, 10*time.Minute, reqParams)
	if err != nil {
		return "", err
	}
	return presignedUrl.String(), nil
}

func getenvDefault(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return ""
}

func getenvBool(key string) bool {
	v := os.Getenv(key)
	if v == "" {
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}
