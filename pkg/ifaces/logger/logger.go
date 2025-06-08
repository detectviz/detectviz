package logger

import "context"

// Logger defines the abstract interface for structured logging in Detectviz.
// This interface is designed to support OpenTelemetry-compatible logging,
// allowing correlation with trace and metric data.
// zh: Logger 定義 Detectviz 中結構化日誌的抽象介面。此介面設計支援 OpenTelemetry 相容的日誌功能，允許與追蹤與指標數據做關聯。
type Logger interface {
	// Info logs a message at the info level with optional structured fields.
	// zh: 記錄 info 級別的日誌，可附帶結構化欄位。
	Info(msg string, fields ...any)

	// Warn logs a message at the warning level.
	// zh: 記錄 warning 級別的日誌。
	Warn(msg string, fields ...any)

	// Error logs a message at the error level.
	// zh: 記錄 error 級別的日誌。
	Error(msg string, fields ...any)

	// Debug logs a message at the debug level.
	// zh: 記錄 debug 級別的日誌。
	Debug(msg string, fields ...any)

	// WithFields returns a logger with the provided structured fields included.
	// zh: 回傳一個包含指定結構化欄位的新 logger 實例。
	WithFields(fields map[string]any) Logger

	// WithContext returns a logger that extracts trace context from the given context.
	// zh: 從傳入的 context 中提取追蹤資訊（如 trace_id、span_id），並回傳帶有 context 的 logger。
	WithContext(ctx context.Context) Logger

	// Named returns a logger with an assigned name, typically per module or component.
	// zh: 回傳一個命名的 logger，常用於模組或元件名稱區分。
	Named(name string) Logger

	// Sync flushes any buffered log entries, if supported.
	// zh: 若 logger 支援緩衝區，則強制將其清空寫出。
	Sync() error
}
