package contracts

import (
	"context"
	"time"
)

// Authenticator defines the interface for authentication providers.
// zh: Authenticator 定義認證提供者介面。
type Authenticator interface {
	Plugin
	Authenticate(ctx context.Context, credentials Credentials) (*UserInfo, error)
	ValidateToken(ctx context.Context, token string) (*UserInfo, error)
}

// Authorizer defines the interface for authorization providers.
// zh: Authorizer 定義授權提供者介面。
type Authorizer interface {
	Plugin
	Authorize(ctx context.Context, user *UserInfo, action, resource string) (bool, error)
	GetPermissions(ctx context.Context, user *UserInfo) ([]Permission, error)
}

// TokenManager defines the interface for token management.
// zh: TokenManager 定義令牌管理介面。
type TokenManager interface {
	Plugin
	GenerateToken(ctx context.Context, user *UserInfo) (*Token, error)
	RefreshToken(ctx context.Context, refreshToken string) (*Token, error)
	RevokeToken(ctx context.Context, token string) error
}

// Credentials represents authentication credentials.
// zh: Credentials 代表認證憑證。
type Credentials struct {
	Type     string            `json:"type"` // basic, jwt, oauth, etc.
	Username string            `json:"username,omitempty"`
	Password string            `json:"password,omitempty"`
	Token    string            `json:"token,omitempty"`
	Extra    map[string]string `json:"extra,omitempty"`
}

// UserInfo represents authenticated user information.
// zh: UserInfo 代表已認證的使用者資訊。
type UserInfo struct {
	ID          string            `json:"id"`
	Username    string            `json:"username"`
	Email       string            `json:"email,omitempty"`
	DisplayName string            `json:"display_name,omitempty"`
	Roles       []string          `json:"roles"`
	Permissions []Permission      `json:"permissions"`
	Attributes  map[string]string `json:"attributes,omitempty"`
	ExpiresAt   *time.Time        `json:"expires_at,omitempty"`
}

// Permission represents a user permission.
// zh: Permission 代表使用者權限。
type Permission struct {
	Action   string   `json:"action"`   // read, write, execute, admin
	Resource string   `json:"resource"` // alerts, rules, users, etc.
	Scope    []string `json:"scope"`    // specific resource IDs or patterns
}

// Token represents an authentication token.
// zh: Token 代表認證令牌。
type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// AuthConfig defines the configuration structure for authentication.
// zh: AuthConfig 定義認證的配置結構。
type AuthConfig struct {
	Provider     string         `yaml:"provider" json:"provider"`
	Endpoint     string         `yaml:"endpoint" json:"endpoint"`
	ClientID     string         `yaml:"client_id" json:"client_id"`
	ClientSecret string         `yaml:"client_secret" json:"client_secret"`
	RedirectURL  string         `yaml:"redirect_url" json:"redirect_url"`
	Scopes       []string       `yaml:"scopes" json:"scopes"`
	ExtraConfig  map[string]any `yaml:"extra_config" json:"extra_config"`
}

// AuthRegistry defines the interface for authentication provider registration.
// zh: AuthRegistry 定義認證提供者註冊介面。
type AuthRegistry interface {
	RegisterAuthenticator(name string, factory func(config AuthConfig) (Authenticator, error)) error
	RegisterAuthorizer(name string, factory func(config AuthConfig) (Authorizer, error)) error
	RegisterTokenManager(name string, factory func(config AuthConfig) (TokenManager, error)) error
	GetAuthenticator(name string) (Authenticator, error)
	GetAuthorizer(name string) (Authorizer, error)
	GetTokenManager(name string) (TokenManager, error)
}
