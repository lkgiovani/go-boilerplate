package fx

import (
	"github.com/lkgiovani/go-boilerplate/infra/config"
	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"go.uber.org/fx"
)

var configModule = fx.Module(
	"config",
	fx.Provide(
		config.LoadConfig,
		provideLogger,
		provideJWTConfig,
	),
)

func provideJWTConfig(cfg *config.Config) config.JWTConfig {
	return cfg.JWT
}

func provideLogger(cfg *config.Config) (logger.Logger, error) {
	return logger.NewLogger(cfg.Server.Mode, cfg.Server.LogLevel)
}
