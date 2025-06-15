package prometheus

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"detectviz/pkg/platform/contracts"
	"detectviz/pkg/shared/log"
)

// PrometheusImporter implements Prometheus metrics importer.
// zh: PrometheusImporter 實作 Prometheus 指標匯入器。
type PrometheusImporter struct {
	name        string
	version     string
	description string
	config      *PrometheusConfig
	initialized bool
	streaming   bool
	httpClient  *http.Client
	metrics     []PrometheusMetric
}

// PrometheusConfig defines the configuration for Prometheus importer.
// zh: PrometheusConfig 定義 Prometheus 匯入器的配置。
type PrometheusConfig struct {
	Endpoint       string `yaml:"endpoint" json:"endpoint" mapstructure:"endpoint"`
	ScrapeInterval string `yaml:"scrape_interval" json:"scrape_interval" mapstructure:"scrape_interval"`
	Timeout        string `yaml:"timeout" json:"timeout" mapstructure:"timeout"`
	MetricsPath    string `yaml:"metrics_path" json:"metrics_path" mapstructure:"metrics_path"`
	Username       string `yaml:"username" json:"username" mapstructure:"username"`
	Password       string `yaml:"password" json:"password" mapstructure:"password"`
	BearerToken    string `yaml:"bearer_token" json:"bearer_token" mapstructure:"bearer_token"`
}

// PrometheusMetric represents a single Prometheus metric.
// zh: PrometheusMetric 代表單個 Prometheus 指標。
type PrometheusMetric struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Labels    map[string]string `json:"labels"`
	Timestamp time.Time         `json:"timestamp"`
}

// NewPrometheusImporter creates a new Prometheus importer instance.
// zh: NewPrometheusImporter 建立新的 Prometheus 匯入器實例。
func NewPrometheusImporter(config any) (contracts.Plugin, error) {
	promConfig := &PrometheusConfig{
		Endpoint:       "http://localhost:9090",
		ScrapeInterval: "15s",
		Timeout:        "10s",
		MetricsPath:    "/metrics",
	}

	// Parse config from the provided config parameter
	if config != nil {
		if err := parsePrometheusConfig(config, promConfig); err != nil {
			return nil, fmt.Errorf("failed to parse Prometheus config: %w", err)
		}
	}

	return &PrometheusImporter{
		name:        "prometheus-importer",
		version:     "1.0.0",
		description: "Import metrics from Prometheus",
		config:      promConfig,
		initialized: false,
		streaming:   false,
		httpClient:  &http.Client{},
	}, nil
}

// parsePrometheusConfig parses the plugin configuration from various formats
// zh: parsePrometheusConfig 從各種格式解析插件配置
func parsePrometheusConfig(config any, target *PrometheusConfig) error {
	if config == nil {
		return nil
	}

	// Handle map[string]any format
	if configMap, ok := config.(map[string]any); ok {
		if endpoint, exists := configMap["endpoint"]; exists {
			if str, ok := endpoint.(string); ok {
				target.Endpoint = str
			}
		}
		if scrapeInterval, exists := configMap["scrape_interval"]; exists {
			if str, ok := scrapeInterval.(string); ok {
				target.ScrapeInterval = str
			}
		}
		if timeout, exists := configMap["timeout"]; exists {
			if str, ok := timeout.(string); ok {
				target.Timeout = str
			}
		}
		if metricsPath, exists := configMap["metrics_path"]; exists {
			if str, ok := metricsPath.(string); ok {
				target.MetricsPath = str
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
		if bearerToken, exists := configMap["bearer_token"]; exists {
			if str, ok := bearerToken.(string); ok {
				target.BearerToken = str
			}
		}
		return nil
	}

	// Handle struct format using reflection
	if reflect.TypeOf(config).Kind() == reflect.Struct {
		configValue := reflect.ValueOf(config)
		targetValue := reflect.ValueOf(target).Elem()

		for i := 0; i < configValue.NumField(); i++ {
			field := configValue.Type().Field(i)
			value := configValue.Field(i)

			// Find corresponding field in target
			if targetField := targetValue.FieldByName(field.Name); targetField.IsValid() && targetField.CanSet() {
				if value.Type().AssignableTo(targetField.Type()) {
					targetField.Set(value)
				}
			}
		}
	}

	return nil
}

// Name returns the plugin name.
// zh: Name 回傳插件名稱。
func (p *PrometheusImporter) Name() string {
	return p.name
}

// Version returns the plugin version.
// zh: Version 回傳插件版本。
func (p *PrometheusImporter) Version() string {
	return p.version
}

// Description returns the plugin description.
// zh: Description 回傳插件描述。
func (p *PrometheusImporter) Description() string {
	return p.description
}

// Init initializes the Prometheus importer.
// zh: Init 初始化 Prometheus 匯入器。
func (p *PrometheusImporter) Init(config any) error {
	if p.initialized {
		return nil
	}

	// Validate configuration
	if p.config.Endpoint == "" {
		return fmt.Errorf("Prometheus endpoint is required")
	}

	// Validate scrape interval format
	if _, err := time.ParseDuration(p.config.ScrapeInterval); err != nil {
		return fmt.Errorf("invalid scrape interval format: %w", err)
	}

	// Validate timeout format
	if _, err := time.ParseDuration(p.config.Timeout); err != nil {
		return fmt.Errorf("invalid timeout format: %w", err)
	}

	p.initialized = true
	return nil
}

// Shutdown shuts down the Prometheus importer.
// zh: Shutdown 關閉 Prometheus 匯入器。
func (p *PrometheusImporter) Shutdown() error {
	p.initialized = false
	p.streaming = false
	return nil
}

// LifecycleAware interface implementation
// zh: LifecycleAware 介面實作

// OnRegister is called when the plugin is registered.
// zh: OnRegister 在插件註冊時呼叫。
func (p *PrometheusImporter) OnRegister() error {
	return nil
}

// OnStart is called when the plugin is started.
// zh: OnStart 在插件啟動時呼叫。
func (p *PrometheusImporter) OnStart() error {
	if !p.initialized {
		return fmt.Errorf("Prometheus importer not initialized")
	}
	return nil
}

// OnStop is called when the plugin is stopped.
// zh: OnStop 在插件停止時呼叫。
func (p *PrometheusImporter) OnStop() error {
	return p.StopStreaming()
}

// OnShutdown is called when the plugin is shutdown.
// zh: OnShutdown 在插件關閉時呼叫。
func (p *PrometheusImporter) OnShutdown() error {
	return p.Shutdown()
}

// HealthChecker interface implementation
// zh: HealthChecker 介面實作

// CheckHealth checks the health of the Prometheus importer.
// zh: CheckHealth 檢查 Prometheus 匯入器的健康狀況。
func (p *PrometheusImporter) CheckHealth(ctx context.Context) contracts.HealthStatus {
	status := contracts.HealthStatus{
		Timestamp: time.Now(),
		Details:   make(map[string]any),
	}

	if !p.initialized {
		status.Status = "unhealthy"
		status.Message = "Prometheus importer not initialized"
		return status
	}

	// Check configuration validity
	if p.config.Endpoint == "" {
		status.Status = "unhealthy"
		status.Message = "Prometheus endpoint not configured"
		return status
	}

	// Check scrape interval format
	if _, err := time.ParseDuration(p.config.ScrapeInterval); err != nil {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("Invalid scrape interval format: %v", err)
		return status
	}

	// Check timeout format
	if _, err := time.ParseDuration(p.config.Timeout); err != nil {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("Invalid timeout format: %v", err)
		return status
	}

	// TODO: Add actual connectivity check to Prometheus endpoint

	status.Status = "healthy"
	status.Message = "Prometheus importer is healthy"
	status.Details["endpoint"] = p.config.Endpoint
	status.Details["scrape_interval"] = p.config.ScrapeInterval
	status.Details["timeout"] = p.config.Timeout
	status.Details["streaming"] = p.streaming

	return status
}

// GetHealthMetrics returns health metrics for the Prometheus importer.
// zh: GetHealthMetrics 回傳 Prometheus 匯入器的健康指標。
func (p *PrometheusImporter) GetHealthMetrics() map[string]any {
	return map[string]any{
		"initialized":     p.initialized,
		"streaming":       p.streaming,
		"endpoint":        p.config.Endpoint,
		"scrape_interval": p.config.ScrapeInterval,
		"timeout":         p.config.Timeout,
	}
}

// Import performs a data import from Prometheus.
// zh: Import 從 Prometheus 執行資料匯入。
func (p *PrometheusImporter) Import(ctx context.Context) error {
	if !p.initialized || !p.streaming {
		return fmt.Errorf("prometheus importer not initialized or started")
	}

	log.L(ctx).Info("Importing metrics from Prometheus endpoint", "endpoint", p.config.Endpoint)

	// Build query URL
	queryURL, err := p.buildQueryURL()
	if err != nil {
		return fmt.Errorf("failed to build query URL: %w", err)
	}

	// Make HTTP request
	resp, err := p.httpClient.Get(queryURL)
	if err != nil {
		return fmt.Errorf("failed to query Prometheus: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("prometheus query failed with status: %d", resp.StatusCode)
	}

	// TODO: Parse Prometheus response and convert to ImportData
	// This is a simplified implementation
	p.metrics = append(p.metrics, PrometheusMetric{
		Name:      "sample_metric",
		Value:     42.0,
		Labels:    map[string]string{"instance": "localhost:9090"},
		Timestamp: time.Now(),
	})

	return nil
}

// StartStreaming starts streaming metrics from Prometheus.
// zh: StartStreaming 開始從 Prometheus 串流指標。
func (p *PrometheusImporter) StartStreaming(ctx context.Context) (<-chan contracts.ImportData, error) {
	if !p.initialized {
		return nil, fmt.Errorf("Prometheus importer not initialized")
	}

	if p.streaming {
		return nil, fmt.Errorf("streaming already started")
	}

	dataChan := make(chan contracts.ImportData, 100)

	// Parse scrape interval
	scrapeInterval, err := time.ParseDuration(p.config.ScrapeInterval)
	if err != nil {
		return nil, fmt.Errorf("invalid scrape interval: %w", err)
	}

	go func() {
		defer close(dataChan)
		ticker := time.NewTicker(scrapeInterval)
		defer ticker.Stop()

		p.streaming = true
		defer func() { p.streaming = false }()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// TODO: Scrape metrics and send to channel
				data := contracts.ImportData{
					Type:      "metrics",
					Source:    p.name,
					Timestamp: time.Now(),
					Labels: map[string]string{
						"job":      "prometheus",
						"instance": p.config.Endpoint,
					},
					Fields: map[string]interface{}{
						"up": 1,
					},
				}

				select {
				case dataChan <- data:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return dataChan, nil
}

// StopStreaming stops streaming metrics.
// zh: StopStreaming 停止串流指標。
func (p *PrometheusImporter) StopStreaming() error {
	p.streaming = false
	return nil
}

// Register registers the Prometheus importer plugin.
// zh: Register 註冊 Prometheus 匯入器插件。
func Register(registry contracts.Registry) error {
	return registry.RegisterPlugin("prometheus-importer", NewPrometheusImporter)
}

// RegisterImporter registers the Prometheus importer with importer registry.
// zh: RegisterImporter 向匯入器註冊表註冊 Prometheus 匯入器。
func RegisterImporter(registry contracts.ImporterRegistry) error {
	factory := func(config contracts.ImporterConfig) (contracts.Importer, error) {
		plugin, err := NewPrometheusImporter(config)
		if err != nil {
			return nil, err
		}

		importer, ok := plugin.(contracts.Importer)
		if !ok {
			return nil, fmt.Errorf("plugin does not implement Importer interface")
		}

		return importer, nil
	}

	return registry.RegisterImporter("prometheus", factory)
}

func (p *PrometheusImporter) buildQueryURL() (string, error) {
	baseURL, err := url.Parse(p.config.Endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to parse Prometheus endpoint: %w", err)
	}

	queryParams := url.Values{}
	queryParams.Add("query", "up")
	queryParams.Add("start", strconv.FormatInt(time.Now().Add(-time.Minute).Unix(), 10))
	queryParams.Add("end", strconv.FormatInt(time.Now().Unix(), 10))
	queryParams.Add("step", p.config.ScrapeInterval)

	baseURL.RawQuery = queryParams.Encode()

	return baseURL.String(), nil
}
