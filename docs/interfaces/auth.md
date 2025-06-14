# Auth Interface 認證介面

> **檔案位置**: `pkg/platform/contracts/auth.go`

## 概述

`Authenticator` 介面定義了使用者認證和授權功能。支援多種認證方式，包括 JWT、OAuth、LDAP、SAML 等，提供統一的認證抽象層。

## 介面定義

```go
type Authenticator interface {
    Plugin
    Authenticate(ctx context.Context, credentials Credentials) (*UserInfo, error)
    ValidateToken(ctx context.Context, token string) (*UserInfo, error)
    RefreshToken(ctx context.Context, refreshToken string) (*TokenInfo, error)
    Logout(ctx context.Context, token string) error
}
```

## 方法說明

### Authenticate(ctx context.Context, credentials Credentials) (*UserInfo, error)
- **用途**: 使用憑證進行使用者認證
- **參數**:
  - `ctx` - 上下文，用於取消和超時控制
  - `credentials` - 使用者憑證（用戶名密碼、JWT token 等）
- **回傳值**: 認證成功回傳使用者資訊，失敗回傳錯誤
- **使用場景**: 登入流程

### ValidateToken(ctx context.Context, token string) (*UserInfo, error)
- **用途**: 驗證存取令牌的有效性
- **參數**:
  - `ctx` - 上下文
  - `token` - 要驗證的令牌
- **回傳值**: 令牌有效回傳使用者資訊，無效回傳錯誤
- **使用場景**: API 請求認證中介層

### RefreshToken(ctx context.Context, refreshToken string) (*TokenInfo, error)
- **用途**: 使用刷新令牌獲取新的存取令牌
- **參數**:
  - `ctx` - 上下文
  - `refreshToken` - 刷新令牌
- **回傳值**: 新的令牌資訊
- **使用場景**: 令牌續約

### Logout(ctx context.Context, token string) error
- **用途**: 登出並撤銷令牌
- **參數**:
  - `ctx` - 上下文
  - `token` - 要撤銷的令牌
- **回傳值**: 登出失敗時回傳錯誤
- **使用場景**: 安全登出

## 資料結構

```go
type Credentials struct {
    Type     string                 `json:"type"`     // password, token, certificate
    Username string                 `json:"username,omitempty"`
    Password string                 `json:"password,omitempty"`
    Token    string                 `json:"token,omitempty"`
    Extra    map[string]interface{} `json:"extra,omitempty"`
}

type UserInfo struct {
    ID          string            `json:"id"`
    Username    string            `json:"username"`
    Email       string            `json:"email,omitempty"`
    DisplayName string            `json:"display_name,omitempty"`
    Roles       []string          `json:"roles"`
    Permissions []string          `json:"permissions"`
    Groups      []string          `json:"groups,omitempty"`
    Attributes  map[string]string `json:"attributes,omitempty"`
    ExpiresAt   *time.Time        `json:"expires_at,omitempty"`
}

type TokenInfo struct {
    AccessToken  string     `json:"access_token"`
    RefreshToken string     `json:"refresh_token,omitempty"`
    TokenType    string     `json:"token_type"`
    ExpiresIn    int64      `json:"expires_in"`
    ExpiresAt    time.Time  `json:"expires_at"`
    Scope        []string   `json:"scope,omitempty"`
}
```

## 配置結構

```go
type AuthConfig struct {
    Enabled        bool              `yaml:"enabled" json:"enabled"`
    Provider       string            `yaml:"provider" json:"provider"`         // jwt, oauth, ldap, saml
    TokenLifetime  string            `yaml:"token_lifetime" json:"token_lifetime"`
    RefreshEnabled bool              `yaml:"refresh_enabled" json:"refresh_enabled"`
    SessionTimeout string            `yaml:"session_timeout" json:"session_timeout"`
    Config         map[string]any    `yaml:"config" json:"config"`             // Provider 特定配置
}
```

## 實作範例

```go
type JWTAuthenticator struct {
    name        string
    version     string
    description string
    config      *JWTConfig
    secretKey   []byte
    initialized bool
}

func (j *JWTAuthenticator) Authenticate(ctx context.Context, credentials Credentials) (*UserInfo, error) {
    if credentials.Type != "password" {
        return nil, fmt.Errorf("unsupported credential type: %s", credentials.Type)
    }
    
    // 驗證用戶名密碼（這裡可能需要查詢資料庫或 LDAP）
    if !j.validateUser(credentials.Username, credentials.Password) {
        return nil, fmt.Errorf("invalid credentials")
    }
    
    // 創建 JWT token
    token, err := j.createJWT(credentials.Username)
    if err != nil {
        return nil, fmt.Errorf("failed to create token: %w", err)
    }
    
    return &UserInfo{
        ID:       credentials.Username,
        Username: credentials.Username,
        Roles:    []string{"user"},
        // ... 其他資訊
    }, nil
}

func (j *JWTAuthenticator) ValidateToken(ctx context.Context, token string) (*UserInfo, error) {
    claims, err := j.parseJWT(token)
    if err != nil {
        return nil, fmt.Errorf("invalid token: %w", err)
    }
    
    return &UserInfo{
        ID:       claims.Subject,
        Username: claims.Subject,
        // ... 從 claims 提取的資訊
    }, nil
}
```

## 內建認證器類型

### JWT 認證
- **本地 JWT**: 使用對稱密鑰簽名
- **RS256 JWT**: 使用 RSA 非對稱密鑰
- **JWKS**: 支援 JSON Web Key Set

### OAuth 2.0 / OpenID Connect
- **Authorization Code**: 標準 OAuth 2.0 流程
- **Client Credentials**: 服務間認證
- **Implicit**: 單頁應用程式
- **PKCE**: 公開客戶端安全擴展

### 企業整合
- **LDAP**: 輕量級目錄存取協議
- **Active Directory**: Microsoft AD 整合
- **SAML**: 安全斷言標記語言
- **Kerberos**: 網路認證協議

### 多因子認證
- **TOTP**: 基於時間的一次性密碼
- **SMS**: 簡訊驗證碼
- **Email**: 電子郵件驗證碼
- **Hardware Token**: 硬體令牌

## 安全最佳實務

1. **密碼安全**: 
   - 使用強密碼策略
   - 密碼雜湊（bcrypt, scrypt, Argon2）
   - 鹽值防止彩虹表攻擊

2. **令牌管理**:
   - 短暫的存取令牌生命週期
   - 安全的刷新令牌存儲
   - 令牌撤銷機制

3. **會話安全**:
   - 安全的會話存儲
   - 會話超時機制
   - 會話固定攻擊防護

4. **傳輸安全**:
   - HTTPS 加密傳輸
   - HTTP 安全標頭
   - CSRF 防護

## 配置範例

### JWT 認證配置
```yaml
plugins:
  - name: jwt-authenticator
    type: auth
    enabled: true
    config:
      secret_key: "your-secret-key"
      issuer: "detectviz"
      expiry_time: "24h"
      signing_method: "HS256"
```

### OAuth 2.0 配置
```yaml
plugins:
  - name: oauth-authenticator
    type: auth
    enabled: true
    config:
      provider: "google"
      client_id: "your-client-id"
      client_secret: "your-client-secret"
      redirect_uri: "http://localhost:8080/auth/callback"
      scopes: ["openid", "profile", "email"]
```

### LDAP 配置
```yaml
plugins:
  - name: ldap-authenticator
    type: auth
    enabled: true
    config:
      server: "ldap://localhost:389"
      bind_dn: "cn=admin,dc=example,dc=com"
      bind_password: "admin-password"
      user_base_dn: "ou=users,dc=example,dc=com"
      user_filter: "(uid=%s)"
```

## 中介層整合

```go
func AuthMiddleware(auth contracts.Authenticator) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            token := extractToken(c.Request())
            if token == "" {
                return c.JSON(401, map[string]string{"error": "missing token"})
            }
            
            userInfo, err := auth.ValidateToken(c.Request().Context(), token)
            if err != nil {
                return c.JSON(401, map[string]string{"error": "invalid token"})
            }
            
            c.Set("user", userInfo)
            return next(c)
        }
    }
}
```

## 技術棧要求

1. **加密安全**: 使用經過驗證的加密函式庫
2. **上下文管理**: 支援 `context.Context` 超時和取消
3. **錯誤處理**: 不洩漏敏感資訊的錯誤訊息
4. **日誌記錄**: 記錄認證事件和安全事件
5. **配置管理**: 支援安全的配置管理

## 相關文件

- [Plugin Interface](./plugin.md)
- [Lifecycle Interface](./lifecycle.md)
- [Security Best Practices](../security/auth-guide.md)
- [Plugin Development Guide](../develop-guide.md) 