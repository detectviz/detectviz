package jwt

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"detectviz/pkg/platform/contracts"
)

// JWTAuthenticator implements JWT-based authentication.
// zh: JWTAuthenticator 實作基於 JWT 的認證功能。
type JWTAuthenticator struct {
	name        string
	version     string
	description string
	config      *JWTConfig
	initialized bool
}

// JWTConfig defines the configuration for JWT authentication.
// zh: JWTConfig 定義 JWT 認證的配置。
type JWTConfig struct {
	SecretKey     string `yaml:"secret_key" json:"secret_key" mapstructure:"secret_key"`
	Issuer        string `yaml:"issuer" json:"issuer" mapstructure:"issuer"`
	ExpiryTime    string `yaml:"expiry_time" json:"expiry_time" mapstructure:"expiry_time"`
	SigningMethod string `yaml:"signing_method" json:"signing_method" mapstructure:"signing_method"`
}

// NewJWTAuthenticator creates a new JWT authenticator instance.
// zh: NewJWTAuthenticator 建立新的 JWT 認證器實例。
func NewJWTAuthenticator(config any) (contracts.Plugin, error) {
	jwtConfig := &JWTConfig{
		SecretKey:     "default-secret-key",
		Issuer:        "detectviz",
		ExpiryTime:    "24h",
		SigningMethod: "HS256",
	}

	// Parse config from the provided config parameter
	if config != nil {
		if err := parsePluginConfig(config, jwtConfig); err != nil {
			return nil, fmt.Errorf("failed to parse JWT config: %w", err)
		}
	}

	return &JWTAuthenticator{
		name:        "jwt-authenticator",
		version:     "1.0.0",
		description: "JWT-based authentication plugin",
		config:      jwtConfig,
		initialized: false,
	}, nil
}

// parsePluginConfig parses the plugin configuration from various formats
// zh: parsePluginConfig 從各種格式解析插件配置
func parsePluginConfig(config any, target *JWTConfig) error {
	if config == nil {
		return nil
	}

	// Handle map[string]any format
	if configMap, ok := config.(map[string]any); ok {
		if secretKey, exists := configMap["secret_key"]; exists {
			if str, ok := secretKey.(string); ok {
				target.SecretKey = str
			}
		}
		if issuer, exists := configMap["issuer"]; exists {
			if str, ok := issuer.(string); ok {
				target.Issuer = str
			}
		}
		if expiryTime, exists := configMap["expiry_time"]; exists {
			if str, ok := expiryTime.(string); ok {
				target.ExpiryTime = str
			}
		}
		if signingMethod, exists := configMap["signing_method"]; exists {
			if str, ok := signingMethod.(string); ok {
				target.SigningMethod = str
			}
		}
		return nil
	}

	// Handle struct format using reflection
	if reflect.TypeOf(config).Kind() == reflect.Struct {
		configValue := reflect.ValueOf(config)
		targetValue := reflect.ValueOf(target).Elem()

		for i := 0; i < configValue.NumField(); i++ {
			field := configValue.Type().Field(i)
			value := configValue.Field(i)

			// Find corresponding field in target
			if targetField := targetValue.FieldByName(field.Name); targetField.IsValid() && targetField.CanSet() {
				if value.Type().AssignableTo(targetField.Type()) {
					targetField.Set(value)
				}
			}
		}
	}

	return nil
}

// Name returns the plugin name.
// zh: Name 回傳插件名稱。
func (j *JWTAuthenticator) Name() string {
	return j.name
}

// Version returns the plugin version.
// zh: Version 回傳插件版本。
func (j *JWTAuthenticator) Version() string {
	return j.version
}

// Description returns the plugin description.
// zh: Description 回傳插件描述。
func (j *JWTAuthenticator) Description() string {
	return j.description
}

// Init initializes the JWT authenticator.
// zh: Init 初始化 JWT 認證器。
func (j *JWTAuthenticator) Init(config any) error {
	if j.initialized {
		return nil
	}

	// Validate configuration
	if j.config.SecretKey == "" {
		return fmt.Errorf("JWT secret key is required")
	}

	if j.config.Issuer == "" {
		return fmt.Errorf("JWT issuer is required")
	}

	// Validate expiry time format
	if _, err := time.ParseDuration(j.config.ExpiryTime); err != nil {
		return fmt.Errorf("invalid expiry time format: %w", err)
	}

	j.initialized = true
	return nil
}

// Shutdown shuts down the JWT authenticator.
// zh: Shutdown 關閉 JWT 認證器。
func (j *JWTAuthenticator) Shutdown() error {
	j.initialized = false
	return nil
}

// LifecycleAware interface implementation
// zh: LifecycleAware 介面實作

// OnRegister is called when the plugin is registered.
// zh: OnRegister 在插件註冊時呼叫。
func (j *JWTAuthenticator) OnRegister() error {
	return nil
}

// OnStart is called when the plugin is started.
// zh: OnStart 在插件啟動時呼叫。
func (j *JWTAuthenticator) OnStart() error {
	if !j.initialized {
		return fmt.Errorf("JWT authenticator not initialized")
	}
	return nil
}

// OnStop is called when the plugin is stopped.
// zh: OnStop 在插件停止時呼叫。
func (j *JWTAuthenticator) OnStop() error {
	return nil
}

// OnShutdown is called when the plugin is shutdown.
// zh: OnShutdown 在插件關閉時呼叫。
func (j *JWTAuthenticator) OnShutdown() error {
	return j.Shutdown()
}

// HealthChecker interface implementation
// zh: HealthChecker 介面實作

// CheckHealth checks the health of the JWT authenticator.
// zh: CheckHealth 檢查 JWT 認證器的健康狀況。
func (j *JWTAuthenticator) CheckHealth(ctx context.Context) contracts.HealthStatus {
	status := contracts.HealthStatus{
		Timestamp: time.Now(),
		Details:   make(map[string]any),
	}

	if !j.initialized {
		status.Status = "unhealthy"
		status.Message = "JWT authenticator not initialized"
		return status
	}

	// Check configuration validity
	if j.config.SecretKey == "" {
		status.Status = "unhealthy"
		status.Message = "JWT secret key not configured"
		return status
	}

	// Check expiry time format
	if _, err := time.ParseDuration(j.config.ExpiryTime); err != nil {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("Invalid expiry time format: %v", err)
		return status
	}

	status.Status = "healthy"
	status.Message = "JWT authenticator is healthy"
	status.Details["issuer"] = j.config.Issuer
	status.Details["signing_method"] = j.config.SigningMethod
	status.Details["expiry_time"] = j.config.ExpiryTime

	return status
}

// GetHealthMetrics returns health metrics for the JWT authenticator.
// zh: GetHealthMetrics 回傳 JWT 認證器的健康指標。
func (j *JWTAuthenticator) GetHealthMetrics() map[string]any {
	return map[string]any{
		"initialized":    j.initialized,
		"issuer":         j.config.Issuer,
		"signing_method": j.config.SigningMethod,
		"expiry_time":    j.config.ExpiryTime,
	}
}

// Authenticate authenticates user credentials and returns user info.
// zh: Authenticate 認證使用者憑證並回傳使用者資訊。
func (j *JWTAuthenticator) Authenticate(ctx context.Context, credentials contracts.Credentials) (*contracts.UserInfo, error) {
	if !j.initialized {
		return nil, fmt.Errorf("JWT authenticator not initialized")
	}

	// TODO: Implement JWT token validation
	if credentials.Type != "jwt" {
		return nil, fmt.Errorf("unsupported credential type: %s", credentials.Type)
	}

	// Mock implementation for now
	return &contracts.UserInfo{
		ID:          "user-123",
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
		Roles:       []string{"user"},
		Permissions: []contracts.Permission{
			{
				Action:   "read",
				Resource: "alerts",
				Scope:    []string{"*"},
			},
		},
	}, nil
}

// ValidateToken validates a JWT token and returns user info.
// zh: ValidateToken 驗證 JWT 令牌並回傳使用者資訊。
func (j *JWTAuthenticator) ValidateToken(ctx context.Context, token string) (*contracts.UserInfo, error) {
	if !j.initialized {
		return nil, fmt.Errorf("JWT authenticator not initialized")
	}

	// TODO: Implement JWT token validation logic
	if token == "" {
		return nil, fmt.Errorf("empty token")
	}

	// Mock implementation for now
	return &contracts.UserInfo{
		ID:          "user-123",
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
		Roles:       []string{"user"},
		Permissions: []contracts.Permission{
			{
				Action:   "read",
				Resource: "alerts",
				Scope:    []string{"*"},
			},
		},
	}, nil
}

// Register registers the JWT authenticator plugin.
// zh: Register 註冊 JWT 認證器插件。
func Register(registry contracts.Registry) error {
	return registry.RegisterPlugin("jwt-authenticator", NewJWTAuthenticator)
}
