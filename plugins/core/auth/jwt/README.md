# JWT Authentication Plugin

JWT 認證插件提供基於 JSON Web Token 的認證功能。

## 功能特性

- 支援 JWT token 生成與驗證
- 可配置的密鑰和簽名演算法
- 支援 token 過期時間設定
- 整合 DetectViz 權限系統

## 配置範例

```yaml
name: jwt-authenticator
type: auth
config:
  secret_key: "your-secret-key"
  issuer: "detectviz"
  expiry_time: "24h"
  refresh_time: "168h"  # 7 days
  signing_method: "HS256"
```

## 使用方式

### 註冊插件

```go
import "detectviz/plugins/core/auth/jwt"

// 註冊到 registry
err := jwt.Register(registry)
```

### 認證流程

1. 使用者提供憑證
2. 插件驗證憑證
3. 生成 JWT token
4. 回傳使用者資訊

## 支援的認證類型

- `jwt`: JWT token 認證

## 依賴項目

- 無外部依賴
- 使用標準 Go 函式庫

## TODO

- [ ] 實作實際的 JWT 簽名/驗證邏輯
- [ ] 支援 RSA 簽名演算法
- [ ] 整合外部身份提供者
- [ ] 新增 token 撤銷機制 