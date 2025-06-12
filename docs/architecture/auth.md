# Authenticator 模組架構設計

本文件說明 detectviz 中 `authenticator` 模組的架構設計，參考自 Grafana 的模組化實作風格，以支援多種登入策略並具備高可用性（HA）、可插拔與易測試的特性。

---

## 設計目標

- 支援多種驗證來源：OAuth2、OIDC、LDAP、SAML、靜態帳號、測試帳號
- 插拔式策略實作（Strategy Pattern）
- 可注入 mock/fake 策略以利測試與模擬
- 配合 middleware 使用，支援 JWT 驗證與 Claim 解析
- 高可用性：多策略容錯、註冊多組 Authenticator 並行驗證

---

## 架構概覽

```plaintext
                       +--------------------+
                       |   HTTP Middleware  |
                       +--------+-----------+
                                |
                                v
+--------------------------+    +--------------------------+
| Authenticator Registry   |<-->| []AuthStrategy Interface |
+--------------------------+    +--------------------------+
     |        |       |
     v        v       v
+-------+ +--------+ +--------+
|OAuth  | |Keycloak| |Static  |
+-------+ +--------+ +--------+

     ↑        ↑
     |        |
+------------+--------------+
| configtypes/auth_config.go|
+---------------------------+
```

---

## 相關目錄架構

```bash
internal/adapters/auth/
  ├── oauth.go
  ├── keycloak_adapter.go
	├── nop.go                 # 不執行驗證（for dev/test）
  ├── ldap.go
  ├── static.go # 開發測試用、靜態帳號驗證
  ├── strategies/
		OAuthStrategy
	•	LDAPStrategy
	•	SAMLStrategy
	•	KeycloakStrategy（可改寫 OAuth extend）
  └── mock.go

pkg/ifaces/auth/
  └── auth.go      ← 定義 AuthStrategy, UserInfo 等

internal/registry/auth/
  └── registry.go  ← 註冊所有策略

pkg/configtypes/auth_config.go
```

---

## 模組角色定位

- `Authenticator` 為供 middleware 使用的抽象介面，可驗證並回傳 UserInfo
- `AuthStrategy` 為登入策略模組，實作每種具體驗證邏輯
- `Registry` 為註冊中心，集中註冊所有策略與初始化 Authenticator 集合
- `internal/adapters/auth/strategies/` 為每個策略對應 adapter 的實作位置
- `pkg/configtypes/auth_config.go` 為 config driven 驅動註冊的設定來源
- 所有策略實作皆可由 `plugins/auth/` 擴充，支援熱插拔架構

---

## Interface 定義

定義於 `pkg/ifaces/auth/auth.go`：

```go
type Authenticator interface {
    Authenticate(ctx context.Context, token string) (*UserInfo, error)
    RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
    Name() string
}
```

---

## AuthStrategy Interface 定義（Strategy Definitions）

作為策略抽象層級，`AuthStrategy` 提供每一種登入機制的具體實作。

```go
type AuthStrategy interface {
    Init(config map[string]any) error
    Authenticate(ctx context.Context, token string) (*UserInfo, error)
    RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
    Name() string
}
```

應由註冊中心進行統一包裝轉換為 `Authenticator`。

---

## 策略實作位置

每種策略皆對應一個 Adapter 檔案，放置於 `internal/adapters/auth/`：

- `oauth.go`：通用 OAuth2/OIDC 實作
- `keycloak.go`：特殊處理 Keycloak 權限與 group 欄位
- `ldap.go`：LDAP 認證（未來實作）
- `static.go`：開發環境靜態帳號
- `nop.go`：空策略，永遠允許通過

---

## Plugin 擴充介面

外部可擴充策略模組，需實作以下介面：

```go
type AuthenticatorPlugin interface {
    ID() string
    Strategy() AuthStrategy
}
```

預設註冊於：
- `plugins/auth/`：每個策略模組實作
- `internal/registry/auth/registry.go`：集中註冊並包裝成 Authenticator

用途：
- 可從外部加入自定義認證策略（如：企業 SSO、第三方 Token 驗證等）

---

## 設定與註冊

透過 `pkg/configtypes/auth_config.go` 定義：

```go
type AuthConfig struct {
    Name   string
    Type   string
    Config map[string]any
}
```

註冊處位於 `internal/registry/auth/registry.go`：

- 採用工廠模式註冊每個策略
- 允許一個以上 `Authenticator` 被註冊
- 每次驗證可依名稱選擇或進行 fallback

---

## 多策略抽象（AuthStrategy）

為了支援多種驗證來源與高可用性策略切換，`authenticator` 模組內部採用策略模式（Strategy Pattern），將各種驗證邏輯抽象為 `AuthStrategy`，每一個策略負責單一驗證方法的實作（如 Keycloak、OAuth、Static）。

介面定義：

```go
type AuthStrategy interface {
    Init(config map[string]any) error
    Authenticate(ctx context.Context, token string) (*UserInfo, error)
    RefreshToken(ctx context.Context, refreshToken string) (*TokenPair, error)
    Name() string
}
```

每個 `AuthStrategy` 實作應封裝該驗證來源的初始化、驗證與 token 處理邏輯。所有策略會由註冊中心包裝為 `Authenticator`，供外部調用。

建議實作位置：
- `internal/adapters/auth/strategies/*.go`

註冊於：
- `internal/registry/auth/registry.go`

---

## 驗證流程與 Middleware 注入

完整驗證流程如下：

```text
[HTTP Request]
  └──→ [middleware.Auth()]
         └──→ registry.LookupAuthenticator(name)
                  └──→ strategy.Authenticate(token)
                         → 成功後產生 *UserInfo
                         → 注入 context 中給 downstream 使用
```

- 若驗證失敗會回傳 401 錯誤
- 成功驗證後 context 中包含 `auth.UserKey` 欄位

---

## 高可用設計

- **多策略容錯**：當 A 策略失敗，可 fallback 至 B
- **模擬策略**：開發或測試時可載入 `static.go`, `fake.go`
- **註冊多組策略**：支援同時存在 Google + Keycloak + Static
- **適用於分散式環境**：Middleware 可根據 context 使用 local JWT 驗證而非每次查詢後端

- 可設定驗證順序，預設依註冊順序進行
- 亦可根據 Token prefix 或 Header 條件自動切換策略（如 `Bearer keycloak:xxx` → Keycloak）

---

## 使用範例

1. HTTP Middleware 驗證 token 並回傳 `UserInfo`
2. API 呼叫驗證 Access Token 並根據 Role 權限篩選回應
3. 支援 `Authorization: Bearer ...` 與 Cookie-based session
4. 前端透過 `/me` endpoint 取得登入資訊

---

## 測試建議

- 每一策略對應：
  - `*_test.go`
  - `mock.go` / `fake.go` / `fallback.go`
- 覆蓋：
  - 正常驗證、token 過期、claim 錯誤、無效格式等情境

---

## 與其他模組整合關係

| 模組               | 關聯說明                                               |
|--------------------|--------------------------------------------------------|
| `internal/middleware/` | middleware auth 驗證依賴 `Authenticator` 並注入 `UserInfo` |
| `internal/registry/auth/` | 託管所有註冊策略與預設策略注入邏輯                         |
| `internal/api/`     | 提供 `/me`, `/login`, `/logout` 等 API endpoint 使用驗證結果 |
| `pkg/ifaces/auth/`  | 提供 interface 對外公開，供 middleware / CLI / service 呼叫 |
| `plugins/auth/`     | 可選擇性新增策略模組，支援外部擴充，例如 Google、LDAP、Mock |

---

## 後續擴充方向

- 實作 LDAP 與 SAML 支援
- 支援 Token Revocation 機制
- 引入 Session Store 與 Token Caching

---

## 文件狀態

本文件已對齊 detectviz 全域設計標準，下一階段預計實作以下：

- [ ] `auth.go` middleware 驗證器
- [ ] `keycloak.go` 與 `oauth.go` 策略
- [ ] `registry.go` 註冊流程與 fallback
- [ ] `/me` endpoint handler
- [ ] 測試模組：mock/fake 實作與驗證

參考文件：
- [docs/interfaces/middleware.md](../interfaces/middleware.md)
- [docs/web-architecture.md](../web-architecture.md)
- [docs/module-architecture-overview.md](../module-architecture-overview.md)
