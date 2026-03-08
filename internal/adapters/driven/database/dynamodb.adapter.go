package database

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"video_processor_service/internal/core/domain"
)

type VideoChunkModel struct {
	VideoID   string `dynamodbav:"video_id"`
	ChunkPart int    `dynamodbav:"chunk_part"`
	Status    string `dynamodbav:"status"`
	UserID    string `dynamodbav:"user_id"`
	UserName  string `dynamodbav:"user_name"`
	UserEmail string `dynamodbav:"user_email"`
	UpdatedAt string `dynamodbav:"updated_at"`
}

type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBRepository(tableName string) (*DynamoDBRepository, error) {
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

	return &DynamoDBRepository{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: tableName,
	}, nil
}

func (r *DynamoDBRepository) SaveChunk(ctx context.Context, chunk domain.VideoChunk) error {
	model := VideoChunkModel{
		VideoID:   chunk.VideoID,
		ChunkPart: chunk.ChunkPart,
		Status:    chunk.Status,
		UserID:    chunk.UserID,
		UserName:  chunk.UserName,
		UserEmail: chunk.UserEmail,
		UpdatedAt: chunk.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	item, err := attributevalue.MarshalMap(model)
	if err != nil {
		return fmt.Errorf("failed to marshal chunk: %w", err)
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("failed to save chunk: %w", err)
	}

	return nil
}

func (r *DynamoDBRepository) GetChunksByVideoID(ctx context.Context, videoID string) ([]domain.VideoChunk, error) {
	result, err := r.client.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		KeyConditionExpression: aws.String("video_id = :video_id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":video_id": &types.AttributeValueMemberS{Value: videoID},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query chunks: %w", err)
	}

	var chunks []domain.VideoChunk
	for _, item := range result.Items {
		var model VideoChunkModel
		if err := attributevalue.UnmarshalMap(item, &model); err != nil {
			return nil, fmt.Errorf("failed to unmarshal chunk: %w", err)
		}

		chunk := domain.VideoChunk{
			VideoID:   model.VideoID,
			ChunkPart: model.ChunkPart,
			Status:    model.Status,
			UserID:    model.UserID,
			UserName:  model.UserName,
			UserEmail: model.UserEmail,
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}
