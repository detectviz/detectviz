# DetectViz 專案進度事項

- [Cursor Scaffold 指引補充](docs/README.md)

> 完成了實作項目需標記為完成狀態：✅

## logger scaffold: Logger 管理模組

- [x] 導入全域 Logger 管理模組 ✅
  - 路徑建議：`pkg/shared/log/logger.go`
  - 功能：
    - ✅ 採用標準庫實作統一平台日誌行為（避免外部依賴）
    - ✅ 支援 log rotation 配置（大小、備份、壓縮）
    - ✅ 可透過 `log.L(ctx)` 取得帶有 trace context 的 logger 實例
    - ✅ 提供 `log.SetGlobalLogger()` 於 scaffold 啟動時設置預設 logger
  - **已更新文件**：
    - ✅ `docs/develop-guide.md` → Logger 初始化與使用方式
    - ✅ `docs/interfaces/logger.md`（若定義介面則需新增）
  - **已新增測試**：
    - ✅ `internal/test/integration/logging_test.go` 驗證初始化與 trace context

- [x] 擴充 LoggingConfig 結構 ✅
  - 路徑建議：`pkg/config/loader/config_loader.go`
  - 說明：
    - ✅ 新增 `Type`, `OTEL`, `FileConfig` 等欄位
    - ✅ 支援 log 輸出選項（stdout/file）、格式（json/text）
  - **已更新文件**：
    - ✅ `docs/develop-guide.md` → LoggingConfig YAML 範例與欄位說明

- [x] 升級 plugin 與 scaffold 中所有 `fmt.Printf` 為 `log.L(ctx).Info/Debug(...)` ✅
  - 範圍包含：
    - ✅ lifecycle.go
    - ✅ prometheus, jwt, influxdb 等 plugin
    - ✅ scaffold_test 與 webui_test 中的 log 輸出
  - **已新增測試**：
    - ✅ plugin 單元測試或整合測試觀察 log 輸出是否符合格式（於現有測試中補強）

- [x] 撰寫日誌初始化與整合測試 ✅
  - ✅ 確認 logger 能從 `LoggingConfig` 正確初始化
  - ✅ 驗證 trace ID 是否成功注入至 log 訊息中
  - ✅ 測試檔案位置建議：`internal/test/integration/logging_test.go`

---

## scaffold: 尚未實作功能與建議補強項目

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

---

### 建議補強的測試與假件（如進入 plugin 邏輯測試階段）

- [x] 製作 Plugin Registry 假件 ✅
  - 檔案位置：`internal/test/fake/fakeregistry.go`
  - 功能：
    - ✅ 模擬 plugin 註冊、解析、依賴注入
    - ✅ 提供簡化的 `Register()`、`Resolve()` 方法
    - ✅ 可用於 plugin 單元測試與組合驗證
    - ✅ 支援錯誤模擬和統計功能
  - **已新增測試**：
    - ✅ `internal/test/fake/fakeregistry_test.go`

- [x] 製作 Plugin Composition 假件 ✅
  - 檔案位置：`internal/test/fake/fakecomposition.go`
  - 功能：
    - ✅ 提供組合配置與 metadata 模擬行為
    - ✅ 可模擬多 plugin config 情境與 lifecycle 順序
    - ✅ 支援插件初始化、啟動、停止流程
  - **已新增測試**：
    - ✅ `internal/test/fake/fakecomposition_test.go`

- [x] 製作 Logger / Context 假件 ✅
  - 檔案位置：
    - ✅ `internal/test/fake/fakelogger.go`
    - ✅ `internal/test/fake/fakectx.go`
  - 功能：
    - ✅ 提供包含 traceId 的 logger 實例（模擬 `otelzap` logger）
    - ✅ 提供帶有 trace context 的 `context.Context`
    - ✅ 支援日誌等級過濾、條目統計、輸出控制
    - ✅ 支援上下文超時、取消、值傳遞
  - **已新增測試**：
    - ✅ `internal/test/fake/fakelogger_test.go`
    - ✅ `internal/test/fake/fakectx_test.go`

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

# DetectViz 專案進度總結

## 🎉 主要功能實作完成狀態

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

## 🎯 所有項目已完成

DetectViz 專案的所有主要功能實作和文檔補充工作已全部完成！
