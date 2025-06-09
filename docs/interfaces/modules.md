# 模組生命週期控制接口說明（modules）

模組啟動系統負責註冊、依賴排序、健康監控與整體關閉流程。對齊 Grafana `pkg/modules` 模組化生命週期設計。

---

## 模組總覽

此模組分為五個控制層級：

- **Engine**：集中註冊與執行模組
- **Registry**：管理具名模組（提供註冊/查詢）
- **DependencyGraph**：模組依賴圖與拓撲排序
- **Runner**：根據依賴啟動模組並反向關閉
- **Listener**：定期監控模組健康狀態，異常時觸發停機

---

## Interface 一覽

### `LifecycleModule`

```go
type LifecycleModule interface {
  Run(ctx context.Context) error
  Shutdown(ctx context.Context) error
}
```

模組基本生命週期：Run 啟動、Shutdown 關閉。

---

### `HealthCheckableModule`

```go
type HealthCheckableModule interface {
  LifecycleModule
  Healthy() bool
}
```

可支援健康檢查的模組，配合 `Listener` 使用。

---

### `ModuleEngine`

```go
type ModuleEngine interface {
  Register(m LifecycleModule)
  RunAll(ctx context.Context) error
  ShutdownAll(ctx context.Context) error
}
```

管理匿名模組的註冊與執行。

---

### `ModuleRegistry`

```go
type ModuleRegistry interface {
  Register(name string, m LifecycleModule) error
  Get(name string) (LifecycleModule, bool)
  List() []string
}
```

管理具名模組，可供依賴圖與健康監控查找使用。

---

### `ModuleRunner`

```go
type ModuleRunner interface {
  Start(ctx context.Context) error
  Stop(ctx context.Context) error
}
```

根據 `DependencyGraph` 執行模組啟動與反向關閉流程。

---

### `ModuleListener`

```go
type ModuleListener interface {
  Start(ctx context.Context)
  Stop()
}
```

執行定時健康檢查，如有異常觸發全域關閉。

---

## 使用範例

模組啟動流程建議如下：

```go
engine := modules.NewEngine()
registry := modules.NewRegistry()
graph := modules.NewDependencyGraph()
runner := modules.NewRunner(engine, registry, graph)
listener := modules.NewListener(engine, registry, 5*time.Second)

ctx := context.Background()

listener.Start(ctx)
if err := runner.Start(ctx); err != nil {
  log.Fatal(err)
}
```

---

## 延伸說明

模組系統與其他模組配合如下：

| 模組 | 說明 |
|------|------|
| `bootstrap/init.go` | 整合所有模組 interface 與 wiring |
| `internal/server` | 最終由 server 啟動模組生命週期 |
| `plugins/` | 未來 plugins 也可註冊為模組 |
