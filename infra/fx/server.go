package fx

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/lkgiovani/go-boilerplate/infra/config"
	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ServerParams struct {
	fx.In

	Config *config.Config
	Router *fiber.App
	Logger logger.Logger
}

var ServerModule = fx.Module("server",
	fx.Provide(
		newRouter,
	),
	fx.Invoke(
		StartServer,
	),
)

func StartServer(lc fx.Lifecycle, p ServerParams) {

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := p.Router.Listen(":" + fmt.Sprintf("%d", p.Config.Server.Port)); err != nil {
					p.Logger.Fatal("Failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			p.Logger.Info("Shutting down server...")
			return p.Router.Shutdown()
		},
	})
}

func newRouter(cfg *config.Config) *fiber.App {
	app := fiber.New()

	origins := cfg.Server.AllowedOrigins
	if origins == "" || origins == "*" {

		if cfg.Server.Mode == "development" {
			origins = "http://localhost:3000,http://localhost:5173,http://localhost:4000"
		} else {
			origins = ""
		}
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH, OPTIONS",
		AllowCredentials: true,
	}))

	return app
}
