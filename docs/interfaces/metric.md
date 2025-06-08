# Metric 介面說明（Metric Interfaces Overview）

> 本文件說明 `pkg/ifaces/metric` 模組內各種與指標（metric）資料相關的抽象介面，並對應實際的 Adapter 實作。設計依據 Clean Architecture 原則，所有具體實作皆置於 `internal/adapters/metrics/`。透過明確分層設計，可支援多種資料來源（Prometheus、InfluxDB、Mock 等）並能彈性擴充。

---

## Interface 一覽

| Interface 名稱             | 說明（中文）                           | 對應實作位置與範例                    |
|----------------------------|----------------------------------------|---------------------------------------|
| `MetricQueryAdapter`       | 查詢單一即時值（如 CPU 使用率）         | `query_adapter.go`：Flux / Prom / Mock |
| `MetricWriter`             | 寫入單一指標資料點（如遞送到 InfluxDB） | `writer_adapter.go`：Influx / Pushgateway / Mock |
| `MetricSeriesReader`       | 查詢時間序列資料（歷史趨勢）            | `series_reader_adapter.go`：Influx / Mock |
| `MetricTransformer`        | 資料轉換（單位轉換、欄位標籤處理）      | `transformer_adapter.go`：Noop       |
| `MetricEvaluator`          | 判斷查詢結果是否符合條件（門檻值比較）  | `evaluator_adapter.go`：ThresholdEvaluator |
| `MetricFetcher`            | 封裝整體查詢、轉換與比對邏輯             | `fetcher_adapter.go`：高階統整查詢流程 |

---

## Aggregator

除了上述 interface，本模組另包含聚合邏輯：

- `SimpleAggregator`：支援 `sum`、`avg`、`min`、`max` 等基本統計運算
- 實作位置：`aggregator.go`

---

## Adapter 結構對應（實作檔案說明）

```
internal/adapters/metrics/
├── aggregator.go                 # 含 SimpleAggregator
├── query_adapter.go             # 含 FluxQueryAdapter / PromQueryAdapter / MockQueryAdapter
├── writer_adapter.go            # 含 InfluxMetricWriter / PushgatewayMetricWriter / MockMetricWriter
├── series_reader_adapter.go     # 含 InfluxSeriesReader / MockSeriesReader
├── transformer_adapter.go       # 含 NoopTransformer
├── evaluator_adapter.go         # 含 ThresholdEvaluator
├── fetcher_adapter.go           # 預期整合多個元件進行查詢與分析（開發中）
```

---

## 擴充建議

若需支援其他資料來源（如 OpenTSDB、VictoriaMetrics、Graphite 等）或自定轉換邏輯，可：

- 實作對應 interface 並放入 `internal/adapters/metrics/`
- 保持一個 adapter 僅實作一種行為（查詢、寫入、轉換、聚合）
- 測試邏輯可使用 Mock adapter 快速驗證（建議放置於同檔案末或 `internal/test`）

---