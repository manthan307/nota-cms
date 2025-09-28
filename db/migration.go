package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func RunMigrations(lc fx.Lifecycle, logger *zap.Logger) {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbname := os.Getenv("POSTGRES_DB")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, dbname)

	m, err := migrate.New(
		"file://db/schema",
		dsn,
	)
	if err != nil {
		logger.Fatal("failed to create migrate instance", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Info("running migrations...")
			if err := m.Up(); err != nil && err != migrate.ErrNoChange {
				logger.Fatal("migration failed", zap.Error(err))
				return err
			}
			logger.Info("migrations applied ✅")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// close connections if needed
			return nil
		},
	})
}
