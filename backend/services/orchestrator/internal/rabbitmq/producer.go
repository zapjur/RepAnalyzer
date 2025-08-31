package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"orchestrator/types"
	"time"
)

type RabbitClient struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func ConnectToRabbitMQ(uri string) (*RabbitClient, error) {
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
				return &RabbitClient{conn, channel}, nil
			}
		}

		log.Printf("RabbitMQ connection failed (%d/%d): %v", i+1, retries, err)
		time.Sleep(retryDelay)
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d retries: %w", retries, err)
}

func (r *RabbitClient) DeclareQueues(queues []string) error {
	for _, queue := range queues {
		_, err := r.Channel.QueueDeclare(
			queue, // queue name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			log.Printf("Failed to declare RabbitMQ queue '%s': %v", queue, err)
			return err
		}
		log.Printf("Queue '%s' declared successfully.", queue)
	}
	return nil
}

func (r *RabbitClient) PublishToQueue(queueName string, msg types.TaskMessage) error {
	taskBody, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal task message: %v", err)
		return fmt.Errorf("marshal error: %w", err)
	}

	if err := r.Channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         taskBody,
			DeliveryMode: amqp.Persistent,
		},
	); err != nil {
		log.Printf("Failed to publish message to RabbitMQ: %v", err)
		return fmt.Errorf("publish error: %w", err)
	}

	return nil
}
