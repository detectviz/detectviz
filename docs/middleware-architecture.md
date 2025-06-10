# Middleware Architecture

本文件說明 detectviz 中 `internal/middleware/` 模組的整體架構與責任劃分，包含 middleware 鏈結順序、plugin 擴充點、context 注入方式，以及 observability 的整合方式。

---

## 設計目標

- 統一管理 HTTP 請求進入流程
- 模組化中介邏輯，具備獨立責任與可測試性
- 可透過 plugin 擴充註冊 middleware（如限速器、CORS）
- 與 auth、metrics、logger、tracing 模組整合

---

## 基本原則

每個 middleware 應符合 Go 標準接口：

```go
type Middleware func(next http.Handler) http.Handler
```

middleware 組合順序應遵守：

```go
chain := Recovery(
           Tracing(
             Metrics(
               Logger(
                 Authenticator(handler),
               ),
             ),
           ),
         )
```

---

## 目錄結構建議

```
internal/middleware/
├── auth.go                # JWT / Session 驗證並注入 UserInfo
├── recovery.go            # panic 捕捉並統一錯誤處理
├── logger.go              # 請求與回應日誌記錄
├── metrics.go             # HTTP 指標記錄（status, duration）
├── tracing.go             # OpenTelemetry trace context 傳遞
├── csrf/csrf.go           # CSRF token 檢查與產生
├── cookies/cookies.go     # Cookie 管理（HttpOnly, Secure 設定）
├── requestmeta/request_metadata.go # X-Request-ID, User-Agent 等處理
├── gziper.go              # 回應壓縮處理
├── testing.go             # middleware 測試輔助工具
```

---

## 與其他模組的整合

| 模組           | 整合說明 |
|----------------|----------|
| `auth`         | 調用 `Authenticator` 並將結果注入 context |
| `infra/metrics`| 調用 Prometheus exporter，記錄 HTTP duration |
| `infra/tracing`| 設定 trace ID，並注入至 downstream context |
| `plugins/middleware/` | 提供 plugin 註冊機制，可注入自定義中介層邏輯 |

---

## Plugin 擴充模式

外部 plugin 可透過以下方式註冊：

```go
type MiddlewarePlugin interface {
    ID() string
    Handler() middleware.Middleware
}
```

註冊於 `internal/registry/middleware/registry.go`。

---

## 使用範例：註冊中介層

```go
e := echo.New()
e.Use(middleware.Recovery())
e.Use(middleware.Auth(authManager))
e.Use(middleware.Tracing(tracerProvider))
e.Use(middleware.Metrics(prometheusExporter))
```

---

## 未來擴充計畫

- [ ] 限速器（RateLimiter）支援 plugin 註冊
- [ ] CORS middleware 設定抽象
- [ ] 請求 profiling middleware（整合 pprof 或 tracing label）
- [ ] middleware 組合預設鏈結匯出為 config 可控

---
