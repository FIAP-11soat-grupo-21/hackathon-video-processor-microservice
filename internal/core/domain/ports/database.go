package ports

import (
	"context"
	"video_processor_service/internal/core/domain"
)

type IVideoChunkRepository interface {
	SaveChunk(ctx context.Context, chunk domain.VideoChunk) error
	GetChunksByVideoID(ctx context.Context, videoID string) ([]domain.VideoChunk, error)
}
