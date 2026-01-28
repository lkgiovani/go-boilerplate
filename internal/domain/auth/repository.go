package auth

import (
	"context"

	"github.com/lkgiovani/go-boilerplate/pkg/utils"
	"gorm.io/gorm"
)

type Repository interface {
	CreateRefreshToken(ctx context.Context, token *RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, hash string) (*RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, hash string) error
	RevokeAllUserRefreshTokens(ctx context.Context, userID int64) error
	RevokeAllUserRefreshTokensExcept(ctx context.Context, userID int64, exceptHash string) error
	MarkAsUsed(ctx context.Context, hash string) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) CreateRefreshToken(ctx context.Context, token *RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *GormRepository) GetRefreshTokenByHash(ctx context.Context, hash string) (*RefreshToken, error) {
	var t RefreshToken
	if err := r.db.WithContext(ctx).Where("token_hash = ?", hash).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *GormRepository) RevokeRefreshToken(ctx context.Context, hash string) error {
	return r.db.WithContext(ctx).Model(&RefreshToken{}).
		Where("token_hash = ?", hash).
		Update("revoked_at", utils.Now()).Error
}

func (r *GormRepository) RevokeAllUserRefreshTokens(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Model(&RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", utils.Now()).Error
}

func (r *GormRepository) RevokeAllUserRefreshTokensExcept(ctx context.Context, userID int64, exceptHash string) error {
	return r.db.WithContext(ctx).Model(&RefreshToken{}).
		Where("user_id = ? AND token_hash != ? AND revoked_at IS NULL", userID, exceptHash).
		Update("revoked_at", utils.Now()).Error
}

func (r *GormRepository) MarkAsUsed(ctx context.Context, hash string) error {
	return r.db.WithContext(ctx).Model(&RefreshToken{}).
		Where("token_hash = ?", hash).
		Updates(map[string]interface{}{
			"used":    true,
			"used_at": utils.Now(),
		}).Error
}
