package use_cases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"video_processor_service/internal/core/domain/ports"
	"video_processor_service/internal/core/dto"
)

type ProcessChunkUseCase struct {
	videoProcessor ports.IVideoProcessor
	storageService ports.IStorageService
	messageBroker  ports.IMessageBroker
	snsTopicArn    string
}

func NewProcessChunkUseCase(
	videoProcessor ports.IVideoProcessor,
	storageService ports.IStorageService,
	messageBroker ports.IMessageBroker,
	snsTopicArn string,
) *ProcessChunkUseCase {
	return &ProcessChunkUseCase{
		videoProcessor: videoProcessor,
		storageService: storageService,
		messageBroker:  messageBroker,
		snsTopicArn:    snsTopicArn,
	}
}

func (uc *ProcessChunkUseCase) Execute(ctx context.Context, message dto.ChunkUploadedDTO) error {
	log.Printf("Processing chunk %d for video %s", message.ChunkPart, message.VideoObjectID)

	tmpDir := "/tmp/video-processing"
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	videoPath := filepath.Join(tmpDir, "video.mp4")
	outputPattern := filepath.Join(tmpDir, "frame_%04d.jpg")

	log.Printf("Downloading video from s3://%s/%s", message.Bucket, message.VideoObjectID)
	videoReader, err := uc.storageService.GetObject(ctx, message.Bucket, message.VideoObjectID)
	if err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}
	defer videoReader.Close()

	videoFile, err := os.Create(videoPath)
	if err != nil {
		return fmt.Errorf("failed to create video file: %w", err)
	}
	defer videoFile.Close()

	if _, err := videoFile.ReadFrom(videoReader); err != nil {
		return fmt.Errorf("failed to write video file: %w", err)
	}

	log.Printf("Extracting frames at %d fps", message.FramePerSecond)
	if err := uc.videoProcessor.ExtractFramesFromVideo(ctx, videoPath, outputPattern, message.FramePerSecond); err != nil {
		return fmt.Errorf("failed to extract frames: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(tmpDir, "frame_*.jpg"))
	if err != nil {
		return fmt.Errorf("failed to list frames: %w", err)
	}

	log.Printf("Uploading %d frames to S3", len(files))
	for _, file := range files {
		frameFile, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("failed to open frame: %w", err)
		}

		frameName := filepath.Base(file)
		s3Key := filepath.Join(message.ImageDestination, fmt.Sprintf("chunk_%d", message.ChunkPart), frameName)

		err = uc.storageService.UploadObject(ctx, message.Bucket, s3Key, frameFile, "image/jpeg")
		frameFile.Close()
		
		if err != nil {
			return fmt.Errorf("failed to upload frame: %w", err)
		}
	}

	chunkProcessed := dto.ChunkProcessedDTO{
		VideoID:   message.VideoObjectID,
		User:      message.User,
		ChunkPart: message.ChunkPart,
		Status:    "completed",
	}

	messageBytes, err := json.Marshal(chunkProcessed)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := uc.messageBroker.Publish(ctx, uc.snsTopicArn, string(messageBytes)); err != nil {
		return fmt.Errorf("failed to publish to SNS: %w", err)
	}

	log.Printf("Successfully processed chunk %d with %d frames", message.ChunkPart, len(files))
	return nil
}
