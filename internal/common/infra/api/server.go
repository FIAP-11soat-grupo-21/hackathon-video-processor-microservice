package api

import (
	"log"

	"github.com/gin-gonic/gin"

	"video_processor_service/internal/adapters/driver/api/routes"
	"video_processor_service/internal/common/config/env"
)

func Init() {
	config := env.GetConfig()

	if config.IsProduction() {
		log.Printf("Running in production mode on [%s]", config.API.URL)
		gin.SetMode(gin.ReleaseMode)
	}

	ginRouter := gin.Default()

	ginRouter.Use(gin.Logger())
	ginRouter.Use(gin.Recovery())

	v1Routes := ginRouter.Group("/v1")

	routes.RegisterVideoRoutes(v1Routes.Group("/videos"))

	ginRouter.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"status": "healthy"})
	})

	if err := ginRouter.Run(config.API.URL); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
