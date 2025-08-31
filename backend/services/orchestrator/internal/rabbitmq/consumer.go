package rabbitmq

import (
	"context"
	"encoding/json"
	"log"
	"orchestrator/internal/client"
	"orchestrator/internal/redis"
	dbPb "orchestrator/proto/db"
	"orchestrator/types"
	"strconv"
)

func (r *RabbitClient) StartConsumers(redisManager *redis.RedisManager, grpcClient *client.Client) {
	if err := r.ConsumeBarpathResults(redisManager, grpcClient); err != nil {
		log.Fatalf("Failed to start Barpath consumer: %v", err)
	}

	if err := r.ConsumePoseResults(redisManager, grpcClient); err != nil {
		log.Fatalf("Failed to start Pose consumer: %v", err)
	}

}

func (r *RabbitClient) ConsumeBarpathResults(redisManager *redis.RedisManager, grpcClient *client.Client) error {
	msgs, err := r.Channel.Consume(
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
			var result types.TaskResult
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
					Bucket:    result.Bucket,
					ObjectKey: result.ObjectKey,
					Type:      "barpath",
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

				ready, err := checkReadiness(redisManager, result.VideoID)
				if err != nil {
					log.Printf("Failed to check readiness: %v", err)
					continue
				}
				if ready {
					err = r.PublishToQueue("analysis_queue", types.TaskMessage{
						VideoID:      result.VideoID,
						Bucket:       result.Bucket,
						ObjectKey:    result.ObjectKey,
						ReplyQueue:   "analysis_results_queue",
						Auth0Id:      result.Auth0Id,
						ExerciseName: result.ExerciseName,
					})
					if err != nil {
						log.Printf("Failed to publish to analysis queue: %v", err)
					} else {
						log.Printf("Published video %s to analysis queue", result.VideoID)
					}
				}
			}
		}
	}()

	return nil
}

func (r *RabbitClient) ConsumePoseResults(redisManager *redis.RedisManager, grpcClient *client.Client) error {
	msgs, err := r.Channel.Consume(
		"pose_results_queue",
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
			var result types.TaskResult
			if err = json.Unmarshal(msg.Body, &result); err != nil {
				log.Printf("Failed to parse pose result: %v", err)
				continue
			}

			log.Printf("[pose] Result: %+v", result)

			err = redisManager.SetTaskStatus(result.VideoID, "pose", result.Status)
			if err != nil {
				log.Printf("Failed to update Redis for pose: %v", err)
			}

			if result.Status == "success" {
				videoID, err := strconv.ParseInt(result.VideoID, 10, 64)
				if err != nil {
					log.Printf("Failed to parse video ID: %v", err)
					continue
				}
				res, err := grpcClient.DBService.SaveAnalysis(context.Background(), &dbPb.VideoAnalysisRequest{
					VideoId:   videoID,
					Bucket:    result.Bucket,
					ObjectKey: result.ObjectKey,
					Type:      "pose",
				})
				if err != nil {
					log.Printf("Failed to save analysis in DB: %v", err)
					continue
				}
				if !res.Success {
					log.Printf("Failed to save analysis in DB: %s", res.Message)
				} else {
					log.Printf("Pose analysis saved successfully in DB")
				}

				ready, err := checkReadiness(redisManager, result.VideoID)
				if err != nil {
					log.Printf("Failed to check readiness: %v", err)
					continue
				}
				if ready {
					err = r.PublishToQueue("analysis_queue", types.TaskMessage{
						VideoID:      result.VideoID,
						Bucket:       result.Bucket,
						ObjectKey:    result.ObjectKey,
						ReplyQueue:   "analysis_results_queue",
						Auth0Id:      result.Auth0Id,
						ExerciseName: result.ExerciseName,
					})
					if err != nil {
						log.Printf("Failed to publish to analysis queue: %v", err)
					} else {
						log.Printf("Published video %s to analysis queue", result.VideoID)
					}
				}
			}
		}
	}()

	return nil
}
