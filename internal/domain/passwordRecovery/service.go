package passwordRecovery

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/lkgiovani/go-boilerplate/internal/domain/email"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/internal/errors"
	"github.com/lkgiovani/go-boilerplate/pkg/encrypt"
	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"github.com/lkgiovani/go-boilerplate/pkg/utils"
	"go.uber.org/zap"
)

const (
	// TokenExpirationHours defines how long a password reset token is valid
	TokenExpirationHours = 1
)

type Service struct {
	tokenRepo   Repository
	userRepo    user.UserService
	emailSender email.EmailSender
	frontendURL string
	logger      logger.Logger
}

func NewService(
	tokenRepo Repository,
	userRepo user.UserService,
	emailSender email.EmailSender,
	frontendURL string,
	logger logger.Logger,
) *Service {
	return &Service{
		tokenRepo:   tokenRepo,
		userRepo:    userRepo,
		emailSender: emailSender,
		frontendURL: frontendURL,
		logger:      logger,
	}
}

// CreateAndSendRecoveryToken generates a new token and sends a recovery email
func (s *Service) CreateAndSendRecoveryToken(ctx context.Context, u *user.User) (*PasswordResetToken, error) {
	s.logger.Debug("Creating password recovery token for user", zap.Int64("userId", u.ID))

	// Mark all existing tokens as used
	if err := s.tokenRepo.MarkAllAsUsedByUserID(ctx, u.ID); err != nil {
		s.logger.Warn("Failed to mark existing tokens as used", zap.Error(err))
	}

	// Generate secure token
	tokenCode, err := generateSecureToken()
	if err != nil {
		return nil, errors.Errorf(errors.EINTERNAL, "failed to generate secure token")
	}

	// Create token entity
	token := &PasswordResetToken{
		UserID:    u.ID,
		Email:     u.Email,
		Token:     tokenCode,
		ExpiresAt: utils.Now().Add(TokenExpirationHours * time.Hour),
		Used:      false,
	}

	// Save token to database
	if err := s.tokenRepo.Create(ctx, token); err != nil {
		s.logger.Error("Failed to save recovery token", zap.Error(err))
		return nil, errors.Errorf(errors.EINTERNAL, "failed to create recovery token")
	}

	// Send recovery email asynchronously
	go func() {
		// Use a detached context for the background job
		sendCtx := context.Background()
		if err := s.sendPasswordResetEmail(sendCtx, u, tokenCode); err != nil {
			s.logger.Error("Failed to send recovery email",
				zap.Int64("userId", u.ID),
				zap.Error(err),
			)
		}
	}()

	s.logger.Info("Password recovery token created and email queued", zap.Int64("userId", u.ID))
	return token, nil
}

// VerifyToken verifies the password reset token without consuming it
func (s *Service) VerifyToken(ctx context.Context, tokenCode string) (*PasswordResetToken, error) {
	s.logger.Debug("Verifying password reset token")

	token, err := s.tokenRepo.FindByToken(ctx, tokenCode)
	if err != nil {
		s.logger.Warn("Token not found or already used", zap.String("token", truncateToken(tokenCode)))
		return nil, errors.Errorf(errors.ENOTFOUND, "Token inválido ou já utilizado")
	}

	if token.IsExpired() {
		s.logger.Warn("Token expired", zap.Int64("tokenId", token.ID))
		return nil, errors.Errorf(errors.EINVALID, "Token expirado")
	}

	return token, nil
}

// ResetPassword verifies the token and updates the user's password
func (s *Service) ResetPassword(ctx context.Context, tokenCode, newPassword string) error {
	s.logger.Debug("Resetting password with token")

	// 1. Verify token
	token, err := s.VerifyToken(ctx, tokenCode)
	if err != nil {
		return err
	}

	// 2. Hash new password
	hashedPassword, err := encrypt.HashPassword(newPassword)
	if err != nil {
		return errors.Errorf(errors.EINTERNAL, "Erro ao processar nova senha")
	}

	// 3. Update user password
	u, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return errors.Errorf(errors.ENOTFOUND, "Usuário não encontrado")
	}

	u.Password = &hashedPassword
	if err := s.userRepo.Update(ctx, u); err != nil {
		return errors.Errorf(errors.EINTERNAL, "Erro ao atualizar senha")
	}

	// 4. Mark token as used
	token.MarkAsUsed()
	if err := s.tokenRepo.Save(ctx, token); err != nil {
		s.logger.Error("Failed to mark token as used after reset", zap.Error(err))
		// We don't fail the operation here as the password was already updated
	}

	// 5. Revoke all previous tokens for this user
	_ = s.tokenRepo.MarkAllAsUsedByUserID(ctx, u.ID)

	s.logger.Info("Password reset successfully", zap.Int64("userId", u.ID))
	return nil
}

// sendPasswordResetEmail sends the link to the frontend for password reset
func (s *Service) sendPasswordResetEmail(ctx context.Context, u *user.User, tokenCode string) error {
	// Link points to the frontend reset-password page
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, tokenCode)

	subject := "Recuperação de Senha"

	// Simple and direct email body as requested
	body := fmt.Sprintf(`
		<div style="font-family: sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
			<h2>Olá, %s</h2>
			<p>Recebemos uma solicitação para redefinir a sua senha.</p>
			<p>Clique no link abaixo para prosseguir:</p>
			<div style="margin: 30px 0;">
				<a href="%s" style="background-color: #007bff; color: white; padding: 12px 24px; text-decoration: none; border-radius: 4px; font-weight: bold;">
					Redefinir Senha
				</a>
			</div>
			<p>Se você não solicitou isso, pode ignorar este email.</p>
			<hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
			<p style="font-size: 12px; color: #666;">
				Ou copie e cole este link no seu navegador:<br>
				%s
			</p>
		</div>
	`, u.Name, resetLink, resetLink)

	return s.emailSender.SendEmail(ctx, u.Email, subject, body)
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(bytes), nil
}

func truncateToken(token string) string {
	if len(token) <= 8 {
		return "***"
	}
	return token[:8] + "..."
}
