package logger

import "context"

// NopLogger implements Logger with no-op methods.
// zh: NopLogger 是 Logger 的空實作，不會輸出任何日誌。
type NopLogger struct{}

// Debug does nothing.
// zh: Debug 級別，不輸出內容。
func (NopLogger) Debug(msg string, args ...any) {}

// Sync does nothing and returns nil.
// zh: 不需 flush，直接回傳 nil。
func (NopLogger) Sync() error { return nil }

// WithFields returns NopLogger.
// zh: 回傳自身，不套用任何欄位。
func (NopLogger) WithFields(fields map[string]any) Logger {
	return NopLogger{}
}

// WithContext returns NopLogger.
// zh: 回傳自身，不注入 context。
func (NopLogger) WithContext(ctx context.Context) Logger {
	return NopLogger{}
}

// Named returns NopLogger.
// zh: 回傳自身，不套用 logger 名稱。
func (NopLogger) Named(name string) Logger {
	return NopLogger{}
}

// Error does nothing.
// zh: Error 級別，不輸出內容。
func (NopLogger) Error(msg string, args ...any) {}

// Info does nothing.
// zh: Info 級別，不輸出內容。
func (NopLogger) Info(msg string, args ...any) {}

// Warn does nothing.
// zh: Warn 級別，不輸出內容。
func (NopLogger) Warn(msg string, args ...any) {}
