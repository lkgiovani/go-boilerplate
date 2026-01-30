package fx

import (
	"github.com/lkgiovani/go-boilerplate/infra/config"
	"github.com/lkgiovani/go-boilerplate/internal/domain/auth"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/internal/security/googleauth"
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
		fx.Annotate(
			func(cfg *config.Config) *googleauth.GoogleGateway {
				return googleauth.NewGoogleGateway(cfg.OAuth2.GoogleAndroidClientID, cfg.OAuth2.GoogleIosClientID)
			},
			fx.As(new(auth.GoogleTokenGateway)),
		),
	),
)
