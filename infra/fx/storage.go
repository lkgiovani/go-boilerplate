package fx

import (
	"github.com/lkgiovani/go-boilerplate/infra/config"
	"github.com/lkgiovani/go-boilerplate/internal/delivery"
	"github.com/lkgiovani/go-boilerplate/internal/domain/storage"
	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var StorageModule = fx.Module("storage",
	fx.Provide(
		storage.NewFileRepository,
		NewStorageProvider,
		provideStorageService,
		delivery.NewUploadHandler,
	),
)

func provideStorageService(provider storage.StorageProvider, repo storage.FileRepository, cfg *config.Config, log logger.Logger) *storage.Service {
	return storage.NewService(
		provider,
		repo,
		cfg.Storage.PresignedUrlDuration,
		cfg.Storage.PublicBaseURL,
		log,
	)
}

func NewStorageProvider(cfg *config.Config, log logger.Logger) (storage.StorageProvider, error) {
	var storageCfg any

	switch storage.ProviderType(cfg.Storage.Provider) {
	case storage.ProviderS3:
		storageCfg = storage.S3Config{
			AccessKeyID:     cfg.Storage.S3AccessKey,
			SecretAccessKey: cfg.Storage.S3SecretKey,
			Region:          cfg.Storage.S3Region,
			BucketName:      cfg.Storage.S3BucketName,
			Endpoint:        cfg.Storage.S3Endpoint,
		}
	case storage.ProviderR2:
		storageCfg = storage.R2Config{
			AccountID:       cfg.Storage.R2AccountID,
			AccessKeyID:     cfg.Storage.R2AccessKey,
			SecretAccessKey: cfg.Storage.R2SecretKey,
			BucketName:      cfg.Storage.R2BucketName,
			PublicURL:       cfg.Storage.R2PublicURL,
		}
	case storage.ProviderLocal, "":
		storageCfg = storage.LocalConfig{
			BasePath: cfg.Storage.LocalDir,
		}
	}

	provider, err := storage.NewStorageProvider(cfg.Storage.Provider, storageCfg, log)
	if err != nil {
		log.Error("Failed to create storage provider", zap.Error(err))
		return nil, err
	}

	return provider, nil
}
