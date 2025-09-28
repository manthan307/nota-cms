package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func Connect(logger *zap.Logger) *pgxpool.Pool {
	ctx := context.Background()

	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")

	if user == "" || pass == "" || dbName == "" || host == "" || port == "" {
		logger.Fatal("database environment variables not set")
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, dbName,
	)

	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		logger.Fatal("failed to parse pgxpool config", zap.Error(err))
	}

	cfg.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		logger.Fatal("unable to create pgxpool", zap.Error(err))
	}

	// Quick health check
	if err := pool.Ping(ctx); err != nil {
		logger.Fatal("cannot ping database", zap.Error(err))
	}

	logger.Info("connected to db ✅")
	return pool
}
