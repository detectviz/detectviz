package gzip

import (
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"detectviz/pkg/platform/contracts"
)

// GzipMiddleware implements HTTP response compression middleware.
// zh: GzipMiddleware 實作 HTTP 回應壓縮中介層。
type GzipMiddleware struct {
	name        string
	version     string
	description string
	config      *GzipConfig
	initialized bool
}

// GzipConfig defines the configuration for gzip middleware.
// zh: GzipConfig 定義 gzip 中介層的配置。
type GzipConfig struct {
	Level         int      `yaml:"level" json:"level" mapstructure:"level"`                            // Compression level (1-9)
	MinLength     int      `yaml:"min_length" json:"min_length" mapstructure:"min_length"`             // Minimum response size to compress
	ExcludedTypes []string `yaml:"excluded_types" json:"excluded_types" mapstructure:"excluded_types"` // MIME types to exclude
	ExcludedPaths []string `yaml:"excluded_paths" json:"excluded_paths" mapstructure:"excluded_paths"` // Paths to exclude
	IncludedTypes []string `yaml:"included_types" json:"included_types" mapstructure:"included_types"` // MIME types to include (if specified, only these will be compressed)
}

// gzipResponseWriter wraps http.ResponseWriter to provide gzip compression.
// zh: gzipResponseWriter 包裝 http.ResponseWriter 以提供 gzip 壓縮。
type gzipResponseWriter struct {
	http.ResponseWriter
	gzipWriter *gzip.Writer
	config     *GzipConfig
	written    bool
	size       int
}

// NewGzipMiddleware creates a new gzip middleware instance.
// zh: NewGzipMiddleware 建立新的 gzip 中介層實例。
func NewGzipMiddleware(config any) (contracts.Plugin, error) {
	gzipConfig := &GzipConfig{
		Level:     gzip.DefaultCompression,
		MinLength: 1024, // 1KB minimum
		ExcludedTypes: []string{
			"image/jpeg",
			"image/png",
			"image/gif",
			"image/webp",
			"video/*",
			"audio/*",
			"application/zip",
			"application/gzip",
			"application/x-gzip",
		},
		ExcludedPaths: []string{
			"/health",
			"/metrics",
		},
	}

	// Parse config from the provided config parameter
	if config != nil {
		if err := parseGzipConfig(config, gzipConfig); err != nil {
			return nil, fmt.Errorf("failed to parse gzip config: %w", err)
		}
	}

	return &GzipMiddleware{
		name:        "gzip-middleware",
		version:     "1.0.0",
		description: "HTTP response compression middleware",
		config:      gzipConfig,
		initialized: false,
	}, nil
}

// parseGzipConfig parses the plugin configuration from various formats
// zh: parseGzipConfig 從各種格式解析插件配置
func parseGzipConfig(config any, target *GzipConfig) error {
	if config == nil {
		return nil
	}

	// Handle map[string]any format
	if configMap, ok := config.(map[string]any); ok {
		if level, exists := configMap["level"]; exists {
			if intVal, ok := level.(int); ok {
				target.Level = intVal
			}
		}
		if minLength, exists := configMap["min_length"]; exists {
			if intVal, ok := minLength.(int); ok {
				target.MinLength = intVal
			}
		}
		if excludedTypes, exists := configMap["excluded_types"]; exists {
			if slice, ok := excludedTypes.([]any); ok {
				target.ExcludedTypes = make([]string, len(slice))
				for i, v := range slice {
					if str, ok := v.(string); ok {
						target.ExcludedTypes[i] = str
					}
				}
			}
		}
		if excludedPaths, exists := configMap["excluded_paths"]; exists {
			if slice, ok := excludedPaths.([]any); ok {
				target.ExcludedPaths = make([]string, len(slice))
				for i, v := range slice {
					if str, ok := v.(string); ok {
						target.ExcludedPaths[i] = str
					}
				}
			}
		}
		if includedTypes, exists := configMap["included_types"]; exists {
			if slice, ok := includedTypes.([]any); ok {
				target.IncludedTypes = make([]string, len(slice))
				for i, v := range slice {
					if str, ok := v.(string); ok {
						target.IncludedTypes[i] = str
					}
				}
			}
		}
		return nil
	}

	return fmt.Errorf("unsupported config format: %T", config)
}

// Plugin interface implementation
// zh: Plugin 介面實作

// Name returns the plugin name.
// zh: Name 回傳插件名稱。
func (g *GzipMiddleware) Name() string {
	return g.name
}

// Version returns the plugin version.
// zh: Version 回傳插件版本。
func (g *GzipMiddleware) Version() string {
	return g.version
}

// Description returns the plugin description.
// zh: Description 回傳插件描述。
func (g *GzipMiddleware) Description() string {
	return g.description
}

// Init initializes the gzip middleware.
// zh: Init 初始化 gzip 中介層。
func (g *GzipMiddleware) Init(config any) error {
	if g.initialized {
		return nil
	}

	// Validate compression level
	if g.config.Level < gzip.NoCompression || g.config.Level > gzip.BestCompression {
		return fmt.Errorf("invalid compression level: %d (must be between %d and %d)",
			g.config.Level, gzip.NoCompression, gzip.BestCompression)
	}

	// Validate minimum length
	if g.config.MinLength < 0 {
		return fmt.Errorf("minimum length cannot be negative: %d", g.config.MinLength)
	}

	g.initialized = true
	return nil
}

// Shutdown shuts down the gzip middleware.
// zh: Shutdown 關閉 gzip 中介層。
func (g *GzipMiddleware) Shutdown() error {
	g.initialized = false
	return nil
}

// LifecycleAware interface implementation
// zh: LifecycleAware 介面實作

// OnRegister is called when the plugin is registered.
// zh: OnRegister 在插件註冊時呼叫。
func (g *GzipMiddleware) OnRegister() error {
	return nil
}

// OnStart is called when the plugin is started.
// zh: OnStart 在插件啟動時呼叫。
func (g *GzipMiddleware) OnStart() error {
	if !g.initialized {
		return fmt.Errorf("gzip middleware not initialized")
	}
	return nil
}

// OnStop is called when the plugin is stopped.
// zh: OnStop 在插件停止時呼叫。
func (g *GzipMiddleware) OnStop() error {
	return nil
}

// OnShutdown is called when the plugin is shutdown.
// zh: OnShutdown 在插件關閉時呼叫。
func (g *GzipMiddleware) OnShutdown() error {
	return g.Shutdown()
}

// HealthChecker interface implementation
// zh: HealthChecker 介面實作

// CheckHealth checks the health of the gzip middleware.
// zh: CheckHealth 檢查 gzip 中介層的健康狀況。
func (g *GzipMiddleware) CheckHealth(ctx context.Context) contracts.HealthStatus {
	status := contracts.HealthStatus{
		Timestamp: time.Now(),
		Details:   make(map[string]any),
	}

	if !g.initialized {
		status.Status = "unhealthy"
		status.Message = "Gzip middleware not initialized"
		return status
	}

	status.Status = "healthy"
	status.Message = "Gzip middleware is healthy"
	status.Details["compression_level"] = g.config.Level
	status.Details["min_length"] = g.config.MinLength
	status.Details["excluded_types_count"] = len(g.config.ExcludedTypes)
	status.Details["excluded_paths_count"] = len(g.config.ExcludedPaths)

	return status
}

// GetHealthMetrics returns health metrics for the gzip middleware.
// zh: GetHealthMetrics 回傳 gzip 中介層的健康指標。
func (g *GzipMiddleware) GetHealthMetrics() map[string]any {
	return map[string]any{
		"initialized":       g.initialized,
		"compression_level": g.config.Level,
		"min_length":        g.config.MinLength,
		"excluded_types":    g.config.ExcludedTypes,
		"excluded_paths":    g.config.ExcludedPaths,
		"included_types":    g.config.IncludedTypes,
	}
}

// Middleware functionality
// zh: 中介層功能

// Handler returns the HTTP middleware handler.
// zh: Handler 回傳 HTTP 中介層處理器。
func (g *GzipMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if client accepts gzip encoding
			if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
				next.ServeHTTP(w, r)
				return
			}

			// Check if path should be excluded
			if g.shouldExcludePath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Create gzip response writer
			gzipWriter := &gzipResponseWriter{
				ResponseWriter: w,
				config:         g.config,
			}

			// Set content encoding header
			w.Header().Set("Content-Encoding", "gzip")
			w.Header().Set("Vary", "Accept-Encoding")

			// Call next handler
			next.ServeHTTP(gzipWriter, r)

			// Close gzip writer if it was created
			if gzipWriter.gzipWriter != nil {
				gzipWriter.gzipWriter.Close()
			}
		})
	}
}

// shouldExcludePath checks if a path should be excluded from compression.
// zh: shouldExcludePath 檢查路徑是否應該排除壓縮。
func (g *GzipMiddleware) shouldExcludePath(path string) bool {
	for _, excludedPath := range g.config.ExcludedPaths {
		if strings.HasPrefix(path, excludedPath) {
			return true
		}
	}
	return false
}

// shouldCompressContentType checks if a content type should be compressed.
// zh: shouldCompressContentType 檢查內容類型是否應該壓縮。
func (g *GzipMiddleware) shouldCompressContentType(contentType string) bool {
	// If included types are specified, only compress those
	if len(g.config.IncludedTypes) > 0 {
		for _, includedType := range g.config.IncludedTypes {
			if strings.Contains(contentType, includedType) {
				return true
			}
		}
		return false
	}

	// Otherwise, compress unless excluded
	for _, excludedType := range g.config.ExcludedTypes {
		if strings.Contains(contentType, excludedType) {
			return false
		}
	}

	return true
}

// gzipResponseWriter methods
// zh: gzipResponseWriter 方法

// Write writes data to the response, applying compression if appropriate.
// zh: Write 將資料寫入回應，如果適當則應用壓縮。
func (grw *gzipResponseWriter) Write(data []byte) (int, error) {
	if !grw.written {
		grw.written = true

		// Check content type
		contentType := grw.Header().Get("Content-Type")
		if contentType == "" {
			contentType = http.DetectContentType(data)
			grw.Header().Set("Content-Type", contentType)
		}

		// Check if we should compress this content type
		if !grw.shouldCompressContentType(contentType) {
			// Remove compression headers and write directly
			grw.Header().Del("Content-Encoding")
			grw.Header().Del("Vary")
			return grw.ResponseWriter.Write(data)
		}

		// Check minimum length
		if len(data) < grw.config.MinLength {
			// Remove compression headers and write directly
			grw.Header().Del("Content-Encoding")
			grw.Header().Del("Vary")
			return grw.ResponseWriter.Write(data)
		}

		// Initialize gzip writer
		var err error
		grw.gzipWriter, err = gzip.NewWriterLevel(grw.ResponseWriter, grw.config.Level)
		if err != nil {
			// Fallback to no compression
			grw.Header().Del("Content-Encoding")
			grw.Header().Del("Vary")
			return grw.ResponseWriter.Write(data)
		}
	}

	if grw.gzipWriter != nil {
		n, err := grw.gzipWriter.Write(data)
		grw.size += n
		return n, err
	}

	return grw.ResponseWriter.Write(data)
}

// shouldCompressContentType checks if a content type should be compressed.
// zh: shouldCompressContentType 檢查內容類型是否應該壓縮。
func (grw *gzipResponseWriter) shouldCompressContentType(contentType string) bool {
	// If included types are specified, only compress those
	if len(grw.config.IncludedTypes) > 0 {
		for _, includedType := range grw.config.IncludedTypes {
			if strings.Contains(contentType, includedType) {
				return true
			}
		}
		return false
	}

	// Otherwise, compress unless excluded
	for _, excludedType := range grw.config.ExcludedTypes {
		if strings.Contains(contentType, excludedType) {
			return false
		}
	}

	return true
}

// WriteHeader writes the status code to the response.
// zh: WriteHeader 將狀態碼寫入回應。
func (grw *gzipResponseWriter) WriteHeader(code int) {
	grw.ResponseWriter.WriteHeader(code)
}

// Flush flushes the response writer.
// zh: Flush 刷新回應寫入器。
func (grw *gzipResponseWriter) Flush() {
	if grw.gzipWriter != nil {
		grw.gzipWriter.Flush()
	}

	if flusher, ok := grw.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

// Close closes the gzip writer.
// zh: Close 關閉 gzip 寫入器。
func (grw *gzipResponseWriter) Close() error {
	if grw.gzipWriter != nil {
		return grw.gzipWriter.Close()
	}
	return nil
}

// Plugin registration
// zh: 插件註冊

// Register registers the gzip middleware plugin.
// zh: Register 註冊 gzip 中介層插件。
func Register(registry contracts.Registry) error {
	return registry.RegisterPlugin("gzip-middleware", NewGzipMiddleware)
}
