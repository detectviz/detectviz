package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"detectviz/pkg/platform/contracts"
)

// Router implements REST API routing for DetectViz platform.
// zh: Router 實作 DetectViz 平台的 REST API 路由。
type Router struct {
	registry         contracts.Registry
	lifecycleManager contracts.LifecycleManager
	routes           map[string]*APIRoute
	middleware       []MiddlewareFunc
	mutex            sync.RWMutex
}

// APIRoute represents a REST API route.
// zh: APIRoute 代表 REST API 路由。
type APIRoute struct {
	Method      string
	Path        string
	Handler     APIHandler
	Middleware  []MiddlewareFunc
	Description string
}

// APIHandler defines the interface for API request handlers.
// zh: APIHandler 定義 API 請求處理器介面。
type APIHandler func(ctx *APIContext) error

// MiddlewareFunc defines the interface for API middleware.
// zh: MiddlewareFunc 定義 API 中介層介面。
type MiddlewareFunc func(APIHandler) APIHandler

// APIContext provides context for API requests.
// zh: APIContext 為 API 請求提供上下文。
type APIContext struct {
	Request  *http.Request
	Response http.ResponseWriter
	Params   map[string]string
	User     *contracts.UserInfo
}

// APIResponse represents a standard API response.
// zh: APIResponse 代表標準 API 回應。
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// NewRouter creates a new API router instance.
// zh: NewRouter 建立新的 API 路由實例。
func NewRouter(registry contracts.Registry, lifecycleManager contracts.LifecycleManager) *Router {
	router := &Router{
		registry:         registry,
		lifecycleManager: lifecycleManager,
		routes:           make(map[string]*APIRoute),
		middleware:       make([]MiddlewareFunc, 0),
	}

	// Register default API routes
	router.registerDefaultRoutes()

	return router
}

// registerDefaultRoutes registers the default system API routes.
// zh: registerDefaultRoutes 註冊預設的系統 API 路由。
func (r *Router) registerDefaultRoutes() {
	// Health check endpoint
	r.GET("/health", r.handleHealth, "System health check")

	// Plugin management endpoints
	r.GET("/plugins", r.handleListPlugins, "List all registered plugins")
	r.GET("/plugins/{name}", r.handleGetPlugin, "Get plugin details")
	r.GET("/plugins/{name}/health", r.handlePluginHealth, "Get plugin health status")

	// System status endpoints
	r.GET("/status", r.handleSystemStatus, "Get system status")
	r.GET("/status/lifecycle", r.handleLifecycleStatus, "Get lifecycle manager status")

	// Registry endpoints
	r.GET("/registry/stats", r.handleRegistryStats, "Get registry statistics")

	// Configuration endpoints
	r.GET("/config/plugins", r.handlePluginConfigs, "Get plugin configurations")
}

// HTTP method handlers
// zh: HTTP 方法處理器

// GET registers a GET route.
// zh: GET 註冊 GET 路由。
func (r *Router) GET(path string, handler APIHandler, description string) {
	r.addRoute("GET", path, handler, nil, description)
}

// POST registers a POST route.
// zh: POST 註冊 POST 路由。
func (r *Router) POST(path string, handler APIHandler, description string) {
	r.addRoute("POST", path, handler, nil, description)
}

// PUT registers a PUT route.
// zh: PUT 註冊 PUT 路由。
func (r *Router) PUT(path string, handler APIHandler, description string) {
	r.addRoute("PUT", path, handler, nil, description)
}

// DELETE registers a DELETE route.
// zh: DELETE 註冊 DELETE 路由。
func (r *Router) DELETE(path string, handler APIHandler, description string) {
	r.addRoute("DELETE", path, handler, nil, description)
}

// addRoute adds a route to the router.
// zh: addRoute 向路由器添加路由。
func (r *Router) addRoute(method, path string, handler APIHandler, middleware []MiddlewareFunc, description string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	routeKey := fmt.Sprintf("%s:%s", method, path)
	route := &APIRoute{
		Method:      method,
		Path:        path,
		Handler:     handler,
		Middleware:  middleware,
		Description: description,
	}

	r.routes[routeKey] = route
}

// Use adds global middleware to the router.
// zh: Use 向路由器添加全域中介層。
func (r *Router) Use(middleware MiddlewareFunc) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.middleware = append(r.middleware, middleware)
}

// BuildHTTPHandler converts the router to an http.Handler.
// zh: BuildHTTPHandler 將路由器轉換為 http.Handler。
func (r *Router) BuildHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, route := range r.routes {
		// Build handler chain with middleware
		handler := route.Handler

		// Apply route-specific middleware
		for i := len(route.Middleware) - 1; i >= 0; i-- {
			handler = route.Middleware[i](handler)
		}

		// Apply global middleware
		for i := len(r.middleware) - 1; i >= 0; i-- {
			handler = r.middleware[i](handler)
		}

		// Wrap with HTTP handler
		httpHandler := r.wrapAPIHandler(handler)

		// Register with pattern that includes method
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		mux.HandleFunc(pattern, httpHandler)
	}

	return mux
}

// wrapAPIHandler wraps an APIHandler to work with http.HandlerFunc.
// zh: wrapAPIHandler 包裝 APIHandler 以與 http.HandlerFunc 相容。
func (r *Router) wrapAPIHandler(handler APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Create API context
		ctx := &APIContext{
			Request:  req,
			Response: w,
			Params:   r.extractParams(req),
		}

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the API handler
		if err := handler(ctx); err != nil {
			ctx.ErrorResponse(http.StatusInternalServerError, err.Error())
		}
	}
}

// extractParams extracts URL parameters from the request.
// zh: extractParams 從請求中提取 URL 參數。
func (r *Router) extractParams(req *http.Request) map[string]string {
	params := make(map[string]string)

	// Simple parameter extraction from URL path
	// In a real implementation, you would use a more sophisticated router
	path := req.URL.Path
	segments := strings.Split(path, "/")

	// Extract parameters based on route patterns
	// This is a simplified implementation
	for i, segment := range segments {
		if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			paramName := segment[1 : len(segment)-1]
			if i < len(segments) {
				params[paramName] = segments[i]
			}
		}
	}

	return params
}

// Default API handlers
// zh: 預設 API 處理器

// handleHealth handles the health check endpoint.
// zh: handleHealth 處理健康檢查端點。
func (r *Router) handleHealth(ctx *APIContext) error {
	healthData := map[string]interface{}{
		"status":    "healthy",
		"timestamp": "2024-01-01T00:00:00Z", // Should use actual timestamp
		"version":   "1.0.0",
		"uptime":    "0s", // Should calculate actual uptime
	}

	return ctx.JSONResponse(http.StatusOK, healthData)
}

// handleListPlugins handles listing all plugins.
// zh: handleListPlugins 處理列出所有插件。
func (r *Router) handleListPlugins(ctx *APIContext) error {
	pluginNames := r.registry.ListPlugins()

	plugins := make([]map[string]interface{}, 0, len(pluginNames))
	for _, name := range pluginNames {
		plugin, err := r.registry.GetPlugin(name)
		if err != nil {
			continue
		}

		pluginInfo := map[string]interface{}{
			"name":        plugin.Name(),
			"version":     plugin.Version(),
			"description": plugin.Description(),
		}

		// Add metadata if available
		if metadata, err := r.registry.GetPluginMetadata(name); err == nil {
			pluginInfo["type"] = metadata.Type
			pluginInfo["category"] = metadata.Category
			pluginInfo["enabled"] = metadata.Enabled
		}

		plugins = append(plugins, pluginInfo)
	}

	return ctx.JSONResponse(http.StatusOK, plugins)
}

// handleGetPlugin handles getting a specific plugin.
// zh: handleGetPlugin 處理取得特定插件。
func (r *Router) handleGetPlugin(ctx *APIContext) error {
	name := ctx.Params["name"]
	if name == "" {
		return ctx.ErrorResponse(http.StatusBadRequest, "Plugin name is required")
	}

	plugin, err := r.registry.GetPlugin(name)
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotFound, fmt.Sprintf("Plugin %s not found", name))
	}

	pluginInfo := map[string]interface{}{
		"name":        plugin.Name(),
		"version":     plugin.Version(),
		"description": plugin.Description(),
	}

	// Add metadata if available
	if metadata, err := r.registry.GetPluginMetadata(name); err == nil {
		pluginInfo["type"] = metadata.Type
		pluginInfo["category"] = metadata.Category
		pluginInfo["enabled"] = metadata.Enabled
		pluginInfo["dependencies"] = metadata.Dependencies
		pluginInfo["config"] = metadata.Config
	}

	return ctx.JSONResponse(http.StatusOK, pluginInfo)
}

// handlePluginHealth handles plugin health check.
// zh: handlePluginHealth 處理插件健康檢查。
func (r *Router) handlePluginHealth(ctx *APIContext) error {
	name := ctx.Params["name"]
	if name == "" {
		return ctx.ErrorResponse(http.StatusBadRequest, "Plugin name is required")
	}

	plugin, err := r.registry.GetPlugin(name)
	if err != nil {
		return ctx.ErrorResponse(http.StatusNotFound, fmt.Sprintf("Plugin %s not found", name))
	}

	// Check if plugin implements HealthChecker
	if healthChecker, ok := plugin.(contracts.HealthChecker); ok {
		healthStatus := healthChecker.CheckHealth(ctx.Request.Context())
		return ctx.JSONResponse(http.StatusOK, healthStatus)
	}

	// Default health status
	healthStatus := map[string]interface{}{
		"status":    "unknown",
		"message":   "Plugin does not implement health checking",
		"timestamp": "2024-01-01T00:00:00Z", // Should use actual timestamp
	}

	return ctx.JSONResponse(http.StatusOK, healthStatus)
}

// handleSystemStatus handles system status endpoint.
// zh: handleSystemStatus 處理系統狀態端點。
func (r *Router) handleSystemStatus(ctx *APIContext) error {
	status := map[string]interface{}{
		"platform": map[string]interface{}{
			"name":    "DetectViz",
			"version": "1.0.0",
			"status":  "running",
		},
		"plugins": map[string]interface{}{
			"total":   len(r.registry.ListPlugins()),
			"enabled": len(r.registry.ListPlugins()), // Simplified
		},
		"timestamp": "2024-01-01T00:00:00Z", // Should use actual timestamp
	}

	return ctx.JSONResponse(http.StatusOK, status)
}

// handleLifecycleStatus handles lifecycle manager status.
// zh: handleLifecycleStatus 處理生命週期管理器狀態。
func (r *Router) handleLifecycleStatus(ctx *APIContext) error {
	if r.lifecycleManager == nil {
		return ctx.ErrorResponse(http.StatusServiceUnavailable, "Lifecycle manager not available")
	}

	status := map[string]interface{}{
		"status":    r.lifecycleManager.GetStatus(),
		"timestamp": "2024-01-01T00:00:00Z", // Should use actual timestamp
	}

	return ctx.JSONResponse(http.StatusOK, status)
}

// handleRegistryStats handles registry statistics.
// zh: handleRegistryStats 處理註冊表統計。
func (r *Router) handleRegistryStats(ctx *APIContext) error {
	pluginNames := r.registry.ListPlugins()

	stats := map[string]interface{}{
		"total_plugins": len(pluginNames),
		"plugins":       pluginNames,
		"timestamp":     "2024-01-01T00:00:00Z", // Should use actual timestamp
	}

	return ctx.JSONResponse(http.StatusOK, stats)
}

// handlePluginConfigs handles plugin configurations.
// zh: handlePluginConfigs 處理插件配置。
func (r *Router) handlePluginConfigs(ctx *APIContext) error {
	pluginNames := r.registry.ListPlugins()

	configs := make(map[string]interface{})
	for _, name := range pluginNames {
		if metadata, err := r.registry.GetPluginMetadata(name); err == nil {
			configs[name] = metadata.Config
		}
	}

	return ctx.JSONResponse(http.StatusOK, configs)
}

// APIContext helper methods
// zh: APIContext 輔助方法

// JSONResponse sends a JSON response.
// zh: JSONResponse 發送 JSON 回應。
func (ctx *APIContext) JSONResponse(code int, data interface{}) error {
	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.WriteHeader(code)

	response := APIResponse{
		Success: code < 400,
		Data:    data,
	}

	encoder := json.NewEncoder(ctx.Response)
	return encoder.Encode(response)
}

// ErrorResponse sends an error response.
// zh: ErrorResponse 發送錯誤回應。
func (ctx *APIContext) ErrorResponse(code int, message string) error {
	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.WriteHeader(code)

	response := APIResponse{
		Success: false,
		Error:   message,
	}

	encoder := json.NewEncoder(ctx.Response)
	return encoder.Encode(response)
}

// MessageResponse sends a message response.
// zh: MessageResponse 發送訊息回應。
func (ctx *APIContext) MessageResponse(code int, message string) error {
	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.WriteHeader(code)

	response := APIResponse{
		Success: code < 400,
		Message: message,
	}

	encoder := json.NewEncoder(ctx.Response)
	return encoder.Encode(response)
}

// GetRoutes returns all registered routes.
// zh: GetRoutes 回傳所有已註冊的路由。
func (r *Router) GetRoutes() map[string]*APIRoute {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make(map[string]*APIRoute)
	for key, route := range r.routes {
		result[key] = route
	}
	return result
}

// GetStats returns router statistics.
// zh: GetStats 回傳路由器統計資訊。
func (r *Router) GetStats() map[string]interface{} {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	stats := map[string]interface{}{
		"total_routes":     len(r.routes),
		"total_middleware": len(r.middleware),
		"routes_by_method": r.getRoutesByMethod(),
	}

	return stats
}

// getRoutesByMethod returns route count by HTTP method.
// zh: getRoutesByMethod 回傳按 HTTP 方法分類的路由數量。
func (r *Router) getRoutesByMethod() map[string]int {
	methodCount := make(map[string]int)

	for _, route := range r.routes {
		methodCount[route.Method]++
	}

	return methodCount
}
