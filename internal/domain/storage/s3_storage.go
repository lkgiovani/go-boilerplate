package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

type S3StorageProvider struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	uploader      *manager.Uploader
	bucket        string
	region        string
	logger        logger.Logger
}

func NewS3StorageProvider(cfg S3Config, logger logger.Logger) (*S3StorageProvider, error) {
	if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.BucketName == "" {
		return nil, fmt.Errorf("S3 credentials and bucket name are required")
	}

	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}

	creds := credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")

	opts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(creds),
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
	})

	return &S3StorageProvider{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		uploader:      manager.NewUploader(client),
		bucket:        cfg.BucketName,
		region:        cfg.Region,
		logger:        logger,
	}, nil
}

func (s *S3StorageProvider) Upload(ctx context.Context, key string, reader io.Reader, contentType string, size int64) (string, error) {
	s.logger.Debug("Uploading file to S3", zap.String("key", key), zap.String("bucket", s.bucket))

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	}

	_, err := s.uploader.Upload(ctx, input)
	if err != nil {
		s.logger.Error("Failed to upload file to S3", zap.String("key", key), zap.Error(err))
		return "", err
	}

	return fmt.Sprintf("s3://%s/%s", s.bucket, key), nil
}

func (s *S3StorageProvider) GetPresignedUrl(ctx context.Context, key string, duration time.Duration) (string, error) {
	s.logger.Debug("Generating presigned GET URL", zap.String("key", key))

	request, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(duration))

	if err != nil {
		return "", err
	}

	return request.URL, nil
}

func (s *S3StorageProvider) GeneratePresignedUploadUrl(ctx context.Context, key string, contentType string, contentLength int64, duration time.Duration) (string, error) {
	s.logger.Debug("Generating presigned PUT URL", zap.String("key", key), zap.String("contentType", contentType))

	request, err := s.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(contentLength),
	}, s3.WithPresignExpires(duration))

	if err != nil {
		return "", err
	}

	return request.URL, nil
}

func (s *S3StorageProvider) Delete(ctx context.Context, key string) error {
	s.logger.Debug("Deleting file from S3", zap.String("key", key))

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	return err
}

func (s *S3StorageProvider) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return false, nil // Should check for specific error types if needed
	}

	return true, nil
}

func (s *S3StorageProvider) GetProviderName() string {
	return "AWS_S3"
}

func (s *S3StorageProvider) GetPublicUrl(key string) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucket, key)
}
