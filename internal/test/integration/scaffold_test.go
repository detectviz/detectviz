package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"detectviz/internal/platform/composition"
	"detectviz/internal/platform/registry"
	"detectviz/pkg/config/loader"
	"detectviz/pkg/config/schema"
	"detectviz/pkg/platform/contracts"

	// Import plugins for testing
	prometheusPlugin "detectviz/plugins/community/importers/prometheus"
	jwtPlugin "detectviz/plugins/core/auth/jwt"
)

// TestScaffoldIntegration tests the complete scaffold integration flow.
// zh: TestScaffoldIntegration 測試完整的 scaffold 整合流程。
func TestScaffoldIntegration(t *testing.T) {
	// Setup test environment
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create test registry
	registryManager := registry.NewManager()

	// Create dependency resolver
	resolver := composition.NewDependencyResolver()

	// Create lifecycle manager
	lifecycleManager := composition.NewLifecycleManager(resolver)

	// Create config loader
	configLoader := loader.NewConfigLoader(".")

	// Create config validator
	configValidator := schema.NewConfigValidator()

	t.Run("RegisterPlugins", func(t *testing.T) {
		// Register plugins manually (in real scenario, they would be auto-discovered)
		err := jwtPlugin.Register(registryManager)
		if err != nil {
			t.Fatalf("Failed to register JWT plugin: %v", err)
		}

		err = prometheusPlugin.Register(registryManager)
		if err != nil {
			t.Fatalf("Failed to register Prometheus plugin: %v", err)
		}

		// Verify plugins are registered
		plugins := registryManager.ListPlugins()
		expectedPlugins := []string{"jwt-authenticator", "prometheus-importer"}

		for _, expected := range expectedPlugins {
			found := false
			for _, plugin := range plugins {
				if plugin == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected plugin %s not found in registry", expected)
			}
		}
	})

	t.Run("LoadConfiguration", func(t *testing.T) {
		// Create test composition file
		testCompositionPath := createTestCompositionFile(t)
		defer os.Remove(testCompositionPath)

		// Load composition
		config, err := configLoader.LoadComposition(testCompositionPath)
		if err != nil {
			t.Fatalf("Failed to load composition: %v", err)
		}

		if config.Metadata.Name != "test-composition" {
			t.Errorf("Expected composition name 'test-composition', got '%s'", config.Metadata.Name)
		}

		// Validate configuration
		err = configLoader.ValidateConfiguration()
		if err != nil {
			t.Fatalf("Configuration validation failed: %v", err)
		}
	})

	t.Run("ValidatePluginConfigs", func(t *testing.T) {
		// Register schemas for validation
		registerTestSchemas(t, configValidator)

		// Create test plugin metadata
		jwtMetadata := &contracts.PluginMetadata{
			Name:     "jwt-authenticator",
			Version:  "1.0.0",
			Type:     "auth",
			Category: "core",
			Config: map[string]any{
				"secret_key":     "test-secret",
				"issuer":         "test-issuer",
				"expiry_time":    "24h",
				"signing_method": "HS256",
			},
		}

		result, err := configValidator.ValidateMetadata(jwtMetadata)
		if err != nil {
			t.Fatalf("Failed to validate JWT metadata: %v", err)
		}

		if !result.Valid {
			t.Errorf("JWT metadata validation failed: %v", result.Errors)
		}
	})

	t.Run("LifecycleManagement", func(t *testing.T) {
		// Initialize lifecycle manager
		err := lifecycleManager.Initialize(ctx)
		if err != nil {
			t.Fatalf("Failed to initialize lifecycle manager: %v", err)
		}

		// Start all plugins
		err = lifecycleManager.StartAll(ctx, registryManager)
		if err != nil {
			t.Fatalf("Failed to start plugins: %v", err)
		}

		// Check lifecycle manager status
		status := lifecycleManager.GetStatus()
		if status != contracts.StatusRunning {
			t.Errorf("Expected lifecycle manager status to be running, got %s", status)
		}

		// Perform health check
		healthStatus := lifecycleManager.HealthCheck(ctx, registryManager)

		for pluginName, health := range healthStatus {
			if health.Status != "healthy" {
				t.Errorf("Plugin %s is not healthy: %s", pluginName, health.Message)
			}
		}

		// Shutdown all plugins
		err = lifecycleManager.ShutdownAll(ctx, registryManager)
		if err != nil {
			t.Fatalf("Failed to shutdown plugins: %v", err)
		}
	})

	t.Run("PluginInstanceCreation", func(t *testing.T) {
		// Test JWT plugin creation and initialization
		testConfig := map[string]any{
			"secret_key":     "test-key-123",
			"issuer":         "test-detectviz",
			"expiry_time":    "12h",
			"signing_method": "HS256",
		}

		jwtPluginInstance, err := jwtPlugin.NewJWTAuthenticator(testConfig)
		if err != nil {
			t.Fatalf("Failed to create JWT plugin: %v", err)
		}

		// Initialize the plugin first
		err = jwtPluginInstance.Init(testConfig)
		if err != nil {
			t.Fatalf("Failed to initialize JWT plugin: %v", err)
		}

		// Complete full lifecycle sequence
		if lifecycleAware, ok := jwtPluginInstance.(contracts.LifecycleAware); ok {
			err = lifecycleAware.OnRegister()
			if err != nil {
				t.Errorf("OnRegister failed: %v", err)
			}

			err = lifecycleAware.OnStart()
			if err != nil {
				t.Errorf("OnStart failed: %v", err)
			}

			// Test health checker interface after complete initialization
			if healthChecker, ok := jwtPluginInstance.(contracts.HealthChecker); ok {
				health := healthChecker.CheckHealth(ctx)
				if health.Status != "healthy" {
					t.Logf("JWT plugin health details: %+v", health)
					t.Errorf("JWT plugin health check failed: %s", health.Message)
				}
			}

			err = lifecycleAware.OnShutdown()
			if err != nil {
				t.Errorf("OnShutdown failed: %v", err)
			}
		}
	})

	t.Run("ConfigParsing", func(t *testing.T) {
		// Test JWT plugin config parsing
		testConfig := map[string]any{
			"secret_key":     "test-key-123",
			"issuer":         "test-detectviz",
			"expiry_time":    "12h",
			"signing_method": "HS256",
		}

		jwtPluginInstance, err := jwtPlugin.NewJWTAuthenticator(testConfig)
		if err != nil {
			t.Fatalf("Failed to create JWT plugin with config: %v", err)
		}

		err = jwtPluginInstance.Init(testConfig)
		if err != nil {
			t.Fatalf("Failed to initialize JWT plugin: %v", err)
		}

		// Test Prometheus plugin config parsing
		promConfig := map[string]any{
			"endpoint":        "http://test:9090",
			"scrape_interval": "30s",
			"timeout":         "15s",
		}

		promPluginInstance, err := prometheusPlugin.NewPrometheusImporter(promConfig)
		if err != nil {
			t.Fatalf("Failed to create Prometheus plugin with config: %v", err)
		}

		err = promPluginInstance.Init(promConfig)
		if err != nil {
			t.Fatalf("Failed to initialize Prometheus plugin: %v", err)
		}
	})

	t.Run("PluginDiscovery", func(t *testing.T) {
		// Create plugin discovery
		discoveryConfig := &registry.DiscoveryConfig{
			BasePath:       "../../..",
			ScanPaths:      []string{"plugins"},
			AutoRegister:   false, // Don't auto-register in tests
			LoadSharedLibs: false,
			ScanGoSources:  true,
		}

		discovery := registry.NewPluginDiscovery(registryManager, discoveryConfig)

		// Discover plugins
		discoveredPlugins, err := discovery.DiscoverPlugins()
		if err != nil {
			t.Fatalf("Plugin discovery failed: %v", err)
		}

		if len(discoveredPlugins) == 0 {
			t.Errorf("No plugins discovered")
		}

		// Check if expected plugins were discovered
		discoveredNames := make(map[string]bool)
		for _, plugin := range discoveredPlugins {
			discoveredNames[plugin.Name] = true
			t.Logf("Discovered plugin: %s (type: %s, category: %s)", plugin.Name, plugin.Type, plugin.Category)
		}

		expectedPlugins := []string{"jwt-authenticator", "prometheus-importer"}
		for _, expected := range expectedPlugins {
			if !discoveredNames[expected] {
				t.Errorf("Expected plugin %s was not discovered. Available plugins: %v", expected, discoveredNames)
			}
		}
	})
}

// TestEnabledFlag tests the enabled/disabled plugin functionality.
// zh: TestEnabledFlag 測試啟用/停用插件功能。
func TestEnabledFlag(t *testing.T) {
	configLoader := loader.NewConfigLoader(".")

	// Create test composition with disabled plugin
	testCompositionPath := createTestCompositionWithDisabledPlugin(t)
	defer os.Remove(testCompositionPath)

	// Load composition
	_, err := configLoader.LoadComposition(testCompositionPath)
	if err != nil {
		t.Fatalf("Failed to load composition: %v", err)
	}

	// Get plugin configs (should only return enabled plugins)
	pluginConfigs, err := configLoader.GetPluginConfigs()
	if err != nil {
		t.Fatalf("Failed to get plugin configs: %v", err)
	}

	// Should only have JWT plugin (Prometheus is disabled)
	if len(pluginConfigs) != 1 {
		t.Errorf("Expected 1 enabled plugin, got %d", len(pluginConfigs))
	}

	if pluginConfigs[0].Name != "jwt-authenticator" {
		t.Errorf("Expected jwt-authenticator to be enabled, got %s", pluginConfigs[0].Name)
	}
}

// Helper functions

// createTestCompositionFile creates a test composition file for testing.
// zh: createTestCompositionFile 建立測試用的組合檔案。
func createTestCompositionFile(t *testing.T) string {
	content := `
apiVersion: v1
kind: Composition
metadata:
  name: test-composition
  description: Test composition for integration tests
  version: 1.0.0

spec:
  platform:
    registry:
      enabled: true
      type: memory
    lifecycle:
      enabled: true
      timeout: 30s
    composition:
      enabled: true
      validation: strict

  core_plugins:
    - name: jwt-authenticator
      type: auth
      enabled: true
      config:
        secret_key: "test-secret-key"
        issuer: "test-detectviz"
        expiry_time: "24h"
        signing_method: "HS256"

  community_plugins:
    - name: prometheus-importer
      type: importer
      enabled: true
      config:
        endpoint: "http://localhost:9090"
        scrape_interval: "15s"
        timeout: "10s"

  applications:
    - name: server
      enabled: true
      config:
        port: 8080
        host: "0.0.0.0"
`

	tmpFile, err := os.CreateTemp("", "test-composition-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}

// createTestCompositionWithDisabledPlugin creates a test composition with disabled plugin.
// zh: createTestCompositionWithDisabledPlugin 建立包含停用插件的測試組合。
func createTestCompositionWithDisabledPlugin(t *testing.T) string {
	content := `
apiVersion: v1
kind: Composition
metadata:
  name: test-composition-disabled
  description: Test composition with disabled plugin
  version: 1.0.0

spec:
  platform:
    registry:
      enabled: true
      type: memory

  core_plugins:
    - name: jwt-authenticator
      type: auth
      enabled: true
      config:
        secret_key: "test-secret-key"
        issuer: "test-detectviz"

  community_plugins:
    - name: prometheus-importer
      type: importer
      enabled: false
      config:
        endpoint: "http://localhost:9090"
`

	tmpFile, err := os.CreateTemp("", "test-composition-disabled-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	return tmpFile.Name()
}

// registerTestSchemas registers validation schemas for testing.
// zh: registerTestSchemas 註冊測試用的驗證模式。
func registerTestSchemas(t *testing.T, validator *schema.ConfigValidator) {
	// JWT plugin schema
	jwtSchema := &schema.PluginSchema{
		Name:        "jwt-authenticator",
		Version:     "1.0.0",
		Description: "JWT authentication plugin schema",
		Fields: map[string]*schema.FieldSchema{
			"secret_key": {
				Type:        "string",
				Description: "Secret key for JWT signing",
				Required:    true,
				MinLength:   &[]int{8}[0],
			},
			"issuer": {
				Type:        "string",
				Description: "JWT issuer",
				Required:    true,
			},
			"expiry_time": {
				Type:        "string",
				Description: "Token expiry time",
				Required:    false,
				Default:     "24h",
				Format:      "duration",
			},
			"signing_method": {
				Type:        "string",
				Description: "JWT signing method",
				Required:    false,
				Default:     "HS256",
				Enum:        []string{"HS256", "HS384", "HS512", "RS256"},
			},
		},
		Required: []string{"secret_key", "issuer"},
		Defaults: map[string]any{
			"expiry_time":    "24h",
			"signing_method": "HS256",
		},
	}

	err := validator.RegisterSchema("jwt-authenticator", jwtSchema)
	if err != nil {
		t.Fatalf("Failed to register JWT schema: %v", err)
	}

	// Prometheus plugin schema
	promSchema := &schema.PluginSchema{
		Name:        "prometheus-importer",
		Version:     "1.0.0",
		Description: "Prometheus importer plugin schema",
		Fields: map[string]*schema.FieldSchema{
			"endpoint": {
				Type:        "string",
				Description: "Prometheus endpoint URL",
				Required:    true,
			},
			"scrape_interval": {
				Type:        "string",
				Description: "Scrape interval",
				Required:    false,
				Default:     "15s",
				Format:      "duration",
			},
			"timeout": {
				Type:        "string",
				Description: "Request timeout",
				Required:    false,
				Default:     "10s",
				Format:      "duration",
			},
		},
		Required: []string{"endpoint"},
		Defaults: map[string]any{
			"scrape_interval": "15s",
			"timeout":         "10s",
		},
	}

	err = validator.RegisterSchema("prometheus-importer", promSchema)
	if err != nil {
		t.Fatalf("Failed to register Prometheus schema: %v", err)
	}
}
