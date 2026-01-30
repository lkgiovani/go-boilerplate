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

func (h *Handler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

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

	accessToken, _, cookies, err := h.AuthService.CreateSession(ctx, userEntity, userAgent, ipAddress, deviceID)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	for _, cookie := range cookies {
		setHTTPCookieToFiber(c, cookie)
	}

	expiresIn := h.JwtService.GetAccessTokenExpirationSeconds()

	response := dto.LoginResponseDTO{
		UserID:      userEntity.ID,
		Email:       userEntity.Email,
		AccessToken: accessToken,
		ExpiresIn:   expiresIn,
	}

	_ = deviceID
	_ = userAgent
	_ = ipAddress

	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handler) Refresh(c *fiber.Ctx) error {

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

func (h *Handler) Logout(c *fiber.Ctx) error {

	refreshToken := c.Cookies(jwt.RefreshTokenCookieName)

	ctx := c.Context()
	err := h.AuthService.RevokeRefreshToken(ctx, refreshToken)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	cookies := h.JwtService.CleanAllFromHeader(c.Get("Cookie"))
	for _, cookie := range cookies {
		setHTTPCookieToFiber(c, cookie)
	}

	return c.Status(fiber.StatusOK).JSON(dto.MessageResponse{
		Message: "Logged out successfully",
	})
}

func (h *Handler) LogoutAll(c *fiber.Ctx) error {

	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return errors.Errorf(errors.EUNAUTHORIZED, "User not authenticated")
	}

	refreshToken := c.Cookies(jwt.RefreshTokenCookieName)

	ctx := c.Context()
	err := h.AuthService.RevokeAllRefreshTokens(ctx, userID, refreshToken)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.MessageResponse{
		Message: "Logged out from all other devices",
	})
}

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

func extractDeviceID(c *fiber.Ctx) string {

	deviceID := c.Get("X-Device-ID")
	if deviceID == "" {

		deviceID = c.Get("User-Agent") + "_" + c.IP()
	}
	return deviceID
}

func convertFiberToHTTPRequest(c *fiber.Ctx) *http.Request {
	req := &http.Request{
		Header: make(http.Header),
	}

	c.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})
	return req
}

func setHTTPCookieToFiber(c *fiber.Ctx, httpCookie *http.Cookie) {

	c.Response().Header.Add("Set-Cookie", httpCookie.String())
}

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
