package lambda

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"

	"video_processor_service/internal/core/dto"
	"video_processor_service/internal/core/factory"
	"video_processor_service/internal/core/use_cases"
)

type UpdateVideoChunkStatusHandler struct {
	useCase *use_cases.UpdateChunkStatusUseCase
}

func NewUpdateVideoChunkStatusHandler() (*UpdateVideoChunkStatusHandler, error) {
	tableName := os.Getenv("DYNAMODB_TABLE_NAME")
	snsTopicArn := os.Getenv("SNS_ALL_CHUNKS_PROCESSED_TOPIC")
	s3Bucket := os.Getenv("S3_BUCKET")

	repository, err := factory.NewVideoChunkRepository(tableName)
	if err != nil {
		return nil, err
	}

	messageBroker, err := factory.NewMessageBroker()
	if err != nil {
		return nil, err
	}

	useCase := use_cases.NewUpdateChunkStatusUseCase(
		repository,
		messageBroker,
		snsTopicArn,
		s3Bucket,
	)

	return &UpdateVideoChunkStatusHandler{useCase: useCase}, nil
}

func (h *UpdateVideoChunkStatusHandler) Handle(ctx context.Context, snsEvent events.SNSEvent) error {
	log.Printf("Processing %d SNS messages", len(snsEvent.Records))

	for _, record := range snsEvent.Records {
		var message dto.ChunkProcessedDTO
		if err := json.Unmarshal([]byte(record.SNS.Message), &message); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		if err := h.useCase.Execute(ctx, message); err != nil {
			log.Printf("Error updating chunk status: %v", err)
			return err
		}
	}

	return nil
}
