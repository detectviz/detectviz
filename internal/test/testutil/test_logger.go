package testutil

import (
	"context"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// TestLogger 是一個可收集訊息的 logger 實作。
// zh: 測試用 logger，可用於單元測試驗證是否記錄特定訊息。
type TestLogger struct {
	mu       sync.Mutex
	messages []string
}

// NewTestLogger 回傳可收集訊息的測試 logger 實例。
// zh: 適用於需要驗證 log 行為的單元測試。
func NewTestLogger() *TestLogger {
	return &TestLogger{}
}

func (l *TestLogger) Debug(msg string, args ...any) {
	l.append(msg)
}
func (l *TestLogger) Info(msg string, args ...any) {
	l.append(msg)
}
func (l *TestLogger) Warn(msg string, args ...any) {
	l.append(msg)
}
func (l *TestLogger) Error(msg string, args ...any) {
	l.append(msg)
}
func (l *TestLogger) Named(name string) logger.Logger {
	return l
}
func (l *TestLogger) WithContext(ctx context.Context) logger.Logger {
	return l
}
func (l *TestLogger) WithFields(fields map[string]any) logger.Logger {
	return l
}
func (l *TestLogger) Sync() error {
	return nil
}

// Messages 回傳所有記錄的訊息。
// zh: 供測試時使用，驗證 logger 是否有記錄預期內容。
func (l *TestLogger) Messages() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.messages...)
}

func (l *TestLogger) append(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, msg)
}
