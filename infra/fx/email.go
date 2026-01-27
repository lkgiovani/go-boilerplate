package fx

import (
	"log/slog"
	"os"

	"github.com/lkgiovani/go-boilerplate/infra/config"
	"github.com/lkgiovani/go-boilerplate/internal/delivery"
	"github.com/lkgiovani/go-boilerplate/internal/domain/email"
	"github.com/lkgiovani/go-boilerplate/internal/domain/emailverification"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

var EmailModule = fx.Module(
	"email",
	fx.Provide(
		// Provide Logger for email components (using slog for internal domains as requested)
		func() *slog.Logger {
			return slog.New(slog.NewJSONHandler(os.Stdout, nil))
		},
		// Provide Email Sender
		func(cfg *config.Config, logger *slog.Logger) (email.EmailSender, error) {
			var messagingConfig any
			switch email.ProviderType(cfg.Email.Provider) {
			case email.ProviderSMTP:
				messagingConfig = email.SMTPConfig{
					Host:     cfg.Email.SMTPHost,
					Port:     cfg.Email.SMTPPort,
					User:     cfg.Email.SMTPUser,
					Password: cfg.Email.SMTPPassword,
				}
			case email.ProviderSendGrid:
				messagingConfig = email.SendGridConfig{APIKey: cfg.Email.APIKey}
			case email.ProviderResend:
				messagingConfig = email.ResendConfig{APIKey: cfg.Email.APIKey}
			}

			return email.NewEmailSender(
				cfg.Email.Provider,
				cfg.Email.FromEmail,
				cfg.Email.FromName,
				messagingConfig,
				logger,
			)
		},
		// Provide Email Verification Repository
		func(db *gorm.DB) emailverification.Repository {
			return emailverification.NewGormRepository(db)
		},
		// Provide Email Verification Service
		func(
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
		},
		// Provide Email Verification Handler
		delivery.NewEmailVerificationHandler,
	),
)
