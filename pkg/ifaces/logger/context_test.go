package logger_test

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

type testLogger struct {
	called bool
}

func (t *testLogger) Debug(msg string, args ...any)                 { t.called = true }
func (t *testLogger) Info(msg string, args ...any)                  { t.called = true }
func (t *testLogger) Warn(msg string, args ...any)                  { t.called = true }
func (t *testLogger) Error(msg string, args ...any)                 { t.called = true }
func (t *testLogger) Sync() error                                   { return nil }
func (t *testLogger) WithFields(map[string]any) logger.Logger       { return t }
func (t *testLogger) WithContext(ctx context.Context) logger.Logger { return t }
func (t *testLogger) Named(string) logger.Logger                    { return t }

func TestFromContext_Default(t *testing.T) {
	ctx := context.Background()
	l := logger.FromContext(ctx)

	if l == nil {
		t.Fatal("expected non-nil fallback logger")
	}
}

func TestWithContext_InjectAndRetrieve(t *testing.T) {
	ctx := context.Background()
	mock := &testLogger{}
	ctx = logger.WithContext(ctx, mock)

	l := logger.FromContext(ctx)
	l.Info("test")

	if !mock.called {
		t.Error("expected injected logger to be called")
	}
}
