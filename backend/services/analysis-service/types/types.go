package types

type AnalysisRequest struct {
	VideoID      string `json:"video_id"`
	Bucket       string `json:"bucket"`
	ObjectKey    string `json:"object_key"`
	Auth0Id      string `json:"auth0_id"`
	ReplyQueue   string `json:"reply_queue"`
	ExerciseName string `json:"exercise_name"`
}
