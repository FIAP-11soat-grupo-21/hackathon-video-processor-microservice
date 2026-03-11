package use_cases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"video_processor_service/internal/core/domain"
	"video_processor_service/internal/core/domain/ports"
	"video_processor_service/internal/core/dto"
)

type UpdateChunkStatusUseCase struct {
	repository    ports.IVideoChunkRepository
	messageBroker ports.IMessageBroker
	snsTopicArn   string
	s3Bucket      string
}

func NewUpdateChunkStatusUseCase(
	repository ports.IVideoChunkRepository,
	messageBroker ports.IMessageBroker,
	snsTopicArn string,
	s3Bucket string,
) *UpdateChunkStatusUseCase {
	return &UpdateChunkStatusUseCase{
		repository:    repository,
		messageBroker: messageBroker,
		snsTopicArn:   snsTopicArn,
		s3Bucket:      s3Bucket,
	}
}

func (uc *UpdateChunkStatusUseCase) Execute(ctx context.Context, message dto.ChunkProcessedDTO) error {
	log.Printf("Updating chunk %d status for video %s", message.ChunkPart, message.VideoID)

	chunk := domain.VideoChunk{
		VideoID:   message.VideoID,
		ChunkPart: message.ChunkPart,
		Status:    message.Status,
		UserID:    message.User.ID,
		UserName:  message.User.Name,
		UserEmail: message.User.Email,
		UpdatedAt: time.Now().UTC(),
	}

	if err := uc.repository.SaveChunk(ctx, chunk); err != nil {
		return fmt.Errorf("failed to save chunk to repository: %w", err)
	}

	log.Printf("Chunk %d saved to repository", message.ChunkPart)

	allCompleted, err := uc.checkAllChunksCompleted(ctx, message.VideoID)
	if err != nil {
		return fmt.Errorf("failed to check chunks status: %w", err)
	}

	if allCompleted {
		log.Printf("All chunks completed for video %s, publishing to SNS", message.VideoID)

		videoIDParts := strings.Split(message.VideoID, "/")
		videoID := message.VideoID
		if len(videoIDParts) >= 2 {
			videoID = videoIDParts[1]
		}

		allChunksMessage := dto.AllChunksProcessedDTO{
			VideoID:        message.VideoID,
			User:           message.User,
			ImagesLocation: fmt.Sprintf("videos/%s/frames", videoID),
			Bucket:         uc.s3Bucket,
		}

		messageBytes, err := json.Marshal(allChunksMessage)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		if err := uc.messageBroker.Publish(ctx, uc.snsTopicArn, string(messageBytes)); err != nil {
			return fmt.Errorf("failed to publish to SNS: %w", err)
		}

		log.Printf("Published all-chunks-processed event for video %s", message.VideoID)
	}

	return nil
}

func (uc *UpdateChunkStatusUseCase) checkAllChunksCompleted(ctx context.Context, videoID string) (bool, error) {
	chunks, err := uc.repository.GetChunksByVideoID(ctx, videoID)
	if err != nil {
		return false, fmt.Errorf("failed to get chunks: %w", err)
	}

	if len(chunks) == 0 {
		return false, nil
	}

	for _, chunk := range chunks {
		if chunk.Status != "completed" {
			return false, nil
		}
	}

	log.Printf("All %d chunks are completed for video %s", len(chunks), videoID)
	return true, nil
}
