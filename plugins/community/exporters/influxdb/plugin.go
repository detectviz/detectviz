package influxdb

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"detectviz/pkg/platform/contracts"
)

// InfluxDBExporter implements InfluxDB metrics exporter.
// zh: InfluxDBExporter 實作 InfluxDB 指標匯出器。
type InfluxDBExporter struct {
	name        string
	version     string
	description string
	config      *Config
	httpClient  *http.Client
	initialized bool
	started     bool
	exportCount int64
}

// Config contains configuration for the InfluxDB exporter.
// zh: Config 包含 InfluxDB 匯出器的配置。
type Config struct {
	URL          string            `yaml:"url" json:"url"`
	Token        string            `yaml:"token" json:"token"`
	Organization string            `yaml:"organization" json:"organization"`
	Bucket       string            `yaml:"bucket" json:"bucket"`
	Database     string            `yaml:"database" json:"database"` // For InfluxDB v1.x compatibility
	Username     string            `yaml:"username" json:"username"` // For InfluxDB v1.x
	Password     string            `yaml:"password" json:"password"` // For InfluxDB v1.x
	Precision    string            `yaml:"precision" json:"precision"`
	Timeout      int               `yaml:"timeout" json:"timeout"` // seconds
	BatchSize    int               `yaml:"batch_size" json:"batch_size"`
	FlushTimeout int               `yaml:"flush_timeout" json:"flush_timeout"` // seconds
	RetryCount   int               `yaml:"retry_count" json:"retry_count"`
	RetryDelay   int               `yaml:"retry_delay" json:"retry_delay"` // seconds
	Tags         map[string]string `yaml:"tags" json:"tags"`
	InsecureTLS  bool              `yaml:"insecure_tls" json:"insecure_tls"`
}

// ExportData represents data to be exported to InfluxDB.
// zh: ExportData 代表要匯出到 InfluxDB 的資料。
type ExportData struct {
	Measurement string                 `json:"measurement"`
	Tags        map[string]string      `json:"tags"`
	Fields      map[string]interface{} `json:"fields"`
	Timestamp   time.Time              `json:"timestamp"`
}

// NewInfluxDBExporter creates a new InfluxDB exporter instance.
// zh: NewInfluxDBExporter 建立新的 InfluxDB 匯出器實例。
func NewInfluxDBExporter(config any) (contracts.Plugin, error) {
	var cfg Config

	// Set default configuration
	cfg = Config{
		URL:          "http://localhost:8086",
		Precision:    "ns",
		Timeout:      30,
		BatchSize:    100,
		FlushTimeout: 10,
		RetryCount:   3,
		RetryDelay:   1,
		Tags:         make(map[string]string),
		InsecureTLS:  false,
	}

	// Parse configuration if provided
	if config != nil {
		if err := parseInfluxDBConfig(config, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse InfluxDB config: %w", err)
		}
	}

	// Validate required configuration
	if cfg.URL == "" {
		return nil, fmt.Errorf("url is required")
	}

	// For InfluxDB v2.x, token and organization are required
	// For InfluxDB v1.x, database is required
	if cfg.Token == "" && cfg.Database == "" {
		return nil, fmt.Errorf("either token (v2.x) or database (v1.x) must be specified")
	}

	if cfg.Token != "" && cfg.Organization == "" {
		return nil, fmt.Errorf("organization is required when using token authentication")
	}

	if cfg.Token != "" && cfg.Bucket == "" {
		return nil, fmt.Errorf("bucket is required when using token authentication")
	}

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	return &InfluxDBExporter{
		name:        "influxdb-exporter",
		version:     "1.0.0",
		description: "Export metrics to InfluxDB v1.x or v2.x",
		config:      &cfg,
		httpClient:  httpClient,
		initialized: false,
		started:     false,
		exportCount: 0,
	}, nil
}

// parseInfluxDBConfig manually parses configuration from map[string]any
// zh: parseInfluxDBConfig 手動從 map[string]any 解析配置
func parseInfluxDBConfig(config any, target *Config) error {
	if config == nil {
		return nil
	}

	configMap, ok := config.(map[string]any)
	if !ok {
		return fmt.Errorf("config must be a map[string]any")
	}

	if url, exists := configMap["url"]; exists {
		if str, ok := url.(string); ok {
			target.URL = str
		}
	}
	if token, exists := configMap["token"]; exists {
		if str, ok := token.(string); ok {
			target.Token = str
		}
	}
	if organization, exists := configMap["organization"]; exists {
		if str, ok := organization.(string); ok {
			target.Organization = str
		}
	}
	if bucket, exists := configMap["bucket"]; exists {
		if str, ok := bucket.(string); ok {
			target.Bucket = str
		}
	}
	if database, exists := configMap["database"]; exists {
		if str, ok := database.(string); ok {
			target.Database = str
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
	if precision, exists := configMap["precision"]; exists {
		if str, ok := precision.(string); ok {
			target.Precision = str
		}
	}
	if timeout, exists := configMap["timeout"]; exists {
		if intVal, ok := timeout.(int); ok {
			target.Timeout = intVal
		}
	}
	if batchSize, exists := configMap["batch_size"]; exists {
		if intVal, ok := batchSize.(int); ok {
			target.BatchSize = intVal
		}
	}
	if flushTimeout, exists := configMap["flush_timeout"]; exists {
		if intVal, ok := flushTimeout.(int); ok {
			target.FlushTimeout = intVal
		}
	}
	if retryCount, exists := configMap["retry_count"]; exists {
		if intVal, ok := retryCount.(int); ok {
			target.RetryCount = intVal
		}
	}
	if retryDelay, exists := configMap["retry_delay"]; exists {
		if intVal, ok := retryDelay.(int); ok {
			target.RetryDelay = intVal
		}
	}
	if tags, exists := configMap["tags"]; exists {
		if tagMap, ok := tags.(map[string]interface{}); ok {
			strTags := make(map[string]string)
			for k, v := range tagMap {
				if str, ok := v.(string); ok {
					strTags[k] = str
				}
			}
			target.Tags = strTags
		}
	}
	if insecureTLS, exists := configMap["insecure_tls"]; exists {
		if boolVal, ok := insecureTLS.(bool); ok {
			target.InsecureTLS = boolVal
		}
	}

	return nil
}

// Plugin interface implementation
// zh: Plugin 介面實作

// Name returns the plugin name.
// zh: Name 回傳插件名稱。
func (e *InfluxDBExporter) Name() string {
	return e.name
}

// Version returns the plugin version.
// zh: Version 回傳插件版本。
func (e *InfluxDBExporter) Version() string {
	return e.version
}

// Description returns the plugin description.
// zh: Description 回傳插件描述。
func (e *InfluxDBExporter) Description() string {
	return e.description
}

// Init initializes the plugin.
// zh: Init 初始化插件。
func (e *InfluxDBExporter) Init(config any) error {
	if e.initialized {
		return nil
	}

	// Test connection to validate configuration
	if err := e.testConnection(); err != nil {
		return fmt.Errorf("failed to validate InfluxDB connection: %w", err)
	}

	e.initialized = true
	return nil
}

// Shutdown shuts down the plugin.
// zh: Shutdown 關閉插件。
func (e *InfluxDBExporter) Shutdown() error {
	e.started = false
	e.initialized = false
	return nil
}

// LifecycleAware interface implementation
// zh: LifecycleAware 介面實作

// OnRegister is called when the plugin is registered.
// zh: OnRegister 在插件註冊時被呼叫。
func (e *InfluxDBExporter) OnRegister() error {
	fmt.Printf("InfluxDB exporter registered for URL: %s\n", e.config.URL)
	return nil
}

// OnStart is called when the plugin is started.
// zh: OnStart 在插件啟動時被呼叫。
func (e *InfluxDBExporter) OnStart() error {
	if !e.initialized {
		return fmt.Errorf("plugin not initialized")
	}

	// Test connection to InfluxDB
	if err := e.testConnection(); err != nil {
		return fmt.Errorf("failed to connect to InfluxDB: %w", err)
	}

	e.started = true
	fmt.Printf("InfluxDB exporter started successfully\n")
	return nil
}

// OnStop is called when the plugin is stopped.
// zh: OnStop 在插件停止時被呼叫。
func (e *InfluxDBExporter) OnStop() error {
	e.started = false
	fmt.Printf("InfluxDB exporter stopped\n")
	return nil
}

// OnShutdown is called when the plugin is shut down.
// zh: OnShutdown 在插件關閉時被呼叫。
func (e *InfluxDBExporter) OnShutdown() error {
	return e.Shutdown()
}

// HealthChecker interface implementation
// zh: HealthChecker 介面實作

// CheckHealth checks the health of the InfluxDB connection.
// zh: CheckHealth 檢查 InfluxDB 連接的健康狀態。
func (e *InfluxDBExporter) CheckHealth(ctx context.Context) contracts.HealthStatus {
	status := contracts.HealthStatus{
		Timestamp: time.Now(),
		Details:   make(map[string]any),
	}

	if !e.initialized {
		status.Status = "unhealthy"
		status.Message = "Plugin not initialized"
		return status
	}

	if !e.started {
		status.Status = "unhealthy"
		status.Message = "Plugin not started"
		return status
	}

	// Test connection to InfluxDB
	if err := e.testConnection(); err != nil {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("Connection failed: %v", err)
		return status
	}

	status.Status = "healthy"
	status.Message = "Connected to InfluxDB successfully"
	status.Details["url"] = e.config.URL
	status.Details["export_count"] = e.exportCount
	if e.config.Token != "" {
		status.Details["version"] = "2.x"
		status.Details["organization"] = e.config.Organization
		status.Details["bucket"] = e.config.Bucket
	} else {
		status.Details["version"] = "1.x"
		status.Details["database"] = e.config.Database
	}

	return status
}

// GetHealthMetrics returns health metrics for the plugin.
// zh: GetHealthMetrics 回傳插件的健康指標。
func (e *InfluxDBExporter) GetHealthMetrics() map[string]any {
	return map[string]any{
		"initialized":   e.initialized,
		"started":       e.started,
		"url":           e.config.URL,
		"export_count":  e.exportCount,
		"batch_size":    e.config.BatchSize,
		"flush_timeout": e.config.FlushTimeout,
		"retry_count":   e.config.RetryCount,
	}
}

// Exporter interface implementation
// zh: Exporter 介面實作

// Export exports a single data point to InfluxDB.
// zh: Export 匯出單個資料點到 InfluxDB。
func (e *InfluxDBExporter) Export(ctx context.Context, data interface{}) error {
	if !e.started {
		return fmt.Errorf("plugin not started")
	}

	exportData, ok := data.(ExportData)
	if !ok {
		return fmt.Errorf("invalid data type, expected ExportData")
	}

	// Convert to line protocol
	line := e.toLineProtocol(exportData)

	// Send to InfluxDB
	if err := e.writeData(ctx, []string{line}); err != nil {
		return fmt.Errorf("failed to export data: %w", err)
	}

	e.exportCount++
	return nil
}

// BatchExport exports multiple data points to InfluxDB.
// zh: BatchExport 批次匯出多個資料點到 InfluxDB。
func (e *InfluxDBExporter) BatchExport(ctx context.Context, batch []interface{}) error {
	if !e.started {
		return fmt.Errorf("plugin not started")
	}

	if len(batch) == 0 {
		return nil
	}

	// Convert all data to line protocol
	lines := make([]string, 0, len(batch))
	for _, data := range batch {
		exportData, ok := data.(ExportData)
		if !ok {
			return fmt.Errorf("invalid data type in batch, expected ExportData")
		}
		lines = append(lines, e.toLineProtocol(exportData))
	}

	// Send batch to InfluxDB
	if err := e.writeData(ctx, lines); err != nil {
		return fmt.Errorf("failed to export batch: %w", err)
	}

	e.exportCount += int64(len(batch))
	return nil
}

// GetMetrics returns export metrics (placeholder implementation).
// zh: GetMetrics 回傳匯出指標（佔位符實作）。
func (e *InfluxDBExporter) GetMetrics() map[string]interface{} {
	return map[string]interface{}{
		"export_count":  e.exportCount,
		"batch_size":    e.config.BatchSize,
		"flush_timeout": e.config.FlushTimeout,
		"retry_count":   e.config.RetryCount,
		"url":           e.config.URL,
	}
}

// Helper methods
// zh: 輔助方法

// testConnection tests the connection to InfluxDB.
// zh: testConnection 測試與 InfluxDB 的連接。
func (e *InfluxDBExporter) testConnection() error {
	var testURL string

	if e.config.Token != "" {
		// InfluxDB v2.x - test health endpoint
		testURL = fmt.Sprintf("%s/health", strings.TrimRight(e.config.URL, "/"))
	} else {
		// InfluxDB v1.x - test ping endpoint
		testURL = fmt.Sprintf("%s/ping", strings.TrimRight(e.config.URL, "/"))
	}

	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create test request: %w", err)
	}

	if e.config.Token != "" {
		req.Header.Set("Authorization", "Token "+e.config.Token)
	}

	resp, err := e.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("connection test failed with status: %d", resp.StatusCode)
	}

	return nil
}

// toLineProtocol converts ExportData to InfluxDB line protocol format.
// zh: toLineProtocol 將 ExportData 轉換為 InfluxDB 行協議格式。
func (e *InfluxDBExporter) toLineProtocol(data ExportData) string {
	var line strings.Builder

	// Measurement name
	line.WriteString(escapeKey(data.Measurement))

	// Tags (merge global tags with data tags)
	allTags := make(map[string]string)
	for k, v := range e.config.Tags {
		allTags[k] = v
	}
	for k, v := range data.Tags {
		allTags[k] = v
	}

	if len(allTags) > 0 {
		line.WriteString(",")
		first := true
		for key, value := range allTags {
			if !first {
				line.WriteString(",")
			}
			line.WriteString(escapeKey(key))
			line.WriteString("=")
			line.WriteString(escapeValue(value))
			first = false
		}
	}

	// Fields
	line.WriteString(" ")
	first := true
	for key, value := range data.Fields {
		if !first {
			line.WriteString(",")
		}
		line.WriteString(escapeKey(key))
		line.WriteString("=")
		line.WriteString(formatFieldValue(value))
		first = false
	}

	// Timestamp
	if !data.Timestamp.IsZero() {
		line.WriteString(" ")
		switch e.config.Precision {
		case "s":
			line.WriteString(strconv.FormatInt(data.Timestamp.Unix(), 10))
		case "ms":
			line.WriteString(strconv.FormatInt(data.Timestamp.UnixMilli(), 10))
		case "us":
			line.WriteString(strconv.FormatInt(data.Timestamp.UnixMicro(), 10))
		case "ns":
			line.WriteString(strconv.FormatInt(data.Timestamp.UnixNano(), 10))
		default:
			line.WriteString(strconv.FormatInt(data.Timestamp.UnixNano(), 10))
		}
	}

	return line.String()
}

// writeData sends data to InfluxDB.
// zh: writeData 將資料發送到 InfluxDB。
func (e *InfluxDBExporter) writeData(ctx context.Context, lines []string) error {
	data := strings.Join(lines, "\n")

	var writeURL string
	var req *http.Request
	var err error

	if e.config.Token != "" {
		// InfluxDB v2.x API
		writeURL = fmt.Sprintf("%s/api/v2/write", strings.TrimRight(e.config.URL, "/"))
		params := url.Values{}
		params.Set("org", e.config.Organization)
		params.Set("bucket", e.config.Bucket)
		params.Set("precision", e.config.Precision)
		writeURL += "?" + params.Encode()

		req, err = http.NewRequestWithContext(ctx, "POST", writeURL, bytes.NewBufferString(data))
		if err != nil {
			return fmt.Errorf("failed to create write request: %w", err)
		}

		req.Header.Set("Authorization", "Token "+e.config.Token)
		req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	} else {
		// InfluxDB v1.x API
		writeURL = fmt.Sprintf("%s/write", strings.TrimRight(e.config.URL, "/"))
		params := url.Values{}
		params.Set("db", e.config.Database)
		params.Set("precision", e.config.Precision)
		if e.config.Username != "" {
			params.Set("u", e.config.Username)
		}
		if e.config.Password != "" {
			params.Set("p", e.config.Password)
		}
		writeURL += "?" + params.Encode()

		req, err = http.NewRequestWithContext(ctx, "POST", writeURL, bytes.NewBufferString(data))
		if err != nil {
			return fmt.Errorf("failed to create write request: %w", err)
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	// Retry logic
	for attempt := 0; attempt <= e.config.RetryCount; attempt++ {
		resp, err := e.httpClient.Do(req)
		if err != nil {
			if attempt == e.config.RetryCount {
				return fmt.Errorf("write request failed after %d attempts: %w", e.config.RetryCount+1, err)
			}
			time.Sleep(time.Duration(e.config.RetryDelay) * time.Second)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent || resp.StatusCode == http.StatusOK {
			return nil
		}

		if attempt == e.config.RetryCount {
			return fmt.Errorf("write failed after %d attempts with status: %d", e.config.RetryCount+1, resp.StatusCode)
		}

		time.Sleep(time.Duration(e.config.RetryDelay) * time.Second)
	}

	return nil
}

// Helper functions for line protocol escaping
// zh: 行協議轉義的輔助函式

// escapeKey escapes keys in line protocol.
// zh: escapeKey 在行協議中轉義鍵。
func escapeKey(key string) string {
	key = strings.ReplaceAll(key, " ", "\\ ")
	key = strings.ReplaceAll(key, ",", "\\,")
	key = strings.ReplaceAll(key, "=", "\\=")
	return key
}

// escapeValue escapes tag values in line protocol.
// zh: escapeValue 在行協議中轉義標籤值。
func escapeValue(value string) string {
	value = strings.ReplaceAll(value, " ", "\\ ")
	value = strings.ReplaceAll(value, ",", "\\,")
	value = strings.ReplaceAll(value, "=", "\\=")
	return value
}

// formatFieldValue formats field values according to InfluxDB requirements.
// zh: formatFieldValue 根據 InfluxDB 要求格式化欄位值。
func formatFieldValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		// String values must be quoted
		escaped := strings.ReplaceAll(v, "\"", "\\\"")
		return fmt.Sprintf("\"%s\"", escaped)
	case int, int8, int16, int32, int64:
		return fmt.Sprintf("%di", v)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%du", v)
	case float32, float64:
		return fmt.Sprintf("%g", v)
	case bool:
		if v {
			return "t"
		}
		return "f"
	default:
		// Fallback to string representation
		escaped := strings.ReplaceAll(fmt.Sprintf("%v", v), "\"", "\\\"")
		return fmt.Sprintf("\"%s\"", escaped)
	}
}

// Register registers the InfluxDB exporter plugin.
// zh: Register 註冊 InfluxDB 匯出器插件。
func Register(registry contracts.Registry) error {
	return registry.RegisterPlugin("influxdb-exporter", NewInfluxDBExporter)
}
