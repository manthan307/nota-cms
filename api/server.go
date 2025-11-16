package api

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func StartServer(lc fx.Lifecycle, logger *zap.Logger, pool *pgxpool.Pool) *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout:           5 * time.Second,
		DisableStartupMessage: true,
		EnableIPValidation:    true,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// Start server
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go app.Listen(":8000")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			pool.Close()
			return app.Shutdown()
		},
	})

	return app
}
