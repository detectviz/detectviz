# DetectViz 專案進度事項

> ✅ 所有已完成項目與範例已移至 [`todo.done.md`](todo.done.md)，
> 提供 Cursor 與 Codex 參考實作格式與邏輯。
> 實作前請先確認是否已有類似完成項目，避免重複開發。

- [Cursor Scaffold 指引補充](docs/README.md)

> 完成了實作項目需標記為完成狀態：✅

---

## scaffold: Logger Plugin 與全域管理模組

- [x] 定義 LoggerProvider interface ✅
  - 路徑：`pkg/platform/contracts/logger.go`
  - 說明：提供平台內部註冊與取得 logger 實例的標準介面
  - 方法建議：
    - ✅ `Logger() contracts.Logger` 統一介面
    - ✅ `Flush() error`（供 shutdown 時釋放資源）
    - ✅ `WithContext(ctx context.Context) contracts.Logger` 做 trace ID 注入
    - ✅ `SetLevel(level string) error` 動態等級變更
    - ✅ `Close() error` 資源釋放

- [x] 建立 Logger plugin 實作範例 ✅
  - 路徑：`plugins/core/logging/otelzap/plugin.go`
  - 說明：實作 contracts.LoggerProvider 並註冊至 Registry
  - 附加實作建議：
    - ✅ 支援 LoggingConfig（輸出型態、等級、格式等）
    - ✅ 注入時設為全域 logger，供 `log.L(ctx)` 呼叫回傳
    - ✅ 使用 `otelzap.New()` 初始化可觀測性日誌系統
    - ✅ 支援檔案輪轉、trace ID 注入、結構化日誌

- [x] 補充 interface 文件 ✅
  - 目標：`docs/interfaces/logger.md`
  - 說明：描述 LoggerProvider 的方法、plugin config 結構與平台引用方式
  - ✅ 完整的配置範例與使用指南
  - ✅ 健康檢查與故障排除說明

- [x] Plugin 化 LoggerProvider 實作 ✅
  - ✅ 建立 `plugins/core/logging/otelzap/` 作為標準日誌插件
    - ✅ 實作 `contracts.LoggerProvider` 接口
    - ✅ 使用 `go.uber.org/zap` 初始化可觀測性日誌系統
    - ✅ 支援 LoggingConfig（包含 log level、輸出格式與路徑）
    - ✅ 說明：已依據 `contracts.LoggerProvider` 介面建立 plugin 並註冊至 Registry，平台其他模組可透過該介面統一取得 logger 實例。
  - ✅ 調整 `pkg/shared/log/logger.go`：
    - ✅ 改為 fallback 預設實作（plugin 尚未載入時使用）
    - ✅ 保留 GetGlobalLogger(), SetGlobalLogger() 等全域函式，但支援 plugin 替代註冊
    - ✅ 提供 LoggerInterface 介面以支援 plugin 適配
  - 整合至 `internal/infrastructure/observability/otelwrapper`：
    - 可注入 logger plugin 實例作為 `LoggerProvider`
    - 初始化時優先透過 Registry 解析註冊的 LoggerPlugin
  - ✅ 撰寫測試：
    - ✅ `otelzap/plugin_test.go` 驗證是否可正確初始化並被掛載使用
    - ✅ 建立適配器(`adapter.go`)確認 plugin 與 shared/log 整合

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

- [x] 掛載 plugin 實例至 scaffold 初始化流程 ✅
  - 路徑：`internal/platform/composition/lifecycle.go`
  - 說明：在 RegisterModules 過程中，自動載入 logging plugin（若存在）
  - ✅ 附加實作：可透過 registry 或 composition 判斷是否啟用 otelzap plugin
  - ✅ 已實作 RegisterLoggerPlugin 方法
  - ✅ 已建立自動註冊機制 (plugins/core/logging/otelzap/init.go)
  - ✅ 已集成到 StartAll 流程中

- [x] 測試 `internal/test/integration/logging_test.go` ✅：
    - ✅ 驗證初始化與 trace context 注入是否正確
    - ✅ 新增 TestOtelZapPluginIntegration 測試
    - ✅ 新增 TestOtelZapTraceContextIntegration 測試
    - ✅ 測試插件註冊、初始化、生命周期管理、健康檢查等功能

---

## scaffold: observability/otelwrapper 模組整合

- [x] 將 `plugins/community/integrations/observability/sdk-wrapper/plugin.go` 移至 `internal/infrastructure/observability/otelwrapper/otelwrapper.go` ✅
  - 重構 `otelwrapper.go`，統一初始化：
    - [x] TracerProvider ✅
    - [x] MeterProvider ✅
    - [x] LoggerProvider（透過 plugin 註冊）✅
    - [x] Resource 設定與 Propagator ✅
  - [x] 移除冗餘設定，合併 context 工具函式與初始化流程 ✅
  - 說明：已建立完整的 OtelWrapper 結構，支援 TracerProvider、MeterProvider、LoggerProvider 整合，並提供全域封裝器實例與便利函式。包含完整的配置結構和初始化流程。

- [x] 移除 `plugins/community/integrations/observability/` 目錄 ✅
  - 預計 alloy 與 otel exporter 相關邏輯移至 deployer 管理，移除目錄跟以下 plugin：
    - [x] `alloy/otlp-receiver/` ✅
    - [x] `tempo-exporter/`, `loki-exporter/`, `mimir-exporter/` ✅
    - [x] `opentelemetry/collector-bridge/`, `auto-instrument/` ✅
  - 說明：已移除整個 observability 插件目錄，相關功能已整合到 otelwrapper 模組中。

- [x] 合併 `internal/infrastructure/{metrics, logging, tracing}` 至 `observability/otelwrapper/` ✅
  - 分析原始封裝內容皆為 OTel 初始化，可統一重構
  - [x] 將三個模組共用的 OTel 資源與 context 傳遞邏輯整合 ✅
  - [x] 調整引用模組：統一透過 `otelwrapper.Tracer()`, `Meter()` 取得對應實例 ✅
  - [x] 提供 `otelwrapper.Shutdown()`，在 lifecycle 終止階段釋放資源 ✅
  - 說明：原始目錄為空，已移除並整合功能到 otelwrapper 模組中。所有 OTel 相關功能現在統一透過 otelwrapper 提供。

- [x] 撰寫對應測試： ✅
  - [x] `otelwrapper_test.go` 驗證三大 provider 是否正常啟動 ✅
  - [x] 測試 traceID、spanID 是否會正確注入 context 並寫入 log ✅
  - 說明：已建立完整的測試套件，包含建立、初始化、全域存取、LoggerProvider 注入等功能測試。

---

## composition: dev-platform 組合重構與管理

- 目前的 DetectViz `compositions/`：
- 是一種「平台 profile 組合包」
- 未來可能透過：
    - GUI 選單列出
    - CLI 顯示描述

- 在每個 compositions/{name}/ 目錄下：
```bash
composition.yaml      # 實際 config，DetectViz 執行時讀取
meta.yaml             # metadata 檔，供 UI / Codex 使用
README.md             # 說明用途與結構
```

#### 使用建議：

- `composition.yaml` 是執行的主體
- `meta.yaml` 是 UI、工具、Codex 用來顯示列表、說明用
* 可在 `apps/server` 的組合選單中整合兩者：
```go
type CompositionBundle struct {
  Meta        MetaDefinition
  Composition CompositionDefinition
}
```
- [x] 將 `minimal-platform/` 改名為 `dev-platform/` ✅
  - 目標：作為 DetectViz 第一版開發平台與 scaffold 動態掛載的預設組合範例來源
  - 變更後路徑：`compositions/dev-platform/`
  - 說明：已成功重命名並更新組合配置

- [x] 移除其他暫未啟用組合平台目錄 ✅
  - 目錄移除：`monitoring-stack/`、`observability-platform/`
  - 說明：已清理未使用的組合平台目錄

- [x] 建立與補充 `composition.yaml` 範例 ✅
  - 路徑：`compositions/dev-platform/composition.yaml`
  - 範例內容：
    ```yaml
    name: dev-platform
    version: 1.0.0
    description: DetectViz 預設開發平台組合
    plugins:
      core_plugins:
        - name: core-auth-jwt
          type: auth
          enabled: true
          config:
            secret: dev-secret
    health:
      plugins:
        include: all
    ```
  - 補充說明：
    - `plugins.core_plugins`/`community_plugins` 依據 plugin 類型分類
    - `health.plugins.include: all` 代表健康檢查涵蓋所有已註冊 plugin
  - 說明：已完成 composition.yaml 重構，包含 core-auth-jwt 和 otelzap-logger 核心插件配置

- [x] health.plugins.include YAML 示意 ✅
  - `composition.yaml` 內 health 欄位：
    ```yaml
    health:
      plugins:
        include: all # 或指定 plugin 名稱陣列
    ```
  - [x] 支援 `all` 或指定 plugin 名稱清單 ✅
  - 說明：已在 composition.yaml 中實現此功能

- [x] 建立 meta.yaml 補充 composition 描述 ✅
  - 路徑：`compositions/dev-platform/meta.yaml`
  - 範例格式：
    ```yaml
    name: dev-platform
    title: DetectViz 開發預設平台
    version: 1.0.0
    description: >
      DetectViz 預設 scaffold 組合平台，適用於本地開發、plugin 驗證與平台整合測試。
    tags: [dev, scaffold, jwt, logger, grafana]
    icon: dev-platform.svg
    ```
  - 說明：已建立完整的 meta.yaml，包含功能特色、系統需求、使用場景等詳細資訊。同時創建了詳細的 README.md 說明文件。

## README.md 與 develop-guide.md 應補充說明的內容清單

  - [x] `README.md`： ✅
    - [x] 組合平台簡介與選擇方式 ✅
    - [x] composition.yaml 與 meta.yaml 結構說明 ✅
    - [x] plugins 區塊與 health.plugins.include 用法 ✅
    - [x] alloy-devkit 路徑與用途 ✅
    - 說明：已在 README.md 中添加完整的組合平台說明章節，包含結構介紹、dev-platform 特色、使用方式、配置範例和健康檢查配置。
  - [x] `docs/develop-guide.md`： ✅
    - [x] dev-platform 作為 scaffold 預設組合平台的角色 ✅
    - [x] composition.yaml 與 meta.yaml 欄位說明 ✅
    - [x] plugin 啟用/停用與 config 覆寫實例 ✅
    - [x] health.plugins.include 支援細節 ✅
    - 說明：已在 develop-guide.md 中添加完整的 dev-platform 說明章節，包含設計目標、核心組件配置、使用場景和啟動方式。
  - [x] 檢查 `README.md` 跟 `docs/develop-guide.md` 內容是否有需要更新或是移除，進行統整重構 ✅
  - 說明：已檢查並更新兩個文件，確保內容一致性和準確性。README.md 已添加組合平台說明，develop-guide.md 已更新 dev-platform 相關內容。

