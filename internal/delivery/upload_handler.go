package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lkgiovani/go-boilerplate/internal/delivery/dto"
	"github.com/lkgiovani/go-boilerplate/internal/domain/storage"
	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

type UploadHandler struct {
	storageService *storage.Service
	logger         logger.Logger
}

func NewUploadHandler(storageService *storage.Service, logger logger.Logger) *UploadHandler {
	return &UploadHandler{
		storageService: storageService,
		logger:         logger,
	}
}

func (h *UploadHandler) GetUploadUrl(c *fiber.Ctx) error {
	var req dto.UploadImageRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// In a real scenario, you would validate the request here
	// using a validator library if integrated.

	result, err := h.storageService.GetPresignedUploadUrl(c.Context(), req.FileName, req.ContentType, req.ContentLength)
	if err != nil {
		h.logger.Error("Failed to generate presigned URL", zap.Error(err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate upload URL")
	}

	return c.Status(fiber.StatusOK).JSON(dto.UploadResponseDTO{
		UploadSignedURL: result.SignedUrl,
		PublicURL:       result.FinalUrl,
	})
}
