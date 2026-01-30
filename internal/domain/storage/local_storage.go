package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

type LocalStorageProvider struct {
	basePath string
	logger   logger.Logger
}

func NewLocalStorageProvider(cfg LocalConfig, logger logger.Logger) (*LocalStorageProvider, error) {
	if cfg.BasePath == "" {
		cfg.BasePath = "./uploads"
	}

	absPath, err := filepath.Abs(cfg.BasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for %s: %w", cfg.BasePath, err)
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", absPath, err)
	}

	logger.Info("Local storage initialized", zap.String("path", absPath))

	return &LocalStorageProvider{
		basePath: absPath,
		logger:   logger,
	}, nil
}

func (l *LocalStorageProvider) Upload(ctx context.Context, key string, reader io.Reader, contentType string, size int64) (string, error) {
	l.logger.Debug("Saving file locally", zap.String("key", key))

	filePath := filepath.Join(l.basePath, key)
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return "", err
	}

	return "file://" + filePath, nil
}

func (l *LocalStorageProvider) GetPresignedUrl(ctx context.Context, key string, duration time.Duration) (string, error) {
	filePath := filepath.Join(l.basePath, key)
	return "file://" + filePath, nil
}

func (l *LocalStorageProvider) GeneratePresignedUploadUrl(ctx context.Context, key string, contentType string, contentLength int64, duration time.Duration) (string, error) {
	filePath := filepath.Join(l.basePath, key)
	return "file://" + filePath, nil
}

func (l *LocalStorageProvider) Delete(ctx context.Context, key string) error {
	filePath := filepath.Join(l.basePath, key)
	return os.Remove(filePath)
}

func (l *LocalStorageProvider) Exists(ctx context.Context, key string) (bool, error) {
	filePath := filepath.Join(l.basePath, key)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func (l *LocalStorageProvider) GetProviderName() string {
	return "LOCAL"
}

func (l *LocalStorageProvider) GetPublicUrl(key string) string {
	return "/api/files/" + key
}
