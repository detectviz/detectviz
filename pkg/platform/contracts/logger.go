package contracts

import (
	"context"
)

// LoggerProvider defines the interface for providing logger instances.
// zh: LoggerProvider 定義提供日誌記錄器實例的介面。
type LoggerProvider interface {
	// Logger returns the logger instance
	// zh: Logger 回傳日誌記錄器實例
	Logger() Logger

	// WithContext returns a logger with context information for trace injection
	// zh: WithContext 回傳帶有上下文資訊的日誌記錄器，用於 trace 注入
	WithContext(ctx context.Context) Logger

	// Flush flushes any buffered log entries
	// zh: Flush 刷新任何緩衝的日誌條目
	Flush() error

	// SetLevel dynamically changes the log level
	// zh: SetLevel 動態變更日誌等級
	SetLevel(level string) error

	// Close closes the logger and releases resources
	// zh: Close 關閉日誌記錄器並釋放資源
	Close() error
}

// Logger defines the basic logging interface that plugins can use.
// zh: Logger 定義插件可以使用的基本日誌記錄介面。
type Logger interface {
	// Debug logs a debug message
	// zh: Debug 記錄除錯訊息
	Debug(msg string, fields ...interface{})

	// Info logs an info message
	// zh: Info 記錄資訊訊息
	Info(msg string, fields ...interface{})

	// Warn logs a warning message
	// zh: Warn 記錄警告訊息
	Warn(msg string, fields ...interface{})

	// Error logs an error message
	// zh: Error 記錄錯誤訊息
	Error(msg string, fields ...interface{})

	// Fatal logs a fatal message and exits
	// zh: Fatal 記錄致命錯誤訊息並退出
	Fatal(msg string, fields ...interface{})
}

// LoggerPlugin defines the interface for logger plugins.
// zh: LoggerPlugin 定義日誌記錄器插件的介面。
type LoggerPlugin interface {
	Plugin
	LoggerProvider
	HealthChecker
	LifecycleAware
}
