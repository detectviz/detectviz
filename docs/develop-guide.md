# 開發者導引（Develop Guide）

本文件提供 detectviz 開發者進行模組實作、擴充、測試時的設計參考與作業流程說明，幫助快速理解整體邏輯與模組邊界。

---

## 架構分層與開發原則

DetectViz 採用**組合式架構 (Composable Architecture)**，具備以下分層結構：

### 框架穩定層（核心不變）
- `apps/`：應用組合層（server、cli、agent、testkit）
- `pkg/platform/`：平台核心抽象（registry、composition、contracts）
- `internal/platform/`：平台實作層（註冊機制、組合引擎、生命週期管理）
- `plugins/core/`：核心插件（內建必要功能）

### 應用擴展層（隨組合擴展）
- `compositions/`：平台組合方案（不同應用組合定義）
- `pkg/domain/`：領域模型（業務實體）
- `internal/services/`：服務層（業務邏輯實作）
- `internal/adapters/`：適配器層（外部系統介接）
- `plugins/community/`：社群插件（功能擴展）

### 開發約束原則
1. **依賴方向**：apps/ → internal/ → pkg/，禁止反向依賴
2. **插件隔離**：plugin 間不可直接依賴，透過 contracts 介面互動
3. **介面契約**：所有模組對應 `docs/interfaces/*.md` 定義規格
4. **組合透明**：模組組合邏輯在 `compositions/` 中明確定義

---

## Scaffold 架構與技術原則

Detectviz 採用 Plugin 為核心的可擴充架構，設計重點如下：

### 插件設計原則

- 每個 plugin 應符合 SRP 原則，職責單一明確
- Interface 駝峰定義於 `pkg/platform/contracts/`，並支援組合式註冊
- Plugin 需支援：
  - 動態註冊與 lifecycle 控制
  - 配置 (`config`) 與 schema 驗證
  - 開關控制（enabled）
- plugin 類型範例如下（皆支援組合式註冊與 enable 控制）：
  - Importer（如：`prometheus`, `telegraf`, `redis`, `csv`）
  - Exporter（如：`influxdb`, `alertmanager`, `slack`）
  - Notifier（如：`email`, `webhook`, `grafana`）
  - Authenticator（如：`keycloak`, `ldap`, `saml`, `oauth2`）
  - Middleware（如：`jwt`, `gzip`, `requestmeta`）
  - WebUIPlugin（如：系統狀態頁、plugin 設定頁）
  - Tools（如：`supportbundles`, `inject-debug-id`, 開發用中介插件）

### 開發技術棧（後端）

- 語言與框架：Go 1.22+、spf13/cobra（支援子指令與 flags）、Echo v5  
- 伺服器 middleware：Echo middleware（Recover, Logger, Gzip 等）
- 架構與模組：Clean Architecture 分層、支援 plugin + lifecycle 組裝
- 設定管理：Viper（支援 hot reload）、`pkg/config/loader`, `pkg/config/schema`
- 設定解析：mapstructure（cfg 解碼為 struct，配合 plugin config 使用）
- Plugin 管理：Registry、Composition 組合、Lifecycle 控制器
- 任務排程與佇列：Redis Streams、Cron 排程器（預留）
- 日誌整合：otelzap（Zap logger with OTEL context，可結合 lumberjack 實作 log rotation）
  - 支援 log level、json/text 格式、trace context、log sink 組合
  - 建議全系統統一使用 otelzap 作為唯一 logger
  - 若需 log rotate，建議結合 lumberjack 套件（如 app.log 輪替、壓縮、保留）
- OTEL 資源註解：`go.opentelemetry.io/otel/sdk/resource`（提供 org, host, platform tag）
- Observability：OTEL SDK（trace, log, metric），整合 Prometheus、Tempo、Loki、Alloy DevKit
- OTEL 導入元件：
  - `otelecho`：Echo middleware trace wrapper（主要 HTTP entrypoint）
  - `opentelemetry-net-http`：標準 http.Client 外部呼叫追蹤
  - `opentelemetry-database-sql`：若使用 `database/sql` 操作資料庫，支援 SQL trace
  - `opentelemetry-gorm`：若 plugin 使用 GORM，可導入 DB ORM trace
  - `opentelemetry-go-contrib/instrumentation/google.golang.org/grpc/otelgrpc`：gRPC interceptor，支援 trace context 傳遞與 metrics 自動導出


### 開發技術棧（前端）

- 框架與樣板：HTMX + Echo SSR + Templ（頁面產生）
- UI 組件：AdminLTE, Tabulator.js（表格視覺化）
- 導覽與權限：WebUIPlugin 註冊 nav node + JWT 權限驗證
- Plugin 注入：支援前端 plugin 載入自定頁面與元件
- iframe 整合：Grafana iframe（支援 var 組織切換、token 傳遞）

### Scaffold 設計示意

```text
[DetectViz Applications]
   ↓ OTLP (gRPC/HTTP)
[Grafana Alloy Agent]
   ├─→ Traces → Grafana Tempo
   ├─→ Logs → Grafana Loki  
   ├─→ Metrics → Grafana Mimir
   └─→ Dashboard → Grafana iframe 嵌入
```

### Alloy 可觀測性整合

DetectViz 整合 **Grafana Alloy** 作為統一的可觀測性代理：

- **Config-Driven 監控**：透過 `alloy-config.river` 統一管理監控配置
- **OTLP 原生支援**：完整支援 OpenTelemetry Protocol
- **多語言 SDK**：提供 Go、Python 等語言的整合範例
- **系統服務化**：支援 systemd 等系統服務管理
- **自動化部署**：透過 `internal/services/observability/alloy_manager.go` 管理

---

## 開發流程建議

1. **閱讀接口設計**
   - 依據 `/docs/interfaces/*.md` 確認模組輸入/輸出、依賴關係
   - 了解相依的 service, store, plugin 使用方式

2. **撰寫模組 scaffold**
   - 每個模組應包含 `.go` 主檔與 `_test.go` 測試
   - 路徑建議依 `/todo.md` 所列方式建立

3. **撰寫對應文件**
   - 模組設計應記錄於 `docs/architecture/*.md`
   - 包含設計目標、資料流、可插拔點與測試建議

4. **依照分層原則設計**
   - handler → service → store → plugin 不可跨層耦合
   - interface 必須抽象於 `pkg/`，實作於 `internal/`

5. **模組測試與 mock**
   - 所有服務模組皆應有 interface 測試與實作測試
   - 可使用 `internal/test/` 中 fake/mock 測資

---

## 依賴管理與架構約束

### 依賴方向規則

為避免循環依賴，嚴格遵循以下依賴方向：

```bash
# 允許的依賴方向 (A → B 表示 A 可以依賴 B)
plugins/ → pkg/platform/contracts/     # Plugin 實作契約介面
plugins/ → pkg/shared/                 # Plugin 使用共用工具
internal/ → pkg/                       # 內部實作依賴公共介面
apps/ → internal/                      # 應用層依賴內部實作
apps/ → pkg/                          # 應用層依賴公共介面

# 禁止的依賴方向
pkg/ ❌→ internal/                     # 公共介面不可依賴內部實作
pkg/ ❌→ plugins/                      # 公共介面不可依賴具體插件
internal/platform/ ❌→ plugins/        # 平台核心不可依賴具體插件
```

### Plugin 隔離約束

```go
// ✅ 正確：Plugin 透過契約介面互動
type PrometheusExporter struct {
    registry contracts.Registry  // 透過 registry 取得其他服務
    logger   shared.Logger       // 使用共用工具
}

// ❌ 錯誤：Plugin 直接依賴其他 Plugin
type PrometheusExporter struct {
    influxImporter *influxdb.Importer  // 不可直接依賴其他 plugin
}

// ✅ 正確：透過事件或 registry 間接互動
func (p *PrometheusExporter) Export(data any) error {
    // 透過事件匯流排通知其他 plugin
    p.eventBus.Publish("data.exported", data)
    return nil
}
```

---

## 命名與版本化建議

- handler 分支版本：`v1/`, `v1beta1/` 路徑區分
- interface 檔案：以業務意圖命名，例如 `AlertNotifier`, `RuleEvaluator`
- plugin 命名：註冊名稱需唯一，例如 `"importers.prometheus"`, `"notifier.slack"`

---

## 插件路徑與分類說明

所有 plugins 依照功能與信任層級分類為：

### 核心插件：`plugins/core/`（框架穩定層，平台啟動必需）
- `auth/`：認證策略（basic、jwt、session）
- `middleware/`：HTTP 中介層（cors、ratelimit、logging、recovery、metrics）
- `hooks/`：平台級事件 hook 系統

### 社群插件：`plugins/community/`（應用擴展層，依 composition 載入）
- `importers/`：資料匯入插件（Telegraf Input 模式）
- `exporters/`：資料匯出插件（Telegraf Output 模式）
- `integrations/`：第三方整合（observability、notification、security、system、processors）
- `visualizers/`：視覺化整合（trace、topology、dashboard-builder）
- `web/`：Web UI 擴展（themes、components、pages、navtree）

### 工具插件：`plugins/tools/`（開發除錯使用）
- `generators/`、`validators/`、`converters/`
- `supportbundles/`、`middleware/inject-debug-id/`

### 組合定義：`compositions/`（平台組合方案）
- `minimal-platform/`：框架機制測試用組合
- `monitoring-stack/`：監控堆疊組合
- `observability-platform/`：完整可觀測性平台組合
- `alloy-devkit/`：Alloy 可觀測性開發套件

---

## Logger

- [ ] 導入全域 Logger 管理模組
  - 路徑建議：`pkg/shared/log/logger.go`
  - 功能：
    - 採用 `otelzap` 為唯一 logger，統一平台日誌行為
    - 搭配 lumberjack 支援 log rotation（大小、備份、壓縮）
    - 可透過 `log.L(ctx)` 取得帶有 trace context 的 logger 實例
    - 提供 `log.SetGlobalLogger()` 於 scaffold 啟動時設置預設 logger

- [ ] 擴充 LoggingConfig 結構
  - 路徑建議：`pkg/config/schema/types.go`
  - 說明：
    - 新增 `Type`, `OTEL`, `FileConfig` 等欄位
    - 支援 log 輸出選項（stdout/file）、格式（json/text）

- [ ] 升級 plugin 與 scaffold 中所有 `fmt.Printf` 為 `log.L(ctx).Info/Debug(...)`
  - 範圍包含：
    - lifecycle.go
    - prometheus, jwt, influxdb 等 plugin
    - scaffold_test 與 webui_test 中的 log 輸出

- [ ] 撰寫日誌初始化與整合測試
  - 確認 logger 能從 `LoggingConfig` 正確初始化
  - 驗證 trace ID 是否成功注入至 log 訊息中
  - 測試檔案位置建議：`internal/test/integration/logging_test.go`
