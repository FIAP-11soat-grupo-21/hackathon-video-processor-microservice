package use_cases

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"video_processor_service/internal/core/dto"
)

type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) GetObjectMetadata(ctx context.Context, bucket, key string) (int64, error) {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorageService) GetObjectRange(ctx context.Context, bucket, key string, start, end int64) (io.ReadCloser, error) {
	args := m.Called(ctx, bucket, key, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockStorageService) UploadObject(ctx context.Context, bucket, key string, body io.Reader, contentType string) error {
	args := m.Called(ctx, bucket, key, body, contentType)
	return args.Error(0)
}

func TestExtractFrameUseCase_Execute_Success(t *testing.T) {
	mockVideoProcessor := new(MockVideoProcessor)
	mockStorageService := new(MockStorageService)
	bucketName := "output-bucket"

	useCase := NewExtractFrameUseCase(
		mockVideoProcessor,
		mockStorageService,
		bucketName,
	)

	ctx := context.Background()
	message := dto.FrameExtractionMessageDTO{
		JobID:     "test-job-123",
		Bucket:    "input-bucket",
		Key:       "video.mp4",
		Timestamp: 5.5,
		Index:     1,
	}

	frameData := []byte("fake-jpeg-data")
	mockVideoProcessor.On("ExtractFrame", ctx, "input-bucket", "video.mp4", 5.5).Return(frameData, nil)
	mockStorageService.On("UploadObject", ctx, bucketName, "video-processor/frames/test-job-123/frame_1_5.50s.jpg", mock.Anything, "image/jpeg").Return(nil)

	err := useCase.Execute(ctx, message)

	assert.NoError(t, err)
	mockVideoProcessor.AssertExpectations(t)
	mockStorageService.AssertExpectations(t)
}

func TestExtractFrameUseCase_Execute_ExtractFrameError(t *testing.T) {
	mockVideoProcessor := new(MockVideoProcessor)
	mockStorageService := new(MockStorageService)
	bucketName := "output-bucket"

	useCase := NewExtractFrameUseCase(
		mockVideoProcessor,
		mockStorageService,
		bucketName,
	)

	ctx := context.Background()
	message := dto.FrameExtractionMessageDTO{
		JobID:     "test-job-123",
		Bucket:    "input-bucket",
		Key:       "video.mp4",
		Timestamp: 5.5,
		Index:     1,
	}

	expectedError := errors.New("ffmpeg error")
	mockVideoProcessor.On("ExtractFrame", ctx, "input-bucket", "video.mp4", 5.5).Return([]byte(nil), expectedError)

	err := useCase.Execute(ctx, message)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to extract frame")
	mockVideoProcessor.AssertExpectations(t)
	mockStorageService.AssertNotCalled(t, "UploadObject")
}

func TestExtractFrameUseCase_Execute_UploadError(t *testing.T) {
	mockVideoProcessor := new(MockVideoProcessor)
	mockStorageService := new(MockStorageService)
	bucketName := "output-bucket"

	useCase := NewExtractFrameUseCase(
		mockVideoProcessor,
		mockStorageService,
		bucketName,
	)

	ctx := context.Background()
	message := dto.FrameExtractionMessageDTO{
		JobID:     "test-job-123",
		Bucket:    "input-bucket",
		Key:       "video.mp4",
		Timestamp: 10.0,
		Index:     2,
	}

	frameData := []byte("fake-jpeg-data")
	uploadError := errors.New("s3 upload failed")
	mockVideoProcessor.On("ExtractFrame", ctx, "input-bucket", "video.mp4", 10.0).Return(frameData, nil)
	mockStorageService.On("UploadObject", ctx, bucketName, "video-processor/frames/test-job-123/frame_2_10.00s.jpg", mock.Anything, "image/jpeg").Return(uploadError)

	err := useCase.Execute(ctx, message)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to upload frame")
	mockVideoProcessor.AssertExpectations(t)
	mockStorageService.AssertExpectations(t)
}
