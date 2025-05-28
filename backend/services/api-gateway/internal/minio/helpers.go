package minio

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
)

func WaitForObject(client *minio.Client, bucket, objectKey string, maxRetries int, delay time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		_, err := client.StatObject(context.Background(), bucket, objectKey, minio.StatObjectOptions{})
		if err == nil {
			return nil
		}
		time.Sleep(delay)
	}
	return fmt.Errorf("object %s/%s not available after %d retries", bucket, objectKey, maxRetries)
}
