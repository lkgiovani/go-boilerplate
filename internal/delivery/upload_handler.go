package delivery

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/lkgiovani/go-boilerplate/internal/delivery/dto"
	"github.com/lkgiovani/go-boilerplate/internal/domain/storage"
)

type UploadHandler struct {
	storageService *storage.Service
	logger         *slog.Logger
}

func NewUploadHandler(storageService *storage.Service, logger *slog.Logger) *UploadHandler {
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
		h.logger.Error("Failed to generate presigned URL", slog.Any("error", err))
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate upload URL")
	}

	return c.Status(fiber.StatusOK).JSON(dto.UploadResponseDTO{
		UploadSignedURL: result.SignedUrl,
		PublicURL:       result.FinalUrl,
	})
}
