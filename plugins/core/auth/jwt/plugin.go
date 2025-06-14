package jwt

import (
	"context"
	"fmt"
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
}

// JWTConfig defines the configuration for JWT authentication.
// zh: JWTConfig 定義 JWT 認證的配置。
type JWTConfig struct {
	SecretKey     string        `yaml:"secret_key" json:"secret_key"`
	Issuer        string        `yaml:"issuer" json:"issuer"`
	ExpiryTime    time.Duration `yaml:"expiry_time" json:"expiry_time"`
	RefreshTime   time.Duration `yaml:"refresh_time" json:"refresh_time"`
	SigningMethod string        `yaml:"signing_method" json:"signing_method"`
}

// NewJWTAuthenticator creates a new JWT authenticator instance.
// zh: NewJWTAuthenticator 建立新的 JWT 認證器實例。
func NewJWTAuthenticator(config any) (contracts.Plugin, error) {
	jwtConfig := &JWTConfig{
		SecretKey:     "default-secret-key",
		Issuer:        "detectviz",
		ExpiryTime:    time.Hour * 24,
		RefreshTime:   time.Hour * 24 * 7,
		SigningMethod: "HS256",
	}

	// TODO: Parse actual config from the provided config parameter
	if config != nil {
		// Parse config here
	}

	return &JWTAuthenticator{
		name:        "jwt-authenticator",
		version:     "1.0.0",
		description: "JWT-based authentication plugin",
		config:      jwtConfig,
	}, nil
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
	// TODO: Initialize JWT signing/verification components
	return nil
}

// Shutdown shuts down the JWT authenticator.
// zh: Shutdown 關閉 JWT 認證器。
func (j *JWTAuthenticator) Shutdown() error {
	// TODO: Cleanup resources
	return nil
}

// Authenticate authenticates user credentials and returns user info.
// zh: Authenticate 認證使用者憑證並回傳使用者資訊。
func (j *JWTAuthenticator) Authenticate(ctx context.Context, credentials contracts.Credentials) (*contracts.UserInfo, error) {
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
