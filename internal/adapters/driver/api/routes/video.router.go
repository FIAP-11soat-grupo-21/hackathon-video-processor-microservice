package routes

import (
	"video_processor_service/internal/adapters/driver/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterVideoRoutes(router *gin.RouterGroup) {
	videoHandler := handlers.NewVideoHandler()

	router.POST("/process", videoHandler.ProcessVideo)
}
