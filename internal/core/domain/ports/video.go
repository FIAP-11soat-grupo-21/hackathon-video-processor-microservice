package ports

import "context"

type IVideoProcessor interface {
	ExtractFrame(ctx context.Context, bucket, key string, timestamp float64) ([]byte, error)
	ExtractFramesFromVideo(ctx context.Context, videoPath, outputPattern string, fps int) error
}
