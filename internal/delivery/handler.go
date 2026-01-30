package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lkgiovani/go-boilerplate/internal/domain/auth"
	"github.com/lkgiovani/go-boilerplate/internal/domain/emailverification"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/internal/security/jwt"
)

type Handler struct {
	AuthService              *auth.Service
	UserService              *user.Service
	EmailVerificationService *emailverification.Service
	EmailVerificationHandler *EmailVerificationHandler
	PasswordRecoveryHandler  *PasswordRecoveryHandler
	UploadHandler            *UploadHandler
	MobileAuthHandler        *MobileAuthHandler
	JwtService               *jwt.JwtService

	ErrorHandler func(c *fiber.Ctx, err error) error
}

func NewHandler(
	AuthService *auth.Service,
	UserService *user.Service,
	EmailVerificationService *emailverification.Service,
	EmailVerificationHandler *EmailVerificationHandler,
	PasswordRecoveryHandler *PasswordRecoveryHandler,
	UploadHandler *UploadHandler,
	MobileAuthHandler *MobileAuthHandler,
	JwtService *jwt.JwtService,

	ErrorHandler func(c *fiber.Ctx, err error) error,
) *Handler {
	return &Handler{
		AuthService:              AuthService,
		UserService:              UserService,
		EmailVerificationService: EmailVerificationService,
		EmailVerificationHandler: EmailVerificationHandler,
		PasswordRecoveryHandler:  PasswordRecoveryHandler,
		UploadHandler:            UploadHandler,
		MobileAuthHandler:        MobileAuthHandler,
		JwtService:               JwtService,

		ErrorHandler: ErrorHandler,
	}
}

func (h *Handler) HealthCheckHandler(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "ok",
	})
}
