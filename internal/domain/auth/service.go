package auth

import (
	"context"
	"strings"
	"unicode"

	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/internal/errors"
)

type Service struct {
	Repository  user.UserService
	UserService *user.Service
}

func NewService(
	repo user.UserService,
	userSvc *user.Service,
) *Service {
	return &Service{
		Repository:  repo,
		UserService: userSvc,
	}
}

func (s *Service) Login(ctx context.Context, login *Login) (*user.User, error) {
	user, err := s.Repository.GetByEmail(ctx, login.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Register(ctx context.Context, user *user.User) error {
	exists, err := s.Repository.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}

	if exists != nil {
		return errors.Errorf(errors.EDUPLICATION, "user already exists")
	}

	if error := PasswordRequirements(*user.Password); error != nil {
		return errors.Errorf(errors.EINVALID, "password must meet requirements")
	}

	if error := s.Repository.Create(ctx, user); error != nil {
		return errors.Errorf(errors.EINTERNAL, "failed to create user")
	}

	return nil
}

func PasswordRequirements(password string) error {
	if len(password) < 8 {
		return errors.Errorf(errors.EINVALID, "password must be at least 8 characters long")
	}

	var hasUpper, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case strings.ContainsRune("@$!%*?&", char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.Errorf(errors.EINVALID, "password must contain at least one uppercase letter")
	}
	if !hasSpecial {
		return errors.Errorf(errors.EINVALID, "password must contain at least one special character")
	}

	return nil
}
