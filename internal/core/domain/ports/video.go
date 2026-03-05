package ports

import "context"

type IVideoProcessor interface {
	ExtractFrame(ctx context.Context, bucket, key string, timestamp float64) ([]byte, error)
}
