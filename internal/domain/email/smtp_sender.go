package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

// SMTPSender implements EmailSender using SMTP
type SMTPSender struct {
	config *EmailConfig[SMTPConfig]
	logger logger.Logger
}

// NewSMTPSender creates a new SMTP email sender
func NewSMTPSender(config *EmailConfig[SMTPConfig], logger logger.Logger) (EmailSender, error) {
	if config.ConfigMessaging.Host == "" || config.ConfigMessaging.Port == 0 || config.ConfigMessaging.User == "" || config.ConfigMessaging.Password == "" {
		return nil, fmt.Errorf("SMTP host and port are required")
	}

	return &SMTPSender{
		config: config,
		logger: logger,
	}, nil
}

func (s *SMTPSender) SendEmail(ctx context.Context, to, subject, body string) error {
	s.logger.Debug("Sending email via SMTP",
		zap.String("to", to),
		zap.String("subject", subject),
	)

	msg := s.buildMessage(to, subject, body)
	smtpCfg := s.config.ConfigMessaging
	addr := fmt.Sprintf("%s:%d", smtpCfg.Host, smtpCfg.Port)

	var err error
	if smtpCfg.Port == 465 {
		err = s.sendWithTLS(addr, to, msg)
	} else {
		err = s.sendWithStartTLS(addr, to, msg)
	}

	if err != nil {
		s.logger.Error("Failed to send email",
			zap.String("to", to),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("Email sent successfully", zap.String("to", to))
	return nil
}

func (s *SMTPSender) SendEmailWithAttachment(ctx context.Context, to, subject, body string, attachment Attachment) error {
	s.logger.Warn("SendEmailWithAttachment not fully implemented for SMTP, sending without attachment")
	return s.SendEmail(ctx, to, subject, body)
}

func (s *SMTPSender) SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error {
	var lastErr error
	for _, recipient := range recipients {
		if err := s.SendEmail(ctx, recipient, subject, body); err != nil {
			s.logger.Error("Failed to send bulk email to recipient",
				zap.String("to", recipient),
				zap.Error(err),
			)
			lastErr = err
		}
	}
	return lastErr
}

func (s *SMTPSender) buildMessage(to, subject, body string) []byte {
	from := s.config.FromEmail
	if s.config.FromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	}

	headers := fmt.Sprintf("From: %s\r\n"+
		"To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n%s", from, to, subject, body)

	return []byte(headers)
}

func (s *SMTPSender) sendWithStartTLS(addr, to string, msg []byte) error {
	smtpCfg := s.config.ConfigMessaging
	auth := smtp.PlainAuth("", smtpCfg.User, smtpCfg.Password, smtpCfg.Host)
	return smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, msg)
}

func (s *SMTPSender) sendWithTLS(addr, to string, msg []byte) error {
	smtpCfg := s.config.ConfigMessaging
	tlsConfig := &tls.Config{
		ServerName: smtpCfg.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect with TLS: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpCfg.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", smtpCfg.User, smtpCfg.Password, smtpCfg.Host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	if err = client.Mail(s.config.FromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer w.Close()

	if _, err = w.Write(msg); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}
