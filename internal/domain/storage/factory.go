package storage

import (
	"fmt"
	"log/slog"
)

type ProviderType string

const (
	ProviderS3    ProviderType = "s3"
	ProviderLocal ProviderType = "local"
	ProviderR2    ProviderType = "r2"
)

func NewStorageProvider(
	provider string,
	storageConfig any,
	logger *slog.Logger,
) (StorageProvider, error) {
	switch ProviderType(provider) {
	case ProviderS3:
		cfg, ok := storageConfig.(S3Config)
		if !ok {
			return nil, fmt.Errorf("invalid storage config for S3")
		}
		return NewS3StorageProvider(cfg, logger)

	case ProviderR2:
		cfg, ok := storageConfig.(R2Config)
		if !ok {
			return nil, fmt.Errorf("invalid storage config for R2")
		}
		return NewR2StorageProvider(cfg, logger)

	case ProviderLocal, "":
		cfg, ok := storageConfig.(LocalConfig)
		if !ok {
			// fallback to default if ""
			cfg = LocalConfig{BasePath: "./uploads"}
		}
		return NewLocalStorageProvider(cfg, logger)

	default:
		return nil, fmt.Errorf("unknown storage provider: %s", provider)
	}
}
