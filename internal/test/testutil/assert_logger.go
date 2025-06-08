package testutil

import (
	"context"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// LogEntry 表示一筆記錄的 log 項目。
// zh: 儲存 log 訊息與相關參數。
type LogEntry struct {
	Level string
	Msg   string
	Args  []any
}

// AssertLogger 是可驗證 log 輸出的 logger 實作。
// zh: 提供簡易記錄與查詢功能，方便在測試中斷言。
type AssertLogger struct {
	mu     sync.Mutex
	events []LogEntry
}

// NewAssertLogger 建立一個新的 AssertLogger 實例。
func NewAssertLogger() *AssertLogger {
	return &AssertLogger{}
}

// Entries 回傳所有記錄的 log 項目。
func (l *AssertLogger) Entries() []LogEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]LogEntry(nil), l.events...)
}

// 以下為 logger.Logger 實作

func (l *AssertLogger) record(level, msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = append(l.events, LogEntry{Level: level, Msg: msg, Args: args})
}

func (l *AssertLogger) Debug(msg string, args ...any) { l.record("DEBUG", msg, args...) }
func (l *AssertLogger) Info(msg string, args ...any)  { l.record("INFO", msg, args...) }
func (l *AssertLogger) Warn(msg string, args ...any)  { l.record("WARN", msg, args...) }
func (l *AssertLogger) Error(msg string, args ...any) { l.record("ERROR", msg, args...) }

func (l *AssertLogger) Named(name string) logger.Logger                { return l }
func (l *AssertLogger) WithContext(ctx context.Context) logger.Logger  { return l }
func (l *AssertLogger) WithFields(fields map[string]any) logger.Logger { return l }
func (l *AssertLogger) Sync() error                                    { return nil }
