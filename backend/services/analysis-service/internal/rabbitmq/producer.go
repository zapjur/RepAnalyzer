package rabbitmq

import (
	"analysis-service/internal/minio"
	"context"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

type RabbitClient struct {
	Connection  *amqp.Connection
	Channel     *amqp.Channel
	Context     context.Context
	MinioClient *minio.Client
}

func ConnectToRabbitMQ(uri string, ctx context.Context, mino *minio.Client) (*RabbitClient, error) {
	const retries = 5
	const retryDelay = 5 * time.Second

	var conn *amqp.Connection
	var channel *amqp.Channel
	var err error

	for i := 0; i < retries; i++ {
		conn, err = amqp.Dial(uri)
		if err == nil {
			channel, err = conn.Channel()
			if err == nil {
				log.Println("Successfully connected to RabbitMQ.")
				return &RabbitClient{conn, channel, ctx, mino}, nil
			}
		}

		log.Printf("RabbitMQ connection failed (%d/%d): %v", i+1, retries, err)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d retries: %w", retries, err)
}
