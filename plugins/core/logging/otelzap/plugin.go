package otelzap

import (
	"context"
	"fmt"
	"time"

	"detectviz/pkg/platform/contracts"
	"detectviz/pkg/shared/log"
)

// OtelZapPlugin implements the contracts.LoggerPlugin interface.
// zh: OtelZapPlugin 實作 contracts.LoggerPlugin 介面。
type OtelZapPlugin struct {
	name        string
	version     string
	description string
	config      *Config
	logger      *OtelZapLogger
	initialized bool
	started     bool
}

// NewOtelZapPlugin creates a new OtelZap plugin instance.
// zh: NewOtelZapPlugin 建立新的 OtelZap 插件實例。
func NewOtelZapPlugin(config any) (contracts.Plugin, error) {
	cfg := DefaultConfig()

	// Parse configuration if provided
	if config != nil {
		if err := parseConfig(config, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse otelzap config: %w", err)
		}

		// Validate configuration immediately upon creation
		if err := cfg.Validate(); err != nil {
			return nil, fmt.Errorf("invalid otelzap configuration: %w", err)
		}
	}

	return &OtelZapPlugin{
		name:        "otelzap-logger",
		version:     "1.0.0",
		description: "OpenTelemetry integrated Zap logger plugin for DetectViz",
		config:      cfg,
		initialized: false,
		started:     false,
	}, nil
}

// parseConfig parses configuration from various formats.
// zh: parseConfig 從各種格式解析配置。
func parseConfig(config any, target *Config) error {
	switch c := config.(type) {
	case map[string]any:
		return parseMapConfig(c, target)
	case *Config:
		*target = *c
		return nil
	default:
		return fmt.Errorf("unsupported config type: %T", config)
	}
}

// parseMapConfig parses configuration from map format.
// zh: parseMapConfig 從 map 格式解析配置。
func parseMapConfig(configMap map[string]any, target *Config) error {
	if enabled, exists := configMap["enabled"]; exists {
		if boolVal, ok := enabled.(bool); ok {
			target.Enabled = boolVal
		}
	}

	if level, exists := configMap["level"]; exists {
		if str, ok := level.(string); ok {
			target.Level = str
		}
	}

	if format, exists := configMap["format"]; exists {
		if str, ok := format.(string); ok {
			target.Format = str
		}
	}

	if outputType, exists := configMap["output_type"]; exists {
		if str, ok := outputType.(string); ok {
			target.OutputType = str
		}
	}

	if output, exists := configMap["output"]; exists {
		if str, ok := output.(string); ok {
			target.Output = str
		}
	}

	if serviceName, exists := configMap["service_name"]; exists {
		if str, ok := serviceName.(string); ok {
			target.ServiceName = str
		}
	}

	if serviceVersion, exists := configMap["service_version"]; exists {
		if str, ok := serviceVersion.(string); ok {
			target.ServiceVersion = str
		}
	}

	if environment, exists := configMap["environment"]; exists {
		if str, ok := environment.(string); ok {
			target.Environment = str
		}
	}

	// Parse file_config
	if fileConfig, exists := configMap["file_config"]; exists {
		if fileConfigMap, ok := fileConfig.(map[string]any); ok {
			if target.FileConfig == nil {
				target.FileConfig = &FileConfig{}
			}
			parseFileConfig(fileConfigMap, target.FileConfig)
		}
	}

	// Parse otel config
	if otelConfig, exists := configMap["otel"]; exists {
		if otelConfigMap, ok := otelConfig.(map[string]any); ok {
			if target.OTEL == nil {
				target.OTEL = &OTELConfig{}
			}
			parseOTELConfig(otelConfigMap, target.OTEL)
		}
	}

	// Parse attributes
	if attributes, exists := configMap["attributes"]; exists {
		if attrMap, ok := attributes.(map[string]any); ok {
			if target.Attributes == nil {
				target.Attributes = make(map[string]string)
			}
			for k, v := range attrMap {
				if str, ok := v.(string); ok {
					target.Attributes[k] = str
				}
			}
		}
	}

	return nil
}

// parseFileConfig parses file configuration.
// zh: parseFileConfig 解析檔案配置。
func parseFileConfig(configMap map[string]any, target *FileConfig) {
	if filename, exists := configMap["filename"]; exists {
		if str, ok := filename.(string); ok {
			target.Filename = str
		}
	}

	if maxSize, exists := configMap["max_size"]; exists {
		if intVal, ok := maxSize.(int); ok {
			target.MaxSize = intVal
		}
	}

	if maxBackups, exists := configMap["max_backups"]; exists {
		if intVal, ok := maxBackups.(int); ok {
			target.MaxBackups = intVal
		}
	}

	if maxAge, exists := configMap["max_age"]; exists {
		if intVal, ok := maxAge.(int); ok {
			target.MaxAge = intVal
		}
	}

	if compress, exists := configMap["compress"]; exists {
		if boolVal, ok := compress.(bool); ok {
			target.Compress = boolVal
		}
	}
}

// parseOTELConfig parses OTEL configuration.
// zh: parseOTELConfig 解析 OTEL 配置。
func parseOTELConfig(configMap map[string]any, target *OTELConfig) {
	if enabled, exists := configMap["enabled"]; exists {
		if boolVal, ok := enabled.(bool); ok {
			target.Enabled = boolVal
		}
	}

	if traceIDField, exists := configMap["trace_id_field"]; exists {
		if str, ok := traceIDField.(string); ok {
			target.TraceIDField = str
		}
	}

	if spanIDField, exists := configMap["span_id_field"]; exists {
		if str, ok := spanIDField.(string); ok {
			target.SpanIDField = str
		}
	}

	if includeTrace, exists := configMap["include_trace"]; exists {
		if boolVal, ok := includeTrace.(bool); ok {
			target.IncludeTrace = boolVal
		}
	}

	if correlationID, exists := configMap["correlation_id"]; exists {
		if boolVal, ok := correlationID.(bool); ok {
			target.CorrelationID = boolVal
		}
	}
}

// Plugin interface implementation
// zh: Plugin 介面實作

// Name returns the plugin name.
// zh: Name 回傳插件名稱。
func (p *OtelZapPlugin) Name() string {
	return p.name
}

// Version returns the plugin version.
// zh: Version 回傳插件版本。
func (p *OtelZapPlugin) Version() string {
	return p.version
}

// Description returns the plugin description.
// zh: Description 回傳插件描述。
func (p *OtelZapPlugin) Description() string {
	return p.description
}

// Init initializes the plugin.
// zh: Init 初始化插件。
func (p *OtelZapPlugin) Init(config any) error {
	if p.initialized {
		return nil
	}

	// Parse configuration
	if config != nil {
		if err := parseConfig(config, p.config); err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}
	}

	// Validate configuration
	if err := p.config.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Create logger instance
	logger, err := NewOtelZapLogger(p.config)
	if err != nil {
		return fmt.Errorf("failed to create otelzap logger: %w", err)
	}

	p.logger = logger
	p.initialized = true

	// Set as global logger if enabled
	if p.config.Enabled {
		p.SetAsGlobalLogger()
	}

	// Log initialization
	ctx := context.Background()
	log.L(ctx).Info("OtelZap logger plugin initialized",
		"enabled", p.config.Enabled,
		"level", p.config.Level,
		"format", p.config.Format,
		"output_type", p.config.OutputType)

	return nil
}

// Shutdown shuts down the plugin.
// zh: Shutdown 關閉插件。
func (p *OtelZapPlugin) Shutdown() error {
	if p.logger != nil {
		if err := p.logger.Close(); err != nil {
			return fmt.Errorf("failed to close logger: %w", err)
		}
	}

	p.started = false
	p.initialized = false

	ctx := context.Background()
	log.L(ctx).Info("OtelZap logger plugin shutdown")

	return nil
}

// LoggerProvider interface implementation
// zh: LoggerProvider 介面實作

// Logger returns the logger instance.
// zh: Logger 回傳日誌記錄器實例。
func (p *OtelZapPlugin) Logger() contracts.Logger {
	if p.logger == nil {
		// Return a fallback logger if not initialized
		return log.GetGlobalLogger()
	}
	return p.logger
}

// WithContext returns a logger with context information.
// zh: WithContext 回傳帶有上下文資訊的日誌記錄器。
func (p *OtelZapPlugin) WithContext(ctx context.Context) contracts.Logger {
	if p.logger == nil {
		return log.L(ctx)
	}
	return p.logger.WithContext(ctx)
}

// Flush flushes any buffered log entries.
// zh: Flush 刷新任何緩衝的日誌條目。
func (p *OtelZapPlugin) Flush() error {
	if p.logger == nil {
		return nil
	}
	return p.logger.Sync()
}

// SetLevel dynamically changes the log level.
// zh: SetLevel 動態變更日誌等級。
func (p *OtelZapPlugin) SetLevel(level string) error {
	if p.logger == nil {
		return fmt.Errorf("logger not initialized")
	}
	return p.logger.SetLevel(level)
}

// Close closes the logger and releases resources.
// zh: Close 關閉日誌記錄器並釋放資源。
func (p *OtelZapPlugin) Close() error {
	if p.logger == nil {
		return nil
	}
	return p.logger.Close()
}

// LifecycleAware interface implementation
// zh: LifecycleAware 介面實作

// OnRegister is called when the plugin is registered.
// zh: OnRegister 在插件註冊時被呼叫。
func (p *OtelZapPlugin) OnRegister() error {
	ctx := context.Background()
	log.L(ctx).Info("OtelZap logger plugin registered")
	return nil
}

// OnStart is called when the plugin is started.
// zh: OnStart 在插件啟動時被呼叫。
func (p *OtelZapPlugin) OnStart() error {
	if !p.initialized {
		return fmt.Errorf("plugin not initialized")
	}

	if !p.config.Enabled {
		ctx := context.Background()
		log.L(ctx).Info("OtelZap logger plugin disabled, skipping start")
		return nil
	}

	p.started = true

	ctx := context.Background()
	log.L(ctx).Info("OtelZap logger plugin started",
		"service_name", p.config.ServiceName,
		"service_version", p.config.ServiceVersion,
		"environment", p.config.Environment)

	return nil
}

// OnStop is called when the plugin is stopped.
// zh: OnStop 在插件停止時被呼叫。
func (p *OtelZapPlugin) OnStop() error {
	if p.logger != nil {
		if err := p.logger.Sync(); err != nil {
			ctx := context.Background()
			log.L(ctx).Error("Failed to sync logger", "error", err)
		}
	}

	p.started = false

	ctx := context.Background()
	log.L(ctx).Info("OtelZap logger plugin stopped")

	return nil
}

// OnShutdown is called when the plugin is completely shut down.
// zh: OnShutdown 在插件完全關閉時被呼叫。
func (p *OtelZapPlugin) OnShutdown() error {
	ctx := context.Background()
	log.L(ctx).Info("OtelZap logger plugin shutdown complete")
	return nil
}

// HealthChecker interface implementation
// zh: HealthChecker 介面實作

// CheckHealth checks the health of the logger plugin.
// zh: CheckHealth 檢查日誌記錄器插件的健康狀態。
func (p *OtelZapPlugin) CheckHealth(ctx context.Context) contracts.HealthStatus {
	status := contracts.HealthStatus{
		Status:    "healthy",
		Message:   "OtelZap logger plugin is running normally",
		Timestamp: time.Now(),
		Details:   make(map[string]any),
	}

	// Check if plugin is initialized and started
	if !p.initialized {
		status.Status = "unhealthy"
		status.Message = "Plugin not initialized"
		return status
	}

	if p.config.Enabled && !p.started {
		status.Status = "degraded"
		status.Message = "Plugin initialized but not started"
	}

	// Add configuration details
	status.Details["enabled"] = p.config.Enabled
	status.Details["level"] = p.config.Level
	status.Details["format"] = p.config.Format
	status.Details["output_type"] = p.config.OutputType
	status.Details["service_name"] = p.config.ServiceName
	status.Details["service_version"] = p.config.ServiceVersion
	status.Details["environment"] = p.config.Environment

	if p.config.OTEL != nil {
		status.Details["otel_enabled"] = p.config.OTEL.Enabled
		status.Details["include_trace"] = p.config.OTEL.IncludeTrace
	}

	return status
}

// GetHealthMetrics returns health metrics for the logger plugin.
// zh: GetHealthMetrics 回傳日誌記錄器插件的健康指標。
func (p *OtelZapPlugin) GetHealthMetrics() map[string]any {
	metrics := make(map[string]any)

	metrics["plugin_initialized"] = p.initialized
	metrics["plugin_started"] = p.started
	metrics["plugin_enabled"] = p.config.Enabled

	// Configuration metrics
	metrics["log_level"] = p.config.Level
	metrics["log_format"] = p.config.Format
	metrics["output_type"] = p.config.OutputType

	// OTEL metrics
	if p.config.OTEL != nil {
		metrics["otel_enabled"] = p.config.OTEL.Enabled
		metrics["include_trace"] = p.config.OTEL.IncludeTrace
		metrics["correlation_id"] = p.config.OTEL.CorrelationID
	}

	return metrics
}

// Helper methods
// zh: 輔助方法

// IsEnabled returns whether the plugin is enabled.
// zh: IsEnabled 回傳插件是否已啟用。
func (p *OtelZapPlugin) IsEnabled() bool {
	return p.config.Enabled
}

// IsStarted returns whether the plugin is started.
// zh: IsStarted 回傳插件是否已啟動。
func (p *OtelZapPlugin) IsStarted() bool {
	return p.started
}

// GetConfig returns the current configuration.
// zh: GetConfig 回傳當前配置。
func (p *OtelZapPlugin) GetConfig() *Config {
	return p.config
}

// Register registers the OtelZap logger plugin.
// zh: Register 註冊 OtelZap 日誌記錄器插件。
func Register(registry contracts.Registry) error {
	return registry.RegisterPlugin("otelzap-logger", NewOtelZapPlugin)
}
