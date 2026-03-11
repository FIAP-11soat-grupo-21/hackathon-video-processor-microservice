package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"video_processor_service/internal/common/config/env"
	"video_processor_service/internal/core/domain/ports"
)

type s3API interface {
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

type S3StorageService struct {
	client s3API
	cfg    *env.Config
}

func NewS3StorageService() *S3StorageService {
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

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	return &S3StorageService{
		client: client,
		cfg:    appCfg,
	}
}

func (s *S3StorageService) GetObjectMetadata(ctx context.Context, bucket, key string) (int64, error) {
	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return 0, fmt.Errorf("failed to get object metadata: %w", err)
	}

	return *result.ContentLength, nil
}

func (s *S3StorageService) GetObjectRange(ctx context.Context, bucket, key string, start, end int64) (io.ReadCloser, error) {
	rangeHeader := fmt.Sprintf("bytes=%d-%d", start, end)

	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Range:  aws.String(rangeHeader),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object range: %w", err)
	}

	return result.Body, nil
}

func (s *S3StorageService) UploadObject(ctx context.Context, bucket, key string, body io.Reader, contentType string) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	return nil
}

func (s *S3StorageService) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return result.Body, nil
}

func (s *S3StorageService) ListObjects(ctx context.Context, bucket, prefix string) ([]ports.S3Object, error) {
	result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var objects []ports.S3Object
	for _, obj := range result.Contents {
		objects = append(objects, ports.S3Object{
			Key:  *obj.Key,
			Size: *obj.Size,
		})
	}

	return objects, nil
}

func (s *S3StorageService) GeneratePresignedURL(ctx context.Context, bucket, key string, expirationSeconds int) (string, error) {
	appCfg := env.GetConfig()
	
	optFns := []func(*config.LoadOptions) error{
		config.WithRegion(appCfg.AWS.Region),
	}

	if appCfg.AWS.Endpoint != "" {
		optFns = append(optFns, config.WithBaseEndpoint(appCfg.AWS.Endpoint))
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return "", fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	presignClient := s3.NewPresignClient(client)
	
	presignResult, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(time.Duration(expirationSeconds)*time.Second))
	
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignResult.URL, nil
}
