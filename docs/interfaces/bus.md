

# Bus Interface 說明文件

> 本文件說明 Detectviz 專案中 `Bus`（事件總線）模組的介面設計、用途、實作與測試方式。此模組用於統一分派各類型事件（如 Alert、Host、Metric、Task）至對應處理器。

---

## 設計目的（Design Purpose）

- 將事件處理機制從各模組中抽離，統一交由 `Bus` 處理
- 可支援同步或非同步的事件投遞模式
- 降低模組間耦合性，提升擴充性與測試性
- 支援註冊多個處理器與事件分類

---

## Interface 定義（pkg/ifaces/bus/types.go）

```go
type EventDispatcher interface {
    RegisterAlertHandler(handler event.AlertEventHandler)
    RegisterHostHandler(handler event.HostEventHandler)
    RegisterMetricHandler(handler event.MetricEventHandler)
    RegisterTaskHandler(handler event.TaskEventHandler)
    Dispatch(ctx context.Context, e event.Event) error
}
```

---

## 實作位置（Implementations）

| 檔案位置                                         | 描述                         |
|--------------------------------------------------|------------------------------|
| `internal/adapters/eventbus/inmemory.go`         | 預設實作，將事件分派至註冊的 handler |
| `internal/registry/eventbus/registry_inmemory.go`| 註冊 InMemoryDispatcher 為預設 Bus |
| `internal/plugins/eventbus/alertlog/alert_handler.go` | Plugin 實作 AlertEvent 處理器 |

---

## 使用情境（Usage Scenarios）

- 發送 `AlertTriggeredEvent` 通知對應模組產生告警紀錄
- 發送 `MetricOverflowEvent` 進行即時告警評估
- 發送 `TaskCompletedEvent` 通知後續任務系統
- 可用於事件擴充如：Webhook、Kafka、Slack 插件

---

## 註冊與擴充方式（Registration & Extension）

- 可透過 `internal/registry/eventbus/plugins.go` 自動註冊 plugins
- plugin 可實作 `AlertEventHandler` 並透過 `eventbus.RegisterAlertHandler()` 註冊
- 預設可替換為其他實作（如 Kafka、Channel-based）

---

## 測試建議（Testing Strategy）

- 實作 `MockEventDispatcher` 測試發送與註冊行為
- 測試註冊多個 handler 並驗證是否全部觸發
- 搭配 `TestLogger` 驗證事件投遞流程

---

## 關聯模組（Related Modules）

- `event`：定義各類事件的資料結構
- `logger`：事件觸發與失敗時紀錄操作狀態
- `notifier`：最終可連接通知模組

---