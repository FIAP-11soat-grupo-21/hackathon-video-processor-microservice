package dto

type FrameExtractionMessageDTO struct {
	JobID     string  `json:"jobId"`
	Bucket    string  `json:"bucket"`
	Key       string  `json:"key"`
	Timestamp float64 `json:"timestamp"`
	Index     int     `json:"index"`
}

type UserDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ChunkUploadedDTO struct {
	Bucket            string   `json:"bucket"`
	VideoObjectID     string   `json:"video_object_id"`
	User              UserDTO  `json:"user"`
	ImageDestination  string   `json:"image_destination"`
	FramePerSecond    int      `json:"frame_per_second"`
	ChunkPart         int      `json:"chunk_part"`
}

type ChunkProcessedDTO struct {
	VideoID   string  `json:"video_id"`
	User      UserDTO `json:"user"`
	ChunkPart int     `json:"chunk_part"`
	Status    string  `json:"status"`
}

type AllChunksProcessedDTO struct {
	VideoID        string  `json:"video_id"`
	User           UserDTO `json:"user"`
	ImagesLocation string  `json:"images_location"`
	Bucket         string  `json:"bucket"`
}

type VideoProcessingCompleteDTO struct {
	VideoID     string  `json:"video_id"`
	User        UserDTO `json:"user"`
	DownloadURL string  `json:"download_url"`
	Status      string  `json:"status"`
}
