package user

import (
	"time"
)

type User struct {
	ID         int64      `gorm:"primaryKey;autoIncrement"`
	Name       string     `gorm:"not null"`
	Email      string     `gorm:"uniqueIndex;not null"`
	Password   *string    `gorm:"column:password"`
	ImgURL     *string    `gorm:"column:img_url;size:500"`
	Admin      bool       `gorm:"not null;default:false"`
	Active     bool       `gorm:"not null;default:true;index"`
	Source     string     `gorm:"not null;default:'LOCAL';size:50"`
	Metadata   *string    `gorm:"type:jsonb"`
	LastAccess *time.Time `gorm:"column:last_access"`
	CreatedAt  time.Time  `gorm:"not null"`
	UpdatedAt  time.Time
}
