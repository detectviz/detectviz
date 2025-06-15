package fake

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// FakeLogger is a mock implementation of logger for testing.
// zh: FakeLogger 是用於測試的日誌記錄器模擬實作。
type FakeLogger struct {
	mu       sync.Mutex
	entries  []LogEntry
	level    LogLevel
	output   io.Writer
	traceID  string
	spanID   string
	disabled bool
}

// LogEntry represents a single log entry.
// zh: LogEntry 代表單一日誌條目。
type LogEntry struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     LogLevel               `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields"`
	TraceID   string                 `json:"trace_id,omitempty"`
	SpanID    string                 `json:"span_id,omitempty"`
}

// LogLevel represents the logging level for fake logger.
// zh: LogLevel 代表假日誌記錄器的日誌等級。
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// String returns the string representation of LogLevel.
func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case FatalLevel:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// NewFakeLogger creates a new fake logger instance.
// zh: NewFakeLogger 建立新的假日誌記錄器實例。
func NewFakeLogger() *FakeLogger {
	return &FakeLogger{
		entries:  make([]LogEntry, 0),
		level:    InfoLevel,
		output:   os.Stdout,
		disabled: false,
	}
}

// NewFakeLoggerWithTracing creates a new fake logger with tracing information.
// zh: NewFakeLoggerWithTracing 建立帶有追蹤資訊的新假日誌記錄器。
func NewFakeLoggerWithTracing(traceID, spanID string) *FakeLogger {
	return &FakeLogger{
		entries:  make([]LogEntry, 0),
		level:    InfoLevel,
		output:   os.Stdout,
		traceID:  traceID,
		spanID:   spanID,
		disabled: false,
	}
}

// SetLevel sets the logging level.
// zh: SetLevel 設置日誌記錄等級。
func (fl *FakeLogger) SetLevel(level LogLevel) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.level = level
}

// SetOutput sets the output writer.
// zh: SetOutput 設置輸出寫入器。
func (fl *FakeLogger) SetOutput(output io.Writer) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.output = output
}

// SetTracing sets the tracing information.
// zh: SetTracing 設置追蹤資訊。
func (fl *FakeLogger) SetTracing(traceID, spanID string) {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.traceID = traceID
	fl.spanID = spanID
}

// Disable disables logging output.
// zh: Disable 禁用日誌輸出。
func (fl *FakeLogger) Disable() {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.disabled = true
}

// Enable enables logging output.
// zh: Enable 啟用日誌輸出。
func (fl *FakeLogger) Enable() {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.disabled = false
}

// log is the internal logging method.
// zh: log 是內部日誌記錄方法。
func (fl *FakeLogger) log(level LogLevel, msg string, fields ...interface{}) {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	// Check if logging is disabled or level is below threshold
	if fl.disabled || level < fl.level {
		return
	}

	// Parse fields into map
	fieldMap := make(map[string]interface{})
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key := fmt.Sprintf("%v", fields[i])
			value := fields[i+1]
			fieldMap[key] = value
		}
	}

	// Create log entry
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    fieldMap,
		TraceID:   fl.traceID,
		SpanID:    fl.spanID,
	}

	// Store entry
	fl.entries = append(fl.entries, entry)

	// Write to output if not disabled
	if !fl.disabled && fl.output != nil {
		fl.writeEntry(entry)
	}
}

// writeEntry writes a log entry to the output.
// zh: writeEntry 將日誌條目寫入輸出。
func (fl *FakeLogger) writeEntry(entry LogEntry) {
	var parts []string

	// Add timestamp
	parts = append(parts, fmt.Sprintf("[%s]", entry.Timestamp.Format("2006-01-02 15:04:05.000")))

	// Add level
	parts = append(parts, entry.Level.String()+":")

	// Add trace information if available
	if entry.TraceID != "" {
		parts = append(parts, fmt.Sprintf("trace_id=%s", entry.TraceID))
	}
	if entry.SpanID != "" {
		parts = append(parts, fmt.Sprintf("span_id=%s", entry.SpanID))
	}

	// Add message
	parts = append(parts, entry.Message)

	// Add fields
	for key, value := range entry.Fields {
		parts = append(parts, fmt.Sprintf("%s=%v", key, value))
	}

	// Write to output
	line := strings.Join(parts, " ") + "\n"
	fl.output.Write([]byte(line))
}

// Debug logs a debug message.
// zh: Debug 記錄除錯訊息。
func (fl *FakeLogger) Debug(msg string, fields ...interface{}) {
	fl.log(DebugLevel, msg, fields...)
}

// Info logs an info message.
// zh: Info 記錄資訊訊息。
func (fl *FakeLogger) Info(msg string, fields ...interface{}) {
	fl.log(InfoLevel, msg, fields...)
}

// Warn logs a warning message.
// zh: Warn 記錄警告訊息。
func (fl *FakeLogger) Warn(msg string, fields ...interface{}) {
	fl.log(WarnLevel, msg, fields...)
}

// Error logs an error message.
// zh: Error 記錄錯誤訊息。
func (fl *FakeLogger) Error(msg string, fields ...interface{}) {
	fl.log(ErrorLevel, msg, fields...)
}

// Fatal logs a fatal message.
// zh: Fatal 記錄致命錯誤訊息。
func (fl *FakeLogger) Fatal(msg string, fields ...interface{}) {
	fl.log(FatalLevel, msg, fields...)
}

// GetEntries returns all logged entries.
// zh: GetEntries 返回所有記錄的條目。
func (fl *FakeLogger) GetEntries() []LogEntry {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	entries := make([]LogEntry, len(fl.entries))
	copy(entries, fl.entries)
	return entries
}

// GetEntriesByLevel returns entries filtered by level.
// zh: GetEntriesByLevel 返回按等級過濾的條目。
func (fl *FakeLogger) GetEntriesByLevel(level LogLevel) []LogEntry {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	var filtered []LogEntry
	for _, entry := range fl.entries {
		if entry.Level == level {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// GetEntriesWithTrace returns entries that have trace information.
// zh: GetEntriesWithTrace 返回包含追蹤資訊的條目。
func (fl *FakeLogger) GetEntriesWithTrace() []LogEntry {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	var filtered []LogEntry
	for _, entry := range fl.entries {
		if entry.TraceID != "" || entry.SpanID != "" {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// CountEntries returns the total number of logged entries.
// zh: CountEntries 返回記錄條目的總數。
func (fl *FakeLogger) CountEntries() int {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	return len(fl.entries)
}

// CountEntriesByLevel returns the number of entries for a specific level.
// zh: CountEntriesByLevel 返回特定等級的條目數量。
func (fl *FakeLogger) CountEntriesByLevel(level LogLevel) int {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	count := 0
	for _, entry := range fl.entries {
		if entry.Level == level {
			count++
		}
	}
	return count
}

// HasEntry checks if an entry with the given message exists.
// zh: HasEntry 檢查是否存在具有給定訊息的條目。
func (fl *FakeLogger) HasEntry(message string) bool {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	for _, entry := range fl.entries {
		if entry.Message == message {
			return true
		}
	}
	return false
}

// HasEntryWithField checks if an entry with the given field exists.
// zh: HasEntryWithField 檢查是否存在具有給定欄位的條目。
func (fl *FakeLogger) HasEntryWithField(key string, value interface{}) bool {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	for _, entry := range fl.entries {
		if fieldValue, exists := entry.Fields[key]; exists && fieldValue == value {
			return true
		}
	}
	return false
}

// Clear clears all logged entries.
// zh: Clear 清除所有記錄的條目。
func (fl *FakeLogger) Clear() {
	fl.mu.Lock()
	defer fl.mu.Unlock()
	fl.entries = make([]LogEntry, 0)
}

// Reset resets the logger to initial state.
// zh: Reset 重置日誌記錄器到初始狀態。
func (fl *FakeLogger) Reset() {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	fl.entries = make([]LogEntry, 0)
	fl.level = InfoLevel
	fl.output = os.Stdout
	fl.traceID = ""
	fl.spanID = ""
	fl.disabled = false
}

// FakeLoggerFromContext creates a fake logger that extracts trace information from context.
// zh: FakeLoggerFromContext 建立從上下文提取追蹤資訊的假日誌記錄器。
func FakeLoggerFromContext(ctx context.Context) *FakeLogger {
	logger := NewFakeLogger()

	// Extract trace information from context if available
	if traceID := ctx.Value("trace_id"); traceID != nil {
		if traceIDStr, ok := traceID.(string); ok {
			logger.traceID = traceIDStr
		}
	}

	if spanID := ctx.Value("span_id"); spanID != nil {
		if spanIDStr, ok := spanID.(string); ok {
			logger.spanID = spanIDStr
		}
	}

	return logger
}

// WithTracing returns a new fake logger with the specified tracing information.
// The new logger shares the same entries slice as the original logger.
// zh: WithTracing 返回帶有指定追蹤資訊的新假日誌記錄器。
// 新的日誌記錄器與原始記錄器共享相同的條目切片。
func (fl *FakeLogger) WithTracing(traceID, spanID string) *FakeLogger {
	newLogger := &FakeLogger{
		entries:  fl.entries, // Share the same entries slice
		level:    fl.level,
		output:   fl.output,
		traceID:  traceID,
		spanID:   spanID,
		disabled: fl.disabled,
		mu:       fl.mu, // Share the same mutex
	}
	return newLogger
}

// GetStats returns statistics about the logged entries.
// zh: GetStats 返回記錄條目的統計資訊。
func (fl *FakeLogger) GetStats() map[string]int {
	fl.mu.Lock()
	defer fl.mu.Unlock()

	stats := map[string]int{
		"total":      len(fl.entries),
		"debug":      0,
		"info":       0,
		"warn":       0,
		"error":      0,
		"fatal":      0,
		"with_trace": 0,
	}

	for _, entry := range fl.entries {
		switch entry.Level {
		case DebugLevel:
			stats["debug"]++
		case InfoLevel:
			stats["info"]++
		case WarnLevel:
			stats["warn"]++
		case ErrorLevel:
			stats["error"]++
		case FatalLevel:
			stats["fatal"]++
		}

		if entry.TraceID != "" || entry.SpanID != "" {
			stats["with_trace"]++
		}
	}

	return stats
}
