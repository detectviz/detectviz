# Infra Architecture

本文件說明 `pkg/infra/` 模組在 detectviz 專案中的角色與設計原則，明確區分其與 `internal/system/` 的責任邊界。

---

## 模組定位

`pkg/infra/` 為專案提供靜態的技術基礎建設封裝，主要負責初始化、封裝與提供以下功能：

- 全域 logger 初始化與注入（ex: zap logger）
- metrics 註冊與 export handler
- config 解析與驗證
- DB 驅動註冊（如 mysql, sqlite）
- caching 與 registry adapter
- OpenTelemetry 或 trace context 工具
- signal/context 管理

---

## 與 internal/system 的分工

| 項目         | pkg/infra/                      | internal/system/                         |
|--------------|----------------------------------|------------------------------------------|
| 設計階段注入 | ✅ ex: config, logger, trace     | ❌                                      |
| Runtime 管理 | ❌                               | ✅ lifecycle, hooks, eventbus, quota 等  |
| 通用封裝     | ✅                               | ❌                                      |
| 模組入口     | 由外部或 bootstrap 呼叫           | 被服務模組調用                           |

---

## 建議目錄結構

```
pkg/infra/
├── log/             # 統一 logger 封裝（zap）
├── config/          # 組態載入與驗證（env/yaml/flag）
├── trace/           # trace context 與 OpenTelemetry helper
├── db/              # 資料庫驅動註冊與 sqlx 管理
├── cache/           # redis 或 memory cache client adapter
├── metrics/         # Prometheus 註冊器與 middleware
├── httpclient/      # 統一 HTTP 調用封裝（支援 middleware, retry, auth）
├── signal/          # 處理 SIGINT、SIGTERM 等系統訊號
```

---

## 使用規範

- 所有模組皆可安全依賴 `pkg/infra` 元件
- 不得反向依賴 `internal/*`
- 不得混入業務邏輯，只封裝基礎設施
- 每個元件皆應支援 mock 或 wrapper

---

## 延伸建議

- 每個封裝模組皆可提供 interface，利於測試注入與 plugin 替換
- 支援 lazy loading / auto init pattern
- 日後可依主機平台調整實作（如支持 cloud config provider）

---

## HTTP Client 模組說明

建議將 `http.Client` 功能封裝為 `pkg/infra/httpclient`，提供：

- 統一 timeout、retry、headers 設定
- 中介行為注入（如 metrics、trace、auth header）
- 對外使用泛型 `Doer` interface，利於測試與替換

常見建構方式：

```go
cli := httpclient.New(
    httpclient.WithTimeout(5 * time.Second),
    httpclient.WithRetry(3),
    httpclient.WithAuth("Bearer token"),
)
```

模組不應直接依賴外部配置，可由 config 模組注入。
