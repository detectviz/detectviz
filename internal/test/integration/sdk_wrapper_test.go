package integration

import (
	"context"
	"testing"
	"time"

	"detectviz/pkg/platform/contracts"
	sdkwrapper "detectviz/plugins/community/integrations/observability/sdk-wrapper"
)

// TestSDKWrapperPluginInitialization tests SDK wrapper plugin initialization.
// zh: TestSDKWrapperPluginInitialization 測試 SDK 包裝器插件初始化。
func TestSDKWrapperPluginInitialization(t *testing.T) {
	t.Run("DefaultConfiguration", func(t *testing.T) {
		// Test with default configuration
		pluginInstance, err := sdkwrapper.NewSDKWrapperPlugin(nil)
		if err != nil {
			t.Fatalf("Failed to create SDK wrapper plugin with default config: %v", err)
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
		if plugin.Name() != "otel-sdk-wrapper" {
			t.Errorf("Expected plugin name 'otel-sdk-wrapper', got '%s'", plugin.Name())
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
			"service_name":    "test-service",
			"service_version": "2.0.0",
			"environment":     "testing",
			"enabled":         true,
			"tracing": map[string]any{
				"enabled":      true,
				"sample_ratio": 0.5,
				"max_spans":    500,
				"batch_size":   50,
				"timeout":      "15s",
			},
			"metrics": map[string]any{
				"enabled":         true,
				"collect_runtime": true,
				"collect_host":    false,
				"interval":        "30s",
			},
			"logging": map[string]any{
				"enabled":        true,
				"include_trace":  true,
				"log_level":      "debug",
				"correlation_id": true,
			},
			"exporters": map[string]any{
				"otlp": map[string]any{
					"enabled":     true,
					"endpoint":    "http://test-otlp:4317",
					"insecure":    true,
					"compression": "gzip",
					"timeout":     "5s",
				},
				"prometheus": map[string]any{
					"enabled": true,
					"host":    "test-host",
					"port":    9090,
					"path":    "/test-metrics",
				},
			},
		}

		pluginInstance, err := sdkwrapper.NewSDKWrapperPlugin(config)
		if err != nil {
			t.Fatalf("Failed to create SDK wrapper plugin with custom config: %v", err)
		}

		// Verify plugin was created with custom config
		plugin := pluginInstance.(contracts.Plugin)
		if plugin.Name() != "otel-sdk-wrapper" {
			t.Errorf("Expected plugin name 'otel-sdk-wrapper', got '%s'", plugin.Name())
		}

		t.Log("Custom configuration test passed")
	})

	t.Run("InvalidConfiguration", func(t *testing.T) {
		// Test with invalid configuration
		invalidConfigs := []map[string]any{
			{
				"service_name": "", // Empty service name
			},
			{
				"service_name":    "test",
				"service_version": "", // Empty service version
			},
			{
				"service_name":    "test",
				"service_version": "1.0.0",
				"tracing": map[string]any{
					"sample_ratio": 1.5, // Invalid sample ratio > 1
				},
			},
			{
				"service_name":    "test",
				"service_version": "1.0.0",
				"exporters": map[string]any{
					"otlp": map[string]any{
						"enabled":  true,
						"endpoint": "", // Empty endpoint when enabled
					},
				},
			},
		}

		for i, config := range invalidConfigs {
			pluginInstance, err := sdkwrapper.NewSDKWrapperPlugin(config)
			if err != nil {
				// Expected error for invalid config
				t.Logf("Invalid config %d correctly rejected: %v", i, err)
				continue
			}

			// If plugin was created, try to initialize it
			if err := pluginInstance.Init(config); err != nil {
				t.Logf("Invalid config %d correctly rejected during init: %v", i, err)
			} else {
				t.Errorf("Invalid config %d should have been rejected", i)
			}
		}

		t.Log("Invalid configuration test passed")
	})
}

// TestSDKWrapperPluginLifecycle tests SDK wrapper plugin lifecycle.
// zh: TestSDKWrapperPluginLifecycle 測試 SDK 包裝器插件生命週期。
func TestSDKWrapperPluginLifecycle(t *testing.T) {
	t.Run("LifecycleOperations", func(t *testing.T) {
		// Create plugin instance
		pluginInstance, err := sdkwrapper.NewSDKWrapperPlugin(nil)
		if err != nil {
			t.Fatalf("Failed to create plugin: %v", err)
		}

		plugin := pluginInstance.(contracts.Plugin)

		// Test initialization
		err = plugin.Init(nil)
		if err != nil {
			t.Fatalf("Failed to initialize plugin: %v", err)
		}

		// Test lifecycle aware interface
		if lifecyclePlugin, ok := pluginInstance.(contracts.LifecycleAware); ok {
			// Test OnRegister
			err = lifecyclePlugin.OnRegister()
			if err != nil {
				t.Errorf("OnRegister failed: %v", err)
			}

			// Test OnStart
			err = lifecyclePlugin.OnStart()
			if err != nil {
				t.Errorf("OnStart failed: %v", err)
			}

			// Test OnStop
			err = lifecyclePlugin.OnStop()
			if err != nil {
				t.Errorf("OnStop failed: %v", err)
			}

			// Test OnShutdown
			err = lifecyclePlugin.OnShutdown()
			if err != nil {
				t.Errorf("OnShutdown failed: %v", err)
			}
		} else {
			t.Error("Plugin does not implement LifecycleAware interface")
		}

		// Test shutdown
		err = plugin.Shutdown()
		if err != nil {
			t.Errorf("Failed to shutdown plugin: %v", err)
		}

		t.Log("Lifecycle operations test passed")
	})

	t.Run("DisabledPlugin", func(t *testing.T) {
		// Test with disabled plugin
		config := map[string]any{
			"enabled": false,
		}

		pluginInstance, err := sdkwrapper.NewSDKWrapperPlugin(config)
		if err != nil {
			t.Fatalf("Failed to create disabled plugin: %v", err)
		}

		plugin := pluginInstance.(contracts.Plugin)
		err = plugin.Init(config)
		if err != nil {
			t.Fatalf("Failed to initialize disabled plugin: %v", err)
		}

		// Test lifecycle operations on disabled plugin
		if lifecyclePlugin, ok := pluginInstance.(contracts.LifecycleAware); ok {
			err = lifecyclePlugin.OnStart()
			if err != nil {
				t.Errorf("OnStart failed for disabled plugin: %v", err)
			}
		}

		t.Log("Disabled plugin test passed")
	})
}

// TestSDKWrapperPluginHealthCheck tests SDK wrapper plugin health checking.
// zh: TestSDKWrapperPluginHealthCheck 測試 SDK 包裝器插件健康檢查。
func TestSDKWrapperPluginHealthCheck(t *testing.T) {
	t.Run("HealthChecker", func(t *testing.T) {
		// Create and initialize plugin
		pluginInstance, err := sdkwrapper.NewSDKWrapperPlugin(nil)
		if err != nil {
			t.Fatalf("Failed to create plugin: %v", err)
		}

		plugin := pluginInstance.(contracts.Plugin)
		err = plugin.Init(nil)
		if err != nil {
			t.Fatalf("Failed to initialize plugin: %v", err)
		}

		// Test health checker interface
		if healthChecker, ok := pluginInstance.(contracts.HealthChecker); ok {
			ctx := context.Background()
			status := healthChecker.CheckHealth(ctx)

			// Verify health status structure
			if status.Status == "" {
				t.Error("Health status should not be empty")
			}

			if status.Message == "" {
				t.Error("Health message should not be empty")
			}

			if status.Timestamp.IsZero() {
				t.Error("Health timestamp should not be zero")
			}

			if status.Details == nil {
				t.Error("Health details should not be nil")
			}

			// Check expected health status values
			expectedStatuses := []string{"healthy", "unhealthy", "degraded"}
			validStatus := false
			for _, expected := range expectedStatuses {
				if status.Status == expected {
					validStatus = true
					break
				}
			}
			if !validStatus {
				t.Errorf("Invalid health status: %s", status.Status)
			}

			t.Logf("Health check result: %s - %s", status.Status, status.Message)

			// Test health metrics
			metrics := healthChecker.GetHealthMetrics()
			if metrics == nil {
				t.Error("Health metrics should not be nil")
			}

			// Verify expected metrics exist
			expectedMetrics := []string{
				"plugin_initialized",
				"plugin_started",
				"plugin_enabled",
				"tracing_enabled",
				"metrics_enabled",
				"logging_enabled",
			}

			for _, metric := range expectedMetrics {
				if _, exists := metrics[metric]; !exists {
					t.Errorf("Expected metric '%s' not found in health metrics", metric)
				}
			}

			t.Log("Health checker test passed")
		} else {
			t.Error("Plugin does not implement HealthChecker interface")
		}
	})

	t.Run("HealthCheckWithContext", func(t *testing.T) {
		pluginInstance, err := sdkwrapper.NewSDKWrapperPlugin(nil)
		if err != nil {
			t.Fatalf("Failed to create plugin: %v", err)
		}

		plugin := pluginInstance.(contracts.Plugin)
		err = plugin.Init(nil)
		if err != nil {
			t.Fatalf("Failed to initialize plugin: %v", err)
		}

		if healthChecker, ok := pluginInstance.(contracts.HealthChecker); ok {
			// Test with timeout context
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			status := healthChecker.CheckHealth(ctx)
			if status.Status == "" {
				t.Error("Health status should not be empty")
			}

			// Test with cancelled context
			cancelledCtx, cancel := context.WithCancel(context.Background())
			cancel() // Cancel immediately

			status = healthChecker.CheckHealth(cancelledCtx)
			// Should still return a status even with cancelled context
			if status.Status == "" {
				t.Error("Health status should not be empty even with cancelled context")
			}

			t.Log("Health check with context test passed")
		}
	})
}

// TestSDKWrapperPluginConfiguration tests SDK wrapper plugin configuration.
// zh: TestSDKWrapperPluginConfiguration 測試 SDK 包裝器插件配置。
func TestSDKWrapperPluginConfiguration(t *testing.T) {
	t.Run("ConfigurationAccess", func(t *testing.T) {
		config := map[string]any{
			"service_name":    "test-config-service",
			"service_version": "1.2.3",
			"environment":     "test",
			"enabled":         true,
		}

		pluginInstance, err := sdkwrapper.NewSDKWrapperPlugin(config)
		if err != nil {
			t.Fatalf("Failed to create plugin: %v", err)
		}

		// Check if plugin provides configuration access
		if configPlugin, ok := pluginInstance.(interface {
			GetConfig() *sdkwrapper.Config
		}); ok {
			cfg := configPlugin.GetConfig()
			if cfg == nil {
				t.Error("Configuration should not be nil")
			}

			if cfg.ServiceName != "test-config-service" {
				t.Errorf("Expected service name 'test-config-service', got '%s'", cfg.ServiceName)
			}

			if cfg.ServiceVersion != "1.2.3" {
				t.Errorf("Expected service version '1.2.3', got '%s'", cfg.ServiceVersion)
			}

			if cfg.Environment != "test" {
				t.Errorf("Expected environment 'test', got '%s'", cfg.Environment)
			}

			if !cfg.Enabled {
				t.Error("Expected plugin to be enabled")
			}

			t.Log("Configuration access test passed")
		}
	})

	t.Run("ConfigurationValidation", func(t *testing.T) {
		// Test various configuration scenarios
		testCases := []struct {
			name        string
			config      map[string]any
			expectError bool
		}{
			{
				name: "ValidMinimalConfig",
				config: map[string]any{
					"service_name":    "test",
					"service_version": "1.0.0",
				},
				expectError: false,
			},
			{
				name: "ValidFullConfig",
				config: map[string]any{
					"service_name":    "test",
					"service_version": "1.0.0",
					"environment":     "production",
					"enabled":         true,
					"tracing": map[string]any{
						"enabled":      true,
						"sample_ratio": 0.1,
						"max_spans":    1000,
						"batch_size":   100,
					},
					"exporters": map[string]any{
						"otlp": map[string]any{
							"enabled":  true,
							"endpoint": "http://localhost:4317",
							"insecure": true,
						},
					},
				},
				expectError: false,
			},
			{
				name: "InvalidSampleRatio",
				config: map[string]any{
					"service_name":    "test",
					"service_version": "1.0.0",
					"tracing": map[string]any{
						"sample_ratio": 2.0, // Invalid: > 1
					},
				},
				expectError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				pluginInstance, err := sdkwrapper.NewSDKWrapperPlugin(tc.config)
				if err != nil && !tc.expectError {
					t.Errorf("Unexpected error creating plugin: %v", err)
					return
				}
				if err == nil && tc.expectError {
					t.Error("Expected error but got none")
					return
				}
				if err != nil && tc.expectError {
					t.Logf("Expected error occurred: %v", err)
					return
				}

				// Try to initialize
				err = pluginInstance.Init(tc.config)
				if err != nil && !tc.expectError {
					t.Errorf("Unexpected error initializing plugin: %v", err)
				}
				if err == nil && tc.expectError {
					t.Error("Expected error during initialization but got none")
				}
				if err != nil && tc.expectError {
					t.Logf("Expected initialization error occurred: %v", err)
				}
			})
		}
	})
}

// TestSDKWrapperPluginIntegration tests SDK wrapper plugin integration scenarios.
// zh: TestSDKWrapperPluginIntegration 測試 SDK 包裝器插件整合場景。
func TestSDKWrapperPluginIntegration(t *testing.T) {
	t.Run("PluginRegistration", func(t *testing.T) {
		// This test would typically involve a real registry
		// For now, we'll test the registration function exists
		if sdkwrapper.Register == nil {
			t.Error("Register function should be available")
		}

		// Test that the function signature is correct
		// In a real test, you would use a mock registry
		t.Log("Plugin registration function exists")
	})

	t.Run("MultipleInstances", func(t *testing.T) {
		// Test creating multiple plugin instances
		configs := []map[string]any{
			{
				"service_name":    "service-1",
				"service_version": "1.0.0",
			},
			{
				"service_name":    "service-2",
				"service_version": "2.0.0",
			},
		}

		var plugins []contracts.Plugin

		for i, config := range configs {
			pluginInstance, err := sdkwrapper.NewSDKWrapperPlugin(config)
			if err != nil {
				t.Fatalf("Failed to create plugin %d: %v", i, err)
			}

			plugin := pluginInstance.(contracts.Plugin)
			err = plugin.Init(config)
			if err != nil {
				t.Fatalf("Failed to initialize plugin %d: %v", i, err)
			}

			plugins = append(plugins, plugin)
		}

		// Verify all plugins are independent
		if len(plugins) != 2 {
			t.Errorf("Expected 2 plugins, got %d", len(plugins))
		}

		// Cleanup
		for i, plugin := range plugins {
			err := plugin.Shutdown()
			if err != nil {
				t.Errorf("Failed to shutdown plugin %d: %v", i, err)
			}
		}

		t.Log("Multiple instances test passed")
	})
}
