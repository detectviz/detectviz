# System Architecture

本文件說明 detectviz 系統的整體架構設計，包括模組邊界、執行流程、模組注入與插件擴充機制，以及核心模組之間的協作關係。

---

## 核心目標

- 建立可觀測、可擴充的監控分析平台
- 以 Plugin 為擴充基礎，支援各模組自訂邏輯與資料來源
- 提供事件驅動、組件化的模組調度架構
- 可支援 Web UI / CLI / API 介面共用服務層與資源

---

## 系統模組總覽

```
┌───────────────┐
│  web          │ ← HTMX Web UI
│  api          │ ← REST API handler
└────┬──────────┘
     ↓
┌───────────────┐
│  middleware   │ ← auth, logging, tracing, cors...
└────┬──────────┘
     ↓
┌───────────────┐
│  services     │ ← 業務邏輯處理（rule, notifier 等）
└────┬──────────┘
     ↓
┌───────────────┐
│  store        │ ← 資料存取與快取（支援多後端）
└───────────────┘
     ↑
┌───────────────┐
│  plugins      │ ← middleware, auth, handler 等插件
└───────────────┘
```

---

## 系統支援模組（internal/system/）

作為跨模組技術支援區域，提供診斷、整合與平台能力，推薦結構如下：

```
internal/system/
├── http/              # 提供底層 HTTP 能力，如 apiserver, proxy, cors
├── diagnostics/       # 系統健康與診斷模組，如 supportbundle, stats
├── integration/       # 外部整合介接，如 grpcserver, search, live update
├── platform/          # 通用平台能力，如配額、快取、鉤子 hooks
```

此區模組不承擔商業邏輯，為 services/plugin 提供能力與環境。

---

## 應用程式進入點（apps/）

detectviz 專案將各類應用程式（CLI 工具、API Server、測試工具等）獨立放置於 `apps/` 目錄中，每個子目錄對應一個獨立執行目標（binary）。建議結構如下：

```
apps/
├── cli/
│   └── main.go       # 執行 detectviz CLI 工具，呼叫 pkg/cmd/root.Execute()
├── server/
│   └── main.go       # 啟動 Web Server 與 API 路由，呼叫 bootstrap.Run()
├── testkit/          # （選用）整合測試與模擬工具入口
│   └── main.go
```

- 每個子目錄皆為一個獨立應用程式，可獨立編譯或部署
- `apps/cli/` 與 `apps/server/` 共用 `pkg/cmd/`、`internal/bootstrap/`、`services/` 等模組
- CLI 執行檔與 Web Server 可使用不同的設定檔或啟動參數

---

## 模組層級職責

| 層級         | 說明 |
|--------------|------|
| `web`        | HTMX 前端頁面渲染，使用 templ 組件 |
| `api`        | Echo handler、API 版本化與 middleware 掛載 |
| `middleware` | 共用中介層處理：Auth、CORS、Log、Tenant |
| `services`   | 業務核心邏輯，處理 rules、alerts、metrics |
| `store`      | 資料存取層，支援 MySQL / InfluxDB / Redis 等 |
| `plugins`    | 可插拔擴充元件，依據模組註冊託管實作 |
| `system`     | 支援性元件與平台整合能力，劃分為 http、diagnostics、integration、platform |

---

## 核心擴充點

- `internal/plugins/plugins/auth/`：登入策略擴充
- `internal/plugins/plugins/middleware/`：middleware 插件
- `internal/plugins/plugins/apihooks/`：API 擴充路由
- `internal/store/{mod}/{backend}/`：資料來源切換與註冊
- `internal/services/{mod}/`：service 注入並整合其他模組

---

## Plugin Lifecycle

- 所有 plugin 註冊由 `internal/plugins/manager/` 控制
- 每個 plugin 實作註冊介面，於 `init()` 時自註冊
- Registry 與 Loader 控制整體初始化順序與依賴

---

## 多介面支持

- Web：HTMX + templ + go embed 頁面
- API：RESTful 結構，對應 Echo router 與 handler
- CLI（預留）：future goal，可重用同一套 service/store 架構

---

## 可觀測性與事件流程

- 所有 API 請求皆經 middleware 記錄 trace 與 context
- Service 可透過 `eventbus` 發出事件，串接 `logger` 或 `notifier`
- 支援事件串流進 AlertLog / LogStore / Remote sink

---

## 系統整合測試建議

- Init → API Router → Middleware → Service → Store → Plugin 全流程測試
- 可搭配 fake component 完成 mock 與流程驗證
- 支援 e2e 模式執行測試與 UI 流程模擬

---