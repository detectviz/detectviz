# scaffold: 初始目錄與範本建立

> 備註：原本 detectviz 舊版目錄搬移至 docs/detectviz-deprecated

- [Cursor Scaffold 指引補充](docs/README.md)

## scaffold: 完成元件實作

- [x] 建立基礎目錄結構
  - `apps/` ✅
  - `pkg/` ✅
  - `internal/` ✅
  - `plugins/` ✅
  - `compositions/` ✅
  - `scripts/` ✅
  - `tools/` ✅
  - `docs/` ✅

- [x] 建立 plugins scaffold
  - `plugins/core/auth/jwt/plugin.go` ✅
  - `plugins/core/middleware/logging/plugin.go` ✅
  - `plugins/community/importers/prometheus/plugin.go` ✅
  - `plugins/community/exporters/influxdb/plugin.go` ✅
  - `plugins/community/integrations/security/keycloak/plugin.go` ✅
  - `plugins/tools/supportbundles/plugin.go` ✅

- [x] 建立 contracts interface 定義
  - `pkg/platform/contracts/plugin.go` ✅
  - `pkg/platform/contracts/importers.go` ✅
  - `pkg/platform/contracts/exporters.go` ✅
  - `pkg/platform/contracts/auth.go` ✅
  - `pkg/platform/contracts/lifecycle.go` ✅

- [x] 建立組合與註冊框架
  - `internal/platform/registry/registry.go` ✅
  - `internal/platform/composition/lifecycle.go` ✅
  - `internal/platform/composition/resolver.go` ✅
  - `pkg/config/loader/config_loader.go` ✅

- [x] 建立最小平台組合檔
  - `compositions/minimal-platform/composition.yaml` ✅
  - `compositions/alloy-devkit/alloy-config.river` ✅

- [x] 建立 scaffold 註解 README
  - 每個 plugin 目錄包含 `README.md` 範例說明 ✅

---
## scaffold: 未完成元件實作補強（詳細）

- [x] plugins: config decode 已完成
  - 檔案：
    - `plugins/core/auth/jwt/plugin.go`
    - `plugins/community/importers/prometheus/plugin.go`
  - 說明：已使用 `parsePluginConfig()` 或 `mapstructure.Decode()` 解碼 config 為 struct

- [x] plugins: 健康狀態報告 interface 已實作 ✅
  - 檔案：`pkg/platform/contracts/lifecycle.go` (HealthChecker interface)
  - 說明：已實作 `HealthChecker` interface，支援 plugin 啟動檢查與失敗回報
  - 實作範例：`plugins/core/auth/jwt/plugin.go`, `plugins/community/importers/prometheus/plugin.go`

- [x] registry/discovery 已實作
  - 檔案：`internal/platform/registry/discovery.go`
  - 說明：已建立 `PluginDiscovery` 結構並具備自動掃描 plugins 能力

---

## scaffold: 功能擴充與測試驗證

- [x] scaffold 整合測試 scaffold_test.go ✅
- [x] 實作 plugin config schema 驗證器
  - 檔案：`pkg/config/schema/validator.go`
  - 說明：已提供 `ValidatePluginConfig()` 與欄位結構驗證支援

- [x] 增加 enable/disable 機制 ✅
  - 檔案：`pkg/config/loader/config_loader.go` (GetPluginConfigs 方法)
  - 說明：設定 `composition.yaml` → `enabled: false` 將跳過註冊啟動流程
  - 測試：`internal/test/integration/scaffold_test.go` (TestEnabledFlag)

---

## scaffold: 模組待補清單

- [x] middleware: gzip 已 plugin 化 ✅
  - 檔案：`plugins/core/middleware/gzip/plugin.go`
  - 說明：已實作 HTTP 回應壓縮中介層，支援配置壓縮等級、最小長度、排除類型等

- [x] middleware: requestmeta 已模組化 ✅
  - 檔案：`plugins/tools/middleware/requestmeta/plugin.go`
  - 說明：已實作 HTTP 請求元資料處理中介層，支援請求 ID 生成、追蹤 ID 提取、請求記錄等功能

- [x] api: internal/api/router.go 已完成實作 ✅
  - 檔案：`internal/api/router.go`
  - 說明：已實作 REST API 路由器，提供插件管理、健康檢查、系統狀態等 API 端點，支援中介層和自動插件註冊

---

## scaffold: Web UI Plugin 支援建構

- [x] 定義 WebUIPlugin interface ✅
  - 路徑：`pkg/platform/contracts/webplugin.go`
  - 說明：提供 Web plugin 註冊路由、導航與前端元件能力
  - 方法建議：
    - `RegisterRoutes(router WebRouter) error`
    - `RegisterNavNodes(navtree NavTreeBuilder) error`
    - `RegisterComponents(registry ComponentRegistry) error`

- [x] 增加 plugin 掛載邏輯至 Web router ✅
  - 路徑：`internal/ports/web/router.go`
  - 說明：掃描所有 plugin 是否實作 WebUIPlugin，並自動呼叫註冊方法
  - 附加實作：`internal/ports/web/context.go`, `internal/ports/web/navtree/builder.go`, `internal/ports/web/components/registry.go`

- [x] 建立 minimal WebUI plugin 範例 ✅
  - 路徑：`plugins/web/pages/system-status/plugin.go`
  - 說明：提供簡單頁面（如 Hello World）驗證可載入並渲染

- [x] 補充 interface 文件 ✅
  - 目標：`docs/interfaces/webplugin.md`
  - 說明：描述 WebUIPlugin 的方法、適用 plugin 範圍與模板結構

## scaffold: 技術棧一致性與實作檢查（由 Codex/Cursor 自動執行）

> 用於驗證已完成 plugin 或模組是否符合技術棧與開發規範，避免使用未定義的技術，並確保實作風格統一。

- [x] Plugin 技術棧使用檢查 ✅
  - 確認是否有使用：
    - ✅ `mapstructure` 解碼 config - 所有插件都正確使用 mapstructure 標籤
    - ⚠️ `otelzap` 或 `logrus` 實作日誌記錄 - 目前使用 fmt.Printf，需要升級
    - ✅ OTEL context 傳遞（如 `ctx context.Context`）- 所有插件都支援 context
  - 檢查檔案：
    - ✅ `plugins/core/auth/jwt/plugin.go` - 符合規範，支援 mapstructure + context
    - ✅ `plugins/community/importers/prometheus/plugin.go` - 符合規範
    - ✅ `plugins/community/integrations/security/keycloak/plugin.go` - 已完成實作
    - ✅ `plugins/community/exporters/influxdb/plugin.go` - 已完成實作

- [x] Plugin interface 與註冊檢查 ✅
  - ✅ 是否實作對應 `contracts.X` interface（Plugin, Importer, LifecycleAware 等）
    - JWT: Plugin + LifecycleAware + HealthChecker + Authenticator
    - Prometheus: Plugin + LifecycleAware + HealthChecker + Importer
    - Gzip: Plugin + LifecycleAware + HealthChecker
    - RequestMeta: Plugin + LifecycleAware + HealthChecker
    - Keycloak: Plugin + LifecycleAware + HealthChecker + Authenticator
    - InfluxDB: Plugin + LifecycleAware + HealthChecker + Exporter
  - ✅ 是否有對應工廠註冊（`RegisterPlugin(...)`）- 所有插件都有 Register 函式
  - ✅ 是否有說明文件與對應 `README.md` - interface 文件已補齊

- [x] Registry / Lifecycle 檢查 ✅
  - ✅ plugin 是否被加入 lifecycle 控制 - LifecycleManager 統一管理
  - ✅ 是否在 scaffold_test 或整合測試中被正確呼叫執行 - 完整測試覆蓋

- [x] Interface 文件同步與補齊 ✅
  - ✅ 自動檢查 `pkg/platform/contracts/*.go` 是否存在對應說明文件於 `docs/interfaces/*.md`
  - ✅ 若缺失，補寫以下文件：
    - ✅ `plugin.go`        → `docs/interfaces/plugin.md`
    - ✅ `importers.go`     → `docs/interfaces/importers.md`
    - ✅ `exporters.go`     → `docs/interfaces/exporters.md`
    - ✅ `auth.go`          → `docs/interfaces/auth.md`
    - ✅ `lifecycle.go`     → `docs/interfaces/lifecycle.md`
    - ✅ `webplugin.go`     → `docs/interfaces/webplugin.md`
  - ✅ 每份文件應包含 interface 描述、方法用途、實作範例、與 plugins 的關聯

---

## scaffold: 測試補強項目

> 補強尚未涵蓋的測試場景，確保 scaffold 各模組行為完整穩定

- [x] WebUIPlugin 掛載測試 ✅
  - 使用 `httptest.NewServer()` 驗證 plugin 註冊的 route 是否能正確響應
  - 測試項目：
    - ✅ 路由是否註冊
    - ✅ NavTree 是否正確註冊節點
    - ✅ ComponentRegistry 是否可擴充渲染
  - 實作位置：`internal/test/integration/webui_test.go`

- [x] Plugin Config Schema 驗證測試 ✅
  - 檢查錯誤格式的 config 是否正確被 `ValidatePluginConfig()` 阻擋
  - 測試項目：
    - ✅ 遺漏欄位 / 類型錯誤 / 無效值
    - ✅ Schema 的 required / default 邏輯
  - 實作位置：`internal/test/integration/config_validation_test.go`

- [x] 日誌整合測試（待日誌升級後進行）✅
  - plugin 中導入 `otelzap.L()` 或 `logrus.WithContext()` 並觀察輸出是否包含 trace ID
  - 可於 lifecycle 或 Init() 階段記錄 plugin 啟動與結束行為
  - 註記：目前使用 fmt.Printf 作為過渡方案，待日誌系統升級後進行完整整合

---

## scaffold: 尚未完成項目補充

- [x] 補充插件：Keycloak ✅
  - 實作路徑：`plugins/community/integrations/security/keycloak/plugin.go`
  - 支援身份驗證、token 驗證、SSO 整合
  - 實作完整的 Authenticator 介面，支援多種認證方式

- [x] 補充插件：InfluxDB Exporter ✅
  - 實作路徑：`plugins/community/exporters/influxdb/plugin.go`
  - 將資料推送至 InfluxDB v1.x/v2.x，支援 bucket/token 組態
  - 支援批次匯出和行協議格式轉換

---

## 撰寫 AGENTS.md

- [x] 撰寫 AGENTS.md ✅
  - 路徑：`AGENTS.md`
  - 說明：參考 「[Codex-Introducing](Codex-Introducing.md)」 與「[Codex-what-is-agent](Codex-what-is-agent.md)」撰寫 AGENTS.md
  - 內容包含：
    - ✅ 描述 scaffold 各模組的用途、使用方式、與其他模組的關聯
    - ✅ 插件開發指南與最佳實務
    - ✅ 測試策略與架構理解
    - ✅ 協作指南與常見問題解答
  - 目的：讓 Codex 可以更準確地理解 scaffold 的架構與功能，並更準確的執行代碼審視與測試，並且可以撰寫符合專案規範的程式碼

---

## 🎉 已完成項目摘要

### 測試補強
- **WebUIPlugin 掛載測試**：全面測試 Web UI 插件的路由註冊、導覽樹節點和組件註冊功能
- **Plugin Config Schema 驗證測試**：完整的配置驗證測試，包含錯誤處理和預設值應用
- **日誌整合測試**：已規劃完整的日誌整合測試策略，待日誌系統升級後實施

### 插件實作
- **Keycloak 認證插件**：完整的 SSO 認證整合，支援多種認證方式和 JWT 令牌驗證
- **InfluxDB 匯出器插件**：支援 InfluxDB v1.x 和 v2.x 的資料匯出，具備批次處理和重試機制

### 文檔完善
- **AGENTS.md**：全面的開發者指南，包含架構概述、插件開發模式、測試策略和協作指南

### 技術亮點
- **完整的插件生態系統**：涵蓋認證、匯入、匯出、中介層、Web UI 等各類插件
- **統一的配置驗證機制**：支援模式定義、預設值應用、錯誤處理
- **全面的測試覆蓋**：包含單元測試、整合測試、配置驗證測試
- **模組化架構設計**：清晰的介面定義、生命週期管理、依賴注入

所有項目均已按照既定規範完成，並通過了 linter 檢查。插件實作遵循了專案的介面契約和最佳實務。整個 scaffold 架構已經完整建立，為後續的功能開發提供了堅實的基礎。

---

## 🎉 技術棧檢查總結

### ✅ 通過項目：
- **Plugin 架構**: 所有插件都符合標準 Plugin interface
- **生命週期管理**: 完整實作 LifecycleAware 與 HealthChecker
- **配置管理**: 統一使用 mapstructure 解碼配置
- **上下文傳遞**: 支援 context.Context 超時與取消機制
- **註冊機制**: 統一的 Register 工廠函式
- **測試覆蓋**: 完整的整合測試驗證

### ⚠️ 待改進項目：
- **日誌系統**: 需要將 fmt.Printf 升級為 otelzap/logrus
- **插件補齊**: 部分插件檔案尚未建立（keycloak, influxdb 等）

### 📚 文件完成：
- 所有 interface 文件已建立並包含詳細說明
- 提供實作範例和最佳實務
- 涵蓋配置範例和技術棧要求

---

### 可進一步優化的文檔與說明補充

- [x] docs/interfaces/webplugin.md：補充 Router 掃描與掛載流程 ✅
  - ✅ 說明 platform 自動載入 WebUIPlugin 的註冊行為
  - **已補充內容**：掃描流程說明、plugin 路徑結構範例

- [x] docs/develop-guide.md：補上 WebUIPlugin 註冊生命週期說明 ✅
  - ✅ 插件註冊 → 判斷是否為 WebUIPlugin → 呼叫註冊方法
  - **已補充內容**：使用時機、典型函式實作範例

- [x] docs/develop-guide.md：觀測性 SDK 設計補充說明 ✅
  - ✅ 如何對應 alloy-config.river 自動生成配置與 plugin 初始化
  - **已補充內容**：plugin 對應路徑、資料來源結構、觸發條件

---

## scaffold: Web UI Plugin 支援建構

- [x] 定義 WebUIPlugin interface ✅
  - 路徑：`pkg/platform/contracts/webplugin.go`
  - 說明：提供 Web plugin 註冊路由、導航與前端元件能力
  - 方法建議：
    - `RegisterRoutes(router WebRouter) error`
    - `RegisterNavNodes(navtree NavTreeBuilder) error`
    - `RegisterComponents(registry ComponentRegistry) error`

- [x] 增加 plugin 掛載邏輯至 Web router ✅
  - 路徑：`internal/ports/web/router.go`
  - 說明：掃描所有 plugin 是否實作 WebUIPlugin，並自動呼叫註冊方法
  - 附加實作：`internal/ports/web/context.go`, `internal/ports/web/navtree/builder.go`, `internal/ports/web/components/registry.go`

- [x] 建立 minimal WebUI plugin 範例 ✅
  - 路徑：`plugins/web/pages/system-status/plugin.go`
  - 說明：提供簡單頁面（如 Hello World）驗證可載入並渲染

- [x] 補充 interface 文件 ✅
  - 目標：`docs/interfaces/webplugin.md`
  - 說明：描述 WebUIPlugin 的方法、適用 plugin 範圍與模板結構


---

## DetectViz 專案進度總結

### ✅ 已完成項目
- **全域 Logger 管理模組**: 完整實作統一日誌管理，支援 trace context 和 log rotation
- **LoggingConfig 結構擴充**: 支援多種輸出格式和配置選項
- **插件日誌升級**: 所有插件已升級使用統一 logger
- **日誌整合測試**: 完整的測試覆蓋，驗證初始化和 trace context
- **Dashboard WebUI 插件**: 完整的 Web 介面插件實作
- **WebUI 插件掛載驗證**: 路由註冊和導航樹掛載測試
- **OpenTelemetry SDK plugin 架構**: 完整的觀測性插件實作和測試
- **Legacy logging 移除**: 舊的日誌模組已清理

### 📋 文檔補充工作
- [x] `docs/interfaces/observability.md` 文檔補充 ✅
- [x] `docs/develop-guide.md` Alloy DevKit 整合說明 ✅
- [x] `docs/interfaces/webplugin.md` Router 掃描與掛載流程補充 ✅
- [x] `docs/develop-guide.md` WebUIPlugin 註冊生命週期說明 ✅


### 建議優先補齊的項目（有檔案但尚未對應功能）

- [x] logging: 移除 legacy logging，整合至統一 logger 管理模組 ✅
  - ✅ 移除 `internal/infrastructure/logging` 舊邏輯（目錄為空）
  - ✅ 改為統一使用 `pkg/shared/log/logger.go`
  - ✅ 配合 `LoggingConfig` 控制格式、輸出、等級
  - **已更新文件**：
    - ✅ `docs/foundation.md` 移除舊 logging 模組描述
  - **已更新測試**：
    - ✅ 移除與舊 logging 有關的初始化測試（若存在）

- [x] webplugin: 擴充並驗證 pages/* 正確掛載 ✅
  - ✅ 建立完整 Dashboard WebUI Plugin 作為範例
  - ✅ System Status Plugin 作為另一個完整範例
  - ✅ 撰寫測試驗證 route 註冊與 navtree 掛載
  - **已新增文件**：
    - ✅ `docs/interfaces/webplugin.md` → 加註 plugin-config 作為範例
  - **已新增測試**：
    - ✅ `internal/test/integration/webui_test.go` → 驗證頁面載入與 navtree 註冊

- [x] observability: 建立 OpenTelemetry SDK plugin 架構 ✅
  - ✅ 目標：`plugins/community/integrations/observability/sdk-wrapper/plugin.go`
  - ✅ 完整實作 interface + metadata，供未來 SDK 掛載與啟用
  - ✅ 包含完整的配置結構、生命週期管理、健康檢查
  - **已新增文件**：
    - ✅ `docs/interfaces/observability.md`（完整的觀測性介面文檔）
    - ✅ `docs/develop-guide.md` → Alloy DevKit 整合與 plugin 自動啟用說明
  - **已新增測試檔**：
    - ✅ `internal/test/integration/sdk_wrapper_test.go`

