package emailverification

import (
	"time"
)

// EmailVerificationToken represents a token for email verification
// Maps to email_verification_tokens table (see V4 migration)
type EmailVerificationToken struct {
	ID         int64      `gorm:"primaryKey;autoIncrement"`
	UserID     int64      `gorm:"not null;index"`
	Email      string     `gorm:"not null;size:255"`
	Token      string     `gorm:"uniqueIndex;not null;size:255"`
	ExpiresAt  time.Time  `gorm:"not null"`
	VerifiedAt *time.Time `gorm:"column:verified_at"`
	Used       bool       `gorm:"not null;default:false"`
	CreatedAt  time.Time  `gorm:"not null;autoCreateTime"`
}

// TableName specifies the table name for GORM
func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}

// IsExpired checks if the token has expired
func (e *EmailVerificationToken) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// MarkAsUsed marks the token as used
func (e *EmailVerificationToken) MarkAsUsed() {
	now := time.Now()
	e.Used = true
	e.VerifiedAt = &now
}

// VerifyEmailResult represents the result of email verification
type VerifyEmailResult struct {
	Success bool   `json:"success"`
	UserID  int64  `json:"userId,omitempty"`
	Email   string `json:"email,omitempty"`
	Message string `json:"message"`
}

// NewSuccessResult creates a successful verification result
func NewSuccessResult(userID int64, email, message string) VerifyEmailResult {
	return VerifyEmailResult{
		Success: true,
		UserID:  userID,
		Email:   email,
		Message: message,
	}
}

// NewFailureResult creates a failed verification result
func NewFailureResult(message string) VerifyEmailResult {
	return VerifyEmailResult{
		Success: false,
		Message: message,
	}
}
