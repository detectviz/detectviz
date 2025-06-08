

# Event Interface 說明文件

> 本文件說明 Detectviz 專案中 `Event` 類型與對應 Handler 的定義與用途。事件模組負責封裝跨模組交換的資料結構，作為 EventBus 分派的核心單位。

---

## 設計目的（Design Purpose）

- 定義各種監控情境下可能產生的事件類型
- 每種事件皆有對應 Handler interface，供 EventBus 呼叫
- 可支援未來擴充更多事件類型（如 UserEvent, ReportEvent）
- 避免硬編碼與不一致資料傳遞格式

---

## Interface 定義（pkg/ifaces/event/*.go）

目前事件定義共分為四類：

### AlertEvent

```go
type AlertTriggeredEvent struct {
	ID      string
	Level   string
	Message string
	Time    time.Time
}

type AlertEventHandler interface {
	HandleAlertEvent(ctx context.Context, event AlertTriggeredEvent) error
}
```

### HostEvent

```go
type HostRegisteredEvent struct {
	InstanceID string
	IP         string
	Time       time.Time
}

type HostEventHandler interface {
	HandleHostEvent(ctx context.Context, event HostRegisteredEvent) error
}
```

### MetricEvent

```go
type MetricOverflowEvent struct {
	Target    string
	Metric    string
	Threshold float64
	Value     float64
	Time      time.Time
}

type MetricEventHandler interface {
	HandleMetricEvent(ctx context.Context, event MetricOverflowEvent) error
}
```

### TaskEvent

```go
type TaskCompletedEvent struct {
	TaskID   string
	Success  bool
	Time     time.Time
	Message  string
}

type TaskEventHandler interface {
	HandleTaskEvent(ctx context.Context, event TaskCompletedEvent) error
}
```

---

## 使用情境（Usage Scenarios）

- 評估告警時產生 `AlertTriggeredEvent` 並通知處理器
- 掃描設備新增時觸發 `HostRegisteredEvent` 記錄並回報
- 發現異常值時觸發 `MetricOverflowEvent` 串聯至告警流程
- 任務排程執行結束時觸發 `TaskCompletedEvent` 進行分析或通知

---

## 測試建議（Testing Strategy）

- 實作 Fake EventBus 驗證每個事件類型是否能正常被觸發與處理
- 使用 `testLogger` 驗證事件是否被正確記錄
- 單元測試每個 Handler 接收到事件後的行為是否符合預期

---

## 擴充建議（Extension Notes）

- 可新增 `UserEvent`, `ReportEvent` 等事件類型
- 可補充每種事件對應 JSON Schema 供外部 API 使用
- 可設計通用 event validator 檢查結構正確性

---