package passwordRecovery

import (
	"time"

	"github.com/lkgiovani/go-boilerplate/pkg/utils"
)

type PasswordResetToken struct {
	ID        int64      `gorm:"primaryKey;autoIncrement"`
	UserID    int64      `gorm:"not null;index"`
	Email     string     `gorm:"not null;size:255"`
	Token     string     `gorm:"uniqueIndex;not null;size:255"`
	ExpiresAt time.Time  `gorm:"not null"`
	Used      bool       `gorm:"not null;default:false"`
	UsedAt    *time.Time `gorm:"column:used_at"`
	CreatedAt time.Time  `gorm:"not null;autoCreateTime"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

func (p *PasswordResetToken) IsExpired() bool {
	return utils.Now().Unix() > p.ExpiresAt.UTC().Unix()
}

func (p *PasswordResetToken) MarkAsUsed() {
	now := utils.Now()
	p.Used = true
	p.UsedAt = &now
}
