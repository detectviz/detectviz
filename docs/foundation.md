# Detectviz 平台基礎設計（Foundation Overview）

本文件說明 detectviz 架構設計背後的理念與開發哲學，並作為多專案重用、快速 scaffold 的核心依據。

---

## 設計目標

- 建立可抽換、可測試、可擴充的模組平台
- 將架構設計本身作為 scaffold，使任一 side project 都可快速成型
- 將 Web / CLI / API / plugin / config / service 各層分離，並清楚定義注入方式
- 提供完整可組合平台，減少重複開發與耦合

---

## 結構核心與應用邊界

detectviz 並非單一應用，而是一組 **技術中台與運行時組件**，可應對以下使用場景：

- Web UI 平台（儀表板、設定頁、分析畫面）
- CLI 自動化工具（批次套件、同步工具）
- Plugin 驅動的模組（importer/exporter/middleware/authenticator）
- 抽象化儲存後端（file, Redis, Influx, MySQL）
- API 服務統一注入與版本管理（v1, v1beta1）

---

## 如何從 detectviz 快速開發一個 Side Project？

1. **複製 core 架構**
   ```bash
   cp -r detectviz/{internal,pkg,apps,plugins,compositions,docs} my-project/
   ```

2. **選擇應用組合**
   - 需要 web UI？使用 `internal/web/` + `apps/server/`
   - 需要 CLI？使用 `pkg/cmd/` + `apps/cli/`
   - 需要事件與儲存？註冊 plugin → `internal/plugins/`

3. **撰寫功能模組**
   - handler → service → store
   - 可選擇 plugin 化（放入 plugins/community/...）或寫在 core/internal（platform 內建）

4. **啟動應用**
   ```bash
   go run apps/server/main.go
   ```

---

## 模組重用策略

每個模組皆符合以下條件，即可於多個 side project 中重用：

| 條件 | 描述 |
|------|------|
| interface/implementation 分離 | interface 在 `pkg/`，實作在 `internal/` |
| 支援 DI 或 plugin 注入 | 如 AuthStrategy, Store, Middleware |
| 測試與錯誤獨立處理 | 使用 `_test.go` 與 `pkg/errors` 規範 |
| 狀態管理去中心化 | 不依賴全域變數或 context 預設值 |

> 若模組實作為 plugin，請遵循：`docs/interfaces/plugins.md`、`pkg/platform/contracts/` 中介面定義。

---

## 建議模組搭配範例

| 專案類型 | 模組組合 |
|----------|----------|
| 規則審查工具 | `rule.handler` + `rule.service` + `logfile.store` |
| webhook router | `notifier.handler` + `service/dispatcher` + `plugin/middleware/hook` |
| 登入整合平台 | `auth.handler` + `plugins/community/integrations/security/keycloak` |
| Prometheus 整合平台 | `metrics.handler` + `plugins/community/importers/prometheus` + `plugins/core/middleware/logging` |

---

## 搭配開發指南與工具

- 每個模組對應：`docs/interfaces/*.md`、`docs/architecture/*.md`
- 開發順序與依賴：見 `/todo.md`，依據優先順序與實作建議進行
- 可配合 `golangci-lint`, `gofumpt`, `testify` 等工具進行驗證與風格統一

---

## 整體收斂

detectviz 並不僅是架構，而是：
> 一套支援你所有專案的「模組化平台骨架」

你可以任意將其拆解、組合、擴充，讓每一個 side project 都能只專注在你要解決的「業務問題」，不必再處理重複的 CLI、Web、API scaffold。