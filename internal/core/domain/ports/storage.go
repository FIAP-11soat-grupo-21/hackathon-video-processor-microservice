package ports

import (
	"context"
	"io"
)

type S3Object struct {
	Key  string
	Size int64
}

type IStorageService interface {
	GetObjectMetadata(ctx context.Context, bucket, key string) (int64, error)
	GetObjectRange(ctx context.Context, bucket, key string, start, end int64) (io.ReadCloser, error)
	GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	UploadObject(ctx context.Context, bucket, key string, body io.Reader, contentType string) error
	ListObjects(ctx context.Context, bucket, prefix string) ([]S3Object, error)
	GeneratePresignedURL(ctx context.Context, bucket, key string, expirationSeconds int) (string, error)
}
