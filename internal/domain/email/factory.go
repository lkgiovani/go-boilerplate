package email

import (
	"fmt"

	"github.com/lkgiovani/go-boilerplate/pkg/logger"
)

type ProviderType string

const (
	ProviderSMTP   ProviderType = "smtp"
	ProviderNoOp   ProviderType = "noop"
	ProviderResend ProviderType = "resend"
	ProviderSES    ProviderType = "ses"
)

func NewEmailSender(
	provider, fromEmail, fromName string,
	messagingConfig any,
	logger logger.Logger,
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
		cfg, ok := messagingConfig.(SESConfig)
		if !ok {
			return nil, fmt.Errorf("invalid messaging config for SES")
		}
		return NewSESSender(&EmailConfig[SESConfig]{
			Provider:        provider,
			FromEmail:       fromEmail,
			FromName:        fromName,
			ConfigMessaging: cfg,
		}, logger)

	default:
		return nil, fmt.Errorf("unknown email provider: %s", provider)
	}
}
