# AGENTS.md

DetectViz 是一個可組合的可觀測性平台，採用 Plugin 架構，並透過 Registry、Lifecycle 與 Composition 系統支援擴展與整合。

本文件提供 Codex / Cursor 在開發與驗證過程中所需的架構說明與檢查指引，專注於問題診斷、插件協調與系統整合。

---

## Codex 檢查任務說明

當 Codex 接收到「請根據專案目錄中 `todo.md` 文件列出的項目，進行檢查有沒有確實完成」時，請依下列步驟執行：

1. 逐項解析 `todo.md` 中定義的任務。
2. 檢查對應的：
   - 實作檔案（plugins/, internal/, pkg/）
   - 測試檔案（internal/test/）
   - 說明文件（docs/）
3. 驗證以下核心結構是否齊備：
   - Plugin interface 是否存在於 `pkg/platform/contracts/`
   - Plugin 是否有對應 `plugin.go` 並註冊於 `internal/platform/registry`
   - 若為 WebPlugin，是否已註冊路由與導覽節點
   - 測試檔案是否覆蓋主要 plugin 功能，並可通過 `go test ./...`

### 檢查結果分類

- ✅ 已完成項目（含路徑）
- ⚠️ 部分完成但缺測試/文件
- ❌ 尚未實作項目（僅存在於 `todo.md`）

---

## DetectViz 核心架構

```
detectviz/
├── apps/                 # 應用組合 (server, cli, testkit)
├── pkg/platform/         # Plugin contracts, registry, lifecycle
├── internal/platform/    # 註冊實作、組合、診斷與掃描
├── plugins/              # Plugin 生態系 (core, community, extensions)
├── compositions/         # 實際平台組合與 devkit 建構
```

## Plugin 類型

| 類型         | 介面定義位置                       | 功能說明                         |
|--------------|------------------------------------|----------------------------------|
| Importer     | pkg/platform/contracts/importers.go | 從外部資料來源匯入資料          |
| Exporter     | pkg/platform/contracts/exporters.go | 將資料輸出至外部儲存            |
| Authenticator| pkg/platform/contracts/auth.go      | 身份驗證、SSO、token 驗證       |
| WebPlugin    | pkg/platform/contracts/webplugin.go | 註冊頁面路由與導覽節點          |
| Middleware   | 無明確 interface（透過 echo.Wrap） | 跨模組處理如日誌、追蹤、壓縮    |

每個 Plugin 應提供 `plugin.go` 與註冊函式 `Register(reg contracts.Registry) error`，Plugin 需支援 context 傳遞與結構化日誌輸出。

---

## 系統整合與診斷重點

### Registry

- 插件註冊與掃描：`internal/platform/registry/registry.go`、`discovery.go`
- Plugin metadata 必須包含 `name`, `version`, `type`, `enabled`, `config`

### Lifecycle

- 實作 `LifecycleAware` 與 `HealthChecker` 可啟用 OnStart、Shutdown、HealthCheck 等
- 註冊生命週期見 `pkg/platform/contracts/lifecycle.go`

### Web Plugin 掛載

- Platform 會自動載入實作 `WebUIPlugin` 的 plugin，掛載於 `internal/ports/web/router.go`
- 插件可透過 `RegisterRoutes()` 註冊 API 與畫面入口
- iframe URL 可嵌入外部 Grafana Cloud 或本地頁面

### 日誌與觀測性

- 採用 `otelzap + lumberjack` 管理日誌
- 所有 Plugin 皆應接受 trace context 並寫入結構化 log
- OTEL SDK 支援 trace/span 記錄，註冊方式見 `sdk-wrapper` plugin（空殼可先行）

---

## 整合測試與驗證建議

- 所有插件應有對應整合測試
- 測試覆蓋 Plugin 註冊、路由掛載、資料處理、配置解析
- 建議測試檔案路徑：`internal/test/integration/{plugin}_test.go`

---

## 開發建議與流程（由 Cursor 執行）

1. 根據 `todo.md` 建立對應目錄與 scaffold
2. Plugin 應配合撰寫 `plugin_test.go`（單元）、整合測試（可後補）
3. 註冊至 registry，確保 `Register()` 被掃描與載入
4. 配置需可被 `composition.yaml` 中解析與驗證

---


---

## 其他補充說明（給 Codex 與 Cursor）

- 檢查任務應聚焦於：
  - registry 掛載狀況是否成功
  - plugin 結構與 interface 是否對應完整
  - 整合測試是否有驗證 lifecycle 運作或基本行為
  - plugin config 結構是否正確對應 schema 並可由 composition 載入