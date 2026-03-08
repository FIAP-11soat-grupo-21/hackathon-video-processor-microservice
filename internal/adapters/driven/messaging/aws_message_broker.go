package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"

	"video_processor_service/internal/core/domain/ports"
)

type AWSMessageBroker struct {
	snsClient *sns.Client
	sqsClient *sqs.Client
	ctx       context.Context
	cancelFn  context.CancelFunc
	isRunning bool
}

type SNSNotification struct {
	Type      string `json:"Type"`
	MessageId string `json:"MessageId"`
	Message   string `json:"Message"`
}

func NewAWSMessageBroker() (*AWSMessageBroker, error) {
	ctx := context.Background()

	optFns := []func(*config.LoadOptions) error{}

	if region := os.Getenv("AWS_REGION"); region != "" {
		optFns = append(optFns, config.WithRegion(region))
	}

	if endpoint := os.Getenv("AWS_ENDPOINT"); endpoint != "" {
		optFns = append(optFns, config.WithBaseEndpoint(endpoint))
	}

	cfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	consumerCtx, cancel := context.WithCancel(ctx)

	return &AWSMessageBroker{
		snsClient: sns.NewFromConfig(cfg),
		sqsClient: sqs.NewFromConfig(cfg),
		ctx:       consumerCtx,
		cancelFn:  cancel,
		isRunning: false,
	}, nil
}

func (b *AWSMessageBroker) Publish(ctx context.Context, destination string, message string) error {
	_, err := b.snsClient.Publish(ctx, &sns.PublishInput{
		TopicArn: aws.String(destination),
		Message:  aws.String(message),
	})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (b *AWSMessageBroker) Subscribe(destination string, handler ports.MessageHandler) error {
	if b.isRunning {
		log.Printf("Consumer already running for queue: %s", destination)
		return nil
	}

	b.isRunning = true
	log.Printf(" [*] Starting to consume messages from SQS queue: %s", destination)

	go b.pollMessages(destination, handler)

	return nil
}

func (b *AWSMessageBroker) pollMessages(queueURL string, handler ports.MessageHandler) {
	for {
		select {
		case <-b.ctx.Done():
			log.Printf("Stopping consumer for queue: %s", queueURL)
			b.isRunning = false
			return
		default:
			result, err := b.sqsClient.ReceiveMessage(b.ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(queueURL),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
				VisibilityTimeout:   30,
			})

			if err != nil {
				log.Printf("Error receiving messages from SQS: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			for _, message := range result.Messages {
				err = b.processMessage(queueURL, message, handler)
				if err != nil {
					log.Printf("Error processing message ID %s: %v", *message.MessageId, err)
				}
			}
		}
	}
}

func (b *AWSMessageBroker) processMessage(queueURL string, message types.Message, handler ports.MessageHandler) error {
	if message.Body == nil {
		log.Printf("Received message with nil body")
		return nil
	}

	messageBody, err := b.unmarshallMessage(message)
	if err != nil {
		log.Printf("Error unmarshaling message ID %s: %v", *message.MessageId, err)
		return err
	}

	err = handler(messageBody)
	if err != nil {
		return err
	}

	_, err = b.sqsClient.DeleteMessage(b.ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: message.ReceiptHandle,
	})

	if err != nil {
		return err
	}

	return nil
}

func (b *AWSMessageBroker) unmarshallMessage(message types.Message) ([]byte, error) {
	var snsNotification SNSNotification
	if err := json.Unmarshal([]byte(*message.Body), &snsNotification); err == nil && snsNotification.Type != "" {
		return []byte(snsNotification.Message), nil
	}

	return []byte(*message.Body), nil
}
