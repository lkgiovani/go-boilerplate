package email

import (
	"context"
	"fmt"

	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"github.com/resend/resend-go/v3"
	"go.uber.org/zap"
)

type ResendSender struct {
	config *EmailConfig[ResendConfig]
	logger logger.Logger
	client *resend.Client
}

func NewResendSender(config *EmailConfig[ResendConfig], logger logger.Logger) (EmailSender, error) {
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
		zap.String("to", to),
		zap.String("subject", subject),
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
			zap.String("to", to),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("Email sent successfully via Resend", zap.String("to", to))
	return nil
}

func (s *ResendSender) SendEmailWithAttachment(ctx context.Context, to, subject, body string, attachment Attachment) error {
	s.logger.Debug("Sending email with attachment via Resend",
		zap.String("to", to),
		zap.String("subject", subject),
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
			zap.String("to", to),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("Email with attachment sent successfully via Resend", zap.String("to", to))
	return nil
}

func (s *ResendSender) SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error {
	s.logger.Debug("Sending bulk email via Resend",
		zap.Int("recipientCount", len(recipients)),
		zap.String("subject", subject),
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
			zap.Int("recipientCount", len(recipients)),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("Bulk email sent successfully via Resend",
		zap.Int("recipientCount", len(recipients)),
	)
	return nil
}
