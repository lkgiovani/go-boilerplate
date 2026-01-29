package email

import "context"

// EmailSender is the main interface for sending emails.
type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, body string) error
	SendEmailWithAttachment(ctx context.Context, to, subject, body string, attachment Attachment) error
	SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error
}

// Attachment represents an email attachment
type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

// EmailTemplate represents a templated email
type EmailTemplate struct {
	TemplateName string
	Subject      string
	Variables    map[string]interface{}
}

// TemplatedEmailSender extends EmailSender with template support
type TemplatedEmailSender interface {
	EmailSender
	SendTemplatedEmail(ctx context.Context, to string, template EmailTemplate) error
}

// SMTPConfig holds configuration for SMTP provider
type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

// ResendConfig holds configuration for Resend provider
type ResendConfig struct {
	APIKey string
}

// SESConfig holds configuration for AWS SES provider
type SESConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Endpoint        string // Optional
}

// EmailProviderConfig is a constraint for email provider configurations
type EmailProviderConfig interface {
	SMTPConfig | ResendConfig | SESConfig
}

// EmailConfig holds configuration for email providers with a generic messaging config
type EmailConfig[T EmailProviderConfig] struct {
	Provider  string
	FromEmail string
	FromName  string

	ConfigMessaging T
}
