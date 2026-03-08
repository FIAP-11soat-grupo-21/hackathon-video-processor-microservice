package factory

import (
	"video_processor_service/internal/adapters/driven/database"
	"video_processor_service/internal/adapters/driven/messaging"
	"video_processor_service/internal/adapters/driven/storage"
	"video_processor_service/internal/adapters/driven/video"
	"video_processor_service/internal/core/domain/ports"
)

func NewStorageService() ports.IStorageService {
	return storage.NewS3StorageService()
}

func NewVideoProcessor() ports.IVideoProcessor {
	storageService := NewStorageService()
	return video.NewFFmpegVideoProcessor(storageService)
}

func NewMessageBroker() (ports.IMessageBroker, error) {
	return messaging.NewAWSMessageBroker()
}

func NewVideoChunkRepository(tableName string) (ports.IVideoChunkRepository, error) {
	return database.NewDynamoDBRepository(tableName)
}
