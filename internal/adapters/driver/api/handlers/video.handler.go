package handlers

import (
	"net/http"
	"video_processor_service/internal/common/config/env"
	"video_processor_service/internal/core/domain/ports"
	"video_processor_service/internal/core/dto"
	"video_processor_service/internal/core/factory"
	"video_processor_service/internal/core/use_cases"

	"github.com/gin-gonic/gin"
)

type VideoHandler struct {
	videoProcessor ports.IVideoProcessor
	queuePublisher ports.IQueuePublisher
	storageService ports.IStorageService
}

func NewVideoHandler() *VideoHandler {
	return &VideoHandler{
		videoProcessor: factory.NewVideoProcessor(),
		queuePublisher: factory.NewQueuePublisher(),
		storageService: factory.NewStorageService(),
	}
}

// @Summary Process Video
// @Tags Videos
// @Accept json
// @Produce json
// @Param request body dto.ProcessVideoRequestDTO true "Video processing request"
// @Success 202 {object} dto.ProcessVideoResponseDTO
// @Failure 400
// @Failure 500
// @Router /videos/process [post]
func (h *VideoHandler) ProcessVideo(ctx *gin.Context) {
	var request dto.ProcessVideoRequestDTO

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cfg := env.GetConfig()

	useCase := use_cases.NewOrchestrateVideoProcessingUseCase(
		h.videoProcessor,
		h.queuePublisher,
		cfg.AWS.SQS.Queues.FrameExtraction,
	)

	response, err := useCase.Execute(ctx, request)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusAccepted, response)
}
