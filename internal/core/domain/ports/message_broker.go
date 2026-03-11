package ports

import "context"

type MessageHandler func(msgBody []byte) error

type IMessageBroker interface {
	Publish(ctx context.Context, destination string, message string) error
	Subscribe(destination string, handler MessageHandler) error
}
