package storage

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	provider          StorageProvider
	presignedDuration time.Duration
	publicBaseUrl     string
	logger            *slog.Logger
}

type PresignedUpload struct {
	SignedUrl string `json:"uploadSignedUrl"`
	FinalUrl  string `json:"finalUrl"`
}

func NewService(
	provider StorageProvider,
	presignedDurationMinutes int,
	publicBaseUrl string,
	logger *slog.Logger,
) *Service {
	if presignedDurationMinutes <= 0 {
		presignedDurationMinutes = 60
	}

	return &Service{
		provider:          provider,
		presignedDuration: time.Duration(presignedDurationMinutes) * time.Minute,
		publicBaseUrl:     publicBaseUrl,
		logger:            logger,
	}
}

func (s *Service) GetPresignedUploadUrl(ctx context.Context, fileName, contentType string, contentLength int64) (*PresignedUpload, error) {
	extension := filepath.Ext(fileName)
	key := fmt.Sprintf("users/avatars/%s%s", uuid.New().String(), extension)

	s.logger.Debug("Generating upload URL", slog.String("key", key), slog.String("contentType", contentType))

	signedUrl, err := s.provider.GeneratePresignedUploadUrl(ctx, key, contentType, contentLength, s.presignedDuration)
	if err != nil {
		s.logger.Error("Failed to generate presigned upload URL", slog.Any("error", err))
		return nil, err
	}

	finalUrl := s.provider.GetPublicUrl(key)
	if s.publicBaseUrl != "" {
		finalUrl = strings.TrimSuffix(s.publicBaseUrl, "/") + "/" + key
	}

	return &PresignedUpload{
		SignedUrl: signedUrl,
		FinalUrl:  finalUrl,
	}, nil
}
