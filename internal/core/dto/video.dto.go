package dto

type ProcessVideoRequestDTO struct {
	Bucket               string  `json:"bucket" binding:"required"`
	Key                  string  `json:"key" binding:"required"`
	ChunkIntervalSeconds float64 `json:"chunkIntervalSeconds" binding:"required,gt=0"`
}

type ProcessVideoResponseDTO struct {
	JobID           string  `json:"jobId"`
	Status          string  `json:"status"`
	EstimatedFrames int     `json:"estimatedFrames"`
	VideoDuration   float64 `json:"videoDuration"`
}

type FrameExtractionMessageDTO struct {
	JobID     string  `json:"jobId"`
	Bucket    string  `json:"bucket"`
	Key       string  `json:"key"`
	Timestamp float64 `json:"timestamp"`
	Index     int     `json:"index"`
}
