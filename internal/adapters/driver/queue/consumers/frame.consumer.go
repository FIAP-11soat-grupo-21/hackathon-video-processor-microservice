package consumers

import (
	"log"
	"video_processor_service/internal/adapters/driver/queue/handlers"
	"video_processor_service/internal/common/config/env"
	"video_processor_service/internal/core/domain/ports"
	"video_processor_service/internal/core/factory"
)

type FrameConsumer struct {
	consumer ports.IQueueConsumer
}

func NewFrameConsumer(consumer ports.IQueueConsumer) *FrameConsumer {
	return &FrameConsumer{
		consumer: consumer,
	}
}

func (fc *FrameConsumer) RegisterConsumers() {
	cfg := env.GetConfig()

	err := fc.consumer.ConsumeQueue(cfg.AWS.SQS.Queues.FrameExtraction, handlers.ExtractFrame)

	if err != nil {
		log.Fatalf("Failed to register frame extraction consumer: %v", err)
	}

	log.Println("Frame extraction consumer registered successfully")
}

func RegisterConsumers() {
	consumer := NewFrameConsumer(factory.NewQueueConsumer())
	consumer.RegisterConsumers()
}
