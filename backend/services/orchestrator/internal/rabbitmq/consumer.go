package rabbitmq

import (
	"encoding/json"
	"log"
	"orchestrator/internal/redis"

	"github.com/streadway/amqp"
)

type BarpathResult struct {
	VideoID   string `json:"video_id"`
	Status    string `json:"status"`
	ResultURL string `json:"result_url,omitempty"`
	Message   string `json:"message,omitempty"`
}

func StartConsumers(ch *amqp.Channel, redisManager *redis.RedisManager) {
	if err := ConsumeBarpathResults(ch, redisManager); err != nil {
		log.Fatalf("Failed to start Barpath consumer: %v", err)
	}

}

func ConsumeBarpathResults(ch *amqp.Channel, redisManager *redis.RedisManager) error {
	msgs, err := ch.Consume(
		"bar_path_results_queue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var result BarpathResult
			if err = json.Unmarshal(msg.Body, &result); err != nil {
				log.Printf("Failed to parse barpath result: %v", err)
				continue
			}

			log.Printf("[barpath] Result: %+v", result)

			err = redisManager.SetTaskStatus(result.VideoID, "bar_path", result.Status)
			if err != nil {
				log.Printf("Failed to update Redis for barpath: %v", err)
			}
		}
	}()

	return nil
}
