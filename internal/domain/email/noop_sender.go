package email

import (
	"context"
	"log/slog"
)

// NoOpSender is a no-operation email sender for development/testing
type NoOpSender struct {
	logger *slog.Logger
}

// NewNoOpSender creates a new no-op email sender that just logs emails
func NewNoOpSender(logger *slog.Logger) EmailSender {
	return &NoOpSender{logger: logger}
}

func (s *NoOpSender) SendEmail(ctx context.Context, to, subject, body string) error {
	s.logger.Info("[NoOpSender] Would send email",
		slog.String("to", to),
		slog.String("subject", subject),
		slog.Int("bodyLength", len(body)),
	)
	return nil
}

func (s *NoOpSender) SendEmailWithAttachment(ctx context.Context, to, subject, body string, attachment Attachment) error {
	s.logger.Info("[NoOpSender] Would send email with attachment",
		slog.String("to", to),
		slog.String("subject", subject),
		slog.String("filename", attachment.Filename),
	)
	return nil
}

func (s *NoOpSender) SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error {
	s.logger.Info("[NoOpSender] Would send bulk email",
		slog.Int("recipientCount", len(recipients)),
		slog.String("subject", subject),
	)
	return nil
}
