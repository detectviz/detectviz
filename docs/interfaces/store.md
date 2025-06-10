# Store Interface

本文件定義 detectviz 中 `internal/store/` 模組的介面設計標準，並說明多模組、多資料源下的實作擴充規則與 plugin 註冊模式。

---

## 模組概念

store 模組負責抽象不同儲存後端的存取邏輯，每個功能模組（如 rule、notifier、metrics）會有對應的 Store 介面。

---

## 命名與檔案規則

介面定義放於：

```
internal/store/{module}/{backend}/{module}_store.go
```

例如：

- `internal/store/rule/mysql/rule_store.go`
- `internal/store/notifier/logfile/notifier_store.go`

---

## 範例介面定義

```go
type RuleStore interface {
    Create(ctx context.Context, rule *model.Rule) error
    Update(ctx context.Context, rule *model.Rule) error
    Delete(ctx context.Context, id string) error
    GetByID(ctx context.Context, id string) (*model.Rule, error)
    List(ctx context.Context, filter RuleFilter) ([]*model.Rule, error)
}
```

---

## 各模組 Store Interface 定義

### RuleStore

```go
type RuleStore interface {
    Create(ctx context.Context, rule *model.Rule) error
    Update(ctx context.Context, rule *model.Rule) error
    Delete(ctx context.Context, id string) error
    GetByID(ctx context.Context, id string) (*model.Rule, error)
    List(ctx context.Context, filter RuleFilter) ([]*model.Rule, error)
    Enable(ctx context.Context, id string) error
    Disable(ctx context.Context, id string) error
}
```

---

### NotifierStore

```go
type NotifierStore interface {
    Register(ctx context.Context, notif *model.Notifier) error
    Update(ctx context.Context, notif *model.Notifier) error
    Delete(ctx context.Context, id string) error
    GetByID(ctx context.Context, id string) (*model.Notifier, error)
    ListByType(ctx context.Context, typ string) ([]*model.Notifier, error)
    GetActiveByTarget(ctx context.Context, target string) ([]*model.Notifier, error)
}
```

---

### LoggerStore

```go
type LoggerStore interface {
    Append(ctx context.Context, entry *model.LogEntry) error
    Query(ctx context.Context, filter LogQuery) ([]*model.LogEntry, error)
    DeleteBefore(ctx context.Context, timestamp time.Time) error
}
```

---

### MetricsStore

```go
type MetricsStore interface {
    WritePoints(ctx context.Context, points []model.MetricPoint) error
    QueryRange(ctx context.Context, query model.MetricQuery) ([]*model.MetricSeries, error)
    ListMeasurements(ctx context.Context) ([]string, error)
}
```

---

> 其他模組（如 EventBusStore、LoggerStore、MetricsStore）可依照此格式擴充定義。各模組 interface 建議單獨定義於對應模組子目錄 `interfaces.go`。

---

## Plugin 註冊介面

所有實作應註冊至 resolver：

```go
func RegisterRuleStore(source string, impl RuleStore)
```

註冊位置：

```
internal/store/resolver.go
```

資料來源識別建議為：

- `"memory"`, `"mysql"`, `"influxdb"`, `"logfile"`, `"redis"`...

---

## 多模組支援

目前支援的模組：

| 模組      | Interface 名稱      | 預設後端          |
|-----------|---------------------|-------------------|
| rule      | `RuleStore`         | memory / mysql / influxdb / cache |
| notifier  | `NotifierStore`     | logfile / mysql   |
| eventbus  | `EventBusStore`     | memory            |
| logger    | `LoggerStore`       | logfile           |
| metrics   | `MetricsStore`      | influxdb          |

---

## 測試與 mock 建議

每個模組應有 `memory` backend 實作，用於單元測試或 dev 模式。

測試時可透過：

```go
resolver.UseMockRuleStore(memory.New())
```

---

## 擴充建議

- 支援 caching wrapper（ex: `cache.NewWrapper(inner RuleStore)`）
- 提供讀寫分離 decorator
- Streaming store 接口：支援 Grafana Loki、Kafka 等 append-only 系統

---

### EventBusStore

```go
type EventBusStore interface {
    Publish(ctx context.Context, event *model.EventMessage) error
    Subscribe(ctx context.Context, topic string) (<-chan *model.EventMessage, error)
    Unsubscribe(ctx context.Context, topic string) error
}
```