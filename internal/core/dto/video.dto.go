package dto

type FrameExtractionMessageDTO struct {
	JobID     string  `json:"jobId"`
	Bucket    string  `json:"bucket"`
	Key       string  `json:"key"`
	Timestamp float64 `json:"timestamp"`
	Index     int     `json:"index"`
}
