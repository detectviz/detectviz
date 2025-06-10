# System Components Interface

本文件定義 detectviz 中 `internal/system/` 區塊的子模組架構與 interface 設計原則。此區域屬於 **cross-cutting technical support modules**，提供平台功能、整合支援與基礎服務，不直接承擔業務邏輯。

---

## 系統支援模組分類與介面說明

### 1. internal/system/http/

- **目的**：處理底層 HTTP 功能，如 proxy, apiserver, cors handler 等。
- **範例介面**：

```go
type CORSProvider interface {
    Wrap(h http.Handler) http.Handler
}
```

---

### 2. internal/system/diagnostics/

- **目的**：系統健康、統計與支援診斷，例如 stats、support bundle。
- **範例介面**：

```go
type HealthChecker interface {
    Report() map[string]string
}
```

---

### 3. internal/system/integration/

- **目的**：與外部系統整合，例如 gRPC server、search API、live updates。
- **範例介面**：

```go
type Integration interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
}
```

---

### 4. internal/system/platform/

- **目的**：通用平台功能，如 quota 控制、全域快取、lifecycle hook。
- **範例介面**：

```go
type QuotaManager interface {
    Check(ctx context.Context, resource string) error
}
```

---

## 設計與擴充原則

- 每個子模組皆可定義對應 interface，並由其他模組透過注入/registry 使用。
- 不應直接依賴 handler 或 store，保持技術底層定位。
- 可由 plugin 或 service 模組依需求調用。