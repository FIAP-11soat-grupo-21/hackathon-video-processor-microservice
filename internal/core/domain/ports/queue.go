package ports

import "context"

type MessageHandler func(msgBody []byte) error

type IQueueConsumer interface {
	ConsumeQueue(queueURL string, handler MessageHandler) error
}

type IQueuePublisher interface {
	SendMessage(ctx context.Context, queueURL string, messageBody string) error
}
