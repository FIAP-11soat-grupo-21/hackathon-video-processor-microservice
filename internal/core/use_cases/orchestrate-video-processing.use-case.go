package use_cases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/google/uuid"

	"video_processor_service/internal/core/domain/ports"
	"video_processor_service/internal/core/dto"
)

type OrchestrateVideoProcessingUseCase struct {
	videoProcessor ports.IVideoProcessor
	queuePublisher ports.IQueuePublisher
	queueURL       string
}

func NewOrchestrateVideoProcessingUseCase(
	videoProcessor ports.IVideoProcessor,
	queuePublisher ports.IQueuePublisher,
	queueURL string,
) *OrchestrateVideoProcessingUseCase {
	return &OrchestrateVideoProcessingUseCase{
		videoProcessor: videoProcessor,
		queuePublisher: queuePublisher,
		queueURL:       queueURL,
	}
}

func (uc *OrchestrateVideoProcessingUseCase) Execute(
	ctx context.Context,
	request dto.ProcessVideoRequestDTO,
) (*dto.ProcessVideoResponseDTO, error) {
	duration, err := uc.videoProcessor.GetVideoDuration(ctx, request.Bucket, request.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to get video duration: %w", err)
	}

	timestamps := uc.calculateTimestamps(duration, request.ChunkIntervalSeconds)

	jobID := uuid.New().String()

	for i, timestamp := range timestamps {
		message := dto.FrameExtractionMessageDTO{
			JobID:     jobID,
			Bucket:    request.Bucket,
			Key:       request.Key,
			Timestamp: timestamp,
			Index:     i,
		}

		messageJSON, err := json.Marshal(message)
		if err != nil {
			log.Printf("Failed to marshal message for timestamp %.2f: %v", timestamp, err)
			continue
		}

		err = uc.queuePublisher.SendMessage(ctx, uc.queueURL, string(messageJSON))
		if err != nil {
			log.Printf("Failed to enqueue message for timestamp %.2f: %v", timestamp, err)
			continue
		}
	}

	log.Printf("Job %s: Enqueued %d frame extraction tasks", jobID, len(timestamps))

	return &dto.ProcessVideoResponseDTO{
		JobID:           jobID,
		Status:          "processing",
		EstimatedFrames: len(timestamps),
		VideoDuration:   duration,
	}, nil
}

func (uc *OrchestrateVideoProcessingUseCase) calculateTimestamps(duration, interval float64) []float64 {
	var timestamps []float64

	numFrames := int(math.Ceil(duration / interval))

	for i := 0; i < numFrames; i++ {
		timestamp := float64(i) * interval
		if timestamp <= duration {
			timestamps = append(timestamps, timestamp)
		}
	}

	return timestamps
}
