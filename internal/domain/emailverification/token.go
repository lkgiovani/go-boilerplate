package emailverification

import (
	"time"

	"github.com/lkgiovani/go-boilerplate/pkg/utils"
)

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

func (EmailVerificationToken) TableName() string {
	return "email_verification_tokens"
}

func (e *EmailVerificationToken) IsExpired() bool {
	return utils.Now().Unix() > e.ExpiresAt.UTC().Unix()
}

func (e *EmailVerificationToken) MarkAsUsed() {
	now := utils.Now()
	e.Used = true
	e.VerifiedAt = &now
}

type VerifyEmailResult struct {
	Success bool   `json:"success"`
	UserID  int64  `json:"userId,omitempty"`
	Email   string `json:"email,omitempty"`
	Message string `json:"message"`
}

func NewSuccessResult(userID int64, email, message string) VerifyEmailResult {
	return VerifyEmailResult{
		Success: true,
		UserID:  userID,
		Email:   email,
		Message: message,
	}
}

func NewFailureResult(message string) VerifyEmailResult {
	return VerifyEmailResult{
		Success: false,
		Message: message,
	}
}
