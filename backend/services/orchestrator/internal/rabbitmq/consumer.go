package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"orchestrator/internal/client"
	"orchestrator/internal/redis"
	dbPb "orchestrator/proto/db"
	"strconv"
)

type BarpathResult struct {
	VideoID   string `json:"video_id"`
	Status    string `json:"status"`
	ResultURL string `json:"result_url,omitempty"`
	Message   string `json:"message,omitempty"`
}

func StartConsumers(ch *amqp.Channel, redisManager *redis.RedisManager, grpcClient *client.Client) {
	if err := ConsumeBarpathResults(ch, redisManager, grpcClient); err != nil {
		log.Fatalf("Failed to start Barpath consumer: %v", err)
	}

}

func ConsumeBarpathResults(ch *amqp.Channel, redisManager *redis.RedisManager, grpcClient *client.Client) error {
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

			if result.Status == "success" {
				videoID, err := strconv.ParseInt(result.VideoID, 10, 64)
				if err != nil {
					log.Printf("Failed to parse video ID: %v", err)
					continue
				}
				res, err := grpcClient.DBService.SaveAnalysis(context.Background(), &dbPb.VideoAnalysisRequest{
					VideoId:   videoID,
					ResultUrl: result.ResultURL,
					Type:      "bar_path",
				})
				if err != nil {
					log.Printf("Failed to save analysis in DB: %v", err)
					continue
				}
				if !res.Success {
					log.Printf("Failed to save analysis in DB: %s", res.Message)
				} else {
					log.Printf("Analysis saved successfully in DB")
				}
			}
		}
	}()

	return nil
}
