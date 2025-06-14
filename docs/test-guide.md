# 測試設計與實作規範

本文件說明 Detectviz 平台的測試設計策略、命名原則、Mock/Fake 管理方式，對齊 Clean Architecture + Plugin 設計，確保可維護性與擴充性。

## 各目錄詳細說明

本專案遵循 Clean Architecture 與 Go 社群常見設計慣例，建議測試檔案與 mock/fake 實作依據責任分工放置在以下位置：

| 類型 | 責任 | 目錄位置 | 適用範例 |
| --- | --- | --- | --- |
| **單元測試** | 測試模組本身的邏輯與行為（純 function or adapter） | **與被測模組同層**（`foo.go` → `foo_test.go`） | logger、notifier、scheduler adapter |
| **整合測試** | 驗證多模組組合（註冊、依賴注入、流程運作） | `/internal/test/` | scheduler + notifier + logger 測通流程 |
| **Fake 實作** | 提供固定邏輯模擬 interface，用於非驗證式測試 | `/internal/test/fakes/` 或 `/test/fakes/` | `fake_scheduler.go`、`fake_notifier.go` |
| **Mock 實作** | 驗證呼叫方法/次數/參數是否正確 | `/internal/test/mocks/` 或 `/pkg/mocks/` | `mock_scheduler.go`（mockery 產出） |
| **共用測試工具** | 建立可複用 logger、clock、context、資料構造等 | `/internal/test/testutil/` | `test_logger.go`、`assert_logger.go`、`fake_clock.go` |

備註：所有 interface 定義皆應集中於 `pkg/platform/contracts/`，原有 `pkg/ifaces/` 已廢除。

### 統整建議

- 模組行為測試 → 模組內
- 整合流程驗證 → internal/test/integration/
- 可共用邏輯模擬 → internal/test/fakes/
- 呼叫驗證測試用 → internal/test/mocks/
- log/config 等測試工具 → internal/test/testutil/

## 測試撰寫原則

- 每個模組皆應具備 `_test.go` 測試檔。
- 單元測試應以 `t.Run(...)` 切分子情境。
- Fake 用於模擬邏輯流程，Mock 用於驗證呼叫與參數。
- 測試涵蓋正常與異常情境（happy path / error path）。
- 禁止在非測試模組內定義 logger 或 config 的 mock 實作。


## 測試目錄結構

```bash
detectviz/
├── internal/
│   ├── adapters/
│   │   └── logger/
│   │       ├── zap_adapter.go
│   │       └── zap_adapter_test.go    # 單元測試：與 adapter 同層
│   ├── platform/
│   │   └── registry/
│   │       └── scheduler/
│   │           ├── registry.go
│   │           └── registry_test.go       # 單元測試（可與 adapter 區隔）
│   ├── test/
│   │   ├── fakes/
│   │   │   └── fake_notifier.go
│   │   ├── mocks/
│   │   │   └── mock_scheduler.go
│   │   ├── testutil/
│   │   │   ├── test_logger.go
│   │   │   ├── assert_logger.go
│   │   │   └── fake_clock.go
│   │   └── integration/
│   │       └── alert_pipeline_test.go
├── pkg/
│   ├── config/
│   │   ├── default.go
│   │   └── default_test.go
│   ├── platform/
│   │   └── contracts/
│   └── mocks/       # mockery 自動產生可選集中放這
```

## 實際範例說明

### `internal/test/fakes/`

- 用途：手動實作的 Fake 類別，符合 interface，但不執行真實邏輯。
- 用於整合測試、模擬流程、不進行斷言。
- 範例：
  - `fake_notifier.go`：模擬發送通知但不實際執行，供 alert 流程測試。
  - `fake_scheduler.go`：模擬任務排程器，用於測試 alert 判斷後的任務注入。

### `internal/test/testutil/`

- 用途：測試輔助函式與共用模組（logger、時間模擬、預設 context 等）。
- 範例：
  - `test_logger.go`：空實作 Logger，適用不需驗證 log 的測試場景。
  - `assert_logger.go`：具 log 記錄功能，可斷言是否有記錄 log 與內容。
  - `fake_clock.go`：固定時間點模擬，便於測試時間觸發邏輯。
  - `test_context.go`：產生標準 context，內含 metadata 標記。

### `internal/test/integration`

- 用途：放置跨模組整合測試（如初始化流程、完整事件觸發測試）。
- 測試真實模組的註冊與協作行為。
- 範例：
  - `registry_init_test.go`：測試 logger/notifier/scheduler 等模組在註冊時是否成功注入。
  - `alert_pipeline_test.go`：模擬告警流程從 metric 觸發 → alert 判斷 → notifier 發送的完整流程。

### `pkg/mocks/`

- 用途：自動產生的 Mock 類別，使用 `mockery` 等工具產生。
- 用於驗證行為、呼叫次數、傳入參數。
- 對應指令範例：
  ```bash
  mockery --name=Notifier --output=./pkg/mocks/
  ```
- 實例用途：
  - `mock_notifier.go`：驗證是否正確呼叫 `.Send()` 並傳入預期 alert message。
  - `mock_scheduler.go`：驗證是否正確註冊任務並執行。
