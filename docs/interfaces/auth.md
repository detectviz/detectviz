# Authenticator Interface

本文件說明 detectviz 專案中 `auth` 模組的介面設計，目標為支援多種 SSO/OAuth 驗證來源，如 Keycloak、LDAP、Google、GitHub 等，具備可擴充、可測試、可插拔的設計風格。

---

## 目標

- 抽象出統一的 `Authenticator` 驗證介面
- 支援多種登入策略與第三方驗證服務（SSO/OAuth2/SAML）
- 提供使用者資訊、Token 解析、Refresh 等功能
- 支援測試與模擬驗證邏輯（mock/fake）

---

## Interface 定義

介面定義位於 `pkg/ifaces/auth/auth.go`：

```go
type Authenticator interface {
    Authenticate(ctx context.Context, token string) (*UserInfo, error)
    RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
    Name() string
}
```

### UserInfo

```go
type UserInfo struct {
    Subject string
    Name    string
    Email   string
    Roles   []string
}
```

### TokenPair

```go
type TokenPair struct {
    AccessToken  string
    RefreshToken string
    Expiry       time.Time
}
```

---

## 策略實作建議

策略實作應放置於 `internal/adapters/auth/`，例如：

- `keycloak.go`
- `oauth.go`
- `ldap.go`
- `static.go`（靜態帳號供開發測試使用）

---

## 設定與註冊

建議透過以下方式註冊：

- 設定結構：`pkg/configtypes/auth_config.go`
- 註冊模組：`internal/registry/auth/registry.go`
- 多策略支援：可定義為 `[]AuthConfig` 以支援多組 Authenticator 同時啟用

---

## 典型使用情境

- HTTP 請求的 Middleware 驗證 token
- WebSocket 或 gRPC 的使用者身份驗證
- API 根據 token 回傳使用者資訊與權限範圍
- 授權判斷中介層依據 `UserInfo.Roles` 決定可見範圍

---

## 測試建議

請為每一個策略實作測試模組：

- `*_test.go` 基本驗證流程
- `mock.go`, `fake.go`, `fallback.go` 供單元測試與整合測試使用
- 驗證 edge case，例如失效 token、非法 claim、過期 refresh token

---