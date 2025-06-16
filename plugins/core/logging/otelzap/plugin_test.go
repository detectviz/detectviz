package otelzap

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"detectviz/pkg/platform/contracts"
)

// TestOtelZapPluginCreation tests plugin creation with various configurations.
// zh: TestOtelZapPluginCreation 測試各種配置下的插件建立。
func TestOtelZapPluginCreation(t *testing.T) {
	t.Run("DefaultConfiguration", func(t *testing.T) {
		// Test with default configuration
		pluginInstance, err := NewOtelZapPlugin(nil)
		if err != nil {
			t.Fatalf("Failed to create OtelZap plugin with default config: %v", err)
		}

		if pluginInstance == nil {
			t.Fatal("Plugin instance should not be nil")
		}

		// Verify plugin implements required interfaces
		plugin, ok := pluginInstance.(contracts.Plugin)
		if !ok {
			t.Fatal("Plugin does not implement Plugin interface")
		}

		// Check basic plugin information
		if plugin.Name() != "otelzap-logger" {
			t.Errorf("Expected plugin name 'otelzap-logger', got '%s'", plugin.Name())
		}

		if plugin.Version() != "1.0.0" {
			t.Errorf("Expected plugin version '1.0.0', got '%s'", plugin.Version())
		}

		if plugin.Description() == "" {
			t.Error("Plugin description should not be empty")
		}

		t.Log("Default configuration test passed")
	})

	t.Run("CustomConfiguration", func(t *testing.T) {
		// Test with custom configuration
		config := map[string]any{
			"enabled":         true,
			"level":           "debug",
			"format":          "json",
			"output_type":     "console",
			"output":          "stdout",
			"service_name":    "test-service",
			"service_version": "2.0.0",
			"environment":     "testing",
			"otel": map[string]any{
				"enabled":        true,
				"include_trace":  true,
				"trace_id_field": "trace_id",
				"span_id_field":  "span_id",
			},
		}

		pluginInstance, err := NewOtelZapPlugin(config)
		if err != nil {
			t.Fatalf("Failed to create OtelZap plugin with custom config: %v", err)
		}

		// Verify plugin was created with custom config
		plugin := pluginInstance.(contracts.Plugin)
		if plugin.Name() != "otelzap-logger" {
			t.Errorf("Expected plugin name 'otelzap-logger', got '%s'", plugin.Name())
		}

		// Check if we can access the plugin configuration
		otelZapPlugin := pluginInstance.(*OtelZapPlugin)
		if otelZapPlugin.config.ServiceName != "test-service" {
			t.Errorf("Expected service name 'test-service', got '%s'", otelZapPlugin.config.ServiceName)
		}

		if otelZapPlugin.config.Level != "debug" {
			t.Errorf("Expected level 'debug', got '%s'", otelZapPlugin.config.Level)
		}

		t.Log("Custom configuration test passed")
	})

	t.Run("InvalidConfiguration", func(t *testing.T) {
		// Test with invalid configuration
		config := map[string]any{
			"level": "invalid_level",
		}

		_, err := NewOtelZapPlugin(config)
		if err == nil {
			t.Fatal("Expected error for invalid configuration")
		}

		t.Logf("Got expected error: %v", err)
		t.Log("Expected error for invalid configuration")
	})
}

// TestOtelZapPluginInitialization tests plugin initialization.
// zh: TestOtelZapPluginInitialization 測試插件初始化。
func TestOtelZapPluginInitialization(t *testing.T) {
	t.Run("InitWithDefaultConfig", func(t *testing.T) {
		pluginInstance, err := NewOtelZapPlugin(nil)
		if err != nil {
			t.Fatalf("Failed to create plugin: %v", err)
		}

		plugin := pluginInstance.(*OtelZapPlugin)

		// Initialize the plugin
		err = plugin.Init(nil)
		if err != nil {
			t.Fatalf("Failed to initialize plugin: %v", err)
		}

		if !plugin.initialized {
			t.Error("Plugin should be initialized")
		}

		if plugin.logger == nil {
			t.Error("Logger should be created after initialization")
		}

		t.Log("Plugin initialization test passed")
	})

	t.Run("InitWithCustomConfig", func(t *testing.T) {
		// Create temporary directory for log files
		tempDir := t.TempDir()
		logFile := filepath.Join(tempDir, "test.log")

		config := map[string]any{
			"enabled":     true,
			"level":       "info",
			"format":      "json",
			"output_type": "file",
			"file_config": map[string]any{
				"filename":    logFile,
				"max_size":    1,
				"max_backups": 2,
				"max_age":     1,
				"compress":    false,
			},
		}

		pluginInstance, err := NewOtelZapPlugin(nil)
		if err != nil {
			t.Fatalf("Failed to create plugin: %v", err)
		}

		plugin := pluginInstance.(*OtelZapPlugin)

		// Initialize with custom config
		err = plugin.Init(config)
		if err != nil {
			t.Fatalf("Failed to initialize plugin with custom config: %v", err)
		}

		if !plugin.initialized {
			t.Error("Plugin should be initialized")
		}

		// Test logging to file
		logger := plugin.Logger()
		logger.Info("Test message to file")

		// Give some time for file write
		time.Sleep(100 * time.Millisecond)

		// Check if log file was created
		if _, err := os.Stat(logFile); os.IsNotExist(err) {
			t.Fatalf("Log file %s was not created", logFile)
		}

		content, err := os.ReadFile(logFile)
		if err != nil {
			t.Fatalf("Failed to read log file: %v", err)
		}

		if !strings.Contains(string(content), "Test message to file") {
			t.Errorf("Log file does not contain expected message")
		}

		t.Log("Custom config initialization test passed")
	})
}

// TestOtelZapPluginLifecycle tests plugin lifecycle management.
// zh: TestOtelZapPluginLifecycle 測試插件生命週期管理。
func TestOtelZapPluginLifecycle(t *testing.T) {
	t.Run("FullLifecycle", func(t *testing.T) {
		pluginInstance, err := NewOtelZapPlugin(nil)
		if err != nil {
			t.Fatalf("Failed to create plugin: %v", err)
		}

		plugin := pluginInstance.(*OtelZapPlugin)

		// Test OnRegister
		err = plugin.OnRegister()
		if err != nil {
			t.Errorf("OnRegister failed: %v", err)
		}

		// Test Init
		err = plugin.Init(nil)
		if err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		// Test OnStart
		err = plugin.OnStart()
		if err != nil {
			t.Errorf("OnStart failed: %v", err)
		}

		if !plugin.started {
			t.Error("Plugin should be started")
		}

		// Test OnStop
		err = plugin.OnStop()
		if err != nil {
			t.Errorf("OnStop failed: %v", err)
		}

		if plugin.started {
			t.Error("Plugin should be stopped")
		}

		// Test OnShutdown
		err = plugin.OnShutdown()
		if err != nil {
			t.Errorf("OnShutdown failed: %v", err)
		}

		// Test Shutdown
		err = plugin.Shutdown()
		if err != nil {
			t.Errorf("Shutdown failed: %v", err)
		}

		if plugin.initialized {
			t.Error("Plugin should be uninitialized after shutdown")
		}

		t.Log("Full lifecycle test passed")
	})
}

// TestOtelZapPluginLoggerProvider tests LoggerProvider interface.
// zh: TestOtelZapPluginLoggerProvider 測試 LoggerProvider 介面。
func TestOtelZapPluginLoggerProvider(t *testing.T) {
	t.Run("LoggerProviderInterface", func(t *testing.T) {
		pluginInstance, err := NewOtelZapPlugin(nil)
		if err != nil {
			t.Fatalf("Failed to create plugin: %v", err)
		}

		plugin := pluginInstance.(*OtelZapPlugin)

		// Initialize the plugin
		err = plugin.Init(nil)
		if err != nil {
			t.Fatalf("Failed to initialize plugin: %v", err)
		}

		// Test Logger method
		logger := plugin.Logger()
		if logger == nil {
			t.Error("Logger should not be nil")
		}

		// Test WithContext method
		ctx := context.Background()
		ctxLogger := plugin.WithContext(ctx)
		if ctxLogger == nil {
			t.Error("Context logger should not be nil")
		}

		// Test logging with different levels
		logger.Debug("Debug message")
		logger.Info("Info message")
		logger.Warn("Warning message")
		logger.Error("Error message")

		// Test SetLevel
		err = plugin.SetLevel("debug")
		if err != nil {
			t.Errorf("SetLevel failed: %v", err)
		}

		// Test invalid level
		err = plugin.SetLevel("invalid")
		if err == nil {
			t.Error("Expected error for invalid level")
		}

		// Test Flush
		err = plugin.Flush()
		if err != nil {
			t.Errorf("Flush failed: %v", err)
		}

		// Test Close
		err = plugin.Close()
		if err != nil {
			t.Errorf("Close failed: %v", err)
		}

		t.Log("LoggerProvider interface test passed")
	})
}

// TestOtelZapPluginHealthChecker tests HealthChecker interface.
// zh: TestOtelZapPluginHealthChecker 測試 HealthChecker 介面。
func TestOtelZapPluginHealthChecker(t *testing.T) {
	t.Run("HealthCheck", func(t *testing.T) {
		pluginInstance, err := NewOtelZapPlugin(nil)
		if err != nil {
			t.Fatalf("Failed to create plugin: %v", err)
		}

		plugin := pluginInstance.(*OtelZapPlugin)
		ctx := context.Background()

		// Test health check before initialization
		status := plugin.CheckHealth(ctx)
		if status.Status != "unhealthy" {
			t.Errorf("Expected status 'unhealthy' before init, got '%s'", status.Status)
		}

		// Initialize and test health check
		err = plugin.Init(nil)
		if err != nil {
			t.Fatalf("Failed to initialize plugin: %v", err)
		}

		status = plugin.CheckHealth(ctx)
		if status.Status != "degraded" {
			t.Errorf("Expected status 'degraded' after init but before start, got '%s'", status.Status)
		}

		// Start and test health check
		err = plugin.OnStart()
		if err != nil {
			t.Fatalf("Failed to start plugin: %v", err)
		}

		status = plugin.CheckHealth(ctx)
		if status.Status != "healthy" {
			t.Errorf("Expected status 'healthy' after start, got '%s'", status.Status)
		}

		// Check status details
		if status.Details == nil {
			t.Error("Status details should not be nil")
		}

		if enabled, exists := status.Details["enabled"]; !exists || enabled != true {
			t.Error("Status details should include enabled=true")
		}

		// Test health metrics
		metrics := plugin.GetHealthMetrics()
		if metrics == nil {
			t.Error("Health metrics should not be nil")
		}

		if initialized, exists := metrics["plugin_initialized"]; !exists || initialized != true {
			t.Error("Health metrics should include plugin_initialized=true")
		}

		t.Log("Health check test passed")
	})
}

// TestOtelZapPluginConfiguration tests configuration parsing and validation.
// zh: TestOtelZapPluginConfiguration 測試配置解析和驗證。
func TestOtelZapPluginConfiguration(t *testing.T) {
	t.Run("ConfigValidation", func(t *testing.T) {
		testCases := []struct {
			name        string
			config      map[string]any
			expectError bool
		}{
			{
				name: "ValidConfig",
				config: map[string]any{
					"enabled":     true,
					"level":       "info",
					"format":      "json",
					"output_type": "console",
				},
				expectError: false,
			},
			{
				name: "InvalidLevel",
				config: map[string]any{
					"level": "invalid_level",
				},
				expectError: true,
			},
			{
				name: "InvalidFormat",
				config: map[string]any{
					"format": "invalid_format",
				},
				expectError: true,
			},
			{
				name: "InvalidOutputType",
				config: map[string]any{
					"output_type": "invalid_output",
				},
				expectError: true,
			},
			{
				name: "FileOutputWithoutFileConfig",
				config: map[string]any{
					"output_type": "file",
				},
				expectError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				pluginInstance, err := NewOtelZapPlugin(tc.config)
				if tc.expectError {
					if err == nil {
						t.Errorf("Expected error for %s, but got none", tc.name)
					}
				} else {
					if err != nil {
						t.Errorf("Unexpected error for %s: %v", tc.name, err)
					}
					if pluginInstance == nil {
						t.Errorf("Plugin instance should not be nil for %s", tc.name)
					}
				}
			})
		}

		t.Log("Configuration validation test passed")
	})
}
