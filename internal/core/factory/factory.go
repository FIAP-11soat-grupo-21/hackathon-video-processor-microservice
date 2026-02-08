package factory

import (
	"video_processor_service/internal/adapters/driven/storage"
	"video_processor_service/internal/adapters/driven/video"
	"video_processor_service/internal/common/infra/queue"
	"video_processor_service/internal/core/domain/ports"
)

func NewQueueConsumer() ports.IQueueConsumer {
	return queue.NewSQSConsumer()
}

func NewQueuePublisher() ports.IQueuePublisher {
	return queue.NewSQSPublisher()
}

func NewStorageService() ports.IStorageService {
	return storage.NewS3StorageService()
}

func NewVideoProcessor() ports.IVideoProcessor {
	storageService := NewStorageService()
	return video.NewFFmpegVideoProcessor(storageService)
}
