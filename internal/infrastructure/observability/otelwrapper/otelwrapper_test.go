package otelwrapper

import (
	"context"
	"detectviz/pkg/platform/contracts"
	"testing"
	"time"
)

// TestOtelWrapperCreation tests the creation of OtelWrapper.
// zh: TestOtelWrapperCreation 測試 OtelWrapper 的建立。
func TestOtelWrapperCreation(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		wrapper := NewOtelWrapper(nil)
		if wrapper == nil {
			t.Fatal("Wrapper should not be nil")
		}

		if wrapper.config == nil {
			t.Fatal("Config should not be nil")
		}

		if wrapper.config.ServiceName != "detectviz" {
			t.Errorf("Expected service name 'detectviz', got '%s'", wrapper.config.ServiceName)
		}

		t.Log("Default config creation test passed")
	})

	t.Run("CustomConfig", func(t *testing.T) {
		config := &Config{
			ServiceName:    "test-service",
			ServiceVersion: "1.0.0",
			Environment:    "test",
			Enabled:        true,
		}

		wrapper := NewOtelWrapper(config)
		if wrapper == nil {
			t.Fatal("Wrapper should not be nil")
		}

		if wrapper.config.ServiceName != "test-service" {
			t.Errorf("Expected service name 'test-service', got '%s'", wrapper.config.ServiceName)
		}

		t.Log("Custom config creation test passed")
	})
}

// TestOtelWrapperInitialization tests the initialization of OtelWrapper.
// zh: TestOtelWrapperInitialization 測試 OtelWrapper 的初始化。
func TestOtelWrapperInitialization(t *testing.T) {
	t.Run("EnabledWrapper", func(t *testing.T) {
		config := &Config{
			ServiceName:    "test-service",
			ServiceVersion: "1.0.0",
			Environment:    "test",
			Enabled:        true,
			Tracing: TracingConfig{
				Enabled: true,
			},
			Metrics: MetricsConfig{
				Enabled: true,
			},
			Logging: LoggingConfig{
				Enabled: true,
			},
			Resource: ResourceConfig{
				DetectHost:    false,
				DetectProcess: false,
				DetectRuntime: false,
				CustomAttrs:   make(map[string]string),
			},
		}

		wrapper := NewOtelWrapper(config)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := wrapper.Initialize(ctx)
		if err != nil {
			t.Fatalf("Failed to initialize wrapper: %v", err)
		}

		if !wrapper.initialized {
			t.Error("Wrapper should be initialized")
		}

		// Test provider access
		tracer := wrapper.Tracer("test")
		if tracer == nil {
			t.Error("Tracer should not be nil")
		}

		meter := wrapper.Meter("test")
		if meter == nil {
			t.Error("Meter should not be nil")
		}

		// Test shutdown
		err = wrapper.Shutdown(ctx)
		if err != nil {
			t.Errorf("Failed to shutdown wrapper: %v", err)
		}

		t.Log("Enabled wrapper initialization test passed")
	})

	t.Run("DisabledWrapper", func(t *testing.T) {
		config := &Config{
			ServiceName: "test-service",
			Enabled:     false,
		}

		wrapper := NewOtelWrapper(config)

		ctx := context.Background()

		err := wrapper.Initialize(ctx)
		if err != nil {
			t.Fatalf("Failed to initialize disabled wrapper: %v", err)
		}

		// Disabled wrapper should not be marked as initialized
		// since it skips actual initialization

		t.Log("Disabled wrapper initialization test passed")
	})
}

// TestOtelWrapperGlobalAccess tests global wrapper functionality.
// zh: TestOtelWrapperGlobalAccess 測試全域包裝器功能。
func TestOtelWrapperGlobalAccess(t *testing.T) {
	t.Run("GlobalWrapper", func(t *testing.T) {
		// Test global wrapper access
		wrapper := GetGlobalWrapper()
		if wrapper == nil {
			t.Fatal("Global wrapper should not be nil")
		}

		// Test convenience functions
		tracer := Tracer("test")
		if tracer == nil {
			t.Error("Global tracer should not be nil")
		}

		meter := Meter("test")
		if meter == nil {
			t.Error("Global meter should not be nil")
		}

		t.Log("Global wrapper access test passed")
	})

	t.Run("SetGlobalWrapper", func(t *testing.T) {
		originalWrapper := GetGlobalWrapper()

		// Create a new wrapper
		newWrapper := NewOtelWrapper(&Config{
			ServiceName: "custom-global",
			Enabled:     true,
		})

		// Set as global
		SetGlobalWrapper(newWrapper)

		// Verify the change
		currentWrapper := GetGlobalWrapper()
		if currentWrapper != newWrapper {
			t.Error("Global wrapper should be updated")
		}

		if currentWrapper.config.ServiceName != "custom-global" {
			t.Errorf("Expected service name 'custom-global', got '%s'",
				currentWrapper.config.ServiceName)
		}

		// Restore original
		SetGlobalWrapper(originalWrapper)

		t.Log("Set global wrapper test passed")
	})
}

// TestOtelWrapperLoggerProviderInjection tests logger provider injection.
// zh: TestOtelWrapperLoggerProviderInjection 測試日誌提供器注入。
func TestOtelWrapperLoggerProviderInjection(t *testing.T) {
	t.Run("InjectLoggerProvider", func(t *testing.T) {
		wrapper := NewOtelWrapper(nil)

		// Create a mock logger provider
		mockProvider := &mockLoggerProvider{}

		// Inject the provider
		wrapper.InjectLoggerProvider(mockProvider)

		// Verify injection
		injectedProvider := wrapper.Logger()
		if injectedProvider != mockProvider {
			t.Error("Logger provider should be injected correctly")
		}

		t.Log("Logger provider injection test passed")
	})
}

// TestDefaultConfig tests the default configuration.
// zh: TestDefaultConfig 測試預設配置。
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	if config == nil {
		t.Fatal("Default config should not be nil")
	}

	if config.ServiceName != "detectviz" {
		t.Errorf("Expected service name 'detectviz', got '%s'", config.ServiceName)
	}

	if config.ServiceVersion != "1.0.0" {
		t.Errorf("Expected service version '1.0.0', got '%s'", config.ServiceVersion)
	}

	if config.Environment != "development" {
		t.Errorf("Expected environment 'development', got '%s'", config.Environment)
	}

	if !config.Enabled {
		t.Error("Default config should be enabled")
	}

	if !config.Tracing.Enabled {
		t.Error("Default tracing should be enabled")
	}

	if !config.Metrics.Enabled {
		t.Error("Default metrics should be enabled")
	}

	if !config.Logging.Enabled {
		t.Error("Default logging should be enabled")
	}

	t.Log("Default config test passed")
}

// mockLoggerProvider is a mock implementation for testing.
// zh: mockLoggerProvider 是用於測試的模擬實作。
type mockLoggerProvider struct{}

func (m *mockLoggerProvider) Logger() contracts.Logger {
	return &mockLogger{}
}

func (m *mockLoggerProvider) WithContext(ctx context.Context) contracts.Logger {
	return &mockLogger{}
}

func (m *mockLoggerProvider) Flush() error {
	return nil
}

func (m *mockLoggerProvider) SetLevel(level string) error {
	return nil
}

func (m *mockLoggerProvider) Close() error {
	return nil
}

// mockLogger is a mock logger implementation for testing.
// zh: mockLogger 是用於測試的模擬日誌器實作。
type mockLogger struct{}

func (m *mockLogger) Debug(msg string, fields ...interface{}) {}
func (m *mockLogger) Info(msg string, fields ...interface{})  {}
func (m *mockLogger) Warn(msg string, fields ...interface{})  {}
func (m *mockLogger) Error(msg string, fields ...interface{}) {}
func (m *mockLogger) Fatal(msg string, fields ...interface{}) {}
