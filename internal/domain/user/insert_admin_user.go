package user

import (
	"context"

	"github.com/lkgiovani/go-boilerplate/infra/config"
	"github.com/lkgiovani/go-boilerplate/pkg/encrypt"
	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"github.com/lkgiovani/go-boilerplate/pkg/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type InsertAdminUser struct {
	userService UserService
	config      *config.Config
	logger      logger.Logger
}

func NewInsertAdminUser(userService UserService, cfg *config.Config, log logger.Logger) *InsertAdminUser {
	return &InsertAdminUser{
		userService: userService,
		config:      cfg,
		logger:      log,
	}
}

func (i *InsertAdminUser) Execute(ctx context.Context) error {
	adminEmail := i.config.Admin.Email

	// Verificar se o usuário admin já existe
	existingUser, err := i.userService.GetByEmail(ctx, adminEmail)
	if err != nil && err != gorm.ErrRecordNotFound {
		i.logger.Error("[InsertAdminUser] Error checking if admin user exists", zap.Error(err))
		return err
	}

	// Se o usuário já existe, não fazer nada
	if existingUser != nil {
		i.logger.Debug("[InsertAdminUser] Administrator user already exists")
		return nil
	}

	// Criar o usuário administrador
	i.logger.Info("[InsertAdminUser] Administrator user not found, creating with email", zap.String("email", adminEmail))

	// Hash da senha
	hashedPassword, err := encrypt.HashPassword(i.config.Admin.Password)
	if err != nil {
		i.logger.Error("[InsertAdminUser] Error hashing admin password", zap.Error(err))
		return err
	}

	now := utils.Now()
	user := &User{
		Name:     "Administrator",
		Email:    adminEmail,
		Admin:    true,
		Active:   true,
		Password: &hashedPassword,
		Source:   "LOCAL",
		Metadata: func() UserMetadata {
			m := NewDefaultMetadata()
			m.EmailVerified = true
			return m
		}(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := i.userService.Create(ctx, user); err != nil {
		i.logger.Error("[InsertAdminUser] Error creating admin user", zap.Error(err))
		return err
	}

	i.logger.Info("[InsertAdminUser] Administrator user created successfully!")
	return nil
}
