package types

type TaskMessage struct {
	Bucket       string `json:"bucket"`
	ObjectKey    string `json:"object_key"`
	ExerciseName string `json:"exercise_name"`
	VideoID      string `json:"video_id"`
	Auth0Id      string `json:"auth0_id"`
	ReplyQueue   string `json:"reply_queue"`
}

type TaskResult struct {
	VideoID      string `json:"video_id"`
	Status       string `json:"status"`
	Bucket       string `json:"bucket"`
	ObjectKey    string `json:"object_key"`
	Message      string `json:"message,omitempty"`
	Auth0Id      string `json:"auth0_id"`
	ExerciseName string `json:"exercise_name"`
}
