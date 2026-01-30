package email

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"go.uber.org/zap"
)

// SESSender implements EmailSender using AWS SES V2
type SESSender struct {
	config *EmailConfig[SESConfig]
	logger logger.Logger
	client *sesv2.Client
}

// NewSESSender creates a new AWS SES email sender
func NewSESSender(cfg *EmailConfig[SESConfig], logger logger.Logger) (EmailSender, error) {
	if cfg.ConfigMessaging.AccessKeyID == "" || cfg.ConfigMessaging.SecretAccessKey == "" {
		return nil, fmt.Errorf("AWS SES credentials (access key and secret key) are required")
	}

	if cfg.ConfigMessaging.Region == "" {
		cfg.ConfigMessaging.Region = "us-east-1" // Default region
	}

	creds := credentials.NewStaticCredentialsProvider(
		cfg.ConfigMessaging.AccessKeyID,
		cfg.ConfigMessaging.SecretAccessKey,
		"",
	)

	opts := []func(*config.LoadOptions) error{
		config.WithRegion(cfg.ConfigMessaging.Region),
		config.WithCredentialsProvider(creds),
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := sesv2.NewFromConfig(awsCfg, func(o *sesv2.Options) {
		if cfg.ConfigMessaging.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.ConfigMessaging.Endpoint)
		}
	})

	return &SESSender{
		config: cfg,
		logger: logger,
		client: client,
	}, nil
}

func (s *SESSender) SendEmail(ctx context.Context, to, subject, body string) error {
	s.logger.Debug("Sending email via AWS SES",
		zap.String("to", to),
		zap.String("subject", subject),
	)

	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(s.fromAddress()),
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(subject),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(body),
					},
				},
			},
		},
	}

	_, err := s.client.SendEmail(ctx, input)
	if err != nil {
		s.logger.Error("Failed to send email via AWS SES",
			zap.String("to", to),
			zap.Error(err),
		)
		return err
	}

	s.logger.Info("Email sent successfully via AWS SES", zap.String("to", to))
	return nil
}

func (s *SESSender) SendEmailWithAttachment(ctx context.Context, to, subject, body string, attachment Attachment) error {
	s.logger.Warn("SendEmailWithAttachment not fully implemented for SES yet - sending without attachment", zap.String("to", to))
	return s.SendEmail(ctx, to, subject, body)
}

func (s *SESSender) SendBulkEmail(ctx context.Context, recipients []string, subject, body string) error {
	for _, to := range recipients {
		if err := s.SendEmail(ctx, to, subject, body); err != nil {
			return err
		}
	}
	return nil
}

func (s *SESSender) fromAddress() string {
	if s.config.FromName != "" {
		return fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	}
	return s.config.FromEmail
}
