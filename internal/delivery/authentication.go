package delivery

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/lkgiovani/go-boilerplate/internal/delivery/dto"
	"github.com/lkgiovani/go-boilerplate/internal/domain/auth"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/internal/errors"
)

// Login handles user authentication
// POST /v1/auth/login
func (h *Handler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	// Extract device information for security monitoring
	deviceID := extractDeviceID(c)
	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()

	login := auth.Login{
		Email:    req.Email,
		Password: req.Password,
	}

	ctx := c.Context()
	userEntity, err := h.AuthService.Login(ctx, &login)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	// Generate access token using JwtService
	accessToken, err := h.JwtService.GenerateTokenFromUser(ctx, userEntity)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	// Generate and set cookie using JwtService
	httpCookie, err := h.JwtService.GenerateCookie(userEntity, convertFiberToHTTPRequest(c))
	if err != nil {
		return h.ErrorHandler(c, err)
	}
	setHTTPCookieToFiber(c, httpCookie)

	// Get token expiration time
	expiresIn := h.JwtService.GetAccessTokenExpirationSeconds()

	response := dto.LoginResponseDTO{
		UserID:      userEntity.ID,
		Email:       userEntity.Email,
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	}

	// Log device info for security monitoring (TODO: implement actual logging)
	_ = deviceID
	_ = userAgent
	_ = ipAddress

	return c.Status(fiber.StatusOK).JSON(response)
}

// Refresh handles access token refresh
// POST /v1/auth/refresh
func (h *Handler) Refresh(c *fiber.Ctx) error {
	// Extract refresh token from cookie
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return errors.Errorf(errors.EUNAUTHORIZED, "Refresh token not found")
	}

	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()

	// TODO: Validate refresh token and generate new access token
	// For now, we'll return a placeholder response
	_ = userAgent
	_ = ipAddress

	// TODO: Implement refresh token validation and rotation
	// This should:
	// 1. Validate the refresh token from database
	// 2. Check if it's expired or revoked
	// 3. Generate new access token
	// 4. Optionally rotate the refresh token
	// newAccessToken, newRefreshToken, err := h.AuthService.RefreshToken(ctx, refreshToken, userAgent, ipAddress)

	response := dto.RefreshResponseDTO{
		AccessToken: "new_access_token_placeholder",
		ExpiresIn:   h.JwtService.GetAccessTokenExpirationSeconds(),
	}

	// TODO: Set new refresh token cookie if rotated
	// addRefreshTokenCookie(c, newRefreshToken)

	return c.Status(fiber.StatusOK).JSON(response)
}

// Logout handles user logout from current device
// POST /v1/auth/logout
func (h *Handler) Logout(c *fiber.Ctx) error {
	// Extract refresh token from cookie
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return errors.Errorf(errors.EUNAUTHORIZED, "Refresh token not found")
	}

	// TODO: Revoke refresh token in database
	// ctx := c.Context()
	// err := h.AuthService.RevokeRefreshToken(ctx, refreshToken)
	// if err != nil {
	//     return h.ErrorHandler(c, err)
	// }

	// Clear authentication cookie using JwtService
	cleanCookie := h.JwtService.CleanCookie()
	setHTTPCookieToFiber(c, cleanCookie)

	return c.Status(fiber.StatusOK).JSON(dto.MessageResponse{
		Message: "Logged out successfully",
	})
}

// LogoutAll handles user logout from all devices
// POST /v1/auth/logout-all
func (h *Handler) LogoutAll(c *fiber.Ctx) error {
	// Get current user ID from context (set by auth middleware)
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return errors.Errorf(errors.EUNAUTHORIZED, "User not authenticated")
	}

	// Extract current refresh token from cookie
	refreshToken := c.Cookies("refresh_token")

	// TODO: Revoke all refresh tokens for this user
	// ctx := c.Context()
	// err = h.AuthService.RevokeAllRefreshTokens(ctx, userID, refreshToken)
	// if err != nil {
	//     return h.ErrorHandler(c, err)
	// }

	_ = userID
	_ = refreshToken

	// Clear authentication cookie using JwtService
	cleanCookie := h.JwtService.CleanCookie()
	setHTTPCookieToFiber(c, cleanCookie)

	return c.Status(fiber.StatusOK).JSON(dto.MessageResponse{
		Message: "Logged out from all devices",
	})
}

// Signup handles user registration
// POST /v1/auth/signup
func (h *Handler) Signup(c *fiber.Ctx) error {
	var req dto.SignupUserRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	newUser := &user.User{
		Name:   req.Name,
		Email:  req.Email,
		Admin:  false,
		Active: true,
		Source: "LOCAL",
	}

	newUser.Password = &req.Password

	ctx := c.Context()
	if err := h.AuthService.Register(ctx, newUser); err != nil {
		return h.ErrorHandler(c, err)
	}

	mapper := NewUserMapper()
	return c.Status(fiber.StatusCreated).JSON(mapper.ToResponseDTO(newUser))
}

// Helper functions

func extractDeviceID(c *fiber.Ctx) string {
	// Try to get device ID from header
	deviceID := c.Get("X-Device-ID")
	if deviceID == "" {
		// Generate a device ID based on user agent and IP
		// In production, you might want to use a more sophisticated method
		deviceID = c.Get("User-Agent") + "_" + c.IP()
	}
	return deviceID
}

// convertFiberToHTTPRequest creates a minimal http.Request for JwtService compatibility
func convertFiberToHTTPRequest(c *fiber.Ctx) *http.Request {
	req := &http.Request{
		Header: make(http.Header),
	}
	// Copy headers if needed
	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})
	return req
}

// setHTTPCookieToFiber converts http.Cookie to fiber.Cookie and sets it
func setHTTPCookieToFiber(c *fiber.Ctx, httpCookie *http.Cookie) {
	fiberCookie := &fiber.Cookie{
		Name:     httpCookie.Name,
		Value:    httpCookie.Value,
		Path:     httpCookie.Path,
		Domain:   httpCookie.Domain,
		MaxAge:   httpCookie.MaxAge,
		Expires:  httpCookie.Expires,
		Secure:   httpCookie.Secure,
		HTTPOnly: httpCookie.HttpOnly,
		SameSite: convertHTTPSameSiteToFiber(httpCookie.SameSite),
	}
	c.Cookie(fiberCookie)
}

// convertHTTPSameSiteToFiber converts http.SameSite to fiber SameSite
func convertHTTPSameSiteToFiber(sameSite http.SameSite) string {
	switch sameSite {
	case http.SameSiteStrictMode:
		return "Strict"
	case http.SameSiteLaxMode:
		return "Lax"
	case http.SameSiteNoneMode:
		return "None"
	default:
		return ""
	}
}
