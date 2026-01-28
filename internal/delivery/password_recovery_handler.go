package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lkgiovani/go-boilerplate/internal/domain/passwordRecovery"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
)

type PasswordRecoveryHandler struct {
	service      *passwordRecovery.Service
	userRepo     user.UserService
	ErrorHandler func(c *fiber.Ctx, err error) error
}

func NewPasswordRecoveryHandler(
	service *passwordRecovery.Service,
	userRepo user.UserService,
	errorHandler func(c *fiber.Ctx, err error) error,
) *PasswordRecoveryHandler {
	return &PasswordRecoveryHandler{
		service:      service,
		userRepo:     userRepo,
		ErrorHandler: errorHandler,
	}
}

// RequestPasswordRecovery handles the first step: requesting a reset link
func (h *PasswordRecoveryHandler) RequestPasswordRecovery(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return h.ErrorHandler(c, fiber.ErrBadRequest)
	}

	if req.Email == "" {
		return h.ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "Email é obrigatório"))
	}

	ctx := c.UserContext()
	u, err := h.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		// Security: don't reveal if email exists
		return c.JSON(fiber.Map{"message": "Se o email estiver cadastrado, um link de recuperação será enviado."})
	}

	_, err = h.service.CreateAndSendRecoveryToken(ctx, u)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.JSON(fiber.Map{"message": "Se o email estiver cadastrado, um link de recuperação será enviado."})
}

// VerifyPasswordRecovery matches POST /verify
func (h *PasswordRecoveryHandler) VerifyPasswordRecovery(c *fiber.Ctx) error {
	var req struct {
		Token string `json:"token"`
	}

	if err := c.BodyParser(&req); err != nil {
		return h.ErrorHandler(c, fiber.ErrBadRequest)
	}

	if req.Token == "" {
		return h.ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "Token é obrigatório"))
	}

	_, err := h.service.VerifyToken(c.UserContext(), req.Token)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.JSON(fiber.Map{"valid": true})
}

// VerifyPasswordRecoveryByQuery matches GET /verify?token=...
func (h *PasswordRecoveryHandler) VerifyPasswordRecoveryByQuery(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return h.ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "Token é obrigatório"))
	}

	_, err := h.service.VerifyToken(c.UserContext(), token)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.JSON(fiber.Map{"valid": true})
}

// ResetPassword matches POST /reset
func (h *PasswordRecoveryHandler) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return h.ErrorHandler(c, fiber.ErrBadRequest)
	}

	if req.Token == "" || req.Password == "" {
		return h.ErrorHandler(c, fiber.NewError(fiber.StatusBadRequest, "Token e senha são obrigatórios"))
	}

	err := h.service.ResetPassword(c.UserContext(), req.Token, req.Password)
	if err != nil {
		return h.ErrorHandler(c, err)
	}

	return c.JSON(fiber.Map{"message": "Senha redefinida com sucesso!"})
}
