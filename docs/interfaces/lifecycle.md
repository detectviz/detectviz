# Lifecycle Interface 生命週期介面

> **檔案位置**: `pkg/platform/contracts/lifecycle.go`

## 概述

`LifecycleAware` 介面定義了插件生命週期管理功能，提供註冊、啟動、停止和關閉等階段的回調方法。`HealthChecker` 介面提供健康檢查功能，用於監控插件運行狀態。

## 介面定義

### LifecycleAware 介面
```go
type LifecycleAware interface {
    OnRegister() error
    OnStart() error
    OnStop() error
    OnShutdown() error
}
```

### HealthChecker 介面
```go
type HealthChecker interface {
    CheckHealth(ctx context.Context) HealthStatus
    GetHealthMetrics() map[string]any
}
```

## 生命週期階段

### 1. OnRegister() error
- **用途**: 插件註冊時呼叫
- **時機**: 插件被註冊到 Registry 時
- **用途**: 執行註冊前的準備工作
- **注意**: 此時插件尚未初始化

### 2. OnStart() error
- **用途**: 插件啟動時呼叫
- **時機**: 系統啟動或插件被動態載入時
- **用途**: 啟動背景服務、建立連接等
- **前提**: 插件必須已初始化 (`Init()` 已呼叫)

### 3. OnStop() error
- **用途**: 插件停止時呼叫
- **時機**: 系統關閉或插件被動態卸載時
- **用途**: 停止背景服務、清理暫存資源
- **注意**: 插件仍保持初始化狀態

### 4. OnShutdown() error
- **用途**: 插件完全關閉時呼叫
- **時機**: 系統最終關閉階段
- **用途**: 完全清理所有資源、關閉連接
- **結果**: 插件回到未初始化狀態

## 生命週期流程

```
未註冊 → [註冊] → 已註冊 → [初始化] → 已初始化 → [啟動] → 運行中
                                                              ↓
已關閉 ← [關閉] ← 已停止 ← [停止] ← 運行中
```

## 健康檢查

### HealthStatus 結構
```go
type HealthStatus struct {
    Status    string                 `json:"status"`    // healthy, unhealthy, degraded
    Message   string                 `json:"message"`
    Timestamp time.Time              `json:"timestamp"`
    Details   map[string]interface{} `json:"details"`
}
```

### CheckHealth(ctx context.Context) HealthStatus
- **用途**: 檢查插件當前健康狀態
- **參數**: `ctx` - 上下文，用於超時控制
- **回傳值**: 健康狀態資訊
- **頻率**: 由系統定期呼叫（預設 30 秒）

### GetHealthMetrics() map[string]any
- **用途**: 取得插件的健康指標數據
- **回傳值**: 包含各種指標的 map
- **用途**: 提供給監控系統使用

## 實作範例

```go
type MyPlugin struct {
    name        string
    version     string
    description string
    config      *MyConfig
    initialized bool
    running     bool
    server      *http.Server
}

// LifecycleAware 實作
func (p *MyPlugin) OnRegister() error {
    // 註冊階段的準備工作
    log.Info().Str("plugin", p.name).Msg("Plugin registered")
    return nil
}

func (p *MyPlugin) OnStart() error {
    if !p.initialized {
        return fmt.Errorf("plugin not initialized")
    }
    
    // 啟動背景服務
    if p.config.EnableHTTPServer {
        p.server = &http.Server{
            Addr:    p.config.HTTPAddress,
            Handler: p.createHandler(),
        }
        
        go func() {
            if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                log.Error().Err(err).Msg("HTTP server error")
            }
        }()
    }
    
    p.running = true
    log.Info().Str("plugin", p.name).Msg("Plugin started")
    return nil
}

func (p *MyPlugin) OnStop() error {
    p.running = false
    
    // 停止 HTTP 服務
    if p.server != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        if err := p.server.Shutdown(ctx); err != nil {
            log.Error().Err(err).Msg("HTTP server shutdown error")
            return err
        }
    }
    
    log.Info().Str("plugin", p.name).Msg("Plugin stopped")
    return nil
}

func (p *MyPlugin) OnShutdown() error {
    // 完全清理資源
    p.running = false
    p.initialized = false
    p.server = nil
    
    log.Info().Str("plugin", p.name).Msg("Plugin shutdown")
    return nil
}

// HealthChecker 實作
func (p *MyPlugin) CheckHealth(ctx context.Context) HealthStatus {
    status := HealthStatus{
        Timestamp: time.Now(),
        Details:   make(map[string]any),
    }
    
    if !p.initialized {
        status.Status = "unhealthy"
        status.Message = "Plugin not initialized"
        return status
    }
    
    if !p.running {
        status.Status = "unhealthy"
        status.Message = "Plugin not running"
        return status
    }
    
    // 檢查服務連接
    if p.server != nil {
        // 簡單的健康檢查
        status.Details["http_server"] = "running"
        status.Details["address"] = p.config.HTTPAddress
    }
    
    status.Status = "healthy"
    status.Message = "Plugin is healthy"
    return status
}

func (p *MyPlugin) GetHealthMetrics() map[string]any {
    return map[string]any{
        "initialized": p.initialized,
        "running":     p.running,
        "uptime":      time.Since(p.startTime).Seconds(),
        "version":     p.version,
    }
}
```

## 生命週期管理器

系統提供 `LifecycleManager` 來統一管理所有插件的生命週期：

```go
type LifecycleManager interface {
    Initialize(ctx context.Context) error
    StartAll(ctx context.Context, registry Registry) error
    StopAll(ctx context.Context, registry Registry) error
    ShutdownAll(ctx context.Context, registry Registry) error
    HealthCheck(ctx context.Context, registry Registry) map[string]HealthStatus
    GetStatus() string
}
```

## 錯誤處理策略

### 註冊階段失敗
- 記錄錯誤但不阻止系統啟動
- 插件標記為不可用
- 繼續處理其他插件

### 啟動階段失敗
- 記錄錯誤並嘗試重試（最多 3 次）
- 如果是關鍵插件，可能需要停止系統
- 非關鍵插件標記為不可用

### 運行時失敗
- 通過健康檢查偵測
- 嘗試重啟插件
- 發送告警通知

### 停止階段失敗
- 記錄錯誤但繼續關閉流程
- 強制終止相關資源
- 確保系統能完整關閉

## 最佳實務

1. **資源管理**:
   - 在 `OnStart()` 中建立資源
   - 在 `OnStop()` 中清理資源
   - 在 `OnShutdown()` 中完全釋放

2. **錯誤處理**:
   - 生命週期方法應該快速返回
   - 使用有意義的錯誤訊息
   - 記錄詳細的操作日誌

3. **健康檢查**:
   - 檢查關鍵資源的可用性
   - 提供詳細的診斷資訊
   - 避免執行耗時操作

4. **併發安全**:
   - 生命週期方法可能並發呼叫
   - 使用適當的同步機制
   - 避免競態條件

## 配置範例

```yaml
lifecycle:
  enabled: true
  timeout: "30s"
  health_check:
    interval: "30s"
    timeout: "10s"
    retries: 3
    failure_threshold: 3
```

## 監控整合

生命週期事件會自動發送到監控系統：

- **指標**: 插件狀態、健康檢查結果、生命週期轉換
- **日誌**: 生命週期事件、錯誤資訊、性能數據
- **告警**: 插件失敗、健康檢查異常、資源用盡

## 相關文件

- [Plugin Interface](./plugin.md)
- [Registry Interface](../platform/registry.md)
- [Monitoring Guide](../monitoring/health-checks.md)
- [Plugin Development Guide](../develop-guide.md) 