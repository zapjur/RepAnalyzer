package types

type VideoWithURL struct {
	Id           int64  `json:"id"`
	Bucket       string `json:"bucket"`
	ObjectKey    string `json:"object_key"`
	ExerciseName string `json:"exercise_name"`
	CreatedAt    string `json:"created_at"`
	Url          string `json:"url"`
}
