package email

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// SendGridSender implements EmailSender using SendGrid API
type SendGridSender struct {
	config     *EmailConfig[SendGridConfig]
	logger     *slog.Logger
	httpClient *http.Client
}

// NewSendGridSender creates a new SendGrid email sender
func NewSendGridSender(config *EmailConfig[SendGridConfig], logger *slog.Logger) (EmailSender, error) {
	if config.ConfigMessaging.APIKey == "" {
		return nil, fmt.Errorf("SendGrid API key is required")
	}
	return &SendGridSender{
		config:     config,
		logger:     logger,
		httpClient: &http.Client{},
	}, nil
}

type sendGridRequest struct {
	Personalizations []sendGridPersonalization `json:"personalizations"`
	From             sendGridEmail             `json:"from"`
	Subject          string                    `json:"subject"`
	Content          []sendGridContent         `json:"content"`
}

type sendGridPersonalization struct {
	To      []sendGridEmail        `json:"to"`
	Subject string                 `json:"subject,omitempty"`
	Dynamic map[string]interface{} `json:"dynamic_template_data,omitempty"`
}

type sendGridEmail struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type sendGridContent struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (s *SendGridSender) SendEmail(ctx context.Context, to, subject, body string) error {
	s.logger.Debug("Sending email via SendGrid",
		slog.String("to", to),
		slog.String("subject", subject),
	)

	req := sendGridRequest{
		Personalizations: []sendGridPersonalization{
			{
				To: []sendGridEmail{{Email: to}},
			},
		},
		From: sendGridEmail{
			Email: s.config.FromEmail,
			Name:  s.config.FromName,
		},
		Subject: subject,
		Content: []sendGridContent{
			{Type: "text/html", Value: body},
		},
	}

	if err := s.sendRequest(ctx, req); err != nil {
		s.logger.Error("Failed to send email via SendGrid",
			slog.String("to", to),
			slog.Any("error", err),
		)
		return err
	}

	s.logger.Info("Email sent successfully via SendGrid", slog.String("to", to))
	return nil
}

func (s *SendGridSender) SendEmailWithAttachment(ctx context.Context, to, subject, body string, attachment Attachment) error {
	s.logger.Warn("SendEmailWithAttachment not fully implemented for SendGrid, sending without attachment")
	return s.SendEmail(ctx, to, subject, body)
}

func (s *SendGridSender) SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error {
	s.logger.Debug("Sending bulk email via SendGrid",
		slog.Int("recipientCount", len(recipients)),
		slog.String("subject", subject),
	)

	toList := make([]sendGridEmail, len(recipients))
	for i, r := range recipients {
		toList[i] = sendGridEmail{Email: r}
	}

	req := sendGridRequest{
		Personalizations: []sendGridPersonalization{
			{To: toList},
		},
		From: sendGridEmail{
			Email: s.config.FromEmail,
			Name:  s.config.FromName,
		},
		Subject: subject,
		Content: []sendGridContent{
			{Type: "text/html", Value: body},
		},
	}

	if err := s.sendRequest(ctx, req); err != nil {
		s.logger.Error("Failed to send bulk email via SendGrid", slog.Any("error", err))
		return err
	}

	s.logger.Info("Bulk email sent successfully via SendGrid",
		slog.Int("recipientCount", len(recipients)),
	)
	return nil
}

func (s *SendGridSender) sendRequest(ctx context.Context, payload sendGridRequest) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.config.ConfigMessaging.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("SendGrid API returned status %d", resp.StatusCode)
	}

	return nil
}
