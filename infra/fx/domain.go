package fx

import (
	"github.com/lkgiovani/go-boilerplate/internal/domain/auth"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/internal/security/jwt"
	"go.uber.org/fx"
)

var DomainModule = fx.Module(
	"domain",
	fx.Provide(
		user.NewUserRepository,
		user.NewService,
		user.NewInsertAdminUser,
		auth.NewAuthRepository,
		auth.NewService,
		jwt.NewJwtService,
	),
)
