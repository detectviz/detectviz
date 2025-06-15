# Observability Interface

## 概述

DetectViz 的觀測性介面提供了統一的監控、追蹤和日誌收集功能。透過 OpenTelemetry SDK wrapper 插件，平台可以整合各種觀測性工具和服務。

## 核心介面

### HealthChecker 介面

```go
type HealthChecker interface {
    CheckHealth(ctx context.Context) HealthStatus
    GetHealthMetrics() map[string]any
}
```

#### HealthStatus 結構

```go
type HealthStatus struct {
    Status    string    `json:"status"`    // healthy, degraded, unhealthy
    Message   string    `json:"message"`   // 狀態描述
    Timestamp any       `json:"timestamp"` // 檢查時間
    Details   map[string]any `json:"details"` // 詳細指標
}
```

### ObservabilityPlugin 介面

觀測性插件需要實作以下介面：

```go
type ObservabilityPlugin interface {
    Plugin
    HealthChecker
    LifecycleAware
    
    // 配置追蹤
    ConfigureTracing(config TracingConfig) error
    
    // 配置指標收集
    ConfigureMetrics(config MetricsConfig) error
    
    // 配置日誌
    ConfigureLogging(config LoggingConfig) error
    
    // 獲取觀測性資料
    GetObservabilityData() ObservabilityData
}
```

## OpenTelemetry SDK Wrapper

### 插件配置

```yaml
name: "otel-sdk-wrapper"
type: "integration"
category: "observability"
enabled: true
config:
  # 服務資訊
  service:
    name: "detectviz"
    version: "1.0.0"
    environment: "production"
  
  # 追蹤配置
  tracing:
    enabled: true
    endpoint: "http://jaeger:14268/api/traces"
    sampler:
      type: "probabilistic"
      param: 0.1
    headers:
      authorization: "Bearer token"
  
  # 指標配置
  metrics:
    enabled: true
    endpoint: "http://prometheus:9090/api/v1/write"
    interval: "30s"
    headers:
      x-api-key: "metrics-key"
  
  # 日誌配置
  logging:
    enabled: true
    endpoint: "http://loki:3100/loki/api/v1/push"
    level: "info"
    format: "json"
    headers:
      x-scope-orgid: "tenant1"
  
  # 資源屬性
  resource:
    attributes:
      deployment.environment: "production"
      service.namespace: "detectviz"
      service.instance.id: "instance-1"
  
  # 批次處理配置
  batch:
    timeout: "5s"
    max_export_batch_size: 512
    max_queue_size: 2048
```

### 主要功能

#### 1. 追蹤 (Tracing)

- **分散式追蹤**: 支援跨服務的請求追蹤
- **Span 管理**: 自動建立和管理 span
- **上下文傳播**: 在服務間傳播追蹤上下文
- **取樣策略**: 支援多種取樣策略

```go
// 使用範例
ctx := trace.ContextWithSpan(ctx, span)
logger := log.L(ctx) // 自動包含 trace ID
```

#### 2. 指標 (Metrics)

- **系統指標**: CPU、記憶體、磁碟使用率
- **應用指標**: 請求數量、回應時間、錯誤率
- **自定義指標**: 業務相關的指標
- **指標聚合**: 支援計數器、直方圖、量表

```go
// 指標範例
counter := meter.NewCounter("requests_total")
histogram := meter.NewHistogram("request_duration")
gauge := meter.NewGauge("active_connections")
```

#### 3. 日誌 (Logging)

- **結構化日誌**: JSON 格式的日誌輸出
- **日誌等級**: Debug、Info、Warn、Error、Fatal
- **上下文注入**: 自動注入 trace ID 和 span ID
- **日誌轉發**: 支援多種日誌後端

```go
// 日誌範例
logger.Info("Processing request", 
    "user_id", userID,
    "request_id", requestID,
    "duration", duration)
```

## 健康檢查

### 健康狀態

- **healthy**: 所有組件正常運作
- **degraded**: 部分組件有問題但仍可運作
- **unhealthy**: 關鍵組件故障，無法正常運作

### 健康檢查指標

```go
type HealthMetrics struct {
    // 系統指標
    CPUUsage    float64 `json:"cpu_usage"`
    MemoryUsage float64 `json:"memory_usage"`
    DiskUsage   float64 `json:"disk_usage"`
    
    // 網路指標
    NetworkLatency time.Duration `json:"network_latency"`
    NetworkErrors  int64         `json:"network_errors"`
    
    // 應用指標
    ActiveConnections int64 `json:"active_connections"`
    RequestsPerSecond float64 `json:"requests_per_second"`
    ErrorRate         float64 `json:"error_rate"`
    
    // 依賴服務狀態
    Dependencies map[string]string `json:"dependencies"`
}
```

## 配置範例

### 基本配置

```yaml
observability:
  enabled: true
  plugins:
    - name: "otel-sdk-wrapper"
      config:
        service:
          name: "detectviz"
          version: "1.0.0"
        tracing:
          enabled: true
          endpoint: "http://jaeger:14268/api/traces"
        metrics:
          enabled: true
          endpoint: "http://prometheus:9090/api/v1/write"
        logging:
          enabled: true
          level: "info"
```

### 進階配置

```yaml
observability:
  enabled: true
  health_check:
    interval: "30s"
    timeout: "10s"
    endpoints:
      - "/health"
      - "/metrics"
  
  plugins:
    - name: "otel-sdk-wrapper"
      config:
        service:
          name: "detectviz"
          version: "1.0.0"
          environment: "production"
        
        tracing:
          enabled: true
          endpoint: "http://jaeger:14268/api/traces"
          sampler:
            type: "probabilistic"
            param: 0.1
          headers:
            authorization: "Bearer ${JAEGER_TOKEN}"
        
        metrics:
          enabled: true
          endpoint: "http://prometheus:9090/api/v1/write"
          interval: "30s"
          headers:
            x-api-key: "${PROMETHEUS_API_KEY}"
        
        logging:
          enabled: true
          endpoint: "http://loki:3100/loki/api/v1/push"
          level: "info"
          format: "json"
          headers:
            x-scope-orgid: "${LOKI_TENANT_ID}"
        
        resource:
          attributes:
            deployment.environment: "production"
            service.namespace: "detectviz"
            service.instance.id: "${INSTANCE_ID}"
        
        batch:
          timeout: "5s"
          max_export_batch_size: 512
          max_queue_size: 2048
```

## 使用指南

### 1. 啟用觀測性

在主配置檔案中啟用觀測性功能：

```yaml
# config.yaml
observability:
  enabled: true
  plugins:
    - name: "otel-sdk-wrapper"
      enabled: true
```

### 2. 配置追蹤

```go
// 在應用程式中使用追蹤
func handleRequest(ctx context.Context, req *Request) error {
    // 建立 span
    ctx, span := tracer.Start(ctx, "handle_request")
    defer span.End()
    
    // 設置 span 屬性
    span.SetAttributes(
        attribute.String("user.id", req.UserID),
        attribute.String("request.method", req.Method),
    )
    
    // 使用帶有追蹤上下文的 logger
    logger := log.L(ctx)
    logger.Info("Processing request", "request_id", req.ID)
    
    // 處理請求...
    
    return nil
}
```

### 3. 收集指標

```go
// 定義指標
var (
    requestCounter = meter.NewCounter("http_requests_total")
    requestDuration = meter.NewHistogram("http_request_duration_seconds")
)

// 在處理器中使用指標
func handler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    
    // 增加請求計數
    requestCounter.Add(1, attribute.String("method", r.Method))
    
    // 處理請求...
    
    // 記錄請求時間
    duration := time.Since(start).Seconds()
    requestDuration.Record(duration, attribute.String("method", r.Method))
}
```

### 4. 健康檢查

```go
// 實作健康檢查
func (p *ObservabilityPlugin) CheckHealth(ctx context.Context) contracts.HealthStatus {
    // 檢查各組件狀態
    tracingOK := p.checkTracingHealth()
    metricsOK := p.checkMetricsHealth()
    loggingOK := p.checkLoggingHealth()
    
    status := "healthy"
    if !tracingOK || !metricsOK || !loggingOK {
        status = "degraded"
    }
    
    return contracts.HealthStatus{
        Status:  status,
        Message: "Observability components status",
        Details: map[string]any{
            "tracing": tracingOK,
            "metrics": metricsOK,
            "logging": loggingOK,
        },
    }
}
```

## 最佳實踐

### 1. 追蹤最佳實踐

- **適當的 Span 命名**: 使用描述性的 span 名稱
- **設置有意義的屬性**: 包含有助於除錯的資訊
- **控制 Span 數量**: 避免建立過多細粒度的 span
- **處理錯誤**: 在 span 中記錄錯誤資訊

### 2. 指標最佳實踐

- **選擇合適的指標類型**: Counter、Histogram、Gauge
- **使用標籤**: 為指標添加有意義的標籤
- **控制基數**: 避免高基數標籤
- **定期清理**: 清理不再使用的指標

### 3. 日誌最佳實踐

- **結構化日誌**: 使用 JSON 格式
- **適當的日誌等級**: 根據重要性選擇等級
- **包含上下文**: 添加 trace ID 和相關資訊
- **避免敏感資訊**: 不要記錄密碼等敏感資料

### 4. 效能考量

- **取樣策略**: 在生產環境中使用適當的取樣率
- **批次處理**: 使用批次處理減少網路開銷
- **資源限制**: 設置適當的記憶體和 CPU 限制
- **監控開銷**: 監控觀測性工具本身的效能影響

## 故障排除

### 常見問題

1. **追蹤資料遺失**
   - 檢查網路連接
   - 確認端點配置正確
   - 檢查取樣設置

2. **指標不準確**
   - 驗證指標定義
   - 檢查標籤使用
   - 確認聚合設置

3. **日誌格式問題**
   - 檢查日誌格式配置
   - 確認編碼設置
   - 驗證結構化日誌

4. **效能問題**
   - 調整批次大小
   - 優化取樣率
   - 檢查資源使用

### 除錯工具

- **健康檢查端點**: `/health`
- **指標端點**: `/metrics`
- **追蹤檢視**: Jaeger UI
- **日誌查詢**: Grafana/Loki

## 相關文檔

- [開發指南](../develop-guide.md)
- [插件開發](../plugin-development.md)
- [配置參考](../configuration.md)
- [部署指南](../deployment.md) 