package ports

import "context"

type IVideoProcessor interface {
	GetVideoDuration(ctx context.Context, bucket, key string) (float64, error)
	ExtractFrame(ctx context.Context, bucket, key string, timestamp float64) ([]byte, error)
}
