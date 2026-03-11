package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"

	lambdaHandler "video_processor_service/internal/adapters/driver/lambda"
)

func main() {
	handler, err := lambdaHandler.NewUpdateVideoChunkStatusHandler()
	if err != nil {
		log.Fatalf("Failed to create handler: %v", err)
	}

	lambda.Start(handler.Handle)
}
