# scaffold: 初始目錄與範本建立

- [Cursor Scaffold 指引補充](docs/README.md)

## logger scaffold: Logger 管理模組

- [ ] 導入全域 Logger 管理模組
  - 路徑建議：`pkg/shared/log/logger.go`
  - 功能：
    - 採用 `otelzap` 為唯一 logger，統一平台日誌行為
    - 搭配 lumberjack 支援 log rotation（大小、備份、壓縮）
    - 可透過 `log.L(ctx)` 取得帶有 trace context 的 logger 實例
    - 提供 `log.SetGlobalLogger()` 於 scaffold 啟動時設置預設 logger
  - **需更新文件**：
    - `docs/develop-guide.md` → Logger 初始化與使用方式
    - `docs/interfaces/logger.md`（若定義介面則需新增）
  - **需新增測試**：
    - `internal/test/integration/logging_test.go` 驗證初始化與 trace context

- [ ] 擴充 LoggingConfig 結構
  - 路徑建議：`pkg/config/schema/types.go`
  - 說明：
    - 新增 `Type`, `OTEL`, `FileConfig` 等欄位
    - 支援 log 輸出選項（stdout/file）、格式（json/text）
  - **需更新文件**：
    - `docs/develop-guide.md` → LoggingConfig YAML 範例與欄位說明

- [ ] 升級 plugin 與 scaffold 中所有 `fmt.Printf` 為 `log.L(ctx).Info/Debug(...)`
  - 範圍包含：
    - lifecycle.go
    - prometheus, jwt, influxdb 等 plugin
    - scaffold_test 與 webui_test 中的 log 輸出
  - **需新增測試**：
    - plugin 單元測試或整合測試觀察 log 輸出是否符合格式（於現有測試中補強）

- [ ] 撰寫日誌初始化與整合測試
  - 確認 logger 能從 `LoggingConfig` 正確初始化
  - 驗證 trace ID 是否成功注入至 log 訊息中
  - 測試檔案位置建議：`internal/test/integration/logging_test.go`

---

## scaffold: 尚未實作功能與建議補強項目

### 建議優先補齊的項目（有檔案但尚未對應功能）

- [ ] logging: 移除 legacy logging，整合至統一 logger 管理模組
  - 移除 `internal/infrastructure/logging` 舊邏輯
  - 改為統一使用 `pkg/shared/log/logger.go`
  - 配合 `LoggingConfig` 控制格式、輸出、等級
  - **需更新文件**：
    - `docs/foundation.md` 移除舊 logging 模組描述
  - **需更新測試**：
    - 移除與舊 logging 有關的初始化測試（若存在）

- [ ] webplugin: 擴充並驗證 pages/* 正確掛載
  - 建立最少一個完整 `plugin.go` 作為 `/plugin-config` 範例
  - 其他頁面維持空實作（含 interface method 空函式）
  - 撰寫測試驗證 route 註冊與 navtree 掛載
  - **需新增文件**：
    - `docs/interfaces/webplugin.md` → 加註 plugin-config 作為範例
  - **需新增測試**：
    - `internal/test/integration/webui_test.go` → 驗證頁面載入與 navtree 註冊

- [ ] observability: 建立 OpenTelemetry SDK plugin 架構
  - 目標：`plugins/community/integrations/observability/sdk-wrapper/plugin.go`
  - 初期可空實作 interface + metadata，供未來 SDK 掛載與啟用
  - **需新增文件**：
    - `docs/interfaces/observability.md`（若 interface 有定義）
    - `docs/develop-guide.md` → Alloy DevKit 整合與 plugin 自動啟用說明
  - **可預留測試檔**：
    - `internal/test/integration/sdk_wrapper_test.go`

---

### 建議補強的測試與假件（如進入 plugin 邏輯測試階段）

- [ ] 製作 Plugin Registry 假件
  - 檔案位置：`internal/test/fake/fakeregistry.go`
  - 用於 plugin 單元測試不依賴真實 registry
  - **可搭配測試**：plugin 載入流程、Registry 組合驗證

- [ ] 製作 Plugin Composition 假件
  - 檔案位置：`internal/test/fake/fakecomposition.go`
  - 模擬 config 加載與 plugin 組合流程
  - **可搭配測試**：composition scaffold 驗證

- [ ] 製作 Logger / Context 假件
  - 提供 `fakelogger.go` 與 `fakectx.go`
  - 測試 trace context 傳遞與 log 記錄功能
  - **可搭配測試**：plugin log 行為是否依 ctx 正確輸出 traceID

---

### 可進一步優化的文檔與說明補充

- [ ] docs/interfaces/webplugin.md：補充 Router 掃描與掛載流程
  - 說明 platform 自動載入 WebUIPlugin 的註冊行為
  - **補充內容**：掃描流程說明、plugin 路徑結構範例

- [ ] docs/develop-guide.md：補上 WebUIPlugin 註冊生命週期說明
  - 插件註冊 → 判斷是否為 WebUIPlugin → 呼叫註冊方法
  - **補充內容**：使用時機、典型函式實作範例

- [ ] docs/develop-guide.md：觀測性 SDK 設計補充說明
  - 如何對應 alloy-config.river 自動生成配置與 plugin 初始化
  - **補充內容**：plugin 對應路徑、資料來源結構、觸發條件
