package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lkgiovani/go-boilerplate/internal/delivery/dto"
	"github.com/lkgiovani/go-boilerplate/internal/domain/auth"
	"github.com/lkgiovani/go-boilerplate/internal/errors"
)

type MobileAuthHandler struct {
	AuthService  *auth.Service
	ErrorHandler func(c *fiber.Ctx, err error) error
}

func NewMobileAuthHandler(authService *auth.Service, errorHandler func(c *fiber.Ctx, err error) error) *MobileAuthHandler {
	return &MobileAuthHandler{
		AuthService:  authService,
		ErrorHandler: errorHandler,
	}
}

// AuthenticateWithGoogleMobile handles Google authentication for mobile devices
// POST /v1/auth/mobile/oauth2/google
func (h *MobileAuthHandler) AuthenticateWithGoogleMobile(c *fiber.Ctx) error {
	var req dto.MobileOAuth2RequestDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	deviceID := req.DeviceId
	if deviceID == "" {
		deviceID = extractDeviceID(c)
	}
	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()

	ctx := c.UserContext()
	result, err := h.AuthService.AuthenticateWithGoogleMobile(ctx, req.IdToken, deviceID, userAgent, ipAddress)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	response := dto.MobileLoginResponseDTO{
		UserID:       result.UserID,
		Email:        result.Email,
		Name:         result.Name,
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		IsNewUser:    result.IsNewUser,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// RefreshMobileToken handles token refresh for mobile devices
// POST /v1/auth/mobile/refresh
func (h *MobileAuthHandler) RefreshMobileToken(c *fiber.Ctx) error {
	var req dto.MobileRefreshRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()
	deviceID := extractDeviceID(c)

	ctx := c.UserContext()
	result, err := h.AuthService.RefreshMobileToken(ctx, req.RefreshToken, userAgent, ipAddress, deviceID)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	response := dto.MobileRefreshResponseDTO{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
