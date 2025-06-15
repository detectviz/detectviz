package fake

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"detectviz/pkg/platform/contracts"
)

// TestFakeRegistry tests the fake registry functionality.
// zh: TestFakeRegistry 測試假註冊表功能。
func TestFakeRegistry(t *testing.T) {
	t.Run("NewFakeRegistry", func(t *testing.T) {
		registry := NewFakeRegistry()
		if registry == nil {
			t.Fatal("NewFakeRegistry should not return nil")
		}

		// Check initial state
		plugins := registry.ListPlugins()
		if len(plugins) != 0 {
			t.Errorf("Expected 0 plugins initially, got %d", len(plugins))
		}

		stats := registry.GetStats()
		if stats["total_plugins"] != 0 {
			t.Errorf("Expected 0 total plugins, got %d", stats["total_plugins"])
		}

		t.Log("NewFakeRegistry test passed")
	})

	t.Run("RegisterPlugin", func(t *testing.T) {
		registry := NewFakeRegistry()
		plugin := NewFakePlugin("test-plugin", "1.0.0", "Test plugin")
		metadata := &contracts.PluginMetadata{
			Name:        "test-plugin",
			Version:     "1.0.0",
			Description: "Test plugin",
			Author:      "Test Author",
		}

		// Register plugin
		err := registry.Register(plugin, metadata)
		if err != nil {
			t.Fatalf("Failed to register plugin: %v", err)
		}

		// Check plugin is registered
		plugins := registry.ListPlugins()
		if len(plugins) != 1 {
			t.Errorf("Expected 1 plugin, got %d", len(plugins))
		}

		if plugins[0] != "test-plugin" {
			t.Errorf("Expected plugin name 'test-plugin', got '%s'", plugins[0])
		}

		// Try to register same plugin again (should fail)
		err = registry.Register(plugin, metadata)
		if err == nil {
			t.Error("Expected error when registering duplicate plugin")
		}

		t.Log("RegisterPlugin test passed")
	})

	t.Run("ResolvePlugin", func(t *testing.T) {
		registry := NewFakeRegistry()
		plugin := NewFakePlugin("test-plugin", "1.0.0", "Test plugin")
		metadata := &contracts.PluginMetadata{
			Name:    "test-plugin",
			Version: "1.0.0",
		}

		// Register plugin
		err := registry.Register(plugin, metadata)
		if err != nil {
			t.Fatalf("Failed to register plugin: %v", err)
		}

		// Resolve plugin
		resolvedPlugin, err := registry.Resolve("test-plugin")
		if err != nil {
			t.Fatalf("Failed to resolve plugin: %v", err)
		}

		if resolvedPlugin.Name() != "test-plugin" {
			t.Errorf("Expected resolved plugin name 'test-plugin', got '%s'", resolvedPlugin.Name())
		}

		// Try to resolve non-existent plugin
		_, err = registry.Resolve("non-existent")
		if err == nil {
			t.Error("Expected error when resolving non-existent plugin")
		}

		t.Log("ResolvePlugin test passed")
	})

	t.Run("GetMetadata", func(t *testing.T) {
		registry := NewFakeRegistry()
		plugin := NewFakePlugin("test-plugin", "1.0.0", "Test plugin")
		metadata := &contracts.PluginMetadata{
			Name:        "test-plugin",
			Version:     "1.0.0",
			Description: "Test plugin",
			Author:      "Test Author",
		}

		// Register plugin
		err := registry.Register(plugin, metadata)
		if err != nil {
			t.Fatalf("Failed to register plugin: %v", err)
		}

		// Get metadata
		retrievedMetadata, err := registry.GetMetadata("test-plugin")
		if err != nil {
			t.Fatalf("Failed to get metadata: %v", err)
		}

		if retrievedMetadata.Name != "test-plugin" {
			t.Errorf("Expected metadata name 'test-plugin', got '%s'", retrievedMetadata.Name)
		}

		if retrievedMetadata.Author != "Test Author" {
			t.Errorf("Expected metadata author 'Test Author', got '%s'", retrievedMetadata.Author)
		}

		// Try to get metadata for non-existent plugin
		_, err = registry.GetMetadata("non-existent")
		if err == nil {
			t.Error("Expected error when getting metadata for non-existent plugin")
		}

		t.Log("GetMetadata test passed")
	})

	t.Run("InitializePlugin", func(t *testing.T) {
		registry := NewFakeRegistry()
		plugin := NewFakePlugin("test-plugin", "1.0.0", "Test plugin")
		metadata := &contracts.PluginMetadata{
			Name:    "test-plugin",
			Version: "1.0.0",
		}

		// Register plugin
		err := registry.Register(plugin, metadata)
		if err != nil {
			t.Fatalf("Failed to register plugin: %v", err)
		}

		// Initialize plugin
		config := map[string]any{
			"setting1": "value1",
			"setting2": 42,
		}

		err = registry.InitializePlugin("test-plugin", config)
		if err != nil {
			t.Fatalf("Failed to initialize plugin: %v", err)
		}

		// Check plugin is initialized
		if !plugin.IsInitialized() {
			t.Error("Plugin should be initialized")
		}

		// Check config is stored
		storedConfig, exists := registry.GetPluginConfig("test-plugin")
		if !exists {
			t.Error("Plugin config should exist")
		}

		if storedConfig["setting1"] != "value1" {
			t.Errorf("Expected setting1 'value1', got '%v'", storedConfig["setting1"])
		}

		if storedConfig["setting2"] != 42 {
			t.Errorf("Expected setting2 42, got %v", storedConfig["setting2"])
		}

		t.Log("InitializePlugin test passed")
	})

	t.Run("StartStopPlugin", func(t *testing.T) {
		registry := NewFakeRegistry()
		plugin := NewFakePlugin("test-plugin", "1.0.0", "Test plugin")
		metadata := &contracts.PluginMetadata{
			Name:    "test-plugin",
			Version: "1.0.0",
		}

		// Register plugin
		err := registry.Register(plugin, metadata)
		if err != nil {
			t.Fatalf("Failed to register plugin: %v", err)
		}

		// Check initial state
		if registry.IsPluginStarted("test-plugin") {
			t.Error("Plugin should not be started initially")
		}

		// Start plugin
		err = registry.StartPlugin("test-plugin")
		if err != nil {
			t.Fatalf("Failed to start plugin: %v", err)
		}

		if !registry.IsPluginStarted("test-plugin") {
			t.Error("Plugin should be started")
		}

		// Stop plugin
		err = registry.StopPlugin("test-plugin")
		if err != nil {
			t.Fatalf("Failed to stop plugin: %v", err)
		}

		if registry.IsPluginStarted("test-plugin") {
			t.Error("Plugin should be stopped")
		}

		t.Log("StartStopPlugin test passed")
	})

	t.Run("SimulateError", func(t *testing.T) {
		registry := NewFakeRegistry()
		plugin := NewFakePlugin("test-plugin", "1.0.0", "Test plugin")
		metadata := &contracts.PluginMetadata{
			Name:    "test-plugin",
			Version: "1.0.0",
		}

		// Register plugin
		err := registry.Register(plugin, metadata)
		if err != nil {
			t.Fatalf("Failed to register plugin: %v", err)
		}

		// Simulate error
		testError := errors.New("simulated error")
		registry.SimulateError("test-plugin", testError)

		// Try to resolve plugin (should fail)
		_, err = registry.Resolve("test-plugin")
		if err == nil {
			t.Error("Expected error when resolving plugin with simulated error")
		}

		if err.Error() != "simulated error" {
			t.Errorf("Expected 'simulated error', got '%s'", err.Error())
		}

		// Clear error
		registry.ClearError("test-plugin")

		// Try to resolve plugin again (should succeed)
		_, err = registry.Resolve("test-plugin")
		if err != nil {
			t.Errorf("Expected no error after clearing, got: %v", err)
		}

		t.Log("SimulateError test passed")
	})

	t.Run("GetStats", func(t *testing.T) {
		registry := NewFakeRegistry()

		// Add multiple plugins
		for i := 0; i < 3; i++ {
			plugin := NewFakePlugin(
				fmt.Sprintf("plugin-%d", i),
				"1.0.0",
				fmt.Sprintf("Test plugin %d", i),
			)
			metadata := &contracts.PluginMetadata{
				Name:    fmt.Sprintf("plugin-%d", i),
				Version: "1.0.0",
			}

			err := registry.Register(plugin, metadata)
			if err != nil {
				t.Fatalf("Failed to register plugin %d: %v", i, err)
			}
		}

		// Start some plugins
		err := registry.StartPlugin("plugin-0")
		if err != nil {
			t.Fatalf("Failed to start plugin-0: %v", err)
		}

		err = registry.StartPlugin("plugin-1")
		if err != nil {
			t.Fatalf("Failed to start plugin-1: %v", err)
		}

		// Check stats
		stats := registry.GetStats()
		if stats["total_plugins"] != 3 {
			t.Errorf("Expected 3 total plugins, got %d", stats["total_plugins"])
		}

		if stats["started_plugins"] != 2 {
			t.Errorf("Expected 2 started plugins, got %d", stats["started_plugins"])
		}

		if stats["stopped_plugins"] != 1 {
			t.Errorf("Expected 1 stopped plugin, got %d", stats["stopped_plugins"])
		}

		t.Log("GetStats test passed")
	})

	t.Run("Reset", func(t *testing.T) {
		registry := NewFakeRegistry()
		plugin := NewFakePlugin("test-plugin", "1.0.0", "Test plugin")
		metadata := &contracts.PluginMetadata{
			Name:    "test-plugin",
			Version: "1.0.0",
		}

		// Register and start plugin
		err := registry.Register(plugin, metadata)
		if err != nil {
			t.Fatalf("Failed to register plugin: %v", err)
		}

		err = registry.StartPlugin("test-plugin")
		if err != nil {
			t.Fatalf("Failed to start plugin: %v", err)
		}

		// Check plugin exists
		plugins := registry.ListPlugins()
		if len(plugins) != 1 {
			t.Errorf("Expected 1 plugin before reset, got %d", len(plugins))
		}

		// Reset registry
		registry.Reset()

		// Check registry is empty
		plugins = registry.ListPlugins()
		if len(plugins) != 0 {
			t.Errorf("Expected 0 plugins after reset, got %d", len(plugins))
		}

		stats := registry.GetStats()
		if stats["total_plugins"] != 0 {
			t.Errorf("Expected 0 total plugins after reset, got %d", stats["total_plugins"])
		}

		t.Log("Reset test passed")
	})
}

// TestFakePlugin tests the fake plugin functionality.
// zh: TestFakePlugin 測試假插件功能。
func TestFakePlugin(t *testing.T) {
	t.Run("NewFakePlugin", func(t *testing.T) {
		plugin := NewFakePlugin("test-plugin", "1.0.0", "Test plugin")
		if plugin == nil {
			t.Fatal("NewFakePlugin should not return nil")
		}

		if plugin.Name() != "test-plugin" {
			t.Errorf("Expected name 'test-plugin', got '%s'", plugin.Name())
		}

		if plugin.Version() != "1.0.0" {
			t.Errorf("Expected version '1.0.0', got '%s'", plugin.Version())
		}

		if plugin.Description() != "Test plugin" {
			t.Errorf("Expected description 'Test plugin', got '%s'", plugin.Description())
		}

		if plugin.IsInitialized() {
			t.Error("Plugin should not be initialized initially")
		}

		t.Log("NewFakePlugin test passed")
	})

	t.Run("InitShutdown", func(t *testing.T) {
		plugin := NewFakePlugin("test-plugin", "1.0.0", "Test plugin")

		// Initialize plugin
		config := map[string]any{
			"setting": "value",
		}

		err := plugin.Init(config)
		if err != nil {
			t.Fatalf("Failed to initialize plugin: %v", err)
		}

		if !plugin.IsInitialized() {
			t.Error("Plugin should be initialized")
		}

		// Check config is stored
		storedConfig := plugin.GetConfig()
		if storedConfig["setting"] != "value" {
			t.Errorf("Expected setting 'value', got '%v'", storedConfig["setting"])
		}

		// Shutdown plugin
		err = plugin.Shutdown()
		if err != nil {
			t.Fatalf("Failed to shutdown plugin: %v", err)
		}

		if plugin.IsInitialized() {
			t.Error("Plugin should not be initialized after shutdown")
		}

		t.Log("InitShutdown test passed")
	})

	t.Run("LifecycleHandlers", func(t *testing.T) {
		plugin := NewFakePlugin("test-plugin", "1.0.0", "Test plugin")

		// Track lifecycle calls
		var calls []string

		plugin.SetLifecycleHandlers(
			func() error {
				calls = append(calls, "register")
				return nil
			},
			func() error {
				calls = append(calls, "start")
				return nil
			},
			func() error {
				calls = append(calls, "stop")
				return nil
			},
			func() error {
				calls = append(calls, "shutdown")
				return nil
			},
		)

		// Call lifecycle methods
		err := plugin.OnRegister()
		if err != nil {
			t.Fatalf("OnRegister failed: %v", err)
		}

		err = plugin.OnStart()
		if err != nil {
			t.Fatalf("OnStart failed: %v", err)
		}

		err = plugin.OnStop()
		if err != nil {
			t.Fatalf("OnStop failed: %v", err)
		}

		err = plugin.OnShutdown()
		if err != nil {
			t.Fatalf("OnShutdown failed: %v", err)
		}

		// Check calls were made in order
		expectedCalls := []string{"register", "start", "stop", "shutdown"}
		if len(calls) != len(expectedCalls) {
			t.Errorf("Expected %d calls, got %d", len(expectedCalls), len(calls))
		}

		for i, expected := range expectedCalls {
			if i >= len(calls) || calls[i] != expected {
				t.Errorf("Expected call %d to be '%s', got '%s'", i, expected, calls[i])
			}
		}

		t.Log("LifecycleHandlers test passed")
	})
}

// TestFakeHealthChecker tests the fake health checker functionality.
// zh: TestFakeHealthChecker 測試假健康檢查器功能。
func TestFakeHealthChecker(t *testing.T) {
	t.Run("NewFakeHealthChecker", func(t *testing.T) {
		checker := NewFakeHealthChecker("healthy", "All systems operational")
		if checker == nil {
			t.Fatal("NewFakeHealthChecker should not return nil")
		}

		ctx := context.Background()
		status := checker.CheckHealth(ctx)

		if status.Status != "healthy" {
			t.Errorf("Expected status 'healthy', got '%s'", status.Status)
		}

		if status.Message != "All systems operational" {
			t.Errorf("Expected message 'All systems operational', got '%s'", status.Message)
		}

		if status.Details == nil {
			t.Error("Details should not be nil")
		}

		t.Log("NewFakeHealthChecker test passed")
	})

	t.Run("SetHealthMetric", func(t *testing.T) {
		checker := NewFakeHealthChecker("healthy", "Test message")

		// Set metrics
		checker.SetHealthMetric("cpu_usage", 45.2)
		checker.SetHealthMetric("memory_usage", 78)
		checker.SetHealthMetric("disk_space", "available")

		// Get metrics
		metrics := checker.GetHealthMetrics()
		if metrics["cpu_usage"] != 45.2 {
			t.Errorf("Expected cpu_usage 45.2, got %v", metrics["cpu_usage"])
		}

		if metrics["memory_usage"] != 78 {
			t.Errorf("Expected memory_usage 78, got %v", metrics["memory_usage"])
		}

		if metrics["disk_space"] != "available" {
			t.Errorf("Expected disk_space 'available', got %v", metrics["disk_space"])
		}

		t.Log("SetHealthMetric test passed")
	})

	t.Run("SetHealthStatus", func(t *testing.T) {
		checker := NewFakeHealthChecker("healthy", "Initial message")

		// Change status
		checker.SetHealthStatus("degraded", "Performance issues detected")

		ctx := context.Background()
		status := checker.CheckHealth(ctx)

		if status.Status != "degraded" {
			t.Errorf("Expected status 'degraded', got '%s'", status.Status)
		}

		if status.Message != "Performance issues detected" {
			t.Errorf("Expected message 'Performance issues detected', got '%s'", status.Message)
		}

		t.Log("SetHealthStatus test passed")
	})
}
