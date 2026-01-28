package delivery

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/lkgiovani/go-boilerplate/internal/delivery/dto"
	"github.com/lkgiovani/go-boilerplate/internal/domain/emailverification"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/internal/errors"
)

// EmailVerificationHandler handles email verification endpoints
type EmailVerificationHandler struct {
	service      *emailverification.Service
	userRepo     user.UserService
	ErrorHandler func(c *fiber.Ctx, err error) error
}

// NewEmailVerificationHandler creates a new email verification handler
func NewEmailVerificationHandler(
	service *emailverification.Service,
	userRepo user.UserService,
	errorHandler func(c *fiber.Ctx, err error) error,
) *EmailVerificationHandler {
	return &EmailVerificationHandler{
		service:      service,
		userRepo:     userRepo,
		ErrorHandler: errorHandler,
	}
}

// SendVerificationEmail sends a verification email to the authenticated user
// POST /v1/email-verification/send
func (h *EmailVerificationHandler) SendVerificationEmail(c *fiber.Ctx) error {
	// Get current user from context (set by auth middleware)
	userID, ok := c.Locals("userID").(int64)
	if !ok {
		return errors.Errorf(errors.EUNAUTHORIZED, "User not authenticated")
	}

	// Get user from database
	ctx := c.UserContext()
	u, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		return h.ErrorHandler(c, errors.Errorf(errors.ENOTFOUND, "User not found"))
	}

	// Check if already verified
	if u.Active {
		return c.Status(fiber.StatusOK).JSON(dto.EmailVerificationResponse{
			Success: true,
			Message: "Email já está verificado",
		})
	}

	_, err = h.service.CreateAndSendVerificationToken(ctx, u)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.EmailVerificationResponse{
		Success: true,
		Message: "Email de verificação enviado com sucesso!",
	})
}

// VerifyEmail verifies an email using the token in the request body
// POST /v1/email-verification/verify
func (h *EmailVerificationHandler) VerifyEmail(c *fiber.Ctx) error {
	var req dto.VerifyEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	if req.Token == "" {
		return errors.Errorf(errors.EBADREQUEST, "Token is required")
	}

	ctx := c.UserContext()
	result := h.service.VerifyToken(ctx, req.Token)

	return c.Status(fiber.StatusOK).JSON(dto.EmailVerificationResponse{
		Success: result.Success,
		Message: result.Message,
		UserID:  result.UserID,
		Email:   result.Email,
	})
}

// VerifyEmailByQuery verifies an email using the token in query parameter
// GET /v1/email-verification/verify?token=xxx
func (h *EmailVerificationHandler) VerifyEmailByQuery(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return errors.Errorf(errors.EBADREQUEST, "Token is required")
	}

	ctx := c.UserContext()
	result := h.service.VerifyToken(ctx, token)

	// If it's a GET request, it's likely a human clicking the email link
	// Return a simple HTML page instead of JSON
	c.Set("Content-Type", "text/html")
	if result.Success {
		return c.Status(fiber.StatusOK).SendString(`
			<!DOCTYPE html>
			<html lang="pt-BR">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Email Verificado</title>
				<style>
					body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #f0f2f5; }
					.card { background: white; padding: 2.5rem; border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.08); text-align: center; max-width: 450px; width: 90%; }
					.icon { font-size: 4rem; margin-bottom: 1rem; }
					h1 { color: #1a1a1a; margin-bottom: 1rem; font-size: 1.5rem; }
					p { color: #666; line-height: 1.6; margin-bottom: 1.5rem; }
					.success-icon { color: #4CAF50; }
					.btn { background-color: #007bff; color: white; border: none; padding: 0.75rem 1.5rem; border-radius: 6px; font-weight: 600; cursor: pointer; text-decoration: none; display: inline-block; transition: background-color 0.2s; }
					.btn:hover { background-color: #0056b3; }
				</style>
			</head>
			<body>
				<div class="card">
					<div class="icon success-icon">✅</div>
					<h1>Email Verificado!</h1>
					<p>Excelente! Seu endereço de email foi confirmado com sucesso. Agora você tem acesso total a todas as funcionalidades do sistema.</p>
					<p>Você já pode fechar esta aba e voltar para o aplicativo.</p>
				</div>
			</body>
			</html>
		`)
	}

	return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="pt-BR">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Erro na Verificação</title>
			<style>
				body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background-color: #f0f2f5; }
				.card { background: white; padding: 2.5rem; border-radius: 12px; box-shadow: 0 4px 20px rgba(0,0,0,0.08); text-align: center; max-width: 450px; width: 90%; }
				.icon { font-size: 4rem; margin-bottom: 1rem; }
				h1 { color: #1a1a1a; margin-bottom: 1rem; font-size: 1.5rem; }
				p { color: #666; line-height: 1.6; margin-bottom: 1.5rem; }
				.error-icon { color: #f44336; }
			</style>
		</head>
		<body>
			<div class="card">
				<div class="icon error-icon">❌</div>
				<h1>Ops! Algo deu errado</h1>
				<p>%s</p>
				<p>Por favor, verifique se o link está correto ou tente solicitar um novo código de verificação no aplicativo.</p>
			</div>
		</body>
		</html>
	`, result.Message))
}

// ResendVerificationEmail resends verification email to the specified email address
// POST /v1/email-verification/resend
func (h *EmailVerificationHandler) ResendVerificationEmail(c *fiber.Ctx) error {
	var req dto.ResendEmailVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return errors.Errorf(errors.EBADREQUEST, "Invalid request body")
	}

	if req.Email == "" {
		return errors.Errorf(errors.EBADREQUEST, "Email is required")
	}

	ctx := c.UserContext()
	if err := h.service.ResendVerification(ctx, req.Email); err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.EmailVerificationResponse{
		Success: true,
		Message: "Email de verificação reenviado com sucesso!",
	})
}
