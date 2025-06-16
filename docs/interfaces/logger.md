# Logger Interface

## 概述

DetectViz 的日誌介面提供了統一的日誌記錄功能。透過 OtelZap 插件，平台可以整合 OpenTelemetry 和 Zap 日誌系統，支援結構化日誌、檔案輪轉、trace ID 注入等功能。

## 核心介面

### LoggerProvider 介面

```go
type LoggerProvider interface {
    // Logger returns the logger instance
    Logger() Logger

    // WithContext returns a logger with context information for trace injection
    WithContext(ctx context.Context) Logger

    // Flush flushes any buffered log entries
    Flush() error

    // SetLevel dynamically changes the log level
    SetLevel(level string) error

    // Close closes the logger and releases resources
    Close() error
}
```

### Logger 介面

```go
type Logger interface {
    // Debug logs a debug message
    Debug(msg string, fields ...interface{})

    // Info logs an info message
    Info(msg string, fields ...interface{})

    // Warn logs a warning message
    Warn(msg string, fields ...interface{})

    // Error logs an error message
    Error(msg string, fields ...interface{})

    // Fatal logs a fatal message and exits
    Fatal(msg string, fields ...interface{})
}
```

### LoggerPlugin 介面

```go
type LoggerPlugin interface {
    Plugin
    LoggerProvider
    HealthChecker
    LifecycleAware
}
```

## OtelZap Logger Plugin

### 插件配置

```yaml
name: "otelzap-logger"
type: "core"
category: "logging"
enabled: true
config:
  enabled: true
  level: "info"                    # debug, info, warn, error, fatal
  format: "json"                   # json, text
  output_type: "both"              # console, file, both
  output: "stdout"                 # stdout, stderr, file path
  service_name: "detectviz"
  service_version: "1.0.0"
  environment: "production"
  
  # 檔案輪轉配置
  file_config:
    filename: "/var/log/detectviz/app.log"
    max_size: 100                  # MB
    max_backups: 3
    max_age: 30                    # days
    compress: true
  
  # OpenTelemetry 配置
  otel:
    enabled: true
    trace_id_field: "trace_id"
    span_id_field: "span_id"
    include_trace: true
    correlation_id: true
  
  # 自定義屬性
  attributes:
    deployment.environment: "production"
    service.namespace: "detectviz"
    team: "platform"
```

### 詳細配置說明

#### 基本配置

- **enabled**: 是否啟用日誌插件
- **level**: 日誌等級，支援 debug、info、warn、error、fatal
- **format**: 日誌格式，支援 json 和 text
- **output_type**: 輸出類型
  - `console`: 僅輸出到控制台
  - `file`: 僅輸出到檔案
  - `both`: 同時輸出到控制台和檔案
- **output**: 控制台輸出目標，支援 stdout、stderr 或檔案路徑
- **service_name**: 服務名稱
- **service_version**: 服務版本
- **environment**: 環境名稱

#### 檔案配置 (file_config)

- **filename**: 日誌檔案路徑
- **max_size**: 單個日誌檔案最大大小（MB）
- **max_backups**: 保留的備份檔案數量
- **max_age**: 檔案保留天數
- **compress**: 是否壓縮舊檔案

#### OpenTelemetry 配置 (otel)

- **enabled**: 是否啟用 OpenTelemetry 整合
- **trace_id_field**: trace ID 欄位名稱
- **span_id_field**: span ID 欄位名稱
- **include_trace**: 是否在日誌中包含 trace 資訊
- **correlation_id**: 是否啟用關聯 ID

#### 自定義屬性 (attributes)

可以新增任意鍵值對作為日誌的額外屬性。

## 使用指南

### 1. 啟用日誌插件

在主配置檔案中啟用日誌功能：

```yaml
# config.yaml
plugins:
  core_plugins:
    - name: "otelzap-logger"
      enabled: true
      config:
        level: "info"
        format: "json"
        output_type: "both"
```

### 2. 在程式碼中使用

```go
import (
    "context"
    "detectviz/pkg/shared/log"
    "detectviz/pkg/platform/contracts"
)

func main() {
    ctx := context.Background()
    
    // 使用全域日誌記錄器
    log.L(ctx).Info("Application starting", "version", "1.0.0")
    
    // 使用帶 context 的日誌記錄器（會自動注入 trace ID）
    logger := log.L(ctx)
    logger.Debug("Debug information", "user_id", 123, "action", "login")
    logger.Info("User logged in", "user_id", 123)
    logger.Warn("Invalid request", "path", "/api/unknown")
    logger.Error("Database connection failed", "error", err)
}
```

### 3. 在插件中使用

```go
// 在插件中取得 logger provider
func (p *MyPlugin) Init(config any) error {
    ctx := context.Background()
    
    // 使用平台提供的 logger
    logger := log.L(ctx)
    logger.Info("Plugin initialized", "plugin", p.Name())
    
    return nil
}
```

### 4. 動態變更日誌等級

```go
import "detectviz/pkg/platform/contracts"

// 透過 registry 取得 logger plugin
loggerPlugin := registry.GetPlugin("otelzap-logger")
if provider, ok := loggerPlugin.(contracts.LoggerProvider); ok {
    err := provider.SetLevel("debug")
    if err != nil {
        log.Error("Failed to set log level", "error", err)
    }
}
```

## 日誌格式範例

### JSON 格式

```json
{
  "timestamp": "2024-01-15T10:30:45.123Z",
  "level": "info",
  "message": "User logged in",
  "caller": "main.go:45",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7",
  "user_id": 123,
  "action": "login",
  "service": {
    "name": "detectviz",
    "version": "1.0.0"
  }
}
```

### Text 格式

```
2024-01-15T10:30:45.123Z INFO main.go:45 User logged in trace_id=4bf92f3577b34da6a3ce929d0e0e4736 span_id=00f067aa0ba902b7 user_id=123 action=login
```

## 健康檢查

日誌插件提供健康檢查功能：

```go
// 檢查插件健康狀態
status := loggerPlugin.CheckHealth(ctx)
fmt.Printf("Status: %s, Message: %s\n", status.Status, status.Message)

// 取得健康指標
metrics := loggerPlugin.GetHealthMetrics()
fmt.Printf("Initialized: %v, Started: %v\n", 
    metrics["plugin_initialized"], 
    metrics["plugin_started"])
```

## 最佳實踐

### 1. 結構化日誌

建議使用結構化欄位而非字串格式化：

```go
// 好的做法
logger.Info("User operation completed", 
    "user_id", userID, 
    "operation", "update_profile",
    "duration_ms", duration.Milliseconds())

// 避免的做法
logger.Info(fmt.Sprintf("User %d completed operation update_profile in %dms", 
    userID, duration.Milliseconds()))
```

### 2. 適當的日誌等級

- **Debug**: 詳細的除錯資訊，僅在開發環境使用
- **Info**: 一般的資訊性訊息，記錄重要的業務流程
- **Warn**: 警告訊息，表示潛在問題但不影響運行
- **Error**: 錯誤訊息，表示發生了錯誤但程式可以繼續
- **Fatal**: 致命錯誤，程式無法繼續運行

### 3. 敏感資訊處理

避免在日誌中記錄敏感資訊：

```go
// 好的做法
logger.Info("User authenticated", "user_id", userID)

// 避免的做法
logger.Info("User authenticated", "password", password, "token", token)
```

### 4. 效能考量

在高頻呼叫的程式碼中，使用適當的日誌等級：

```go
// 檢查日誌等級以避免不必要的處理
if logger.Level() <= log.DebugLevel {
    logger.Debug("Detailed processing info", "data", expensiveOperation())
}
```

## 故障排除

### 常見問題

1. **日誌檔案無法建立**
   - 檢查目錄權限
   - 確認磁碟空間
   - 檢查檔案路徑是否正確

2. **日誌等級不生效**
   - 確認配置檔案語法正確
   - 檢查插件是否正確初始化
   - 驗證動態等級變更是否成功

3. **Trace ID 未顯示**
   - 確認 OpenTelemetry 配置已啟用
   - 檢查 context 是否正確傳遞
   - 驗證 trace provider 是否正確設定

### 除錯方法

```go
// 檢查插件狀態
if plugin, ok := registry.GetPlugin("otelzap-logger"); ok {
    if healthChecker, ok := plugin.(contracts.HealthChecker); ok {
        status := healthChecker.CheckHealth(ctx)
        fmt.Printf("Logger status: %+v\n", status)
    }
}

// 檢查配置
if otelzapPlugin, ok := plugin.(*otelzap.OtelZapPlugin); ok {
    config := otelzapPlugin.GetConfig()
    fmt.Printf("Logger config: %+v\n", config)
}
``` 