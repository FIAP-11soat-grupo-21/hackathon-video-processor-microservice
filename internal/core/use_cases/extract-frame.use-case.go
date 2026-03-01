package use_cases

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"video_processor_service/internal/core/domain/ports"
	"video_processor_service/internal/core/dto"
)

type ExtractFrameUseCase struct {
	videoProcessor ports.IVideoProcessor
	storageService ports.IStorageService
	outputBucket   string
}

func NewExtractFrameUseCase(
	videoProcessor ports.IVideoProcessor,
	storageService ports.IStorageService,
	outputBucket string,
) *ExtractFrameUseCase {
	return &ExtractFrameUseCase{
		videoProcessor: videoProcessor,
		storageService: storageService,
		outputBucket:   outputBucket,
	}
}

func (uc *ExtractFrameUseCase) Execute(
	ctx context.Context,
	message dto.FrameExtractionMessageDTO,
) error {
	log.Printf("Extracting frame at %.2fs for job %s (index %d)", message.Timestamp, message.JobID, message.Index)

	frameData, err := uc.videoProcessor.ExtractFrame(ctx, message.Bucket, message.Key, message.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to extract frame: %w", err)
	}

	frameKey := fmt.Sprintf("frames/%s/frame_%d_%.2fs.jpg", message.JobID, message.Index, message.Timestamp)

	err = uc.storageService.UploadObject(
		ctx,
		uc.outputBucket,
		frameKey,
		bytes.NewReader(frameData),
		"image/jpeg",
	)

	if err != nil {
		return fmt.Errorf("failed to upload frame: %w", err)
	}

	log.Printf("Frame uploaded successfully: %s", frameKey)

	return nil
}
