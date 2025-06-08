package scheduler

import (
	"context"

	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// testLogger implements logger.Logger for unit testing.
// zh: 測試用 logger 實作，符合 logger.Logger 介面。
type testLogger struct{}

// Debug logs debug messages (no-op).
// zh: 輸出除錯訊息（測試中無行為）。
func (l *testLogger) Debug(msg string, args ...any) {}

// Info logs informational messages (no-op).
// zh: 輸出一般訊息（測試中無行為）。
func (l *testLogger) Info(msg string, args ...any) {}

// Warn logs warning messages (no-op).
// zh: 輸出警告訊息（測試中無行為）。
func (l *testLogger) Warn(msg string, args ...any) {}

// Error logs error messages (no-op).
// zh: 輸出錯誤訊息（測試中無行為）。
func (l *testLogger) Error(msg string, args ...any) {}

// Named returns a child logger with name (no-op).
// zh: 回傳具名 logger（測試中無行為）。
func (l *testLogger) Named(name string) ifacelogger.Logger {
	return l
}

// WithFields returns a logger with fields (no-op).
// zh: 回傳附帶欄位的 logger（測試中無行為）。
func (l *testLogger) WithFields(fields map[string]interface{}) ifacelogger.Logger {
	return l
}

// WithContext returns a logger tied to a context (no-op).
// zh: 回傳與 context 綁定的 logger（測試中無行為）。
func (l *testLogger) WithContext(ctx context.Context) ifacelogger.Logger {
	return l
}

// Sync flushes logs (no-op).
// zh: 強制 flush 日誌（測試中無行為）。
func (l *testLogger) Sync() error {
	return nil
}
