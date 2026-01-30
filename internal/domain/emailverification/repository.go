package emailverification

import (
	"context"

	"github.com/lkgiovani/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, token *EmailVerificationToken) error

	FindByToken(ctx context.Context, token string) (*EmailVerificationToken, error)

	FindByTokenIncludingUsed(ctx context.Context, token string) (*EmailVerificationToken, error)

	MarkAllAsUsedByUserID(ctx context.Context, userID int64) error

	Save(ctx context.Context, token *EmailVerificationToken) error

	DeleteExpired(ctx context.Context) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(ctx context.Context, token *EmailVerificationToken) error {

	return r.db.WithContext(ctx).Create(token).Error
}

func (r *GormRepository) FindByToken(ctx context.Context, token string) (*EmailVerificationToken, error) {
	var t EmailVerificationToken
	if err := r.db.WithContext(ctx).
		Where("token = ? AND used = false", token).
		First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *GormRepository) FindByTokenIncludingUsed(ctx context.Context, token string) (*EmailVerificationToken, error) {
	var t EmailVerificationToken
	if err := r.db.WithContext(ctx).
		Where("token = ?", token).
		First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *GormRepository) MarkAllAsUsedByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).
		Model(&EmailVerificationToken{}).
		Where("user_id = ? AND used = false", userID).
		Updates(map[string]interface{}{
			"used":        true,
			"verified_at": utils.Now(),
		}).Error
}

func (r *GormRepository) Save(ctx context.Context, token *EmailVerificationToken) error {
	return r.db.WithContext(ctx).Save(token).Error
}

func (r *GormRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ? AND used = true", utils.Now()).
		Delete(&EmailVerificationToken{}).Error
}
