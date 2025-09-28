package postgres

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
	"github.com/manthan307/nota-cms/logger"
)

func TestConnect(t *testing.T) {
	_ = godotenv.Load("../.env")

	logger := logger.InitLogger()
	pool := Connect(logger)
	defer pool.Close() // close after test

	ctx := context.Background()

	// simple ping test
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("ping failed: %v", err)
	}
}
