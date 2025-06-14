package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"detectviz/internal/platform/registry"
	"detectviz/internal/ports/web"
	"detectviz/internal/ports/web/components"
	"detectviz/internal/ports/web/navtree"
	"detectviz/pkg/platform/contracts"

	// Import WebUI plugins for testing
	systemStatusPlugin "detectviz/plugins/web/pages/system-status"
)

// TestWebUIPluginMounting tests WebUIPlugin mounting functionality.
// zh: TestWebUIPluginMounting 測試 WebUIPlugin 掛載功能。
func TestWebUIPluginMounting(t *testing.T) {
	// Create test registry
	registryManager := registry.NewManager()

	// Create navigation tree builder
	navTreeBuilder := navtree.NewBuilder()

	// Create component registry
	componentRegistry := components.NewRegistry()

	// Create web router
	webRouter := web.NewRouter(navTreeBuilder, componentRegistry)

	t.Run("RegisterWebUIPlugin", func(t *testing.T) {
		// Create system status plugin
		pluginInstance, err := systemStatusPlugin.NewSystemStatusPlugin(map[string]any{
			"title":        "測試系統狀態",
			"refresh_rate": 10,
			"show_memory":  true,
			"show_cpu":     true,
			"show_plugins": true,
		})
		if err != nil {
			t.Fatalf("Failed to create system status plugin: %v", err)
		}

		// Verify plugin implements WebUIPlugin interface
		webUIPlugin, ok := pluginInstance.(contracts.WebUIPlugin)
		if !ok {
			t.Fatalf("Plugin does not implement WebUIPlugin interface")
		}

		// Initialize plugin
		err = webUIPlugin.Init(nil)
		if err != nil {
			t.Fatalf("Failed to initialize plugin: %v", err)
		}

		// Register the WebUI plugin
		err = webRouter.RegisterWebUIPlugin("system-status", webUIPlugin)
		if err != nil {
			t.Fatalf("Failed to register WebUI plugin: %v", err)
		}

		// Verify plugin was registered
		registeredPlugins := webRouter.GetRegisteredWebUIPlugins()
		if _, exists := registeredPlugins["system-status"]; !exists {
			t.Errorf("WebUI plugin was not registered")
		}
	})

	t.Run("RouteRegistration", func(t *testing.T) {
		// Get registered routes
		routes := webRouter.GetRoutes()

		// Expected routes from system status plugin
		expectedRoutes := []string{
			"GET:/system/status",
			"GET:/api/system/status",
			"GET:/system/status/refresh",
		}

		for _, expectedRoute := range expectedRoutes {
			if _, exists := routes[expectedRoute]; !exists {
				t.Errorf("Expected route %s was not registered. Available routes: %v", expectedRoute, getRouteKeys(routes))
			}
		}

		t.Logf("Registered routes: %v", getRouteKeys(routes))
	})

	t.Run("NavTreeRegistration", func(t *testing.T) {
		// Verify navigation node was registered
		navNode, err := navTreeBuilder.GetNode("system-status")
		if err != nil {
			t.Errorf("Navigation node 'system-status' was not registered: %v", err)
		} else {
			if navNode.Title != "測試系統狀態" {
				t.Errorf("Expected nav node title '測試系統狀態', got '%s'", navNode.Title)
			}
			if navNode.URL != "/system/status" {
				t.Errorf("Expected nav node URL '/system/status', got '%s'", navNode.URL)
			}
			if navNode.Icon != "fas fa-heartbeat" {
				t.Errorf("Expected nav node icon 'fas fa-heartbeat', got '%s'", navNode.Icon)
			}
		}
	})

	t.Run("ComponentRegistration", func(t *testing.T) {
		// Verify widgets were registered
		expectedWidgets := []string{
			"system-status-card",
			"memory-usage",
			"plugin-list",
		}

		registeredWidgets := componentRegistry.ListWidgets()
		registeredWidgetMap := make(map[string]bool)
		for _, widget := range registeredWidgets {
			registeredWidgetMap[widget] = true
		}

		for _, widgetName := range expectedWidgets {
			if !registeredWidgetMap[widgetName] {
				t.Errorf("Widget '%s' was not registered", widgetName)
			}
		}

		t.Logf("Registered widgets: %v", registeredWidgets)
	})

	t.Run("HTTPServerIntegration", func(t *testing.T) {
		// Build HTTP handler from web router
		httpHandler := webRouter.BuildHTTPHandler()

		// Create test server
		testServer := httptest.NewServer(httpHandler)
		defer testServer.Close()

		// Test system status page route
		resp, err := http.Get(testServer.URL + "/system/status")
		if err != nil {
			t.Fatalf("Failed to request system status page: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Test system status API route
		resp, err = http.Get(testServer.URL + "/api/system/status")
		if err != nil {
			t.Fatalf("Failed to request system status API: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		// Test content type for API endpoint
		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			t.Errorf("Expected JSON content type, got %s", contentType)
		}
	})

	t.Run("PluginDiscoveryAndRegistration", func(t *testing.T) {
		// Reset web router for clean test
		webRouter2 := web.NewRouter(navtree.NewBuilder(), components.NewRegistry())

		// Register plugin in registry first
		err := systemStatusPlugin.Register(registryManager)
		if err != nil {
			t.Fatalf("Failed to register plugin in registry: %v", err)
		}

		// Discover and register WebUI plugins
		err = webRouter2.DiscoverAndRegisterWebUIPlugins(registryManager)
		if err != nil {
			t.Fatalf("Failed to discover and register WebUI plugins: %v", err)
		}

		// Verify plugin was discovered and registered
		registeredPlugins := webRouter2.GetRegisteredWebUIPlugins()
		if len(registeredPlugins) == 0 {
			t.Errorf("No WebUI plugins were discovered and registered")
		}

		// Note: The system-status plugin implements both Plugin and WebUIPlugin interfaces
		// but may not be automatically discovered as WebUIPlugin in this test setup
		t.Logf("Discovered WebUI plugins: %v", getPluginNames(registeredPlugins))
	})

	t.Run("RouterStatistics", func(t *testing.T) {
		// Get router statistics
		stats := webRouter.GetStats()

		totalRoutes, ok := stats["total_routes"].(int)
		if !ok || totalRoutes == 0 {
			t.Errorf("Expected total_routes > 0, got %v", stats["total_routes"])
		}

		totalWebUIPlugins, ok := stats["total_webui_plugins"].(int)
		if !ok || totalWebUIPlugins == 0 {
			t.Errorf("Expected total_webui_plugins > 0, got %v", stats["total_webui_plugins"])
		}

		t.Logf("Router statistics: %+v", stats)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		// Test duplicate plugin registration
		pluginInstance, _ := systemStatusPlugin.NewSystemStatusPlugin(nil)
		webUIPlugin := pluginInstance.(contracts.WebUIPlugin)
		webUIPlugin.Init(nil)

		err := webRouter.RegisterWebUIPlugin("system-status", webUIPlugin)
		if err == nil {
			t.Errorf("Expected error when registering duplicate plugin")
		}

		expectedErrorMsg := "WebUI plugin system-status already registered"
		if !strings.Contains(err.Error(), expectedErrorMsg) {
			t.Errorf("Expected error message to contain '%s', got '%s'", expectedErrorMsg, err.Error())
		}
	})
}

// TestWebUIPluginUnregistration tests WebUIPlugin unregistration functionality.
// zh: TestWebUIPluginUnregistration 測試 WebUIPlugin 取消註冊功能。
func TestWebUIPluginUnregistration(t *testing.T) {
	// Create web router
	navTreeBuilder := navtree.NewBuilder()
	componentRegistry := components.NewRegistry()
	webRouter := web.NewRouter(navTreeBuilder, componentRegistry)

	// Register a plugin first
	pluginInstance, err := systemStatusPlugin.NewSystemStatusPlugin(nil)
	if err != nil {
		t.Fatalf("Failed to create plugin: %v", err)
	}

	webUIPlugin := pluginInstance.(contracts.WebUIPlugin)
	webUIPlugin.Init(nil)

	err = webRouter.RegisterWebUIPlugin("test-plugin", webUIPlugin)
	if err != nil {
		t.Fatalf("Failed to register plugin: %v", err)
	}

	// Verify plugin is registered
	registeredPlugins := webRouter.GetRegisteredWebUIPlugins()
	if _, exists := registeredPlugins["test-plugin"]; !exists {
		t.Fatalf("Plugin was not registered")
	}

	// Unregister the plugin
	err = webRouter.UnregisterWebUIPlugin("test-plugin")
	if err != nil {
		t.Errorf("Failed to unregister plugin: %v", err)
	}

	// Verify plugin is no longer registered
	registeredPlugins = webRouter.GetRegisteredWebUIPlugins()
	if _, exists := registeredPlugins["test-plugin"]; exists {
		t.Errorf("Plugin was not unregistered")
	}

	// Test unregistering non-existent plugin
	err = webRouter.UnregisterWebUIPlugin("non-existent")
	if err == nil {
		t.Errorf("Expected error when unregistering non-existent plugin")
	}
}

// Helper functions
// zh: 輔助函式

// getRouteKeys returns all route keys from routes map.
// zh: getRouteKeys 從路由 map 回傳所有路由鍵。
func getRouteKeys(routes map[string]*web.Route) []string {
	keys := make([]string, 0, len(routes))
	for key := range routes {
		keys = append(keys, key)
	}
	return keys
}

// getPluginNames returns all plugin names from registered plugins map.
// zh: getPluginNames 從已註冊插件 map 回傳所有插件名稱。
func getPluginNames(plugins map[string]contracts.WebUIPlugin) []string {
	names := make([]string, 0, len(plugins))
	for name := range plugins {
		names = append(names, name)
	}
	return names
}

// MockWebUIPlugin is a simple mock plugin for testing.
// zh: MockWebUIPlugin 是用於測試的簡單模擬插件。
type MockWebUIPlugin struct {
	name        string
	version     string
	description string
	initialized bool
}

// NewMockWebUIPlugin creates a new mock WebUI plugin.
// zh: NewMockWebUIPlugin 建立新的模擬 WebUI 插件。
func NewMockWebUIPlugin(name string) *MockWebUIPlugin {
	return &MockWebUIPlugin{
		name:        name,
		version:     "1.0.0",
		description: fmt.Sprintf("Mock WebUI plugin: %s", name),
		initialized: false,
	}
}

// Plugin interface implementation
func (m *MockWebUIPlugin) Name() string        { return m.name }
func (m *MockWebUIPlugin) Version() string     { return m.version }
func (m *MockWebUIPlugin) Description() string { return m.description }

func (m *MockWebUIPlugin) Init(config any) error {
	m.initialized = true
	return nil
}

func (m *MockWebUIPlugin) Shutdown() error {
	m.initialized = false
	return nil
}

// WebUIPlugin interface implementation
func (m *MockWebUIPlugin) RegisterRoutes(router contracts.WebRouter) error {
	return router.GET(fmt.Sprintf("/%s", m.name), m.handlePage)
}

func (m *MockWebUIPlugin) RegisterNavNodes(navtree contracts.NavTreeBuilder) error {
	node := contracts.NavNode{
		ID:      m.name,
		Title:   m.name,
		Icon:    "fas fa-test",
		URL:     fmt.Sprintf("/%s", m.name),
		Order:   100,
		Visible: true,
		Enabled: true,
	}
	return navtree.AddNode(m.name, node)
}

func (m *MockWebUIPlugin) RegisterComponents(registry contracts.ComponentRegistry) error {
	return registry.RegisterWidget(fmt.Sprintf("%s-widget", m.name), m.testWidget)
}

func (m *MockWebUIPlugin) handlePage(ctx contracts.WebContext) error {
	return ctx.HTML(200, "test-template", map[string]any{
		"title": m.name,
		"data":  "test data",
	})
}

func (m *MockWebUIPlugin) testWidget(ctx contracts.WebContext, params map[string]any) (any, error) {
	return map[string]any{
		"widget": m.name,
		"data":   "test widget data",
	}, nil
}

// TestMockWebUIPlugin tests the mock WebUI plugin functionality.
// zh: TestMockWebUIPlugin 測試模擬 WebUI 插件功能。
func TestMockWebUIPlugin(t *testing.T) {
	// Create mock plugin
	mockPlugin := NewMockWebUIPlugin("test-mock")

	// Test basic plugin interface
	if mockPlugin.Name() != "test-mock" {
		t.Errorf("Expected name 'test-mock', got '%s'", mockPlugin.Name())
	}

	if mockPlugin.Version() != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", mockPlugin.Version())
	}

	// Test initialization
	err := mockPlugin.Init(nil)
	if err != nil {
		t.Errorf("Failed to initialize mock plugin: %v", err)
	}

	if !mockPlugin.initialized {
		t.Errorf("Plugin was not marked as initialized")
	}

	// Test with web router
	navTreeBuilder := navtree.NewBuilder()
	componentRegistry := components.NewRegistry()
	webRouter := web.NewRouter(navTreeBuilder, componentRegistry)

	err = webRouter.RegisterWebUIPlugin("test-mock", mockPlugin)
	if err != nil {
		t.Errorf("Failed to register mock plugin: %v", err)
	}

	// Verify route registration
	routes := webRouter.GetRoutes()
	if _, exists := routes["GET:/test-mock"]; !exists {
		t.Errorf("Mock plugin route was not registered")
	}

	// Verify navigation node registration
	navNode, err := navTreeBuilder.GetNode("test-mock")
	if err != nil {
		t.Errorf("Mock plugin nav node was not registered: %v", err)
	} else {
		if navNode.Title != "test-mock" {
			t.Errorf("Expected nav node title 'test-mock', got '%s'", navNode.Title)
		}
	}

	// Verify widget registration by checking if it exists in the list
	registeredWidgets := componentRegistry.ListWidgets()
	widgetFound := false
	for _, widget := range registeredWidgets {
		if widget == "test-mock-widget" {
			widgetFound = true
			break
		}
	}
	if !widgetFound {
		t.Errorf("Mock plugin widget was not registered")
	}
}
