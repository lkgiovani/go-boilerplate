package email

import (
	"fmt"
	"log/slog"
)

// ProviderType defines supported email provider types
type ProviderType string

const (
	ProviderSMTP     ProviderType = "smtp"
	ProviderNoOp     ProviderType = "noop"
	ProviderSendGrid ProviderType = "sendgrid"
	ProviderResend   ProviderType = "resend"
	ProviderSES      ProviderType = "ses"
)

// NewEmailSender creates an EmailSender based on the configuration providers.
// Since EmailConfig is now generic, this factory receives specific configurations
// for each provider type via type assertions.
func NewEmailSender(
	provider, fromEmail, fromName string,
	messagingConfig any,
	logger *slog.Logger,
) (EmailSender, error) {
	switch ProviderType(provider) {
	case ProviderSMTP:
		cfg, ok := messagingConfig.(SMTPConfig)
		if !ok {
			return nil, fmt.Errorf("invalid messaging config for SMTP")
		}
		return NewSMTPSender(&EmailConfig[SMTPConfig]{
			Provider:        provider,
			FromEmail:       fromEmail,
			FromName:        fromName,
			ConfigMessaging: cfg,
		}, logger)

	case ProviderSendGrid:
		cfg, ok := messagingConfig.(SendGridConfig)
		if !ok {
			return nil, fmt.Errorf("invalid messaging config for SendGrid")
		}
		return NewSendGridSender(&EmailConfig[SendGridConfig]{
			Provider:        provider,
			FromEmail:       fromEmail,
			FromName:        fromName,
			ConfigMessaging: cfg,
		}, logger)

	case ProviderResend:
		cfg, ok := messagingConfig.(ResendConfig)
		if !ok {
			return nil, fmt.Errorf("invalid messaging config for Resend")
		}
		return NewResendSender(&EmailConfig[ResendConfig]{
			Provider:        provider,
			FromEmail:       fromEmail,
			FromName:        fromName,
			ConfigMessaging: cfg,
		}, logger)

	case ProviderNoOp, "":
		logger.Warn("Using NoOp email sender - emails will be logged but not sent")
		return NewNoOpSender(logger), nil

	case ProviderSES:
		return nil, fmt.Errorf("AWS SES email provider not implemented yet")

	default:
		return nil, fmt.Errorf("unknown email provider: %s", provider)
	}
}
