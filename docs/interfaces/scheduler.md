# Scheduler Interface 設計說明

> 本文件說明 `pkg/ifaces/scheduler` 模組中的排程器介面設計與使用方式，並對應 detectviz 專案的實作結構與測試策略。

---

## 設計目標（Design Goals）

Scheduler 模組負責執行背景任務與週期性排程，核心目標包括：

- 抽象化任務調度邏輯
- 支援多種排程策略（如 Cron、Worker Pool）
- 任務具備命名、週期（spec）與執行邏輯
- 支援注入 logger 與 retry 機制
- 易於單元測試與擴充其他排程器實作

---

## Interface 定義（Interface Definitions）

```go
type Job interface {
	Name() string              // 任務名稱
	Spec() string              // 排程規則，如 cron 表達式
	Run(ctx context.Context) error // 執行邏輯
}

type Scheduler interface {
	Register(job Job)                // 註冊任務
	Start(ctx context.Context) error // 啟動排程器
	Stop(ctx context.Context) error  // 停止排程器
}
```

---

## 使用情境（Usage Scenarios）

Scheduler 適用於以下場景：

- 定時資料清理任務
- 定期健康檢查與狀態回報
- 指標彙整與轉拋（例如送出至其他系統）
- 延遲與重試任務（例如任務失敗時重新排程）

---

## 實作位置（Implementations）

| 類型                | 檔案路徑                                                 | 描述                                          |
|---------------------|----------------------------------------------------------|-----------------------------------------------|
| CronScheduler        | `internal/adapters/scheduler/cron_adapter.go`           | 使用 robfig/cron 實現基於 Spec 的排程器       |
| WorkerPoolScheduler  | `internal/adapters/scheduler/workerpool_adapter.go`     | 使用 goroutine pool 實現具備 retry 與 logger |
| MockScheduler        | `internal/adapters/scheduler/mock_adapter.go`           | 測試用，模擬註冊與執行行為                   |

---

## 測試位置（Testing Files）

| 測試檔案路徑                                               | 測試內容                                   |
|------------------------------------------------------------|--------------------------------------------|
| `workerpool_adapter_test.go`                               | 測試排程、併發執行與 retry 行為             |
| `cron_adapter_test.go`                                     | 測試 Cron 表達式排程是否正確執行            |
| `mock_adapter_test.go`                                     | 驗證任務註冊與啟動流程                      |
| `testlogger.go`                                            | 提供測試環境下使用的 logger 實作            |

---

## 実作註冊（Registry）

於 `internal/registry/scheduler/registry.go` 註冊 Scheduler 實作：

```go
func ProvideScheduler(log logger.Logger) scheduler.Scheduler {
	return adapters.NewWorkerPoolScheduler(4, log)
}
```

如需使用其他排程方式，例如 Cron：

```go
return adapters.NewCronScheduler(log)
```

---

## 擴充方式（How to Add a New Scheduler）

若需新增其他排程器類型（如 DelayQueueScheduler）：

1. 建立實作檔案：`internal/adapters/scheduler/{name}_adapter.go`
2. 實作 `Scheduler` 介面三個方法：`Register`、`Start`、`Stop`
3. 可注入 logger，並實作 retry 行為（依需求）
4. 加入單元測試：`{name}_adapter_test.go`
5. 註冊於 `registry.go` 供主程式使用

---

## 關聯模組（Related Modules）

- `logger`：可注入日誌紀錄任務執行與錯誤
- `eventbus`：任務完成後可發出事件
- `config`：可支援未來以動態設定排程行為

---
