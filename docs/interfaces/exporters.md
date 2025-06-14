# Exporter Interface 匯出器介面

> **檔案位置**: `pkg/platform/contracts/exporters.go`

## 概述

`Exporter` 介面定義了資料匯出功能，類似於 Telegraf Output 插件的模式。匯出器負責將處理後的資料發送到外部系統或存儲位置。

## 介面定義

```go
type Exporter interface {
    Plugin
    Export(ctx context.Context, data ExportData) error
    BatchExport(ctx context.Context, batch []ExportData) error
    GetMetrics() ExportMetrics
}
```

## 方法說明

### Export(ctx context.Context, data ExportData) error
- **用途**: 匯出單筆資料
- **參數**: 
  - `ctx` - 上下文，用於取消和超時控制
  - `data` - 要匯出的資料
- **回傳值**: 匯出失敗時回傳錯誤
- **使用場景**: 即時資料處理

### BatchExport(ctx context.Context, batch []ExportData) error
- **用途**: 批次匯出多筆資料
- **參數**:
  - `ctx` - 上下文，用於控制生命週期
  - `batch` - 資料批次陣列
- **回傳值**: 匯出失敗時回傳錯誤
- **使用場景**: 批次處理提升效能

### GetMetrics() ExportMetrics
- **用途**: 取得匯出器效能指標
- **回傳值**: 包含匯出統計的結構體

## 資料結構

```go
type ExportData struct {
    Type      string                 `json:"type"`      // metrics, logs, traces, alerts
    Target    string                 `json:"target"`    // 目標系統識別
    Timestamp time.Time              `json:"timestamp"`
    Data      map[string]interface{} `json:"data"`
    Metadata  map[string]string      `json:"metadata"`
    Priority  int                    `json:"priority"`  // 優先級 (1-10)
}

type ExportMetrics struct {
    TotalExported     int64         `json:"total_exported"`
    TotalErrors       int64         `json:"total_errors"`
    LastExportTime    time.Time     `json:"last_export_time"`
    ExportRate        float64       `json:"export_rate"` // per second
    AverageLatency    time.Duration `json:"average_latency"`
    QueueSize         int           `json:"queue_size"`
}
```

## 配置結構

```go
type ExporterConfig struct {
    Enabled       bool              `yaml:"enabled" json:"enabled"`
    BatchSize     int               `yaml:"batch_size" json:"batch_size"`       // 批次大小
    FlushInterval string            `yaml:"flush_interval" json:"flush_interval"` // 刷新間隔
    Timeout       string            `yaml:"timeout" json:"timeout"`             // 逾時設定
    RetryCount    int               `yaml:"retry_count" json:"retry_count"`     // 重試次數
    RetryDelay    string            `yaml:"retry_delay" json:"retry_delay"`     // 重試延遲
    Target        TargetConfig      `yaml:"target" json:"target"`               // 目標配置
    Transform     TransformConfig   `yaml:"transform" json:"transform"`         // 轉換配置
    Filter        FilterConfig      `yaml:"filter" json:"filter"`               // 過濾配置
}
```

## 實作範例

```go
type InfluxDBExporter struct {
    name        string
    version     string
    description string
    config      *InfluxDBConfig
    client      influxdb2.Client
    writeAPI    api.WriteAPI
    batchQueue  []ExportData
    mutex       sync.Mutex
}

func (e *InfluxDBExporter) Export(ctx context.Context, data ExportData) error {
    // 轉換為 InfluxDB Point
    point, err := e.convertToPoint(data)
    if err != nil {
        return fmt.Errorf("failed to convert data: %w", err)
    }
    
    // 寫入 InfluxDB
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        e.writeAPI.WritePoint(point)
        return nil
    }
}

func (e *InfluxDBExporter) BatchExport(ctx context.Context, batch []ExportData) error {
    points := make([]*write.Point, 0, len(batch))
    
    for _, data := range batch {
        point, err := e.convertToPoint(data)
        if err != nil {
            continue // 記錄錯誤但繼續處理其他資料
        }
        points = append(points, point)
    }
    
    // 批次寫入
    for _, point := range points {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            e.writeAPI.WritePoint(point)
        }
    }
    
    return nil
}
```

## 內建匯出器類型

### 時序資料庫匯出器
- **InfluxDB**: 時序資料庫
- **Prometheus**: Prometheus 時序資料庫
- **TimescaleDB**: PostgreSQL 擴展時序資料庫

### 日誌系統匯出器
- **Elasticsearch**: 全文搜索引擎
- **Loki**: Grafana 日誌聚合系統  
- **Fluentd**: 統一日誌層

### 通知系統匯出器
- **Slack**: 團隊協作通知
- **Email**: 電子郵件通知
- **Webhook**: HTTP 回調通知
- **AlertManager**: Prometheus 告警管理

### 雲端服務匯出器
- **AWS CloudWatch**: AWS 監控服務
- **Azure Monitor**: Azure 監控服務
- **Google Cloud Monitoring**: GCP 監控服務

## 技術棧要求

1. **上下文管理**: 必須支援 `context.Context` 取消和超時
2. **錯誤處理**: 使用有意義的錯誤訊息和錯誤包裝
3. **指標收集**: 實作 `GetMetrics()` 方法提供監控資訊
4. **批次處理**: 支援批次匯出以提升效能
5. **重試機制**: 實作重試邏輯處理臨時失敗

## 最佳實務

1. **連接池管理**: 合理管理連接資源
2. **批次優化**: 根據目標系統特性調整批次大小
3. **錯誤分類**: 區分可重試和不可重試的錯誤
4. **背壓處理**: 處理下游系統的背壓情況
5. **監控整合**: 提供詳細的匯出指標和健康檢查

## 資料流範例

```
Importer → Transform → Filter → Exporter
                                    ↓
                            [Target System]
```

## 配置範例

```yaml
plugins:
  - name: influxdb-exporter
    type: exporter
    enabled: true
    config:
      batch_size: 1000
      flush_interval: "10s"
      timeout: "30s"
      retry_count: 3
      retry_delay: "5s"
      target:
        url: "http://localhost:8086"
        database: "detectviz"
        username: "admin"
        password: "password"
```

## 相關文件

- [Plugin Interface](./plugin.md)
- [Importer Interface](./importers.md)
- [Configuration Schema](../config/schema.md)
- [Plugin Development Guide](../develop-guide.md) 