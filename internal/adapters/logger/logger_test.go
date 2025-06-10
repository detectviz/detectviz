package loggeradapter_test

import (
	"context"
	"testing"

	loggeradapter "github.com/detectviz/detectviz/internal/adapters/logger"
	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestZapLogger_Info(t *testing.T) {
	core, observedLogs := observer.New(zapcore.InfoLevel)
	baseLogger := zap.New(core) // zh: 正確建立 zap.Logger
	sugar := baseLogger.Sugar()

	zapLogger := loggeradapter.NewZapLogger(sugar)
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
	var log ifacelogger.Logger = loggeradapter.NewNopLogger()

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

func TestZapLogger_WithFieldsAndNamed(t *testing.T) {
	core, observedLogs := observer.New(zapcore.InfoLevel)
	base := zap.New(core).Sugar()
	log := loggeradapter.NewZapLogger(base)
	child := log.Named("child").WithFields(map[string]any{"k": "v"})
	child.Info("msg")
	entries := observedLogs.All()
	if len(entries) == 0 {
		t.Fatal("no log entry recorded")
	}
	entry := entries[len(entries)-1]
	if entry.Message != "msg" {
		t.Errorf("unexpected message: %s", entry.Message)
	}
	if entry.LoggerName != "child" {
		t.Errorf("expected logger name child, got %s", entry.LoggerName)
	}
}
