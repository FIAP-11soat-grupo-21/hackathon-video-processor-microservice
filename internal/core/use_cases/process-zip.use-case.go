package use_cases

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"

	"video_processor_service/internal/core/domain/ports"
	"video_processor_service/internal/core/dto"
)

type ProcessZipUseCase struct {
	storageService ports.IStorageService
	messageBroker  ports.IMessageBroker
	snsTopicArn    string
}

func NewProcessZipUseCase(
	storageService ports.IStorageService,
	messageBroker ports.IMessageBroker,
	snsTopicArn string,
) *ProcessZipUseCase {
	return &ProcessZipUseCase{
		storageService: storageService,
		messageBroker:  messageBroker,
		snsTopicArn:    snsTopicArn,
	}
}

func (uc *ProcessZipUseCase) Execute(ctx context.Context, message dto.AllChunksProcessedDTO) error {
	log.Printf("Creating zip for video %s from location %s", message.VideoID, message.ImagesLocation)

	objects, err := uc.storageService.ListObjects(ctx, message.Bucket, message.ImagesLocation)
	if err != nil {
		return fmt.Errorf("failed to list images: %w", err)
	}

	if len(objects) == 0 {
		return fmt.Errorf("no images found at %s", message.ImagesLocation)
	}

	log.Printf("Found %d images to zip", len(objects))

	zipBuffer := new(bytes.Buffer)
	zipWriter := zip.NewWriter(zipBuffer)

	for _, obj := range objects {
		if obj.Size == 0 {
			continue
		}

		log.Printf("Adding %s to zip", obj.Key)

		reader, err := uc.storageService.GetObject(ctx, message.Bucket, obj.Key)
		if err != nil {
			return fmt.Errorf("failed to download image %s: %w", obj.Key, err)
		}

		fileName := filepath.Base(obj.Key)
		writer, err := zipWriter.Create(fileName)
		if err != nil {
			reader.Close()
			return fmt.Errorf("failed to create zip entry: %w", err)
		}

		if _, err := io.Copy(writer, reader); err != nil {
			reader.Close()
			return fmt.Errorf("failed to write to zip: %w", err)
		}
		reader.Close()
	}

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("failed to close zip: %w", err)
	}

	videoIDParts := strings.Split(message.VideoID, "/")
	videoID := message.VideoID
	if len(videoIDParts) >= 2 {
		videoID = videoIDParts[1]
	}
	
	zipKey := fmt.Sprintf("videos/%s/output.zip", videoID)
	log.Printf("Uploading zip to s3://%s/%s", message.Bucket, zipKey)

	if err := uc.storageService.UploadObject(ctx, message.Bucket, zipKey, bytes.NewReader(zipBuffer.Bytes()), "application/zip"); err != nil {
		return fmt.Errorf("failed to upload zip: %w", err)
	}

	presignedURL, err := uc.storageService.GeneratePresignedURL(ctx, message.Bucket, zipKey, 7*24*3600)
	if err != nil {
		return fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	completeMessage := dto.VideoProcessingCompleteDTO{
		VideoID:     message.VideoID,
		User:        message.User,
		DownloadURL: presignedURL,
		Status:      "completed",
	}

	messageBytes, err := json.Marshal(completeMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if err := uc.messageBroker.Publish(ctx, uc.snsTopicArn, string(messageBytes)); err != nil {
		return fmt.Errorf("failed to publish to SNS: %w", err)
	}

	log.Printf("Successfully created zip with %d images for video %s", len(objects), message.VideoID)
	return nil
}
