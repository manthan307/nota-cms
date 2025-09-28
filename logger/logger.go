package logger

import (
	"os"

	"go.uber.org/zap"
)

func InitLogger() *zap.Logger {
	cfg := zap.NewProductionConfig()
	if os.Getenv("ENV") == "development" {
		cfg = zap.NewDevelopmentConfig()
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return logger
}
