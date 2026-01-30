package emailverification

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"time"

	"github.com/lkgiovani/go-boilerplate/internal/domain/email"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/internal/errors"
	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"github.com/lkgiovani/go-boilerplate/pkg/utils"
	"go.uber.org/zap"
)

const (
	TokenExpirationHours = 24
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

func (s *Service) CreateAndSendVerificationToken(ctx context.Context, u *user.User) (*EmailVerificationToken, error) {
	s.logger.Debug("Creating verification token for user", zap.Int64("userId", u.ID))

	if err := s.tokenRepo.MarkAllAsUsedByUserID(ctx, u.ID); err != nil {
		s.logger.Error("Failed to mark existing tokens as used", zap.Error(err))

	}

	tokenCode, err := generateSecureToken()
	if err != nil {
		return nil, errors.Errorf(errors.EINTERNAL, "failed to generate token")
	}

	token := &EmailVerificationToken{
		UserID:    u.ID,
		Email:     u.Email,
		Token:     tokenCode,
		ExpiresAt: utils.Now().Add(TokenExpirationHours * time.Hour),
		Used:      false,
	}

	if err := s.tokenRepo.Create(ctx, token); err != nil {
		s.logger.Error("Failed to save verification token", zap.Error(err))
		return nil, errors.Errorf(errors.EINTERNAL, "failed to create verification token")
	}

	go func() {
		sendCtx := context.Background()
		if err := s.sendVerificationEmail(sendCtx, u.Email, tokenCode); err != nil {
			s.logger.Error("Failed to send verification email",
				zap.Int64("userId", u.ID),
				zap.Error(err),
			)
		}
	}()

	s.logger.Info("Verification token created and email queued", zap.Int64("userId", u.ID))
	return token, nil
}

func (s *Service) VerifyToken(ctx context.Context, tokenCode string) VerifyEmailResult {
	s.logger.Debug("Verifying token")

	token, err := s.tokenRepo.FindByTokenIncludingUsed(ctx, tokenCode)
	if err != nil {
		s.logger.Warn("Token not found", zap.String("token", truncateToken(tokenCode)))
		return NewFailureResult("Token inválido ou não encontrado")
	}

	if token.Used {
		s.logger.Warn("Token already used", zap.Int64("tokenId", token.ID))

		if token.VerifiedAt != nil && token.VerifiedAt.Before(utils.Now().Add(-1*time.Hour)) {
			return NewFailureResult("Token expirado. Solicite um novo código")
		}

		u, errUser := s.userRepo.GetByID(ctx, token.UserID)
		if errUser == nil && u.Metadata.EmailVerified {
			return NewSuccessResult(token.UserID, token.Email, "Email já foi verificado anteriormente")
		}

		return NewFailureResult("Token já foi utilizado")
	}

	if token.IsExpired() {
		s.logger.Warn("Token expired", zap.Int64("tokenId", token.ID))
		return NewFailureResult("Token expirado. Solicite um novo código")
	}

	token.MarkAsUsed()
	if err := s.tokenRepo.Save(ctx, token); err != nil {
		s.logger.Error("Failed to mark token as used", zap.Error(err))
		return NewFailureResult("Erro interno ao verificar email")
	}

	u, errGet := s.userRepo.GetByID(ctx, token.UserID)
	if errGet != nil {
		s.logger.Error("Failed to find user", zap.Error(errGet))
		return NewFailureResult("Usuário não encontrado")
	}

	u.Active = true
	u.Metadata.EmailVerified = true
	if err := s.userRepo.Update(ctx, u); err != nil {
		s.logger.Error("Failed to update user verification status", zap.Error(err))
		return NewFailureResult("Erro ao atualizar status de verificação")
	}

	s.logger.Info("Email verified successfully", zap.Int64("userId", token.UserID))
	return NewSuccessResult(token.UserID, token.Email, "Email verificado com sucesso!")
}

func (s *Service) ResendVerification(ctx context.Context, emailAddr string) error {
	s.logger.Debug("Resending verification email", zap.String("email", emailAddr))

	u, err := s.userRepo.GetByEmail(ctx, emailAddr)
	if err != nil {
		return errors.Errorf(errors.ENOTFOUND, "Usuário não encontrado")
	}

	if u.Metadata.EmailVerified {
		return errors.Errorf(errors.EINVALID, "Email já verificado")
	}

	_, err = s.CreateAndSendVerificationToken(ctx, u)
	if err != nil {
		return err
	}

	s.logger.Info("Verification email resent", zap.String("email", emailAddr))
	return nil
}

func (s *Service) sendVerificationEmail(ctx context.Context, toEmail, tokenCode string) error {
	verificationURL := s.frontendURL + "/v1/email-verification/verify?token=" + tokenCode

	subject := "Verificação de Email"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h1>Verificação de Email</h1>
			<p>Clique no link abaixo para verificar seu email:</p>
			<a href="%s">Verificar Email</a>
			<p>Ou copie e cole este link no seu navegador:</p>
			<p>%s</p>
			<p>Este link expira em 24 horas.</p>
		</body>
		</html>
	`, verificationURL, verificationURL)

	return s.emailSender.SendEmail(ctx, toEmail, subject, body)
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
