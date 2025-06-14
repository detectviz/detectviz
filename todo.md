# scaffold: 初始目錄與範本建立

> 備註：原本 detectviz 舊版目錄搬移至 docs/detectviz-deprecated

- [Cursor Scaffold 指引補充](docs/README.md)

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
    - ⚠️ `plugins/community/integrations/security/keycloak/plugin.go` - 檔案不存在
    - ⚠️ `plugins/community/exporters/influxdb/plugin.go` - 檔案不存在

- [x] Plugin interface 與註冊檢查 ✅
  - ✅ 是否實作對應 `contracts.X` interface（Plugin, Importer, LifecycleAware 等）
    - JWT: Plugin + LifecycleAware + HealthChecker + Authenticator
    - Prometheus: Plugin + LifecycleAware + HealthChecker + Importer
    - Gzip: Plugin + LifecycleAware + HealthChecker
    - RequestMeta: Plugin + LifecycleAware + HealthChecker
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

- [ ] WebUIPlugin 掛載測試
  - 使用 `httptest.NewServer()` 驗證 plugin 註冊的 route 是否能正確響應
  - 測試項目：
    - 路由是否註冊
    - NavTree 是否正確註冊節點
    - ComponentRegistry 是否可擴充渲染

- [ ] Plugin Config Schema 驗證測試
  - 檢查錯誤格式的 config 是否正確被 `ValidatePluginConfig()` 阻擋
  - 測試項目：
    - 遺漏欄位 / 類型錯誤 / 無效值
    - Schema 的 required / default 邏輯

- [ ] 日誌整合測試（待日誌升級後進行）
  - plugin 中導入 `otelzap.L()` 或 `logrus.WithContext()` 並觀察輸出是否包含 trace ID
  - 可於 lifecycle 或 Init() 階段記錄 plugin 啟動與結束行為

---

## scaffold: 尚未完成項目補充

- [ ] 補充插件：Keycloak
  - 預定路徑：`plugins/community/integrations/security/keycloak/plugin.go`
  - 支援身份驗證、token 驗證、SSO 整合

- [ ] 補充插件：InfluxDB Exporter
  - 預定路徑：`plugins/community/exporters/influxdb/plugin.go`
  - 將資料推送至 InfluxDB v2/v3，支援 bucket/token 組態

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
