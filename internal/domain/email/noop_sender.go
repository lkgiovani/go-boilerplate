package email

import (
	"context"

	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

type NoOpSender struct {
	logger logger.Logger
}

func NewNoOpSender(logger logger.Logger) EmailSender {
	return &NoOpSender{logger: logger}
}

func (s *NoOpSender) SendEmail(ctx context.Context, to, subject, body string) error {
	s.logger.Info("[NoOpSender] Would send email",
		zap.String("to", to),
		zap.String("subject", subject),
		zap.Int("bodyLength", len(body)),
	)
	return nil
}

func (s *NoOpSender) SendEmailWithAttachment(ctx context.Context, to, subject, body string, attachment Attachment) error {
	s.logger.Info("[NoOpSender] Would send email with attachment",
		zap.String("to", to),
		zap.String("subject", subject),
		zap.String("filename", attachment.Filename),
	)
	return nil
}

func (s *NoOpSender) SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error {
	s.logger.Info("[NoOpSender] Would send bulk email",
		zap.Int("recipientCount", len(recipients)),
		zap.String("subject", subject),
	)
	return nil
}
