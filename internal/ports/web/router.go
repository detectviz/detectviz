package web

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"detectviz/pkg/platform/contracts"
	"detectviz/pkg/shared/log"
)

// Router implements WebRouter interface for plugin route registration.
// zh: Router 實作 WebRouter 介面用於插件路由註冊。
type Router struct {
	routes       map[string]*Route
	groups       map[string]*RouterGroup
	navTree      contracts.NavTreeBuilder
	componentReg contracts.ComponentRegistry
	webUIPlugins map[string]contracts.WebUIPlugin
	mutex        sync.RWMutex
}

// Route represents a registered route.
// zh: Route 代表已註冊的路由。
type Route struct {
	Method  string
	Path    string
	Handler contracts.WebHandler
	Group   string
}

// RouterGroup represents a group of routes with common prefix.
// zh: RouterGroup 代表具有共同前綴的路由群組。
type RouterGroup struct {
	prefix     string
	middleware []func(contracts.WebHandler) contracts.WebHandler
	router     *Router
}

// NewRouter creates a new web router instance.
// zh: NewRouter 建立新的 Web 路由實例。
func NewRouter(navTree contracts.NavTreeBuilder, componentReg contracts.ComponentRegistry) *Router {
	return &Router{
		routes:       make(map[string]*Route),
		groups:       make(map[string]*RouterGroup),
		navTree:      navTree,
		componentReg: componentReg,
		webUIPlugins: make(map[string]contracts.WebUIPlugin),
	}
}

// WebRouter interface implementation
// zh: WebRouter 介面實作

// GET registers a GET route.
// zh: GET 註冊 GET 路由。
func (r *Router) GET(path string, handler contracts.WebHandler) error {
	return r.addRoute("GET", path, handler, "")
}

// POST registers a POST route.
// zh: POST 註冊 POST 路由。
func (r *Router) POST(path string, handler contracts.WebHandler) error {
	return r.addRoute("POST", path, handler, "")
}

// PUT registers a PUT route.
// zh: PUT 註冊 PUT 路由。
func (r *Router) PUT(path string, handler contracts.WebHandler) error {
	return r.addRoute("PUT", path, handler, "")
}

// DELETE registers a DELETE route.
// zh: DELETE 註冊 DELETE 路由。
func (r *Router) DELETE(path string, handler contracts.WebHandler) error {
	return r.addRoute("DELETE", path, handler, "")
}

// Group creates a new router group with the given prefix.
// zh: Group 使用給定前綴建立新的路由群組。
func (r *Router) Group(prefix string) contracts.WebRouter {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	group := &RouterGroup{
		prefix: prefix,
		router: r,
	}

	r.groups[prefix] = group
	return group
}

// RouterGroup implementation
// zh: RouterGroup 實作

// GET registers a GET route in the group.
// zh: GET 在群組中註冊 GET 路由。
func (rg *RouterGroup) GET(path string, handler contracts.WebHandler) error {
	fullPath := rg.prefix + path
	return rg.router.addRoute("GET", fullPath, handler, rg.prefix)
}

// POST registers a POST route in the group.
// zh: POST 在群組中註冊 POST 路由。
func (rg *RouterGroup) POST(path string, handler contracts.WebHandler) error {
	fullPath := rg.prefix + path
	return rg.router.addRoute("POST", fullPath, handler, rg.prefix)
}

// PUT registers a PUT route in the group.
// zh: PUT 在群組中註冊 PUT 路由。
func (rg *RouterGroup) PUT(path string, handler contracts.WebHandler) error {
	fullPath := rg.prefix + path
	return rg.router.addRoute("PUT", fullPath, handler, rg.prefix)
}

// DELETE registers a DELETE route in the group.
// zh: DELETE 在群組中註冊 DELETE 路由。
func (rg *RouterGroup) DELETE(path string, handler contracts.WebHandler) error {
	fullPath := rg.prefix + path
	return rg.router.addRoute("DELETE", fullPath, handler, rg.prefix)
}

// Group creates a nested group.
// zh: Group 建立嵌套群組。
func (rg *RouterGroup) Group(prefix string) contracts.WebRouter {
	fullPrefix := rg.prefix + prefix
	return rg.router.Group(fullPrefix)
}

// Plugin management methods
// zh: 插件管理方法

// RegisterWebUIPlugin registers a WebUIPlugin with the router.
// zh: RegisterWebUIPlugin 註冊 WebUIPlugin 到路由器。
func (r *Router) RegisterWebUIPlugin(name string, plugin contracts.WebUIPlugin) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.webUIPlugins[name]; exists {
		ctx := context.Background()
		log.L(ctx).Warn("WebUI plugin already registered, overwriting", "plugin", name)
	}

	r.webUIPlugins[name] = plugin

	// Register routes
	if err := plugin.RegisterRoutes(r); err != nil {
		ctx := context.Background()
		log.L(ctx).Error("Failed to register WebUI plugin", "plugin", name, "error", err)
		return err
	}

	// Register navigation tree items
	if err := plugin.RegisterNavNodes(r.navTree); err != nil {
		ctx := context.Background()
		log.L(ctx).Error("Failed to register navigation for WebUI plugin", "plugin", name, "error", err)
		return err
	}

	// Register components
	if err := plugin.RegisterComponents(r.componentReg); err != nil {
		ctx := context.Background()
		log.L(ctx).Error("Failed to register components for WebUI plugin", "plugin", name, "error", err)
		return err
	}

	ctx := context.Background()
	log.L(ctx).Info("Registered WebUI plugin", "plugin", name)
	return nil
}

// DiscoverAndRegisterWebUIPlugins scans registry for WebUI plugins and registers them.
// zh: DiscoverAndRegisterWebUIPlugins 掃描註冊表中的 WebUI 插件並註冊它們。
func (r *Router) DiscoverAndRegisterWebUIPlugins(registry contracts.Registry) error {
	pluginNames := registry.ListPlugins()

	for _, name := range pluginNames {
		plugin, err := registry.GetPlugin(name)
		if err != nil {
			continue // Skip plugins that can't be loaded
		}

		// Check if plugin implements WebUIPlugin interface
		if webUIPlugin, ok := plugin.(contracts.WebUIPlugin); ok {
			if err := r.RegisterWebUIPlugin(name, webUIPlugin); err != nil {
				// Log error but continue with other plugins
				fmt.Printf("Failed to register WebUI plugin %s: %v\n", name, err)
				continue
			}
			fmt.Printf("Registered WebUI plugin: %s\n", name)
		}
	}

	return nil
}

// GetRegisteredWebUIPlugins returns all registered WebUI plugins.
// zh: GetRegisteredWebUIPlugins 回傳所有已註冊的 WebUI 插件。
func (r *Router) GetRegisteredWebUIPlugins() map[string]contracts.WebUIPlugin {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make(map[string]contracts.WebUIPlugin)
	for name, plugin := range r.webUIPlugins {
		result[name] = plugin
	}
	return result
}

// Route management methods
// zh: 路由管理方法

// addRoute adds a route to the router.
// zh: addRoute 向路由器添加路由。
func (r *Router) addRoute(method, path string, handler contracts.WebHandler, group string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	routeKey := fmt.Sprintf("%s:%s", method, path)
	if _, exists := r.routes[routeKey]; exists {
		return fmt.Errorf("route %s %s already registered", method, path)
	}

	route := &Route{
		Method:  method,
		Path:    path,
		Handler: handler,
		Group:   group,
	}

	r.routes[routeKey] = route
	return nil
}

// GetRoutes returns all registered routes.
// zh: GetRoutes 回傳所有已註冊的路由。
func (r *Router) GetRoutes() map[string]*Route {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	result := make(map[string]*Route)
	for key, route := range r.routes {
		result[key] = route
	}
	return result
}

// GetRoutesByGroup returns routes for a specific group.
// zh: GetRoutesByGroup 回傳特定群組的路由。
func (r *Router) GetRoutesByGroup(group string) []*Route {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var routes []*Route
	for _, route := range r.routes {
		if route.Group == group {
			routes = append(routes, route)
		}
	}
	return routes
}

// HTTP integration methods
// zh: HTTP 整合方法

// BuildHTTPHandler converts the router to an http.Handler.
// zh: BuildHTTPHandler 將路由器轉換為 http.Handler。
func (r *Router) BuildHTTPHandler() http.Handler {
	mux := http.NewServeMux()

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, route := range r.routes {
		// Wrap WebHandler to http.HandlerFunc
		httpHandler := r.wrapWebHandler(route.Handler)

		// Register with pattern that includes method
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		mux.HandleFunc(pattern, httpHandler)
	}

	return mux
}

// wrapWebHandler wraps a WebHandler to work with http.HandlerFunc.
// zh: wrapWebHandler 包裝 WebHandler 以與 http.HandlerFunc 相容。
func (r *Router) wrapWebHandler(handler contracts.WebHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Create WebContext from http request
		ctx := NewWebContext(w, req)

		// Call the WebHandler
		if err := handler(ctx); err != nil {
			// Handle error - could be customized based on error type
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Statistics and debugging methods
// zh: 統計與除錯方法

// GetStats returns router statistics.
// zh: GetStats 回傳路由器統計資訊。
func (r *Router) GetStats() map[string]any {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	stats := map[string]any{
		"total_routes":        len(r.routes),
		"total_groups":        len(r.groups),
		"total_webui_plugins": len(r.webUIPlugins),
		"routes_by_method":    r.getRoutesByMethod(),
		"routes_by_group":     r.getRoutesByGroupStats(),
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

// getRoutesByGroupStats returns route count by group.
// zh: getRoutesByGroupStats 回傳按群組分類的路由數量。
func (r *Router) getRoutesByGroupStats() map[string]int {
	groupCount := make(map[string]int)

	for _, route := range r.routes {
		if route.Group == "" {
			groupCount["root"]++
		} else {
			groupCount[route.Group]++
		}
	}

	return groupCount
}

// UnregisterWebUIPlugin unregisters a WebUI plugin.
// zh: UnregisterWebUIPlugin 取消註冊 WebUI 插件。
func (r *Router) UnregisterWebUIPlugin(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	plugin, exists := r.webUIPlugins[name]
	if !exists {
		return fmt.Errorf("WebUI plugin %s not found", name)
	}

	// TODO: Remove routes, nav nodes, and components registered by this plugin
	// This would require tracking which resources belong to which plugin

	delete(r.webUIPlugins, name)

	// Call plugin shutdown if it supports lifecycle
	if lifecyclePlugin, ok := plugin.(contracts.LifecycleAware); ok {
		return lifecyclePlugin.OnShutdown()
	}

	return nil
}
