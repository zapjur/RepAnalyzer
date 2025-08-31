package rabbitmq

import "orchestrator/internal/redis"

func checkReadiness(redisManager *redis.RedisManager, videoID string) (bool, error) {
	barpathStatus, err := redisManager.GetTaskStatus(videoID, "bar_path")
	if err != nil {
		return false, err
	}

	poseStatus, err := redisManager.GetTaskStatus(videoID, "pose")
	if err != nil {
		return false, err
	}
	return barpathStatus.Status == "success" && poseStatus.Status == "success", nil
}
