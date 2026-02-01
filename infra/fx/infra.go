package fx

import (
	"context"
	"time"

	"github.com/lkgiovani/go-boilerplate/infra/config"
	"github.com/lkgiovani/go-boilerplate/infra/database"
	"github.com/lkgiovani/go-boilerplate/internal/domain/user"
	"github.com/lkgiovani/go-boilerplate/pkg/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var infraModule = fx.Module(
	"infra",
	fx.Provide(
		NewDatabase,
	),
	fx.Invoke(
		InitializeAdminUser,
	),
)

func NewDatabase(lc fx.Lifecycle, cfg *config.Config, log logger.Logger) (*gorm.DB, error) {
	if err := database.Connect(); err != nil {
		log.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	db := database.GetDB()

	// Configure connection pool
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(cfg.Database.MinIdle)
		sqlDB.SetMaxOpenConns(cfg.Database.MaxPoolSize)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.MaxLifetime) * time.Millisecond)
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.Database.IdleTimeout) * time.Millisecond)
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Database connected successfully")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			sqlDB, err := db.DB()
			if err != nil {
				log.Error("Failed to get database instance", zap.Error(err))
				return err
			}
			log.Info("Closing database connection...")
			return sqlDB.Close()
		},
	})

	return db, nil
}

func InitializeAdminUser(lc fx.Lifecycle, insertAdminUser *user.InsertAdminUser, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Initializing admin user...")
			if err := insertAdminUser.Execute(ctx); err != nil {
				log.Error("Failed to initialize admin user", zap.Error(err))
				return err
			}
			return nil
		},
	})
}
