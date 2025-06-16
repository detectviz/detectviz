package otelwrapper

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"

	"detectviz/pkg/platform/contracts"
	"detectviz/pkg/shared/log"
)

// OtelWrapper provides unified OpenTelemetry integration for DetectViz.
// zh: OtelWrapper 為 DetectViz 提供統一的 OpenTelemetry 整合。
type OtelWrapper struct {
	config         *Config
	tracerProvider trace.TracerProvider
	meterProvider  metric.MeterProvider
	loggerProvider contracts.LoggerProvider
	resource       *resource.Resource
	shutdownFuncs  []func(context.Context) error
	initialized    bool
	mutex          sync.RWMutex
}

// Config contains configuration for the OpenTelemetry wrapper.
// zh: Config 包含 OpenTelemetry 包裝器的配置。
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

// Global wrapper instance
var (
	globalWrapper *OtelWrapper
	wrapperOnce   sync.Once
)

// NewOtelWrapper creates a new OpenTelemetry wrapper instance.
// zh: NewOtelWrapper 建立新的 OpenTelemetry 包裝器實例。
func NewOtelWrapper(config *Config) *OtelWrapper {
	if config == nil {
		config = DefaultConfig()
	}

	return &OtelWrapper{
		config:        config,
		shutdownFuncs: make([]func(context.Context) error, 0),
		initialized:   false,
	}
}

// DefaultConfig returns the default configuration.
// zh: DefaultConfig 回傳預設配置。
func DefaultConfig() *Config {
	return &Config{
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
}

// Initialize initializes the OpenTelemetry wrapper.
// zh: Initialize 初始化 OpenTelemetry 包裝器。
func (w *OtelWrapper) Initialize(ctx context.Context) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.initialized {
		return nil
	}

	if !w.config.Enabled {
		log.L(ctx).Info("OpenTelemetry wrapper disabled, skipping initialization")
		return nil
	}

	// Initialize resource
	if err := w.initializeResource(ctx); err != nil {
		return fmt.Errorf("failed to initialize resource: %w", err)
	}

	// Initialize tracing
	if w.config.Tracing.Enabled {
		if err := w.initializeTracing(ctx); err != nil {
			return fmt.Errorf("failed to initialize tracing: %w", err)
		}
	}

	// Initialize metrics
	if w.config.Metrics.Enabled {
		if err := w.initializeMetrics(ctx); err != nil {
			return fmt.Errorf("failed to initialize metrics: %w", err)
		}
	}

	// Initialize logger provider (inject from plugin if available)
	if w.config.Logging.Enabled {
		if err := w.initializeLogging(ctx); err != nil {
			return fmt.Errorf("failed to initialize logging: %w", err)
		}
	}

	// Set up propagators
	w.setupPropagators()

	w.initialized = true
	log.L(ctx).Info("OpenTelemetry wrapper initialized successfully",
		"service_name", w.config.ServiceName,
		"service_version", w.config.ServiceVersion,
		"environment", w.config.Environment)

	return nil
}

// Shutdown shuts down the OpenTelemetry wrapper and releases resources.
// zh: Shutdown 關閉 OpenTelemetry 包裝器並釋放資源。
func (w *OtelWrapper) Shutdown(ctx context.Context) error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if !w.initialized {
		return nil
	}

	var errors []error

	// Call all shutdown functions in reverse order
	for i := len(w.shutdownFuncs) - 1; i >= 0; i-- {
		if err := w.shutdownFuncs[i](ctx); err != nil {
			errors = append(errors, err)
			log.L(ctx).Error("Error during shutdown", "error", err)
		}
	}

	w.initialized = false
	w.shutdownFuncs = w.shutdownFuncs[:0]

	if len(errors) > 0 {
		return fmt.Errorf("shutdown completed with %d errors", len(errors))
	}

	log.L(ctx).Info("OpenTelemetry wrapper shutdown complete")
	return nil
}

// Tracer returns a tracer with the given name.
// zh: Tracer 回傳指定名稱的追蹤器。
func (w *OtelWrapper) Tracer(name string) trace.Tracer {
	if w.tracerProvider == nil {
		return otel.Tracer(name)
	}
	return w.tracerProvider.Tracer(name)
}

// Meter returns a meter with the given name.
// zh: Meter 回傳指定名稱的指標器。
func (w *OtelWrapper) Meter(name string) metric.Meter {
	if w.meterProvider == nil {
		return otel.Meter(name)
	}
	return w.meterProvider.Meter(name)
}

// Logger returns the logger provider.
// zh: Logger 回傳日誌提供器。
func (w *OtelWrapper) Logger() contracts.LoggerProvider {
	return w.loggerProvider
}

// InjectLoggerProvider injects a logger provider from a plugin.
// zh: InjectLoggerProvider 從插件注入日誌提供器。
func (w *OtelWrapper) InjectLoggerProvider(provider contracts.LoggerProvider) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.loggerProvider = provider
}

// GetGlobalWrapper returns the global wrapper instance.
// zh: GetGlobalWrapper 回傳全域包裝器實例。
func GetGlobalWrapper() *OtelWrapper {
	wrapperOnce.Do(func() {
		globalWrapper = NewOtelWrapper(nil)
	})
	return globalWrapper
}

// SetGlobalWrapper sets the global wrapper instance.
// zh: SetGlobalWrapper 設定全域包裝器實例。
func SetGlobalWrapper(wrapper *OtelWrapper) {
	globalWrapper = wrapper
}

// Convenience functions for global access
// zh: 全域存取的便利函式

// Tracer returns a tracer from the global wrapper.
// zh: Tracer 從全域包裝器回傳追蹤器。
func Tracer(name string) trace.Tracer {
	return GetGlobalWrapper().Tracer(name)
}

// Meter returns a meter from the global wrapper.
// zh: Meter 從全域包裝器回傳指標器。
func Meter(name string) metric.Meter {
	return GetGlobalWrapper().Meter(name)
}

// Logger returns the logger provider from the global wrapper.
// zh: Logger 從全域包裝器回傳日誌提供器。
func Logger() contracts.LoggerProvider {
	return GetGlobalWrapper().Logger()
}

// Shutdown shuts down the global wrapper.
// zh: Shutdown 關閉全域包裝器。
func Shutdown(ctx context.Context) error {
	if globalWrapper != nil {
		return globalWrapper.Shutdown(ctx)
	}
	return nil
}
