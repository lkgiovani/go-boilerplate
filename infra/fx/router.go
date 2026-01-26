package fx

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lkgiovani/go-boilerplate/internal/delivery"
	"github.com/lkgiovani/go-boilerplate/internal/domain/auth"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/internal/security/jwt"
	"go.uber.org/fx"
)

var RoutesModule = fx.Module("routes",
	fx.Provide(
		newHandler,
	),
)

func newHandler(
	userSvc *user.Service,
	jwtSvc *jwt.JwtService,
	authSvc *auth.Service,
	errHandler func(c *fiber.Ctx, err error) error,
) *delivery.Handler {
	return &delivery.Handler{
		UserService:  userSvc,
		JwtService:   jwtSvc,
		AuthService:  authSvc,
		ErrorHandler: errHandler,
	}
}
