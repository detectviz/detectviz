package keycloak

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"detectviz/pkg/platform/contracts"
)

// KeycloakPlugin provides Keycloak-based authentication.
// zh: KeycloakPlugin 提供基於 Keycloak 的認證。
type KeycloakPlugin struct {
	name        string
	version     string
	description string
	config      *Config
	httpClient  *http.Client
	initialized bool
	started     bool
}

// Config contains configuration for the Keycloak plugin.
// zh: Config 包含 Keycloak 插件的配置。
type Config struct {
	ServerURL    string            `yaml:"server_url" json:"server_url"`
	Realm        string            `yaml:"realm" json:"realm"`
	ClientID     string            `yaml:"client_id" json:"client_id"`
	ClientSecret string            `yaml:"client_secret" json:"client_secret"`
	RedirectURI  string            `yaml:"redirect_uri" json:"redirect_uri"`
	Scopes       []string          `yaml:"scopes" json:"scopes"`
	Timeout      int               `yaml:"timeout" json:"timeout"` // seconds
	CertPath     string            `yaml:"cert_path" json:"cert_path"`
	InsecureTLS  bool              `yaml:"insecure_tls" json:"insecure_tls"`
	UserMapping  map[string]string `yaml:"user_mapping" json:"user_mapping"`
	RoleMapping  map[string]string `yaml:"role_mapping" json:"role_mapping"`
}

// KeycloakUserInfo represents user information from Keycloak.
// zh: KeycloakUserInfo 代表來自 Keycloak 的使用者資訊。
type KeycloakUserInfo struct {
	Sub               string                 `json:"sub"`
	Name              string                 `json:"name"`
	PreferredUsername string                 `json:"preferred_username"`
	GivenName         string                 `json:"given_name"`
	FamilyName        string                 `json:"family_name"`
	Email             string                 `json:"email"`
	EmailVerified     bool                   `json:"email_verified"`
	Groups            []string               `json:"groups"`
	RealmAccess       map[string]interface{} `json:"realm_access"`
	ResourceAccess    map[string]interface{} `json:"resource_access"`
}

// KeycloakTokenResponse represents token response from Keycloak.
// zh: KeycloakTokenResponse 代表來自 Keycloak 的令牌回應。
type KeycloakTokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	Scope            string `json:"scope"`
}

// NewKeycloakPlugin creates a new Keycloak plugin instance.
// zh: NewKeycloakPlugin 建立新的 Keycloak 插件實例。
func NewKeycloakPlugin(config any) (contracts.Plugin, error) {
	var cfg Config

	// Set default configuration
	cfg = Config{
		Timeout:     30,
		Scopes:      []string{"openid", "profile", "email"},
		InsecureTLS: false,
		UserMapping: map[string]string{
			"id":       "sub",
			"username": "preferred_username",
			"email":    "email",
			"name":     "name",
		},
		RoleMapping: map[string]string{
			"admin":     "realm-admin",
			"user":      "realm-user",
			"moderator": "realm-moderator",
		},
	}

	// Parse configuration if provided using manual parsing
	if config != nil {
		if err := parseKeycloakConfig(config, &cfg); err != nil {
			return nil, fmt.Errorf("failed to decode Keycloak config: %w", err)
		}
	}

	// Validate required configuration
	if cfg.ServerURL == "" {
		return nil, fmt.Errorf("server_url is required")
	}
	if cfg.Realm == "" {
		return nil, fmt.Errorf("realm is required")
	}
	if cfg.ClientID == "" {
		return nil, fmt.Errorf("client_id is required")
	}

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: time.Duration(cfg.Timeout) * time.Second,
	}

	return &KeycloakPlugin{
		name:        "keycloak-authenticator",
		version:     "1.0.0",
		description: "Keycloak integration for SSO authentication and authorization",
		config:      &cfg,
		httpClient:  httpClient,
		initialized: false,
		started:     false,
	}, nil
}

// parseKeycloakConfig manually parses configuration from map[string]any
// zh: parseKeycloakConfig 手動從 map[string]any 解析配置
func parseKeycloakConfig(config any, target *Config) error {
	if config == nil {
		return nil
	}

	configMap, ok := config.(map[string]any)
	if !ok {
		return fmt.Errorf("config must be a map[string]any")
	}

	if serverURL, exists := configMap["server_url"]; exists {
		if str, ok := serverURL.(string); ok {
			target.ServerURL = str
		}
	}
	if realm, exists := configMap["realm"]; exists {
		if str, ok := realm.(string); ok {
			target.Realm = str
		}
	}
	if clientID, exists := configMap["client_id"]; exists {
		if str, ok := clientID.(string); ok {
			target.ClientID = str
		}
	}
	if clientSecret, exists := configMap["client_secret"]; exists {
		if str, ok := clientSecret.(string); ok {
			target.ClientSecret = str
		}
	}
	if redirectURI, exists := configMap["redirect_uri"]; exists {
		if str, ok := redirectURI.(string); ok {
			target.RedirectURI = str
		}
	}
	if timeout, exists := configMap["timeout"]; exists {
		if intVal, ok := timeout.(int); ok {
			target.Timeout = intVal
		}
	}
	if scopes, exists := configMap["scopes"]; exists {
		if scopeList, ok := scopes.([]interface{}); ok {
			strScopes := make([]string, len(scopeList))
			for i, scope := range scopeList {
				if str, ok := scope.(string); ok {
					strScopes[i] = str
				}
			}
			target.Scopes = strScopes
		}
	}

	return nil
}

// Plugin interface implementation
// zh: Plugin 介面實作

// Name returns the plugin name.
// zh: Name 回傳插件名稱。
func (p *KeycloakPlugin) Name() string {
	return p.name
}

// Version returns the plugin version.
// zh: Version 回傳插件版本。
func (p *KeycloakPlugin) Version() string {
	return p.version
}

// Description returns the plugin description.
// zh: Description 回傳插件描述。
func (p *KeycloakPlugin) Description() string {
	return p.description
}

// Init initializes the plugin.
// zh: Init 初始化插件。
func (p *KeycloakPlugin) Init(config any) error {
	if p.initialized {
		return nil
	}

	// Test connection to validate configuration
	if err := p.testConnection(); err != nil {
		return fmt.Errorf("failed to validate Keycloak connection: %w", err)
	}

	p.initialized = true
	return nil
}

// Shutdown shuts down the plugin.
// zh: Shutdown 關閉插件。
func (p *KeycloakPlugin) Shutdown() error {
	p.started = false
	p.initialized = false
	return nil
}

// LifecycleAware interface implementation
// zh: LifecycleAware 介面實作

// OnRegister is called when the plugin is registered.
// zh: OnRegister 在插件註冊時被呼叫。
func (p *KeycloakPlugin) OnRegister() error {
	fmt.Printf("Keycloak plugin registered for realm: %s\n", p.config.Realm)
	return nil
}

// OnStart is called when the plugin is started.
// zh: OnStart 在插件啟動時被呼叫。
func (p *KeycloakPlugin) OnStart() error {
	if !p.initialized {
		return fmt.Errorf("plugin not initialized")
	}

	// Test connection to Keycloak
	if err := p.testConnection(); err != nil {
		return fmt.Errorf("failed to connect to Keycloak: %w", err)
	}

	p.started = true
	fmt.Printf("Keycloak plugin started successfully\n")
	return nil
}

// OnStop is called when the plugin is stopped.
// zh: OnStop 在插件停止時被呼叫。
func (p *KeycloakPlugin) OnStop() error {
	p.started = false
	fmt.Printf("Keycloak plugin stopped\n")
	return nil
}

// OnShutdown is called when the plugin is shut down.
// zh: OnShutdown 在插件關閉時被呼叫。
func (p *KeycloakPlugin) OnShutdown() error {
	return p.Shutdown()
}

// HealthChecker interface implementation
// zh: HealthChecker 介面實作

// CheckHealth checks the health of the Keycloak connection.
// zh: CheckHealth 檢查 Keycloak 連接的健康狀態。
func (p *KeycloakPlugin) CheckHealth(ctx context.Context) contracts.HealthStatus {
	status := contracts.HealthStatus{
		Timestamp: time.Now(),
		Details:   make(map[string]any),
	}

	if !p.initialized {
		status.Status = "unhealthy"
		status.Message = "Plugin not initialized"
		return status
	}

	if !p.started {
		status.Status = "unhealthy"
		status.Message = "Plugin not started"
		return status
	}

	// Test connection to Keycloak
	if err := p.testConnection(); err != nil {
		status.Status = "unhealthy"
		status.Message = fmt.Sprintf("Connection failed: %v", err)
		return status
	}

	status.Status = "healthy"
	status.Message = "Connected to Keycloak successfully"
	status.Details["server_url"] = p.config.ServerURL
	status.Details["realm"] = p.config.Realm
	status.Details["client_id"] = p.config.ClientID

	return status
}

// GetHealthMetrics returns health metrics for the plugin.
// zh: GetHealthMetrics 回傳插件的健康指標。
func (p *KeycloakPlugin) GetHealthMetrics() map[string]any {
	return map[string]any{
		"initialized": p.initialized,
		"started":     p.started,
		"server_url":  p.config.ServerURL,
		"realm":       p.config.Realm,
		"timeout":     p.config.Timeout,
	}
}

// Authenticator interface implementation
// zh: Authenticator 介面實作

// Authenticate authenticates a user using Keycloak.
// zh: Authenticate 使用 Keycloak 認證使用者。
func (p *KeycloakPlugin) Authenticate(ctx context.Context, credentials contracts.Credentials) (*contracts.UserInfo, error) {
	if !p.started {
		return nil, fmt.Errorf("plugin not started")
	}

	switch credentials.Type {
	case "password":
		return p.authenticateWithPassword(ctx, credentials.Username, credentials.Password)
	case "token":
		return p.authenticateWithToken(ctx, credentials.Token)
	case "oauth_code":
		if code, ok := credentials.Extra["code"]; ok {
			state := credentials.Extra["state"]
			return p.authenticateWithAuthCode(ctx, code, state)
		}
		return nil, fmt.Errorf("missing authorization code")
	default:
		return nil, fmt.Errorf("unsupported credential type: %s", credentials.Type)
	}
}

// ValidateToken validates an access token.
// zh: ValidateToken 驗證存取令牌。
func (p *KeycloakPlugin) ValidateToken(ctx context.Context, tokenString string) (*contracts.UserInfo, error) {
	if !p.started {
		return nil, fmt.Errorf("plugin not started")
	}

	// Validate token by calling Keycloak userinfo endpoint
	userInfoURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo",
		strings.TrimRight(p.config.ServerURL, "/"), p.config.Realm)

	req, err := http.NewRequestWithContext(ctx, "GET", userInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create userinfo request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+tokenString)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("userinfo request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token validation failed with status: %d", resp.StatusCode)
	}

	var keycloakUser KeycloakUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&keycloakUser); err != nil {
		return nil, fmt.Errorf("failed to decode userinfo response: %w", err)
	}

	// Convert to UserInfo
	return p.keycloakUserToUserInfo(&keycloakUser), nil
}

// Helper methods
// zh: 輔助方法

// testConnection tests the connection to Keycloak.
// zh: testConnection 測試與 Keycloak 的連接。
func (p *KeycloakPlugin) testConnection() error {
	realmURL := fmt.Sprintf("%s/realms/%s",
		strings.TrimRight(p.config.ServerURL, "/"), p.config.Realm)

	resp, err := p.httpClient.Get(realmURL)
	if err != nil {
		return fmt.Errorf("connection test failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("connection test failed with status: %d", resp.StatusCode)
	}

	return nil
}

// authenticateWithPassword authenticates using username/password.
// zh: authenticateWithPassword 使用用戶名/密碼進行認證。
func (p *KeycloakPlugin) authenticateWithPassword(ctx context.Context, username, password string) (*contracts.UserInfo, error) {
	// Prepare token request
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", username)
	data.Set("password", password)
	data.Set("client_id", p.config.ClientID)
	if p.config.ClientSecret != "" {
		data.Set("client_secret", p.config.ClientSecret)
	}
	data.Set("scope", strings.Join(p.config.Scopes, " "))

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		strings.TrimRight(p.config.ServerURL, "/"), p.config.Realm)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("authentication request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("authentication failed with status: %d", resp.StatusCode)
	}

	var tokenResp KeycloakTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Validate and extract user info from access token
	return p.ValidateToken(ctx, tokenResp.AccessToken)
}

// authenticateWithToken authenticates using an access token.
// zh: authenticateWithToken 使用存取令牌進行認證。
func (p *KeycloakPlugin) authenticateWithToken(ctx context.Context, token string) (*contracts.UserInfo, error) {
	return p.ValidateToken(ctx, token)
}

// authenticateWithAuthCode authenticates using authorization code.
// zh: authenticateWithAuthCode 使用授權碼進行認證。
func (p *KeycloakPlugin) authenticateWithAuthCode(ctx context.Context, code, state string) (*contracts.UserInfo, error) {
	// Exchange authorization code for tokens
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", p.config.RedirectURI)
	data.Set("client_id", p.config.ClientID)
	if p.config.ClientSecret != "" {
		data.Set("client_secret", p.config.ClientSecret)
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		strings.TrimRight(p.config.ServerURL, "/"), p.config.Realm)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status: %d", resp.StatusCode)
	}

	var tokenResp KeycloakTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	// Validate and extract user info from access token
	return p.ValidateToken(ctx, tokenResp.AccessToken)
}

// keycloakUserToUserInfo converts Keycloak user info to UserInfo.
// zh: keycloakUserToUserInfo 將 Keycloak 使用者資訊轉換為 UserInfo。
func (p *KeycloakPlugin) keycloakUserToUserInfo(keycloakUser *KeycloakUserInfo) *contracts.UserInfo {
	userInfo := &contracts.UserInfo{
		ID:          keycloakUser.Sub,
		Username:    keycloakUser.PreferredUsername,
		Email:       keycloakUser.Email,
		DisplayName: keycloakUser.Name,
		Attributes:  make(map[string]string),
	}

	// Map custom claims using user mapping
	for userField, keycloakField := range p.config.UserMapping {
		switch keycloakField {
		case "sub":
			userInfo.Attributes[userField] = keycloakUser.Sub
		case "preferred_username":
			userInfo.Attributes[userField] = keycloakUser.PreferredUsername
		case "email":
			userInfo.Attributes[userField] = keycloakUser.Email
		case "name":
			userInfo.Attributes[userField] = keycloakUser.Name
		case "given_name":
			userInfo.Attributes[userField] = keycloakUser.GivenName
		case "family_name":
			userInfo.Attributes[userField] = keycloakUser.FamilyName
		}
	}

	// Extract roles from realm access
	if realmAccess, ok := keycloakUser.RealmAccess["roles"].([]interface{}); ok {
		for _, role := range realmAccess {
			if roleStr, ok := role.(string); ok {
				userInfo.Roles = append(userInfo.Roles, roleStr)

				// Map roles to permissions using role mapping
				if permission, exists := p.config.RoleMapping[roleStr]; exists {
					userInfo.Permissions = append(userInfo.Permissions, contracts.Permission{
						Action:   "allow",
						Resource: permission,
						Scope:    []string{"*"},
					})
				}
			}
		}
	}

	return userInfo
}

// GetAuthorizationURL generates an authorization URL for OAuth flow.
// zh: GetAuthorizationURL 為 OAuth 流程生成授權 URL。
func (p *KeycloakPlugin) GetAuthorizationURL(state string) string {
	authURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth",
		strings.TrimRight(p.config.ServerURL, "/"), p.config.Realm)

	params := url.Values{}
	params.Set("client_id", p.config.ClientID)
	params.Set("redirect_uri", p.config.RedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", strings.Join(p.config.Scopes, " "))
	params.Set("state", state)

	return fmt.Sprintf("%s?%s", authURL, params.Encode())
}

// RefreshToken refreshes an access token using refresh token (for compatibility).
// zh: RefreshToken 使用刷新令牌刷新存取令牌（為了相容性）。
func (p *KeycloakPlugin) RefreshToken(ctx context.Context, refreshToken string) (*contracts.Token, error) {
	if !p.started {
		return nil, fmt.Errorf("plugin not started")
	}

	// Prepare refresh token request
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", p.config.ClientID)
	if p.config.ClientSecret != "" {
		data.Set("client_secret", p.config.ClientSecret)
	}

	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token",
		strings.TrimRight(p.config.ServerURL, "/"), p.config.Realm)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("refresh token request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("refresh token failed with status: %d", resp.StatusCode)
	}

	var tokenResp KeycloakTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	expiresAt := time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return &contracts.Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		TokenType:    tokenResp.TokenType,
		ExpiresAt:    expiresAt,
	}, nil
}

// Logout logs out a user (for compatibility).
// zh: Logout 使用者登出（為了相容性）。
func (p *KeycloakPlugin) Logout(ctx context.Context, token string) error {
	if !p.started {
		return fmt.Errorf("plugin not started")
	}

	// Prepare logout request
	data := url.Values{}
	data.Set("client_id", p.config.ClientID)
	if p.config.ClientSecret != "" {
		data.Set("client_secret", p.config.ClientSecret)
	}
	data.Set("refresh_token", token)

	logoutURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/logout",
		strings.TrimRight(p.config.ServerURL, "/"), p.config.Realm)

	req, err := http.NewRequestWithContext(ctx, "POST", logoutURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create logout request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("logout request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("logout failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Register registers the Keycloak plugin.
// zh: Register 註冊 Keycloak 插件。
func Register(registry contracts.Registry) error {
	return registry.RegisterPlugin("keycloak-authenticator", NewKeycloakPlugin)
}
