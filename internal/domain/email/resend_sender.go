package email

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/resend/resend-go/v3"
)

// ResendSender implements EmailSender using Resend API
type ResendSender struct {
	config *EmailConfig[ResendConfig]
	logger *slog.Logger
	client *resend.Client
}

// NewResendSender creates a new Resend email sender
func NewResendSender(config *EmailConfig[ResendConfig], logger *slog.Logger) (EmailSender, error) {
	if config.ConfigMessaging.APIKey == "" {
		return nil, fmt.Errorf("Resend API key is required")
	}

	client := resend.NewClient(config.ConfigMessaging.APIKey)

	return &ResendSender{
		config: config,
		logger: logger,
		client: client,
	}, nil
}

func (s *ResendSender) SendEmail(ctx context.Context, to, subject, body string) error {
	s.logger.Debug("Sending email via Resend",
		slog.String("to", to),
		slog.String("subject", subject),
	)

	from := s.config.FromEmail
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	}

	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}

	_, err := s.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		s.logger.Error("Failed to send email via Resend",
			slog.String("to", to),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("Email sent successfully via Resend", slog.String("to", to))
	return nil
}

func (s *ResendSender) SendEmailWithAttachment(ctx context.Context, to, subject, body string, attachment Attachment) error {
	s.logger.Debug("Sending email with attachment via Resend",
		slog.String("to", to),
		slog.String("subject", subject),
	)

	from := s.config.FromEmail
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	}

	params := &resend.SendEmailRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		Html:    body,
		Attachments: []*resend.Attachment{
			{
				Content:  attachment.Data,
				Filename: attachment.Filename,
			},
		},
	}

	_, err := s.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		s.logger.Error("Failed to send email with attachment via Resend",
			slog.String("to", to),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("Email with attachment sent successfully via Resend", slog.String("to", to))
	return nil
}

func (s *ResendSender) SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error {
	s.logger.Debug("Sending bulk email via Resend",
		slog.Int("recipientCount", len(recipients)),
		slog.String("subject", subject),
	)

	from := s.config.FromEmail
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	}

	params := &resend.SendEmailRequest{
		From:    from,
		To:      recipients,
		Subject: subject,
		Html:    body,
	}

	_, err := s.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		s.logger.Error("Failed to send bulk email via Resend",
			slog.Int("recipientCount", len(recipients)),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("Bulk email sent successfully via Resend",
		slog.Int("recipientCount", len(recipients)),
	)
	return nil
}
