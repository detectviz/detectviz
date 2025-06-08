# Scheduler Interface 設計說明

本文件說明 `pkg/ifaces/scheduler` 中的排程器介面定義與使用方式，並對應目前 detectviz 的實作結構。

---

## 設計目標（Design Goals）

Scheduler 模組負責執行背景任務與週期性排程，目標包括：

- 抽象化任務調度機制
- 支援多種排程策略（如 Cron、Worker Pool）
- 任務具備命名與週期定義能力
- 支援 logger、retry 機制
- 易於單元測試與擴充新排程方式

---

## Interface 定義（Interface Definitions）

```go
type Job interface {
	Name() string
	Spec() string
	Run(ctx context.Context) error
}

type Scheduler interface {
	Register(job Job)
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
```

- `Job` 定義任務名稱、週期與執行邏輯。
- `Scheduler` 負責註冊、啟動與停止任務。

---

## 使用情境（Usage）

適用於以下場景：

- 定時資料清理
- 健康檢查回報
- 指標彙整與轉拋任務
- 延遲執行的重試任務佇列

---

## 實作位置（Implementations）

| 類型               | 路徑                                               | 描述                                 |
|--------------------|----------------------------------------------------|--------------------------------------|
| CronScheduler       | `internal/adapters/scheduler/cron_adapter.go`      | 使用 robfig/cron 實現 Spec-based 排程 |
| WorkerPoolScheduler | `internal/adapters/scheduler/workerpool_adapter.go`| goroutine pool 輪詢任務，支援 retry 與 logger |
| MockScheduler       | `internal/adapters/scheduler/mock_adapter.go`      | 測試用，模擬註冊與排程流程             |

---

## 測試位置（Testing Files）

| 檔案路徑                                                  | 測試內容                           |
|-----------------------------------------------------------|------------------------------------|
| `workerpool_adapter_test.go`                              | 測試排程與 retry 行為               |
| `cron_adapter_test.go`                                    | 測試定時排程任務是否依據 Spec 執行 |
| `mock_adapter_test.go`                                    | 測試註冊流程與調度是否正確         |
| `testlogger.go`                                           | 提供共用測試用 logger              |

---

## 實作註冊（Registry）

在 `internal/bootstrap/scheduler_registry.go` 中，可註冊以下實作：

```go
func NewScheduler(log logger.Logger) scheduler.Scheduler {
	return adapters.NewWorkerPoolScheduler(4, log)
}
```

---

## 擴充方式（Extension Tips）

新增其他類型排程器（如事件觸發、Delay Queue）時，僅需實作 `Scheduler` interface：

```go
type MyScheduler struct{}

func (s *MyScheduler) Register(job scheduler.Job) { ... }
func (s *MyScheduler) Start(ctx context.Context) error { ... }
func (s *MyScheduler) Stop(ctx context.Context) error { ... }
```

並加入至 registry 中，即可替換使用。

---
