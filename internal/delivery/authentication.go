package delivery

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/lkgiovani/go-boilerplate/internal/delivery/dto"
	"github.com/lkgiovani/go-boilerplate/internal/domain/auth"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/internal/errors"
	"github.com/lkgiovani/go-boilerplate/internal/security/jwt"
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

	ctx := c.UserContext()
	userEntity, err := h.AuthService.Login(ctx, &login)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	// Generate and set cookies (access and refresh)
	accessToken, _, cookies, err := h.AuthService.CreateSession(ctx, userEntity, userAgent, ipAddress, deviceID)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	for _, cookie := range cookies {
		setHTTPCookieToFiber(c, cookie)
	}

	// Get token expiration times
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
	// Get refresh token strictly from cookie
	refreshToken := c.Cookies(jwt.RefreshTokenCookieName)

	if refreshToken == "" {
		return errors.Errorf(errors.EUNAUTHORIZED, "Refresh token not found in cookies")
	}

	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()
	deviceID := extractDeviceID(c)

	ctx := c.UserContext()
	newAccessToken, _, cookies, err := h.AuthService.RefreshToken(ctx, refreshToken, userAgent, ipAddress, deviceID)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	for _, cookie := range cookies {
		setHTTPCookieToFiber(c, cookie)
	}

	response := dto.RefreshResponseDTO{
		AccessToken: newAccessToken,
		ExpiresIn:   int(h.JwtService.GetAccessTokenExpirationSeconds()),
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// Logout handles user logout from current device
// POST /v1/auth/logout
func (h *Handler) Logout(c *fiber.Ctx) error {
	// Extract refresh token strictly from cookie
	refreshToken := c.Cookies(jwt.RefreshTokenCookieName)

	// Revoke refresh token in database
	ctx := c.Context()
	err := h.AuthService.RevokeRefreshToken(ctx, refreshToken)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	// Thoroughly clear all cookies from the page
	cookies := h.JwtService.CleanAllFromHeader(c.Get("Cookie"))
	for _, cookie := range cookies {
		setHTTPCookieToFiber(c, cookie)
	}

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

	// Extract current refresh token strictly from cookie
	refreshToken := c.Cookies(jwt.RefreshTokenCookieName)

	// Revoke all refresh tokens for this user except the current one
	ctx := c.Context()
	err := h.AuthService.RevokeAllRefreshTokens(ctx, userID, refreshToken)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.MessageResponse{
		Message: "Logged out from all other devices",
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
		Active: false,
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
	// We use the standard cookie string to avoid Fiber's overwriting behavior
	// which happens when using c.Cookie() for multiple cookies with the same name
	// Using c.Response().Header.Add instead of c.Append because c.Append joins with commas
	// and Set-Cookie headers should be separate.
	c.Response().Header.Add("Set-Cookie", httpCookie.String())
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
