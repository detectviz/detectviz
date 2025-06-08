# Logger 介面設計與使用指南

> 本文件說明 Detectviz 專案中 `Logger` 介面的設計目標、主要功能、典型應用場景、實作方式，以及測試與擴充建議。Logger 為專案核心元件，支援結構化日誌輸出與追蹤整合，便於後續分析、監控與除錯。

---

## 介面功能摘要

Logger 為應用層日誌統一抽象，提供：

- info/warn/error/debug 等多層級日誌
- 結構化欄位輸出，便於檢索與分析
- 支援 context 傳遞 trace_id、span_id 等追蹤資訊
- 可擴充為 OTLP 格式，整合 Loki、Tempo 等觀測系統

---

## 典型應用場景

- `bootstrap.Init()` 注入 logger 與對應 adapter
- Middleware 記錄 HTTP 請求與 trace 資訊
- 背景任務、排程、事件處理等執行流程與錯誤日誌
- 搭配 `WithFields()` 記錄告警、業務欄位等結構化資訊

---

## 介面定義

Logger interface 介面如下：

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

### 方法說明

- `Debug/Info/Warn/Error`：輸出對應層級的訊息
- `WithFields`：加入結構化欄位（如 rule_id、host 等）
- `WithContext`：保留 context 資訊（trace_id 等）
- `Named`：指定 logger 模組名稱
- `Sync`：flush buffer（如 zap 需明確呼叫）

---

## Context 整合工具

定義於 `pkg/ifaces/logger/context.go`：

```go
func WithContext(ctx context.Context, l Logger) context.Context
func FromContext(ctx context.Context) Logger
```

說明：

- 可於 middleware 注入 logger 實例，後續流程可從 context 擷取
- `FromContext` 未注入時會回傳 fallback 的 `NopLogger`

---

## 預期實作與關聯模組

| 檔案位置                                       | 說明                             |
|------------------------------------------------|----------------------------------|
| `internal/adapters/logger/zap_adapter.go`      | zap 套件實作，支援欄位與模組命名 |
| `internal/adapters/logger/nop_adapter.go`      | 空實作，靜默略過所有輸出         |
| `pkg/ifaces/logger/nop_logger.go`              | NopLogger 結構，為 fallback 預設 |
| `pkg/ifaces/logger/context.go`                 | context 操作工具函式             |

### 擴充性與整合建議

- 可整合 OpenTelemetry trace context，實現 log-trace 關聯
- 支援 OTLP 匯出，整合 Loki、Tempo 等觀測後端
- plugin 機制支援自訂 logger backend（file、stdout、Redis 等）

---

## 測試與驗證方式

### 測試檔案

| 檔案路徑                                      | 測試內容                             |
|-----------------------------------------------|--------------------------------------|
| `internal/adapters/logger/logger_test.go`     | 測試 ZapLogger 輸出格式與行為        |
|                                               | 測試 NopLogger 是否靜默處理           |
| `pkg/ifaces/logger/context_test.go`           | 驗證 context 工具函式正確行為         |

### 重要測試場景

- logger 實例可注入並從 context 中擷取
- ZapLogger 可輸出結構化欄位與多級訊息
- NopLogger 可調用且不產生錯誤
- 未注入時，FromContext 回傳 fallback logger

### 建議測試方式

- 使用 `zaptest/observer` 驗證 logger 輸出
- 使用標準 `testing` 驗證 interface 契約與 fallback 行為