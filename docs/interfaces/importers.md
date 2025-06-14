# Importer Interface 匯入器介面

> **檔案位置**: `pkg/platform/contracts/importers.go`

## 概述

`Importer` 介面定義了資料匯入功能，類似於 Telegraf Input 插件的模式。匯入器負責從外部來源收集資料並將其轉換為標準格式。

## 介面定義

```go
type Importer interface {
    Plugin
    Import(ctx context.Context) error
    StartStreaming(ctx context.Context) (<-chan ImportData, error)
    StopStreaming() error
    GetMetrics() ImportMetrics
}
```

## 方法說明

### Import(ctx context.Context) error
- **用途**: 執行一次性資料匯入
- **參數**: `ctx` - 上下文，用於取消和超時控制
- **回傳值**: 匯入失敗時回傳錯誤
- **使用場景**: 定時批次匯入

### StartStreaming(ctx context.Context) (<-chan ImportData, error)
- **用途**: 開始持續資料流匯入
- **參數**: `ctx` - 上下文，用於控制生命週期
- **回傳值**: 資料通道和可能的錯誤
- **使用場景**: 即時資料串流

### StopStreaming() error
- **用途**: 停止資料流匯入
- **回傳值**: 停止失敗時回傳錯誤

### GetMetrics() ImportMetrics
- **用途**: 取得匯入器效能指標
- **回傳值**: 包含匯入統計的結構體

## 資料結構

```go
type ImportData struct {
    Type      string                 `json:"type"`      // metrics, logs, traces
    Source    string                 `json:"source"`    // 資料來源識別
    Timestamp time.Time              `json:"timestamp"`
    Data      map[string]interface{} `json:"data"`
    Metadata  map[string]string      `json:"metadata"`
}

type ImportMetrics struct {
    TotalImported     int64         `json:"total_imported"`
    TotalErrors       int64         `json:"total_errors"`
    LastImportTime    time.Time     `json:"last_import_time"`
    ImportRate        float64       `json:"import_rate"` // per second
    AverageLatency    time.Duration `json:"average_latency"`
}
```

## 配置結構

```go
type ImporterConfig struct {
    Enabled       bool              `yaml:"enabled" json:"enabled"`
    Interval      string            `yaml:"interval" json:"interval"`           // 匯入間隔
    Timeout       string            `yaml:"timeout" json:"timeout"`             // 逾時設定
    BatchSize     int               `yaml:"batch_size" json:"batch_size"`       // 批次大小
    RetryCount    int               `yaml:"retry_count" json:"retry_count"`     // 重試次數
    Source        SourceConfig      `yaml:"source" json:"source"`               // 來源配置
    Transform     TransformConfig   `yaml:"transform" json:"transform"`         // 轉換配置
    Output        OutputConfig      `yaml:"output" json:"output"`               // 輸出配置
}
```

## 實作範例

```go
type PrometheusImporter struct {
    name        string
    version     string  
    description string
    config      *PrometheusConfig
    client      *http.Client
    streaming   bool
    stopCh      chan struct{}
}

func (p *PrometheusImporter) Import(ctx context.Context) error {
    // 獲取 metrics
    metrics, err := p.scrapeMetrics(ctx)
    if err != nil {
        return fmt.Errorf("failed to scrape metrics: %w", err)
    }
    
    // 轉換格式
    importData := p.convertToImportData(metrics)
    
    // 發送到輸出管道
    return p.sendToOutput(ctx, importData)
}

func (p *PrometheusImporter) StartStreaming(ctx context.Context) (<-chan ImportData, error) {
    if p.streaming {
        return nil, fmt.Errorf("streaming already started")
    }
    
    dataCh := make(chan ImportData, 100)
    p.stopCh = make(chan struct{})
    p.streaming = true
    
    go p.streamingLoop(ctx, dataCh)
    
    return dataCh, nil
}
```

## 內建匯入器類型

### 監控指標匯入器
- **Prometheus**: 從 Prometheus 端點獲取指標
- **InfluxDB**: 從 InfluxDB 查詢資料  
- **Telegraf**: 兼容 Telegraf 配置

### 日誌匯入器
- **File**: 檔案日誌匯入
- **Syslog**: 系統日誌匯入
- **HTTP**: HTTP API 日誌匯入

### 分散式追蹤匯入器
- **Jaeger**: Jaeger 追蹤資料
- **Zipkin**: Zipkin 追蹤資料
- **OTLP**: OpenTelemetry 協議

## 技術棧要求

1. **上下文管理**: 必須支援 `context.Context` 取消和超時
2. **錯誤處理**: 使用有意義的錯誤訊息和錯誤包裝
3. **指標收集**: 實作 `GetMetrics()` 方法提供監控資訊
4. **設定解析**: 使用 `mapstructure` 解析配置

## 最佳實務

1. **資源管理**: 正確處理連接和資源清理
2. **錯誤恢復**: 實作重試機制和錯誤恢復
3. **效能監控**: 收集並提供詳細的效能指標
4. **配置驗證**: 驗證必要的配置參數

## 相關文件

- [Plugin Interface](./plugin.md)
- [Exporter Interface](./exporters.md)
- [Configuration Schema](../config/schema.md)
- [Plugin Development Guide](../develop-guide.md) 