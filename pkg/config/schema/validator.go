package schema

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"detectviz/pkg/platform/contracts"
)

// ConfigValidator provides configuration validation functionality.
// zh: ConfigValidator 提供配置驗證功能。
type ConfigValidator struct {
	schemas map[string]*PluginSchema
}

// PluginSchema defines the validation schema for a plugin.
// zh: PluginSchema 定義插件的驗證模式。
type PluginSchema struct {
	Name        string                  `json:"name"`
	Version     string                  `json:"version"`
	Description string                  `json:"description"`
	Fields      map[string]*FieldSchema `json:"fields"`
	Required    []string                `json:"required"`
	Optional    []string                `json:"optional"`
	Defaults    map[string]any          `json:"defaults"`
}

// FieldSchema defines the validation schema for a configuration field.
// zh: FieldSchema 定義配置欄位的驗證模式。
type FieldSchema struct {
	Type        string   `json:"type"` // string, int, bool, duration, array, object
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Default     any      `json:"default"`
	MinLength   *int     `json:"min_length,omitempty"`
	MaxLength   *int     `json:"max_length,omitempty"`
	MinValue    *int64   `json:"min_value,omitempty"`
	MaxValue    *int64   `json:"max_value,omitempty"`
	Pattern     string   `json:"pattern,omitempty"`
	Enum        []string `json:"enum,omitempty"`
	Format      string   `json:"format,omitempty"` // email, url, duration, etc.
}

// ValidationError represents a configuration validation error.
// zh: ValidationError 代表配置驗證錯誤。
type ValidationError struct {
	Field   string `json:"field"`
	Value   any    `json:"value"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// Error implements the error interface.
// zh: Error 實作 error 介面。
func (ve *ValidationError) Error() string {
	return fmt.Sprintf("validation error in field '%s': %s", ve.Field, ve.Message)
}

// ValidationResult contains the validation results.
// zh: ValidationResult 包含驗證結果。
type ValidationResult struct {
	Valid      bool               `json:"valid"`
	Errors     []*ValidationError `json:"errors,omitempty"`
	Warnings   []*ValidationError `json:"warnings,omitempty"`
	Normalized map[string]any     `json:"normalized,omitempty"`
}

// NewConfigValidator creates a new configuration validator.
// zh: NewConfigValidator 建立新的配置驗證器。
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		schemas: make(map[string]*PluginSchema),
	}
}

// RegisterSchema registers a validation schema for a plugin.
// zh: RegisterSchema 為插件註冊驗證模式。
func (cv *ConfigValidator) RegisterSchema(pluginName string, schema *PluginSchema) error {
	if pluginName == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}
	if schema == nil {
		return fmt.Errorf("schema cannot be nil")
	}

	cv.schemas[pluginName] = schema
	return nil
}

// ValidatePluginConfig validates plugin configuration against its schema.
// zh: ValidatePluginConfig 根據模式驗證插件配置。
func (cv *ConfigValidator) ValidatePluginConfig(pluginName string, config map[string]any) (*ValidationResult, error) {
	schema, exists := cv.schemas[pluginName]
	if !exists {
		return &ValidationResult{
			Valid:      true,
			Normalized: config,
		}, nil // No schema means no validation
	}

	return cv.validateWithSchema(schema, config)
}

// ValidateMetadata validates plugin metadata.
// zh: ValidateMetadata 驗證插件元資料。
func (cv *ConfigValidator) ValidateMetadata(metadata *contracts.PluginMetadata) (*ValidationResult, error) {
	if metadata == nil {
		return &ValidationResult{
			Valid: false,
			Errors: []*ValidationError{
				{
					Field:   "metadata",
					Message: "metadata cannot be nil",
					Code:    "required",
				},
			},
		}, nil
	}

	var errors []*ValidationError

	// Validate required fields
	if metadata.Name == "" {
		errors = append(errors, &ValidationError{
			Field:   "name",
			Message: "plugin name is required",
			Code:    "required",
		})
	}

	if metadata.Version == "" {
		errors = append(errors, &ValidationError{
			Field:   "version",
			Message: "plugin version is required",
			Code:    "required",
		})
	}

	if metadata.Type == "" {
		errors = append(errors, &ValidationError{
			Field:   "type",
			Message: "plugin type is required",
			Code:    "required",
		})
	}

	// Validate plugin type
	validTypes := []string{"importer", "exporter", "auth", "middleware", "integration", "tool"}
	if metadata.Type != "" && !contains(validTypes, metadata.Type) {
		errors = append(errors, &ValidationError{
			Field:   "type",
			Value:   metadata.Type,
			Message: fmt.Sprintf("invalid plugin type, must be one of: %s", strings.Join(validTypes, ", ")),
			Code:    "invalid_enum",
		})
	}

	// Validate plugin category
	validCategories := []string{"core", "community", "custom"}
	if metadata.Category != "" && !contains(validCategories, metadata.Category) {
		errors = append(errors, &ValidationError{
			Field:   "category",
			Value:   metadata.Category,
			Message: fmt.Sprintf("invalid plugin category, must be one of: %s", strings.Join(validCategories, ", ")),
			Code:    "invalid_enum",
		})
	}

	// Validate config if schema exists
	var configErrors []*ValidationError
	if metadata.Config != nil {
		configResult, err := cv.ValidatePluginConfig(metadata.Name, metadata.Config)
		if err != nil {
			return nil, err
		}
		if !configResult.Valid {
			configErrors = configResult.Errors
		}
	}

	return &ValidationResult{
		Valid:  len(errors) == 0 && len(configErrors) == 0,
		Errors: append(errors, configErrors...),
	}, nil
}

// validateWithSchema performs validation against a plugin schema.
// zh: validateWithSchema 根據插件模式執行驗證。
func (cv *ConfigValidator) validateWithSchema(schema *PluginSchema, config map[string]any) (*ValidationResult, error) {
	var errors []*ValidationError
	var warnings []*ValidationError
	normalized := make(map[string]any)

	// Apply defaults first
	for key, defaultValue := range schema.Defaults {
		if _, exists := config[key]; !exists {
			normalized[key] = defaultValue
		}
	}

	// Copy provided config
	for key, value := range config {
		normalized[key] = value
	}

	// Validate required fields
	for _, requiredField := range schema.Required {
		if _, exists := normalized[requiredField]; !exists {
			errors = append(errors, &ValidationError{
				Field:   requiredField,
				Message: fmt.Sprintf("required field '%s' is missing", requiredField),
				Code:    "required",
			})
		}
	}

	// Validate each field
	for fieldName, fieldSchema := range schema.Fields {
		value, exists := normalized[fieldName]
		if !exists {
			if fieldSchema.Required {
				errors = append(errors, &ValidationError{
					Field:   fieldName,
					Message: fmt.Sprintf("required field '%s' is missing", fieldName),
					Code:    "required",
				})
			} else if fieldSchema.Default != nil {
				normalized[fieldName] = fieldSchema.Default
			}
			continue
		}

		// Validate field value
		fieldErrors := cv.validateField(fieldName, value, fieldSchema)
		errors = append(errors, fieldErrors...)

		// Type conversion and normalization
		if normalizedValue, err := cv.normalizeValue(value, fieldSchema); err == nil {
			normalized[fieldName] = normalizedValue
		}
	}

	// Check for unknown fields
	for fieldName := range config {
		if _, exists := schema.Fields[fieldName]; !exists {
			warnings = append(warnings, &ValidationError{
				Field:   fieldName,
				Value:   config[fieldName],
				Message: fmt.Sprintf("unknown field '%s'", fieldName),
				Code:    "unknown_field",
			})
		}
	}

	return &ValidationResult{
		Valid:      len(errors) == 0,
		Errors:     errors,
		Warnings:   warnings,
		Normalized: normalized,
	}, nil
}

// validateField validates a single field value.
// zh: validateField 驗證單一欄位值。
func (cv *ConfigValidator) validateField(fieldName string, value any, schema *FieldSchema) []*ValidationError {
	var errors []*ValidationError

	// Type validation
	if !cv.isValidType(value, schema.Type) {
		errors = append(errors, &ValidationError{
			Field:   fieldName,
			Value:   value,
			Message: fmt.Sprintf("expected type %s, got %T", schema.Type, value),
			Code:    "invalid_type",
		})
		return errors // Return early if type is wrong
	}

	// String validations
	if schema.Type == "string" {
		if str, ok := value.(string); ok {
			if schema.MinLength != nil && len(str) < *schema.MinLength {
				errors = append(errors, &ValidationError{
					Field:   fieldName,
					Value:   value,
					Message: fmt.Sprintf("string length %d is less than minimum %d", len(str), *schema.MinLength),
					Code:    "min_length",
				})
			}
			if schema.MaxLength != nil && len(str) > *schema.MaxLength {
				errors = append(errors, &ValidationError{
					Field:   fieldName,
					Value:   value,
					Message: fmt.Sprintf("string length %d exceeds maximum %d", len(str), *schema.MaxLength),
					Code:    "max_length",
				})
			}
			if len(schema.Enum) > 0 && !contains(schema.Enum, str) {
				errors = append(errors, &ValidationError{
					Field:   fieldName,
					Value:   value,
					Message: fmt.Sprintf("value must be one of: %s", strings.Join(schema.Enum, ", ")),
					Code:    "invalid_enum",
				})
			}
			if schema.Format == "duration" {
				if _, err := time.ParseDuration(str); err != nil {
					errors = append(errors, &ValidationError{
						Field:   fieldName,
						Value:   value,
						Message: fmt.Sprintf("invalid duration format: %v", err),
						Code:    "invalid_format",
					})
				}
			}
		}
	}

	// Integer validations
	if schema.Type == "int" {
		if num, ok := convertToInt64(value); ok {
			if schema.MinValue != nil && num < *schema.MinValue {
				errors = append(errors, &ValidationError{
					Field:   fieldName,
					Value:   value,
					Message: fmt.Sprintf("value %d is less than minimum %d", num, *schema.MinValue),
					Code:    "min_value",
				})
			}
			if schema.MaxValue != nil && num > *schema.MaxValue {
				errors = append(errors, &ValidationError{
					Field:   fieldName,
					Value:   value,
					Message: fmt.Sprintf("value %d exceeds maximum %d", num, *schema.MaxValue),
					Code:    "max_value",
				})
			}
		}
	}

	return errors
}

// isValidType checks if value matches the expected type.
// zh: isValidType 檢查值是否符合預期類型。
func (cv *ConfigValidator) isValidType(value any, expectedType string) bool {
	if value == nil {
		return true // nil is valid for any type (will be handled by required check)
	}

	switch expectedType {
	case "string":
		_, ok := value.(string)
		return ok
	case "int":
		return isIntegerType(value)
	case "bool":
		_, ok := value.(bool)
		return ok
	case "duration":
		if str, ok := value.(string); ok {
			_, err := time.ParseDuration(str)
			return err == nil
		}
		return false
	case "array":
		rv := reflect.ValueOf(value)
		return rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array
	case "object":
		rv := reflect.ValueOf(value)
		return rv.Kind() == reflect.Map
	default:
		return true // Unknown types are accepted
	}
}

// normalizeValue converts value to the appropriate normalized form.
// zh: normalizeValue 將值轉換為適當的標準化形式。
func (cv *ConfigValidator) normalizeValue(value any, schema *FieldSchema) (any, error) {
	switch schema.Type {
	case "int":
		if num, ok := convertToInt64(value); ok {
			return int(num), nil
		}
	case "duration":
		if str, ok := value.(string); ok {
			if duration, err := time.ParseDuration(str); err == nil {
				return duration, nil
			}
		}
	}
	return value, nil
}

// Helper functions

// contains checks if a slice contains a specific string.
// zh: contains 檢查切片是否包含特定字串。
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// isIntegerType checks if value is any integer type.
// zh: isIntegerType 檢查值是否為任何整數類型。
func isIntegerType(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64:
		// Check if float is actually an integer
		if f, ok := value.(float64); ok {
			return f == float64(int64(f))
		}
		if f, ok := value.(float32); ok {
			return f == float32(int32(f))
		}
	}
	return false
}

// convertToInt64 converts various numeric types to int64.
// zh: convertToInt64 將各種數值類型轉換為 int64。
func convertToInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case float32:
		if v == float32(int32(v)) {
			return int64(v), true
		}
	case float64:
		if v == float64(int64(v)) {
			return int64(v), true
		}
	}
	return 0, false
}
