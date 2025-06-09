

# 核心伺服器模組（server）

`/internal/server` 負責整合設定檔、日誌系統、模組控制與 HTTP Server 啟動，作為 detectviz 系統啟動的主入口點。

---

## 模組結構

| 檔案 | 說明 |
|------|------|
| `server.go` | 定義 `Server` 結構與建構函式，整合各依賴模組 |
| `runner.go` | 實作 `Run()` 與 `Shutdown()`，控制模組啟動與 HTTP Server |
| `instrumentation.go` | 提供 `/metrics`, `/health`, `/debug/pprof/` HTTP 監控端點 |

---

## Server Interface

對應路徑：`pkg/ifaces/server/server.go`

```go
type Server interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
```

- `Run(ctx)`：啟動模組與 HTTP Server，阻塞執行直到 context 結束
- `Shutdown(ctx)`：優雅關閉 HTTP Server 與模組，釋放資源

---

## 使用範例

```go
srv := server.NewServer(cfg, log, engine)
go func() {
  if err := srv.Run(ctx); err != nil {
    log.Error("server run failed", "error", err)
  }
}()

// 接收中斷訊號後關閉
<-signalCtx.Done()
_ = srv.Shutdown(context.Background())
```

---

## 延伸整合

| 元件 | 說明 |
|------|------|
| `bootstrap/init.go` | 注入 Server 所需依賴並建立實例 |
| `modules.ModuleEngine` | 管理模組註冊與 RunAll / ShutdownAll 流程 |
| `config.Provider` | 提供設定值，例如 HTTP Port |
| `logger.Logger` | 日誌輸出，支援結構化 log |

---