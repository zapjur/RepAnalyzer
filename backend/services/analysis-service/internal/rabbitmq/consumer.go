package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type AnalysisRequest struct {
	VideoID   string `json:"video_id"`
	VideoType string `json:"video_type"`
	Bucket    string `json:"bucket"`
	ObjectKey string `json:"object_key"`
}

func StartConsumers(ctx context.Context, ch *amqp.Channel) error {
	return ConsumeAnalysisRequests(ctx, ch)
}

func ConsumeAnalysisRequests(ctx context.Context, ch *amqp.Channel) error {
	if err := ch.Qos(10, 0, false); err != nil {
		return err
	}

	const consumerTag = "analysis-consumer"
	msgs, err := ch.Consume(
		"analysis_queue",
		consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = ch.Cancel(consumerTag, false)
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				var req AnalysisRequest
				if err := json.Unmarshal(msg.Body, &req); err != nil {
					log.Printf("Failed to parse analysis request: %v", err)
					_ = msg.Nack(false, false)
					continue
				}

				log.Printf("[analysis] Request: %+v", req)

				// TODO: logic to process the analysis request
				_ = msg.Ack(false)
			}
		}
	}()

	return nil
}
