package rabbitmq

import (
	"analysis-service/internal/processing"
	"analysis-service/types"
	"encoding/json"
	"log"
)

func (r *RabbitClient) StartConsumers() error {
	return r.ConsumeAnalysisRequests()
}

func (r *RabbitClient) ConsumeAnalysisRequests() error {
	if err := r.Channel.Qos(10, 0, false); err != nil {
		return err
	}

	const consumerTag = "analysis-consumer"
	msgs, err := r.Channel.Consume(
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
			case <-r.Context.Done():
				_ = r.Channel.Cancel(consumerTag, false)
				return
			case msg, ok := <-msgs:
				if !ok {
					return
				}

				var req types.AnalysisRequest
				if err := json.Unmarshal(msg.Body, &req); err != nil {
					log.Printf("Failed to parse analysis request: %v", err)
					_ = msg.Nack(false, false)
					continue
				}

				log.Printf("[analysis] Request: %+v", req)

				err = processing.GenerateAnalysis(r.Context, r.MinioClient, r.grpcClient, req)
				if err != nil {
					log.Printf("Failed to process analysis request: %v", err)
					_ = msg.Nack(false, true)
					continue
				}
				_ = msg.Ack(false)
			}
		}
	}()

	return nil
}
