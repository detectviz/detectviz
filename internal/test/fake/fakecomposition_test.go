package fake

import (
	"context"
	"errors"
	"testing"

	"detectviz/pkg/shared/log"
)

// TestFakeComposition tests the fake composition functionality.
// zh: TestFakeComposition 測試假組合功能。
func TestFakeComposition(t *testing.T) {
	t.Run("NewFakeComposition", func(t *testing.T) {
		composition := NewFakeComposition("test-composition", "1.0.0")
		if composition == nil {
			t.Fatal("NewFakeComposition should not return nil")
		}

		if composition.name != "test-composition" {
			t.Errorf("Expected name 'test-composition', got '%s'", composition.name)
		}

		if composition.version != "1.0.0" {
			t.Errorf("Expected version '1.0.0', got '%s'", composition.version)
		}

		// Check initial state
		plugins := composition.ListPlugins()
		if len(plugins) != 0 {
			t.Errorf("Expected 0 plugins initially, got %d", len(plugins))
		}

		info := composition.GetCompositionInfo()
		if info["total_plugins"] != 0 {
			t.Errorf("Expected 0 total plugins, got %v", info["total_plugins"])
		}

		t.Log("NewFakeComposition test passed")
	})

	t.Run("AddPlugin", func(t *testing.T) {
		composition := NewFakeComposition("test-composition", "1.0.0")

		// Add a plugin
		config := map[string]any{
			"setting1": "value1",
			"setting2": 42,
		}
		dependencies := []string{}

		err := composition.AddPlugin("test-plugin", "importer", config, dependencies)
		if err != nil {
			t.Fatalf("Failed to add plugin: %v", err)
		}

		// Check plugin was added
		plugins := composition.ListPlugins()
		if len(plugins) != 1 {
			t.Errorf("Expected 1 plugin, got %d", len(plugins))
		}

		if plugins[0] != "test-plugin" {
			t.Errorf("Expected plugin name 'test-plugin', got '%s'", plugins[0])
		}

		// Get plugin details
		plugin, err := composition.GetPlugin("test-plugin")
		if err != nil {
			t.Fatalf("Failed to get plugin: %v", err)
		}

		if plugin.Name != "test-plugin" {
			t.Errorf("Expected plugin name 'test-plugin', got '%s'", plugin.Name)
		}

		if plugin.Type != "importer" {
			t.Errorf("Expected plugin type 'importer', got '%s'", plugin.Type)
		}

		if !plugin.Enabled {
			t.Error("Plugin should be enabled by default")
		}

		// Try to add same plugin again (should fail)
		err = composition.AddPlugin("test-plugin", "exporter", config, dependencies)
		if err == nil {
			t.Error("Expected error when adding duplicate plugin")
		}

		t.Log("AddPlugin test passed")
	})

	t.Run("GetPlugin", func(t *testing.T) {
		composition := NewFakeComposition("test-composition", "1.0.0")

		// Add a plugin
		config := map[string]any{"test": "value"}
		err := composition.AddPlugin("test-plugin", "integration", config, []string{})
		if err != nil {
			t.Fatalf("Failed to add plugin: %v", err)
		}

		// Get existing plugin
		plugin, err := composition.GetPlugin("test-plugin")
		if err != nil {
			t.Fatalf("Failed to get plugin: %v", err)
		}

		if plugin.Config["test"] != "value" {
			t.Errorf("Expected config test 'value', got '%v'", plugin.Config["test"])
		}

		// Try to get non-existent plugin
		_, err = composition.GetPlugin("non-existent")
		if err == nil {
			t.Error("Expected error when getting non-existent plugin")
		}

		t.Log("GetPlugin test passed")
	})

	t.Run("InitializePlugins", func(t *testing.T) {
		// Initialize logger for testing
		logger, _ := log.NewLogger(&log.LoggerConfig{
			Type:   "console",
			Level:  "info",
			Format: "text",
		})
		log.SetGlobalLogger(logger)

		composition := NewFakeComposition("test-composition", "1.0.0")

		// Add multiple plugins
		plugins := []struct {
			name    string
			ptype   string
			enabled bool
		}{
			{"plugin-1", "importer", true},
			{"plugin-2", "exporter", true},
			{"plugin-3", "integration", false}, // disabled
		}

		for _, p := range plugins {
			config := map[string]any{"enabled": p.enabled}
			err := composition.AddPlugin(p.name, p.ptype, config, []string{})
			if err != nil {
				t.Fatalf("Failed to add plugin %s: %v", p.name, err)
			}

			// Set enabled state
			plugin, _ := composition.GetPlugin(p.name)
			plugin.Enabled = p.enabled
		}

		// Initialize plugins
		ctx := context.Background()
		err := composition.InitializePlugins(ctx)
		if err != nil {
			t.Fatalf("Failed to initialize plugins: %v", err)
		}

		// Check initialization status
		info := composition.GetCompositionInfo()
		if info["initialized_plugins"] != 2 { // Only enabled plugins
			t.Errorf("Expected 2 initialized plugins, got %v", info["initialized_plugins"])
		}

		// Check individual plugin status
		plugin1, _ := composition.GetPlugin("plugin-1")
		if !plugin1.initialized {
			t.Error("Plugin-1 should be initialized")
		}

		plugin3, _ := composition.GetPlugin("plugin-3")
		if plugin3.initialized {
			t.Error("Plugin-3 should not be initialized (disabled)")
		}

		t.Log("InitializePlugins test passed")
	})

	t.Run("StartPlugins", func(t *testing.T) {
		// Initialize logger for testing
		logger, _ := log.NewLogger(&log.LoggerConfig{
			Type:   "console",
			Level:  "info",
			Format: "text",
		})
		log.SetGlobalLogger(logger)

		composition := NewFakeComposition("test-composition", "1.0.0")

		// Add and initialize plugins
		config := map[string]any{"test": "value"}
		err := composition.AddPlugin("test-plugin", "importer", config, []string{})
		if err != nil {
			t.Fatalf("Failed to add plugin: %v", err)
		}

		ctx := context.Background()
		err = composition.InitializePlugins(ctx)
		if err != nil {
			t.Fatalf("Failed to initialize plugins: %v", err)
		}

		// Start plugins
		err = composition.StartPlugins(ctx)
		if err != nil {
			t.Fatalf("Failed to start plugins: %v", err)
		}

		// Check start status
		info := composition.GetCompositionInfo()
		if info["started_plugins"] != 1 {
			t.Errorf("Expected 1 started plugin, got %v", info["started_plugins"])
		}

		if !info["composition_started"].(bool) {
			t.Error("Composition should be marked as started")
		}

		// Check individual plugin status
		plugin, _ := composition.GetPlugin("test-plugin")
		if !plugin.started {
			t.Error("Plugin should be started")
		}

		t.Log("StartPlugins test passed")
	})

	t.Run("SimulateError", func(t *testing.T) {
		// Initialize logger for testing
		logger, _ := log.NewLogger(&log.LoggerConfig{
			Type:   "console",
			Level:  "info",
			Format: "text",
		})
		log.SetGlobalLogger(logger)

		composition := NewFakeComposition("test-composition", "1.0.0")

		// Add plugin
		config := map[string]any{"test": "value"}
		err := composition.AddPlugin("test-plugin", "importer", config, []string{})
		if err != nil {
			t.Fatalf("Failed to add plugin: %v", err)
		}

		// Simulate error
		testError := errors.New("simulated initialization error")
		composition.SimulateError("test-plugin", testError)

		// Try to initialize plugins (should fail)
		ctx := context.Background()
		err = composition.InitializePlugins(ctx)
		if err == nil {
			t.Error("Expected error during initialization")
		}

		if !errors.Is(err, testError) {
			t.Errorf("Expected simulated error, got: %v", err)
		}

		t.Log("SimulateError test passed")
	})

	t.Run("GetCompositionInfo", func(t *testing.T) {
		// Initialize logger for testing
		logger, _ := log.NewLogger(&log.LoggerConfig{
			Type:   "console",
			Level:  "info",
			Format: "text",
		})
		log.SetGlobalLogger(logger)

		composition := NewFakeComposition("test-composition", "1.0.0")

		// Add multiple plugins with different states
		plugins := []struct {
			name    string
			enabled bool
		}{
			{"plugin-1", true},
			{"plugin-2", true},
			{"plugin-3", false},
		}

		for _, p := range plugins {
			config := map[string]any{"enabled": p.enabled}
			err := composition.AddPlugin(p.name, "importer", config, []string{})
			if err != nil {
				t.Fatalf("Failed to add plugin %s: %v", p.name, err)
			}

			// Set enabled state
			plugin, _ := composition.GetPlugin(p.name)
			plugin.Enabled = p.enabled
		}

		// Initialize and start some plugins
		ctx := context.Background()
		err := composition.InitializePlugins(ctx)
		if err != nil {
			t.Fatalf("Failed to initialize plugins: %v", err)
		}

		err = composition.StartPlugins(ctx)
		if err != nil {
			t.Fatalf("Failed to start plugins: %v", err)
		}

		// Check composition info
		info := composition.GetCompositionInfo()

		if info["name"] != "test-composition" {
			t.Errorf("Expected name 'test-composition', got %v", info["name"])
		}

		if info["version"] != "1.0.0" {
			t.Errorf("Expected version '1.0.0', got %v", info["version"])
		}

		if info["total_plugins"] != 3 {
			t.Errorf("Expected 3 total plugins, got %v", info["total_plugins"])
		}

		if info["enabled_plugins"] != 2 {
			t.Errorf("Expected 2 enabled plugins, got %v", info["enabled_plugins"])
		}

		if info["initialized_plugins"] != 2 {
			t.Errorf("Expected 2 initialized plugins, got %v", info["initialized_plugins"])
		}

		if info["started_plugins"] != 2 {
			t.Errorf("Expected 2 started plugins, got %v", info["started_plugins"])
		}

		if !info["composition_started"].(bool) {
			t.Error("Composition should be started")
		}

		t.Log("GetCompositionInfo test passed")
	})

	t.Run("Reset", func(t *testing.T) {
		composition := NewFakeComposition("test-composition", "1.0.0")

		// Add plugins
		config := map[string]any{"test": "value"}
		err := composition.AddPlugin("plugin-1", "importer", config, []string{})
		if err != nil {
			t.Fatalf("Failed to add plugin: %v", err)
		}

		err = composition.AddPlugin("plugin-2", "exporter", config, []string{})
		if err != nil {
			t.Fatalf("Failed to add plugin: %v", err)
		}

		// Check plugins exist
		plugins := composition.ListPlugins()
		if len(plugins) != 2 {
			t.Errorf("Expected 2 plugins before reset, got %d", len(plugins))
		}

		// Reset composition
		composition.Reset()

		// Check composition is empty
		plugins = composition.ListPlugins()
		if len(plugins) != 0 {
			t.Errorf("Expected 0 plugins after reset, got %d", len(plugins))
		}

		info := composition.GetCompositionInfo()
		if info["total_plugins"] != 0 {
			t.Errorf("Expected 0 total plugins after reset, got %v", info["total_plugins"])
		}

		if info["composition_started"].(bool) {
			t.Error("Composition should not be started after reset")
		}

		t.Log("Reset test passed")
	})
}

// TestFakeCompositionPlugin tests the fake composition plugin functionality.
// zh: TestFakeCompositionPlugin 測試假組合插件功能。
func TestFakeCompositionPlugin(t *testing.T) {
	t.Run("PluginCreation", func(t *testing.T) {
		composition := NewFakeComposition("test-composition", "1.0.0")

		config := map[string]any{
			"host":    "localhost",
			"port":    8080,
			"enabled": true,
			"timeout": 30,
		}

		err := composition.AddPlugin("web-server", "integration", config, []string{"database", "cache"})
		if err != nil {
			t.Fatalf("Failed to add plugin: %v", err)
		}

		plugin, err := composition.GetPlugin("web-server")
		if err != nil {
			t.Fatalf("Failed to get plugin: %v", err)
		}

		// Check plugin properties
		if plugin.Name != "web-server" {
			t.Errorf("Expected name 'web-server', got '%s'", plugin.Name)
		}

		if plugin.Type != "integration" {
			t.Errorf("Expected type 'integration', got '%s'", plugin.Type)
		}

		if len(plugin.Dependencies) != 2 {
			t.Errorf("Expected 2 dependencies, got %d", len(plugin.Dependencies))
		}

		if plugin.Dependencies[0] != "database" || plugin.Dependencies[1] != "cache" {
			t.Errorf("Expected dependencies [database, cache], got %v", plugin.Dependencies)
		}

		if plugin.Config["host"] != "localhost" {
			t.Errorf("Expected host 'localhost', got '%v'", plugin.Config["host"])
		}

		if plugin.Config["port"] != 8080 {
			t.Errorf("Expected port 8080, got %v", plugin.Config["port"])
		}

		// Check metadata
		if plugin.Metadata == nil {
			t.Fatal("Plugin metadata should not be nil")
		}

		if plugin.Metadata.Name != "web-server" {
			t.Errorf("Expected metadata name 'web-server', got '%s'", plugin.Metadata.Name)
		}

		if plugin.Metadata.Type != "integration" {
			t.Errorf("Expected metadata type 'integration', got '%s'", plugin.Metadata.Type)
		}

		if plugin.Metadata.Category != "test" {
			t.Errorf("Expected metadata category 'test', got '%s'", plugin.Metadata.Category)
		}

		// Check instance
		if plugin.Instance == nil {
			t.Fatal("Plugin instance should not be nil")
		}

		if plugin.Instance.Name() != "web-server" {
			t.Errorf("Expected instance name 'web-server', got '%s'", plugin.Instance.Name())
		}

		// Check initial states
		if plugin.initialized {
			t.Error("Plugin should not be initialized initially")
		}

		if plugin.started {
			t.Error("Plugin should not be started initially")
		}

		t.Log("PluginCreation test passed")
	})

	t.Run("PluginStates", func(t *testing.T) {
		composition := NewFakeComposition("test-composition", "1.0.0")

		config := map[string]any{"test": "value"}
		err := composition.AddPlugin("state-test", "importer", config, []string{})
		if err != nil {
			t.Fatalf("Failed to add plugin: %v", err)
		}

		plugin, _ := composition.GetPlugin("state-test")

		// Test enabled/disabled state
		if !plugin.Enabled {
			t.Error("Plugin should be enabled by default")
		}

		plugin.Enabled = false
		if plugin.Enabled {
			t.Error("Plugin should be disabled after setting")
		}

		// Test priority
		if plugin.Priority != 0 {
			t.Errorf("Expected default priority 0, got %d", plugin.Priority)
		}

		plugin.Priority = 10
		if plugin.Priority != 10 {
			t.Errorf("Expected priority 10, got %d", plugin.Priority)
		}

		t.Log("PluginStates test passed")
	})
}
