# Infra Interfaces

本文件定義 detectviz 專案中 `pkg/infra/` 模組的主要介面設計，目的是統一各項平台工具的使用介面，並支援測試注入、mock 與 plugin 替換。

---

## 設計原則

- 所有封裝皆應公開 interface 給核心模組依賴
- 每個 interface 皆可透過 config 或 bootstrap 替換實作
- 不得混入業務邏輯，保持純工具/依賴注入層角色

---

## Logger

```go
type Logger interface {
    Info(args ...any)
    Warn(args ...any)
    Error(args ...any)
    Debug(args ...any)
    With(fields map[string]any) Logger
}
```

> 實作於：`pkg/infra/log/zap_logger.go`

---

## ConfigProvider

```go
type ConfigProvider interface {
    Load() error
    Get(key string) string
    Unmarshal(out any) error
}
```

> 實作於：`pkg/infra/config/provider.go`

---

## Tracer

```go
type Tracer interface {
    Start(ctx context.Context, name string) (context.Context, Span)
}
```

```go
type Span interface {
    End()
    SetTag(key string, value any)
}
```

> 實作於：`pkg/infra/trace/otel_tracer.go`

---

## Cache

```go
type Cache interface {
    Set(ctx context.Context, key string, value any, ttl time.Duration) error
    Get(ctx context.Context, key string) (any, error)
    Delete(ctx context.Context, key string) error
}
```

> 可對應 redis/memory 實作（封裝位置：`pkg/infra/cache/`）

---

## MetricsRegistry

```go
type MetricsRegistry interface {
    RegisterCounter(name string, labels ...string) Counter
    RegisterGauge(name string, labels ...string) Gauge
}
```

```go
type Counter interface {
    Inc(labels map[string]string)
}

type Gauge interface {
    Set(labels map[string]string, value float64)
}
```

> 對應 Prometheus 實作於：`pkg/infra/metrics/prometheus.go`

---

## SignalContext

```go
type SignalContext interface {
    WithSignal(parent context.Context) (context.Context, context.CancelFunc)
}
```

> 封裝 `signal.NotifyContext`，實作於：`pkg/infra/signal/signalctx.go`

---

## HTTPClient

```go
type Doer interface {
    Do(req *http.Request) (*http.Response, error)
}
```

> 實作於：`pkg/infra/httpclient/client.go`

### 延伸建議

- 可封裝 retry, timeout, headers 等設定為 Option 模式
- RoundTripper 可注入 trace, metrics, auth 等中介行為
- 所有模組應使用 `Doer` interface 注入，便於測試與 mock

---

## 擴充建議

- 每個介面應註冊為可由 config/init 注入的元件
- 可搭配 DI 容器或 plugin lifecycle 模組協作