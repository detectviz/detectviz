# Metric 介面說明（Metric Interfaces Overview）

本文件描述 `pkg/ifaces/metric` 所定義的核心介面，以及對應的實作（Adapter）結構，依據 Clean Architecture 原則設計，所有實作皆置於 `internal/adapters/metrics/` 目錄中。

---

## Interface 一覽

| Interface 名稱             | 說明（中文）                         | 對應實作位置與範例                 |
|----------------------------|--------------------------------------|------------------------------------|
| `MetricQueryAdapter`       | 查詢單一即時值，例如 `cpu_usage`     | `query_adapter.go`：Flux / Prom / Mock |
| `MetricWriter`             | 寫入單一指標資料點                   | `writer_adapter.go`：Influx / Pushgateway / Mock |
| `MetricSeriesReader`       | 查詢時間序列資料                     | `series_reader_adapter.go`：Influx / Mock |
| `MetricTransformer`        | 資料轉換（如單位轉換、標籤處理）     | `transformer_adapter.go`：Noop    |

---

## Aggregator

除了上述 interface，本模組另包含聚合邏輯：

- `SimpleAggregator`：可進行 `sum`、`avg`、`min`、`max` 等基本統計運算。
- 實作位置：`aggregator.go`

---

## 實作位置與結構對應

```
internal/adapters/metrics/
├── aggregator.go                 # 含 SimpleAggregator
├── query_adapter.go             # 含 FluxQueryAdapter / PromQueryAdapter / MockQueryAdapter
├── writer_adapter.go            # 含 InfluxMetricWriter / PushgatewayMetricWriter / MockMetricWriter
├── series_reader_adapter.go     # 含 InfluxSeriesReader / MockSeriesReader
├── transformer_adapter.go       # 含 NoopTransformer
```

---

## 擴充建議

若需支援其他資料來源（如 OpenTSDB、VictoriaMetrics 等）或自定轉換邏輯，可：

- 實作對應 interface 並放入 `internal/adapters/metrics/`
- 保持一個 adapter 僅實作一種行為（查詢、寫入、轉換、聚合）
- 測試邏輯可使用 mock 結構快速驗證

---