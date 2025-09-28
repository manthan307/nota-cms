package logger_test

import (
	"os"
	"testing"

	"github.com/manthan307/nota-cms/logger"
	"go.uber.org/zap"
)

func TestInitLogger(t *testing.T) {
	// Set environment variable for test
	os.Setenv("ENV", "development")

	logger := logger.InitLogger()
	defer logger.Sync() // sync when test finishes

	// Simple test: can write logs without panic
	logger.Info("Logger initialized successfully")

	if logger.Core().Enabled(zap.InfoLevel) != true {
		t.Errorf("Info level logging should be enabled")
	}
}
