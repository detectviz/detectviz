package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// LogLevel represents the logging level
// zh: LogLevel 代表日誌記錄等級
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// String returns the string representation of LogLevel
// zh: String 回傳 LogLevel 的字串表示
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

// ParseLogLevel parses string to LogLevel
// zh: ParseLogLevel 解析字串為 LogLevel
func ParseLogLevel(level string) LogLevel {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DebugLevel
	case "INFO":
		return InfoLevel
	case "WARN", "WARNING":
		return WarnLevel
	case "ERROR":
		return ErrorLevel
	case "FATAL":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// GlobalLogger holds the global logger instance
var globalLogger LoggerInterface
var loggerMutex sync.RWMutex

// LoggerInterface defines the interface for logger implementations
// zh: LoggerInterface 定義日誌記錄器實作的介面
type LoggerInterface interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Warn(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
	Fatal(msg string, fields ...interface{})
}

// Logger represents a logger instance (legacy implementation)
// zh: Logger 代表日誌記錄器實例（遺留實作）
type Logger struct {
	level  LogLevel
	output io.Writer
	mu     sync.Mutex
}

// Ensure Logger implements LoggerInterface
var _ LoggerInterface = (*Logger)(nil)

// LoggerConfig defines the configuration for the logger
// zh: LoggerConfig 定義日誌記錄器的配置
type LoggerConfig struct {
	Type       string      `yaml:"type" json:"type"`               // console, file, both
	Level      string      `yaml:"level" json:"level"`             // debug, info, warn, error
	Format     string      `yaml:"format" json:"format"`           // json, text
	Output     string      `yaml:"output" json:"output"`           // stdout, stderr, file path
	FileConfig *FileConfig `yaml:"file_config" json:"file_config"` // file rotation config
	OTEL       *OTELConfig `yaml:"otel" json:"otel"`               // OpenTelemetry config
}

// FileConfig defines file rotation configuration
// zh: FileConfig 定義檔案輪轉配置
type FileConfig struct {
	Filename   string `yaml:"filename" json:"filename"`       // log file path
	MaxSize    int    `yaml:"max_size" json:"max_size"`       // max size in MB (not implemented yet)
	MaxBackups int    `yaml:"max_backups" json:"max_backups"` // max backup files (not implemented yet)
	MaxAge     int    `yaml:"max_age" json:"max_age"`         // max age in days (not implemented yet)
	Compress   bool   `yaml:"compress" json:"compress"`       // compress old files (not implemented yet)
}

// OTELConfig defines OpenTelemetry configuration
// zh: OTELConfig 定義 OpenTelemetry 配置
type OTELConfig struct {
	Enabled        bool   `yaml:"enabled" json:"enabled"`
	ServiceName    string `yaml:"service_name" json:"service_name"`
	ServiceVersion string `yaml:"service_version" json:"service_version"`
	TraceIDField   string `yaml:"trace_id_field" json:"trace_id_field"`
	SpanIDField    string `yaml:"span_id_field" json:"span_id_field"`
}

// DefaultLoggerConfig returns the default logger configuration
// zh: DefaultLoggerConfig 回傳預設的日誌記錄器配置
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Type:   "console",
		Level:  "info",
		Format: "text",
		Output: "stdout",
		FileConfig: &FileConfig{
			Filename:   "logs/detectviz.log",
			MaxSize:    100, // 100MB
			MaxBackups: 3,
			MaxAge:     28, // 28 days
			Compress:   true,
		},
		OTEL: &OTELConfig{
			Enabled:        false, // Disabled for now
			ServiceName:    "detectviz",
			ServiceVersion: "1.0.0",
			TraceIDField:   "trace_id",
			SpanIDField:    "span_id",
		},
	}
}

// NewLogger creates a new logger with the given configuration
// zh: NewLogger 使用給定配置建立新的日誌記錄器
func NewLogger(config *LoggerConfig) (*Logger, error) {
	if config == nil {
		config = DefaultLoggerConfig()
	}

	level := ParseLogLevel(config.Level)

	var output io.Writer
	switch config.Type {
	case "console":
		output = getConsoleWriter(config.Output)
	case "file":
		if config.FileConfig == nil {
			return nil, fmt.Errorf("file_config is required for file output")
		}
		fileWriter, err := getFileWriter(config.FileConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create file writer: %w", err)
		}
		output = fileWriter
	case "both":
		consoleWriter := getConsoleWriter(config.Output)
		var fileWriter io.Writer
		if config.FileConfig != nil {
			fw, err := getFileWriter(config.FileConfig)
			if err != nil {
				return nil, fmt.Errorf("failed to create file writer: %w", err)
			}
			fileWriter = fw
		}
		if fileWriter != nil {
			output = io.MultiWriter(consoleWriter, fileWriter)
		} else {
			output = consoleWriter
		}
	default:
		output = os.Stdout
	}

	return &Logger{
		level:  level,
		output: output,
	}, nil
}

// getConsoleWriter returns console writer based on output config
// zh: getConsoleWriter 根據輸出配置回傳控制台寫入器
func getConsoleWriter(output string) io.Writer {
	switch output {
	case "stderr":
		return os.Stderr
	case "stdout", "":
		return os.Stdout
	default:
		// Try to open as file
		if file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err == nil {
			return file
		}
		return os.Stdout
	}
}

// getFileWriter returns file writer based on file config
// zh: getFileWriter 根據檔案配置回傳檔案寫入器
func getFileWriter(config *FileConfig) (io.Writer, error) {
	// Ensure directory exists
	dir := filepath.Dir(config.Filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory %s: %w", dir, err)
	}

	file, err := os.OpenFile(config.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file %s: %w", config.Filename, err)
	}

	return file, nil
}

// SetGlobalLogger sets the global logger instance
// zh: SetGlobalLogger 設置全域日誌記錄器實例
func SetGlobalLogger(logger LoggerInterface) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	globalLogger = logger
}

// SetGlobalLoggerLegacy sets the global logger instance (legacy API)
// zh: SetGlobalLoggerLegacy 設置全域日誌記錄器實例（遺留 API）
func SetGlobalLoggerLegacy(logger *Logger) {
	SetGlobalLogger(logger)
}

// GetGlobalLogger returns the global logger instance
// zh: GetGlobalLogger 回傳全域日誌記錄器實例
func GetGlobalLogger() LoggerInterface {
	loggerMutex.RLock()
	if globalLogger != nil {
		defer loggerMutex.RUnlock()
		return globalLogger
	}
	loggerMutex.RUnlock()

	// Initialize with default config if not set
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	if globalLogger == nil {
		config := DefaultLoggerConfig()
		logger, err := NewLogger(config)
		if err != nil {
			// Fallback to basic logger
			globalLogger = &Logger{
				level:  InfoLevel,
				output: os.Stdout,
			}
		} else {
			globalLogger = logger
		}
	}

	return globalLogger
}

// GetGlobalLoggerLegacy returns the global logger instance as *Logger (for backward compatibility)
// zh: GetGlobalLoggerLegacy 回傳全域日誌記錄器實例為 *Logger（向後相容）
func GetGlobalLoggerLegacy() *Logger {
	logger := GetGlobalLogger()
	if legacyLogger, ok := logger.(*Logger); ok {
		return legacyLogger
	}
	// If it's not a legacy logger, create a fallback
	return &Logger{
		level:  InfoLevel,
		output: os.Stdout,
	}
}

// L returns a logger (for future context support)
// zh: L 回傳日誌記錄器（為未來的 context 支援做準備）
func L(ctx context.Context) LoggerInterface {
	// For now, ignore context and return global logger
	// In the future, this can extract trace information from context
	return GetGlobalLogger()
}

// LLegacy returns a logger as *Logger (for backward compatibility)
// zh: LLegacy 回傳 *Logger 類型的日誌記錄器（向後相容）
func LLegacy(ctx context.Context) *Logger {
	return GetGlobalLoggerLegacy()
}

// log writes a log message with the given level
// zh: log 以給定等級寫入日誌訊息
func (l *Logger) log(level LogLevel, msg string, fields ...interface{}) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	var logMsg string
	if len(fields) > 0 {
		logMsg = fmt.Sprintf("[%s] %s: %s %v\n", timestamp, level.String(), msg, fields)
	} else {
		logMsg = fmt.Sprintf("[%s] %s: %s\n", timestamp, level.String(), msg)
	}

	fmt.Fprint(l.output, logMsg)
}

// Debug logs a debug message
// zh: Debug 記錄 debug 級別訊息
func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.log(DebugLevel, msg, fields...)
}

// Info logs an info message
// zh: Info 記錄 info 級別訊息
func (l *Logger) Info(msg string, fields ...interface{}) {
	l.log(InfoLevel, msg, fields...)
}

// Warn logs a warning message
// zh: Warn 記錄 warn 級別訊息
func (l *Logger) Warn(msg string, fields ...interface{}) {
	l.log(WarnLevel, msg, fields...)
}

// Error logs an error message
// zh: Error 記錄 error 級別訊息
func (l *Logger) Error(msg string, fields ...interface{}) {
	l.log(ErrorLevel, msg, fields...)
}

// Fatal logs a fatal message and exits
// zh: Fatal 記錄 fatal 級別訊息並退出
func (l *Logger) Fatal(msg string, fields ...interface{}) {
	l.log(FatalLevel, msg, fields...)
	os.Exit(1)
}

// Global logging functions
// zh: 全域日誌記錄函式

// Debug logs a debug message using global logger
// zh: Debug 使用全域日誌記錄器記錄 debug 級別訊息
func Debug(msg string, fields ...interface{}) {
	GetGlobalLogger().Debug(msg, fields...)
}

// Info logs an info message using global logger
// zh: Info 使用全域日誌記錄器記錄 info 級別訊息
func Info(msg string, fields ...interface{}) {
	GetGlobalLogger().Info(msg, fields...)
}

// Warn logs a warning message using global logger
// zh: Warn 使用全域日誌記錄器記錄 warn 級別訊息
func Warn(msg string, fields ...interface{}) {
	GetGlobalLogger().Warn(msg, fields...)
}

// Error logs an error message using global logger
// zh: Error 使用全域日誌記錄器記錄 error 級別訊息
func Error(msg string, fields ...interface{}) {
	GetGlobalLogger().Error(msg, fields...)
}

// Fatal logs a fatal message and exits using global logger
// zh: Fatal 使用全域日誌記錄器記錄 fatal 級別訊息並退出
func Fatal(msg string, fields ...interface{}) {
	GetGlobalLogger().Fatal(msg, fields...)
}

// Printf provides Printf-style logging for compatibility
// zh: Printf 提供 Printf 風格的日誌記錄以保持相容性
func Printf(format string, v ...interface{}) {
	GetGlobalLogger().Info(fmt.Sprintf(format, v...))
}

// Sync flushes any buffered log entries (no-op for current implementation)
// zh: Sync 刷新任何緩衝的日誌條目（目前實作中為空操作）
func Sync() error {
	// Current implementation uses direct writes, no buffering
	return nil
}

// Close closes any open log files (placeholder for future implementation)
// zh: Close 關閉任何開啟的日誌檔案（為未來實作預留的佔位符）
func Close() error {
	// Future implementation may need to close file handles
	return nil
}
