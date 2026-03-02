package queue

import (
	"context"
	"log"
	"video_processor_service/internal/common/config/env"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type sqsPublisherAPI interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type SQSPublisher struct {
	client sqsPublisherAPI
	cfg    *env.Config
}

func NewSQSPublisher() *SQSPublisher {
	ctx := context.Background()
	appCfg := env.GetConfig()

	optFns := []func(*config.LoadOptions) error{
		config.WithRegion(appCfg.AWS.Region),
	}

	if appCfg.AWS.Endpoint != "" {
		optFns = append(optFns, config.WithBaseEndpoint(appCfg.AWS.Endpoint))
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	client := sqs.NewFromConfig(awsCfg)

	return &SQSPublisher{
		client: client,
		cfg:    appCfg,
	}
}

func (p *SQSPublisher) SendMessage(ctx context.Context, queueURL string, messageBody string) error {
	_, err := p.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(queueURL),
		MessageBody: aws.String(messageBody),
	})

	if err != nil {
		log.Printf("Error sending message to SQS: %v", err)
		return err
	}

	return nil
}
