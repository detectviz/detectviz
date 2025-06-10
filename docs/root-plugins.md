


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

所有 plugin 建議放置於：

```
plugins/
├── datasources/
│   ├── rulestore_mysql/
│   ├── loggerstore_logfile/
│   └── metricsstore_influxdb/
├── middleware/
│   ├── cors/
│   └── ratelimit/
├── auth/
│   ├── keycloak/
│   └── static/
├── handlers/
│   └── custom_alerts/
```

---

## Plugin 實作原則

- 每個 plugin 應實作 interface 並呼叫 `RegisterXXX(...)` 完成註冊
- 所有註冊行為集中於 `init()` 時完成
- plugin 應不主動依賴其他 plugin，僅與 interface 接口互動

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