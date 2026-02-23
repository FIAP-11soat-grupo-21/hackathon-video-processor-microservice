package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"video_processor_service/internal/common/config/env"
	"video_processor_service/internal/core/dto"
	"video_processor_service/internal/core/factory"
	"video_processor_service/internal/core/use_cases"
)

func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	log.Printf("Processing %d SQS messages", len(sqsEvent.Records))

	cfg := env.GetConfig()

	videoProcessor := factory.NewVideoProcessor()
	storageService := factory.NewStorageService()

	useCase := use_cases.NewExtractFrameUseCase(
		videoProcessor,
		storageService,
		cfg.S3.BucketName,
	)

	for _, record := range sqsEvent.Records {
		log.Printf("Processing message ID: %s", record.MessageId)

		var message dto.FrameExtractionMessageDTO
		if err := json.Unmarshal([]byte(record.Body), &message); err != nil {
			log.Printf("Error parsing message %s: %v", record.MessageId, err)
			continue
		}

		log.Printf("Extracting frame for job %s at %.2fs (index %d)",
			message.JobID, message.Timestamp, message.Index)

		if err := useCase.Execute(ctx, message); err != nil {
			log.Printf("Error extracting frame for message %s: %v", record.MessageId, err)
			return fmt.Errorf("failed to process message %s: %v", record.MessageId, err)
		}

		log.Printf("Successfully processed frame for job %s", message.JobID)
	}

	log.Printf("Successfully processed all %d messages", len(sqsEvent.Records))
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Starting Lambda handler")
	
	lambda.Start(handler)
}
