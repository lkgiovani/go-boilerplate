package user

import (
	"context"
)

type Service struct {
	Repository UserService
}

func NewService(repo UserService) *Service {
	return &Service{
		Repository: repo,
	}
}

type UserService interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	DeleteByIDs(ctx context.Context, ids []int64) error

	FindAll(ctx context.Context, page, size int) ([]User, int64, error)
	FindAllWithFilter(ctx context.Context, keyword string, page, size int) ([]User, int64, error)

	ToggleStatus(ctx context.Context, id int64, active bool) error

	RequestPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	ChangePassword(ctx context.Context, id int64, currentPassword, newPassword string) error
	ResetUserPassword(ctx context.Context, id int64, newPassword string) error

	UpdateAccessMode(ctx context.Context, id int64, accessMode string) (*User, error)
	UpdateFeatures(ctx context.Context, id int64, canCreateBudgets, canExportData, canUseReports, canUseGoals *bool) (*User, error)
	UpdateLimits(ctx context.Context, id int64, maxAccounts, maxTransactionsPerMonth, maxCategoriesPerAccount *int) (*User, error)
	GrantLifetimePro(ctx context.Context, id int64, reason string) (*User, error)
	RevokeLifetimePro(ctx context.Context, id int64) (*User, error)
	EnsureMetadata(ctx context.Context, id int64) (*User, error)
}
