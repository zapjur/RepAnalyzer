package types

type VideoWithURL struct {
	Id           int64  `json:"id"`
	Bucket       string `json:"bucket"`
	ObjectKey    string `json:"object_key"`
	ExerciseName string `json:"exercise_name"`
	CreatedAt    string `json:"created_at"`
	Url          string `json:"url"`
}

type VideoAnalysisWithURL struct {
	Id        int64   `json:"id"`
	Bucket    string  `json:"bucket"`
	ObjectKey string  `json:"object_key"`
	Type      string  `json:"type"`
	Url       string  `json:"url"`
	CsvUrl    *string `json:"csv_url,omitempty"`
	VideoId   int64   `json:"video_id"`
}
