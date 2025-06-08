# Logger Interface 說明文件

> 本文件為 Detectviz 專案中 `Logger` 介面的設計說明與使用情境整理。

## 介面用途（What it does）

Logger 是一個結構化日誌抽象介面，負責提供統一的應用層日誌輸出能力。其設計目標包含：

- 支援 info/warn/error/debug 等基本層級
- 結構化欄位輸出以利後續分析
- 可透過 context 傳遞 trace_id / span_id 等追蹤資訊
- 可擴充為 OTLP 相容輸出，對應 Grafana Tempo / Loki 等後端

## 使用情境（When and where it's used）

- 於 `bootstrap.Init()` 時注入核心服務
- 在中介層（middleware）中記錄 HTTP 請求 trace 資訊
- 各模組的背景任務、排程、事件流程中記錄異常與執行過程
- 搭配 `WithFields()` 加入欄位，支援告警與結構化分析

## 方法說明（Methods）

Logger interface 定義如下：

```go
type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
    WithFields(fields map[string]any) Logger
    WithContext(ctx context.Context) Logger
    Named(name string) Logger
    Sync() error
}
```

- `Debug/Info/Warn/Error`：依據 log level 輸出訊息與欄位
- `WithFields`：新增自定欄位後回傳 logger 實例
- `WithContext`：擴充 context（可整合 trace_id / span_id）
- `Named`：新增 logger 名稱（模組分區）
- `Sync`：強制 flush（適用於 zap、buffered logger）

## context 工具（Context Tools）

工具函式定義於 `pkg/ifaces/logger/context.go`：

```go
func WithContext(ctx context.Context, l Logger) context.Context
func FromContext(ctx context.Context) Logger
```

- 可於 middleware 注入 logger 到 context，供後續擷取
- 未注入時預設回傳 `NopLogger`

## 預期實作（Expected implementations）

| 檔案位置                                       | 說明                             |
|------------------------------------------------|----------------------------------|
| `internal/adapters/logger/zap_adapter.go`      | 使用 zap 實作，支援欄位與命名    |
| `internal/adapters/logger/nop_adapter.go`      | 空實作，靜默略過所有輸出         |
| `pkg/ifaces/logger/nop_logger.go`              | NopLogger 結構，做為 fallback     |
| `pkg/ifaces/logger/context.go`                 | context 操作工具函式             |

## 關聯模組與擴充性（Related & extensibility）

- 與 OpenTelemetry trace context 整合，支援 log-trace correlation
- 可結合 Loki / Tempo 透過 OTLP 匯出
- 可 plugin 化 logger backend（stdout、Loki、Redis、file 等）

## 測試建議與驗證方式（Testing & Validation）

Logger 為全域使用的基礎元件，建議進行以下測試：

### 1. 單元測試檔案位置

| 檔案路徑                                      | 測試內容                             |
|-----------------------------------------------|--------------------------------------|
| `internal/adapters/logger/logger_test.go`     | 測試 ZapLogger 輸出是否正確          |
|                                               | 測試 NopLogger 是否靜默處理           |
| `pkg/ifaces/logger/context_test.go`           | 測試 WithContext / FromContext 行為  |

### 2. 核心測試情境

- 能正確將 logger 實例注入並從 context 中擷取
- ZapLogger 能輸出結構化欄位，包含 msg 與 key-value pairs
- NopLogger 調用不會 panic，並正確回傳自身
- 未注入時 FromContext 回傳預設 fallback logger（NopLogger）

### 3. 測試建議技術

- 使用 `zaptest/observer` 觀察輸出行為
- 使用標準 `testing` 套件進行介面一致性驗證