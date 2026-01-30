package storage

import (
	"context"
	"io"
	"time"
)

// StorageProvider is the interface for file storage operations
type StorageProvider interface {
	Upload(ctx context.Context, key string, reader io.Reader, contentType string, size int64) (string, error)
	GetPresignedUrl(ctx context.Context, key string, duration time.Duration) (string, error)
	GeneratePresignedUploadUrl(ctx context.Context, key string, contentType string, contentLength int64, duration time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	GetProviderName() string
	GetPublicUrl(key string) string
}

// S3Config holds configuration for AWS S3
type S3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	BucketName      string
	Endpoint        string // Optional (e.g., for LocalStack)
}

// LocalConfig holds configuration for Local storage
type LocalConfig struct {
	BasePath string
}

// R2Config holds configuration for Cloudflare R2
type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	PublicURL       string // Optional: Custom domain/public URL for R2
}

// StorageConfig is a generic config for storage providers
type StorageProviderConfig interface {
	S3Config | LocalConfig | R2Config
}

type Config[T StorageProviderConfig] struct {
	Provider string
	Config   T
}
