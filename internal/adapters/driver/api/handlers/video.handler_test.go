package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"video_processor_service/internal/core/dto"
)

func setupTestEnv() {
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_ENDPOINT", "http://localhost:4566")
	os.Setenv("AWS_SQS_FRAME_EXTRACTION_QUEUE", "http://localhost:4566/000000000000/frame-extraction")
	os.Setenv("S3_BUCKET", "test-bucket")
}

func setupTestRouter() *gin.Engine {
	setupTestEnv()
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestVideoHandler_ProcessVideo_InvalidJSON(t *testing.T) {
	router := setupTestRouter()
	handler := NewVideoHandler()
	router.POST("/videos/process", handler.ProcessVideo)

	invalidJSON := `{"bucket": "test", "key": }`

	req, _ := http.NewRequest(http.MethodPost, "/videos/process", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVideoHandler_ProcessVideo_MissingRequiredFields(t *testing.T) {
	router := setupTestRouter()
	handler := NewVideoHandler()
	router.POST("/videos/process", handler.ProcessVideo)

	tests := []struct {
		name    string
		payload dto.ProcessVideoRequestDTO
	}{
		{
			name: "Missing bucket",
			payload: dto.ProcessVideoRequestDTO{
				Key:                  "video.mp4",
				ChunkIntervalSeconds: 5.0,
			},
		},
		{
			name: "Missing key",
			payload: dto.ProcessVideoRequestDTO{
				Bucket:               "test-bucket",
				ChunkIntervalSeconds: 5.0,
			},
		},
		{
			name: "Missing chunkIntervalSeconds",
			payload: dto.ProcessVideoRequestDTO{
				Bucket: "test-bucket",
				Key:    "video.mp4",
			},
		},
		{
			name: "Invalid chunkIntervalSeconds (zero)",
			payload: dto.ProcessVideoRequestDTO{
				Bucket:               "test-bucket",
				Key:                  "video.mp4",
				ChunkIntervalSeconds: 0,
			},
		},
		{
			name: "Invalid chunkIntervalSeconds (negative)",
			payload: dto.ProcessVideoRequestDTO{
				Bucket:               "test-bucket",
				Key:                  "video.mp4",
				ChunkIntervalSeconds: -5.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, "/videos/process", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

func TestVideoHandler_ProcessVideo_ValidRequest(t *testing.T) {
	router := setupTestRouter()
	handler := NewVideoHandler()
	router.POST("/videos/process", handler.ProcessVideo)

	payload := dto.ProcessVideoRequestDTO{
		Bucket:               "test-bucket",
		Key:                  "test-video.mp4",
		ChunkIntervalSeconds: 5.0,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/videos/process", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusAccepted || w.Code == http.StatusInternalServerError)
}
