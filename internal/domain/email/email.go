package email

import "context"

type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, body string) error
	SendEmailWithAttachment(ctx context.Context, to, subject, body string, attachment Attachment) error
	SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error
}

type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

type EmailTemplate struct {
	TemplateName string
	Subject      string
	Variables    map[string]interface{}
}

type TemplatedEmailSender interface {
	EmailSender
	SendTemplatedEmail(ctx context.Context, to string, template EmailTemplate) error
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

type ResendConfig struct {
	APIKey string
}

type SESConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
	Endpoint        string
}

type EmailProviderConfig interface {
	SMTPConfig | ResendConfig | SESConfig
}

type EmailConfig[T EmailProviderConfig] struct {
	Provider  string
	FromEmail string
	FromName  string

	ConfigMessaging T
}
