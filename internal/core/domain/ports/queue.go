package ports

type MessageHandler func(msgBody []byte) error

type IQueueConsumer interface {
	ConsumeQueue(queueURL string, handler MessageHandler) error
}
