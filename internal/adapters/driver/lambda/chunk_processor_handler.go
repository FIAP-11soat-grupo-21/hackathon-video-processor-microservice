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

type ChunkProcessorHandler struct {
	useCase *use_cases.ProcessChunkUseCase
}

func NewChunkProcessorHandler() (*ChunkProcessorHandler, error) {
	videoProcessor := factory.NewVideoProcessor()
	storageService := factory.NewStorageService()
	messageBroker, err := factory.NewMessageBroker()
	if err != nil {
		return nil, err
	}

	snsTopicArn := os.Getenv("SNS_CHUNK_PROCESSED_TOPIC")

	useCase := use_cases.NewProcessChunkUseCase(
		videoProcessor,
		storageService,
		messageBroker,
		snsTopicArn,
	)

	return &ChunkProcessorHandler{useCase: useCase}, nil
}

func (h *ChunkProcessorHandler) Handle(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Printf("Processing %d SQS messages", len(sqsEvent.Records))

	for _, record := range sqsEvent.Records {
		log.Printf("SQS Record Body: %s", record.Body)

		var message dto.ChunkUploadedDTO
		if err := json.Unmarshal([]byte(record.Body), &message); err != nil {
			log.Printf("Error parsing chunk message: %v", err)
			continue
		}

		log.Printf("Parsed message: bucket=%s, video_object_id=%s, chunk_part=%d", 
			message.Bucket, message.VideoObjectID, message.ChunkPart)

		if err := h.useCase.Execute(ctx, message); err != nil {
			log.Printf("Error processing chunk: %v", err)
			return err
		}
	}

	return nil
}
