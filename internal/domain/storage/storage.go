package storage

import (
	"context"
	"io"
	"time"
)

type FileReference struct {
	ID               int64     `gorm:"primaryKey;autoIncrement"`
	UserID           int64     `gorm:"not null"`
	OriginalFilename string    `gorm:"not null"`
	StorageKey       string    `gorm:"uniqueIndex;not null"`
	ContentType      string    `gorm:"not null"`
	FileSize         int64     `gorm:"not null"`
	FileType         string    `gorm:"not null"`
	StorageProvider  string    `gorm:"not null;default:'S3'"`
	CreatedAt        time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type FileRepository interface {
	Save(ctx context.Context, file *FileReference) error
	GetByID(ctx context.Context, id int64) (*FileReference, error)
	GetByStorageKey(ctx context.Context, key string) (*FileReference, error)
	Delete(ctx context.Context, id int64) error
	FindByUserID(ctx context.Context, userID int64) ([]FileReference, error)
}

type StorageProvider interface {
	Upload(ctx context.Context, key string, reader io.Reader, contentType string, size int64) (string, error)
	GetPresignedUrl(ctx context.Context, key string, duration time.Duration) (string, error)
	GeneratePresignedUploadUrl(ctx context.Context, key string, contentType string, contentLength int64, duration time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	GetProviderName() string
	GetPublicUrl(key string) string
}

type S3Config struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	BucketName      string
	Endpoint        string
}

type LocalConfig struct {
	BasePath string
}

type R2Config struct {
	AccountID       string
	AccessKeyID     string
	SecretAccessKey string
	BucketName      string
	PublicURL       string
}

type StorageProviderConfig interface {
	S3Config | LocalConfig | R2Config
}

type Config[T StorageProviderConfig] struct {
	Provider string
	Config   T
}
