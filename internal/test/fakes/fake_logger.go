package fakes

import (
	"context"
	"sync"

	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// LogEntry records a logging call for verification in tests.
// zh: 用於測試中紀錄 logger 輸出內容。
type LogEntry struct {
	Level string
	Msg   string
	Args  []any
}

// FakeLogger implements logger.Logger and stores log entries.
// zh: 簡易的 logger 假實作，會收集所有輸出訊息供檢查。
type FakeLogger struct {
	mu      sync.Mutex
	Entries []LogEntry
}

func (l *FakeLogger) record(level, msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Entries = append(l.Entries, LogEntry{Level: level, Msg: msg, Args: args})
}

func (l *FakeLogger) Debug(msg string, args ...any) { l.record("DEBUG", msg, args...) }
func (l *FakeLogger) Info(msg string, args ...any)  { l.record("INFO", msg, args...) }
func (l *FakeLogger) Warn(msg string, args ...any)  { l.record("WARN", msg, args...) }
func (l *FakeLogger) Error(msg string, args ...any) { l.record("ERROR", msg, args...) }

func (l *FakeLogger) Named(name string) ifacelogger.Logger                { return l }
func (l *FakeLogger) WithContext(ctx context.Context) ifacelogger.Logger  { return l }
func (l *FakeLogger) WithFields(fields map[string]any) ifacelogger.Logger { return l }
func (l *FakeLogger) Sync() error                                         { return nil }
