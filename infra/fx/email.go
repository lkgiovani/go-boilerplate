package fx

import (
	"log/slog"
	"os"

	"github.com/lkgiovani/go-boilerplate/infra/config"
	"github.com/lkgiovani/go-boilerplate/internal/delivery"
	"github.com/lkgiovani/go-boilerplate/internal/domain/email"
	"github.com/lkgiovani/go-boilerplate/internal/domain/emailverification"
	"github.com/lkgiovani/go-boilerplate/internal/domain/passwordRecovery"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var EmailModule = fx.Module(
	"email",
	fx.Provide(
		provideEmailLogger,
		provideEmailSender,
		provideEmailVerificationRepository,
		provideEmailVerificationService,
		delivery.NewEmailVerificationHandler,
		providePasswordRecoveryRepository,
		providePasswordRecoveryService,
		delivery.NewPasswordRecoveryHandler,
	),
)

// provideEmailLogger provides Logger for email components
func provideEmailLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

// provideEmailSender provides the Email Sender based on configuration
func provideEmailSender(cfg *config.Config, logger *slog.Logger) (email.EmailSender, error) {
	var messagingConfig any
	switch email.ProviderType(cfg.Email.Provider) {
	case email.ProviderSMTP:
		messagingConfig = email.SMTPConfig{
			Host:     cfg.Email.SMTPHost,
			Port:     cfg.Email.SMTPPort,
			User:     cfg.Email.SMTPUser,
			Password: cfg.Email.SMTPPassword,
		}
	case email.ProviderResend:
		messagingConfig = email.ResendConfig{APIKey: cfg.Email.APIKey}
	case email.ProviderSES:
		messagingConfig = email.SESConfig{
			AccessKeyID:     cfg.Email.SESAccessKey,
			SecretAccessKey: cfg.Email.SESSecretKey,
			Region:          cfg.Email.SESRegion,
			Endpoint:        cfg.Email.SESEndpoint,
		}
	}

	return email.NewEmailSender(
		cfg.Email.Provider,
		cfg.Email.FromEmail,
		cfg.Email.FromName,
		messagingConfig,
		logger,
	)
}

// provideEmailVerificationRepository provides the Email Verification Repository
func provideEmailVerificationRepository(db *gorm.DB) emailverification.Repository {
	return emailverification.NewGormRepository(db)
}

// provideEmailVerificationService provides the Email Verification Service
func provideEmailVerificationService(
	repo emailverification.Repository,
	userRepo user.UserService,
	sender email.EmailSender,
	cfg *config.Config,
	logger *slog.Logger,
) *emailverification.Service {
	return emailverification.NewService(
		repo,
		userRepo,
		sender,
		cfg.Email.FrontendURL,
		logger,
	)
}

// providePasswordRecoveryRepository provides the Password Recovery Repository
func providePasswordRecoveryRepository(db *gorm.DB) passwordRecovery.Repository {
	return passwordRecovery.NewGormRepository(db)
}

// providePasswordRecoveryService provides the Password Recovery Service
func providePasswordRecoveryService(
	repo passwordRecovery.Repository,
	userRepo user.UserService,
	sender email.EmailSender,
	cfg *config.Config,
	logger *slog.Logger,
) *passwordRecovery.Service {
	return passwordRecovery.NewService(
		repo,
		userRepo,
		sender,
		cfg.Email.FrontendURL,
		logger,
	)
}
