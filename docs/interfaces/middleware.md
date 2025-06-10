# Middleware Interface

本文件說明 detectviz 的中介層模組介面定義，涵蓋 HTTP 請求處理流程、context 注入、plugin 擴充方式與測試建議。

---

## Interface 定義（Interface Definitions）

本模組定義於 `internal/middleware/`，搭配各模組實作以下 interface：

```go
package middleware

type Middleware func(next http.Handler) http.Handler
```

此外 plugin 擴充時應實作以下介面：

```go
type MiddlewarePlugin interface {
    ID() string
    Handler() Middleware
}
```

註冊位置預設為：
- `internal/registry/middleware/registry.go`：提供註冊與查詢
- 可於 `bootstrap` 或 `apps/*` 中統一掛載

---

## 基本接口定義

每個中介層應實作以下型別：

```go
type Middleware func(next http.Handler) http.Handler
```

此定義符合標準 Go HTTP middleware pattern，可與 Echo 相容。

---

## 主要中介層介面

### Authenticator Middleware

```go
func Auth(manager auth.Authenticator) Middleware
```

- 從 Header / Cookie 擷取 Token
- 呼叫 `Authenticate(ctx, token)`
- 驗證通過後將 `UserInfo` 注入 context，否則回傳 401

### Recovery Middleware

```go
func Recovery() Middleware
```

- 捕捉 panic，輸出日誌並回傳 500 錯誤 JSON 結構

### Logger Middleware

```go
func Logger(logger log.Logger) Middleware
```

- 記錄請求方法、路徑、耗時、回應狀態

### Metrics Middleware

```go
func Metrics(exporter metrics.Exporter) Middleware
```

- 收集 Prometheus metrics（如 duration, code count）

### Tracing Middleware

```go
func Tracing(tp trace.TracerProvider) Middleware
```

- 植入 OpenTelemetry trace context 到每個請求
- 支援與 downstream service trace 串接

---

## Plugin 擴充接口

plugin 可透過註冊以下型別加入自定義 middleware：

```go
type MiddlewarePlugin interface {
  ID() string
  Handler() Middleware
}
```

註冊方式（在 `internal/registry/middleware/registry.go`）：

```go
func Register(m MiddlewarePlugin)
```

---

## Context 與注入內容

所有中介層應支援 context 注入資料，例如：

| 資料類型 | Key 定義                  | 用途                    |
|----------|---------------------------|-------------------------|
| 使用者資訊 | `auth.UserKey`            | 提供後續 handler 使用者身份 |
| trace ID | `traceparent`, `X-Trace-Id`| 用於鏈結 observability |
| request ID | `X-Request-Id`            | 方便除錯與對應日誌     |

---

## 測試建議

- 每個 middleware 應可單元測試與鏈結測試
- 提供測試工具 `middleware/testing.go`，支援模擬 echo.Context 或 http.Request
- 常見測試場景：
  - 未登入 vs 登入權限驗證
  - panic 是否被捕捉
  - 是否產生 trace ID 並 propagate
  - 指標是否正確累加

---

## 參考文件

- [middleware-architecture.md](../middleware-architecture.md)
