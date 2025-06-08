package logger_test

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/internal/adapters/logger"
	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestZapLogger_Info(t *testing.T) {
	core, observedLogs := observer.New(zapcore.InfoLevel)
	baseLogger := zap.New(core) // 正確建立 zap.Logger
	sugar := baseLogger.Sugar()

	zapLogger := logger.NewZapLogger(sugar)
	zapLogger.Info("test message", "key", "value")

	if observedLogs.Len() != 1 {
		t.Fatalf("expected 1 log entry, got %d", observedLogs.Len())
	}

	entry := observedLogs.All()[0]
	if entry.Message != "test message" {
		t.Errorf("unexpected log message: %s", entry.Message)
	}
	if val, ok := entry.ContextMap()["key"]; !ok || val != "value" {
		t.Errorf("missing or incorrect field value: %v", entry.ContextMap())
	}
}

func TestNopLogger_NoPanic(t *testing.T) {
	var log ifacelogger.Logger = logger.NewNopLogger()

	// All calls should be no-op
	log.Debug("debug")
	log.Info("info", "x", 1)
	log.Warn("warn")
	log.Error("error", "y", 2)
	log.Sync()
	log2 := log.WithFields(map[string]any{"test": true})
	log2 = log2.WithContext(context.Background())
	log2 = log2.Named("test")

	if log2 == nil {
		t.Error("NopLogger should return non-nil instance")
	}
}
