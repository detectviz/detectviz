package integration

import (
	"testing"

	"detectviz/pkg/config/schema"
)

// TestPluginConfigValidation tests plugin configuration validation functionality.
// zh: TestPluginConfigValidation 測試插件配置驗證功能。
func TestPluginConfigValidation(t *testing.T) {

	t.Run("ValidJWTConfig", func(t *testing.T) {
		// Valid JWT authenticator configuration
		validConfig := map[string]any{
			"secret_key":      "test-secret-key-123",
			"token_ttl":       3600,
			"refresh_enabled": true,
			"issuer":          "detectviz",
			"audience":        "detectviz-users",
		}

		validator := schema.NewConfigValidator()

		// Register JWT schema for testing
		registerJWTSchema(t, validator)

		result, err := validator.ValidatePluginConfig("jwt-authenticator", validConfig)
		if err != nil {
			t.Fatalf("Validation failed with error: %v", err)
		}

		if !result.Valid {
			t.Errorf("Expected valid JWT config to pass validation, got errors: %v", result.Errors)
		}
	})

	t.Run("InvalidJWTConfig_MissingRequired", func(t *testing.T) {
		// Invalid JWT config - missing required secret_key
		invalidConfig := map[string]any{
			"token_ttl":       3600,
			"refresh_enabled": true,
			"issuer":          "detectviz",
		}

		validator := schema.NewConfigValidator()
		registerJWTSchema(t, validator)

		result, err := validator.ValidatePluginConfig("jwt-authenticator", invalidConfig)
		if err != nil {
			t.Fatalf("Validation failed with error: %v", err)
		}

		if result.Valid {
			t.Errorf("Expected validation error for missing required field, but validation passed")
		}

		// Check that error mentions secret_key
		found := false
		for _, validationErr := range result.Errors {
			if containsString(validationErr.Error(), "secret_key") {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected error message to contain 'secret_key', got errors: %v", result.Errors)
		}
	})

	t.Run("InvalidJWTConfig_TypeMismatch", func(t *testing.T) {
		// Invalid JWT config - wrong type for token_ttl
		invalidConfig := map[string]any{
			"secret_key":      "test-secret-key-123",
			"token_ttl":       "invalid-number", // Should be int
			"refresh_enabled": true,
			"issuer":          "detectviz",
		}

		validator := schema.NewConfigValidator()
		registerJWTSchema(t, validator)

		result, err := validator.ValidatePluginConfig("jwt-authenticator", invalidConfig)
		if err != nil {
			t.Fatalf("Validation failed with error: %v", err)
		}

		if result.Valid {
			t.Errorf("Expected validation error for type mismatch, but validation passed")
		}

		t.Logf("Type mismatch errors: %v", result.Errors)
	})

	t.Run("ValidPrometheusConfig", func(t *testing.T) {
		// Valid Prometheus importer configuration
		validConfig := map[string]any{
			"endpoint":        "http://localhost:9090",
			"timeout":         "30s",
			"metrics":         []string{"cpu_usage", "memory_usage"},
			"scrape_interval": "60s",
		}

		validator := schema.NewConfigValidator()
		registerPrometheusSchema(t, validator)

		result, err := validator.ValidatePluginConfig("prometheus-importer", validConfig)
		if err != nil {
			t.Fatalf("Validation failed with error: %v", err)
		}

		if !result.Valid {
			t.Errorf("Expected valid Prometheus config to pass validation, got errors: %v", result.Errors)
		}
	})

	t.Run("InvalidPrometheusConfig_MissingEndpoint", func(t *testing.T) {
		// Invalid Prometheus config - missing required endpoint
		invalidConfig := map[string]any{
			"timeout":         "30s",
			"metrics":         []string{"cpu_usage"},
			"scrape_interval": "60s",
		}

		validator := schema.NewConfigValidator()
		registerPrometheusSchema(t, validator)

		result, err := validator.ValidatePluginConfig("prometheus-importer", invalidConfig)
		if err != nil {
			t.Fatalf("Validation failed with error: %v", err)
		}

		if result.Valid {
			t.Errorf("Expected validation error for missing endpoint, but validation passed")
		}

		t.Logf("Missing endpoint errors: %v", result.Errors)
	})

	t.Run("ValidSystemStatusConfig", func(t *testing.T) {
		// Valid system status plugin configuration
		validConfig := map[string]any{
			"title":        "系統狀態監控",
			"refresh_rate": 30,
			"show_memory":  true,
			"show_cpu":     true,
			"show_plugins": true,
		}

		validator := schema.NewConfigValidator()
		registerSystemStatusSchema(t, validator)

		result, err := validator.ValidatePluginConfig("system-status", validConfig)
		if err != nil {
			t.Fatalf("Validation failed with error: %v", err)
		}

		if !result.Valid {
			t.Errorf("Expected valid system status config to pass validation, got errors: %v", result.Errors)
		}
	})

	t.Run("InvalidSystemStatusConfig_InvalidRefreshRate", func(t *testing.T) {
		// Invalid system status config - refresh rate too low
		invalidConfig := map[string]any{
			"title":        "系統狀態監控",
			"refresh_rate": 0, // Should be positive
			"show_memory":  true,
			"show_cpu":     true,
			"show_plugins": true,
		}

		validator := schema.NewConfigValidator()
		registerSystemStatusSchema(t, validator)

		result, err := validator.ValidatePluginConfig("system-status", invalidConfig)
		if err != nil {
			t.Fatalf("Validation failed with error: %v", err)
		}

		if result.Valid {
			t.Errorf("Expected validation error for invalid refresh rate, but validation passed")
		}

		t.Logf("Invalid refresh rate errors: %v", result.Errors)
	})

	t.Run("UnknownPlugin", func(t *testing.T) {
		// Test validation with unknown plugin (should pass with no schema)
		config := map[string]any{
			"some_field": "some_value",
		}

		validator := schema.NewConfigValidator()
		result, err := validator.ValidatePluginConfig("unknown-plugin", config)
		if err != nil {
			t.Fatalf("Validation failed with error: %v", err)
		}

		// Unknown plugins should pass validation (no schema = no validation)
		if !result.Valid {
			t.Errorf("Expected unknown plugin to pass validation (no schema), got errors: %v", result.Errors)
		}
	})

	t.Run("DefaultValues", func(t *testing.T) {
		// Test default values application
		minimalConfig := map[string]any{
			"secret_key": "test-secret",
		}

		validator := schema.NewConfigValidator()
		registerJWTSchema(t, validator)

		result, err := validator.ValidatePluginConfig("jwt-authenticator", minimalConfig)
		if err != nil {
			t.Fatalf("Validation failed with error: %v", err)
		}

		if !result.Valid {
			t.Errorf("Expected minimal config with defaults to pass validation, got errors: %v", result.Errors)
		}

		// Check if defaults were applied in normalized config
		if tokenTTL, ok := result.Normalized["token_ttl"]; !ok || tokenTTL != 3600 {
			t.Errorf("Expected default token_ttl=3600, got %v", tokenTTL)
		}

		if refreshEnabled, ok := result.Normalized["refresh_enabled"]; !ok || refreshEnabled != false {
			t.Errorf("Expected default refresh_enabled=false, got %v", refreshEnabled)
		}

		t.Logf("Applied config with defaults: %+v", result.Normalized)
	})

	t.Run("ComplexNestedValidation", func(t *testing.T) {
		// Test complex nested structure validation
		complexConfig := map[string]any{
			"endpoint":        "http://localhost:9090",
			"timeout":         "30s",
			"scrape_interval": "15s",
			"metrics": []string{
				"up",
				"cpu_usage_percent",
				"memory_usage_bytes",
			},
		}

		validator := schema.NewConfigValidator()
		registerPrometheusSchema(t, validator)

		result, err := validator.ValidatePluginConfig("prometheus-importer", complexConfig)
		if err != nil {
			t.Fatalf("Validation failed with error: %v", err)
		}

		if !result.Valid {
			t.Errorf("Expected complex config to pass validation, got errors: %v", result.Errors)
		}
	})

	t.Run("ValidatorReuse", func(t *testing.T) {
		// Test validator reuse across multiple validations
		validator := schema.NewConfigValidator()
		registerJWTSchema(t, validator)

		configs := []map[string]any{
			{
				"secret_key": "test-secret-key-1",
				"token_ttl":  1800,
			},
			{
				"secret_key": "test-secret-key-2",
				"token_ttl":  7200,
			},
			{
				"secret_key": "test-secret-key-3",
				"token_ttl":  3600,
			},
		}

		for i, config := range configs {
			result, err := validator.ValidatePluginConfig("jwt-authenticator", config)
			if err != nil {
				t.Fatalf("Validation %d failed with error: %v", i+1, err)
			}
			if !result.Valid {
				t.Errorf("Validation %d failed: %v", i+1, result.Errors)
			}
		}
	})

	t.Run("ValidationErrorDetails", func(t *testing.T) {
		// Test detailed validation error reporting
		invalidConfig := map[string]any{
			"secret_key":      "",              // Empty required field
			"token_ttl":       "not-a-number",  // Wrong type
			"refresh_enabled": "not-a-boolean", // Wrong type
			"issuer":          123,             // Wrong type
		}

		validator := schema.NewConfigValidator()
		registerJWTSchema(t, validator)

		result, err := validator.ValidatePluginConfig("jwt-authenticator", invalidConfig)
		if err != nil {
			t.Fatalf("Validation failed with error: %v", err)
		}

		if result.Valid {
			t.Errorf("Expected validation error for multiple invalid fields, but validation passed")
		}

		// Check that we have multiple validation errors
		if len(result.Errors) == 0 {
			t.Errorf("Expected multiple validation errors, got none")
		}

		t.Logf("Detailed validation errors: %v", result.Errors)

		// Should contain information about multiple validation failures
		errorMessages := make([]string, len(result.Errors))
		for i, err := range result.Errors {
			errorMessages[i] = err.Error()
		}

		expectedFields := []string{"secret_key", "token_ttl"}
		for _, expectedField := range expectedFields {
			found := false
			for _, errorMsg := range errorMessages {
				if containsString(errorMsg, expectedField) {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected error message to contain '%s', but it didn't. All errors: %v", expectedField, errorMessages)
			}
		}
	})
}

// TestSchemaDefinition tests the schema definition functionality.
// zh: TestSchemaDefinition 測試模式定義功能。
func TestSchemaDefinition(t *testing.T) {
	t.Run("SchemaRegistration", func(t *testing.T) {
		validator := schema.NewConfigValidator()

		// Register a test schema
		testSchema := &schema.PluginSchema{
			Name:        "test-plugin",
			Version:     "1.0.0",
			Description: "Test plugin schema",
			Fields: map[string]*schema.FieldSchema{
				"required_field": {
					Type:        "string",
					Description: "A required field",
					Required:    true,
				},
				"optional_field": {
					Type:        "int",
					Description: "An optional field",
					Required:    false,
					Default:     42,
				},
			},
			Required: []string{"required_field"},
			Defaults: map[string]any{
				"optional_field": 42,
			},
		}

		err := validator.RegisterSchema("test-plugin", testSchema)
		if err != nil {
			t.Errorf("Failed to register schema: %v", err)
		}

		// Test validation with the registered schema
		validConfig := map[string]any{
			"required_field": "test-value",
		}

		result, err := validator.ValidatePluginConfig("test-plugin", validConfig)
		if err != nil {
			t.Fatalf("Validation failed with error: %v", err)
		}

		if !result.Valid {
			t.Errorf("Expected valid config to pass validation, got errors: %v", result.Errors)
		}

		// Check that default was applied
		if optionalField, ok := result.Normalized["optional_field"]; !ok || optionalField != 42 {
			t.Errorf("Expected default optional_field=42, got %v", optionalField)
		}
	})

	t.Run("SchemaValidation", func(t *testing.T) {
		validator := schema.NewConfigValidator()
		registerJWTSchema(t, validator)

		// Test schema-based validation
		testConfig := map[string]any{
			"secret_key": "test-key",
			"token_ttl":  1800,
		}

		result, err := validator.ValidatePluginConfig("jwt-authenticator", testConfig)
		if err != nil {
			t.Fatalf("Schema-based validation failed with error: %v", err)
		}

		if !result.Valid {
			t.Errorf("Schema-based validation failed: %v", result.Errors)
		}
	})
}

// Helper functions
// zh: 輔助函式

// registerJWTSchema registers JWT plugin schema for testing
// zh: registerJWTSchema 為測試註冊 JWT 插件模式
func registerJWTSchema(t *testing.T, validator *schema.ConfigValidator) {
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
			"token_ttl": {
				Type:        "int",
				Description: "Token time to live in seconds",
				Required:    false,
				Default:     3600,
				MinValue:    &[]int64{1}[0],
			},
			"refresh_enabled": {
				Type:        "bool",
				Description: "Enable refresh token",
				Required:    false,
				Default:     false,
			},
			"issuer": {
				Type:        "string",
				Description: "JWT issuer",
				Required:    false,
			},
			"audience": {
				Type:        "string",
				Description: "JWT audience",
				Required:    false,
			},
		},
		Required: []string{"secret_key"},
		Defaults: map[string]any{
			"token_ttl":       3600,
			"refresh_enabled": false,
		},
	}

	err := validator.RegisterSchema("jwt-authenticator", jwtSchema)
	if err != nil {
		t.Fatalf("Failed to register JWT schema: %v", err)
	}
}

// registerPrometheusSchema registers Prometheus plugin schema for testing
// zh: registerPrometheusSchema 為測試註冊 Prometheus 插件模式
func registerPrometheusSchema(t *testing.T, validator *schema.ConfigValidator) {
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
			"metrics": {
				Type:        "array",
				Description: "List of metrics to scrape",
				Required:    false,
			},
		},
		Required: []string{"endpoint"},
		Defaults: map[string]any{
			"scrape_interval": "15s",
			"timeout":         "10s",
		},
	}

	err := validator.RegisterSchema("prometheus-importer", promSchema)
	if err != nil {
		t.Fatalf("Failed to register Prometheus schema: %v", err)
	}
}

// registerSystemStatusSchema registers system status plugin schema for testing
// zh: registerSystemStatusSchema 為測試註冊系統狀態插件模式
func registerSystemStatusSchema(t *testing.T, validator *schema.ConfigValidator) {
	statusSchema := &schema.PluginSchema{
		Name:        "system-status",
		Version:     "1.0.0",
		Description: "System status plugin schema",
		Fields: map[string]*schema.FieldSchema{
			"title": {
				Type:        "string",
				Description: "Page title",
				Required:    false,
				Default:     "System Status",
			},
			"refresh_rate": {
				Type:        "int",
				Description: "Refresh rate in seconds",
				Required:    false,
				Default:     30,
				MinValue:    &[]int64{1}[0],
			},
			"show_memory": {
				Type:        "bool",
				Description: "Show memory usage",
				Required:    false,
				Default:     true,
			},
			"show_cpu": {
				Type:        "bool",
				Description: "Show CPU usage",
				Required:    false,
				Default:     true,
			},
			"show_plugins": {
				Type:        "bool",
				Description: "Show plugin status",
				Required:    false,
				Default:     true,
			},
		},
		Required: []string{},
		Defaults: map[string]any{
			"title":        "System Status",
			"refresh_rate": 30,
			"show_memory":  true,
			"show_cpu":     true,
			"show_plugins": true,
		},
	}

	err := validator.RegisterSchema("system-status", statusSchema)
	if err != nil {
		t.Fatalf("Failed to register system status schema: %v", err)
	}
}

// containsString checks if a string contains a substring.
// zh: containsString 檢查字串是否包含子字串。
func containsString(str, substr string) bool {
	return len(str) >= len(substr) &&
		(len(substr) == 0 || str == substr ||
			len(str) > len(substr) &&
				(str[:len(substr)] == substr ||
					str[len(str)-len(substr):] == substr ||
					containsSubstr(str, substr)))
}

// containsSubstr is a helper for substring checking.
// zh: containsSubstr 是子字串檢查的輔助函式。
func containsSubstr(str, substr string) bool {
	if len(substr) > len(str) {
		return false
	}
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
