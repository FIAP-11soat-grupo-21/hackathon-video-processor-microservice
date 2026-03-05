package use_cases

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"video_processor_service/internal/core/dto"
)

type MockVideoProcessor struct {
	mock.Mock
}

func (m *MockVideoProcessor) GetVideoDuration(ctx context.Context, bucket, key string) (float64, error) {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockVideoProcessor) ExtractFrame(ctx context.Context, bucket, key string, timestamp float64) ([]byte, error) {
	args := m.Called(ctx, bucket, key, timestamp)
	return args.Get(0).([]byte), args.Error(1)
}

type MockQueuePublisher struct {
	mock.Mock
}

func (m *MockQueuePublisher) SendMessage(ctx context.Context, queueURL string, messageBody string) error {
	args := m.Called(ctx, queueURL, messageBody)
	return args.Error(0)
}

func TestOrchestrateVideoProcessingUseCase_Execute_Success(t *testing.T) {
	mockVideoProcessor := new(MockVideoProcessor)
	mockQueuePublisher := new(MockQueuePublisher)
	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789/test-queue"

	useCase := NewOrchestrateVideoProcessingUseCase(
		mockVideoProcessor,
		mockQueuePublisher,
		queueURL,
	)

	ctx := context.Background()
	request := dto.ProcessVideoRequestDTO{
		Bucket:               "test-bucket",
		Key:                  "test-video.mp4",
		ChunkIntervalSeconds: 5.0,
	}

	mockVideoProcessor.On("GetVideoDuration", ctx, "test-bucket", "test-video.mp4").Return(15.0, nil)
	mockQueuePublisher.On("SendMessage", ctx, queueURL, mock.AnythingOfType("string")).Return(nil).Times(3)

	response, err := useCase.Execute(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.JobID)
	assert.Equal(t, "processing", response.Status)
	assert.Equal(t, 3, response.EstimatedFrames)
	assert.Equal(t, 15.0, response.VideoDuration)

	mockVideoProcessor.AssertExpectations(t)
	mockQueuePublisher.AssertExpectations(t)
}

func TestOrchestrateVideoProcessingUseCase_Execute_GetDurationError(t *testing.T) {
	mockVideoProcessor := new(MockVideoProcessor)
	mockQueuePublisher := new(MockQueuePublisher)
	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789/test-queue"

	useCase := NewOrchestrateVideoProcessingUseCase(
		mockVideoProcessor,
		mockQueuePublisher,
		queueURL,
	)

	ctx := context.Background()
	request := dto.ProcessVideoRequestDTO{
		Bucket:               "test-bucket",
		Key:                  "invalid-video.mp4",
		ChunkIntervalSeconds: 5.0,
	}

	expectedError := errors.New("video not found")
	mockVideoProcessor.On("GetVideoDuration", ctx, "test-bucket", "invalid-video.mp4").Return(0.0, expectedError)

	response, err := useCase.Execute(ctx, request)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "failed to get video duration")

	mockVideoProcessor.AssertExpectations(t)
	mockQueuePublisher.AssertNotCalled(t, "SendMessage")
}

func TestOrchestrateVideoProcessingUseCase_CalculateTimestamps(t *testing.T) {
	tests := []struct {
		name     string
		duration float64
		interval float64
		expected []float64
	}{
		{
			name:     "Exact division",
			duration: 10.0,
			interval: 5.0,
			expected: []float64{0.0, 5.0},
		},
		{
			name:     "With remainder",
			duration: 12.0,
			interval: 5.0,
			expected: []float64{0.0, 5.0, 10.0},
		},
		{
			name:     "Single frame",
			duration: 3.0,
			interval: 5.0,
			expected: []float64{0.0},
		},
		{
			name:     "Multiple frames",
			duration: 30.0,
			interval: 10.0,
			expected: []float64{0.0, 10.0, 20.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			useCase := &OrchestrateVideoProcessingUseCase{}

			result := useCase.calculateTimestamps(tt.duration, tt.interval)

			assert.Equal(t, tt.expected, result)
		})
	}
}
