

# Services Interface

本文件定義 detectviz 中 `internal/services/` 模組的標準介面，說明每個功能模組的 Service 接口設計原則與跨模組依賴方式。

---

## 設計原則

- 每個 Service interface 應代表一個功能模組的業務邏輯
- 所有 service interface 應可被 mock 以利單元測試
- Service 應依賴抽象的 store 介面與其他 service，不直接使用外部插件

---

## 接口命名建議

- 模組名稱 + `Service` 結尾
- 定義在 `internal/services/{mod}/service.go`
- 跨模組呼叫可集中定義於 `internal/services/interfaces.go`

---

## 各模組 Service Interface 定義

### RuleService

```go
type RuleService interface {
    Create(ctx context.Context, in *model.Rule) error
    Update(ctx context.Context, in *model.Rule) error
    Delete(ctx context.Context, id string) error
    Get(ctx context.Context, id string) (*model.Rule, error)
    List(ctx context.Context, filter RuleFilter) ([]*model.Rule, error)
    Enable(ctx context.Context, id string) error
    Disable(ctx context.Context, id string) error
}
```

---

### NotifierService

```go
type NotifierService interface {
    Register(ctx context.Context, notif *model.Notifier) error
    Send(ctx context.Context, target string, message string) error
    List(ctx context.Context, typ string) ([]*model.Notifier, error)
}
```

---

### LoggerService

```go
type LoggerService interface {
    Append(ctx context.Context, entry *model.LogEntry) error
    Query(ctx context.Context, filter LogQuery) ([]*model.LogEntry, error)
}
```

---

### MetricsService

```go
type MetricsService interface {
    Record(ctx context.Context, point model.MetricPoint) error
    Query(ctx context.Context, query model.MetricQuery) ([]*model.MetricSeries, error)
}
```

---

### EventBusService

```go
type EventBusService interface {
    Emit(ctx context.Context, topic string, msg *model.EventMessage) error
    On(topic string, handler EventHandler)
}
```

---

## 測試建議

- 可定義 fake service 或 mock 介面供 handler 測試注入
- 建議 interface 層不依賴 handler，保持雙向分離

---

## 擴充建議

- 可加入 WithTenant 接口支援多租戶
- 支援 context.Metadata 傳遞，與 middleware 整合使用