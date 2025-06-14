package requestmeta

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"detectviz/pkg/platform/contracts"
)

// RequestMetaMiddleware implements HTTP request metadata processing middleware.
// zh: RequestMetaMiddleware 實作 HTTP 請求元資料處理中介層。
type RequestMetaMiddleware struct {
	name        string
	version     string
	description string
	config      *RequestMetaConfig
	initialized bool
}

// RequestMetaConfig defines the configuration for request metadata middleware.
// zh: RequestMetaConfig 定義請求元資料中介層的配置。
type RequestMetaConfig struct {
	RequestIDHeader    string   `yaml:"request_id_header" json:"request_id_header" mapstructure:"request_id_header"`          // Header name for request ID
	TraceIDHeader      string   `yaml:"trace_id_header" json:"trace_id_header" mapstructure:"trace_id_header"`                // Header name for trace ID
	UserAgentHeader    string   `yaml:"user_agent_header" json:"user_agent_header" mapstructure:"user_agent_header"`          // Header name for user agent
	ForwardedForHeader string   `yaml:"forwarded_for_header" json:"forwarded_for_header" mapstructure:"forwarded_for_header"` // Header name for forwarded for
	GenerateRequestID  bool     `yaml:"generate_request_id" json:"generate_request_id" mapstructure:"generate_request_id"`    // Whether to generate request ID if not present
	LogRequests        bool     `yaml:"log_requests" json:"log_requests" mapstructure:"log_requests"`                         // Whether to log request metadata
	ExcludedPaths      []string `yaml:"excluded_paths" json:"excluded_paths" mapstructure:"excluded_paths"`                   // Paths to exclude from processing
	IncludeHeaders     []string `yaml:"include_headers" json:"include_headers" mapstructure:"include_headers"`                // Additional headers to include in metadata
}

// RequestMetadata contains extracted request metadata.
// zh: RequestMetadata 包含提取的請求元資料。
type RequestMetadata struct {
	RequestID     string            `json:"request_id"`
	TraceID       string            `json:"trace_id"`
	UserAgent     string            `json:"user_agent"`
	RemoteAddr    string            `json:"remote_addr"`
	ForwardedFor  string            `json:"forwarded_for"`
	Method        string            `json:"method"`
	Path          string            `json:"path"`
	Query         string            `json:"query"`
	ContentType   string            `json:"content_type"`
	ContentLength int64             `json:"content_length"`
	Headers       map[string]string `json:"headers"`
	Timestamp     time.Time         `json:"timestamp"`
}

// contextKey is used for storing request metadata in context.
// zh: contextKey 用於在上下文中儲存請求元資料。
type contextKey string

const (
	RequestMetaContextKey contextKey = "request_metadata"
)

// NewRequestMetaMiddleware creates a new request metadata middleware instance.
// zh: NewRequestMetaMiddleware 建立新的請求元資料中介層實例。
func NewRequestMetaMiddleware(config any) (contracts.Plugin, error) {
	metaConfig := &RequestMetaConfig{
		RequestIDHeader:    "X-Request-ID",
		TraceIDHeader:      "X-Trace-ID",
		UserAgentHeader:    "User-Agent",
		ForwardedForHeader: "X-Forwarded-For",
		GenerateRequestID:  true,
		LogRequests:        false,
		ExcludedPaths: []string{
			"/health",
			"/metrics",
			"/favicon.ico",
		},
		IncludeHeaders: []string{
			"Authorization",
			"Content-Type",
			"Accept",
		},
	}

	// Parse config from the provided config parameter
	if config != nil {
		if err := parseRequestMetaConfig(config, metaConfig); err != nil {
			return nil, fmt.Errorf("failed to parse request metadata config: %w", err)
		}
	}

	return &RequestMetaMiddleware{
		name:        "requestmeta-middleware",
		version:     "1.0.0",
		description: "HTTP request metadata processing middleware",
		config:      metaConfig,
		initialized: false,
	}, nil
}

// parseRequestMetaConfig parses the plugin configuration from various formats
// zh: parseRequestMetaConfig 從各種格式解析插件配置
func parseRequestMetaConfig(config any, target *RequestMetaConfig) error {
	if config == nil {
		return nil
	}

	// Handle map[string]any format
	if configMap, ok := config.(map[string]any); ok {
		if requestIDHeader, exists := configMap["request_id_header"]; exists {
			if str, ok := requestIDHeader.(string); ok {
				target.RequestIDHeader = str
			}
		}
		if traceIDHeader, exists := configMap["trace_id_header"]; exists {
			if str, ok := traceIDHeader.(string); ok {
				target.TraceIDHeader = str
			}
		}
		if userAgentHeader, exists := configMap["user_agent_header"]; exists {
			if str, ok := userAgentHeader.(string); ok {
				target.UserAgentHeader = str
			}
		}
		if forwardedForHeader, exists := configMap["forwarded_for_header"]; exists {
			if str, ok := forwardedForHeader.(string); ok {
				target.ForwardedForHeader = str
			}
		}
		if generateRequestID, exists := configMap["generate_request_id"]; exists {
			if boolVal, ok := generateRequestID.(bool); ok {
				target.GenerateRequestID = boolVal
			}
		}
		if logRequests, exists := configMap["log_requests"]; exists {
			if boolVal, ok := logRequests.(bool); ok {
				target.LogRequests = boolVal
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
		if includeHeaders, exists := configMap["include_headers"]; exists {
			if slice, ok := includeHeaders.([]any); ok {
				target.IncludeHeaders = make([]string, len(slice))
				for i, v := range slice {
					if str, ok := v.(string); ok {
						target.IncludeHeaders[i] = str
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
func (rm *RequestMetaMiddleware) Name() string {
	return rm.name
}

// Version returns the plugin version.
// zh: Version 回傳插件版本。
func (rm *RequestMetaMiddleware) Version() string {
	return rm.version
}

// Description returns the plugin description.
// zh: Description 回傳插件描述。
func (rm *RequestMetaMiddleware) Description() string {
	return rm.description
}

// Init initializes the request metadata middleware.
// zh: Init 初始化請求元資料中介層。
func (rm *RequestMetaMiddleware) Init(config any) error {
	if rm.initialized {
		return nil
	}

	// Validate configuration
	if rm.config.RequestIDHeader == "" {
		return fmt.Errorf("request ID header cannot be empty")
	}

	rm.initialized = true
	return nil
}

// Shutdown shuts down the request metadata middleware.
// zh: Shutdown 關閉請求元資料中介層。
func (rm *RequestMetaMiddleware) Shutdown() error {
	rm.initialized = false
	return nil
}

// LifecycleAware interface implementation
// zh: LifecycleAware 介面實作

// OnRegister is called when the plugin is registered.
// zh: OnRegister 在插件註冊時呼叫。
func (rm *RequestMetaMiddleware) OnRegister() error {
	return nil
}

// OnStart is called when the plugin is started.
// zh: OnStart 在插件啟動時呼叫。
func (rm *RequestMetaMiddleware) OnStart() error {
	if !rm.initialized {
		return fmt.Errorf("request metadata middleware not initialized")
	}
	return nil
}

// OnStop is called when the plugin is stopped.
// zh: OnStop 在插件停止時呼叫。
func (rm *RequestMetaMiddleware) OnStop() error {
	return nil
}

// OnShutdown is called when the plugin is shutdown.
// zh: OnShutdown 在插件關閉時呼叫。
func (rm *RequestMetaMiddleware) OnShutdown() error {
	return rm.Shutdown()
}

// HealthChecker interface implementation
// zh: HealthChecker 介面實作

// CheckHealth checks the health of the request metadata middleware.
// zh: CheckHealth 檢查請求元資料中介層的健康狀況。
func (rm *RequestMetaMiddleware) CheckHealth(ctx context.Context) contracts.HealthStatus {
	status := contracts.HealthStatus{
		Timestamp: time.Now(),
		Details:   make(map[string]any),
	}

	if !rm.initialized {
		status.Status = "unhealthy"
		status.Message = "Request metadata middleware not initialized"
		return status
	}

	status.Status = "healthy"
	status.Message = "Request metadata middleware is healthy"
	status.Details["request_id_header"] = rm.config.RequestIDHeader
	status.Details["trace_id_header"] = rm.config.TraceIDHeader
	status.Details["generate_request_id"] = rm.config.GenerateRequestID
	status.Details["log_requests"] = rm.config.LogRequests
	status.Details["excluded_paths_count"] = len(rm.config.ExcludedPaths)

	return status
}

// GetHealthMetrics returns health metrics for the request metadata middleware.
// zh: GetHealthMetrics 回傳請求元資料中介層的健康指標。
func (rm *RequestMetaMiddleware) GetHealthMetrics() map[string]any {
	return map[string]any{
		"initialized":         rm.initialized,
		"request_id_header":   rm.config.RequestIDHeader,
		"trace_id_header":     rm.config.TraceIDHeader,
		"generate_request_id": rm.config.GenerateRequestID,
		"log_requests":        rm.config.LogRequests,
		"excluded_paths":      rm.config.ExcludedPaths,
		"include_headers":     rm.config.IncludeHeaders,
	}
}

// Middleware functionality
// zh: 中介層功能

// Handler returns the HTTP middleware handler.
// zh: Handler 回傳 HTTP 中介層處理器。
func (rm *RequestMetaMiddleware) Handler() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if path should be excluded
			if rm.shouldExcludePath(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}

			// Extract request metadata
			metadata := rm.extractRequestMetadata(r)

			// Add metadata to request context
			ctx := context.WithValue(r.Context(), RequestMetaContextKey, metadata)
			r = r.WithContext(ctx)

			// Add request ID to response header
			if metadata.RequestID != "" {
				w.Header().Set(rm.config.RequestIDHeader, metadata.RequestID)
			}

			// Log request if enabled
			if rm.config.LogRequests {
				rm.logRequest(metadata)
			}

			// Call next handler
			next.ServeHTTP(w, r)
		})
	}
}

// shouldExcludePath checks if a path should be excluded from processing.
// zh: shouldExcludePath 檢查路徑是否應該排除處理。
func (rm *RequestMetaMiddleware) shouldExcludePath(path string) bool {
	for _, excludedPath := range rm.config.ExcludedPaths {
		if path == excludedPath {
			return true
		}
	}
	return false
}

// extractRequestMetadata extracts metadata from HTTP request.
// zh: extractRequestMetadata 從 HTTP 請求中提取元資料。
func (rm *RequestMetaMiddleware) extractRequestMetadata(r *http.Request) *RequestMetadata {
	metadata := &RequestMetadata{
		Method:        r.Method,
		Path:          r.URL.Path,
		Query:         r.URL.RawQuery,
		ContentType:   r.Header.Get("Content-Type"),
		ContentLength: r.ContentLength,
		RemoteAddr:    r.RemoteAddr,
		Timestamp:     time.Now(),
		Headers:       make(map[string]string),
	}

	// Extract request ID
	metadata.RequestID = r.Header.Get(rm.config.RequestIDHeader)
	if metadata.RequestID == "" && rm.config.GenerateRequestID {
		metadata.RequestID = rm.generateRequestID()
	}

	// Extract trace ID
	metadata.TraceID = r.Header.Get(rm.config.TraceIDHeader)

	// Extract user agent
	metadata.UserAgent = r.Header.Get(rm.config.UserAgentHeader)

	// Extract forwarded for
	metadata.ForwardedFor = r.Header.Get(rm.config.ForwardedForHeader)

	// Extract additional headers
	for _, headerName := range rm.config.IncludeHeaders {
		if value := r.Header.Get(headerName); value != "" {
			metadata.Headers[headerName] = value
		}
	}

	return metadata
}

// generateRequestID generates a unique request ID.
// zh: generateRequestID 產生唯一的請求 ID。
func (rm *RequestMetaMiddleware) generateRequestID() string {
	// Simple implementation using timestamp and random component
	// In production, you might want to use UUID or other more robust methods
	return fmt.Sprintf("req_%d_%d", time.Now().UnixNano(), time.Now().Nanosecond()%10000)
}

// logRequest logs request metadata.
// zh: logRequest 記錄請求元資料。
func (rm *RequestMetaMiddleware) logRequest(metadata *RequestMetadata) {
	// Simple logging implementation
	// In production, you would use a proper logger (otelzap, logrus, etc.)
	fmt.Printf("[RequestMeta] %s %s - RequestID: %s, TraceID: %s, UserAgent: %s\n",
		metadata.Method, metadata.Path, metadata.RequestID, metadata.TraceID, metadata.UserAgent)
}

// Utility functions
// zh: 工具函式

// GetRequestMetadata extracts request metadata from context.
// zh: GetRequestMetadata 從上下文中提取請求元資料。
func GetRequestMetadata(ctx context.Context) (*RequestMetadata, bool) {
	metadata, ok := ctx.Value(RequestMetaContextKey).(*RequestMetadata)
	return metadata, ok
}

// GetRequestID extracts request ID from context.
// zh: GetRequestID 從上下文中提取請求 ID。
func GetRequestID(ctx context.Context) string {
	if metadata, ok := GetRequestMetadata(ctx); ok {
		return metadata.RequestID
	}
	return ""
}

// GetTraceID extracts trace ID from context.
// zh: GetTraceID 從上下文中提取追蹤 ID。
func GetTraceID(ctx context.Context) string {
	if metadata, ok := GetRequestMetadata(ctx); ok {
		return metadata.TraceID
	}
	return ""
}

// Plugin registration
// zh: 插件註冊

// Register registers the request metadata middleware plugin.
// zh: Register 註冊請求元資料中介層插件。
func Register(registry contracts.Registry) error {
	return registry.RegisterPlugin("requestmeta-middleware", NewRequestMetaMiddleware)
}
