package logger

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// NewNopLogger returns a no-op logger implementation.
// zh: 回傳一個不會輸出任何內容的 logger，常用於測試或作為預設 fallback。
func NewNopLogger() logger.Logger {
	return &NopLogger{}
}

// NopLogger 實作 logger.Logger 但不輸出任何 log。
// zh: 此實作常用於測試或禁用 log 時。
type NopLogger struct{}

func (n *NopLogger) Info(msg string, args ...any)  {}
func (n *NopLogger) Warn(msg string, args ...any)  {}
func (n *NopLogger) Error(msg string, args ...any) {}
func (n *NopLogger) Debug(msg string, args ...any) {}

func (n *NopLogger) Named(name string) logger.Logger {
	return n
}

func (n *NopLogger) WithContext(ctx context.Context) logger.Logger {
	return n
}

func (n *NopLogger) WithFields(fields map[string]any) logger.Logger {
	return n
}

func (n *NopLogger) Sync() error {
	return nil
}
