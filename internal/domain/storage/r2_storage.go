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

type R2StorageProvider struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	uploader      *manager.Uploader
	bucket        string
	publicUrl     string
	logger        logger.Logger
}

func NewR2StorageProvider(cfg R2Config, logger logger.Logger) (*R2StorageProvider, error) {
	if cfg.AccountID == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.BucketName == "" {
		return nil, fmt.Errorf("R2 credentials, account ID and bucket name are required")
	}

	// Cloudflare R2 endpoint format: https://<account_id>.r2.cloudflarestorage.com
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)

	creds := credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")

	// R2 uses 'auto' as region, though the SDK might require a placeholder or specific handling
	opts := []func(*config.LoadOptions) error{
		config.WithRegion("auto"),
		config.WithCredentialsProvider(creds),
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load R2 config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &R2StorageProvider{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		uploader:      manager.NewUploader(client),
		bucket:        cfg.BucketName,
		publicUrl:     cfg.PublicURL,
		logger:        logger,
	}, nil
}

func (r *R2StorageProvider) Upload(ctx context.Context, key string, reader io.Reader, contentType string, size int64) (string, error) {
	r.logger.Debug("Uploading file to R2", zap.String("key", key), zap.String("bucket", r.bucket))

	input := &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	}

	_, err := r.uploader.Upload(ctx, input)
	if err != nil {
		r.logger.Error("Failed to upload file to R2", zap.String("key", key), zap.Error(err))
		return "", err
	}

	return fmt.Sprintf("r2://%s/%s", r.bucket, key), nil
}

func (r *R2StorageProvider) GetPresignedUrl(ctx context.Context, key string, duration time.Duration) (string, error) {
	r.logger.Debug("Generating presigned GET URL for R2", zap.String("key", key))

	request, err := r.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(duration))

	if err != nil {
		return "", err
	}

	return request.URL, nil
}

func (r *R2StorageProvider) GeneratePresignedUploadUrl(ctx context.Context, key string, contentType string, contentLength int64, duration time.Duration) (string, error) {
	r.logger.Debug("Generating presigned PUT URL for R2", zap.String("key", key), zap.String("contentType", contentType))

	request, err := r.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(r.bucket),
		Key:           aws.String(key),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(contentLength),
	}, s3.WithPresignExpires(duration))

	if err != nil {
		return "", err
	}

	return request.URL, nil
}

func (r *R2StorageProvider) Delete(ctx context.Context, key string) error {
	r.logger.Debug("Deleting file from R2", zap.String("key", key))

	_, err := r.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})

	return err
}

func (r *R2StorageProvider) Exists(ctx context.Context, key string) (bool, error) {
	_, err := r.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return false, nil
	}

	return true, nil
}

func (r *R2StorageProvider) GetProviderName() string {
	return "CLOUDFLARE_R2"
}

func (r *R2StorageProvider) GetPublicUrl(key string) string {
	if r.publicUrl != "" {
		return fmt.Sprintf("%s/%s", r.publicUrl, key)
	}
	// Fallback or custom logic if needed
	return fmt.Sprintf("r2://%s/%s", r.bucket, key)
}
