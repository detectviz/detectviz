# DetectViz 架構總覽

本文件說明 DetectViz 專案的整體架構思維、模組分層、依賴原則與擴充策略，作為設計、開發與組裝新應用的參考基礎。

---

## 架構設計目標

- **Composable First**：以平台可組合性為核心，支援多種應用組裝與功能裁切。
- **Clean Architecture**：依據分層原則劃分 handler / service / store / domain。
- **Plugin 注入機制**：所有擴充皆透過 interface 實作並由 registry 管理注入。
- **框架穩定、模組擴展**：apps、platform 層穩定，services、plugins 可隨組合擴充。

---

## 架構層級說明（由內而外）

| 層級 | 模組             | 說明 |
|------|------------------|------|
| 1️⃣  | `pkg/platform`    | 提供 contracts, composition, registry 能力 |
| 2️⃣  | `pkg/domain`      | 各業務邏輯的 interface 與模型結構 |
| 3️⃣  | `internal/services` | 各功能模組實作，如 alerting, rule, notifier |
| 4️⃣  | `internal/store`    | 多資料源存取實作（e.g. redis, influx） |
| 5️⃣  | `internal/infrastructure` | logging, eventbus, tracing 等底層支援 |
| 6️⃣  | `plugins/core`、`plugins/community` | plugin 實作模組，向上註冊 |
| 7️⃣  | `compositions/`   | 平台組合與應用場景（如 monitoring-stack） |
| 8️⃣  | `apps/`           | 應用主程式（server, cli, agent） |

---

## 架構呼叫流程圖（簡化）

```
apps/server → compositions/ → platform/composition
    → platform/registry → services/xxx
        → store/, plugins/, infrastructure/
```

- Plugins 僅向上實作 interface，經由 registry 掃描後注入
- Services 僅依賴 domain/interface 與 registry 注入之實作
- Infra 不可反向依賴 services

---

## Plugin 類型與註冊位置

| 類型       | 註冊點                            | 功能範圍 |
|------------|------------------------------------|----------|
| Middleware | `plugins/core/middleware/`         | CORS、限速、自訂中介層 |
| Auth       | `plugins/core/auth/`               | OAuth、Basic、SSO 登入策略 |
| Notifier   | `plugins/community/notifiers/`     | Email, Webhook |
| Datasource | `plugins/community/importers/`     | Redis, MySQL, Influx |
| Exporter   | `plugins/community/exporters/`     | Log, Tracing, Metrics |
| Integration| `plugins/community/integrations/`  | Grafana, Keycloak, Prometheus |

---

## Scaffold 重用與平台組裝

DetectViz 核心以組合平台為主，可快速派生以下應用：

- ✔ 偵測與告警平台
- ✔ 可視化與審查平台
- ✔ 測試資料注入器
- ✔ 匯入匯出工具（log / metrics）
- ✔ 自訂儀表板平台

---

## 開發與維護建議

- 每個模組應補對應 `docs/interfaces/*.md`
- 所有架構建議與順序請參考 `/todo.md`
- 建議每個模組 scaffold 含 `_test.go` 測試與 mock registry

---
