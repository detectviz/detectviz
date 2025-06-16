package otelzap

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"detectviz/pkg/platform/contracts"
)

// OtelZapLogger implements the contracts.Logger interface with OpenTelemetry integration.
// zh: OtelZapLogger 實作 contracts.Logger 介面並整合 OpenTelemetry。
type OtelZapLogger struct {
	zap    *zap.Logger
	sugar  *zap.SugaredLogger
	config *Config
	level  *zap.AtomicLevel
}

// NewOtelZapLogger creates a new OtelZap logger instance.
// zh: NewOtelZapLogger 建立新的 OtelZap 日誌記錄器實例。
func NewOtelZapLogger(config *Config) (*OtelZapLogger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Create atomic level for dynamic level changes
	level := zap.NewAtomicLevelAt(config.ParseLogLevel())

	// Build encoder config
	encoderConfig := buildEncoderConfig(config)

	// Build cores based on output configuration
	cores, err := buildCores(config, encoderConfig, level)
	if err != nil {
		return nil, err
	}

	// Create core
	var core zapcore.Core
	if len(cores) == 1 {
		core = cores[0]
	} else {
		core = zapcore.NewTee(cores...)
	}

	// Create logger with options
	logger := zap.New(core, buildLoggerOptions(config)...)
	sugar := logger.Sugar()

	return &OtelZapLogger{
		zap:    logger,
		sugar:  sugar,
		config: config,
		level:  &level,
	}, nil
}

// buildEncoderConfig builds the zapcore encoder configuration.
// zh: buildEncoderConfig 建構 zapcore 編碼器配置。
func buildEncoderConfig(config *Config) zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()

	// Configure time format
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Configure level
	encoderConfig.LevelKey = "level"
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	// Configure caller
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	// Configure message
	encoderConfig.MessageKey = "message"

	// Configure OpenTelemetry fields
	if config.OTEL != nil && config.OTEL.Enabled {
		// Trace ID and Span ID will be added by hooks
	}

	return encoderConfig
}

// buildCores builds zapcore cores based on configuration.
// zh: buildCores 根據配置建構 zapcore 核心。
func buildCores(config *Config, encoderConfig zapcore.EncoderConfig, level zapcore.LevelEnabler) ([]zapcore.Core, error) {
	var cores []zapcore.Core

	// Console output
	if config.OutputType == "console" || config.OutputType == "both" {
		var encoder zapcore.Encoder
		if config.Format == "json" {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		var writer zapcore.WriteSyncer
		switch config.Output {
		case "stderr":
			writer = zapcore.AddSync(os.Stderr)
		case "stdout", "":
			writer = zapcore.AddSync(os.Stdout)
		default:
			// Try to open as file
			file, err := os.OpenFile(config.Output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, err
			}
			writer = zapcore.AddSync(file)
		}

		cores = append(cores, zapcore.NewCore(encoder, writer, level))
	}

	// File output
	if config.OutputType == "file" || config.OutputType == "both" {
		if config.FileConfig == nil {
			return nil, fmt.Errorf("file_config is required for file output")
		}

		// Ensure directory exists
		dir := filepath.Dir(config.FileConfig.Filename)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory %s: %w", dir, err)
		}

		// Create lumberjack logger for rotation
		lumberjackLogger := &lumberjack.Logger{
			Filename:   config.FileConfig.Filename,
			MaxSize:    config.FileConfig.MaxSize,
			MaxBackups: config.FileConfig.MaxBackups,
			MaxAge:     config.FileConfig.MaxAge,
			Compress:   config.FileConfig.Compress,
		}

		var encoder zapcore.Encoder
		if config.Format == "json" {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		writer := zapcore.AddSync(lumberjackLogger)
		cores = append(cores, zapcore.NewCore(encoder, writer, level))
	}

	return cores, nil
}

// buildLoggerOptions builds zap logger options.
// zh: buildLoggerOptions 建構 zap 日誌記錄器選項。
func buildLoggerOptions(config *Config) []zap.Option {
	var options []zap.Option

	// Add caller information
	options = append(options, zap.AddCaller())
	options = append(options, zap.AddCallerSkip(1))

	// Add stack trace for error and fatal levels
	options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))

	// Add development mode if needed
	if config.Environment == "development" {
		options = append(options, zap.Development())
	}

	return options
}

// Debug implements contracts.Logger.
// zh: Debug 實作 contracts.Logger。
func (l *OtelZapLogger) Debug(msg string, fields ...interface{}) {
	l.logWithFields(zapcore.DebugLevel, msg, fields...)
}

// Info implements contracts.Logger.
// zh: Info 實作 contracts.Logger。
func (l *OtelZapLogger) Info(msg string, fields ...interface{}) {
	l.logWithFields(zapcore.InfoLevel, msg, fields...)
}

// Warn implements contracts.Logger.
// zh: Warn 實作 contracts.Logger。
func (l *OtelZapLogger) Warn(msg string, fields ...interface{}) {
	l.logWithFields(zapcore.WarnLevel, msg, fields...)
}

// Error implements contracts.Logger.
// zh: Error 實作 contracts.Logger。
func (l *OtelZapLogger) Error(msg string, fields ...interface{}) {
	l.logWithFields(zapcore.ErrorLevel, msg, fields...)
}

// Fatal implements contracts.Logger.
// zh: Fatal 實作 contracts.Logger。
func (l *OtelZapLogger) Fatal(msg string, fields ...interface{}) {
	l.logWithFields(zapcore.FatalLevel, msg, fields...)
	os.Exit(1)
}

// logWithFields is the internal method to log with structured fields.
// zh: logWithFields 是使用結構化欄位記錄的內部方法。
func (l *OtelZapLogger) logWithFields(level zapcore.Level, msg string, fields ...interface{}) {
	if !l.level.Enabled(level) {
		return
	}

	// Convert fields to zap fields
	zapFields := l.convertFields(fields...)

	// Log the message
	l.zap.Log(level, msg, zapFields...)
}

// convertFields converts interface{} fields to zap.Field.
// zh: convertFields 將 interface{} 欄位轉換為 zap.Field。
func (l *OtelZapLogger) convertFields(fields ...interface{}) []zap.Field {
	var zapFields []zap.Field

	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key := ""
			if keyStr, ok := fields[i].(string); ok {
				key = keyStr
			} else {
				key = fmt.Sprintf("%v", fields[i])
			}
			value := fields[i+1]
			zapFields = append(zapFields, zap.Any(key, value))
		}
	}

	return zapFields
}

// WithContext creates a new logger with trace context information.
// zh: WithContext 建立帶有追蹤上下文資訊的新日誌記錄器。
func (l *OtelZapLogger) WithContext(ctx context.Context) contracts.Logger {
	if ctx == nil {
		return l
	}

	// Extract trace information from context
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return l
	}

	spanContext := span.SpanContext()
	if !spanContext.IsValid() {
		return l
	}

	// Add trace fields to logger
	var fields []zap.Field
	if l.config.OTEL != nil && l.config.OTEL.Enabled {
		if l.config.OTEL.IncludeTrace {
			fields = append(fields, zap.String(l.config.OTEL.TraceIDField, spanContext.TraceID().String()))
			fields = append(fields, zap.String(l.config.OTEL.SpanIDField, spanContext.SpanID().String()))
		}
	}

	// Return a new logger with trace fields
	newLogger := l.zap.With(fields...)
	newSugar := newLogger.Sugar()

	return &OtelZapLogger{
		zap:    newLogger,
		sugar:  newSugar,
		config: l.config,
		level:  l.level,
	}
}

// Sync flushes any buffered log entries.
// zh: Sync 刷新任何緩衝的日誌條目。
func (l *OtelZapLogger) Sync() error {
	err := l.zap.Sync()
	// Ignore sync errors for stdout/stderr in test environments
	if err != nil && (err.Error() == "sync /dev/stdout: bad file descriptor" ||
		err.Error() == "sync /dev/stderr: bad file descriptor") {
		return nil
	}
	return err
}

// SetLevel dynamically changes the log level.
// zh: SetLevel 動態變更日誌等級。
func (l *OtelZapLogger) SetLevel(level string) error {
	config := &Config{Level: level}
	if !isValidLogLevel(config.Level) {
		return fmt.Errorf("invalid log level: %s", level)
	}

	zapLevel := config.ParseLogLevel()
	l.level.SetLevel(zapLevel)
	return nil
}

// Close closes the logger and releases resources.
// zh: Close 關閉日誌記錄器並釋放資源。
func (l *OtelZapLogger) Close() error {
	return l.Sync()
}
