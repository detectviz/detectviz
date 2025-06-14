# Prometheus Importer Plugin

Prometheus 匯入器插件用於從 Prometheus 伺服器匯入指標資料。

## 功能特性

- 支援 Prometheus HTTP API 查詢
- 可配置的抓取間隔
- 支援認證（基本認證、Bearer Token）
- 串流和批次匯入模式
- 自動重試機制

## 配置範例

```yaml
name: prometheus-importer
type: importer
config:
  endpoint: "http://localhost:9090"
  scrape_interval: "15s"
  timeout: "10s"
  metrics_path: "/metrics"
  username: "admin"
  password: "password"
  bearer_token: "your-token"
```

## 使用方式

### 註冊插件

```go
import "detectviz/plugins/community/importers/prometheus"

// 註冊到 registry
err := prometheus.Register(registry)

// 註冊到 importer registry
err := prometheus.RegisterImporter(importerRegistry)
```

### 匯入模式

#### 單次匯入
```go
importer, err := registry.GetImporter("prometheus")
err = importer.Import(ctx)
```

#### 串流匯入
```go
dataChan, err := importer.StartStreaming(ctx)
for data := range dataChan {
    // 處理匯入的資料
    processData(data)
}
```

## 支援的資料格式

- Prometheus 指標格式
- OpenMetrics 格式
- 自動轉換為 DetectViz 標準格式

## 依賴項目

- Prometheus Go client library
- HTTP client for API calls

## 配置參數

| 參數 | 類型 | 必填 | 預設值 | 說明 |
|------|------|------|--------|------|
| `endpoint` | string | 是 | - | Prometheus 伺服器端點 |
| `scrape_interval` | duration | 否 | 15s | 抓取間隔 |
| `timeout` | duration | 否 | 10s | 請求超時時間 |
| `metrics_path` | string | 否 | /metrics | 指標路徑 |
| `username` | string | 否 | - | 基本認證使用者名稱 |
| `password` | string | 否 | - | 基本認證密碼 |
| `bearer_token` | string | 否 | - | Bearer token |

## TODO

- [ ] 實作實際的 Prometheus API 客戶端
- [ ] 支援 PromQL 查詢
- [ ] 新增指標過濾功能
- [ ] 支援 TLS 連線
- [ ] 新增連線池管理 