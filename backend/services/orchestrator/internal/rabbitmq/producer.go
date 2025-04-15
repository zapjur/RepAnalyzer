package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func ConnectToRabbitMQ(uri string) (*amqp.Connection, *amqp.Channel, error) {
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
				return conn, channel, nil
			}
		}

		log.Printf("RabbitMQ connection failed (%d/%d): %v", i+1, retries, err)
		time.Sleep(retryDelay)
	}

	return nil, nil, fmt.Errorf("failed to connect to RabbitMQ after %d retries: %w", retries, err)
}

func DeclareQueues(channel *amqp.Channel, queues []string) error {
	for _, queue := range queues {
		_, err := channel.QueueDeclare(
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
