# Plugin Architecture

本文件說明 detectviz 中 Plugin 架構設計原則與模組整合方式，支援動態擴充功能、介面註冊與模組注入。

---

## Plugin 模組目的

- 動態擴充特定模組實作（如 store、auth、middleware、handler）
- 保持核心模組精簡、可維護
- 支援不同儲存後端或擴充邏輯模組的載入與註冊

---

## 支援 plugin 的模組

| 模組        | Plugin 類型                | 註冊目錄                      |
|-------------|-----------------------------|-------------------------------|
| store       | RuleStore, NotifierStore... | `internal/store/resolver.go` |
| middleware  | CORS, RateLimit 等           | `internal/registry/middleware` |
| auth        | AuthStrategy 擴充            | `internal/registry/auth`     |
| handlers    | 自訂 API handler             | `internal/api/router.go` 或模組註冊器 |
| metrics     | Exporter 擴充                | `internal/registry/metrics`（待定） |

---

## Plugin 路徑建議

所有 plugin 應統一放置於：

```
internal/plugins/
├── plugins/               # 可插拔模組擴充（註冊 middleware、auth 等）
│   ├── auth/              # 額外擴充的登入策略（如 keycloak, static）
│   ├── middleware/        # 中介層插件（如 cors、ratelimit）
│   ├── apihooks/          # 註冊外部 API handler 或擴充點
│   ├── eventbus/
│   │   └── alertlog/      # 事件匯流處理插件（如 alert log sink）
│   └── plugin.go          # 插件註冊總入口
├── manager/               # plugin lifecycle 載入與管理邏輯
│   ├── loader.go
│   ├── registry.go
│   ├── process.go
│   └── lifecycle.go
```

---

## Plugin 實作原則

- 每個 plugin 應位於 `internal/plugins/plugins/{模組}/` 目錄下，並依照介面實作
- 註冊行為統一於該目錄 `init()` 或透過 `manager/registry.go` 載入
- plugin 不應直接依賴其他 plugin，僅透過 interface 與 registry 溝通

---

## Plugin Lifecycle 管理

`internal/plugins/manager/` 包含完整的 plugin lifecycle 管理流程：

- `loader.go`：載入 plugins 之邏輯
- `registry.go`：註冊並掛載 plugins
- `process.go`：處理啟用流程與初始化
- `lifecycle.go`：控制 plugin 啟動與關閉流程

此區模組未直接依賴外部實作，僅處理注入與調用

---

## 範例：註冊自訂 RuleStore plugin

```go
package rulestore_mysql

func init() {
    store.RegisterRuleStore("mysql", NewMySQLRuleStore())
}
```

---

## Plugin 啟用策略（未來）

- 支援透過 config 啟用特定 plugin：
  ```yaml
  store:
    rule:
      backend: "mysql"
  ```
- 可依環境選擇 plugin 優先順序或 fallback

---

## Plugin 測試建議

- 每個 plugin 應搭配 fake context 單元測試
- plugin 可作為 e2e 測試的組件之一
- 禁止 plugin 內部使用全域變數與副作用行為

---

## 參考文件

- [docs/interfaces/store.md](../interfaces/store.md)
- [docs/interfaces/middleware.md](../interfaces/middleware.md)
- [docs/interfaces/auth.md](../interfaces/auth.md)
- [docs/architecture/store.md](./store.md)