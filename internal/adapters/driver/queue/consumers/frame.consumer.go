package consumers

import (
	"log"
	"video_processor_service/internal/adapters/driver/queue/handlers"
	"video_processor_service/internal/common/config/env"
	"video_processor_service/internal/core/domain/ports"
	"video_processor_service/internal/core/factory"
)

type FrameConsumer struct {
	messageBroker ports.IMessageBroker
}

func NewFrameConsumer(messageBroker ports.IMessageBroker) *FrameConsumer {
	return &FrameConsumer{
		messageBroker: messageBroker,
	}
}

func (fc *FrameConsumer) RegisterConsumers() {
	cfg := env.GetConfig()

	err := fc.messageBroker.Subscribe(cfg.AWS.SQS.Queues.FrameExtraction, handlers.ExtractFrame)

	if err != nil {
		log.Fatalf("Failed to register frame extraction consumer: %v", err)
	}

	log.Println("Frame extraction consumer registered successfully")
}

func RegisterConsumers() {
	broker, err := factory.NewMessageBroker()
	if err != nil {
		log.Fatalf("Failed to create message broker: %v", err)
	}
	consumer := NewFrameConsumer(broker)
	consumer.RegisterConsumers()
}
