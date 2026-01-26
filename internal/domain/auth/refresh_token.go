package auth

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID      int64     `gorm:"not null"`
	UserEmail   string    `gorm:"not null"`
	DeviceID    string    `gorm:"not null"`
	Jti         string    `gorm:"uniqueIndex;not null"`
	FamilyID    uuid.UUID `gorm:"type:uuid;not null"`
	TokenHash   string    `gorm:"uniqueIndex;not null"`
	ExpiresAt   time.Time `gorm:"not null"`
	CreatedAt   time.Time `gorm:"not null"`
	Used        bool      `gorm:"not null;default:false"`
	UsedAt      *time.Time
	RotatedFrom *uuid.UUID `gorm:"type:uuid"`
	RevokedAt   *time.Time
	UserAgent   string `gorm:"size:500"`
	IpAddress   string `gorm:"size:45"`
}
