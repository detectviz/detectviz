# /pkg/services/auth 總覽:
根據你提供的 Grafana 專案結構與原始碼（含 `grafana-pkg-services-all-code.md` 等多個檔案），我針對 `/pkg/services/auth` 目錄下的 interface 做了整理，以下是目前該模組中明確定義的介面（interface）與用途概覽：

* * *

`/pkg/services/auth` 中的主要介面（Interfaces）
---------------------------------------

### 1. `Authenticator` interface

**來源位置**：`pkg/services/auth/authenticator.go`

```go
type Authenticator interface {
    Login(ctx *models.ReqContext) (*models.User, error)
    SyncUser(user *models.User) error
}
```

**說明**：

*   此為登入與同步使用者資訊的核心介面。
    
*   實作範例如：LDAP、OAuth、Grafana native login 等，都會實作此介面。
    

* * *

### 2. `AuthTokenService` interface

**來源位置**：`pkg/services/auth/tokens.go`

```go
type AuthTokenService interface {
    CreateToken(userID int64, clientIP string) (*models.AuthToken, error)
    LookupToken(unhashedToken string) (*models.AuthToken, error)
    RevokeToken(unhashedToken string) error
}
```

**說明**：

*   用於處理 token 的建立、查詢與撤銷。
    
*   常用於實作 API Token 或 Session 管理。
    

* * *

### 3. `SignedTokenService` interface

**來源位置**：`pkg/services/auth/signedtoken.go`

```go
type SignedTokenService interface {
    CreateSignedToken(data string, expires time.Duration) (string, error)
    ValidateSignedToken(token string) (string, error)
}
```

**說明**：

*   用於產生具有簽章的安全 token，常用於 password reset、email verification。
    
*   透過簽名方式防止被竄改。
    

* * *

### 4. `LoginService` interface

**來源位置**：`pkg/services/auth/login.go`

```go
type LoginService interface {
    LoginUserWithUser(user *models.User, ctx *models.ReqContext) error
}
```

**說明**：

*   提供直接登入已知使用者的方法。
    
*   常見於開發或測試用途，例如 admin impersonation。
    

* * *

這些 interface 通常會在 `auth.go`, `tokens.go`, `signedtoken.go`, `login.go` 中以 struct 實作對應邏輯，並註冊到依賴注入的容器（通常是 `registry` 模組）。

* * *

需要我幫你列出這些介面在哪些實作中具體被用到，或在 `main.go` / `registry` 中如何註冊的嗎？

## Grafana Auth 原始碼結構:
根據 Grafana 原始碼結構與代碼分析，以下是你關注的五個目錄的 **功能與定義的 interface** 整理：

* * *

1. `pkg/services/auth`
-----------------------

### 🔧 功能：

統一定義與管理授權（Authentication）相關介面與抽象，包括 Token 管理、ID 簽署、外部 session 存取等。

### 🔌 Interface：

*   `UserTokenService`：管理 user auth token，例如建立、查詢、撤銷等。
    
*   `ExternalSessionStore`：處理外部登入會話（如 OAuth）的 CRUD 操作grafana-pkg-services-al…。
    
*   `IDService`：為 identity 請求者簽發 ID Token。
    
*   `IDSigner`：以 signing key 實作簽發 JWT。
    
*   `IDClaims`：JWT 內容的結構別名。
    

* * *

2. `pkg/services/authimpl`
---------------------------

### 🔧 功能：

此為 `auth` interface 的具體實作（implementation），例如 user token 實作、session 儲存、token 清除等。

### 🔌 Interface 實作：

*   `UserAuthTokenService`：實作 `UserTokenService`，提供 token lifecycle 管理、revoked tokens 查詢、rotation 機制等。
    
*   `store`：實作 `ExternalSessionStore`，負責 DB 交易與密文儲存解密grafana-pkg-services-al…。
    

* * *

3. `pkg/services/gcomsso`
--------------------------

### 🔧 功能：

Grafana Cloud 的單一登出（SSO）hook 實作，提供整合外部系統的 logout 機制。

### 🔌 Interface：

未定義獨立的公開 interface，主要提供 `gcom_logout_hook.go` 中的邏輯，作為 `authn` 登出 hook 的實作模組。

* * *

4. `pkg/services/idimpl`
-------------------------

### 🔧 功能：

提供 ID Token 的本地簽署實作，內部使用 JWT + JOSE + signingkeys，支援組織角色、快取、claim 建構等。

### 🔌 Interface 實作：

*   `LocalSigner`：實作 `IDSigner`，透過 `signingkeys.Service` 提供簽署服務grafana-pkg-services-al…。
    

* * *

5. `pkg/services/jwt`
----------------------

### 🔧 功能：

專責 JWT 的驗證邏輯，包括 key set 檢查、claim 驗證、token 解碼等。

### 🔌 Interface：

*   `JWTService`：
    
    ```go
    type JWTService interface {
        Verify(ctx context.Context, strToken string) (map[string]any, error)
    }
    ```
    
*   `FakeJWTService`：測試用 mock 實作grafana-pkg-services-al…。
    

* * *

如果你有需要，我也可以幫你畫出這些 interface 的 UML 關聯圖，或補上 registry 的註冊方式與注入流程。是否要繼續？
