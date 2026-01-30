package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

type Service struct {
	provider          StorageProvider
	repo              FileRepository
	presignedDuration time.Duration
	publicBaseUrl     string
	logger            logger.Logger
}

type PresignedUpload struct {
	SignedUrl string `json:"uploadSignedUrl"`
	FinalUrl  string `json:"finalUrl"`
}

func NewService(
	provider StorageProvider,
	repo FileRepository,
	presignedDurationMinutes int,
	publicBaseUrl string,
	logger logger.Logger,
) *Service {
	if presignedDurationMinutes <= 0 {
		presignedDurationMinutes = 60
	}

	return &Service{
		provider:          provider,
		repo:              repo,
		presignedDuration: time.Duration(presignedDurationMinutes) * time.Minute,
		publicBaseUrl:     publicBaseUrl,
		logger:            logger,
	}
}

func (s *Service) Upload(ctx context.Context, userID int64, fileType string, reader io.Reader, fileName, contentType string, size int64) (*FileReference, error) {
	extension := filepath.Ext(fileName)

	key := fmt.Sprintf("users/%d/%s/%s%s", userID, strings.ToLower(fileType), uuid.New().String(), extension)

	s.logger.Debug("Uploading file", zap.String("key", key), zap.Int64("size", size))

	_, err := s.provider.Upload(ctx, key, reader, contentType, size)
	if err != nil {
		s.logger.Error("Failed to upload file to provider", zap.Error(err))
		return nil, err
	}

	fileRef := &FileReference{
		UserID:           userID,
		OriginalFilename: fileName,
		StorageKey:       key,
		ContentType:      contentType,
		FileSize:         size,
		FileType:         fileType,
		StorageProvider:  s.provider.GetProviderName(),
	}

	if err := s.repo.Save(ctx, fileRef); err != nil {
		s.logger.Error("Failed to save file reference to DB", zap.Error(err))

		return nil, err
	}

	return fileRef, nil
}

func (s *Service) GetPresignedUploadUrl(ctx context.Context, fileName, contentType string, contentLength int64) (*PresignedUpload, error) {
	extension := filepath.Ext(fileName)
	key := fmt.Sprintf("users/avatars/%s%s", uuid.New().String(), extension)

	s.logger.Debug("Generating upload URL", zap.String("key", key), zap.String("contentType", contentType))

	signedUrl, err := s.provider.GeneratePresignedUploadUrl(ctx, key, contentType, contentLength, s.presignedDuration)
	if err != nil {
		s.logger.Error("Failed to generate presigned upload URL", zap.Error(err))
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
