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

type ZipProcessorHandler struct {
	useCase *use_cases.ProcessZipUseCase
}

func NewZipProcessorHandler() (*ZipProcessorHandler, error) {
	storageService := factory.NewStorageService()
	messageBroker, err := factory.NewMessageBroker()
	if err != nil {
		return nil, err
	}

	snsTopicArn := os.Getenv("SNS_VIDEO_PROCESSING_COMPLETE_TOPIC")

	useCase := use_cases.NewProcessZipUseCase(
		storageService,
		messageBroker,
		snsTopicArn,
	)

	return &ZipProcessorHandler{useCase: useCase}, nil
}

func (h *ZipProcessorHandler) Handle(ctx context.Context, snsEvent events.SNSEvent) error {
	log.Printf("Processing %d SNS messages", len(snsEvent.Records))

	for _, record := range snsEvent.Records {
		var message dto.AllChunksProcessedDTO
		if err := json.Unmarshal([]byte(record.SNS.Message), &message); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		if err := h.useCase.Execute(ctx, message); err != nil {
			log.Printf("Error processing zip: %v", err)
			return err
		}
	}

	return nil
}
