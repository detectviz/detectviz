package sdkwrapper

import (
	"context"
	"fmt"
	"time"

	"detectviz/pkg/platform/contracts"
	"detectviz/pkg/shared/log"
)

// SDKWrapperPlugin provides OpenTelemetry SDK integration for DetectViz.
// zh: SDKWrapperPlugin 為 DetectViz 提供 OpenTelemetry SDK 整合。
type SDKWrapperPlugin struct {
	name        string
	version     string
	description string
	config      *Config
	initialized bool
	started     bool
}

// Config contains configuration for the OpenTelemetry SDK wrapper.
// zh: Config 包含 OpenTelemetry SDK 包裝器的配置。
type Config struct {
	ServiceName    string            `yaml:"service_name" json:"service_name"`
	ServiceVersion string            `yaml:"service_version" json:"service_version"`
	Environment    string            `yaml:"environment" json:"environment"`
	Enabled        bool              `yaml:"enabled" json:"enabled"`
	Tracing        TracingConfig     `yaml:"tracing" json:"tracing"`
	Metrics        MetricsConfig     `yaml:"metrics" json:"metrics"`
	Logging        LoggingConfig     `yaml:"logging" json:"logging"`
	Resource       ResourceConfig    `yaml:"resource" json:"resource"`
	Exporters      ExportersConfig   `yaml:"exporters" json:"exporters"`
	Attributes     map[string]string `yaml:"attributes" json:"attributes"`
}

// TracingConfig contains tracing-specific configuration.
// zh: TracingConfig 包含追蹤特定的配置。
type TracingConfig struct {
	Enabled     bool    `yaml:"enabled" json:"enabled"`
	SampleRatio float64 `yaml:"sample_ratio" json:"sample_ratio"`
	MaxSpans    int     `yaml:"max_spans" json:"max_spans"`
	BatchSize   int     `yaml:"batch_size" json:"batch_size"`
	Timeout     string  `yaml:"timeout" json:"timeout"`
}

// MetricsConfig contains metrics-specific configuration.
// zh: MetricsConfig 包含指標特定的配置。
type MetricsConfig struct {
	Enabled        bool   `yaml:"enabled" json:"enabled"`
	CollectRuntime bool   `yaml:"collect_runtime" json:"collect_runtime"`
	CollectHost    bool   `yaml:"collect_host" json:"collect_host"`
	Interval       string `yaml:"interval" json:"interval"`
}

// LoggingConfig contains logging-specific configuration.
// zh: LoggingConfig 包含日誌特定的配置。
type LoggingConfig struct {
	Enabled       bool   `yaml:"enabled" json:"enabled"`
	IncludeTrace  bool   `yaml:"include_trace" json:"include_trace"`
	LogLevel      string `yaml:"log_level" json:"log_level"`
	CorrelationID bool   `yaml:"correlation_id" json:"correlation_id"`
}

// ResourceConfig contains resource attribute configuration.
// zh: ResourceConfig 包含資源屬性配置。
type ResourceConfig struct {
	DetectHost    bool              `yaml:"detect_host" json:"detect_host"`
	DetectProcess bool              `yaml:"detect_process" json:"detect_process"`
	DetectRuntime bool              `yaml:"detect_runtime" json:"detect_runtime"`
	CustomAttrs   map[string]string `yaml:"custom_attrs" json:"custom_attrs"`
	Organization  string            `yaml:"organization" json:"organization"`
	Team          string            `yaml:"team" json:"team"`
	DeploymentEnv string            `yaml:"deployment_env" json:"deployment_env"`
}

// ExportersConfig contains exporter configuration.
// zh: ExportersConfig 包含匯出器配置。
type ExportersConfig struct {
	OTLP       OTLPConfig       `yaml:"otlp" json:"otlp"`
	Jaeger     JaegerConfig     `yaml:"jaeger" json:"jaeger"`
	Prometheus PrometheusConfig `yaml:"prometheus" json:"prometheus"`
	Console    ConsoleConfig    `yaml:"console" json:"console"`
}

// OTLPConfig contains OTLP exporter configuration.
// zh: OTLPConfig 包含 OTLP 匯出器配置。
type OTLPConfig struct {
	Enabled     bool              `yaml:"enabled" json:"enabled"`
	Endpoint    string            `yaml:"endpoint" json:"endpoint"`
	Insecure    bool              `yaml:"insecure" json:"insecure"`
	Headers     map[string]string `yaml:"headers" json:"headers"`
	Compression string            `yaml:"compression" json:"compression"`
	Timeout     string            `yaml:"timeout" json:"timeout"`
}

// JaegerConfig contains Jaeger exporter configuration.
// zh: JaegerConfig 包含 Jaeger 匯出器配置。
type JaegerConfig struct {
	Enabled  bool   `yaml:"enabled" json:"enabled"`
	Endpoint string `yaml:"endpoint" json:"endpoint"`
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}

// PrometheusConfig contains Prometheus exporter configuration.
// zh: PrometheusConfig 包含 Prometheus 匯出器配置。
type PrometheusConfig struct {
	Enabled bool   `yaml:"enabled" json:"enabled"`
	Host    string `yaml:"host" json:"host"`
	Port    int    `yaml:"port" json:"port"`
	Path    string `yaml:"path" json:"path"`
}

// ConsoleConfig contains console exporter configuration.
// zh: ConsoleConfig 包含控制台匯出器配置。
type ConsoleConfig struct {
	Enabled bool `yaml:"enabled" json:"enabled"`
	Pretty  bool `yaml:"pretty" json:"pretty"`
}

// NewSDKWrapperPlugin creates a new OpenTelemetry SDK wrapper plugin instance.
// zh: NewSDKWrapperPlugin 建立新的 OpenTelemetry SDK 包裝器插件實例。
func NewSDKWrapperPlugin(config any) (contracts.Plugin, error) {
	var cfg Config

	// Set default configuration
	cfg = Config{
		ServiceName:    "detectviz",
		ServiceVersion: "1.0.0",
		Environment:    "development",
		Enabled:        true,
		Tracing: TracingConfig{
			Enabled:     true,
			SampleRatio: 1.0,
			MaxSpans:    1000,
			BatchSize:   100,
			Timeout:     "30s",
		},
		Metrics: MetricsConfig{
			Enabled:        true,
			CollectRuntime: true,
			CollectHost:    true,
			Interval:       "15s",
		},
		Logging: LoggingConfig{
			Enabled:       true,
			IncludeTrace:  true,
			LogLevel:      "info",
			CorrelationID: true,
		},
		Resource: ResourceConfig{
			DetectHost:    true,
			DetectProcess: true,
			DetectRuntime: true,
			CustomAttrs:   make(map[string]string),
			Organization:  "detectviz",
			Team:          "platform",
			DeploymentEnv: "development",
		},
		Exporters: ExportersConfig{
			OTLP: OTLPConfig{
				Enabled:     true,
				Endpoint:    "http://localhost:4317",
				Insecure:    true,
				Headers:     make(map[string]string),
				Compression: "gzip",
				Timeout:     "10s",
			},
			Jaeger: JaegerConfig{
				Enabled:  false,
				Endpoint: "http://localhost:14268/api/traces",
			},
			Prometheus: PrometheusConfig{
				Enabled: true,
				Host:    "localhost",
				Port:    8080,
				Path:    "/metrics",
			},
			Console: ConsoleConfig{
				Enabled: false,
				Pretty:  true,
			},
		},
		Attributes: make(map[string]string),
	}

	// Parse configuration if provided
	if config != nil {
		if err := parseSDKConfig(config, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse SDK wrapper config: %w", err)
		}
	}

	return &SDKWrapperPlugin{
		name:        "otel-sdk-wrapper",
		version:     "1.0.0",
		description: "OpenTelemetry SDK wrapper for unified observability",
		config:      &cfg,
		initialized: false,
		started:     false,
	}, nil
}

// parseSDKConfig parses configuration from various formats.
// zh: parseSDKConfig 從各種格式解析配置。
func parseSDKConfig(config any, target *Config) error {
	configMap, ok := config.(map[string]any)
	if !ok {
		return fmt.Errorf("config must be a map")
	}

	// Parse basic fields
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
	if enabled, exists := configMap["enabled"]; exists {
		if boolVal, ok := enabled.(bool); ok {
			target.Enabled = boolVal
		}
	}

	// Parse tracing config
	if tracing, exists := configMap["tracing"]; exists {
		if tracingMap, ok := tracing.(map[string]any); ok {
			parseTracingConfig(tracingMap, &target.Tracing)
		}
	}

	// Parse metrics config
	if metrics, exists := configMap["metrics"]; exists {
		if metricsMap, ok := metrics.(map[string]any); ok {
			parseMetricsConfig(metricsMap, &target.Metrics)
		}
	}

	// Parse logging config
	if logging, exists := configMap["logging"]; exists {
		if loggingMap, ok := logging.(map[string]any); ok {
			parseLoggingConfig(loggingMap, &target.Logging)
		}
	}

	// Parse resource config
	if resource, exists := configMap["resource"]; exists {
		if resourceMap, ok := resource.(map[string]any); ok {
			parseResourceConfig(resourceMap, &target.Resource)
		}
	}

	// Parse exporters config
	if exporters, exists := configMap["exporters"]; exists {
		if exportersMap, ok := exporters.(map[string]any); ok {
			parseExportersConfig(exportersMap, &target.Exporters)
		}
	}

	// Parse attributes
	if attributes, exists := configMap["attributes"]; exists {
		if attrMap, ok := attributes.(map[string]any); ok {
			strAttrs := make(map[string]string)
			for k, v := range attrMap {
				if str, ok := v.(string); ok {
					strAttrs[k] = str
				}
			}
			target.Attributes = strAttrs
		}
	}

	return nil
}

// Helper parsing functions
// zh: 輔助解析函式

func parseTracingConfig(configMap map[string]any, target *TracingConfig) {
	if enabled, exists := configMap["enabled"]; exists {
		if boolVal, ok := enabled.(bool); ok {
			target.Enabled = boolVal
		}
	}
	if sampleRatio, exists := configMap["sample_ratio"]; exists {
		if floatVal, ok := sampleRatio.(float64); ok {
			target.SampleRatio = floatVal
		}
	}
	if maxSpans, exists := configMap["max_spans"]; exists {
		if intVal, ok := maxSpans.(int); ok {
			target.MaxSpans = intVal
		}
	}
	if batchSize, exists := configMap["batch_size"]; exists {
		if intVal, ok := batchSize.(int); ok {
			target.BatchSize = intVal
		}
	}
	if timeout, exists := configMap["timeout"]; exists {
		if str, ok := timeout.(string); ok {
			target.Timeout = str
		}
	}
}

func parseMetricsConfig(configMap map[string]any, target *MetricsConfig) {
	if enabled, exists := configMap["enabled"]; exists {
		if boolVal, ok := enabled.(bool); ok {
			target.Enabled = boolVal
		}
	}
	if collectRuntime, exists := configMap["collect_runtime"]; exists {
		if boolVal, ok := collectRuntime.(bool); ok {
			target.CollectRuntime = boolVal
		}
	}
	if collectHost, exists := configMap["collect_host"]; exists {
		if boolVal, ok := collectHost.(bool); ok {
			target.CollectHost = boolVal
		}
	}
	if interval, exists := configMap["interval"]; exists {
		if str, ok := interval.(string); ok {
			target.Interval = str
		}
	}
}

func parseLoggingConfig(configMap map[string]any, target *LoggingConfig) {
	if enabled, exists := configMap["enabled"]; exists {
		if boolVal, ok := enabled.(bool); ok {
			target.Enabled = boolVal
		}
	}
	if includeTrace, exists := configMap["include_trace"]; exists {
		if boolVal, ok := includeTrace.(bool); ok {
			target.IncludeTrace = boolVal
		}
	}
	if logLevel, exists := configMap["log_level"]; exists {
		if str, ok := logLevel.(string); ok {
			target.LogLevel = str
		}
	}
	if correlationID, exists := configMap["correlation_id"]; exists {
		if boolVal, ok := correlationID.(bool); ok {
			target.CorrelationID = boolVal
		}
	}
}

func parseResourceConfig(configMap map[string]any, target *ResourceConfig) {
	if detectHost, exists := configMap["detect_host"]; exists {
		if boolVal, ok := detectHost.(bool); ok {
			target.DetectHost = boolVal
		}
	}
	if detectProcess, exists := configMap["detect_process"]; exists {
		if boolVal, ok := detectProcess.(bool); ok {
			target.DetectProcess = boolVal
		}
	}
	if detectRuntime, exists := configMap["detect_runtime"]; exists {
		if boolVal, ok := detectRuntime.(bool); ok {
			target.DetectRuntime = boolVal
		}
	}
	if organization, exists := configMap["organization"]; exists {
		if str, ok := organization.(string); ok {
			target.Organization = str
		}
	}
	if team, exists := configMap["team"]; exists {
		if str, ok := team.(string); ok {
			target.Team = str
		}
	}
	if deploymentEnv, exists := configMap["deployment_env"]; exists {
		if str, ok := deploymentEnv.(string); ok {
			target.DeploymentEnv = str
		}
	}
	if customAttrs, exists := configMap["custom_attrs"]; exists {
		if attrMap, ok := customAttrs.(map[string]any); ok {
			strAttrs := make(map[string]string)
			for k, v := range attrMap {
				if str, ok := v.(string); ok {
					strAttrs[k] = str
				}
			}
			target.CustomAttrs = strAttrs
		}
	}
}

func parseExportersConfig(configMap map[string]any, target *ExportersConfig) {
	if otlp, exists := configMap["otlp"]; exists {
		if otlpMap, ok := otlp.(map[string]any); ok {
			parseOTLPConfig(otlpMap, &target.OTLP)
		}
	}
	if jaeger, exists := configMap["jaeger"]; exists {
		if jaegerMap, ok := jaeger.(map[string]any); ok {
			parseJaegerConfig(jaegerMap, &target.Jaeger)
		}
	}
	if prometheus, exists := configMap["prometheus"]; exists {
		if prometheusMap, ok := prometheus.(map[string]any); ok {
			parsePrometheusConfig(prometheusMap, &target.Prometheus)
		}
	}
	if console, exists := configMap["console"]; exists {
		if consoleMap, ok := console.(map[string]any); ok {
			parseConsoleConfig(consoleMap, &target.Console)
		}
	}
}

func parseOTLPConfig(configMap map[string]any, target *OTLPConfig) {
	if enabled, exists := configMap["enabled"]; exists {
		if boolVal, ok := enabled.(bool); ok {
			target.Enabled = boolVal
		}
	}
	if endpoint, exists := configMap["endpoint"]; exists {
		if str, ok := endpoint.(string); ok {
			target.Endpoint = str
		}
	}
	if insecure, exists := configMap["insecure"]; exists {
		if boolVal, ok := insecure.(bool); ok {
			target.Insecure = boolVal
		}
	}
	if compression, exists := configMap["compression"]; exists {
		if str, ok := compression.(string); ok {
			target.Compression = str
		}
	}
	if timeout, exists := configMap["timeout"]; exists {
		if str, ok := timeout.(string); ok {
			target.Timeout = str
		}
	}
	if headers, exists := configMap["headers"]; exists {
		if headerMap, ok := headers.(map[string]any); ok {
			strHeaders := make(map[string]string)
			for k, v := range headerMap {
				if str, ok := v.(string); ok {
					strHeaders[k] = str
				}
			}
			target.Headers = strHeaders
		}
	}
}

func parseJaegerConfig(configMap map[string]any, target *JaegerConfig) {
	if enabled, exists := configMap["enabled"]; exists {
		if boolVal, ok := enabled.(bool); ok {
			target.Enabled = boolVal
		}
	}
	if endpoint, exists := configMap["endpoint"]; exists {
		if str, ok := endpoint.(string); ok {
			target.Endpoint = str
		}
	}
	if username, exists := configMap["username"]; exists {
		if str, ok := username.(string); ok {
			target.Username = str
		}
	}
	if password, exists := configMap["password"]; exists {
		if str, ok := password.(string); ok {
			target.Password = str
		}
	}
}

func parsePrometheusConfig(configMap map[string]any, target *PrometheusConfig) {
	if enabled, exists := configMap["enabled"]; exists {
		if boolVal, ok := enabled.(bool); ok {
			target.Enabled = boolVal
		}
	}
	if host, exists := configMap["host"]; exists {
		if str, ok := host.(string); ok {
			target.Host = str
		}
	}
	if port, exists := configMap["port"]; exists {
		if intVal, ok := port.(int); ok {
			target.Port = intVal
		}
	}
	if path, exists := configMap["path"]; exists {
		if str, ok := path.(string); ok {
			target.Path = str
		}
	}
}

func parseConsoleConfig(configMap map[string]any, target *ConsoleConfig) {
	if enabled, exists := configMap["enabled"]; exists {
		if boolVal, ok := enabled.(bool); ok {
			target.Enabled = boolVal
		}
	}
	if pretty, exists := configMap["pretty"]; exists {
		if boolVal, ok := pretty.(bool); ok {
			target.Pretty = boolVal
		}
	}
}

// Plugin interface implementation
// zh: Plugin 介面實作

// Name returns the plugin name.
// zh: Name 回傳插件名稱。
func (s *SDKWrapperPlugin) Name() string {
	return s.name
}

// Version returns the plugin version.
// zh: Version 回傳插件版本。
func (s *SDKWrapperPlugin) Version() string {
	return s.version
}

// Description returns the plugin description.
// zh: Description 回傳插件描述。
func (s *SDKWrapperPlugin) Description() string {
	return s.description
}

// Init initializes the OpenTelemetry SDK wrapper plugin.
// zh: Init 初始化 OpenTelemetry SDK 包裝器插件。
func (s *SDKWrapperPlugin) Init(config any) error {
	if s.initialized {
		return nil
	}

	// Parse configuration
	if config != nil {
		if err := parseSDKConfig(config, s.config); err != nil {
			return fmt.Errorf("failed to parse SDK wrapper config: %w", err)
		}
	}

	// Validate configuration
	if err := s.validateConfig(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	s.initialized = true

	ctx := context.Background()
	log.L(ctx).Info("OpenTelemetry SDK wrapper plugin initialized",
		"service_name", s.config.ServiceName,
		"service_version", s.config.ServiceVersion,
		"environment", s.config.Environment,
		"enabled", s.config.Enabled)

	return nil
}

// Shutdown shuts down the OpenTelemetry SDK wrapper plugin.
// zh: Shutdown 關閉 OpenTelemetry SDK 包裝器插件。
func (s *SDKWrapperPlugin) Shutdown() error {
	s.started = false
	s.initialized = false

	ctx := context.Background()
	log.L(ctx).Info("OpenTelemetry SDK wrapper plugin shutdown")

	return nil
}

// LifecycleAware interface implementation
// zh: LifecycleAware 介面實作

// OnRegister is called when the plugin is registered.
// zh: OnRegister 在插件註冊時被呼叫。
func (s *SDKWrapperPlugin) OnRegister() error {
	ctx := context.Background()
	log.L(ctx).Info("OpenTelemetry SDK wrapper plugin registered")
	return nil
}

// OnStart is called when the plugin is started.
// zh: OnStart 在插件啟動時被呼叫。
func (s *SDKWrapperPlugin) OnStart() error {
	if !s.initialized {
		return fmt.Errorf("plugin not initialized")
	}

	if !s.config.Enabled {
		ctx := context.Background()
		log.L(ctx).Info("OpenTelemetry SDK wrapper plugin disabled, skipping start")
		return nil
	}

	// Initialize OpenTelemetry SDK components
	if err := s.initializeSDK(); err != nil {
		return fmt.Errorf("failed to initialize OpenTelemetry SDK: %w", err)
	}

	s.started = true

	ctx := context.Background()
	log.L(ctx).Info("OpenTelemetry SDK wrapper plugin started",
		"tracing_enabled", s.config.Tracing.Enabled,
		"metrics_enabled", s.config.Metrics.Enabled,
		"logging_enabled", s.config.Logging.Enabled)

	return nil
}

// OnStop is called when the plugin is stopped.
// zh: OnStop 在插件停止時被呼叫。
func (s *SDKWrapperPlugin) OnStop() error {
	if s.started {
		// Cleanup OpenTelemetry SDK components
		if err := s.shutdownSDK(); err != nil {
			ctx := context.Background()
			log.L(ctx).Error("Failed to shutdown OpenTelemetry SDK", "error", err)
		}
	}

	s.started = false

	ctx := context.Background()
	log.L(ctx).Info("OpenTelemetry SDK wrapper plugin stopped")

	return nil
}

// OnShutdown is called when the plugin is completely shut down.
// zh: OnShutdown 在插件完全關閉時被呼叫。
func (s *SDKWrapperPlugin) OnShutdown() error {
	ctx := context.Background()
	log.L(ctx).Info("OpenTelemetry SDK wrapper plugin shutdown complete")
	return nil
}

// HealthChecker interface implementation
// zh: HealthChecker 介面實作

// CheckHealth checks the health of the OpenTelemetry SDK wrapper.
// zh: CheckHealth 檢查 OpenTelemetry SDK 包裝器的健康狀態。
func (s *SDKWrapperPlugin) CheckHealth(ctx context.Context) contracts.HealthStatus {
	status := contracts.HealthStatus{
		Status:    "healthy",
		Message:   "OpenTelemetry SDK wrapper is running normally",
		Timestamp: time.Now(),
		Details:   make(map[string]any),
	}

	// Check if plugin is initialized and started
	if !s.initialized {
		status.Status = "unhealthy"
		status.Message = "Plugin not initialized"
		return status
	}

	if !s.started && s.config.Enabled {
		status.Status = "degraded"
		status.Message = "Plugin initialized but not started"
	}

	// Add configuration details
	status.Details["service_name"] = s.config.ServiceName
	status.Details["service_version"] = s.config.ServiceVersion
	status.Details["environment"] = s.config.Environment
	status.Details["enabled"] = s.config.Enabled
	status.Details["tracing_enabled"] = s.config.Tracing.Enabled
	status.Details["metrics_enabled"] = s.config.Metrics.Enabled
	status.Details["logging_enabled"] = s.config.Logging.Enabled

	// Check exporter configurations
	exporters := make(map[string]bool)
	exporters["otlp"] = s.config.Exporters.OTLP.Enabled
	exporters["jaeger"] = s.config.Exporters.Jaeger.Enabled
	exporters["prometheus"] = s.config.Exporters.Prometheus.Enabled
	exporters["console"] = s.config.Exporters.Console.Enabled
	status.Details["exporters"] = exporters

	return status
}

// GetHealthMetrics returns health metrics for the OpenTelemetry SDK wrapper.
// zh: GetHealthMetrics 回傳 OpenTelemetry SDK 包裝器的健康指標。
func (s *SDKWrapperPlugin) GetHealthMetrics() map[string]any {
	metrics := make(map[string]any)

	metrics["plugin_initialized"] = s.initialized
	metrics["plugin_started"] = s.started
	metrics["plugin_enabled"] = s.config.Enabled

	// Configuration metrics
	metrics["tracing_enabled"] = s.config.Tracing.Enabled
	metrics["metrics_enabled"] = s.config.Metrics.Enabled
	metrics["logging_enabled"] = s.config.Logging.Enabled

	// Exporter metrics
	metrics["otlp_exporter_enabled"] = s.config.Exporters.OTLP.Enabled
	metrics["jaeger_exporter_enabled"] = s.config.Exporters.Jaeger.Enabled
	metrics["prometheus_exporter_enabled"] = s.config.Exporters.Prometheus.Enabled
	metrics["console_exporter_enabled"] = s.config.Exporters.Console.Enabled

	// Resource metrics
	metrics["detect_host"] = s.config.Resource.DetectHost
	metrics["detect_process"] = s.config.Resource.DetectProcess
	metrics["detect_runtime"] = s.config.Resource.DetectRuntime

	return metrics
}

// Helper methods
// zh: 輔助方法

// validateConfig validates the plugin configuration.
// zh: validateConfig 驗證插件配置。
func (s *SDKWrapperPlugin) validateConfig() error {
	if s.config.ServiceName == "" {
		return fmt.Errorf("service_name is required")
	}

	if s.config.ServiceVersion == "" {
		return fmt.Errorf("service_version is required")
	}

	// Validate tracing configuration
	if s.config.Tracing.Enabled {
		if s.config.Tracing.SampleRatio < 0 || s.config.Tracing.SampleRatio > 1 {
			return fmt.Errorf("tracing sample_ratio must be between 0 and 1")
		}
		if s.config.Tracing.MaxSpans <= 0 {
			return fmt.Errorf("tracing max_spans must be positive")
		}
		if s.config.Tracing.BatchSize <= 0 {
			return fmt.Errorf("tracing batch_size must be positive")
		}
	}

	// Validate OTLP exporter configuration
	if s.config.Exporters.OTLP.Enabled {
		if s.config.Exporters.OTLP.Endpoint == "" {
			return fmt.Errorf("OTLP endpoint is required when OTLP exporter is enabled")
		}
	}

	// Validate Jaeger exporter configuration
	if s.config.Exporters.Jaeger.Enabled {
		if s.config.Exporters.Jaeger.Endpoint == "" {
			return fmt.Errorf("Jaeger endpoint is required when Jaeger exporter is enabled")
		}
	}

	// Validate Prometheus exporter configuration
	if s.config.Exporters.Prometheus.Enabled {
		if s.config.Exporters.Prometheus.Port <= 0 || s.config.Exporters.Prometheus.Port > 65535 {
			return fmt.Errorf("Prometheus port must be between 1 and 65535")
		}
	}

	return nil
}

// initializeSDK initializes the OpenTelemetry SDK components.
// zh: initializeSDK 初始化 OpenTelemetry SDK 組件。
func (s *SDKWrapperPlugin) initializeSDK() error {
	ctx := context.Background()

	// This is a placeholder implementation
	// In a real implementation, this would:
	// 1. Initialize TracerProvider with configured exporters
	// 2. Initialize MeterProvider with configured exporters
	// 3. Initialize LoggerProvider with configured exporters
	// 4. Set up resource detection
	// 5. Configure sampling strategies
	// 6. Set up automatic instrumentation

	log.L(ctx).Info("Initializing OpenTelemetry SDK components",
		"service_name", s.config.ServiceName,
		"otlp_endpoint", s.config.Exporters.OTLP.Endpoint)

	// TODO: Implement actual SDK initialization
	// Example:
	// - Set up OTLP exporters
	// - Configure resource detection
	// - Initialize tracer provider
	// - Initialize meter provider
	// - Set up automatic instrumentation

	return nil
}

// shutdownSDK shuts down the OpenTelemetry SDK components.
// zh: shutdownSDK 關閉 OpenTelemetry SDK 組件。
func (s *SDKWrapperPlugin) shutdownSDK() error {
	ctx := context.Background()

	// This is a placeholder implementation
	// In a real implementation, this would:
	// 1. Flush and shutdown TracerProvider
	// 2. Flush and shutdown MeterProvider
	// 3. Flush and shutdown LoggerProvider
	// 4. Clean up resources

	log.L(ctx).Info("Shutting down OpenTelemetry SDK components")

	// TODO: Implement actual SDK shutdown
	// Example:
	// - Flush pending spans
	// - Shutdown tracer provider
	// - Shutdown meter provider
	// - Clean up resources

	return nil
}

// GetConfig returns the current configuration.
// zh: GetConfig 回傳當前配置。
func (s *SDKWrapperPlugin) GetConfig() *Config {
	return s.config
}

// IsEnabled returns whether the plugin is enabled.
// zh: IsEnabled 回傳插件是否已啟用。
func (s *SDKWrapperPlugin) IsEnabled() bool {
	return s.config.Enabled
}

// IsStarted returns whether the plugin is started.
// zh: IsStarted 回傳插件是否已啟動。
func (s *SDKWrapperPlugin) IsStarted() bool {
	return s.started
}

// Register registers the OpenTelemetry SDK wrapper plugin.
// zh: Register 註冊 OpenTelemetry SDK 包裝器插件。
func Register(registry contracts.Registry) error {
	return registry.RegisterPlugin("otel-sdk-wrapper", NewSDKWrapperPlugin)
}
