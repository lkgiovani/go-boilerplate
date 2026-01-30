package passwordRecovery

import (
	"context"

	"github.com/lkgiovani/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, token *PasswordResetToken) error

	FindByToken(ctx context.Context, token string) (*PasswordResetToken, error)

	FindByTokenIncludingUsed(ctx context.Context, token string) (*PasswordResetToken, error)

	MarkAllAsUsedByUserID(ctx context.Context, userID int64) error

	Save(ctx context.Context, token *PasswordResetToken) error

	DeleteExpired(ctx context.Context) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, token *PasswordResetToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *GormRepository) FindByToken(ctx context.Context, token string) (*PasswordResetToken, error) {
	var t PasswordResetToken
	if err := r.db.WithContext(ctx).
		Where("token = ? AND used = false", token).
		First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *GormRepository) FindByTokenIncludingUsed(ctx context.Context, token string) (*PasswordResetToken, error) {
	var t PasswordResetToken
	if err := r.db.WithContext(ctx).
		Where("token = ?", token).
		First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *GormRepository) MarkAllAsUsedByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Model(&PasswordResetToken{}).
		Where("user_id = ? AND used = false", userID).
		Updates(map[string]interface{}{
			"used":    true,
			"used_at": utils.Now(),
		}).Error
}

func (r *GormRepository) Save(ctx context.Context, token *PasswordResetToken) error {
	return r.db.WithContext(ctx).Save(token).Error
}

func (r *GormRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", utils.Now()).
		Delete(&PasswordResetToken{}).Error
}
