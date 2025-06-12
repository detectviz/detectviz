# RBAC 架構設計（Role-Based Access Control）

本文件說明 detectviz 專案中的 RBAC 權限模組（internal/rbac）整體架構設計、模組拆分、核心元件職責與插件化策略，供平台內部與外部擴充參考。

---

## 設計目標

- 支援多租戶與多組織切換機制（org）
- 封裝角色與資源範疇控制邏輯（scope, role, evaluator）
- 提供靈活的授權驗證與 middleware 注入
- 可測試、可 mock、支援 plugin 註冊的策略模組（如 service account, temp user）

---

## 模組目錄結構規劃

```
internal/rbac/
├── accesscontrol/       # 權限比對邏輯與授權流程
│   ├── authorizer.go    # 主授權介面與 dispatch
│   ├── checker.go       # 快取與角色比對
│   ├── scope.go         # 定義資源層級與授權範疇
│   ├── middleware.go    # 整合到 HTTP handler 的 middleware 授權
├── org/                 # 組織/租戶管理與刪除邏輯
│   ├── org.go
│   ├── model.go
│   └── org_delete_svc.go
├── team/                # 群組管理（可選）
│   └── team.go
├── mock/                # 測試專用的 mock/fake 實作
│   ├── fake.go
│   └── mock.go
```

---

## 關鍵元件與職責

### Authorizer

- 負責驗證使用者是否擁有某資源的存取權限
- 對應介面：

```go
type Authorizer interface {
    Check(user *User, action string, resource string, scope Scope) (bool, error)
}
```

### Scope

- 表示權限作用範疇，例如：
  - 全域 (`global`)
  - 組織 (`org_id`)
  - 資源層級 (`dashboard:123`)
- 可擴充定義樹狀關係與父層推導

### Checker / Evaluator

- 快取權限查詢結果
- 提供具效能的比對機制
- 可擴充 RBAC policy 判斷模式

### Middleware

- 整合至 Echo/Gin HTTP handler 中
- 可對路由設定 required roles, scope 等屬性

---

## Plugin 支援模組（非必要，可注入）

### `plugins/serviceaccounts/`

- 提供 CI/CD 或自動化流程帳號
- 可掛載於 Auth middleware 供 token 驗證

### `plugins/tempuser/`

- 提供臨時登入憑證與自動銷毀策略
- 常用於邀請機制與 One-time access

---

## 測試與模擬策略

- mock/fake 實作位於 `internal/rbac/mock/`
- 可模擬授權流程、角色配置與 scope 調用
- 適用於 integration 與 e2e 測試流程

---

## 延伸建議

- 每個 scope, role, resource 定義皆應標準化命名與 metadata 備註
- plugin 可註冊動態角色對應邏輯（如 alert reviewer, rule editor）
- 未來可支援 policy 註冊器與 DSL（如：`allow if org_id == x`）

---