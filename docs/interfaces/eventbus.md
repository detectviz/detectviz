# EventDispatcher Interface 說明文件

> 本文件為 Detectviz 專案中事件分派模組的介面與結構說明，包含所有事件型別的 Dispatch 與 Handler 設計。

---

## 介面用途（What it does）

`EventDispatcher` 是模組間透過事件進行解耦通訊的抽象介面，主要用途為：

- 定義所有支援的事件型別與對應的註冊與分派邏輯
- 將事件廣播給所有註冊的 handler
- 支援測試用的 Nop / InMemory 實作

---

## 事件總覽（Supported Events）

| 事件名稱              | 對應介面            | Struct 定義檔案             |
|-----------------------|---------------------|-----------------------------|
| HostDiscoveredEvent   | `HostEventHandler`  | `pkg/ifaces/eventbus/host.go` |
| MetricOverflowEvent   | `MetricEventHandler`| `pkg/ifaces/eventbus/metric.go` |
| AlertTriggeredEvent   | `AlertEventHandler` | `pkg/ifaces/eventbus/alert.go` |
| TaskCompletedEvent    | `TaskEventHandler`  | `pkg/ifaces/eventbus/task.go` |

---

## EventDispatcher 介面定義

### HostDiscoveredEvent

實作位置：`pkg/ifaces/eventbus/host.go`

```go
DispatchHostDiscovered(ctx context.Context, event HostDiscoveredEvent) error
RegisterHostHandler(handler HostEventHandler)
```

- `DispatchHostDiscovered`：發送主機註冊事件
- `RegisterHostHandler`：註冊主機事件的接收者

---

### MetricOverflowEvent

實作位置：`pkg/ifaces/eventbus/metric.go`

```go
DispatchMetricOverflow(ctx context.Context, event MetricOverflowEvent) error
RegisterMetricHandler(handler MetricEventHandler)
```

- `DispatchMetricOverflow`：發送指標過多事件（如分群維度超標）
- `RegisterMetricHandler`：註冊指標事件處理器

---

### AlertTriggeredEvent

實作位置：`pkg/ifaces/eventbus/alert.go`

```go
DispatchAlertTriggered(ctx context.Context, event AlertTriggeredEvent) error
RegisterAlertHandler(handler AlertEventHandler)
```

- `DispatchAlertTriggered`：發送告警觸發事件
- `RegisterAlertHandler`：註冊接收告警事件的處理器

---

### TaskCompletedEvent

實作位置：`pkg/ifaces/eventbus/task.go`

```go
DispatchTaskCompleted(ctx context.Context, event TaskCompletedEvent) error
RegisterTaskHandler(handler TaskEventHandler)
```

- `DispatchTaskCompleted`：任務完成後發送通知
- `RegisterTaskHandler`：註冊接收任務完成事件的模組

---

## Handler 介面定義（Event Handler Interfaces）

| Handler 介面名稱       | 方法名稱                             | 接收的事件結構              |
|------------------------|--------------------------------------|-----------------------------|
| `HostEventHandler`     | `HandleHostDiscovered(ctx, event)`   | `HostDiscoveredEvent`       |
| `MetricEventHandler`   | `HandleMetricOverflow(ctx, event)`   | `MetricOverflowEvent`       |
| `AlertEventHandler`    | `HandleAlertTriggered(ctx, event)`   | `AlertTriggeredEvent`       |
| `TaskEventHandler`     | `HandleTaskCompleted(ctx, event)`    | `TaskCompletedEvent`        |

---

## 事件結構說明（Event Structs）

### HostDiscoveredEvent

| 欄位     | 說明               |
|----------|--------------------|
| Name     | 主機名稱           |
| Labels   | 額外附加標籤資訊   |
| Source   | 資料來源描述字串   |

---

### MetricOverflowEvent

| 欄位     | 說明                     |
|----------|--------------------------|
| Target   | 發生溢出的指標來源       |
| Reason   | 描述溢出的原因（可顯示） |

---

### AlertTriggeredEvent

| 欄位       | 說明                       |
|------------|----------------------------|
| ConditionID | 對應告警條件識別碼       |
| Value      | 觸發時的實際值             |
| Message    | 觸發時的說明訊息           |

---

### TaskCompletedEvent

| 欄位     | 說明               |
|----------|--------------------|
| Name     | 任務名稱           |
| Success  | 是否成功完成       |
| Message  | 額外訊息（如失敗原因） |

---

## EventDispatcher 實作範例（Dispatcher Implementations）

以下為 `EventDispatcher` 介面的具體實作，符合 `pkg/ifaces/eventbus/eventbus.go` 所定義之方法簽章與行為契約。

| 實作檔案                              | 說明                                   |
|---------------------------------------|----------------------------------------|
| `internal/adapters/eventbus/inmemory.go` | 同步事件分派（單元測試與內部模擬用途） |
| `internal/adapters/eventbus/nop.go`      | 空實作，忽略所有事件（禁用或跳過場景） |
| `internal/adapters/eventbus/custom.go`    | 預留擴充用途，可用於注入自定事件邏輯     |

---

### 各事件 Handler 預期實作檔案位置（Planned Handler Implementations）

為使事件接收模組結構清晰，建議將各事件的 `EventHandler` 實作放置於下列位置：

| 事件類型              | Handler Interface        | 建議實作檔案位置                        |
|-----------------------|--------------------------|-----------------------------------------|
| HostDiscoveredEvent   | `HostEventHandler`       | `internal/adapters/eventbus/host.go`    |
| MetricOverflowEvent   | `MetricEventHandler`     | `internal/adapters/eventbus/metric.go`  |
| AlertTriggeredEvent   | `AlertEventHandler`      | `internal/adapters/eventbus/alert.go`   |
| TaskCompletedEvent    | `TaskEventHandler`       | `internal/adapters/eventbus/task.go`    |

上述檔案應包含對應的 `HandleXxx(ctx, event)` 方法實作。

每個實作皆應實作 `pkg/ifaces/eventbus/eventbus.go` 中的 `EventDispatcher` 介面。

---

## 擴充建議（Extensibility）

- 每新增一個事件型別，請建立：
  - 對應的 `Event struct`（如 `XxxEvent`）
  - 對應的 `Handler interface`（如 `XxxEventHandler`）
  - 在 `EventDispatcher` 中定義 `DispatchXxx` 與 `RegisterXxxHandler`
- 建議將每個事件類型放入獨立檔案管理，並在 `eventbus.go` 中用區塊分隔定義。

---
