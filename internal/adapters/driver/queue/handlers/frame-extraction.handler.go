package handlers

import (
	"context"
	"encoding/json"
	"log"
	"video_processor_service/internal/common/config/env"
	"video_processor_service/internal/core/dto"
	"video_processor_service/internal/core/factory"
	"video_processor_service/internal/core/use_cases"
)

func ExtractFrame(msgBody []byte) error {
	log.Println("Processing frame extraction message:", string(msgBody))

	var message dto.FrameExtractionMessageDTO
	if err := json.Unmarshal(msgBody, &message); err != nil {
		log.Printf("Error parsing message: %v", err)
		return err
	}

	cfg := env.GetConfig()
	ctx := context.Background()

	videoProcessor := factory.NewVideoProcessor()
	storageService := factory.NewStorageService()

	useCase := use_cases.NewExtractFrameUseCase(
		videoProcessor,
		storageService,
		cfg.S3.OutputBucket,
	)

	err := useCase.Execute(ctx, message)
	if err != nil {
		log.Printf("Error extracting frame: %v", err)
		return err
	}

	log.Printf("Frame extracted successfully for job %s at %.2fs", message.JobID, message.Timestamp)
	return nil
}
