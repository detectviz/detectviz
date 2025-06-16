package otelzap

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

// Config defines the configuration for the OtelZap logger plugin.
// zh: Config 定義 OtelZap 日誌記錄器插件的配置。
type Config struct {
	Enabled        bool              `yaml:"enabled" json:"enabled"`
	Level          string            `yaml:"level" json:"level"`
	Format         string            `yaml:"format" json:"format"`           // json, text
	OutputType     string            `yaml:"output_type" json:"output_type"` // console, file, both
	Output         string            `yaml:"output" json:"output"`           // stdout, stderr, file path
	FileConfig     *FileConfig       `yaml:"file_config" json:"file_config"`
	OTEL           *OTELConfig       `yaml:"otel" json:"otel"`
	ServiceName    string            `yaml:"service_name" json:"service_name"`
	ServiceVersion string            `yaml:"service_version" json:"service_version"`
	Environment    string            `yaml:"environment" json:"environment"`
	Attributes     map[string]string `yaml:"attributes" json:"attributes"`
}

// FileConfig defines file rotation configuration.
// zh: FileConfig 定義檔案輪轉配置。
type FileConfig struct {
	Filename   string `yaml:"filename" json:"filename"`       // log file path
	MaxSize    int    `yaml:"max_size" json:"max_size"`       // max size in MB
	MaxBackups int    `yaml:"max_backups" json:"max_backups"` // max backup files
	MaxAge     int    `yaml:"max_age" json:"max_age"`         // max age in days
	Compress   bool   `yaml:"compress" json:"compress"`       // compress old files
}

// OTELConfig defines OpenTelemetry configuration.
// zh: OTELConfig 定義 OpenTelemetry 配置。
type OTELConfig struct {
	Enabled       bool   `yaml:"enabled" json:"enabled"`
	TraceIDField  string `yaml:"trace_id_field" json:"trace_id_field"`
	SpanIDField   string `yaml:"span_id_field" json:"span_id_field"`
	IncludeTrace  bool   `yaml:"include_trace" json:"include_trace"`
	CorrelationID bool   `yaml:"correlation_id" json:"correlation_id"`
}

// DefaultConfig returns the default configuration for the OtelZap plugin.
// zh: DefaultConfig 回傳 OtelZap 插件的預設配置。
func DefaultConfig() *Config {
	return &Config{
		Enabled:        true,
		Level:          "info",
		Format:         "json",
		OutputType:     "console",
		Output:         "stdout",
		ServiceName:    "detectviz",
		ServiceVersion: "1.0.0",
		Environment:    "development",
		FileConfig: &FileConfig{
			Filename:   "/var/log/detectviz/app.log",
			MaxSize:    100, // 100MB
			MaxBackups: 3,
			MaxAge:     30, // 30 days
			Compress:   true,
		},
		OTEL: &OTELConfig{
			Enabled:       true,
			TraceIDField:  "trace_id",
			SpanIDField:   "span_id",
			IncludeTrace:  true,
			CorrelationID: true,
		},
		Attributes: make(map[string]string),
	}
}

// Validate validates the configuration.
// zh: Validate 驗證配置。
func (c *Config) Validate() error {
	// Validate log level
	if !isValidLogLevel(c.Level) {
		return fmt.Errorf("invalid log level: %s", c.Level)
	}

	// Validate format
	if c.Format != "json" && c.Format != "text" {
		return fmt.Errorf("invalid format: %s, must be 'json' or 'text'", c.Format)
	}

	// Validate output type
	if c.OutputType != "console" && c.OutputType != "file" && c.OutputType != "both" {
		return fmt.Errorf("invalid output_type: %s, must be 'console', 'file', or 'both'", c.OutputType)
	}

	// Validate file configuration if needed
	if c.OutputType == "file" || c.OutputType == "both" {
		if c.FileConfig == nil {
			return fmt.Errorf("file_config is required when output_type is 'file' or 'both'")
		}
		if c.FileConfig.Filename == "" {
			return fmt.Errorf("filename is required in file_config")
		}
		if c.FileConfig.MaxSize <= 0 {
			return fmt.Errorf("max_size must be positive")
		}
		if c.FileConfig.MaxBackups < 0 {
			return fmt.Errorf("max_backups must be non-negative")
		}
		if c.FileConfig.MaxAge < 0 {
			return fmt.Errorf("max_age must be non-negative")
		}
	}

	// Validate OTEL configuration
	if c.OTEL != nil && c.OTEL.Enabled {
		if c.OTEL.TraceIDField == "" {
			c.OTEL.TraceIDField = "trace_id"
		}
		if c.OTEL.SpanIDField == "" {
			c.OTEL.SpanIDField = "span_id"
		}
	}

	return nil
}

// ParseLogLevel converts string level to zapcore.Level.
// zh: ParseLogLevel 將字串等級轉換為 zapcore.Level。
func (c *Config) ParseLogLevel() zapcore.Level {
	switch strings.ToLower(c.Level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// isValidLogLevel checks if the log level is valid.
// zh: isValidLogLevel 檢查日誌等級是否有效。
func isValidLogLevel(level string) bool {
	validLevels := []string{"debug", "info", "warn", "warning", "error", "fatal"}
	levelLower := strings.ToLower(level)
	for _, valid := range validLevels {
		if levelLower == valid {
			return true
		}
	}
	return false
}

// BuildResourceAttributes builds resource attributes for OpenTelemetry.
// zh: BuildResourceAttributes 為 OpenTelemetry 建構資源屬性。
func (c *Config) BuildResourceAttributes() map[string]string {
	attrs := make(map[string]string)

	// Add service information
	if c.ServiceName != "" {
		attrs["service.name"] = c.ServiceName
	}
	if c.ServiceVersion != "" {
		attrs["service.version"] = c.ServiceVersion
	}
	if c.Environment != "" {
		attrs["deployment.environment"] = c.Environment
	}

	// Add custom attributes
	for k, v := range c.Attributes {
		attrs[k] = v
	}

	// Add timestamp
	attrs["config.created_at"] = time.Now().Format(time.RFC3339)

	return attrs
}
