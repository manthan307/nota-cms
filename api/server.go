package api

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func StartServer(lc fx.Lifecycle, logger *zap.Logger, pool *pgxpool.Pool) *fiber.App {
	app := fiber.New(
		fiber.Config{
			IdleTimeout:           5 * time.Second,
			DisableStartupMessage: true,
			EnableIPValidation:    true,
		},
	)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen(":3000"); err != nil {
					logger.Error("Failed to start server", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Shutdown the Fiber app
			if err := app.Shutdown(); err != nil {
				logger.Error("Failed to shutdown server", zap.Error(err))
				return err
			}

			pool.Close()

			return nil
		},
	})

	return app
}
