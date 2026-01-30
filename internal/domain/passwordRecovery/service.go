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

func (s *Service) CreateAndSendRecoveryToken(ctx context.Context, u *user.User) (*PasswordResetToken, error) {
	s.logger.Debug("Creating password recovery token for user", zap.Int64("userId", u.ID))

	if err := s.tokenRepo.MarkAllAsUsedByUserID(ctx, u.ID); err != nil {
		s.logger.Warn("Failed to mark existing tokens as used", zap.Error(err))
	}

	tokenCode, err := generateSecureToken()
	if err != nil {
		return nil, errors.Errorf(errors.EINTERNAL, "failed to generate secure token")
	}

	token := &PasswordResetToken{
		UserID:    u.ID,
		Email:     u.Email,
		Token:     tokenCode,
		ExpiresAt: utils.Now().Add(TokenExpirationHours * time.Hour),
		Used:      false,
	}

	if err := s.tokenRepo.Create(ctx, token); err != nil {
		s.logger.Error("Failed to save recovery token", zap.Error(err))
		return nil, errors.Errorf(errors.EINTERNAL, "failed to create recovery token")
	}

	go func() {

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

func (s *Service) ResetPassword(ctx context.Context, tokenCode, newPassword string) error {
	s.logger.Debug("Resetting password with token")

	token, err := s.VerifyToken(ctx, tokenCode)
	if err != nil {
		return err
	}

	hashedPassword, err := encrypt.HashPassword(newPassword)
	if err != nil {
		return errors.Errorf(errors.EINTERNAL, "Erro ao processar nova senha")
	}

	u, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return errors.Errorf(errors.ENOTFOUND, "Usuário não encontrado")
	}

	u.Password = &hashedPassword
	if err := s.userRepo.Update(ctx, u); err != nil {
		return errors.Errorf(errors.EINTERNAL, "Erro ao atualizar senha")
	}

	token.MarkAsUsed()
	if err := s.tokenRepo.Save(ctx, token); err != nil {
		s.logger.Error("Failed to mark token as used after reset", zap.Error(err))

	}

	_ = s.tokenRepo.MarkAllAsUsedByUserID(ctx, u.ID)

	s.logger.Info("Password reset successfully", zap.Int64("userId", u.ID))
	return nil
}

func (s *Service) sendPasswordResetEmail(ctx context.Context, u *user.User, tokenCode string) error {

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", s.frontendURL, tokenCode)

	subject := "Recuperação de Senha"

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
