package ports

import (
	"context"
	"io"
)

type IStorageService interface {
	GetObjectMetadata(ctx context.Context, bucket, key string) (int64, error)
	GetObjectRange(ctx context.Context, bucket, key string, start, end int64) (io.ReadCloser, error)
	UploadObject(ctx context.Context, bucket, key string, body io.Reader, contentType string) error
}
