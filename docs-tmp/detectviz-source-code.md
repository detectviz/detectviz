Directory structure:
└── detectviz-detectviz/
    ├── README.md
    ├── go.mod
    ├── go.sum
    ├── docs/
    │   ├── app-template-guide.md
    │   ├── architecture-overview.md
    │   ├── coding-style-guide.md
    │   ├── detectviz-apps.md
    │   ├── develop-guide.md
    │   ├── interface-contract.md
    │   ├── plugins-guide.md
    │   ├── interfaces/
    │   │   ├── README.md
    │   │   ├── alert.md
    │   │   ├── bus.md
    │   │   ├── cachestore.md
    │   │   ├── config.md
    │   │   ├── configtypes.md
    │   │   ├── event.md
    │   │   ├── eventbus.md
    │   │   ├── interface-doc-template.md
    │   │   ├── logger.md
    │   │   ├── metric.md
    │   │   ├── modules.md
    │   │   ├── notifier.md
    │   │   ├── scheduler.md
    │   │   └── server.md
    │   └── testing/
    │       └── test-guide.md
    ├── internal/
    │   ├── adapters/
    │   │   ├── alert/
    │   │   │   ├── evaluator.go
    │   │   │   ├── mock_adapter.go
    │   │   │   ├── mock_adapter_test.go
    │   │   │   ├── flux/
    │   │   │   │   ├── flux.go
    │   │   │   │   └── flux_test.go
    │   │   │   └── prom/
    │   │   │       ├── prom.go
    │   │   │       └── prom_test.go
    │   │   ├── cachestore/
    │   │   │   ├── registry.go
    │   │   │   ├── registry_test.go
    │   │   │   ├── memory/
    │   │   │   │   ├── memory.go
    │   │   │   │   └── memory_test.go
    │   │   │   └── redis/
    │   │   │       ├── redis.go
    │   │   │       └── redis_test.go
    │   │   ├── eventbus/
    │   │   │   ├── alert.go
    │   │   │   ├── host.go
    │   │   │   ├── inmemory.go
    │   │   │   ├── metric.go
    │   │   │   └── task.go
    │   │   ├── logger/
    │   │   │   ├── logger_test.go
    │   │   │   ├── nop_adapter.go
    │   │   │   └── zap_adapter.go
    │   │   ├── metrics/
    │   │   │   ├── aggregator.go
    │   │   │   ├── query_adapter.go
    │   │   │   ├── series_reader_adapter.go
    │   │   │   ├── transformer_adapter.go
    │   │   │   └── writer_adapter.go
    │   │   ├── modules/
    │   │   │   ├── engine_adapter.go
    │   │   │   ├── listener_adapter.go
    │   │   │   ├── registry_adapter.go
    │   │   │   └── runner_adapter.go
    │   │   ├── notifier/
    │   │   │   ├── email_adapter.go
    │   │   │   ├── mock_adapter.go
    │   │   │   ├── multi.go
    │   │   │   ├── nop.go
    │   │   │   ├── slack_adapter.go
    │   │   │   └── webhook_adapter.go
    │   │   ├── scheduler/
    │   │   │   ├── cron_adapter.go
    │   │   │   ├── cron_adapter_test.go
    │   │   │   ├── mock_adapter.go
    │   │   │   ├── workerpool_adapter.go
    │   │   │   └── workerpool_adapter_test.go
    │   │   └── server/
    │   │       └── server_adapter.go
    │   ├── alert/
    │   │   └── alert.go
    │   ├── bootstrap/
    │   │   ├── config_loader.go
    │   │   ├── init.go
    │   │   └── wire.go
    │   ├── modules/
    │   │   ├── dependencies.go
    │   │   ├── engine.go
    │   │   ├── listener.go
    │   │   ├── registry.go
    │   │   └── runner.go
    │   ├── plugins/
    │   │   └── eventbus/
    │   │       └── alertlog/
    │   │           ├── alert_handler.go
    │   │           ├── alert_handler_test.go
    │   │           └── init.go
    │   ├── registry/
    │   │   ├── registry.go
    │   │   ├── alert/
    │   │   │   └── registry.go
    │   │   ├── cachestore/
    │   │   │   └── registry.go
    │   │   ├── config/
    │   │   │   └── registry.go
    │   │   ├── eventbus/
    │   │   │   ├── plugins.go
    │   │   │   ├── providers.go
    │   │   │   ├── registry.go
    │   │   │   └── registry_inmemory.go
    │   │   ├── logger/
    │   │   │   └── registry.go
    │   │   ├── notifier/
    │   │   │   ├── registry.go
    │   │   │   └── registry_test.go
    │   │   └── scheduler/
    │   │       └── registry.go
    │   ├── server/
    │   │   ├── instrumentation.go
    │   │   ├── runner.go
    │   │   └── server.go
    │   └── test/
    │       ├── README.md
    │       ├── fakes/
    │       │   ├── fake_config.go
    │       │   ├── fake_metrics.go
    │       │   ├── fake_modules.go
    │       │   └── fake_server.go
    │       ├── plugins/
    │       │   └── eventbus/
    │       │       └── alertlog/
    │       │           └── alert_plugin_test.go
    │       ├── server/
    │       │   └── server_test.go
    │       └── testutil/
    │           ├── assert_logger.go
    │           └── test_logger.go
    └── pkg/
        ├── config/
        │   ├── README.md
        │   └── default.go
        ├── configtypes/
        │   ├── cache_config.go
        │   └── notifier_config.go
        └── ifaces/
            ├── alert/
            │   ├── evaluate.go
            │   ├── evaluate_test.go
            │   └── evaluator.go
            ├── bus/
            │   ├── alert.go
            │   ├── host.go
            │   ├── metric.go
            │   ├── task.go
            │   └── types.go
            ├── cachestore/
            │   └── cachestore.go
            ├── config/
            │   └── config.go
            ├── event/
            │   ├── alert.go
            │   ├── host.go
            │   ├── metric.go
            │   ├── task.go
            │   └── types.go
            ├── eventbus/
            │   ├── eventbus.go
            │   └── provider.go
            ├── logger/
            │   ├── context.go
            │   ├── context_test.go
            │   ├── logger.go
            │   └── nop_logger.go
            ├── metrics/
            │   ├── metric.go
            │   ├── query.go
            │   └── types.go
            ├── modules/
            │   └── modules.go
            ├── notifier/
            │   └── notifier.go
            ├── scheduler/
            │   ├── mock_adapter_test.go
            │   └── scheduler.go
            └── server/
                └── server.go

================================================
FILE: README.md
================================================
# Detectviz

Detectviz 是一套基於 Clean Architecture 設計的模組化監控與告警平台，支援指標查詢、條件比對、事件發布與通知處理。透過 Plugin 機制整合各種數據來源（如 Prometheus, InfluxDB, Flux）與通知通道（如 Email, Slack, Webhook），並提供可維護、可擴充的事件處理架構。

---

## 專案目錄結構

```
detectviz/
├── apps/                     # 獨立 App 模組（如 alert-app、web-app 等）
├── cmd/                      # CLI 主程式進入點（可選）
├── internal/                 # 核心邏輯模組（僅供 apps 使用）
│   ├── adapters/             # 各模組實作（logger, notifier, scheduler 等）
│   ├── registry/             # 模組註冊中心（alert, notifier, scheduler 等）
│   └── test/                 # 整合測試、fakes、mocks、testutil 工具
├── pkg/                      # 共用抽象（interface、config、domain）
│   ├── config/               # 設定載入與注入模組
│   ├── ifaces/               # 模組介面定義（Logger, Scheduler, Notifier 等）
│   └── mocks/                # 使用 mockery 產出的 mock interface（自動生成）
├── plugins/                  # 插件模組（資料查詢來源、通知管道擴充等）
├── scripts/                  # 輔助腳本（備份、啟動、模擬工具）
├── deploy/                   # Docker 與環境部署相關設定
└── README.md
```

---

## 已實作模組

- **Logger**：支援 Zap 實作與 NopLogger。
- **ConfigProvider**：統一提供全域設定注入。
- **EventBus**：可註冊多種事件處理器（Host, Metric, Alert, Task）。
- **AlertEvaluator**：支援 Prometheus、Flux 查詢條件擴充。
- **Scheduler**：支援 Cron 與 WorkerPool 型任務排程。
- **Notifier**：支援 Email、Slack、Webhook 多種通道。

---

## 啟動方式

請搭配各 `apps/` 內主程式使用 `go run` 或 `make` 指令：

```bash
go run ./apps/alert-app/main.go
make run-scheduler
```

可參考 `scripts/` 或 `Makefile` 中的啟動流程與模擬指令。

---

## 文件導引

- [docs/interfaces/](docs/interfaces/)：介面定義與實作契約說明
- [internal/registry/](internal/registry/)：模組註冊流程（AlertEvaluator、Notifier、Scheduler 等）
- [internal/test/README.md](internal/test/README.md)：測試策略與實際目錄規劃
- [docs/develop-guide.md](docs/develop-guide.md)：設計原則與架構圖
- [docs/coding-style-guide.md](docs/coding-style-guide.md)：程式撰寫風格（命名規則、註解格式、golangci-lint 設定）



================================================
FILE: go.mod
================================================
module github.com/detectviz/detectviz

go 1.24.3

require (
	github.com/robfig/cron/v3 v3.0.1 // @grafana/grafana-backend-group
	github.com/stretchr/testify v1.10.0 // @grafana/grafana-backend-group
	go.uber.org/zap v1.27.0 // @grafana/identity-access-team
	gopkg.in/yaml.v3 v3.0.1 // indirect; @grafana/alerting-backend
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
)

require (
	github.com/kr/pretty v0.3.1 // indirect
	github.com/redis/go-redis/v9 v9.10.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

// Use fork of crewjam/saml with fixes for some issues until changes get merged into upstream
replace github.com/crewjam/saml => github.com/grafana/saml v0.4.15-0.20240917091248-ae3bbdad8a56

// Use our fork of the upstream alertmanagers.
// This is required in order to get notification delivery errors from the receivers API.
replace github.com/prometheus/alertmanager => github.com/grafana/prometheus-alertmanager v0.25.1-0.20250417181314-6d0f5436a1fb

exclude github.com/mattn/go-sqlite3 v2.0.3+incompatible

// lock for mysql tsdb compat
replace github.com/go-sql-driver/mysql => github.com/go-sql-driver/mysql v1.7.1

// v1.* versions were retracted, we need to stick with v0.*. This should work
// without the exclude, but this otherwise gets pulled in as a transitive
// dependency.
exclude github.com/prometheus/prometheus v1.8.2-0.20221021121301-51a44e6657c3

// This was retracted, but seems to be known by the Go module proxy, and is
// otherwise pulled in as a transitive dependency.
exclude k8s.io/client-go v12.0.0+incompatible



================================================
FILE: go.sum
================================================
github.com/bsm/ginkgo/v2 v2.12.0 h1:Ny8MWAHyOepLGlLKYmXG4IEkioBysk6GpaRTLC8zwWs=
github.com/bsm/ginkgo/v2 v2.12.0/go.mod h1:SwYbGRRDovPVboqFv0tPTcG1sN61LM1Z4ARdbAV9g4c=
github.com/bsm/gomega v1.27.10 h1:yeMWxP2pV2fG3FgAODIY8EiRE3dy0aeFYt4l7wh6yKA=
github.com/bsm/gomega v1.27.10/go.mod h1:JyEr/xRbxbtgWNi8tIEVPUYZ5Dzef52k01W3YH0H+O0=
github.com/cespare/xxhash/v2 v2.3.0 h1:UL815xU9SqsFlibzuggzjXhog7bL6oX9BbNZnL2UFvs=
github.com/cespare/xxhash/v2 v2.3.0/go.mod h1:VGX0DQ3Q6kWi7AoAeZDth3/j3BFtOZR5XLFGgcrjCOs=
github.com/creack/pty v1.1.9/go.mod h1:oKZEueFk5CKHvIhNR5MUki03XCEU+Q6VDXinZuGJ33E=
github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc h1:U9qPSI2PIWSS1VwoXQT9A3Wy9MM3WgvqSxFWenqJduM=
github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f h1:lO4WD4F/rVNCu3HqELle0jiPLLBs70cWOduZpkS1E78=
github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f/go.mod h1:cuUVRXasLTGF7a8hSLbxyZXjz+1KgoB3wDUb6vlszIc=
github.com/kr/pretty v0.2.1/go.mod h1:ipq/a2n7PKx3OHsz4KJII5eveXtPO4qwEXGdVfWzfnI=
github.com/kr/pretty v0.3.1 h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=
github.com/kr/pretty v0.3.1/go.mod h1:hoEshYVHaxMs3cyo3Yncou5ZscifuDolrwPKZanG3xk=
github.com/kr/pty v1.1.1/go.mod h1:pFQYn66WHrOpPYNljwOMqo10TkYh1fy3cYio2l3bCsQ=
github.com/kr/text v0.1.0/go.mod h1:4Jbv+DJW3UT/LiOwJeYQe1efqtUx/iVham/4vfdArNI=
github.com/kr/text v0.2.0 h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=
github.com/kr/text v0.2.0/go.mod h1:eLer722TekiGuMkidMxC/pM04lWEeraHUUmBw8l2grE=
github.com/pkg/diff v0.0.0-20210226163009-20ebb0f2a09e/go.mod h1:pJLUxLENpZxwdsKMEsNbx1VGcRFpLqf3715MtcvvzbA=
github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 h1:Jamvg5psRIccs7FGNTlIRMkT8wgtp5eCXdBlqhYGL6U=
github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/redis/go-redis/v9 v9.10.0 h1:FxwK3eV8p/CQa0Ch276C7u2d0eNC9kCmAYQ7mCXCzVs=
github.com/redis/go-redis/v9 v9.10.0/go.mod h1:huWgSWd8mW6+m0VPhJjSSQ+d6Nh1VICQ6Q5lHuCH/Iw=
github.com/robfig/cron/v3 v3.0.1 h1:WdRxkvbJztn8LMz/QEvLN5sBU+xKpSqwwUO1Pjr4qDs=
github.com/robfig/cron/v3 v3.0.1/go.mod h1:eQICP3HwyT7UooqI/z+Ov+PtYAWygg1TEWWzGIFLtro=
github.com/rogpeppe/go-internal v1.9.0/go.mod h1:WtVeX8xhTBvf0smdhujwtBcq4Qrzq/fJaraNFVN+nFs=
github.com/rogpeppe/go-internal v1.14.1 h1:UQB4HGPB6osV0SQTLymcB4TgvyWu6ZyliaW0tI/otEQ=
github.com/rogpeppe/go-internal v1.14.1/go.mod h1:MaRKkUm5W0goXpeCfT7UZI6fk/L7L7so1lCWt35ZSgc=
github.com/stretchr/testify v1.10.0 h1:Xv5erBjTwe/5IxqUQTdXv5kgmIvbHo3QQyRwhJsOfJA=
github.com/stretchr/testify v1.10.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
go.uber.org/goleak v1.3.0 h1:2K3zAYmnTNqV73imy9J1T3WC+gmCePx2hEGkimedGto=
go.uber.org/goleak v1.3.0/go.mod h1:CoHD4mav9JJNrW/WLlf7HGZPjdw8EucARQHekz1X6bE=
go.uber.org/multierr v1.11.0 h1:blXXJkSxSSfBVBlC76pxqeO+LN3aDfLQo+309xJstO0=
go.uber.org/multierr v1.11.0/go.mod h1:20+QtiLqy0Nd6FdQB9TLXag12DsQkrbs3htMFfDN80Y=
go.uber.org/zap v1.27.0 h1:aJMhYGrd5QSmlpLMr2MftRKl7t8J8PTZPA732ud/XR8=
go.uber.org/zap v1.27.0/go.mod h1:GB2qFLM7cTU87MWRP2mPIjqfIDnGu+VIO4V/SdhGo2E=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=
gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c/go.mod h1:JHkPIbrfpd72SG/EVd6muEfDQjcINNoR0C8j2r3qZ4Q=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=



================================================
FILE: docs/app-template-guide.md
================================================



================================================
FILE: docs/architecture-overview.md
================================================



================================================
FILE: docs/coding-style-guide.md
================================================
# Detectviz Coding Style Guide

本文件統一規範 Detectviz 專案的程式撰寫風格，包括：

- 中英文註解格式與範例
- Interface 與 Adapter 命名原則
- Lint 與靜態分析工具建議（golangci-lint）
- 測試與模擬實作的撰寫建議

---

## 註解風格與語言標準

Detectviz 採用中英文對照註解格式，目的是提升可讀性、跨語系協作能力，並支援 IDE 自動產生文件說明。

### 一般原則

- 每個公開的 struct、interface、function 必須包含英文 GoDoc 註解。
- 若有額外中文補充，應以 `zh:` 為前綴，放在同一段或下一行。
- 若為欄位註解，可與變數置於同行。

### 範例：結構與方法

```go
// InfluxSeriesReader implements MetricSeriesReader by querying InfluxDB.
// zh: InfluxSeriesReader 透過查詢 InfluxDB 實作 MetricSeriesReader，用於取得時間序列資料。
type InfluxSeriesReader struct {
    Client influxdb2.QueryAPI // zh: InfluxDB 查詢客戶端（需先透過 config 注入）
}

// ReadSeries retrieves metric data using a Flux query and returns parsed points.
// zh: 透過 Flux 查詢取得指標資料，並回傳解析後的時間點清單。
func (r *InfluxSeriesReader) ReadSeries(ctx context.Context, req *metric.ReadRequest) ([]metric.TimePoint, error) {
    ...
}
```

### 範例：欄位說明

```go
type MockSeriesReader struct {
    Points []metric.TimePoint // zh: 預設回傳的資料點清單
    Err    error              // zh: 若設定錯誤，則每次查詢都會回傳此錯誤
}
```

### TODO 與測試相關註解

```go
// TODO: Inject InfluxDB client via config.Provider
// zh: 待透過組態注入 InfluxDB 查詢元件（influxdb2.QueryAPI）

// TODO: Support dynamic rule evaluation and fallback logic
// zh: 預計支援動態告警規則與備援條件比對機制
```

```go
// testLogger 是簡化用於測試的 logger 實作。
// zh: 用於測試過程中收集 log 輸出內容的模擬 Logger。
type testLogger struct {
    logs []string // zh: 收集輸出的記錄訊息
}
```

### 附錄：使用建議

- 所有可導出的元件皆應撰寫英文註解，並建議補上 `zh:` 中譯。
- 測試用的 fake/mock struct 建議清楚標示用途，例如 `// 用於整合測試，不含驗證邏輯`
- 結構、方法、欄位皆應保持註解風格一致，利於維護與自動文件生成。

---

## 介面與實作命名原則

Detectviz 採用清楚且語意一致的命名規則來區分 interface（抽象）與 adapter（具體實作）。

### 命名原則

- 所有 interface 命名採用具名描述詞結尾，例如：`Logger`、`Scheduler`、`MetricWriter`
- interface 實作（adapter）命名需加入語意前綴或後綴，例如：
  - 前綴：`MockLogger`, `InfluxSeriesReader`, `SlackNotifier`
  - 後綴：`DefaultAlertEvaluator`, `CronScheduler`, `PushgatewayMetricWriter`
- 若為多介面組合的實作，可保留語意一致性，但避免冗長命名
- 建議 interface 檔案命名為：`<功能名>.go`（如 `logger.go`, `scheduler.go`）
- 建議 adapter 檔案命名為：`<類型>_adapter.go` 或 `<具體功能>.go`（如 `zap_adapter.go`, `cron_adapter.go`）

### 命名範例

| 類別       | 命名                        | 說明                         |
|------------|-----------------------------|------------------------------|
| interface  | `Scheduler`                 | 定義任務排程器功能           |
| adapter    | `CronScheduler`             | 使用 cron 實作的 Scheduler  |
| interface  | `Notifier`                  | 定義通知推送介面             |
| adapter    | `SlackNotifier`             | 使用 Slack 實作的 Notifier  |
| adapter    | `MockNotifier`              | 用於測試的假通知實作        |
| interface  | `SeriesReader`              | 定義時間序列讀取邏輯         |
| adapter    | `InfluxSeriesReader`        | 使用 InfluxDB 實作的讀取器  |

### 補充建議

- 測試用實作以 `Mock` / `Fake` 為前綴，明確區隔
- 跨模組共享的實作命名應避免功能模糊，例如避免使用 `DefaultImpl`
- 對應文檔請一併維護於 `/docs/interfaces/` 中

### Package 命名慣例補充

- 為避免與 `pkg/ifaces/metric` 等核心介面衝突，建議 adapter 的實作包使用 `package xxxadapter` 命名，例如 `metricsadapter`, `loggeradapter`
- 這有助於區分：
  - 核心抽象層（如 `metric.Writer`、`logger.Logger`）
  - 對應的實作包（如 `internal/adapters/metrics`）
- 匯入時可維持一致，例如：

```go
import metricsadapter "detectviz/internal/adapters/metrics"
```

- 適用情境：
  - 該模組有清楚對應 interface（如 `pkg/ifaces/metric.Writer`）
  - 有多種實作可能（如 `InfluxWriter`, `PrometheusWriter`）

---

## 靜態分析與 Lint 工具

為維持一致的命名與格式，Detectviz 專案建議透過 `golangci-lint` 搭配下列規則進行開發與 CI 驗證。

### 工具安裝

```bash
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 設定範例：`.golangci.yml`

```yaml
run:
  timeout: 3m
  skip-dirs:
    - "docs/"
    - "scripts/"

linters:
  enable:
    - govet
    - revive
    - gocritic
    - gofumpt
    - misspell
    - revive
    - interfacer
    - godot
    - godox
    - gocognit
    - gocyclo
    - dupl

linters-settings:
  revive:
    rules:
      - name: var-naming
        arguments:
          allow-leading-underscore: false
          allow-underscore: false
      - name: exported
        arguments:
          allow-underscore_prefix: false
          allow-unexported: false
  gocritic:
    enabled-checks:
      - ifElseChain
      - commentFormatting

issues:
  exclude-rules:
    - text: "should have comment"
      path: "mocks/"
```

### 檢查重點

- 所有導出的 struct、interface、function 必須有對應英文註解（godot）
- 不可使用 `DefaultImpl` 等模糊命名
- 測試用物件建議以 `Mock` / `Fake` 為前綴
- 使用 gofumpt 強制格式一致

### 建議整合至 CI/CD

可於 `.gitlab-ci.yml` 或 GitHub Actions 中加入以下命令：

```bash
golangci-lint run ./...
```

如需更進階的 CI 規則，可另建立 `scripts/lint.sh` 作為統一入口點。


================================================
FILE: docs/detectviz-apps.md
================================================
# Detectviz Apps 應用層總覽

## 模組原型範例

- alert 模組 是流程編排者，觸發 analyzer、reporter。
- analyzer 模組 不實作模型，而是根據規則呼叫 anomaly-python 提供的多種 API。
- reporter 模組 可接收標記後的結果，呼叫 llm-python 生成摘要，再整合為報表。
- anomaly-python 是多種異常偵測模型的 REST API 服務。
- llm-python 是單一 LLM API（摘要、分析、重寫）服務。



================================================
FILE: docs/develop-guide.md
================================================



================================================
FILE: docs/interface-contract.md
================================================



================================================
FILE: docs/plugins-guide.md
================================================



================================================
FILE: docs/interfaces/README.md
================================================
# Interface 撰寫準則（Detectviz 專案）

本文件定義 Detectviz 專案中撰寫 interface 的統一風格與格式，供所有開發者遵循。

本規範旨在：
- 統一 interface 撰寫風格，提升維護與閱讀性
- 協助 AI 工具（如 Cursor / ChatGPT）產生一致、可替換的實作
- 促進模組化、可測試性與依賴反轉原則的落實

---

## 命名與結構原則

- Interface 檔案放置於 `pkg/ifaces/{module}/{name}.go`
- 命名使用 PascalCase（如 `Logger`, `CacheStore`, `EventBus`）
- 一個檔案只定義一個主要 interface
- 每個方法數量建議在 3～7 個以內，過多請拆分子模組

---

## 註解風格與格式

### 每個 interface 應具備：

- 英文主註解（簡潔描述用途）
- 對應的繁體中文補充，使用 `// zh:` 開頭

### 方法註解範例：

```go
// Info logs a message at the info level.
// zh: 記錄 info 級別的日誌訊息。
Info(msg string, fields ...any)
```

### Interface 註解範例：

```go
// Logger defines the structured logging interface for Detectviz.
// zh: Logger 定義 Detectviz 中的結構化日誌介面。
```

---

## 設計原則

- 僅定義「Detectviz 需要什麼」，不耦合第三方套件實作（如 zap, Redis）
- interface 為抽象 contract，不包含具體邏輯
- 若需支援 context、trace、TTL、分群等擴充性，應納入 method 設計
- 實作放在 `internal/adapters/{module}` 中

---

## 相關目錄規範

- 所有 interface 的中文說明與使用情境應補充於 `docs/interfaces/{name}.md`
- 若 interface 被核心流程注入，請註記於 `bootstrap.Init()` 文件或架構圖中

- 撰寫對應說明文件時，請參考 [interface-doc-template.md](./interface-doc-template.md)，以統一內容結構與敘述方式



================================================
FILE: docs/interfaces/alert.md
================================================
# AlertEvaluator Interface 說明文件

> 本文件為 Detectviz 專案中 `AlertEvaluator` 告警評估介面的設計說明與擴充建議。該介面採用 Clean Architecture 原則，將「告警邏輯」與「資料查詢來源」解耦，並透過統一的 `MetricQueryAdapter` 擴展不同數據來源的查詢能力。所有判斷邏輯統一由 `evaluate.go` 實作，確保運算一致性與可維護性。

---

## 介面用途（What it does）

`AlertEvaluator` 為告警模組中負責「條件邏輯判斷」的抽象介面，其核心目的是將「告警邏輯」與「資料查詢來源」完全解耦。

- 接收 `AlertCondition` 作為輸入
- 呼叫注入的 `MetricQueryAdapter` 查詢資料
- 根據條件比較結果回傳 `AlertResult`

---

## 使用情境（When and where it's used）

- 由排程器或事件觸發時，批次執行多筆條件檢查
- 可於測試階段透過 Mock Adapter 驗證評估邏輯
- 生產環境可共用一套評估器，但對接不同資料來源（如 Prometheus、InfluxDB 等）

---

## 方法說明（Methods）

### AlertEvaluator

```go
Evaluate(ctx context.Context, cond AlertCondition) (AlertResult, error)
```

- 輸入：`AlertCondition`（條件 ID、Expr、閾值、標籤等）
- 輸出：`AlertResult`（是否觸發、訊息、實際值）

---

## 預期實作（Expected implementations）

| 檔案                                      | 功能描述                                         |
|-------------------------------------------|--------------------------------------------------|
| `internal/adapters/alert/mock_adapter.go` | 測試用模擬實作，可注入任意回傳值與錯誤             |
| `internal/adapters/alert/flux/flux.go`    | 基於 Flux 查詢語言的告警評估邏輯（InfluxDB）       |
| `internal/adapters/alert/prom/prom.go`    | 基於 PromQL 的告警評估邏輯（Prometheus）          |
| `internal/adapters/alert/static.go`       | 靜態值或固定閾值比對邏輯（單元測試常用）           |
| `internal/adapters/alert/nop.go`          | 不進行任何評估的 Noop 實作                         |

> 註：所有資料查詢行為均透過 `MetricQueryAdapter` 處理，`AlertEvaluator` 僅關心運算邏輯。運算邏輯集中實作於 `pkg/ifaces/alert/evaluate.go`，支援如 `ge`, `lt`, `eq`, `ne` 等運算子。

---

## 擴充資料來源時應實作哪些部分（How to add a new data source）

### 1. 實作 MetricQueryAdapter（參見 [`docs/interfaces/metric.md`](../metric.md)）

- 建立檔案：`internal/adapters/metrics/{source}_adapter.go`
- 實作方法：`Query(ctx, expr, labels) (float64, error)`
- 包裝對應查詢語法（PromQL、Flux、SQL、API...）並解析為數值

### 2. 選擇性：提供資料來源專用 Evaluator（必要時）

- 若查詢結果須特殊處理（如比對多欄位、進行統計），可建立專屬 `AlertEvaluator`
- 例如：`fluxEvaluator := NewEvaluator(fluxAdapter)`

### 3. 注入與註冊

- 於 `bootstrap/init.go` 中依資料來源或環境設定注入對應實作

---

## 結構說明（Structs）

### AlertCondition

| 欄位     | 型別               | 說明                                 |
|----------|--------------------|--------------------------------------|
| ID       | `string`           | 條件識別碼                           |
| Expr     | `string`           | 查詢語法（如 Flux、PromQL）         |
| Threshold| `float64`          | 閾值數值                             |
| Labels   | `map[string]string` | 過濾條件的標籤組                     |
| Operator | `string`           | 比對運算子（ge, gt, lt, eq 等）     |

### AlertResult

| 欄位     | 型別      | 說明                                |
|----------|-----------|-------------------------------------|
| Firing   | `bool`    | 是否觸發告警                        |
| Value    | `float64` | 查詢到的實際值                      |
| Message  | `string`  | 描述文字，可供記錄與通知            |

---

## 關聯模組與擴充性（Related & extensibility）

- 評估器由 `AlertScheduler` 控制執行時機
- 結果可拋給 `EventDispatcher` 通知對應管道
- 下層資料查詢使用統一介面 `MetricQueryAdapter`，可自由替換與擴充
- 可於任意 adapter 中共用 `pkg/ifaces/alert/evaluate.go` 實作運算


================================================
FILE: docs/interfaces/bus.md
================================================


# Bus Interface 說明文件

> 本文件說明 Detectviz 專案中 `Bus`（事件總線）模組的介面設計、用途、實作與測試方式。此模組用於統一分派各類型事件（如 Alert、Host、Metric、Task）至對應處理器。

---

## 設計目的（Design Purpose）

- 將事件處理機制從各模組中抽離，統一交由 `Bus` 處理
- 可支援同步或非同步的事件投遞模式
- 降低模組間耦合性，提升擴充性與測試性
- 支援註冊多個處理器與事件分類

---

## Interface 定義（pkg/ifaces/bus/types.go）

```go
type EventDispatcher interface {
    RegisterAlertHandler(handler event.AlertEventHandler)
    RegisterHostHandler(handler event.HostEventHandler)
    RegisterMetricHandler(handler event.MetricEventHandler)
    RegisterTaskHandler(handler event.TaskEventHandler)
    Dispatch(ctx context.Context, e event.Event) error
}
```

---

## 實作位置（Implementations）

| 檔案位置                                         | 描述                         |
|--------------------------------------------------|------------------------------|
| `internal/adapters/eventbus/inmemory.go`         | 預設實作，將事件分派至註冊的 handler |
| `internal/registry/eventbus/registry_inmemory.go`| 註冊 InMemoryDispatcher 為預設 Bus |
| `internal/plugins/eventbus/alertlog/alert_handler.go` | Plugin 實作 AlertEvent 處理器 |

---

## 使用情境（Usage Scenarios）

- 發送 `AlertTriggeredEvent` 通知對應模組產生告警紀錄
- 發送 `MetricOverflowEvent` 進行即時告警評估
- 發送 `TaskCompletedEvent` 通知後續任務系統
- 可用於事件擴充如：Webhook、Kafka、Slack 插件

---

## 註冊與擴充方式（Registration & Extension）

- 可透過 `internal/registry/eventbus/plugins.go` 自動註冊 plugins
- plugin 可實作 `AlertEventHandler` 並透過 `eventbus.RegisterAlertHandler()` 註冊
- 預設可替換為其他實作（如 Kafka、Channel-based）

---

## 測試建議（Testing Strategy）

- 實作 `MockEventDispatcher` 測試發送與註冊行為
- 測試註冊多個 handler 並驗證是否全部觸發
- 搭配 `TestLogger` 驗證事件投遞流程

---

## 關聯模組（Related Modules）

- `event`：定義各類事件的資料結構
- `logger`：事件觸發與失敗時紀錄操作狀態
- `notifier`：最終可連接通知模組

---


================================================
FILE: docs/interfaces/cachestore.md
================================================
# CacheStore Interface 說明文件

> 本文件說明 Detectviz 專案中的 `CacheStore` 介面設計原則、應用情境與實作結構。CacheStore 模組負責管理應用中的短期暫存資料，提升效能並減少重複操作。

---

## 設計目的（Design Purpose）

- 統一快取操作接口，避免不同模組重複實作快取邏輯
- 可支援多種後端（如 memory、Redis、Nop）
- 提供測試用無狀態快取，以簡化測試流程
- 可依據模組需求擴充 TTL、key namespace、動態註冊等能力

---

## Interface 定義（Methods）

```go
type CacheStore interface {
	Get(key string) (string, error)                  // 取得指定 key 對應值，不存在時回傳錯誤
	Set(key string, val string, ttlSeconds int) error // 寫入快取，ttl 為 0 表示永久有效
	Has(key string) bool                             // 檢查 key 是否存在
	Delete(key string) error                         // 刪除指定快取資料
	Keys(prefix string) ([]string, error)            // 回傳所有符合 prefix 的 key（便於批次操作）
}
```

---

## 使用情境（Use Cases）

- 掃描模組快取設備辨識結果，避免重複查詢
- 快取自動部署結果（如 conf 註記）
- Web UI 儲存快照狀態或快取選單資料
- 排程器記憶任務執行狀態避免重複執行

---

## 預期實作（Expected Implementations）

| 類型         | 實作檔案路徑                                             | 描述                     |
|--------------|----------------------------------------------------------|--------------------------|
| Memory       | `internal/adapters/cachestore/memory/memory.go`          | 使用 go-cache 作為記憶體快取 |
| Redis        | `internal/adapters/cachestore/redis/redis.go`            | 使用 Redis 實作跨節點 TTL 快取 |
| Noop         | `internal/adapters/cachestore/nop.go`                    | 空實作，用於禁用或測試用途 |

---

## 設定與註冊方式（Registry and Config）

- 設定由 `pkg/configtypes/cache_config.go` 提供
- 註冊邏輯於 `internal/registry/cachestore/registry.go`
- 可根據 config 中指定類型自動切換對應實作
- 支援 logger 注入以監控快取行為與錯誤

---

## 測試與 Mock（Testing & Mocking）

- 所有實作皆應具備對應 `_test.go` 單元測試
- 可使用 `MockCacheStore` 進行介面驗證
- 建議將測試實作放於與 adapter 相同目錄下或 `internal/test`

---

## 擴充建議（Extension Notes）

- 若需支援分散式環境，可實作 Redis Cluster adapter
- 若需監控命中率，可加入統計與 trace log 機制
- 可考慮支援 JSON / struct 序列化版本，以簡化應用層邏輯


================================================
FILE: docs/interfaces/config.md
================================================
# ConfigProvider Interface 設計說明

> 本文件為 Detectviz 專案中 `ConfigProvider` 介面的設計原則、使用情境與實作擴充方式整理，並統一與其他 interface 文件格式。

---

## 介面用途（What it does）

`ConfigProvider` 是 Detectviz 平台中用於讀取設定值的統一抽象介面，其目標如下：

- 解耦設定讀取來源與核心邏輯（YAML、ENV、Remote 等）
- 避免硬編碼與全域變數污染
- 支援 hot-reload（視實作而定）
- 提供良好測試性（可注入 Mock 設定來源）
- 支援擴充（如套用 Config Schema、版本控制）

---

## 使用情境（When and where it's used）

- 在 `internal/bootstrap/init.go` 中初始化後注入各核心模組
- 子模組可透過 `Get`、`GetInt`、`GetBool` 取得配置參數
- 搭配 logger、notifier、scheduler 等模組進行動態設定注入
- 測試時可替換為 Fake 或 Map-based 實作

---

## 方法定義（Interface Methods）

```go
type ConfigProvider interface {
    Get(key string) string
    GetInt(key string) int
    GetBool(key string) bool
    GetOrDefault(key string, defaultVal string) string
    GetCacheConfig() configtypes.CacheConfig
    GetNotifierConfigs() []configtypes.NotifierConfig
    Logger() logger.Logger
    Reload() error
}
```

- `Get`：傳回指定 key 的字串值，若不存在回空字串
- `GetInt`：傳回整數值，無法解析時預設為 0
- `GetBool`：傳回布林值，支援 "true"/"false" 字串轉換
- `GetOrDefault`：若指定 key 無值，則傳回提供的 default 值
- `GetCacheConfig`：回傳快取模組所需的結構設定
- `GetNotifierConfigs`：回傳通知通道模組的設定清單
- `Logger`：回傳 logger 實例，供模組共用
- `Reload`：重新載入設定來源，若支援 hot-reload 機制，否則為 no-op

---

## 預期實作（Expected Implementations）

| 類型     | 路徑位置                                   | 描述                            |
|----------|--------------------------------------------|---------------------------------|
| 預設     | `pkg/config/default.go`                    | 使用 map + ENV 實作設定讀取器   |
| YAML     | `internal/adapters/config/yaml.go`         | 從指定 YAML 檔案讀取設定         |
| Remote   | `internal/adapters/config/remote.go`       | 支援 HTTP / gRPC 動態設定服務   |
| Nop/Fake | `internal/adapters/config/mock_adapter.go`<br>`internal/test/fakes/fake_config.go` | 提供測試用途的空實作或假資料     |

---

## 測試建議（Testing Strategy）

- 可透過注入 `FakeConfigProvider` 模擬錯誤或邊界值
- 建議單元測試涵蓋 `Reload` 行為與 fallback 邏輯
- 可加入 `Set()` 方法以便測試程式中直接設定參數（建議僅於測試實作中使用）
- 結合整合測試驗證模組是否正確依據 config 行為切換

---

## 擴充與整合建議（Extensions & Integration）

- 可整合 Config Schema，支援欄位驗證與版本轉換
- 搭配熱更新模組（如 fsnotify）實現自動 reload
- 與遠端設定中心（如 Consul、etcd、Spring Cloud Config）整合
- 規劃支援 json-schema-version，可提升未來 JSON 設定相容性管理

---


================================================
FILE: docs/interfaces/configtypes.md
================================================


# ConfigTypes 設定結構說明文件

> 本文件說明 Detectviz 專案中 `pkg/configtypes` 目錄內各設定結構（Config Struct）的用途、對應模組與設定來源。這些結構主要用於從靜態或動態設定載入模組所需參數，並透過 ConfigProvider 傳遞至模組內部。

---

## 設計目的（Design Purpose）

- 將每個功能模組所需的設定集中管理
- 可清楚對應設定來源（如 JSON / YAML / ENV）
- 搭配 ConfigProvider 實現動態注入或切換
- 提供 Registry 與 Bootstrap 使用的統一結構

---

## 設定結構總覽（Defined Config Structs）

| 結構名稱           | 檔案位置                           | 用途與對應模組                       |
|--------------------|------------------------------------|--------------------------------------|
| `LoggerConfig`     | `pkg/configtypes/logger_config.go` | 設定 logger 類型、Level、輸出格式等 |
| `SchedulerConfig`  | `pkg/configtypes/scheduler_config.go` | 設定排程器類型與參數                |
| `NotifierConfig`   | `pkg/configtypes/notifier_config.go` | 定義通知方式（email/slack/webhook） |
| `CacheConfig`      | `pkg/configtypes/cache_config.go`  | 定義快取實作方式與參數（memory/redis） |
| `AlertConfig`      | `pkg/configtypes/alert_config.go`  | 告警模組預設行為與閾值設定           |
| `BusConfig`        | `pkg/configtypes/bus_config.go`    | 指定使用哪一種 EventBus 實作         |
| `MetricsConfig`    | `pkg/configtypes/metrics_config.go` | 設定查詢資料來源（prom/flux）等     |

---

## 使用方式（Usage Pattern）

這些結構通常會出現在 `default.go` 或 `ConfigProvider` 的初始化中，例如：

```go
config := &configtypes.AppConfig{
    Logger:   configtypes.LoggerConfig{Level: "debug"},
    Scheduler: configtypes.SchedulerConfig{Type: "workerpool"},
}
```

或從 YAML / JSON 讀取時：

```yaml
logger:
  type: zap
  level: info

scheduler:
  type: cron
  spec: "0 * * * *"
```

---

## 測試建議（Testing Strategy）

- 可建立對應 `*_test.go` 測試配置反序列化行為
- 搭配 `testutil.LoadConfigFromYAML` 測試整體載入是否成功
- 支援 fallback 至預設值邏輯可單元測試驗證

---

## 擴充建議（Extensions）

- 若新增模組應同步建立對應 ConfigTypes 檔案（命名規則：`{module}_config.go`）
- 結構欄位建議皆標記 yaml/json tag
- 可未來對接 Config Schema 驗證或 UI 編輯介面

---


================================================
FILE: docs/interfaces/event.md
================================================


# Event Interface 說明文件

> 本文件說明 Detectviz 專案中 `Event` 類型與對應 Handler 的定義與用途。事件模組負責封裝跨模組交換的資料結構，作為 EventBus 分派的核心單位。

---

## 設計目的（Design Purpose）

- 定義各種監控情境下可能產生的事件類型
- 每種事件皆有對應 Handler interface，供 EventBus 呼叫
- 可支援未來擴充更多事件類型（如 UserEvent, ReportEvent）
- 避免硬編碼與不一致資料傳遞格式

---

## Interface 定義（pkg/ifaces/event/*.go）

目前事件定義共分為四類：

### AlertEvent

```go
type AlertTriggeredEvent struct {
	ID      string
	Level   string
	Message string
	Time    time.Time
}

type AlertEventHandler interface {
	HandleAlertEvent(ctx context.Context, event AlertTriggeredEvent) error
}
```

### HostEvent

```go
type HostRegisteredEvent struct {
	InstanceID string
	IP         string
	Time       time.Time
}

type HostEventHandler interface {
	HandleHostEvent(ctx context.Context, event HostRegisteredEvent) error
}
```

### MetricEvent

```go
type MetricOverflowEvent struct {
	Target    string
	Metric    string
	Threshold float64
	Value     float64
	Time      time.Time
}

type MetricEventHandler interface {
	HandleMetricEvent(ctx context.Context, event MetricOverflowEvent) error
}
```

### TaskEvent

```go
type TaskCompletedEvent struct {
	TaskID   string
	Success  bool
	Time     time.Time
	Message  string
}

type TaskEventHandler interface {
	HandleTaskEvent(ctx context.Context, event TaskCompletedEvent) error
}
```

---

## 使用情境（Usage Scenarios）

- 評估告警時產生 `AlertTriggeredEvent` 並通知處理器
- 掃描設備新增時觸發 `HostRegisteredEvent` 記錄並回報
- 發現異常值時觸發 `MetricOverflowEvent` 串聯至告警流程
- 任務排程執行結束時觸發 `TaskCompletedEvent` 進行分析或通知

---

## 測試建議（Testing Strategy）

- 實作 Fake EventBus 驗證每個事件類型是否能正常被觸發與處理
- 使用 `testLogger` 驗證事件是否被正確記錄
- 單元測試每個 Handler 接收到事件後的行為是否符合預期

---

## 擴充建議（Extension Notes）

- 可新增 `UserEvent`, `ReportEvent` 等事件類型
- 可補充每種事件對應 JSON Schema 供外部 API 使用
- 可設計通用 event validator 檢查結構正確性

---


================================================
FILE: docs/interfaces/eventbus.md
================================================
# EventDispatcher Interface 說明文件

> 本文件為 Detectviz 專案中事件分派模組的介面與結構說明，包含所有事件型別的 Dispatch 與 Handler 設計。

---

## 介面用途（What it does）

`EventDispatcher` 是模組間透過事件進行解耦通訊的抽象介面，主要用途為：

- 定義所有支援的事件型別與對應的註冊與分派邏輯
- 將事件廣播給所有註冊的 handler
- 支援測試用的 Nop / InMemory 實作

---

## 事件總覽（Supported Events）

| 事件名稱              | 對應介面            | Struct 定義檔案             |
|-----------------------|---------------------|-----------------------------|
| HostDiscoveredEvent   | `HostEventHandler`  | `pkg/ifaces/eventbus/host.go` |
| MetricOverflowEvent   | `MetricEventHandler`| `pkg/ifaces/eventbus/metric.go` |
| AlertTriggeredEvent   | `AlertEventHandler` | `pkg/ifaces/eventbus/alert.go` |
| TaskCompletedEvent    | `TaskEventHandler`  | `pkg/ifaces/eventbus/task.go` |

---

## EventDispatcher 介面定義

### HostDiscoveredEvent

實作位置：`pkg/ifaces/eventbus/host.go`

```go
DispatchHostDiscovered(ctx context.Context, event HostDiscoveredEvent) error
RegisterHostHandler(handler HostEventHandler)
```

- `DispatchHostDiscovered`：發送主機註冊事件
- `RegisterHostHandler`：註冊主機事件的接收者

---

### MetricOverflowEvent

實作位置：`pkg/ifaces/eventbus/metric.go`

```go
DispatchMetricOverflow(ctx context.Context, event MetricOverflowEvent) error
RegisterMetricHandler(handler MetricEventHandler)
```

- `DispatchMetricOverflow`：發送指標過多事件（如分群維度超標）
- `RegisterMetricHandler`：註冊指標事件處理器

---

### AlertTriggeredEvent

實作位置：`pkg/ifaces/eventbus/alert.go`

```go
DispatchAlertTriggered(ctx context.Context, event AlertTriggeredEvent) error
RegisterAlertHandler(handler AlertEventHandler)
```

- `DispatchAlertTriggered`：發送告警觸發事件
- `RegisterAlertHandler`：註冊接收告警事件的處理器

---

### TaskCompletedEvent

實作位置：`pkg/ifaces/eventbus/task.go`

```go
DispatchTaskCompleted(ctx context.Context, event TaskCompletedEvent) error
RegisterTaskHandler(handler TaskEventHandler)
```

- `DispatchTaskCompleted`：任務完成後發送通知
- `RegisterTaskHandler`：註冊接收任務完成事件的模組

---

## Handler 介面定義（Event Handler Interfaces）

| Handler 介面名稱       | 方法名稱                             | 接收的事件結構              |
|------------------------|--------------------------------------|-----------------------------|
| `HostEventHandler`     | `HandleHostDiscovered(ctx, event)`   | `HostDiscoveredEvent`       |
| `MetricEventHandler`   | `HandleMetricOverflow(ctx, event)`   | `MetricOverflowEvent`       |
| `AlertEventHandler`    | `HandleAlertTriggered(ctx, event)`   | `AlertTriggeredEvent`       |
| `TaskEventHandler`     | `HandleTaskCompleted(ctx, event)`    | `TaskCompletedEvent`        |

---

## 事件結構說明（Event Structs）

### HostDiscoveredEvent

| 欄位     | 說明               |
|----------|--------------------|
| Name     | 主機名稱           |
| Labels   | 額外附加標籤資訊   |
| Source   | 資料來源描述字串   |

---

### MetricOverflowEvent

| 欄位     | 說明                     |
|----------|--------------------------|
| Target   | 發生溢出的指標來源       |
| Reason   | 描述溢出的原因（可顯示） |

---

### AlertTriggeredEvent

| 欄位       | 說明                       |
|------------|----------------------------|
| ConditionID | 對應告警條件識別碼       |
| Value      | 觸發時的實際值             |
| Message    | 觸發時的說明訊息           |

---

### TaskCompletedEvent

| 欄位     | 說明               |
|----------|--------------------|
| Name     | 任務名稱           |
| Success  | 是否成功完成       |
| Message  | 額外訊息（如失敗原因） |

---

## EventDispatcher 實作範例（Dispatcher Implementations）

以下為 `EventDispatcher` 介面的具體實作，符合 `pkg/ifaces/eventbus/eventbus.go` 所定義之方法簽章與行為契約。

| 實作檔案                              | 說明                                   |
|---------------------------------------|----------------------------------------|
| `internal/adapters/eventbus/inmemory.go` | 同步事件分派（單元測試與內部模擬用途） |
| `internal/adapters/eventbus/nop.go`      | 空實作，忽略所有事件（禁用或跳過場景） |
| `internal/adapters/eventbus/custom.go`    | 預留擴充用途，可用於注入自定事件邏輯     |

---

### 各事件 Handler 預期實作檔案位置（Planned Handler Implementations）

為使事件接收模組結構清晰，建議將各事件的 `EventHandler` 實作放置於下列位置：

| 事件類型              | Handler Interface        | 建議實作檔案位置                        |
|-----------------------|--------------------------|-----------------------------------------|
| HostDiscoveredEvent   | `HostEventHandler`       | `internal/adapters/eventbus/host.go`    |
| MetricOverflowEvent   | `MetricEventHandler`     | `internal/adapters/eventbus/metric.go`  |
| AlertTriggeredEvent   | `AlertEventHandler`      | `internal/adapters/eventbus/alert.go`   |
| TaskCompletedEvent    | `TaskEventHandler`       | `internal/adapters/eventbus/task.go`    |

上述檔案應包含對應的 `HandleXxx(ctx, event)` 方法實作。

每個實作皆應實作 `pkg/ifaces/eventbus/eventbus.go` 中的 `EventDispatcher` 介面。

---

## 擴充建議（Extensibility）

- 每新增一個事件型別，請建立：
  - 對應的 `Event struct`（如 `XxxEvent`）
  - 對應的 `Handler interface`（如 `XxxEventHandler`）
  - 在 `EventDispatcher` 中定義 `DispatchXxx` 與 `RegisterXxxHandler`
- 建議將每個事件類型放入獨立檔案管理，並在 `eventbus.go` 中用區塊分隔定義。

---

## Logger 設定與測試輔助（Logger Injection and Testing）

### plugin 測試建議

為了讓事件處理器中的 log 行為可被測試與驗證，`EventDispatcher` 支援設定全域預設 logger。可用於單元測試中攔截 log 輸出。

```go
import testlogger "github.com/detectviz/detectviz/internal/test/testutil"
import pluginlog "github.com/detectviz/detectviz/internal/registry/eventbus"

log := testlogger.NewTestLogger()
pluginlog.OverrideDefaultLogger(log)  // 在測試前先注入
```

事件處理器可在內部這樣取得 logger：

```go
log := pluginlog.GetDefaultLogger()
log.Info("plugin triggered")
```

### 測試建議配套

請使用內建 `TestLogger` 並驗證其 `Messages()` 方法是否包含預期內容，詳見：

- `internal/test/testutil/test_logger.go`
- `internal/test/plugins/eventbus/alertlog/alert_plugin_test.go`



================================================
FILE: docs/interfaces/interface-doc-template.md
================================================
# {Interface Name} Interface 說明文件

> 本文件為 Detectviz 專案中 `{Interface Name}` 介面的設計說明與使用情境整理。
> 
> 請將 `{Interface Name}` 替換為實際名稱，並填入每段描述內容。

---

## 介面用途（What it does）

說明此 interface 的抽象職責與在架構中的角色，例如：

- 提供設定存取抽象層
- 支援模組間的訊息廣播
- 統一快取存取與替換機制
- 與 OpenTelemetry 整合之觀測能力抽象

---

## 使用情境（When and where it's used）

列出具體模組或流程中會用到此介面的情況，例如：

- 在 `bootstrap.Init()` 中注入供 core 使用
- 被 HTTP middleware 呼叫以取得 trace context
- 在任務排程器中使用以避免重複執行

---

## 方法說明（Methods）

針對每個方法補充說明其功能與回傳值意圖，例如：

- `Get(key string) string`：回傳對應設定值，若不存在可能為空字串
- `WithContext(ctx)`：將 context 內 trace id 注入 logger 流程
- `Publish(topic string, payload any)`：廣播一筆訊息給訂閱者

---

## 預期實作（Expected implementations）

列出你預期會有哪些實作，以及建議放置目錄與用途：

- `internal/adapters/logger/zaplogger.go`
- `internal/adapters/cachestore/memcache.go`
- `internal/adapters/eventbus/inmemory.go`
- `internal/adapters/{module}/nop.go`（測試用）

---

## 關聯模組與擴充性（Related & extensibility）

如有與其他模組或 interface 的整合需求，或未來可擴充方向，請說明：

- 與 trace / metrics 的整合
- Redis/NATS 替代支援
- 可 plugin 化設計

---


================================================
FILE: docs/interfaces/logger.md
================================================
# Logger 介面設計與使用指南

> 本文件說明 Detectviz 專案中 `Logger` 介面的設計目標、主要功能、典型應用場景、實作方式，以及測試與擴充建議。Logger 為專案核心元件，支援結構化日誌輸出與追蹤整合，便於後續分析、監控與除錯。

---

## 介面功能摘要

Logger 為應用層日誌統一抽象，提供：

- info/warn/error/debug 等多層級日誌
- 結構化欄位輸出，便於檢索與分析
- 支援 context 傳遞 trace_id、span_id 等追蹤資訊
- 可擴充為 OTLP 格式，整合 Loki、Tempo 等觀測系統

---

## 典型應用場景

- `bootstrap.Init()` 注入 logger 與對應 adapter
- Middleware 記錄 HTTP 請求與 trace 資訊
- 背景任務、排程、事件處理等執行流程與錯誤日誌
- 搭配 `WithFields()` 記錄告警、業務欄位等結構化資訊

---

## 介面定義

Logger interface 介面如下：

```go
type Logger interface {
    Debug(msg string, args ...any)
    Info(msg string, args ...any)
    Warn(msg string, args ...any)
    Error(msg string, args ...any)
    WithFields(fields map[string]any) Logger
    WithContext(ctx context.Context) Logger
    Named(name string) Logger
    Sync() error
}
```

### 方法說明

- `Debug/Info/Warn/Error`：輸出對應層級的訊息
- `WithFields`：加入結構化欄位（如 rule_id、host 等）
- `WithContext`：保留 context 資訊（trace_id 等）
- `Named`：指定 logger 模組名稱
- `Sync`：flush buffer（如 zap 需明確呼叫）

---

## Context 整合工具

定義於 `pkg/ifaces/logger/context.go`：

```go
func WithContext(ctx context.Context, l Logger) context.Context
func FromContext(ctx context.Context) Logger
```

說明：

- 可於 middleware 注入 logger 實例，後續流程可從 context 擷取
- `FromContext` 未注入時會回傳 fallback 的 `NopLogger`

---

## 預期實作與關聯模組

| 檔案位置                                       | 說明                             |
|------------------------------------------------|----------------------------------|
| `internal/adapters/logger/zap_adapter.go`      | zap 套件實作，支援欄位與模組命名 |
| `internal/adapters/logger/nop_adapter.go`      | 空實作，靜默略過所有輸出         |
| `pkg/ifaces/logger/nop_logger.go`              | NopLogger 結構，為 fallback 預設 |
| `pkg/ifaces/logger/context.go`                 | context 操作工具函式             |

### 擴充性與整合建議

- 可整合 OpenTelemetry trace context，實現 log-trace 關聯
- 支援 OTLP 匯出，整合 Loki、Tempo 等觀測後端
- plugin 機制支援自訂 logger backend（file、stdout、Redis 等）

---

## 測試與驗證方式

### 測試檔案

| 檔案路徑                                      | 測試內容                             |
|-----------------------------------------------|--------------------------------------|
| `internal/adapters/logger/logger_test.go`     | 測試 ZapLogger 輸出格式與行為        |
|                                               | 測試 NopLogger 是否靜默處理           |
| `pkg/ifaces/logger/context_test.go`           | 驗證 context 工具函式正確行為         |

### 重要測試場景

- logger 實例可注入並從 context 中擷取
- ZapLogger 可輸出結構化欄位與多級訊息
- NopLogger 可調用且不產生錯誤
- 未注入時，FromContext 回傳 fallback logger

### 建議測試方式

- 使用 `zaptest/observer` 驗證 logger 輸出
- 使用標準 `testing` 驗證 interface 契約與 fallback 行為


================================================
FILE: docs/interfaces/metric.md
================================================
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


================================================
FILE: docs/interfaces/modules.md
================================================


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
  StartAll(ctx context.Context) error
  StopAll(ctx context.Context) error
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
if err := runner.StartAll(ctx); err != nil {
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



================================================
FILE: docs/interfaces/notifier.md
================================================
# Notifier Interface 設計說明

> 本文件說明 `pkg/ifaces/notifier` 中的通知模組介面設計、使用情境與實作結構。Notifier 模組是 detectviz 的通用訊息推送介面，負責統一封裝告警與事件訊息，並透過 Email、Slack、Webhook 等方式傳遞至外部。

---

## 設計目標（Design Goals）

- 抽象化通知機制，支援多種傳送通道（如 email, webhook, slack）
- 支援訊息格式、通知等級、推送目標等欄位封裝
- 可由 alert 模組、scheduler 模組、eventbus 模組觸發使用
- 易於擴充與動態註冊
- 搭配 logger 記錄成功與錯誤資訊

---

## Interface 定義（Interface Definitions）

```go
type Notifier interface {
	Send(ctx context.Context, msg Message) error
}

type Message struct {
	Title   string            // 通知標題
	Content string            // 通知內容
	Level   string            // 分類：info / warning / critical
	Target  string            // 接收對象，例如 webhook URL、email
	Time    time.Time         // 發送時間
}
```

---

## 使用情境（When and where it's used）

- 告警觸發時由 AlertEvaluator 推送通知
- 任務完成後由 Scheduler 發送執行結果
- 系統異常時透過 EventBus 發送通知
- 可於 plugin 中擴充事件通知行為

---

## 實作位置與類型（Implementations）

| 名稱             | 檔案路徑                                               | 描述                         |
|------------------|--------------------------------------------------------|------------------------------|
| EmailNotifier     | `internal/adapters/notifier/email_adapter.go`          | 寄送 email 通知               |
| SlackNotifier     | `internal/adapters/notifier/slack_adapter.go`          | 發送 Slack 訊息               |
| WebhookNotifier   | `internal/adapters/notifier/webhook_adapter.go`        | 送出 HTTP POST 通知          |
| MockNotifier      | `internal/adapters/notifier/mock_adapter.go`           | 單元測試用，不實際發送        |
| MultiNotifier     | `internal/adapters/notifier/multi.go`                  | 將一則訊息傳送給多個 Notifier |
| NopNotifier       | `internal/adapters/notifier/nop.go`                    | 無動作通知器（開發或測試使用）|

---

## 設定與註冊方式（Configuration & Registration）

- 設定來源：`pkg/configtypes/notifier_config.go`
- 動態註冊位置：`internal/registry/notifier/registry.go`
- 支援多個 notifier 並以 `[]configtypes.NotifierConfig` 批次註冊
- 可注入 logger 與自定義 http client

---

## 測試結構（Testing Structure）

| 測試檔案位置                                         | 測試目標                        |
|------------------------------------------------------|---------------------------------|
| `mock_adapter_test.go`                               | 驗證介面契約與訊息結構          |
| `email_adapter_test.go`（預定）                      | 測試 email 發送與錯誤處理       |
| `slack_adapter_test.go`（預定）                      | 測試 Slack 訊息格式與傳遞邏輯   |
| `webhook_adapter_test.go`（預定）                   | 測試 HTTP POST 傳送邏輯         |
| `multi_test.go` / `nop_test.go`（補充中）            | 測試組合與靜態 fallback 行為    |

---

## 擴充方式（How to add a new Notifier）

1. 建立檔案於 `internal/adapters/notifier/{name}_adapter.go`
2. 實作 `pkg/ifaces/notifier.Notifier` 介面
3. 可搭配自定義 logger、HTTP client、retry 邏輯
4. 補上單元測試：`{name}_adapter_test.go`
5. 註冊至 `internal/registry/notifier/registry.go`
6. 擴充 config 設定於 `pkg/configtypes/notifier_config.go`

---

## 相關依賴模組（Related Modules）

- `logger`：記錄通知發送與錯誤
- `eventbus`：事件觸發來源
- `alert` / `scheduler`：主要呼叫來源
- `configtypes.NotifierConfig`：指定通道參數與開關

---


================================================
FILE: docs/interfaces/scheduler.md
================================================
# Scheduler Interface 設計說明

> 本文件說明 `pkg/ifaces/scheduler` 模組中的排程器介面設計與使用方式，並對應 detectviz 專案的實作結構與測試策略。

---

## 設計目標（Design Goals）

Scheduler 模組負責執行背景任務與週期性排程，核心目標包括：

- 抽象化任務調度邏輯
- 支援多種排程策略（如 Cron、Worker Pool）
- 任務具備命名、週期（spec）與執行邏輯
- 支援注入 logger 與 retry 機制
- 易於單元測試與擴充其他排程器實作

---

## Interface 定義（Interface Definitions）

```go
type Job interface {
	Name() string              // 任務名稱
	Spec() string              // 排程規則，如 cron 表達式
	Run(ctx context.Context) error // 執行邏輯
}

type Scheduler interface {
	Register(job Job)                // 註冊任務
	Start(ctx context.Context) error // 啟動排程器
	Stop(ctx context.Context) error  // 停止排程器
}
```

---

## 使用情境（Usage Scenarios）

Scheduler 適用於以下場景：

- 定時資料清理任務
- 定期健康檢查與狀態回報
- 指標彙整與轉拋（例如送出至其他系統）
- 延遲與重試任務（例如任務失敗時重新排程）

---

## 實作位置（Implementations）

| 類型                | 檔案路徑                                                 | 描述                                          |
|---------------------|----------------------------------------------------------|-----------------------------------------------|
| CronScheduler        | `internal/adapters/scheduler/cron_adapter.go`           | 使用 robfig/cron 實現基於 Spec 的排程器       |
| WorkerPoolScheduler  | `internal/adapters/scheduler/workerpool_adapter.go`     | 使用 goroutine pool 實現具備 retry 與 logger |
| MockScheduler        | `internal/adapters/scheduler/mock_adapter.go`           | 測試用，模擬註冊與執行行為                   |

---

## 測試位置（Testing Files）

| 測試檔案路徑                                               | 測試內容                                   |
|------------------------------------------------------------|--------------------------------------------|
| `workerpool_adapter_test.go`                               | 測試排程、併發執行與 retry 行為             |
| `cron_adapter_test.go`                                     | 測試 Cron 表達式排程是否正確執行            |
| `mock_adapter_test.go`                                     | 驗證任務註冊與啟動流程                      |
| `testlogger.go`                                            | 提供測試環境下使用的 logger 實作            |

---

## 実作註冊（Registry）

於 `internal/registry/scheduler/registry.go` 註冊 Scheduler 實作：

```go
func ProvideScheduler(log logger.Logger) scheduler.Scheduler {
	return adapters.NewWorkerPoolScheduler(4, log)
}
```

如需使用其他排程方式，例如 Cron：

```go
return adapters.NewCronScheduler(log)
```

---

## 擴充方式（How to Add a New Scheduler）

若需新增其他排程器類型（如 DelayQueueScheduler）：

1. 建立實作檔案：`internal/adapters/scheduler/{name}_adapter.go`
2. 實作 `Scheduler` 介面三個方法：`Register`、`Start`、`Stop`
3. 可注入 logger，並實作 retry 行為（依需求）
4. 加入單元測試：`{name}_adapter_test.go`
5. 註冊於 `registry.go` 供主程式使用

---

## 關聯模組（Related Modules）

- `logger`：可注入日誌紀錄任務執行與錯誤
- `eventbus`：任務完成後可發出事件
- `config`：可支援未來以動態設定排程行為

---



================================================
FILE: docs/interfaces/server.md
================================================


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


================================================
FILE: docs/testing/test-guide.md
================================================
# 測試指南文件已移動

本測試設計與實作規範文件，已移至以下路徑以方便與測試邏輯共存：

**新位置：**  
[/internal/test/README.md](/internal/test/README.md)

該文件涵蓋以下主題：

- 單元測試、整合測試、Fake、Mock、測試工具分類原則
- 測試結構與目錄命名慣例
- 各模組對應測試位置與責任劃分
- 實際範例說明

請前往新位置查看完整內容。



================================================
FILE: internal/adapters/alert/evaluator.go
================================================
// Package alert 提供 AlertEvaluator 的預設實作，用於根據規則比對指標狀態。
package alertadapter

import (
	"context"
	"fmt"

	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	metric "github.com/detectviz/detectviz/pkg/ifaces/metrics"
)

// DefaultAlertEvaluator 是 alert.AlertEvaluator 的預設實作。
// zh: 預設告警判斷器，用於依據查詢結果與規則閾值進行比對。
type DefaultAlertEvaluator struct {
	Log   logger.Logger             // zh: 用於記錄評估過程與錯誤資訊
	Query metric.MetricQueryAdapter // zh: 指標查詢介面（需由外部注入）
}

// NewDefaultAlertEvaluator 建立一個 DefaultAlertEvaluator 實例。
// zh: 可傳入 logger 與 metric.QueryAdapter 元件以追蹤告警觸發過程並查詢指標。
func NewDefaultAlertEvaluator(log logger.Logger, query metric.MetricQueryAdapter) *DefaultAlertEvaluator {
	return &DefaultAlertEvaluator{
		Log:   log.Named("alertevaluator"),
		Query: query,
	}
}

// Evaluate 根據告警條件進行比對。
// zh: 根據 AlertCondition 查詢指標結果，並依據閾值與運算子進行告警判斷。
func (e *DefaultAlertEvaluator) Evaluate(ctx context.Context, cond alert.AlertCondition) (alert.AlertResult, error) {
	e.Log.Debug("evaluating alert condition", "expr", cond.Expr, "threshold", cond.Threshold)

	// 查詢指標資料
	value, err := e.Query.Query(ctx, cond.Expr, cond.Labels)
	if err != nil {
		e.Log.Error("query failed", "error", err)
		return alert.AlertResult{
			Firing:  false,
			Message: fmt.Sprintf("failed to evaluate '%s': %v", cond.Expr, err),
			Value:   0,
		}, err
	}

	e.Log.Debug("evaluation result", "expr", cond.Expr, "value", value, "threshold", cond.Threshold)

	// 呼叫通用比對邏輯
	result, evalErr := alert.Evaluate(value, cond)
	return result, evalErr
}



================================================
FILE: internal/adapters/alert/mock_adapter.go
================================================
package alertadapter

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/alert"
)

// MockAlertEvaluator 是 alert.AlertEvaluator 的模擬實作。
// zh: 用於測試的假告警判斷器，可自訂回傳結果與錯誤。
type MockAlertEvaluator struct {
	MockResult alert.AlertResult // zh: 預設回傳的結果
	MockError  error             // zh: 預設回傳的錯誤
}

// NewMockAlertEvaluator 建立一個 MockAlertEvaluator 實例。
// zh: 可傳入欲回傳的結果與錯誤，模擬告警行為。
func NewMockAlertEvaluator(result alert.AlertResult, err error) *MockAlertEvaluator {
	return &MockAlertEvaluator{
		MockResult: result,
		MockError:  err,
	}
}

// Evaluate 回傳預設的結果與錯誤。
// zh: 不執行任何邏輯，僅回傳初始化時指定的內容。
func (m *MockAlertEvaluator) Evaluate(ctx context.Context, cond alert.AlertCondition) (alert.AlertResult, error) {
	return m.MockResult, m.MockError
}



================================================
FILE: internal/adapters/alert/mock_adapter_test.go
================================================
package alertadapter

import (
	"context"
	"errors"
	"testing"

	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/stretchr/testify/assert"
)

// TestMockAlertEvaluator_ReturnsExpectedResult 驗證模擬評估器能正確回傳指定結果。
// zh: 模擬告警條件成功觸發時的行為。
func TestMockAlertEvaluator_ReturnsExpectedResult(t *testing.T) {
	want := alert.AlertResult{
		Firing:  true,
		Message: "mock triggered",
		Value:   42.0,
	}

	mock := NewMockAlertEvaluator(want, nil)

	got, err := mock.Evaluate(context.Background(), alert.AlertCondition{
		Expr:      "mock_expr",
		Threshold: 40,
	})

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}

// TestMockAlertEvaluator_ReturnsError 驗證模擬評估器在預期錯誤時能正確回傳 error。
// zh: 模擬告警評估過程中出現錯誤時的行為。
func TestMockAlertEvaluator_ReturnsError(t *testing.T) {
	mockErr := errors.New("mock error")
	mock := NewMockAlertEvaluator(alert.AlertResult{}, mockErr)

	_, err := mock.Evaluate(context.Background(), alert.AlertCondition{})

	assert.Error(t, err)
	assert.EqualError(t, err, "mock error")
}



================================================
FILE: internal/adapters/alert/flux/flux.go
================================================
// Package flux 提供基於 Flux 查詢語言的 AlertEvaluator 實作。
package fluxadapter

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// FluxAlertEvaluator 是基於 Flux 查詢語言的告警評估器實作。
// zh: 透過 InfluxDB 的 Flux 語法進行指標查詢與告警條件比對。
type FluxAlertEvaluator struct {
	Log logger.Logger // zh: 日誌紀錄器，用於除錯與追蹤查詢過程
}

// NewEvaluator 建立一個 FluxAlertEvaluator。
// zh: 傳入 logger 實例以支援除錯與紀錄。
func NewEvaluator(log logger.Logger) alert.AlertEvaluator {
	return &FluxAlertEvaluator{
		Log: log.Named("flux"),
	}
}

// Evaluate 根據 AlertCondition 進行查詢與比對。
// zh: 目前尚未實作實際查詢與比對邏輯。
func (f *FluxAlertEvaluator) Evaluate(ctx context.Context, cond alert.AlertCondition) (alert.AlertResult, error) {
	f.Log.Debug("flux evaluator not implemented", "expr", cond.Expr)

	// TODO: 解析 cond.Expr，執行 Flux 查詢，解析回傳值進行閾值比對
	return alert.AlertResult{
		Firing:  false,
		Message: "flux evaluation not implemented",
		Value:   0,
	}, nil
}



================================================
FILE: internal/adapters/alert/flux/flux_test.go
================================================
package fluxadapter

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/internal/test/testutil"
	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"github.com/stretchr/testify/assert"
)

func TestFluxAlertEvaluator_Evaluate(t *testing.T) {
	type fields struct {
		logger logger.Logger
	}
	type args struct {
		ctx  context.Context
		cond alert.AlertCondition
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantFiring bool
		wantMsg    string
		wantErr    bool
	}{
		{
			name: "not implemented",
			fields: fields{
				logger: testutil.NewTestLogger(),
			},
			args: args{
				ctx: context.Background(),
				cond: alert.AlertCondition{
					Expr:      `from(bucket:"test") |> range(start:-1h)`,
					Threshold: 80,
					Labels:    map[string]string{"host": "dev"},
				},
			},
			wantFiring: false,
			wantMsg:    "flux evaluation not implemented",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEvaluator(tt.fields.logger)
			result, err := e.Evaluate(tt.args.ctx, tt.args.cond)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantFiring, result.Firing)
			assert.Equal(t, tt.wantMsg, result.Message)
		})
	}
}



================================================
FILE: internal/adapters/alert/prom/prom.go
================================================
// Package prom 提供基於 Prometheus 查詢語言的 AlertEvaluator 實作。
package promadapter

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// PromAlertEvaluator 是基於 PromQL 的告警評估器實作。
// zh: 使用 Prometheus 查詢語法進行評估與比對。
type PromAlertEvaluator struct {
	Log logger.Logger // zh: 日誌記錄器，用於觀察查詢與比對過程
}

// NewEvaluator 建立一個 PromAlertEvaluator。
// zh: 可傳入 logger 實例以觀察觸發情況。
func NewEvaluator(log logger.Logger) alert.AlertEvaluator {
	return &PromAlertEvaluator{
		Log: log.Named("prom"),
	}
}

// Evaluate 根據 AlertCondition 執行查詢與閾值比對。
// zh: 實際查詢與閾值判斷尚未實作。
func (p *PromAlertEvaluator) Evaluate(ctx context.Context, cond alert.AlertCondition) (alert.AlertResult, error) {
	p.Log.Debug("prom evaluator not implemented", "expr", cond.Expr)

	// TODO: 實作 Prometheus 查詢、解析結果、進行閾值比對
	return alert.AlertResult{
		Firing:  false,
		Message: "prom evaluation not implemented",
		Value:   0,
	}, nil
}



================================================
FILE: internal/adapters/alert/prom/prom_test.go
================================================
package promadapter

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/internal/test/testutil"
	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/stretchr/testify/assert"
)

func TestPromAlertEvaluator_Evaluate(t *testing.T) {
	e := NewEvaluator(testutil.NewTestLogger())

	tests := []struct {
		name        string
		condition   alert.AlertCondition
		wantFiring  bool
		wantMessage string
		wantErr     bool
	}{
		{
			name: "basic evaluation",
			condition: alert.AlertCondition{
				Expr:      "up{job=\"node\"} == 0",
				Threshold: 1,
				Labels:    map[string]string{"job": "node"},
			},
			wantFiring:  false,
			wantMessage: "prom evaluation not implemented",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.Evaluate(context.Background(), tt.condition)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantFiring, result.Firing)
			assert.Equal(t, tt.wantMessage, result.Message)
		})
	}
}



================================================
FILE: internal/adapters/cachestore/registry.go
================================================
package cachestoreadapter

import (
	"errors"
	"fmt"
	"sort"

	memoryadapter "github.com/detectviz/detectviz/internal/adapters/cachestore/memory"
	redisadapter "github.com/detectviz/detectviz/internal/adapters/cachestore/redis"
	"github.com/detectviz/detectviz/pkg/ifaces/cachestore"
	goredis "github.com/redis/go-redis/v9"
)

// DefaultBackend 為預設的快取實作。
// zh: 若無指定 backend，預設使用 memory。
const DefaultBackend = "memory"

// ErrUnknownBackend 表示快取 backend 未註冊。
// zh: 用於回傳尚未註冊的 backend 錯誤。
var ErrUnknownBackend = errors.New("unknown cache backend")

var registry = map[string]func() cachestore.CacheStore{
	"memory": func() cachestore.CacheStore {
		return memoryadapter.NewMemoryCacheStore()
	},
}

// Register registers a new cache backend.
// zh: 註冊一個新的快取後端。
func Register(name string, fn func() cachestore.CacheStore) {
	registry[name] = fn
}

// Get returns a registered cache backend by name.
// zh: 根據名稱取得對應的快取實作，若不存在則回傳錯誤。
func Get(name string) (cachestore.CacheStore, error) {
	if fn, ok := registry[name]; ok {
		return fn(), nil
	}
	return nil, fmt.Errorf("%w: %s", ErrUnknownBackend, name)
}

// GetDefault returns the default cache backend.
// zh: 取得預設的快取後端實作。
func GetDefault() cachestore.CacheStore {
	fn := registry[DefaultBackend]
	return fn()
}

// WithRedisClient registers and returns a Redis-based cache store.
// zh: 註冊並建立 Redis 快取實作，用於注入 redis client。
func WithRedisClient(client *goredis.Client) cachestore.CacheStore {
	Register("redis", func() cachestore.CacheStore {
		return redisadapter.NewRedisCacheStore(client)
	})
	return redisadapter.NewRedisCacheStore(client)
}

// List returns the names of all registered cache backends.
// zh: 回傳所有已註冊的快取後端名稱清單（依字母排序）。
func List() []string {
	var names []string
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}



================================================
FILE: internal/adapters/cachestore/registry_test.go
================================================
package cachestoreadapter_test

import (
	"testing"

	cachestoreadapter "github.com/detectviz/detectviz/internal/adapters/cachestore"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestGetMemoryBackend(t *testing.T) {
	store, err := cachestoreadapter.Get("memory")
	assert.NoError(t, err)
	assert.NotNil(t, store)
}

func TestGetUnknownBackend(t *testing.T) {
	store, err := cachestoreadapter.Get("unknown")
	assert.Error(t, err)
	assert.Nil(t, store)
}

func TestWithRedisClient(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   9,
	})
	cachestoreadapter.WithRedisClient(client)

	store, err := cachestoreadapter.Get("redis")
	assert.NoError(t, err)
	assert.NotNil(t, store)
}

func TestListBackends(t *testing.T) {
	names := cachestoreadapter.List()
	assert.Contains(t, names, "memory")
	assert.Contains(t, names, "redis")
}



================================================
FILE: internal/adapters/cachestore/memory/memory.go
================================================
package memoryadapter

import (
	"strings"
	"sync"
)

// MemoryCacheStore 是基於 sync.Map 的快取實作。
// zh: 適用於單元測試與本地環境的記憶體快取。
type MemoryCacheStore struct {
	store sync.Map
}

// NewMemoryCacheStore 回傳一個新的記憶體快取實例。
func NewMemoryCacheStore() *MemoryCacheStore {
	return &MemoryCacheStore{}
}

// Get 取得指定 key 的快取值。
// zh: 若 key 不存在則回傳空字串與 nil 錯誤。
func (m *MemoryCacheStore) Get(key string) (string, error) {
	val, ok := m.store.Load(key)
	if !ok {
		return "", nil
	}
	return val.(string), nil
}

// Set 寫入 key 對應的快取值，忽略 TTL。
// zh: TTL 僅為佔位參數，目前不實作過期。
func (m *MemoryCacheStore) Set(key, value string, ttl int) error {
	m.store.Store(key, value)
	return nil
}

// Has 檢查指定 key 是否存在於快取中。
// zh: 若 key 存在則回傳 true，否則 false。此實作永不錯誤，僅為符合介面定義。
func (m *MemoryCacheStore) Has(key string) (bool, error) {
	_, ok := m.store.Load(key)
	return ok, nil
}

// Delete 移除指定 key。
// zh: 刪除快取資料，成功不回傳錯誤。
func (m *MemoryCacheStore) Delete(key string) error {
	m.store.Delete(key)
	return nil
}

// Keys 回傳符合指定 prefix 的所有 key。
// zh: 支援 key 模糊查詢，用於群組操作。
func (m *MemoryCacheStore) Keys(prefix string) ([]string, error) {
	var keys []string
	m.store.Range(func(k, _ any) bool {
		s := k.(string)
		if strings.HasPrefix(s, prefix) {
			keys = append(keys, s)
		}
		return true
	})
	return keys, nil
}



================================================
FILE: internal/adapters/cachestore/memory/memory_test.go
================================================
package memoryadapter_test

import (
	"testing"

	memoryadapter "github.com/detectviz/detectviz/internal/adapters/cachestore/memory"
	"github.com/stretchr/testify/assert"
)

func TestMemoryCacheStore_BasicOperations(t *testing.T) {
	store := memoryadapter.NewMemoryCacheStore()

	err := store.Set("foo", "123", 0)
	assert.NoError(t, err)
	err = store.Set("bar", "456", 0)
	assert.NoError(t, err)

	has, err := store.Has("foo")
	assert.NoError(t, err)
	assert.True(t, has)

	has, err = store.Has("bar")
	assert.NoError(t, err)
	assert.True(t, has)

	has, err = store.Has("baz")
	assert.NoError(t, err)
	assert.False(t, has)

	val, err := store.Get("foo")
	assert.NoError(t, err)
	assert.Equal(t, "123", val)

	val, err = store.Get("bar")
	assert.NoError(t, err)
	assert.Equal(t, "456", val)

	val, err = store.Get("baz")
	assert.NoError(t, err)
	assert.Equal(t, "", val)

	err = store.Delete("foo")
	assert.NoError(t, err)

	has, err = store.Has("foo")
	assert.NoError(t, err)
	assert.False(t, has)
}

func TestMemoryCacheStore_Keys(t *testing.T) {
	store := memoryadapter.NewMemoryCacheStore()

	err := store.Set("user:1", "A", 0)
	assert.NoError(t, err)
	err = store.Set("user:2", "B", 0)
	assert.NoError(t, err)
	err = store.Set("device:1", "X", 0)
	assert.NoError(t, err)

	keys, err := store.Keys("user:")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"user:1", "user:2"}, keys)
}



================================================
FILE: internal/adapters/cachestore/redis/redis.go
================================================
package redisadapter

import (
	"context"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// RedisCacheStore 是 Redis 實作的 CacheStore。
// zh: 使用 Redis 作為儲存後端，支援快取存取與過期控制。
type RedisCacheStore struct {
	client *goredis.Client
	ctx    context.Context
}

// NewRedisCacheStore 建立 Redis 快取實例。
// zh: 須傳入 go-redis v9 的 Redis 客戶端。
func NewRedisCacheStore(client *goredis.Client) *RedisCacheStore {
	return &RedisCacheStore{
		client: client,
		ctx:    context.Background(),
	}
}

// Get 取得指定 key 的值。
// zh: 若 key 不存在則回傳空字串與 nil 錯誤。
func (r *RedisCacheStore) Get(key string) (string, error) {
	val, err := r.client.Get(r.ctx, key).Result()
	if err == goredis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return val, nil
}

// Set 寫入 key 的值與過期時間（秒）。
// zh: ttl 以秒為單位，若 ttl <= 0 則為永久存留。
func (r *RedisCacheStore) Set(key, value string, ttl int) error {
	return r.client.Set(r.ctx, key, value, time.Duration(ttl)*time.Second).Err()
}

// Has 檢查 key 是否存在。
// zh: 若存在回傳 true，否則回傳 false，若查詢失敗則回傳 error。
func (r *RedisCacheStore) Has(key string) (bool, error) {
	val, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

// Delete 移除指定 key。
// zh: 若 key 存在則刪除，否則無動作，回傳刪除錯誤（若有）。
func (r *RedisCacheStore) Delete(key string) error {
	_, err := r.client.Del(r.ctx, key).Result()
	return err
}

// Keys 回傳符合 prefix 的所有 key。
// zh: 僅回傳以 prefix 開頭的 key 清單。
func (r *RedisCacheStore) Keys(prefix string) ([]string, error) {
	keys, err := r.client.Keys(r.ctx, prefix+"*").Result()
	if err != nil {
		return nil, err
	}
	var filtered []string
	for _, k := range keys {
		if strings.HasPrefix(k, prefix) {
			filtered = append(filtered, k)
		}
	}
	return filtered, nil
}



================================================
FILE: internal/adapters/cachestore/redis/redis_test.go
================================================
package redisadapter_test

import (
	"testing"
	"time"

	redisadapter "github.com/detectviz/detectviz/internal/adapters/cachestore/redis"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func newTestClient() *redisadapter.RedisCacheStore {
	client := redisadapter.NewRedisCacheStore(redisclient())
	_ = client.Delete("test:key1")
	_ = client.Delete("test:key2")
	_ = client.Delete("test:temp")
	return client
}

func redisclient() *goredis.Client {
	return goredis.NewClient(&goredis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
}

func TestRedisCacheStore_CRUD(t *testing.T) {
	store := newTestClient()

	err := store.Set("test:key1", "value1", 10)
	assert.NoError(t, err)

	has, err := store.Has("test:key1")
	assert.NoError(t, err)
	assert.True(t, has)

	val, err := store.Get("test:key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", val)

	err = store.Delete("test:key1")
	assert.NoError(t, err)

	has, err = store.Has("test:key1")
	assert.NoError(t, err)
	assert.False(t, has)

	val, err = store.Get("test:key1")
	assert.NoError(t, err)
	assert.Equal(t, "", val)
}

func TestRedisCacheStore_TTL(t *testing.T) {
	store := newTestClient()

	err := store.Set("test:temp", "soon-expire", 1)
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)

	has, err := store.Has("test:temp")
	assert.NoError(t, err)
	assert.False(t, has)

	val, err := store.Get("test:temp")
	assert.NoError(t, err)
	assert.Equal(t, "", val)
}

func TestRedisCacheStore_Keys(t *testing.T) {
	store := newTestClient()

	err := store.Set("test:key1", "1", 10)
	assert.NoError(t, err)
	err = store.Set("test:key2", "2", 10)
	assert.NoError(t, err)

	keys, err := store.Keys("test:")
	assert.NoError(t, err)
	assert.ElementsMatch(t, []string{"test:key1", "test:key2"}, keys)
}



================================================
FILE: internal/adapters/eventbus/alert.go
================================================
package eventbusadapter

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/eventbus"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// AlertLoggerHandler is a sample implementation of AlertEventHandler that logs the alert.
// zh: AlertLoggerHandler 是接收到告警觸發事件時記錄日誌的處理器實作範例。
type AlertLoggerHandler struct{}

// HandleAlertTriggered handles the alert triggered event by logging it.
// zh: 接收到告警事件後，透過 logger 模組輸出結構化日誌。
func (h *AlertLoggerHandler) HandleAlertTriggered(ctx context.Context, event event.AlertTriggeredEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"alert_id":   event.AlertID,
		"rule_name":  event.RuleName,
		"level":      event.Level,
		"instance":   event.Instance,
		"metric":     event.Metric,
		"comparison": event.Comparison,
		"value":      event.Value,
		"threshold":  event.Threshold,
		"message":    event.Message,
	}).Info("[ALERT] " + event.Message)
	return nil
}

// alertHandlers 是所有已註冊的 AlertEventHandler 清單
var alertHandlers []eventbus.AlertEventHandler

// RegisterAlertHandler 用於讓 plugin 模組註冊自訂的告警事件處理器。
// zh: 提供 plugin 自動註冊機制，會將處理器加入全域列表。
func RegisterAlertHandler(handler eventbus.AlertEventHandler) {
	alertHandlers = append(alertHandlers, handler)
}

// LoadPluginAlertHandlers 回傳目前已註冊的 plugin AlertEventHandler 清單。
// zh: 在註冊器中載入 plugin 註冊的所有 alert handler。
func LoadPluginAlertHandlers() []eventbus.AlertEventHandler {
	return alertHandlers
}



================================================
FILE: internal/adapters/eventbus/host.go
================================================
package eventbusadapter

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// HostLoggerHandler is a sample implementation of HostEventHandler that logs discovered hosts.
// zh: HostLoggerHandler 是接收到主機註冊事件時記錄日誌的處理器實作範例。
type HostLoggerHandler struct{}

// HandleHostDiscovered handles HostDiscoveredEvent and logs host information.
// zh: 接收到主機註冊事件後，透過 logger 模組輸出結構化日誌。
func (h *HostLoggerHandler) HandleHostDiscovered(ctx context.Context, event event.HostDiscoveredEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"name":   event.Name,
		"source": event.Source,
		"labels": event.Labels,
	}).Info("[HOST] discovered")
	return nil
}

// hostHandlers 是所有已註冊的 HostEventHandler 清單
var hostHandlers []event.HostEventHandler

// RegisterHostHandler 用於讓 plugin 模組註冊自訂的主機事件處理器。
// zh: 提供 plugin 自動註冊機制，會將處理器加入全域列表。
func RegisterHostHandler(handler event.HostEventHandler) {
	hostHandlers = append(hostHandlers, handler)
}

// LoadPluginHostHandlers 回傳目前已註冊的 plugin HostEventHandler 清單。
// zh: 在註冊器中載入 plugin 註冊的所有 host handler。
func LoadPluginHostHandlers() []event.HostEventHandler {
	return hostHandlers
}



================================================
FILE: internal/adapters/eventbus/inmemory.go
================================================
package eventbusadapter

import (
	"context"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/eventbus"
)

// InMemoryDispatcher provides a thread-safe, in-process implementation of the EventDispatcher interface.
// zh: InMemoryDispatcher 提供執行緒安全、僅在記憶體中運作的事件分派器。
type InMemoryDispatcher struct {
	alertHandlers  []event.AlertEventHandler
	taskHandlers   []event.TaskEventHandler
	hostHandlers   []event.HostEventHandler
	metricHandlers []event.MetricEventHandler
	mu             sync.RWMutex
}

// NewInMemoryDispatcher creates a new in-memory event dispatcher instance.
// zh: 建立新的記憶體內事件分派器。
func NewInMemoryDispatcher() eventbus.EventDispatcher {
	return &InMemoryDispatcher{}
}

// DispatchAlertTriggered dispatches AlertTriggeredEvent to all registered alert handlers.
// zh: 將 AlertTriggeredEvent 傳遞給所有已註冊的告警處理器。
func (d *InMemoryDispatcher) DispatchAlertTriggered(ctx context.Context, e event.AlertTriggeredEvent) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, h := range d.alertHandlers {
		if err := h.HandleAlertTriggered(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

// RegisterAlertHandler registers an alert event handler.
// zh: 註冊 Alert 事件處理器。
func (d *InMemoryDispatcher) RegisterAlertHandler(h event.AlertEventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.alertHandlers = append(d.alertHandlers, h)
}

// DispatchTaskCompleted dispatches TaskCompletedEvent to all task handlers.
// zh: 傳遞 TaskCompletedEvent 給所有註冊的任務處理器。
func (d *InMemoryDispatcher) DispatchTaskCompleted(ctx context.Context, e event.TaskCompletedEvent) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, h := range d.taskHandlers {
		if err := h.HandleTaskCompleted(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

// RegisterTaskHandler registers a task event handler.
// zh: 註冊 Task 事件處理器。
func (d *InMemoryDispatcher) RegisterTaskHandler(h event.TaskEventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.taskHandlers = append(d.taskHandlers, h)
}

// DispatchHostDiscovered dispatches HostDiscoveredEvent to all host handlers.
// zh: 傳遞 HostDiscoveredEvent 給所有註冊的主機處理器。
func (d *InMemoryDispatcher) DispatchHostDiscovered(ctx context.Context, e event.HostDiscoveredEvent) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, h := range d.hostHandlers {
		if err := h.HandleHostDiscovered(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

// RegisterHostHandler registers a host event handler.
// zh: 註冊 Host 事件處理器。
func (d *InMemoryDispatcher) RegisterHostHandler(h event.HostEventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.hostHandlers = append(d.hostHandlers, h)
}

// DispatchMetricOverflow dispatches MetricOverflowEvent to all metric handlers.
// zh: 傳遞 MetricOverflowEvent 給所有註冊的指標處理器。
func (d *InMemoryDispatcher) DispatchMetricOverflow(ctx context.Context, e event.MetricOverflowEvent) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	for _, h := range d.metricHandlers {
		if err := h.HandleMetricOverflow(ctx, e); err != nil {
			return err
		}
	}
	return nil
}

// RegisterMetricHandler registers a metric event handler.
// zh: 註冊 Metric 事件處理器。
func (d *InMemoryDispatcher) RegisterMetricHandler(h event.MetricEventHandler) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.metricHandlers = append(d.metricHandlers, h)
}



================================================
FILE: internal/adapters/eventbus/metric.go
================================================
package eventbusadapter

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// MetricLoggerHandler is a sample implementation of MetricEventHandler that logs overflow events.
// zh: MetricLoggerHandler 是接收到指標溢出事件時記錄日誌的處理器實作範例。
type MetricLoggerHandler struct{}

// HandleMetricOverflow handles MetricOverflowEvent by logging the event with structured fields.
// zh: 接收到指標溢出事件後，使用結構化欄位輸出警告日誌
func (h *MetricLoggerHandler) HandleMetricOverflow(ctx context.Context, event event.MetricOverflowEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"metric":    event.MetricName,
		"value":     event.Value,
		"threshold": event.Threshold,
		"instance":  event.Instance,
		"reason":    event.Reason,
	}).Warn("[METRIC] overflow detected")
	return nil
}

// metricHandlers 是所有已註冊的 MetricEventHandler 清單
var metricHandlers []event.MetricEventHandler

// RegisterMetricHandler 用於讓 plugin 模組註冊自訂的指標事件處理器。
// zh: 提供 plugin 自動註冊機制，會將處理器加入全域列表。
func RegisterMetricHandler(handler event.MetricEventHandler) {
	metricHandlers = append(metricHandlers, handler)
}

// LoadPluginMetricHandlers 回傳目前已註冊的 plugin MetricEventHandler 清單。
// zh: 在註冊器中載入 plugin 註冊的所有 metric handler。
func LoadPluginMetricHandlers() []event.MetricEventHandler {
	return metricHandlers
}



================================================
FILE: internal/adapters/eventbus/task.go
================================================
package eventbusadapter

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// TaskLoggerHandler is a sample implementation of TaskEventHandler that logs task results.
// zh: TaskLoggerHandler 是接收到任務完成事件時記錄日誌的處理器實作範例。
type TaskLoggerHandler struct{}

// HandleTaskCompleted handles task completion and logs the result.
// zh: 接收到任務完成事件後，透過 logger 模組輸出結構化日誌。
func (h *TaskLoggerHandler) HandleTaskCompleted(ctx context.Context, event event.TaskCompletedEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"task_id":   event.TaskID,
		"worker_id": event.WorkerID,
		"status":    event.Status,
	}).Info("[TASK] completed")
	return nil
}

// taskHandlers 是所有已註冊的 TaskEventHandler 清單
var taskHandlers []event.TaskEventHandler

// RegisterTaskHandler 用於讓 plugin 模組註冊自訂的任務事件處理器。
// zh: 提供 plugin 自動註冊機制，會將處理器加入全域列表。
func RegisterTaskHandler(handler event.TaskEventHandler) {
	taskHandlers = append(taskHandlers, handler)
}

// LoadPluginTaskHandlers 回傳目前已註冊的 plugin TaskEventHandler 清單。
// zh: 在註冊器中載入 plugin 註冊的所有 task handler。
func LoadPluginTaskHandlers() []event.TaskEventHandler {
	return taskHandlers
}



================================================
FILE: internal/adapters/logger/logger_test.go
================================================
package loggeradapter_test

import (
	"context"
	"testing"

	loggeradapter "github.com/detectviz/detectviz/internal/adapters/logger"
	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestZapLogger_Info(t *testing.T) {
	core, observedLogs := observer.New(zapcore.InfoLevel)
	baseLogger := zap.New(core) // 正確建立 zap.Logger
	sugar := baseLogger.Sugar()

	zapLogger := loggeradapter.NewZapLogger(sugar)
	zapLogger.Info("test message", "key", "value")

	if observedLogs.Len() != 1 {
		t.Fatalf("expected 1 log entry, got %d", observedLogs.Len())
	}

	entry := observedLogs.All()[0]
	if entry.Message != "test message" {
		t.Errorf("unexpected log message: %s", entry.Message)
	}
	if val, ok := entry.ContextMap()["key"]; !ok || val != "value" {
		t.Errorf("missing or incorrect field value: %v", entry.ContextMap())
	}
}

func TestNopLogger_NoPanic(t *testing.T) {
	var log ifacelogger.Logger = loggeradapter.NewNopLogger()

	// All calls should be no-op
	log.Debug("debug")
	log.Info("info", "x", 1)
	log.Warn("warn")
	log.Error("error", "y", 2)
	log.Sync()
	log2 := log.WithFields(map[string]any{"test": true})
	log2 = log2.WithContext(context.Background())
	log2 = log2.Named("test")

	if log2 == nil {
		t.Error("NopLogger should return non-nil instance")
	}
}



================================================
FILE: internal/adapters/logger/nop_adapter.go
================================================
package loggeradapter

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// NewNopLogger returns a no-op logger implementation.
// zh: 回傳一個不會輸出任何內容的 logger，常用於測試或作為預設 fallback。
func NewNopLogger() logger.Logger {
	return &NopLogger{}
}

// NopLogger 實作 logger.Logger 但不輸出任何 log。
// zh: 此實作常用於測試或禁用 log 時。
type NopLogger struct{}

func (n *NopLogger) Info(msg string, args ...any)  {}
func (n *NopLogger) Warn(msg string, args ...any)  {}
func (n *NopLogger) Error(msg string, args ...any) {}
func (n *NopLogger) Debug(msg string, args ...any) {}

func (n *NopLogger) Named(name string) logger.Logger {
	return n
}

func (n *NopLogger) WithContext(ctx context.Context) logger.Logger {
	return n
}

func (n *NopLogger) WithFields(fields map[string]any) logger.Logger {
	return n
}

func (n *NopLogger) Sync() error {
	return nil
}



================================================
FILE: internal/adapters/logger/zap_adapter.go
================================================
package loggeradapter

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
)

type ZapLogger struct {
	l *zap.SugaredLogger
}

// NewZapLogger 建立新的 ZapLogger，接收 zap.SugaredLogger 實例。
// zh: 建立以 zap.SugaredLogger 為基礎的日誌介面實作。
func NewZapLogger(base *zap.SugaredLogger) ifacelogger.Logger {
	return &ZapLogger{l: base}
}

// Info logs a message at InfoLevel.
// zh: 以 InfoLevel 輸出日誌訊息（可選格式化參數）。
func (z *ZapLogger) Info(msg string, args ...any) {
	z.l.Infow(msg, args...)
}

// Warn logs a message at WarnLevel.
// zh: 以 WarnLevel 輸出日誌訊息（可選格式化參數）。
func (z *ZapLogger) Warn(msg string, args ...any) {
	z.l.Warnw(msg, args...)
}

// Error logs a message at ErrorLevel.
// zh: 以 ErrorLevel 輸出日誌訊息（可選格式化參數）。
func (z *ZapLogger) Error(msg string, args ...any) {
	z.l.Errorw(msg, args...)
}

// Debug logs a debug-level message with optional formatting.
// zh: 以 DebugLevel 輸出日誌訊息（可選格式化參數）。
func (z *ZapLogger) Debug(msg string, args ...any) {
	z.l.Debugw(msg, args...)
}

func (z *ZapLogger) Sync() error {
	return z.l.Sync()
}

func (z *ZapLogger) WithFields(fields map[string]any) ifacelogger.Logger {
	return &ZapLogger{l: z.l.With(fields)}
}

func (z *ZapLogger) WithContext(ctx context.Context) ifacelogger.Logger {
	// zh: 可選擇從 context 擷取 trace_id/span_id 並加到欄位
	// 這裡暫以原 logger 回傳
	return z
}

func (z *ZapLogger) Named(name string) ifacelogger.Logger {
	return &ZapLogger{l: z.l.Named(name)}
}

// NewDefaultZap 建立預設設定的 zap.Logger 實例。
// zh: 提供簡易用於本地測試與開發的 logger 初始化方法。
func NewDefaultZap() *zap.SugaredLogger {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)

	logger, _ := cfg.Build()
	return logger.Sugar()
}



================================================
FILE: internal/adapters/metrics/aggregator.go
================================================
package metricsadapter

import (
	"errors"

	metric "github.com/detectviz/detectviz/pkg/ifaces/metrics"
)

// AggregationType defines the available aggregation methods.
// zh: AggregationType 定義可用的統計聚合類型。
type AggregationType string

const (
	SumAggregation     AggregationType = "sum"
	AverageAggregation AggregationType = "avg"
	MaxAggregation     AggregationType = "max"
	MinAggregation     AggregationType = "min"
)

// SimpleAggregator provides basic aggregation logic over time series.
// zh: SimpleAggregator 提供對時間序列資料的基本統計運算邏輯。
type SimpleAggregator struct{}

// Aggregate performs the specified aggregation on a list of TimePoints.
// zh: 根據指定類型對一組時間點進行聚合計算。
func (a *SimpleAggregator) Aggregate(points []metric.TimePoint, aggType AggregationType) (float64, error) {
	if len(points) == 0 {
		return 0, errors.New("no points to aggregate")
	}

	switch aggType {
	case SumAggregation:
		var sum float64
		for _, p := range points {
			sum += p.Value
		}
		return sum, nil

	case AverageAggregation:
		var sum float64
		for _, p := range points {
			sum += p.Value
		}
		return sum / float64(len(points)), nil

	case MaxAggregation:
		max := points[0].Value
		for _, p := range points {
			if p.Value > max {
				max = p.Value
			}
		}
		return max, nil

	case MinAggregation:
		min := points[0].Value
		for _, p := range points {
			if p.Value < min {
				min = p.Value
			}
		}
		return min, nil

	default:
		return 0, errors.New("unsupported aggregation type")
	}
}



================================================
FILE: internal/adapters/metrics/query_adapter.go
================================================
package metricsadapter

import (
	"context"
	"errors"
)

// FluxQueryAdapter implements MetricQueryAdapter for querying InfluxDB via Flux.
// zh: FluxQueryAdapter 是透過 Flux 查詢 InfluxDB 的 MetricQueryAdapter 實作。
type FluxQueryAdapter struct {
	// TODO: 注入 InfluxDB client 或查詢執行器
	// Client influxdb2.Client
}

func (f *FluxQueryAdapter) Query(ctx context.Context, expr string, labels map[string]string) (float64, error) {
	// TODO: 根據 expr 與 labels 組裝 Flux 語法
	// TODO: 呼叫 InfluxDB 查詢 API 並解析回傳資料
	// TODO: 回傳查詢結果中的第一筆數值或適當的錯誤
	return 0, errors.New("flux query not implemented")
}

// PromQueryAdapter implements metric.MetricQueryAdapter for querying Prometheus using PromQL.
// PromQueryAdapter 是使用 PromQL 查詢 Prometheus 的 metric.MetricQueryAdapter 實作。
type PromQueryAdapter struct {
	// TODO: 注入 Prometheus 查詢 client，例如 prometheus.Client
	// Client *promapi.Client
}

func (p *PromQueryAdapter) Query(ctx context.Context, expr string, labels map[string]string) (float64, error) {
	// TODO: 根據 expr 與 labels 組裝 PromQL 語法
	// TODO: 呼叫 Prometheus API 並解析回傳資料
	// TODO: 回傳查詢結果中的第一筆數值或適當的錯誤
	return 0, errors.New("prometheus query not implemented")
}

// MockQueryAdapter is a stub implementation of MetricQueryAdapter for testing.
// zh: MockQueryAdapter 是 MetricQueryAdapter 的測試用假實作，回傳固定值。
type MockQueryAdapter struct {
	FixedValue float64
	Err        error
}

func (m *MockQueryAdapter) Query(ctx context.Context, expr string, labels map[string]string) (float64, error) {
	if m.Err != nil {
		return 0, m.Err
	}
	return m.FixedValue, nil
}



================================================
FILE: internal/adapters/metrics/series_reader_adapter.go
================================================
package metricsadapter

import (
	"context"
	"errors"
	"time"

	metric "github.com/detectviz/detectviz/pkg/ifaces/metrics"
)

// MockSeriesReader implements MetricSeriesReader with dummy time series data.
// zh: MockSeriesReader 是 MetricSeriesReader 的模擬實作，用於回傳虛擬的時間序列資料。
type MockSeriesReader struct {
	Points []metric.TimePoint // zh: 預設回傳的資料點清單
	Err    error              // zh: 若設定錯誤，則每次查詢都會回傳此錯誤
}

// ReadSeries returns a sequence of metric points based on the configured mock data.
// zh: 根據預設資料回傳一段時間序列資料，若未設定則產生 5 筆模擬資料。
func (m *MockSeriesReader) ReadSeries(ctx context.Context, expr string, labels map[string]string, start, end int64) ([]metric.TimePoint, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if len(m.Points) > 0 {
		return m.Points, nil
	}
	// Default mock: generate 5 points, spaced 1 minute apart
	var points []metric.TimePoint
	t := time.Now().Unix()
	for i := 0; i < 5; i++ {
		points = append(points, metric.TimePoint{
			Timestamp: t - int64(60*i),
			Value:     float64(i),
		})
	}
	return points, nil
}

// InfluxSeriesReader implements MetricSeriesReader for reading time series from InfluxDB.
// zh: InfluxSeriesReader 是用於從 InfluxDB 讀取時間序列資料的實作。
type InfluxSeriesReader struct {
	// TODO: 注入 InfluxDB 的查詢 client，例如 influxdb2.QueryAPI
	// Client influxdb2.QueryAPI
}

// ReadSeries reads a sequence of metric points from InfluxDB using Flux.
// zh: 從 InfluxDB 讀取指定條件的時間序列資料，並轉換為 TimePoint 格式回傳。
func (r *InfluxSeriesReader) ReadSeries(ctx context.Context, expr string, labels map[string]string, start, end int64) ([]metric.TimePoint, error) {
	// TODO: 組裝 flux 查詢語法
	// TODO: 執行查詢，解析查詢結果為 TimePoint 陣列
	return nil, errors.New("influx series reader not implemented")
}



================================================
FILE: internal/adapters/metrics/transformer_adapter.go
================================================
package metricsadapter

// NoopTransformer implements MetricTransformer without any modification.
// zh: NoopTransformer 是不進行任何轉換的預設實作。
type NoopTransformer struct{}

// Transform fulfills the MetricTransformer interface without changing input.
// zh: 不進行任何修改，原樣傳回
func (t *NoopTransformer) Transform(measurement *string, value *float64, labels map[string]string) error {
	return nil
}



================================================
FILE: internal/adapters/metrics/writer_adapter.go
================================================
package metricsadapter

import (
	"context"
	"errors"
)

// InfluxMetricWriter implements MetricWriter for sending metrics to InfluxDB.
// zh: InfluxMetricWriter 是將指標寫入 InfluxDB 的 MetricWriter 實作。
type InfluxMetricWriter struct {
	// TODO: 注入 InfluxDB 寫入 client
	// Client influxdb2.WriteAPI
}

func (w *InfluxMetricWriter) WritePoint(ctx context.Context, measurement string, value float64, labels map[string]string) error {
	// TODO: 將 measurement、value、labels 組裝成 Point 資料
	// TODO: 呼叫 InfluxDB client 寫入方法
	return errors.New("influx write not implemented")
}

// PushgatewayMetricWriter implements MetricWriter for pushing metrics to Prometheus Pushgateway.
// zh: PushgatewayMetricWriter 是將指標推送至 Prometheus Pushgateway 的 MetricWriter 實作。
type PushgatewayMetricWriter struct {
	// TODO: 注入 Pushgateway HTTP client 或設定參數
	// Endpoint string
}

func (w *PushgatewayMetricWriter) WritePoint(ctx context.Context, measurement string, value float64, labels map[string]string) error {
	// TODO: 組裝成 Pushgateway 支援的格式並發送 HTTP 請求
	return errors.New("pushgateway write not implemented")
}

// MockMetricWriter is a mock implementation of MetricWriter for testing.
// zh: MockMetricWriter 是測試用的 MetricWriter 實作，可記錄寫入行為或模擬錯誤。
type MockMetricWriter struct {
	LastMeasurement string
	LastValue       float64
	LastLabels      map[string]string
	Err             error
}

func (m *MockMetricWriter) WritePoint(ctx context.Context, measurement string, value float64, labels map[string]string) error {
	if m.Err != nil {
		return m.Err
	}
	m.LastMeasurement = measurement
	m.LastValue = value
	m.LastLabels = labels
	return nil
}



================================================
FILE: internal/adapters/modules/engine_adapter.go
================================================
package modulesadapter

import (
	"context"

	core "github.com/detectviz/detectviz/internal/modules"
	iface "github.com/detectviz/detectviz/pkg/ifaces/modules"
)

// EngineAdapter wraps core.Engine to implement iface.ModuleEngine.
// zh: EngineAdapter 包裝 core.Engine，使其實作 ModuleEngine 介面。
type EngineAdapter struct {
	engine *core.Engine
}

// NewEngineAdapter constructs a new EngineAdapter instance.
// zh: 建立新的 EngineAdapter 實例。
func NewEngineAdapter(e *core.Engine) *EngineAdapter {
	return &EngineAdapter{
		engine: e,
	}
}

// Register adds a module to the core engine.
// zh: 註冊模組至底層 Engine 實例。
func (a *EngineAdapter) Register(m iface.LifecycleModule) {
	a.engine.Register(m)
}

// RunAll starts all registered modules via the core engine.
// zh: 啟動所有已註冊的模組。
func (a *EngineAdapter) RunAll(ctx context.Context) error {
	return a.engine.RunAll(ctx)
}

// ShutdownAll stops all modules via the core engine.
// zh: 關閉所有已註冊的模組。
func (a *EngineAdapter) ShutdownAll(ctx context.Context) error {
	return a.engine.ShutdownAll(ctx)
}



================================================
FILE: internal/adapters/modules/listener_adapter.go
================================================
package modulesadapter

import (
	"context"

	core "github.com/detectviz/detectviz/internal/modules"
)

// ListenerAdapter wraps core.Listener to implement iface.ModuleListener.
// zh: ListenerAdapter 包裝 core.Listener，使其實作 ModuleListener 介面。
type ListenerAdapter struct {
	listener *core.Listener
}

// NewListenerAdapter constructs a new ListenerAdapter instance.
// zh: 建立新的 ListenerAdapter 實例。
func NewListenerAdapter(l *core.Listener) *ListenerAdapter {
	return &ListenerAdapter{
		listener: l,
	}
}

// Start begins the health monitoring loop.
// zh: 啟動健康狀態監控迴圈。
func (a *ListenerAdapter) Start(ctx context.Context) {
	a.listener.Start(ctx)
}

// Stop stops the health monitoring listener.
// zh: 停止健康狀態監聽器。
func (a *ListenerAdapter) Stop() {
	a.listener.Stop()
}



================================================
FILE: internal/adapters/modules/registry_adapter.go
================================================
package modulesadapter

import (
	core "github.com/detectviz/detectviz/internal/modules"
	iface "github.com/detectviz/detectviz/pkg/ifaces/modules"
)

// RegistryAdapter wraps core.Registry to implement iface.ModuleRegistry.
// zh: RegistryAdapter 包裝 core.Registry，使其實作 ModuleRegistry 介面。
type RegistryAdapter struct {
	registry *core.Registry
}

// NewRegistryAdapter constructs a new RegistryAdapter instance.
// zh: 建立新的 RegistryAdapter 實例。
func NewRegistryAdapter(r *core.Registry) *RegistryAdapter {
	return &RegistryAdapter{
		registry: r,
	}
}

// Register adds a named module to the registry.
// zh: 註冊具名模組至底層 registry。
func (a *RegistryAdapter) Register(name string, m iface.LifecycleModule) error {
	return a.registry.Register(name, m)
}

// Get retrieves a registered module by name.
// zh: 根據名稱查詢模組。
func (a *RegistryAdapter) Get(name string) (iface.LifecycleModule, bool) {
	return a.registry.Get(name)
}

// List returns all registered module names.
// zh: 回傳所有已註冊模組名稱。
func (a *RegistryAdapter) List() []string {
	return a.registry.List()
}



================================================
FILE: internal/adapters/modules/runner_adapter.go
================================================
package modulesadapter

import (
	"context"

	core "github.com/detectviz/detectviz/internal/modules"
)

// RunnerAdapter wraps core.Runner to implement iface.ModuleRunner.
// zh: RunnerAdapter 包裝 core.Runner，使其實作 ModuleRunner 介面。
type RunnerAdapter struct {
	runner *core.Runner
}

// NewRunnerAdapter constructs a new RunnerAdapter instance.
// zh: 建立新的 RunnerAdapter 實例。
func NewRunnerAdapter(r *core.Runner) *RunnerAdapter {
	return &RunnerAdapter{
		runner: r,
	}
}

// StartAll starts all modules based on dependency order.
// zh: 依據依賴關係啟動所有模組。
func (a *RunnerAdapter) StartAll(ctx context.Context) error {
	return a.runner.StartAll(ctx)
}

// StopAll stops all modules in reverse order.
// zh: 依照啟動順序反向關閉所有模組。
func (a *RunnerAdapter) StopAll(ctx context.Context) error {
	return a.runner.StopAll(ctx)
}



================================================
FILE: internal/adapters/notifier/email_adapter.go
================================================
package notifieradapter

import (
	"context"

	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
	ifacenotifier "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// EmailNotifier sends notifications via email.
// zh: EmailNotifier 負責透過電子郵件傳送通知。
type EmailNotifier struct {
	name   string
	sender string
	logger ifacelogger.Logger
}

// NewEmailNotifier returns a new instance of EmailNotifier.
// zh: 建立新的 EmailNotifier 實例。
func NewEmailNotifier(name string, sender string, logger ifacelogger.Logger) *EmailNotifier {
	return &EmailNotifier{
		name:   name,
		sender: sender,
		logger: logger,
	}
}

// Name returns the notifier name.
// zh: 回傳通知器名稱。
func (n *EmailNotifier) Name() string {
	return n.name
}

// Send sends the message as an email.
// zh: 傳送通知訊息為 email（尚未實作寄信邏輯）。
func (n *EmailNotifier) Send(ctx context.Context, msg ifacenotifier.Message) error {
	n.logger.WithContext(ctx).Info("EmailNotifier sending",
		"to", msg.Target,
		"title", msg.Title,
		"content", msg.Content,
	)
	// TODO: implement real email sending logic using SMTP or third-party API
	return nil
}

// Notify implements the Notifier interface for simple notification.
// zh: 將簡易 title/message 組裝為完整訊息後透過 Send 傳送。
func (n *EmailNotifier) Notify(title, message string) error {
	// TODO: 實際應從 config 或預設值決定 target
	msg := ifacenotifier.Message{
		Target:  "default@example.com",
		Title:   title,
		Content: message,
	}
	return n.Send(context.Background(), msg)
}



================================================
FILE: internal/adapters/notifier/mock_adapter.go
================================================
package notifieradapter

import (
	"context"
	"sync"

	ifacenotifier "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// MockNotifier is a test implementation of the Notifier interface.
// zh: MockNotifier 為測試用通知器，會記錄所有發送的訊息。
type MockNotifier struct {
	name     string
	messages []ifacenotifier.Message
	mu       sync.Mutex
}

// NewMockNotifier creates a new MockNotifier.
// zh: 建立新的 MockNotifier 實例。
func NewMockNotifier(name string) *MockNotifier {
	return &MockNotifier{name: name}
}

// Name returns the notifier name.
// zh: 回傳通知器名稱。
func (n *MockNotifier) Name() string {
	return n.name
}

// Send appends the message to internal buffer.
// zh: 將通知訊息儲存至內部緩衝區，供測試驗證。
func (n *MockNotifier) Send(ctx context.Context, msg ifacenotifier.Message) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.messages = append(n.messages, msg)
	return nil
}

// Messages returns all messages sent.
// zh: 回傳所有已發送訊息。
func (n *MockNotifier) Messages() []ifacenotifier.Message {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]ifacenotifier.Message(nil), n.messages...)
}



================================================
FILE: internal/adapters/notifier/multi.go
================================================
package notifieradapter

import (
	"context"

	notifieriface "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// MultiNotifier 將多個 Notifier 組合成一個 Notifier。
// zh: 每個通知請求會依序傳遞至所有註冊的 Notifier。
type MultiNotifier struct {
	notifiers []notifieriface.Notifier
}

// NewMultiNotifier 建立 MultiNotifier 實例。
// zh: 可接受多個 Notifier 作為參數，並整合為一個執行單元。
func NewMultiNotifier(list ...notifieriface.Notifier) notifieriface.Notifier {
	return &MultiNotifier{notifiers: list}
}

// Notify 傳送通知，會依序呼叫每個註冊的 Notifier。
// zh: 若某些 Notifier 回傳錯誤，最終僅回傳最後一個錯誤。
func (m *MultiNotifier) Notify(title, message string) error {
	var lastErr error
	for _, n := range m.notifiers {
		if err := n.Notify(title, message); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Send 傳送完整通知訊息，依序傳遞至每個 Notifier。
// zh: 支援包含標籤、時間等欄位的複雜訊息傳送。
func (m *MultiNotifier) Send(ctx context.Context, msg notifieriface.Message) error {
	var lastErr error
	for _, n := range m.notifiers {
		if err := n.Send(ctx, msg); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Name 回傳 Notifier 名稱。
// zh: 表示此組合型通知器的識別名稱。
func (m *MultiNotifier) Name() string {
	return "multi"
}



================================================
FILE: internal/adapters/notifier/nop.go
================================================
package notifieradapter

import (
	"context"

	notifieriface "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// NopNotifier 是一個不執行任何操作的通知器。
// zh: 常用於測試或當無通知功能需求時作為 fallback 實作。
type NopNotifier struct{}

// NewNopNotifier 建立 NopNotifier 實例。
// zh: 回傳不執行通知邏輯的 Notifier 實作。
func NewNopNotifier() notifieriface.Notifier {
	return &NopNotifier{}
}

// Name 回傳通知器名稱。
// zh: 標示此 notifier 為 nop 類型。
func (n *NopNotifier) Name() string {
	return "nop"
}

// Send 實作完整通知方法但不執行任何動作。
// zh: 忽略訊息內容並回傳 nil，適用於測試或預設無通報。
func (n *NopNotifier) Send(ctx context.Context, msg notifieriface.Message) error {
	return nil
}

// Notify 實作簡易通知方法但不執行任何動作。
// zh: 忽略標題與訊息並回傳 nil，適用於測試或不需通知的情境。
func (n *NopNotifier) Notify(title, message string) error {
	return nil
}



================================================
FILE: internal/adapters/notifier/slack_adapter.go
================================================
package notifieradapter

import (
	"context"

	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
	ifacenotifier "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// SlackNotifier sends notifications to Slack.
// zh: SlackNotifier 負責傳送通知至 Slack。
type SlackNotifier struct {
	name       string
	webhookURL string
	logger     ifacelogger.Logger
}

// NewSlackNotifier returns a new instance of SlackNotifier.
// zh: 建立新的 SlackNotifier 實例。
func NewSlackNotifier(name string, webhookURL string, logger ifacelogger.Logger) *SlackNotifier {
	return &SlackNotifier{
		name:       name,
		webhookURL: webhookURL,
		logger:     logger,
	}
}

// Name returns the notifier name.
// zh: 回傳通知器名稱。
func (n *SlackNotifier) Name() string {
	return n.name
}

// Send sends the message to Slack.
// zh: 傳送通知訊息至 Slack（尚未實作實際 API 呼叫）。
func (n *SlackNotifier) Send(ctx context.Context, msg ifacenotifier.Message) error {
	n.logger.WithContext(ctx).Info("SlackNotifier sending",
		"to", msg.Target,
		"title", msg.Title,
		"content", msg.Content,
	)
	// TODO: implement real Slack webhook logic
	return nil
}

// Notify implements simplified notification with title and message.
// zh: 使用簡易標題與訊息格式發送 Slack 通知。
func (n *SlackNotifier) Notify(title, message string) error {
	// TODO: 決定預設或從設定注入 Target
	msg := ifacenotifier.Message{
		Target:  "default-slack-channel",
		Title:   title,
		Content: message,
	}
	return n.Send(context.Background(), msg)
}



================================================
FILE: internal/adapters/notifier/webhook_adapter.go
================================================
package notifieradapter

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
	ifacenotifier "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// WebhookNotifier sends notifications to a specified webhook URL.
// zh: WebhookNotifier 負責將通知訊息送往指定的 webhook URL。
type WebhookNotifier struct {
	name   string
	client *http.Client
	logger ifacelogger.Logger
}

// NewWebhookNotifier creates a new WebhookNotifier.
// zh: 建立新的 WebhookNotifier 實例。
// 若未提供 client 則使用 http.DefaultClient。
func NewWebhookNotifier(name string, logger ifacelogger.Logger, client *http.Client) *WebhookNotifier {
	if client == nil {
		client = http.DefaultClient
	}
	return &WebhookNotifier{
		name:   name,
		client: client,
		logger: logger,
	}
}

// Name returns the notifier name.
// zh: 回傳 notifier 名稱。
func (n *WebhookNotifier) Name() string {
	return n.name
}

// Send transmits the message as a JSON payload via HTTP POST.
// zh: 傳送訊息為 JSON 格式，透過 HTTP POST 傳送至 msg.Target。
func (n *WebhookNotifier) Send(ctx context.Context, msg ifacenotifier.Message) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		n.logger.WithContext(ctx).Error("Failed to marshal webhook message", "error", err)
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, msg.Target, bytes.NewBuffer(payload))
	if err != nil {
		n.logger.WithContext(ctx).Error("Failed to create webhook request", "error", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		n.logger.WithContext(ctx).Error("Failed to send webhook request", "error", err)
		return err
	}
	defer resp.Body.Close()

	n.logger.WithContext(ctx).Info("WebhookNotifier sent", "target", msg.Target, "status", resp.Status)
	return nil
}

// Notify sends a simple title-message notification.
// zh: 傳送簡易的標題與訊息格式至 webhook，將組成 JSON 輸出。
func (n *WebhookNotifier) Notify(title, message string) error {
	// TODO: 後續改為可設定或多通道支援
	msg := ifacenotifier.Message{
		Target:  "http://localhost:8080/webhook",
		Title:   title,
		Content: message,
	}
	return n.Send(context.Background(), msg)
}



================================================
FILE: internal/adapters/scheduler/cron_adapter.go
================================================
package scheduleradapter

import (
	"context"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"github.com/detectviz/detectviz/pkg/ifaces/scheduler"
	"github.com/robfig/cron/v3"
)

// CronScheduler is an implementation of the Scheduler interface using robfig/cron.
// zh: CronScheduler 是使用 robfig/cron 套件實作的排程器。
type CronScheduler struct {
	c     *cron.Cron
	jobs  []scheduler.Job
	mu    sync.Mutex
	start sync.Once
	log   logger.Logger
}

// NewCronScheduler creates a new CronScheduler instance with logger support.
// zh: 建立一個新的 CronScheduler 實例，支援注入 logger。
func NewCronScheduler(log logger.Logger) *CronScheduler {
	return &CronScheduler{
		c:   cron.New(),
		log: log,
	}
}

// Register adds a job to the scheduler. The job must implement Spec() string to define its schedule.
// zh: 註冊一個任務到排程器，排程時間由 job.Spec() 提供。
func (s *CronScheduler) Register(job scheduler.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs = append(s.jobs, job)
}

// Start schedules all registered jobs using their Spec() value.
// zh: 根據每個註冊任務的 Spec() 結果進行排程並啟動排程器。
func (s *CronScheduler) Start(ctx context.Context) error {
	s.start.Do(func() {
		for _, j := range s.jobs {
			spec := j.Spec()
			_, err := s.c.AddFunc(spec, func() {
				if err := j.Run(ctx); err != nil && s.log != nil {
					s.log.Error("job run failed", "name", j.Name(), "error", err)
				}
			})
			if err != nil && s.log != nil {
				s.log.Error("failed to add cron job", "spec", spec, "name", j.Name(), "error", err)
			}
		}
		s.c.Start()
	})
	return nil
}

// Stop stops the cron engine gracefully.
// zh: 優雅地停止排程引擎。
func (s *CronScheduler) Stop(ctx context.Context) error {
	s.c.Stop()
	return nil
}

/*
範例：如何使用 CronScheduler 註冊任務

log := logger.NewZapLogger(...)
sched := NewCronScheduler(log)

job := &MyJob{} // 實作 scheduler.Job 介面，並實作 Spec() string 方法
sched.Register(job)

ctx := context.Background()
_ = sched.Start(ctx)
*/



================================================
FILE: internal/adapters/scheduler/cron_adapter_test.go
================================================
package scheduleradapter_test

import (
	"context"
	"testing"
	"time"

	scheduleradapter "github.com/detectviz/detectviz/internal/adapters/scheduler"
	"github.com/detectviz/detectviz/internal/test/testutil"
	"github.com/stretchr/testify/assert"
)

// mockJob implements scheduler.Job interface.
// zh: 用於測試 CronScheduler 的模擬任務。
type mockJob struct {
	name  string
	spec  string
	calls int
}

func (m *mockJob) Run(ctx context.Context) error {
	m.calls++
	return nil
}

func (m *mockJob) Name() string {
	return m.name
}

func (m *mockJob) Spec() string {
	return m.spec
}

// TestCronScheduler_Run 驗證 CronScheduler 能夠依據時間規則週期執行註冊任務。
// zh: 測試 Cron 型排程器是否能依指定的頻率成功執行任務。
func TestCronScheduler_Run(t *testing.T) {
	ctx := context.Background()
	job := &mockJob{name: "cron-job", spec: "@every 2s"}

	logger := testutil.NewTestLogger()
	sched := scheduleradapter.NewCronScheduler(logger)
	sched.Register(job)

	err := sched.Start(ctx)
	assert.NoError(t, err)

	time.Sleep(5 * time.Second) // 等待任務至少執行兩次

	err = sched.Stop(ctx)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, job.calls, 2)
}



================================================
FILE: internal/adapters/scheduler/mock_adapter.go
================================================
package scheduleradapter

import (
	"context"
	"fmt"

	"github.com/detectviz/detectviz/pkg/ifaces/scheduler"
)

// MockJob is a simple mock implementation of a scheduled job.
// zh: MockJob 是模擬用的任務實作，用於測試排程器功能。
type MockJob struct {
	ID   string
	Logs *[]string // zh: 用於收集執行記錄
}

// Run appends a message to Logs when executed.
// zh: Run 在執行時將訊息寫入 Logs。
func (j *MockJob) Run(ctx context.Context) error {
	if j.Logs != nil {
		*j.Logs = append(*j.Logs, fmt.Sprintf("job %s executed", j.ID))
	}
	return nil
}

// Name returns the job ID.
// zh: 回傳任務名稱。
func (j *MockJob) Name() string {
	return j.ID
}

// MockScheduler is a no-op implementation for testing.
// zh: MockScheduler 是模擬用排程器實作，不會實際執行任務。
type MockScheduler struct {
	Jobs []scheduler.Job
}

// NewMockScheduler returns a new instance.
// zh: 建立一個新的 MockScheduler 實例。
func NewMockScheduler() *MockScheduler {
	return &MockScheduler{}
}

// Register adds a job to internal list.
// zh: 註冊任務至模擬排程器。
func (s *MockScheduler) Register(job scheduler.Job) {
	s.Jobs = append(s.Jobs, job)
}

// Start does nothing.
// zh: 模擬啟動，實際不執行任何任務。
func (s *MockScheduler) Start(ctx context.Context) error {
	return nil
}

// Stop does nothing.
// zh: 模擬停止，實際不執行任何操作。
func (s *MockScheduler) Stop(ctx context.Context) error {
	return nil
}



================================================
FILE: internal/adapters/scheduler/workerpool_adapter.go
================================================
package scheduleradapter

import (
	"context"
	"sync"
	"time"

	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"github.com/detectviz/detectviz/pkg/ifaces/scheduler"
)

// WorkerPoolScheduler executes jobs using a fixed pool of worker goroutines.
// zh: WorkerPoolScheduler 使用固定數量的 worker goroutine 執行排程任務。
type WorkerPoolScheduler struct {
	jobs       []scheduler.Job
	workerSize int
	started    bool
	mu         sync.Mutex
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	log        logger.Logger
}

// NewWorkerPoolScheduler creates a new scheduler with a given worker count and logger.
// zh: 建立一個新的 WorkerPoolScheduler，可指定 worker 數量與注入 logger。
func NewWorkerPoolScheduler(workerSize int, log logger.Logger) *WorkerPoolScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPoolScheduler{
		workerSize: workerSize,
		ctx:        ctx,
		cancel:     cancel,
		log:        log,
	}
}

// Register adds a job to the queue.
// zh: 註冊任務至排程器中。
func (s *WorkerPoolScheduler) Register(job scheduler.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs = append(s.jobs, job)
}

// Start runs the worker pool and schedules all jobs in round-robin.
// zh: 啟動 worker pool 並以輪詢方式執行所有註冊任務。
func (s *WorkerPoolScheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return nil
	}
	s.started = true
	s.mu.Unlock()

	for i := 0; i < s.workerSize; i++ {
		s.wg.Add(1)
		go s.worker()
	}
	return nil
}

// Stop cancels execution and waits for all workers to complete.
// zh: 停止排程器並等待所有 worker 結束。
func (s *WorkerPoolScheduler) Stop(ctx context.Context) error {
	s.cancel()
	s.wg.Wait()
	return nil
}

// worker is the function executed by each worker goroutine.
// zh: worker 是每個 worker goroutine 執行的函式。
func (s *WorkerPoolScheduler) worker() {
	defer s.wg.Done()
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.mu.Lock()
			jobs := make([]scheduler.Job, len(s.jobs))
			copy(jobs, s.jobs)
			s.mu.Unlock()
			for _, job := range jobs {
				go s.runWithRetry(job)
			}
		}
	}
}

// runWithRetry executes a job and retries up to 3 times if it fails.
// zh: 嘗試執行任務，若失敗則最多重試 3 次。
func (s *WorkerPoolScheduler) runWithRetry(job scheduler.Job) {
	const maxRetry = 3
	for i := 0; i < maxRetry; i++ {
		err := job.Run(s.ctx)
		if err == nil {
			return
		}
		s.log.Warn("job failed", "name", job.Name(), "attempt", i+1, "error", err)
		time.Sleep(1 * time.Second)
	}
	s.log.Error("job permanently failed after retries", "name", job.Name())
}



================================================
FILE: internal/adapters/scheduler/workerpool_adapter_test.go
================================================
package scheduleradapter_test

import (
	"context"
	"testing"
	"time"

	scheduleradapter "github.com/detectviz/detectviz/internal/adapters/scheduler"
	"github.com/detectviz/detectviz/internal/test/testutil"
	ifacescheduler "github.com/detectviz/detectviz/pkg/ifaces/scheduler"
	"github.com/stretchr/testify/assert"
)

// mockJobWithRetry implements scheduler.Job and simulates a retryable task.
// zh: 模擬會失敗後成功的任務，驗證 retry 機制。
type mockJobWithRetry struct {
	name      string
	spec      string
	calls     int
	failUntil int
}

func (m *mockJobWithRetry) Run(ctx context.Context) error {
	m.calls++
	if m.calls <= m.failUntil {
		return assert.AnError
	}
	return nil
}

func (m *mockJobWithRetry) Name() string {
	return m.name
}

func (m *mockJobWithRetry) Spec() string {
	return m.spec
}

// TestWorkerPoolSchedulerIntegration 驗證 WorkerPoolScheduler 能夠依據排程與重試邏輯執行工作。
// zh: 測試 worker pool 型排程器是否能正確執行失敗後可重試的任務。
func TestWorkerPoolSchedulerIntegration(t *testing.T) {
	ctx := context.Background()
	job := &mockJobWithRetry{
		name:      "retry-job",
		spec:      "@every 5s",
		failUntil: 1,
	}

	log := testutil.NewTestLogger()
	s := scheduleradapter.NewWorkerPoolScheduler(1, log)

	var _ ifacescheduler.Scheduler = s // 型別符合驗證

	s.Register(job)

	err := s.Start(ctx)
	assert.NoError(t, err)

	time.Sleep(7 * time.Second) // 等待 worker 排程與重試

	err = s.Stop(ctx)
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, job.calls, 2)
}



================================================
FILE: internal/adapters/server/server_adapter.go
================================================
package serveradapter

import (
	"context"

	core "github.com/detectviz/detectviz/internal/server"
)

// ServerAdapter wraps core.Server to implement iface.Server.
// zh: ServerAdapter 包裝 core.Server，使其實作 Server interface。
type ServerAdapter struct {
	srv *core.Server
}

// NewServerAdapter creates a new adapter instance.
// zh: 建立新的 ServerAdapter 實例。
func NewServerAdapter(s *core.Server) *ServerAdapter {
	return &ServerAdapter{
		srv: s,
	}
}

// Run starts the server.
// zh: 啟動伺服器。
func (a *ServerAdapter) Run(ctx context.Context) error {
	return a.srv.Run(ctx)
}

// Shutdown gracefully shuts down the server.
// zh: 優雅地關閉伺服器。
func (a *ServerAdapter) Shutdown(ctx context.Context) error {
	return a.srv.Shutdown(ctx)
}



================================================
FILE: internal/alert/alert.go
================================================
package alert

import (
	ifconfig "github.com/detectviz/detectviz/pkg/ifaces/config"
)

var cfg ifconfig.ConfigProvider

func Init(c ifconfig.ConfigProvider) {
	cfg = c
}

// 範例：使用設定值啟用或停用告警模組
func IsAlertEnabled() bool {
	return cfg.GetBool("alert.enabled")
}



================================================
FILE: internal/bootstrap/config_loader.go
================================================
package bootstrap

import (
	"github.com/detectviz/detectviz/pkg/config"
	// ConfigProvider 介面
	// zh: ConfigProvider 介面
	configiface "github.com/detectviz/detectviz/pkg/ifaces/config"
)

// LoadConfig initializes and returns a ConfigProvider instance.
// zh: 初始化並回傳 ConfigProvider 設定供應器。
func LoadConfig() configiface.ConfigProvider {
	return config.NewDefaultProvider()
}



================================================
FILE: internal/bootstrap/init.go
================================================
package bootstrap

import (
	fluxadapter "github.com/detectviz/detectviz/internal/adapters/alert/flux"
	promadapter "github.com/detectviz/detectviz/internal/adapters/alert/prom"
	"github.com/detectviz/detectviz/internal/alert"
	"github.com/detectviz/detectviz/internal/registry"

	alertregistry "github.com/detectviz/detectviz/internal/registry/alert"

	"github.com/detectviz/detectviz/pkg/config"
	ifconfig "github.com/detectviz/detectviz/pkg/ifaces/config"
)

var (
	Config ifconfig.ConfigProvider
	// Registry is the global runtime component registry.
	// zh: Registry 為平台執行期模組註冊中心。
	Registry            *registry.RegistryContainer
	AlertEvaluatorStore *alertregistry.AlertEvaluatorRegistry
)

// Init initializes the entire system.
// zh: 執行整體系統初始化，包含設定載入與核心元件註冊。
func Init() {
	initConfig()
	initRegistry()
	initAlertEvaluators()
	initAlertModule()
}

// initConfig loads the default config provider.
// zh: 載入預設設定提供者。
func initConfig() {
	Config = config.NewDefaultProvider()
}

// initRegistry sets up the runtime registry.
// zh: 註冊所有平台元件（如 notifier、scheduler 等）。
func initRegistry() {
	Registry = registry.NewRegistry(Config, "cron", Config.Logger())
}

// initAlertEvaluators registers alert evaluators (prometheus, flux).
// zh: 註冊預設告警評估器。
func initAlertEvaluators() {
	AlertEvaluatorStore = alertregistry.NewDefaultAlertEvaluatorRegistry(Config.Logger())
	AlertEvaluatorStore.Register("prometheus", promadapter.NewEvaluator(Config.Logger()))
	AlertEvaluatorStore.Register("flux", fluxadapter.NewEvaluator(Config.Logger()))
}

// initAlertModule injects config into the alert module.
// zh: 將設定注入 alert 模組。
func initAlertModule() {
	alert.Init(Config)
}



================================================
FILE: internal/bootstrap/wire.go
================================================
package bootstrap

import (
	modulesadapter "github.com/detectviz/detectviz/internal/adapters/modules"
	serveradapter "github.com/detectviz/detectviz/internal/adapters/server"
	"github.com/detectviz/detectviz/internal/modules"
	"github.com/detectviz/detectviz/internal/server"
	ifaceserver "github.com/detectviz/detectviz/pkg/ifaces/server"
)

// BuildServer assembles all core components and returns a Server interface.
// zh: 組裝所有系統元件並回傳 Server interface 實例。
func BuildServer() ifaceserver.Server {
	cfg := LoadConfig()
	engine := modules.NewEngine()

	engineAdapter := modulesadapter.NewEngineAdapter(engine)
	srv := server.NewServer(cfg, cfg.Logger(), engineAdapter)

	return serveradapter.NewServerAdapter(srv)
}



================================================
FILE: internal/modules/dependencies.go
================================================
package modules

import (
	"fmt"
)

// DependencyGraph defines module dependencies.
// zh: DependencyGraph 用於描述模組之間的依賴關係。
type DependencyGraph struct {
	nodes map[string][]string
}

// NewDependencyGraph creates a new empty dependency graph.
// zh: 建立新的模組依賴圖。
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		nodes: make(map[string][]string),
	}
}

// AddDependency adds a dependency from module 'from' to module 'to'.
// zh: 新增模組之間的依賴關係，表示 from 依賴 to。
func (g *DependencyGraph) AddDependency(from, to string) {
	g.nodes[from] = append(g.nodes[from], to)
}

// TopologicalSort returns a valid start order of modules based on dependencies.
// It returns an error if there is a cycle.
// zh: 根據依賴關係進行拓撲排序，若有循環依賴則回傳錯誤。
func (g *DependencyGraph) TopologicalSort() ([]string, error) {
	visited := make(map[string]bool)
	temp := make(map[string]bool)
	result := []string{}

	var visit func(string) error
	visit = func(n string) error {
		if temp[n] {
			return fmt.Errorf("cycle detected at module %q", n)
		}
		if !visited[n] {
			temp[n] = true
			for _, dep := range g.nodes[n] {
				if err := visit(dep); err != nil {
					return err
				}
			}
			temp[n] = false
			visited[n] = true
			result = append(result, n)
		}
		return nil
	}

	for n := range g.nodes {
		if !visited[n] {
			if err := visit(n); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}



================================================
FILE: internal/modules/engine.go
================================================
package modules

import (
	"context"
	"fmt"
	"sync"
)

// Module defines the standard lifecycle interface for all modules.
// zh: 定義所有模組的統一生命週期介面。
type Module interface {
	// Run starts the module and blocks until it is stopped or encounters an error.
	// zh: 啟動模組，阻塞直到停止或發生錯誤。
	Run(ctx context.Context) error

	// Shutdown gracefully stops the module.
	// zh: 優雅地關閉模組。
	Shutdown(ctx context.Context) error
}

// Engine is the central controller responsible for managing module lifecycles.
// zh: 控制所有模組註冊與生命週期的核心引擎。
type Engine struct {
	modules []Module
	mu      sync.Mutex
}

// NewEngine creates a new Engine instance to manage modules.
// zh: 建立新的模組控制引擎。
func NewEngine() *Engine {
	return &Engine{
		modules: make([]Module, 0),
	}
}

// Register adds a module to the Engine for lifecycle management.
// zh: 註冊一個模組到引擎中，由引擎管理其生命週期。
func (e *Engine) Register(m Module) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.modules = append(e.modules, m)
}

// RunAll starts all registered modules sequentially.
// It returns an error immediately if any module fails to start.
// zh: 順序啟動所有已註冊模組，若任一模組啟動失敗則立刻回傳錯誤。
func (e *Engine) RunAll(ctx context.Context) error {
	for _, m := range e.modules {
		if err := m.Run(ctx); err != nil {
			return err
		}
	}
	return nil
}

// ShutdownAll gracefully stops all registered modules.
// If multiple modules fail to shutdown, their errors are combined.
// zh: 優雅地關閉所有模組，若有多個錯誤會合併回傳。
func (e *Engine) ShutdownAll(ctx context.Context) error {
	var shutdownErr error

	for _, m := range e.modules {
		if err := m.Shutdown(ctx); err != nil {
			if shutdownErr == nil {
				shutdownErr = err
			} else {
				shutdownErr = fmt.Errorf("%w; %v", shutdownErr, err)
			}
		}
	}

	return shutdownErr
}



================================================
FILE: internal/modules/listener.go
================================================
package modules

import (
	"context"
	"log"
	"time"
)

// HealthCheckable represents modules that can report health status.
// zh: 具備健康狀態回報能力的模組介面。
type HealthCheckable interface {
	Healthy() bool
}

// Listener monitors the health of modules and triggers shutdown if needed.
// zh: Listener 監控所有模組健康狀態，必要時觸發全域停機。
type Listener struct {
	engine   *Engine
	registry *Registry
	interval time.Duration
	cancel   context.CancelFunc
}

// NewListener creates a new health monitoring listener.
// zh: 建立健康狀態監控器。
func NewListener(engine *Engine, registry *Registry, interval time.Duration) *Listener {
	return &Listener{
		engine:   engine,
		registry: registry,
		interval: interval,
	}
}

// Start launches the periodic health check loop.
// zh: 啟動定期健康檢查迴圈。
func (l *Listener) Start(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	l.cancel = cancel

	go func() {
		ticker := time.NewTicker(l.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if !l.checkAll() {
					log.Println("[listener] unhealthy module detected, shutting down all modules")
					_ = l.engine.ShutdownAll(ctx)
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Stop stops the health listener.
// zh: 停止健康監聽器。
func (l *Listener) Stop() {
	if l.cancel != nil {
		l.cancel()
	}
}

// checkAll verifies the health status of all registered modules that support HealthCheckable.
// zh: 檢查所有具備健康檢查功能的模組是否健康。
func (l *Listener) checkAll() bool {
	for _, name := range l.registry.List() {
		m, ok := l.registry.Get(name)
		if !ok {
			continue
		}
		hc, ok := m.(HealthCheckable)
		if ok && !hc.Healthy() {
			log.Printf("[listener] module %q is unhealthy", name)
			return false
		}
	}
	return true
}



================================================
FILE: internal/modules/registry.go
================================================
// Package modules provides global module registration and management.
// zh: modules 套件提供全域模組註冊與管理功能。
package modules

import (
	"fmt"
	"sync"
)

// Registry maintains a list of named modules for global management.
// zh: Registry 用於管理所有具名模組的註冊與查找。
type Registry struct {
	mu      sync.RWMutex
	modules map[string]Module
}

// NewRegistry creates a new module registry.
// zh: 建立一個模組註冊器。
func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

// Register registers a named module to the registry.
// Returns an error if the name already exists.
// zh: 註冊具名模組，若名稱已存在則回傳錯誤。
func (r *Registry) Register(name string, m Module) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.modules[name]; exists {
		return fmt.Errorf("module %q already registered", name)
	}

	r.modules[name] = m
	return nil
}

// Get retrieves a registered module by name.
// zh: 透過名稱查找已註冊模組。
func (r *Registry) Get(name string) (Module, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, ok := r.modules[name]
	return m, ok
}

// List returns all registered module names.
// zh: 回傳所有已註冊的模組名稱。
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	return names
}



================================================
FILE: internal/modules/runner.go
================================================
package modules

import (
	"context"
	"fmt"
)

// Runner coordinates the startup and shutdown of all modules based on dependencies.
// zh: Runner 根據模組依賴關係協調所有模組的啟動與停止。
type Runner struct {
	engine   *Engine
	registry *Registry
	graph    *DependencyGraph
	started  []string
}

// NewRunner creates a new module runner.
// zh: 建立模組啟動與關閉的協調器。
func NewRunner(engine *Engine, registry *Registry, graph *DependencyGraph) *Runner {
	return &Runner{
		engine:   engine,
		registry: registry,
		graph:    graph,
	}
}

// StartAll starts all registered modules in topological order.
// zh: 依拓撲排序啟動所有模組。
func (r *Runner) StartAll(ctx context.Context) error {
	order, err := r.graph.TopologicalSort()
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	for _, name := range order {
		module, ok := r.registry.Get(name)
		if !ok {
			return fmt.Errorf("module %q not found in registry", name)
		}
		if err := module.Run(ctx); err != nil {
			return fmt.Errorf("failed to start module %q: %w", name, err)
		}
		r.started = append(r.started, name)
	}

	return nil
}

// StopAll shuts down all started modules in reverse order.
// zh: 依反向順序關閉所有模組。
func (r *Runner) StopAll(ctx context.Context) error {
	for i := len(r.started) - 1; i >= 0; i-- {
		name := r.started[i]
		module, ok := r.registry.Get(name)
		if ok {
			if err := module.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed to stop module %q: %w", name, err)
			}
		}
	}
	return nil
}



================================================
FILE: internal/plugins/eventbus/alertlog/alert_handler.go
================================================
package alertlog

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// AlertPluginHandler is a sample plugin that handles alert events.
// zh: AlertPluginHandler 是一個範例插件，實作告警事件處理邏輯。
type AlertPluginHandler struct{}

// HandleAlertTriggered logs or processes the alert event.
// zh: 在此處理 AlertTriggeredEvent，例如自訂通知或紀錄。
func (h *AlertPluginHandler) HandleAlertTriggered(ctx context.Context, e event.AlertTriggeredEvent) error {
	log := logger.FromContext(ctx)
	log.WithFields(map[string]any{
		"alert_id":   e.AlertID,
		"rule_name":  e.RuleName,
		"level":      e.Level,
		"instance":   e.Instance,
		"metric":     e.Metric,
		"comparison": e.Comparison,
		"value":      e.Value,
		"threshold":  e.Threshold,
	}).Info("[ALERT] triggered")
	return nil
}



================================================
FILE: internal/plugins/eventbus/alertlog/alert_handler_test.go
================================================
package alertlog

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

func TestAlertPluginHandler_HandleAlertTriggered(t *testing.T) {
	handler := &AlertPluginHandler{}

	t.Run("handle normal alert", func(t *testing.T) {
		event := event.AlertTriggeredEvent{
			AlertID:    "test-alert",
			RuleName:   "test-rule",
			Level:      "critical",
			Instance:   "test-instance",
			Metric:     "test-metric",
			Comparison: ">",
			Value:      100,
			Threshold:  100,
			Message:    "alert test",
		}

		err := handler.HandleAlertTriggered(context.Background(), event)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("handle empty fields", func(t *testing.T) {
		event := event.AlertTriggeredEvent{
			AlertID:    "",
			RuleName:   "",
			Level:      "",
			Instance:   "",
			Metric:     "",
			Comparison: "",
			Value:      0,
			Threshold:  0,
			Message:    "",
		}

		err := handler.HandleAlertTriggered(context.Background(), event)
		if err != nil {
			t.Errorf("expected no error for empty fields, got %v", err)
		}
	})
}



================================================
FILE: internal/plugins/eventbus/alertlog/init.go
================================================
package alertlog

import (
	eventbusadapter "github.com/detectviz/detectviz/internal/adapters/eventbus"
	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

// 確保 AlertPluginHandler 實作 AlertEventHandler
var _ event.AlertEventHandler = (*AlertPluginHandler)(nil)

// init 初始化 alertlog plugin 並註冊其事件處理器
// zh: 在載入此模組時自動註冊 AlertPluginHandler 作為告警事件處理器。
func init() {
	eventbusadapter.RegisterAlertHandler(&AlertPluginHandler{})
}



================================================
FILE: internal/registry/registry.go
================================================
package registry

import (
	eventbusadapter "github.com/detectviz/detectviz/internal/adapters/eventbus"
	cachestoreregistry "github.com/detectviz/detectviz/internal/registry/cachestore"
	notifierregistry "github.com/detectviz/detectviz/internal/registry/notifier"
	scheduleregistry "github.com/detectviz/detectviz/internal/registry/scheduler"
	"github.com/detectviz/detectviz/pkg/ifaces/cachestore"
	ifconfig "github.com/detectviz/detectviz/pkg/ifaces/config"
	eventbusiface "github.com/detectviz/detectviz/pkg/ifaces/eventbus"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"github.com/detectviz/detectviz/pkg/ifaces/notifier"
	"github.com/detectviz/detectviz/pkg/ifaces/scheduler"
)

// RegistryContainer 包含 Detectviz 所有模組註冊項目
// zh: 用於封裝所有可注入與呼叫的服務或核心元件
type RegistryContainer struct {
	EventDispatcher eventbusiface.EventDispatcher
	Scheduler       scheduler.Scheduler
	Notifier        notifier.Notifier
	CacheStore      cachestore.CacheStore
}

// NewRegistry 建立並初始化所有必要的註冊模組
// zh: 使用指定的設定與 logger 建立平台所需資源
func NewRegistry(cfg ifconfig.ConfigProvider, schedulerProvider string, log logger.Logger) *RegistryContainer {
	return &RegistryContainer{
		EventDispatcher: NewInMemoryEventDispatcher(),
		Scheduler:       scheduleregistry.ProvideScheduler(schedulerProvider, log),
		Notifier:        notifierregistry.ProvideNotifier(cfg.GetNotifierConfigs(), log),
		CacheStore:      cachestoreregistry.RegisterCacheStore(cfg.GetCacheConfig()),
	}
}

// NewInMemoryEventDispatcher 建立新的事件總線 Dispatcher（使用記憶體）
// zh: 建立並整合內建與 plugin 註冊的所有事件處理器。
func NewInMemoryEventDispatcher() eventbusiface.EventDispatcher {
	dispatcher := eventbusadapter.NewInMemoryDispatcher()
	registerAllHandlers(dispatcher)
	return dispatcher
}

// NewKafkaEventDispatcher 建立新的事件總線 Dispatcher（使用 Kafka 作為傳輸層）
// zh: 建立 Kafka 型態的事件總線 Dispatcher，預留未來整合 Kafka handler。
func NewKafkaEventDispatcher() eventbusiface.EventDispatcher {
	// TODO: 實作 Kafka-based Dispatcher，例如 eventbus.NewKafkaDispatcher()
	// dispatcher := eventbus.NewKafkaDispatcher()
	// registerAllHandlers(dispatcher)
	panic("KafkaEventDispatcher 尚未實作")
}

// registerAllHandlers 將所有 plugin handler 註冊到指定 Dispatcher
// zh: 將 Alert、Host、Metric、Task 等處理器註冊至事件總線。
func registerAllHandlers(d eventbusiface.EventDispatcher) {
	for _, h := range eventbusadapter.LoadPluginAlertHandlers() {
		d.RegisterAlertHandler(h)
	}
	for _, h := range eventbusadapter.LoadPluginHostHandlers() {
		d.RegisterHostHandler(h)
	}
	for _, h := range eventbusadapter.LoadPluginMetricHandlers() {
		d.RegisterMetricHandler(h)
	}
	for _, h := range eventbusadapter.LoadPluginTaskHandlers() {
		d.RegisterTaskHandler(h)
	}
}



================================================
FILE: internal/registry/alert/registry.go
================================================
package alert

import (
	"fmt"
	"sync"

	fluxadapter "github.com/detectviz/detectviz/internal/adapters/alert/flux"
	promadapter "github.com/detectviz/detectviz/internal/adapters/alert/prom"
	"github.com/detectviz/detectviz/pkg/ifaces/alert"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// AlertEvaluatorRegistry 管理所有註冊的 AlertEvaluator 實例
// zh: 負責註冊與查詢可用的告警評估器實作。
type AlertEvaluatorRegistry struct {
	evaluators map[string]alert.AlertEvaluator // zh: 已註冊的告警評估器映射表
	mu         sync.RWMutex                    // zh: 保護 evaluators 映射表的同步鎖
	log        logger.Logger                   // zh: 用於記錄註冊過程的日誌
}

// NewAlertEvaluatorRegistryWithLogger 建立註冊中心實例並注入 logger
// zh: 初始化並回傳告警評估器註冊中心，並設定日誌記錄器。
func NewAlertEvaluatorRegistryWithLogger(log logger.Logger) *AlertEvaluatorRegistry {
	return &AlertEvaluatorRegistry{
		evaluators: make(map[string]alert.AlertEvaluator),
		log:        log.Named("alert-registry"),
	}
}

// NewDefaultAlertEvaluatorRegistry 建立包含預設實作的註冊中心
// zh: 注入 PromEvaluator 與 FluxEvaluator 供預設使用。
func NewDefaultAlertEvaluatorRegistry(log logger.Logger) *AlertEvaluatorRegistry {
	r := NewAlertEvaluatorRegistryWithLogger(log)
	r.Register("prometheus", promadapter.NewEvaluator(log))
	r.Register("flux", fluxadapter.NewEvaluator(log))
	return r
}

// Register 註冊一個新的 AlertEvaluator
// zh: 以指定名稱註冊一個告警評估器實作，若名稱已存在則覆蓋並記錄警告。
func (r *AlertEvaluatorRegistry) Register(name string, evaluator alert.AlertEvaluator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.evaluators[name]; exists {
		r.log.Warn("overriding existing alert evaluator", "name", name)
	} else {
		r.log.Info("registering new alert evaluator", "name", name)
	}
	r.evaluators[name] = evaluator
}

// Get 根據名稱取得已註冊的 AlertEvaluator
// zh: 回傳對應名稱的告警評估器，若不存在則回傳錯誤。
func (r *AlertEvaluatorRegistry) Get(name string) (alert.AlertEvaluator, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	evaluator, ok := r.evaluators[name]
	if !ok {
		return nil, fmt.Errorf("alert evaluator '%s' not found", name)
	}
	return evaluator, nil
}

// GetDefault 回傳預設的 AlertEvaluator（prometheus）。
// zh: 簡化用法，直接取得預設的告警評估器（目前為 prometheus 實作）。
func (r *AlertEvaluatorRegistry) GetDefault() alert.AlertEvaluator {
	e, _ := r.Get("prometheus")
	return e
}



================================================
FILE: internal/registry/cachestore/registry.go
================================================
package cachestore

import (
	memoryadapter "github.com/detectviz/detectviz/internal/adapters/cachestore/memory"
	redisadapter "github.com/detectviz/detectviz/internal/adapters/cachestore/redis"
	"github.com/detectviz/detectviz/pkg/configtypes"
	cachestoreiface "github.com/detectviz/detectviz/pkg/ifaces/cachestore"
	goredis "github.com/redis/go-redis/v9"
)

// RegisterCacheStore 根據設定註冊 CacheStore adapter。
// zh: 若 config 指定 redis，則註冊 redis adapter；否則預設使用記憶體快取。
func RegisterCacheStore(cfg configtypes.CacheConfig) cachestoreiface.CacheStore {
	if cfg.Backend == "redis" {
		client := goredis.NewClient(&goredis.Options{
			Addr:     cfg.Redis.Address,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
		return redisadapter.NewRedisCacheStore(client)
	}

	return memoryadapter.NewMemoryCacheStore()
}



================================================
FILE: internal/registry/config/registry.go
================================================
package config

import (
	"github.com/detectviz/detectviz/pkg/config"
	configiface "github.com/detectviz/detectviz/pkg/ifaces/config"
)

// RegisterConfigProvider 初始化並回傳預設的 ConfigProvider。
// zh: 用於統一初始化設定模組，供其他模組依賴注入。
func RegisterConfigProvider() configiface.ConfigProvider {
	return config.NewDefaultProvider()
}



================================================
FILE: internal/registry/eventbus/plugins.go
================================================
package eventbus

import (
	// 強制載入 plugin 模組以觸發其 init() 註冊事件處理器
	adapterlogger "github.com/detectviz/detectviz/internal/adapters/logger"
	"github.com/detectviz/detectviz/pkg/ifaces/event"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

var (
	registeredAlertHandlers  []event.AlertEventHandler
	registeredHostHandlers   []event.HostEventHandler
	registeredMetricHandlers []event.MetricEventHandler
	registeredTaskHandlers   []event.TaskEventHandler
	defaultLogger            logger.Logger = adapterlogger.NewNopLogger()
)

// RegisterPluginAlertHandler allows plugin modules to register alert handlers.
// zh: 外部模組註冊自訂 Alert 處理器。
func RegisterPluginAlertHandler(h event.AlertEventHandler) {
	registeredAlertHandlers = append(registeredAlertHandlers, h)
}

// RegisterPluginHostHandler allows plugin modules to register host handlers.
// zh: 外部模組註冊自訂 Host 處理器。
func RegisterPluginHostHandler(h event.HostEventHandler) {
	registeredHostHandlers = append(registeredHostHandlers, h)
}

// RegisterPluginMetricHandler allows plugin modules to register metric handlers.
// zh: 外部模組註冊自訂 Metric 處理器。
func RegisterPluginMetricHandler(h event.MetricEventHandler) {
	registeredMetricHandlers = append(registeredMetricHandlers, h)
}

// RegisterPluginTaskHandler allows plugin modules to register task handlers.
// zh: 外部模組註冊自訂 Task 處理器。
func RegisterPluginTaskHandler(h event.TaskEventHandler) {
	registeredTaskHandlers = append(registeredTaskHandlers, h)
}

// LoadPluginAlertHandlers retrieves all plugin-registered alert handlers.
// zh: 載入所有插件註冊的 Alert 處理器。
func LoadPluginAlertHandlers() []event.AlertEventHandler {
	return registeredAlertHandlers
}

// LoadPluginHostHandlers retrieves all plugin-registered host handlers.
// zh: 載入所有插件註冊的 Host 處理器。
func LoadPluginHostHandlers() []event.HostEventHandler {
	return registeredHostHandlers
}

// LoadPluginMetricHandlers retrieves all plugin-registered metric handlers.
// zh: 載入所有插件註冊的 Metric 處理器。
func LoadPluginMetricHandlers() []event.MetricEventHandler {
	return registeredMetricHandlers
}

// LoadPluginTaskHandlers retrieves all plugin-registered task handlers.
// zh: 載入所有插件註冊的 Task 處理器。
func LoadPluginTaskHandlers() []event.TaskEventHandler {
	return registeredTaskHandlers
}

// ExplorePlugins 探索並初始化所有 event plugin 模組。
// zh: 匯入所有 plugin，以觸發各自 init 函式註冊對應事件處理器。
func ExplorePlugins() {
	// 模組透過上方 import _ 初始化，故此函式本身不需邏輯
}

// OverrideDefaultLogger 設定 plugin 使用的預設 logger。
// zh: 測試時可使用此方法注入 log 攔截器。
func OverrideDefaultLogger(log logger.Logger) {
	defaultLogger = log
}

// GetDefaultLogger 回傳 plugin 使用的 logger。
// zh: 提供 plugin 在 init 時注入的預設 logger 實例。
func GetDefaultLogger() logger.Logger {
	return defaultLogger
}



================================================
FILE: internal/registry/eventbus/providers.go
================================================
package eventbus

import (
	"fmt"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/eventbus"
)

var (
	providersMu sync.RWMutex
	providers   = make(map[string]eventbus.DispatcherProvider)
)

// RegisterProvider 註冊 provider，例如 in-memory / kafka / nats
func RegisterProvider(p eventbus.DispatcherProvider) {
	providersMu.Lock()
	defer providersMu.Unlock()
	providers[p.Name()] = p
}

// GetProvider 依據 provider 名稱取得 EventDispatcher
func GetProvider(name string) (eventbus.EventDispatcher, error) {
	providersMu.RLock()
	defer providersMu.RUnlock()
	p, ok := providers[name]
	if !ok {
		return nil, fmt.Errorf("eventbus provider '%s' not found", name)
	}
	return p.Build(), nil
}



================================================
FILE: internal/registry/eventbus/registry.go
================================================
package eventbus

import (
	"fmt"

	eventbusadapter "github.com/detectviz/detectviz/internal/adapters/eventbus"
	"github.com/detectviz/detectviz/pkg/ifaces/event"
	eventbusiface "github.com/detectviz/detectviz/pkg/ifaces/eventbus"
)

// DispatcherOptions allows injection of custom handlers for each event type.
// zh: DispatcherOptions 提供事件處理器的組態選項，支援外部擴充與測試注入。
type DispatcherOptions struct {
	AlertHandlers  []event.AlertEventHandler
	HostHandlers   []event.HostEventHandler
	MetricHandlers []event.MetricEventHandler
	TaskHandlers   []event.TaskEventHandler
}

// NewInMemoryEventDispatcherWithOptions constructs an in-memory dispatcher with injected handlers.
// zh: 根據傳入的 DispatcherOptions 建立記憶體型事件分派器。
func NewInMemoryEventDispatcherWithOptions(opt DispatcherOptions) eventbusiface.EventDispatcher {
	dispatcher := eventbusadapter.NewInMemoryDispatcher()

	for _, h := range opt.AlertHandlers {
		dispatcher.RegisterAlertHandler(h)
	}
	for _, h := range opt.HostHandlers {
		dispatcher.RegisterHostHandler(h)
	}
	for _, h := range opt.MetricHandlers {
		dispatcher.RegisterMetricHandler(h)
	}
	for _, h := range opt.TaskHandlers {
		dispatcher.RegisterTaskHandler(h)
	}

	return dispatcher
}

// NewEventDispatcher 建立指定 provider 的事件總線實作。
// zh: 根據指定 provider 名稱建立對應的事件處理器，可擴充支援 kafka/nats 等實作。
func NewEventDispatcher(provider string) (eventbusiface.EventDispatcher, error) {
	p, err := GetProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("failed to get eventbus provider: %w", err)
	}
	return p, nil
}



================================================
FILE: internal/registry/eventbus/registry_inmemory.go
================================================
package eventbus

import (
	eventbusadapter "github.com/detectviz/detectviz/internal/adapters/eventbus"
	eventbusiface "github.com/detectviz/detectviz/pkg/ifaces/eventbus"
)

type inMemoryProvider struct{}

func (p *inMemoryProvider) Name() string {
	return "in-memory"
}

func (p *inMemoryProvider) Build() eventbusiface.EventDispatcher {
	dispatcher := eventbusadapter.NewInMemoryDispatcher()

	for _, h := range eventbusadapter.LoadPluginAlertHandlers() {
		dispatcher.RegisterAlertHandler(h)
	}
	for _, h := range eventbusadapter.LoadPluginHostHandlers() {
		dispatcher.RegisterHostHandler(h)
	}
	for _, h := range eventbusadapter.LoadPluginMetricHandlers() {
		dispatcher.RegisterMetricHandler(h)
	}
	for _, h := range eventbusadapter.LoadPluginTaskHandlers() {
		dispatcher.RegisterTaskHandler(h)
	}

	return dispatcher
}

func init() {
	RegisterProvider(&inMemoryProvider{})
}



================================================
FILE: internal/registry/logger/registry.go
================================================
package logger

import (
	loggeradapter "github.com/detectviz/detectviz/internal/adapters/logger"
	ifacelogger "github.com/detectviz/detectviz/pkg/ifaces/logger"
	"go.uber.org/zap"
)

// ProvideLogger creates a logger instance using the Zap adapter implementation.
// zh: 建立 logger 實例，採用 internal/adapters/logger/zap_adapter.go 作為預設實作。
func ProvideLogger() ifacelogger.Logger {
	// Use Zap logger by default.
	// zh: 預設使用 Zap logger，實作於 internal/adapters/logger/zap_adapter.go。
	zapLogger, _ := zap.NewDevelopment()
	return loggeradapter.NewZapLogger(zapLogger.Sugar())
}



================================================
FILE: internal/registry/notifier/registry.go
================================================
package notifier

import (
	"errors"
	"fmt"
	"net/http"

	notifieradapter "github.com/detectviz/detectviz/internal/adapters/notifier"
	"github.com/detectviz/detectviz/pkg/configtypes"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	notifieriface "github.com/detectviz/detectviz/pkg/ifaces/notifier"
)

// NewNotifierRegistry 建立 notifier 實例清單，根據 config 載入。
// zh: 根據設定檔產生可用的 notifier 實作清單。
func NewNotifierRegistry(cfgs []configtypes.NotifierConfig, log logger.Logger) []notifieriface.Notifier {
	var list []notifieriface.Notifier

	for _, cfg := range cfgs {
		if !cfg.Enable {
			continue
		}
		n, err := buildNotifier(cfg, log)
		if err != nil {
			log.Warn("Failed to build notifier", "type", cfg.Type, "name", cfg.Name, "err", err)
			continue
		}
		list = append(list, n)
	}

	return list
}

// buildNotifier 根據設定建立對應的 Notifier 實例。
// zh: 根據 notifier 類型回傳實作，支援 email, slack, webhook。
func buildNotifier(cfg configtypes.NotifierConfig, log logger.Logger) (notifieriface.Notifier, error) {
	switch cfg.Type {
	case "email":
		return notifieradapter.NewEmailNotifier(cfg.Name, cfg.Target, log), nil
	case "slack":
		return notifieradapter.NewSlackNotifier(cfg.Name, cfg.Target, log), nil
	case "webhook":
		return notifieradapter.NewWebhookNotifier(cfg.Name, log, http.DefaultClient), nil
	default:
		return nil, errors.New(fmt.Sprintf("unknown notifier type: %s", cfg.Type))
	}
}

// ProvideNotifier 根據設定與 logger 提供整合後的 Notifier 實例。
// zh: 整合所有 notifier 為一個 Notifier，用於統一發送通知。
func ProvideNotifier(cfgs []configtypes.NotifierConfig, log logger.Logger) notifieriface.Notifier {
	list := NewNotifierRegistry(cfgs, log)
	if len(list) == 0 {
		return notifieradapter.NewNopNotifier()
	}
	return notifieradapter.NewMultiNotifier(list...)
}



================================================
FILE: internal/registry/notifier/registry_test.go
================================================
package notifier_test

import (
	"testing"

	"github.com/detectviz/detectviz/internal/registry/notifier"
	"github.com/detectviz/detectviz/internal/test/testutil"
	"github.com/detectviz/detectviz/pkg/configtypes"
)

// TestNewNotifierRegistry 測試 notifier 註冊邏輯。
// zh: 檢查是否依據設定正確建立 notifier 清單。
func TestNewNotifierRegistry(t *testing.T) {
	cfgs := []configtypes.NotifierConfig{
		{Name: "email", Type: "email", Target: "noreply@example.com", Enable: true},
		{Name: "slack", Type: "slack", Target: "https://hooks.slack.com/xxx", Enable: true},
		{Name: "webhook", Type: "webhook", Target: "https://example.com/webhook", Enable: false}, // zh: 停用
		{Name: "invalid", Type: "foo", Target: "xxx", Enable: true},                              // zh: 無效類型
	}

	log := testutil.NewTestLogger()
	notifiers := notifier.NewNotifierRegistry(cfgs, log)

	if len(notifiers) != 2 {
		t.Errorf("expected 2 notifiers, got %d", len(notifiers))
	}

	names := map[string]bool{}
	for _, n := range notifiers {
		names[n.Name()] = true
	}

	if !names["email"] {
		t.Error("email notifier not registered")
	}
	if !names["slack"] {
		t.Error("slack notifier not registered")
	}
	if names["webhook"] {
		t.Error("disabled notifier should not be registered")
	}
	if names["invalid"] {
		t.Error("invalid type should not be registered")
	}
}



================================================
FILE: internal/registry/scheduler/registry.go
================================================
package scheduler

import (
	adapters "github.com/detectviz/detectviz/internal/adapters/scheduler"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"github.com/detectviz/detectviz/pkg/ifaces/scheduler"
)

// ProvideScheduler returns a scheduler instance based on the given name.
// zh: 根據指定名稱回傳對應的排程器實作。
func ProvideScheduler(name string, log logger.Logger) scheduler.Scheduler {
	switch name {
	case "cron":
		return adapters.NewCronScheduler(log)
	case "workerpool":
		return adapters.NewWorkerPoolScheduler(4, log) // 預設 4 workers
	case "mock":
		return adapters.NewMockScheduler()
	default:
		log.Warn("unknown scheduler type, fallback to mock", "name", name)
		return adapters.NewMockScheduler()
	}
}



================================================
FILE: internal/server/instrumentation.go
================================================
package server

import (
	"net/http"
	"net/http/pprof"
)

// RegisterInstrumentation mounts instrumentation routes into the given ServeMux.
// zh: 將監控與除錯端點註冊至 HTTP multiplexer。
func RegisterInstrumentation(mux *http.ServeMux) {
	// Prometheus metrics endpoint (佔位用，可替換為 promhttp.Handler())
	mux.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("# HELP detectviz_metrics_placeholder Placeholder for Prometheus metrics\n"))
		w.Write([]byte("# TYPE detectviz_metrics_placeholder counter\n"))
		w.Write([]byte("detectviz_metrics_placeholder 1\n"))
	}))

	// Basic health check endpoint
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))

	// Register pprof endpoints
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
}



================================================
FILE: internal/server/runner.go
================================================
package server

import (
	"context"
	"fmt"
	"net/http"
)

// Run starts all modules and the HTTP server.
// zh: 啟動所有模組並啟動 HTTP 伺服器。
func (s *Server) Run(ctx context.Context) error {
	// 啟動模組
	if err := s.ModuleEngine.RunAll(ctx); err != nil {
		return fmt.Errorf("failed to run modules: %w", err)
	}

	// 建立 HTTP Server
	s.HTTPServer = &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux, // TODO: 可替換為自定義 multiplexer
	}

	// 啟動 HTTP Server
	go func() {
		if err := s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Error("HTTP server error", "error", err)
		}
	}()

	s.Logger.Info("Server started on :8080")
	<-ctx.Done()
	return nil
}

// Shutdown gracefully shuts down the HTTP server and all modules.
// zh: 優雅地關閉 HTTP 伺服器與所有模組。
func (s *Server) Shutdown(ctx context.Context) error {
	s.Logger.Info("Shutting down server...")

	if s.HTTPServer != nil {
		if err := s.HTTPServer.Shutdown(ctx); err != nil {
			s.Logger.Warn("HTTP server shutdown error", "error", err)
		}
	}

	if err := s.ModuleEngine.ShutdownAll(ctx); err != nil {
		return fmt.Errorf("failed to shutdown modules: %w", err)
	}

	s.Logger.Info("Server shutdown complete")
	return nil
}



================================================
FILE: internal/server/server.go
================================================
package server

import (
	"net/http"

	"github.com/detectviz/detectviz/pkg/ifaces/config"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"github.com/detectviz/detectviz/pkg/ifaces/modules"
)

// Server represents the core application server.
// zh: Server 為核心伺服器，整合設定、日誌與模組控制。
type Server struct {
	Config       config.ConfigProvider
	Logger       logger.Logger
	ModuleEngine modules.ModuleEngine
	HTTPServer   *http.Server
}

// NewServer creates a new instance of Server.
// zh: 建立新的伺服器實例，注入設定與模組控制元件。
func NewServer(
	cfg config.ConfigProvider,
	log logger.Logger,
	engine modules.ModuleEngine,
) *Server {
	return &Server{
		Config:       cfg,
		Logger:       log,
		ModuleEngine: engine,
		HTTPServer:   nil, // 後續由 runner/init 時注入
	}
}



================================================
FILE: internal/test/README.md
================================================
# 測試設計與實作規範

本文件說明 Detectviz 平台的測試設計策略、命名原則、Mock/Fake 管理方式，對齊 Clean Architecture + Plugin 設計，確保可維護性與擴充性。

## 各目錄詳細說明

本專案遵循 Clean Architecture 與 Go 社群常見設計慣例，建議測試檔案與 mock/fake 實作依據責任分工放置在以下位置：

| 類型 | 責任 | 目錄位置 | 適用範例 |
| --- | --- | --- | --- |
| **單元測試** | 測試模組本身的邏輯與行為（純 function or adapter） | **與被測模組同層**（`foo.go` → `foo_test.go`） | logger、notifier、scheduler adapter |
| **整合測試** | 驗證多模組組合（註冊、依賴注入、流程運作） | `/internal/test/` | scheduler + notifier + logger 測通流程 |
| **Fake 實作** | 提供固定邏輯模擬 interface，用於非驗證式測試 | `/internal/test/fakes/` 或 `/test/fakes/` | `fake_scheduler.go`、`fake_notifier.go` |
| **Mock 實作** | 驗證呼叫方法/次數/參數是否正確 | `/internal/test/mocks/` 或 `/pkg/mocks/` | `mock_scheduler.go`（mockery 產出） |
| **共用測試工具** | 建立可複用 logger、clock、context、資料構造等 | `/internal/test/testutil/` | `test_logger.go`、`assert_logger.go`、`fake_clock.go` |

### 統整建議

- 模組行為測試 → 模組內
- 整合流程驗證 → internal/test/integration/
- 可共用邏輯模擬 → internal/test/fakes/
- 呼叫驗證測試用 → internal/test/mocks/
- log/config 等測試工具 → internal/test/testutil/

## 測試撰寫原則

- 每個模組皆應具備 `_test.go` 測試檔。
- 單元測試應以 `t.Run(...)` 切分子情境。
- Fake 用於模擬邏輯流程，Mock 用於驗證呼叫與參數。
- 測試涵蓋正常與異常情境（happy path / error path）。
- 禁止在非測試模組內定義 logger 或 config 的 mock 實作。


## 測試目錄結構

```bash
detectviz/
├── internal/
│   ├── adapters/
│   │   └── logger/
│   │       ├── zap_adapter.go
│   │       └── zap_adapter_test.go    # 單元測試：與 adapter 同層
│   ├── registry/
│   │   └── scheduler/
│   │       ├── registry.go
│   │       └── registry_test.go       # 單元測試（可與 adapter 區隔）
│   ├── test/
│   │   ├── fakes/
│   │   │   └── fake_notifier.go
│   │   ├── mocks/
│   │   │   └── mock_scheduler.go
│   │   ├── testutil/
│   │   │   ├── test_logger.go
│   │   │   ├── assert_logger.go
│   │   │   └── fake_clock.go
│   │   └── integration/
│   │       └── alert_pipeline_test.go
├── pkg/
│   ├── config/
│   │   ├── default.go
│   │   └── default_test.go
│   ├── ifaces/
│   └── mocks/       # mockery 自動產生可選集中放這
```

## 實際範例說明

### `internal/test/fakes/`

- 用途：手動實作的 Fake 類別，符合 interface，但不執行真實邏輯。
- 用於整合測試、模擬流程、不進行斷言。
- 範例：
  - `fake_notifier.go`：模擬發送通知但不實際執行，供 alert 流程測試。
  - `fake_scheduler.go`：模擬任務排程器，用於測試 alert 判斷後的任務注入。

### `internal/test/testutil/`

- 用途：測試輔助函式與共用模組（logger、時間模擬、預設 context 等）。
- 範例：
  - `test_logger.go`：空實作 Logger，適用不需驗證 log 的測試場景。
  - `assert_logger.go`：具 log 記錄功能，可斷言是否有記錄 log 與內容。
  - `fake_clock.go`：固定時間點模擬，便於測試時間觸發邏輯。
  - `test_context.go`：產生標準 context，內含 metadata 標記。

### `internal/test/integration`

- 用途：放置跨模組整合測試（如初始化流程、完整事件觸發測試）。
- 測試真實模組的註冊與協作行為。
- 範例：
  - `registry_init_test.go`：測試 logger/notifier/scheduler 等模組在註冊時是否成功注入。
  - `alert_pipeline_test.go`：模擬告警流程從 metric 觸發 → alert 判斷 → notifier 發送的完整流程。

### `pkg/mocks/`

- 用途：自動產生的 Mock 類別，使用 `mockery` 等工具產生。
- 用於驗證行為、呼叫次數、傳入參數。
- 對應指令範例：
  ```bash
  mockery --name=Notifier --output=./pkg/mocks/
  ```
- 實例用途：
  - `mock_notifier.go`：驗證是否正確呼叫 `.Send()` 並傳入預期 alert message。
  - `mock_scheduler.go`：驗證是否正確註冊任務並執行。



================================================
FILE: internal/test/fakes/fake_config.go
================================================
// Package fakes 提供測試用的假物件。
// zh: 用於測試的設定假實作，支援 key-value 模擬與錯誤注入。
package fakes

import (
	"errors"
	"strconv"

	"github.com/detectviz/detectviz/pkg/configtypes"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// FakeConfig implements the ConfigProvider interface for testing purposes.
// zh: 用於測試的設定假實作，支援 key-value 模擬與錯誤注入。
type FakeConfig struct {
	// Values 儲存 key-value 字串對應，模擬設定內容。
	Values map[string]string
	// CacheCfg 模擬快取設定。
	CacheCfg configtypes.CacheConfig
	// NotifierCfgs 模擬通知器設定清單。
	NotifierCfgs []configtypes.NotifierConfig
	// ReloadCalled 標記 Reload 是否被呼叫過。
	ReloadCalled bool
	// LoggerInstance 模擬 logger 實例。
	LoggerInstance logger.Logger
	// ReloadShouldFail 控制 Reload 是否回傳錯誤。
	ReloadShouldFail bool
}

// Get returns the string value for a given key.
// zh: 取得指定 key 的字串值，若不存在則回傳空字串。
func (f *FakeConfig) Get(key string) string {
	return f.Values[key]
}

// GetOrDefault returns the value or default if not set.
// zh: 取得指定 key 的字串值，若不存在則回傳預設值。
func (f *FakeConfig) GetOrDefault(key, defaultVal string) string {
	val, ok := f.Values[key]
	if !ok {
		return defaultVal
	}
	return val
}

// GetInt returns the int value for a given key.
// zh: 取得指定 key 的整數值，若格式錯誤則回傳 0。
func (f *FakeConfig) GetInt(key string) int {
	v, _ := strconv.Atoi(f.Values[key])
	return v
}

// GetBool returns the bool value for a given key.
// zh: 取得指定 key 的布林值，若格式錯誤則回傳 false。
func (f *FakeConfig) GetBool(key string) bool {
	v, _ := strconv.ParseBool(f.Values[key])
	return v
}

// GetCacheConfig returns the cache config struct.
// zh: 取得快取設定結構。
func (f *FakeConfig) GetCacheConfig() configtypes.CacheConfig {
	return f.CacheCfg
}

// GetNotifierConfigs returns notifier configuration list.
// zh: 取得通知器設定清單。
func (f *FakeConfig) GetNotifierConfigs() []configtypes.NotifierConfig {
	return f.NotifierCfgs
}

// Logger returns the logger instance.
// zh: 取得 logger 實例。
func (f *FakeConfig) Logger() logger.Logger {
	return f.LoggerInstance
}

// Reload sets ReloadCalled and optionally returns error.
// zh: 重新載入設定，會標記 ReloadCalled，並可模擬錯誤。
func (f *FakeConfig) Reload() error {
	f.ReloadCalled = true
	if f.ReloadShouldFail {
		return errors.New("simulated reload failure")
	}
	return nil
}



================================================
FILE: internal/test/fakes/fake_metrics.go
================================================
package fakes

import (
	"context"

	metric "github.com/detectviz/detectviz/pkg/ifaces/metrics"
)

// FakeMetricWriter is a fake implementation of metric.MetricWriter for testing.
// zh: FakeMetricWriter 是測試用的 MetricWriter 假實作。
type FakeMetricWriter struct {
	WritePointCalls []WritePointCall
	WritePointError error
}

type WritePointCall struct {
	Ctx         context.Context
	Measurement string
	Value       float64
	Labels      map[string]string
}

func (f *FakeMetricWriter) WritePoint(ctx context.Context, measurement string, value float64, labels map[string]string) error {
	f.WritePointCalls = append(f.WritePointCalls, WritePointCall{
		Ctx:         ctx,
		Measurement: measurement,
		Value:       value,
		Labels:      labels,
	})
	return f.WritePointError
}

// FakeMetricTransformer is a fake implementation of metric.MetricTransformer for testing.
// zh: FakeMetricTransformer 是測試用的 MetricTransformer 假實作。
type FakeMetricTransformer struct {
	TransformCalls []TransformCall
	TransformError error
}

type TransformCall struct {
	Measurement string
	Value       float64
	Labels      map[string]string
}

func (f *FakeMetricTransformer) Transform(measurement *string, value *float64, labels map[string]string) error {
	f.TransformCalls = append(f.TransformCalls, TransformCall{
		Measurement: *measurement,
		Value:       *value,
		Labels:      labels,
	})
	return f.TransformError
}

// FakeMetricSeriesReader is a fake implementation of metric.MetricSeriesReader for testing.
// zh: FakeMetricSeriesReader 是測試用的 MetricSeriesReader 假實作。
type FakeMetricSeriesReader struct {
	ReadSeriesCalls  []ReadSeriesCall
	ReadSeriesResult []metric.TimePoint
	ReadSeriesError  error
}

type ReadSeriesCall struct {
	Ctx    context.Context
	Expr   string
	Labels map[string]string
	Start  int64
	End    int64
}

func (f *FakeMetricSeriesReader) ReadSeries(ctx context.Context, expr string, labels map[string]string, start, end int64) ([]metric.TimePoint, error) {
	f.ReadSeriesCalls = append(f.ReadSeriesCalls, ReadSeriesCall{
		Ctx:    ctx,
		Expr:   expr,
		Labels: labels,
		Start:  start,
		End:    end,
	})
	return f.ReadSeriesResult, f.ReadSeriesError
}



================================================
FILE: internal/test/fakes/fake_modules.go
================================================
package fakes

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/modules"
)

// FakeLifecycleModule is a fake implementation of modules.LifecycleModule for testing.
// zh: FakeLifecycleModule 是測試用的 LifecycleModule 假實作。
type FakeLifecycleModule struct {
	RunCalls      []context.Context
	ShutdownCalls []context.Context
	RunError      error
	ShutdownError error
	IsHealthy     bool
}

func (f *FakeLifecycleModule) Run(ctx context.Context) error {
	f.RunCalls = append(f.RunCalls, ctx)
	return f.RunError
}

func (f *FakeLifecycleModule) Shutdown(ctx context.Context) error {
	f.ShutdownCalls = append(f.ShutdownCalls, ctx)
	return f.ShutdownError
}

func (f *FakeLifecycleModule) Healthy() bool {
	return f.IsHealthy
}

// FakeModuleEngine is a fake implementation of modules.ModuleEngine for testing.
// zh: FakeModuleEngine 是測試用的 ModuleEngine 假實作。
type FakeModuleEngine struct {
	RegisterCalls    []modules.LifecycleModule
	RunAllCalls      []context.Context
	ShutdownAllCalls []context.Context
	RunAllError      error
	ShutdownAllError error
}

func (f *FakeModuleEngine) Register(m modules.LifecycleModule) {
	f.RegisterCalls = append(f.RegisterCalls, m)
}

func (f *FakeModuleEngine) RunAll(ctx context.Context) error {
	f.RunAllCalls = append(f.RunAllCalls, ctx)
	return f.RunAllError
}

func (f *FakeModuleEngine) ShutdownAll(ctx context.Context) error {
	f.ShutdownAllCalls = append(f.ShutdownAllCalls, ctx)
	return f.ShutdownAllError
}

// FakeModuleRegistry is a fake implementation of modules.ModuleRegistry for testing.
// zh: FakeModuleRegistry 是測試用的 ModuleRegistry 假實作。
type FakeModuleRegistry struct {
	RegisterCalls []RegisterCall
	GetCalls      []string
	ListCalls     int
	RegisterError error
	GetResult     modules.LifecycleModule
	GetFound      bool
	ListResult    []string
}

type RegisterCall struct {
	Name   string
	Module modules.LifecycleModule
}

func (f *FakeModuleRegistry) Register(name string, m modules.LifecycleModule) error {
	f.RegisterCalls = append(f.RegisterCalls, RegisterCall{
		Name:   name,
		Module: m,
	})
	return f.RegisterError
}

func (f *FakeModuleRegistry) Get(name string) (modules.LifecycleModule, bool) {
	f.GetCalls = append(f.GetCalls, name)
	return f.GetResult, f.GetFound
}

func (f *FakeModuleRegistry) List() []string {
	f.ListCalls++
	return f.ListResult
}



================================================
FILE: internal/test/fakes/fake_server.go
================================================
package fakes

import (
	"context"
)

// FakeServer implements the Server interface for testing purposes.
// zh: 用於測試場景的 Server interface 假實作。
type FakeServer struct {
	RunCalled      bool
	ShutdownCalled bool
	RunError       error
	ShutdownError  error
}

// Run records that it was called and returns the configured error.
// zh: 模擬 Run 行為，並標記已呼叫狀態。
func (f *FakeServer) Run(ctx context.Context) error {
	f.RunCalled = true
	return f.RunError
}

// Shutdown records that it was called and returns the configured error.
// zh: 模擬 Shutdown 行為，並標記已呼叫狀態。
func (f *FakeServer) Shutdown(ctx context.Context) error {
	f.ShutdownCalled = true
	return f.ShutdownError
}



================================================
FILE: internal/test/plugins/eventbus/alertlog/alert_plugin_test.go
================================================
package alertlog_test

import (
	"context"
	"strings"
	"testing"

	// 匯入 plugin，觸發 init 註冊處理器
	_ "github.com/detectviz/detectviz/internal/plugins/eventbus/alertlog"

	"github.com/detectviz/detectviz/internal/registry/eventbus"
	testlogger "github.com/detectviz/detectviz/internal/test/testutil"
	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

func TestAlertPluginHandler_Registration(t *testing.T) {
	// 使用可驗證輸出的測試 logger
	log := testlogger.NewTestLogger()

	// 替換 plugin handler 的 logger（需 plugin 內部支援）
	eventbus.OverrideDefaultLogger(log)

	dispatcher, err := eventbus.NewEventDispatcher("in-memory")
	if err != nil {
		t.Fatalf("failed to create event dispatcher: %v", err)
	}

	err = dispatcher.DispatchAlertTriggered(context.Background(), event.AlertTriggeredEvent{
		AlertID:    "test-alert",
		RuleName:   "test-rule",
		Level:      "warning",
		Instance:   "test-instance",
		Metric:     "test-metric",
		Comparison: ">",
		Value:      100,
		Threshold:  100,
		Message:    "plugin alert test",
	})
	if err != nil {
		t.Errorf("unexpected dispatch error: %v", err)
	}

	// 驗證 logger 是否捕捉到 alert 處理訊息
	found := false
	for _, line := range log.Messages() {
		if strings.Contains(line, "plugin alert test") {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected alert message not found in log")
	}
}



================================================
FILE: internal/test/server/server_test.go
================================================
package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/detectviz/detectviz/internal/test/fakes"
)

// TestServer_RunAndShutdown tests the basic lifecycle of a Server implementation.
// zh: 測試 Server 實作是否能正確執行 Run 和 Shutdown。
func TestServer_RunAndShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server := &fakes.FakeServer{}

	go func() {
		_ = server.Run(ctx)
	}()

	// 模擬短時間後觸發關閉
	time.Sleep(10 * time.Millisecond)
	_ = server.Shutdown(context.Background())

	if !server.RunCalled {
		t.Errorf("Run was not called")
	}
	if !server.ShutdownCalled {
		t.Errorf("Shutdown was not called")
	}
}



================================================
FILE: internal/test/testutil/assert_logger.go
================================================
package testutil

import (
	"context"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// LogEntry 表示一筆記錄的 log 項目。
// zh: 儲存 log 訊息與相關參數。
type LogEntry struct {
	Level string
	Msg   string
	Args  []any
}

// AssertLogger 是可驗證 log 輸出的 logger 實作。
// zh: 提供簡易記錄與查詢功能，方便在測試中斷言。
type AssertLogger struct {
	mu     sync.Mutex
	events []LogEntry
}

// NewAssertLogger 建立一個新的 AssertLogger 實例。
func NewAssertLogger() *AssertLogger {
	return &AssertLogger{}
}

// Entries 回傳所有記錄的 log 項目。
func (l *AssertLogger) Entries() []LogEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]LogEntry(nil), l.events...)
}

// 以下為 logger.Logger 實作

func (l *AssertLogger) record(level, msg string, args ...any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.events = append(l.events, LogEntry{Level: level, Msg: msg, Args: args})
}

func (l *AssertLogger) Debug(msg string, args ...any) { l.record("DEBUG", msg, args...) }
func (l *AssertLogger) Info(msg string, args ...any)  { l.record("INFO", msg, args...) }
func (l *AssertLogger) Warn(msg string, args ...any)  { l.record("WARN", msg, args...) }
func (l *AssertLogger) Error(msg string, args ...any) { l.record("ERROR", msg, args...) }

func (l *AssertLogger) Named(name string) logger.Logger                { return l }
func (l *AssertLogger) WithContext(ctx context.Context) logger.Logger  { return l }
func (l *AssertLogger) WithFields(fields map[string]any) logger.Logger { return l }
func (l *AssertLogger) Sync() error                                    { return nil }



================================================
FILE: internal/test/testutil/test_logger.go
================================================
package testutil

import (
	"context"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// TestLogger 是一個可收集訊息的 logger 實作。
// zh: 測試用 logger，可用於單元測試驗證是否記錄特定訊息。
type TestLogger struct {
	mu       sync.Mutex
	messages []string
}

// NewTestLogger 回傳可收集訊息的測試 logger 實例。
// zh: 適用於需要驗證 log 行為的單元測試。
func NewTestLogger() *TestLogger {
	return &TestLogger{}
}

func (l *TestLogger) Debug(msg string, args ...any) {
	l.append(msg)
}
func (l *TestLogger) Info(msg string, args ...any) {
	l.append(msg)
}
func (l *TestLogger) Warn(msg string, args ...any) {
	l.append(msg)
}
func (l *TestLogger) Error(msg string, args ...any) {
	l.append(msg)
}
func (l *TestLogger) Named(name string) logger.Logger {
	return l
}
func (l *TestLogger) WithContext(ctx context.Context) logger.Logger {
	return l
}
func (l *TestLogger) WithFields(fields map[string]any) logger.Logger {
	return l
}
func (l *TestLogger) Sync() error {
	return nil
}

// Messages 回傳所有記錄的訊息。
// zh: 供測試時使用，驗證 logger 是否有記錄預期內容。
func (l *TestLogger) Messages() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.messages...)
}

func (l *TestLogger) append(msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, msg)
}



================================================
FILE: pkg/config/README.md
================================================
# pkg/setting:

```bash
.
├── config_anonymous.go
├── config_auth_proxy.go
├── config_azure.go
├── config_cloud_migration.go
├── config_data_proxy.go
├── config_feature_toggles.go
├── config_featuremgmt.go
├── config_grafana_javascript_agent.go
├── config_grpc.go
├── config_jwt.go
├── config_k8s_dashboard_cleanup.go
├── config_passwordless_magic_link.go
├── config_plugins.go
├── config_quota.go
├── config_remote_cache.go
├── config_search.go
├── config_secrets_manager.go
├── config_secure_socks_proxy.go
├── config_smtp.go
├── config_storage.go
├── config_unified_alerting.go
├── config_unified_storage.go
├── config.go
├── configs_rbac.go
├── configs_zanzana.go
├── configtest
├── date_formats.go
├── expanders.go
├── provider.go
└── README.md
```


`pkg/setting` 是 Grafana 中負責 統一管理組態設定（設定檔解析、結構化欄位、環境變數覆寫、預設值與驗證） 的模組。它是整個 Grafana 啟動過程的設定核心，提供其他模組依賴的 `Cfg` 結構與各項設定細節。

* * *

✅ 功能總覽
------

| 功能 | 說明 |
| --- | --- |
| 解析 `grafana.ini` 設定檔 | 使用 `ini.v1` 解析器載入與讀取分段設定 |
| 支援環境變數覆寫 | 可透過 `GF_XXX_YYY` 覆蓋 ini 中的值 |
| 結構化設定分類 | 每個主題對應一個 `setting_xxx.go` 檔案，定義專屬 struct 與初始化方法 |
| 驗證與預設值設定 | 多數欄位會做 `MustBool`, `MustInt`, `MustDuration` 等型別轉換與下限檢查 |

* * *

🧱 設定結構 (`Cfg`)
---------------

所有設定會聚合到 `setting.Cfg` 結構中：

```go
type Cfg struct {
  Raw    *ini.File
  Logger log.Logger

  AppUrl string
  Env    string
  Quota  QuotaSettings
  Plugins PluginSettings
  ...
}
```

* * *

子模組說明（部分）
------------

| 子檔案 | 功能 |
| --- | --- |
| `setting_plugins.go` | Plugin 安裝與更新策略設定grafana-pkg-all-code |
| `setting_remote_cache.go` | Redis 等快取設定（加密、prefix）grafana-pkg-all-code |
| `setting_unified_storage.go` | Unified Storage（新儲存引擎）細節控制，如 dual writer、shardinggrafana-pkg-all-code |
| `setting_anonymous.go` | 匿名登入設定，如 org\_role、限制功能grafana-pkg-all-code |
| `setting_search.go` | Dashboard reindex 與查詢效能設定grafana-pkg-all-code |
| `setting_unified_alerting.go` | Alerting cluster 設定與狀態儲存策略grafana-pkg-all-code |
| `expanders.go` | 定義支援 `env:`、`file:` 等自訂變數展開語法grafana-pkg-all-code |

* * *

🔗 與其他模組關係
----------

| 模組 | 說明 |
| --- | --- |
| `pkg/services/...` | 多數 service 會依賴 `setting.Cfg` 中的特定欄位來決定啟動邏輯 |
| `pkg/server` | 在 `main.go` 啟動流程中會先初始化 `setting.NewCfg()` |
| `pkg/plugins` | Plugin 的來源與啟用清單來自 `setting.PluginSettings` |
| `pkg/infra/log` | Logger 初始化時會依據 `setting` 中的等級與輸出位置進行設定 |

* * *

🧠 為什麼重要？
---------

Grafana 採用 設定中心 + 結構對映 + 動態解析 的設計，使得：

- 設定變更更易控管與落地
    
- 對應模組可以直接取結構，不需要再自行解析 ini
    
- 可支援未來動態設定儲存（如 cloud 控制台、GUI 編輯）
    

* * *

如你在 `detectviz` 也需要支援 `.ini`、`.env` 或動態設定系統，這個模組是非常適合參考的。我可以幫你簡化為一套 `setting` 框架版本，是否要我整理一份範本？



================================================
FILE: pkg/config/default.go
================================================
package config

import (
	"os"
	"strconv"
	"sync"

	zap_adapter "github.com/detectviz/detectviz/internal/adapters/logger"
	"github.com/detectviz/detectviz/pkg/configtypes"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
	"go.uber.org/zap"
)

// defaultProvider provides a thread-safe configuration provider that reads from an internal map and falls back to environment variables.
// zh: defaultProvider 是一個具備內部快取與 thread-safe 的設定提供者，若未設定則回傳對應的環境變數值。
type defaultProvider struct {
	configMap       map[string]string
	notifierConfigs []configtypes.NotifierConfig
	cacheConfig     configtypes.CacheConfig // 快取模組的組態設定
	log             logger.Logger
	mu              sync.RWMutex
}

// NewDefaultProvider creates a new defaultProvider instance.
// zh: 建立並回傳一個預設的設定提供者實例。
func NewDefaultProvider() *defaultProvider {
	return &defaultProvider{
		configMap: make(map[string]string),
		notifierConfigs: []configtypes.NotifierConfig{
			{Name: "email", Type: "email", Target: "noreply@example.com", Enable: true},
			{Name: "slack", Type: "slack", Target: "https://hooks.slack.com/services/xxx", Enable: true},
			{Name: "webhook", Type: "webhook", Target: "https://example.com/webhook", Enable: false},
		},
		cacheConfig: configtypes.CacheConfig{
			Backend: "memory",
			Redis: configtypes.RedisConfig{
				Address:  "localhost:6379",
				Password: "",
				DB:       0,
			},
		},
		log: zap_adapter.NewZapLogger(zap.NewNop().Sugar()), // 預設使用 Zap Logger
	}
}

// Get retrieves the value associated with the given key.
// If the key is not present in configMap, it returns the value from the environment.
// zh: 取得指定 key 的設定值，若 configMap 未命中，則回傳對應環境變數值。
func (p *defaultProvider) Get(key string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if val, ok := p.configMap[key]; ok {
		return val
	}
	return os.Getenv(key)
}

// GetInt retrieves the configuration value as an integer.
// zh: 取得指定 key 的設定值，並轉換為 int（若轉換失敗則回傳 0）。
func (p *defaultProvider) GetInt(key string) int {
	val := p.Get(key)
	i, _ := strconv.Atoi(val)
	return i
}

// GetBool retrieves the configuration value as a boolean.
// zh: 取得指定 key 的設定值，並轉換為 bool（若轉換失敗則回傳 false）。
func (p *defaultProvider) GetBool(key string) bool {
	val := p.Get(key)
	b, _ := strconv.ParseBool(val)
	return b
}

// GetOrDefault returns the value associated with the key or a provided default value if not found.
// zh: 取得指定 key 的設定值，若為空字串則回傳提供的預設值。
func (p *defaultProvider) GetOrDefault(key, defaultVal string) string {
	val := p.Get(key)
	if val == "" {
		return defaultVal
	}
	return val
}

// GetNotifierConfigs returns the list of notifier configurations.
// zh: 回傳 notifier 設定的配置列表。
func (p *defaultProvider) GetNotifierConfigs() []configtypes.NotifierConfig {
	return p.notifierConfigs
}

// GetCacheConfig returns the cache configuration.
// zh: 回傳快取模組的組態設定。
func (p *defaultProvider) GetCacheConfig() configtypes.CacheConfig {
	return p.cacheConfig
}

// Reload is a no-op for defaultProvider.
// This method is a placeholder for config reload logic, and always returns nil.
// If reload is not supported, returns nil. (See interface documentation.)
// zh: 預留重新載入設定檔功能，目前尚未實作。若不支援，則回傳 nil。
func (p *defaultProvider) Reload() error {
	return nil
}

// Set assigns a key-value pair to the configMap.
// WARNING: This method is intended for testing purposes only. Do NOT use in production code!
// zh: 寫入設定鍵值對。⚠ 僅供測試用途，請勿於正式執行路徑中使用。
func (p *defaultProvider) Set(key, val string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.configMap[key] = val
}

// Logger returns the configured logger instance.
// zh: 回傳已配置的 logger 實例。
func (p *defaultProvider) Logger() logger.Logger {
	return p.log
}



================================================
FILE: pkg/configtypes/cache_config.go
================================================
package configtypes

// CacheConfig defines the configuration for the cache module.
// zh: 支援記憶體與 Redis 快取，透過 backend 參數選擇。
type CacheConfig struct {
	Backend string      `json:"backend" yaml:"backend"` // The cache backend to use ("memory" or "redis") // 使用的快取後端（memory 或 redis）
	Redis   RedisConfig `json:"redis" yaml:"redis"`     // Redis configuration // Redis 設定
}

// RedisConfig defines the configuration for Redis cache.
// zh: 僅在 backend 設為 redis 時會使用。
type RedisConfig struct {
	Address  string `json:"address" yaml:"address"`   // Redis connection address // Redis 連線位址
	Password string `json:"password" yaml:"password"` // Redis password (can be empty) // Redis 密碼（可留空）
	DB       int    `json:"db" yaml:"db"`             // Redis DB number (default 0) // Redis 使用的 DB 編號（預設 0）
}



================================================
FILE: pkg/configtypes/notifier_config.go
================================================
package configtypes

// NotifierConfig defines the configuration for a notifier.
// zh: NotifierConfig 定義單一通知通道的設定結構。
type NotifierConfig struct {
	Name   string `json:"name"`   // zh: 通道名稱（email, slack, webhook 等）
	Type   string `json:"type"`   // zh: 通道類型
	Target string `json:"target"` // zh: 傳送目標（例如 email address、webhook URL）
	Enable bool   `json:"enable"` // zh: 是否啟用此通道
}



================================================
FILE: pkg/ifaces/alert/evaluate.go
================================================
package alert

import (
	"fmt"
)

// Evaluate 判斷查詢結果是否觸發告警條件。
// zh: 根據 AlertCondition 中指定的運算子（Operator）與閾值（Threshold），比對查詢結果 value 是否觸發告警。
func Evaluate(value float64, cond AlertCondition) (AlertResult, error) {
	operator := cond.Operator
	if operator == "" {
		operator = "ge" // 預設為 >=
	}

	var firing bool
	var message string

	switch operator {
	case "ge":
		firing = value >= cond.Threshold
		message = "threshold exceeded (>=)"
	case "gt":
		firing = value > cond.Threshold
		message = "threshold exceeded (>)"
	case "le":
		firing = value <= cond.Threshold
		message = "threshold under (<=)"
	case "lt":
		firing = value < cond.Threshold
		message = "threshold under (<)"
	case "eq":
		firing = value == cond.Threshold
		message = "threshold matched (==)"
	case "ne":
		firing = value != cond.Threshold
		message = "threshold not matched (!=)"
	default:
		return AlertResult{
			Firing:  false,
			Message: fmt.Sprintf("unsupported operator: %s", cond.Operator),
			Value:   value,
		}, fmt.Errorf("unsupported operator: %s", cond.Operator)
	}

	return AlertResult{
		Firing:  firing,
		Message: message,
		Value:   value,
	}, nil
}



================================================
FILE: pkg/ifaces/alert/evaluate_test.go
================================================
package alert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		cond     AlertCondition
		wantFire bool
		wantMsg  string
		wantErr  bool
	}{
		{"default operator (ge)", 80, AlertCondition{Threshold: 70}, true, "threshold exceeded (>=)", false},
		{"gt operator", 75, AlertCondition{Threshold: 70, Operator: "gt"}, true, "threshold exceeded (>)", false},
		{"lt operator", 60, AlertCondition{Threshold: 70, Operator: "lt"}, true, "threshold under (<)", false},
		{"eq operator", 100, AlertCondition{Threshold: 100, Operator: "eq"}, true, "threshold matched (==)", false},
		{"ne operator", 95, AlertCondition{Threshold: 100, Operator: "ne"}, true, "threshold not matched (!=)", false},
		{"unsupported operator", 42, AlertCondition{Threshold: 50, Operator: "bad"}, false, "unsupported operator: bad", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Evaluate(tt.value, tt.cond)
			assert.Equal(t, tt.wantFire, result.Firing)
			assert.Contains(t, result.Message, tt.wantMsg)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}



================================================
FILE: pkg/ifaces/alert/evaluator.go
================================================
package alert

import "context"

// AlertCondition represents the input to the evaluator, typically a rule or threshold.
// zh: AlertCondition 表示要評估的告警條件，例如規則或閾值。
type AlertCondition struct {
	RuleID    string            // zh: 規則唯一識別碼
	Expr      string            // zh: 查詢語句（如 PromQL 或 Flux）
	Operator  string            // zh: 閾值比較運算子，例如 "ge"（大於等於）、"lt"（小於）
	Threshold float64           // zh: 閾值，用於與查詢結果比對
	Labels    map[string]string // zh: 查詢附加標籤，例如 host, job 等
}

// AlertResult represents the outcome of the evaluation.
// zh: AlertResult 表示評估結果，包含是否觸發及原因。
type AlertResult struct {
	Firing  bool
	Message string
	Value   float64
}

// AlertEvaluator defines the interface for evaluating alert conditions.
// zh: AlertEvaluator 定義告警條件評估器的介面，用於根據輸入條件判斷是否觸發告警。
type AlertEvaluator interface {
	// Evaluate analyzes a condition and returns the result.
	// zh: 根據輸入條件進行評估，回傳是否觸發與原因。
	Evaluate(ctx context.Context, cond AlertCondition) (AlertResult, error)
}



================================================
FILE: pkg/ifaces/bus/alert.go
================================================
package bus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

// AlertDispatcher defines a dispatcher for alert events.
// zh: AlertDispatcher 定義用於分派告警事件的介面。
type AlertDispatcher interface {
	// DispatchAlert sends an AlertTriggeredEvent to registered handlers.
	// zh: 將 AlertTriggeredEvent 傳遞給已註冊的處理器。
	DispatchAlert(ctx context.Context, event event.AlertTriggeredEvent) error
}



================================================
FILE: pkg/ifaces/bus/host.go
================================================
package bus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

// HostDispatcher defines a dispatcher for host discovery events.
// zh: HostDispatcher 定義用於分派主機發現事件的介面。
type HostDispatcher interface {
	// DispatchHostDiscovered sends a HostDiscoveredEvent to registered handlers.
	// zh: 將 HostDiscoveredEvent 傳遞給已註冊的處理器。
	DispatchHostDiscovered(ctx context.Context, event event.HostDiscoveredEvent) error
}



================================================
FILE: pkg/ifaces/bus/metric.go
================================================
package bus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

// MetricDispatcher defines a dispatcher for metric overflow events.
// zh: MetricDispatcher 定義用於分派監控指標溢出事件的介面。
type MetricDispatcher interface {
	// DispatchMetricOverflow sends a MetricOverflowEvent to registered handlers.
	// zh: 將 MetricOverflowEvent 傳遞給已註冊的處理器。
	DispatchMetricOverflow(ctx context.Context, event event.MetricOverflowEvent) error
}



================================================
FILE: pkg/ifaces/bus/task.go
================================================
package bus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

// TaskDispatcher defines a dispatcher for task completed events.
// zh: TaskDispatcher 定義用於分派任務完成事件的介面。
type TaskDispatcher interface {
	// DispatchTaskCompleted sends a TaskCompletedEvent to registered handlers.
	// zh: 將 TaskCompletedEvent 傳遞給已註冊的處理器。
	DispatchTaskCompleted(ctx context.Context, event event.TaskCompletedEvent) error
}



================================================
FILE: pkg/ifaces/bus/types.go
================================================
package bus

import "time"

// Envelope defines a wrapper for dispatched events.
// zh: Envelope 用來包裝傳遞的事件與其附加資訊，供 Dispatcher 使用。
type Envelope struct {
	EventType string      // zh: 事件類型，例如 "host.discovered"
	Payload   interface{} // zh: 真實事件資料內容，可為任意型別
	Timestamp time.Time   // zh: 發送時間，用於排序或延遲處理用途
}



================================================
FILE: pkg/ifaces/cachestore/cachestore.go
================================================
package cachestore

// CacheStore defines the abstract interface for caching used in Detectviz.
// It supports basic key-value operations with TTL (time-to-live) and can be implemented using in-memory or Redis-based stores.
// zh: CacheStore 定義 Detectviz 中的快取抽象介面，支援基本的鍵值操作與 TTL（存活時間），可對接 memory 或 Redis 實作。
type CacheStore interface {
	// Get returns the cached value for the given key, or an error if not found.
	// zh: 根據 key 取得對應快取內容，若 key 不存在應回傳錯誤。
	Get(key string) (string, error)

	// Set sets the cache value for a given key with a TTL in seconds.
	// If ttlSeconds is 0, the value never expires.
	// zh: 設定快取內容與存活時間，ttlSeconds 為 0 時表示永不過期。
	Set(key string, val string, ttlSeconds int) error

	// Has returns true if the given key exists in cache.
	// If an error occurs during the check, it returns false and the error.
	// zh: 檢查 key 是否存在於快取中，非強一致性。若查詢錯誤則回傳 error。
	Has(key string) (bool, error)

	// Delete removes the value associated with the given key.
	// zh: 從快取中移除指定的 key。
	Delete(key string) error

	// Keys returns all cache keys that start with the given prefix.
	// Optional method for grouping or bulk invalidation use cases.
	// zh: 回傳所有以指定 prefix 開頭的 key，常用於群組清除。
	Keys(prefix string) ([]string, error)
}



================================================
FILE: pkg/ifaces/config/config.go
================================================
package config

import (
	"github.com/detectviz/detectviz/pkg/configtypes"
	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

// ConfigProvider defines the interface for retrieving configuration values.
// zh: ConfigProvider 定義取得設定值的抽象介面，可支援從環境變數、檔案或遠端服務載入。
type ConfigProvider interface {
	// Get returns the string value for a given key.
	// zh: 根據指定 key 取得對應的字串設定值。
	Get(key string) string

	// GetOrDefault returns the value for a given key, or returns the provided default if not found.
	// zh: 根據 key 取得設定值，若無對應值則回傳預設值。
	GetOrDefault(key, defaultVal string) string

	// GetInt returns the integer value for a given key.
	// zh: 根據指定 key 取得整數型別的設定值。
	GetInt(key string) int

	// GetBool returns the boolean value for a given key.
	// zh: 根據指定 key 取得布林型別的設定值。
	GetBool(key string) bool

	// GetCacheConfig returns the cache module configuration.
	// zh: 回傳快取模組的組態設定。
	GetCacheConfig() configtypes.CacheConfig

	// GetNotifierConfigs returns the list of notifier configurations.
	// zh: 回傳 notifier 設定的配置列表。
	GetNotifierConfigs() []configtypes.NotifierConfig

	// Logger returns the configured logger instance.
	// zh: 回傳已配置的 logger 實例。
	Logger() logger.Logger

	// Reload refreshes the underlying configuration source, if supported.
	// If hot-reload is unsupported, this may be a no-op or return nil.
	// zh: 重新載入設定來源內容（若支援），常用於檔案或環境變數動態更新；若不支援，可能為空操作。
	Reload() error
}



================================================
FILE: pkg/ifaces/event/alert.go
================================================
package event

import "context"

// AlertEventHandler defines the handler interface for AlertTriggeredEvent.
// zh: AlertEventHandler 定義處理 AlertTriggeredEvent 的介面。
//
// Used in event dispatch systems to register alert-related handlers.
// zh: 用於事件分派系統中註冊處理告警事件的 handler。
//
// Event definition: see pkg/ifaces/event/types.go
// zh: 事件資料結構定義請參考 pkg/ifaces/event/types.go
type AlertEventHandler interface {
	// HandleAlertTriggered processes the AlertTriggeredEvent.
	// zh: 處理告警事件的實作函式。
	HandleAlertTriggered(ctx context.Context, event AlertTriggeredEvent) error
}



================================================
FILE: pkg/ifaces/event/host.go
================================================
package event

import "context"

// HostEventHandler defines the handler interface for HostDiscoveredEvent.
// zh: HostEventHandler 定義處理 HostDiscoveredEvent 的事件處理器介面。
//
// This interface is typically registered to an EventBus dispatcher.
// zh: 本介面通常會註冊至 EventBus 分派器以接收主機註冊事件。
//
// Event definition: see pkg/ifaces/event/types.go
// zh: 事件資料結構定義請參考 pkg/ifaces/event/types.go
type HostEventHandler interface {
	// HandleHostDiscovered processes the HostDiscoveredEvent.
	// zh: 處理主機註冊事件的實作函式。
	HandleHostDiscovered(ctx context.Context, event HostDiscoveredEvent) error
}



================================================
FILE: pkg/ifaces/event/metric.go
================================================
package event

import "context"

// MetricEventHandler defines the handler interface for MetricOverflowEvent.
// zh: MetricEventHandler 定義處理 MetricOverflowEvent 的事件處理器介面。
//
// This interface is used in the EventBus to register handlers for metric overflow conditions.
// zh: 本介面用於 EventBus 中註冊處理指標溢出事件的 handler。
//
// Event definition: see pkg/ifaces/event/types.go
// zh: 事件資料結構定義請參考 pkg/ifaces/event/types.go
type MetricEventHandler interface {
	// HandleMetricOverflow processes the MetricOverflowEvent.
	// zh: 處理監控指標溢出事件的實作函式。
	HandleMetricOverflow(ctx context.Context, event MetricOverflowEvent) error
}



================================================
FILE: pkg/ifaces/event/task.go
================================================
package event

import "context"

// TaskEventHandler defines the handler interface for TaskCompletedEvent.
// zh: TaskEventHandler 定義處理 TaskCompletedEvent 的介面。
//
// This interface is used to handle task execution completion events via the EventBus dispatcher.
// zh: 本介面用於透過 EventBus 分派器處理任務執行完成事件。
//
// Event definition: see pkg/ifaces/event/types.go
// zh: 事件資料結構定義請參考 pkg/ifaces/event/types.go
type TaskEventHandler interface {
	// HandleTaskCompleted processes the TaskCompletedEvent.
	// zh: 處理任務完成事件的實作函式。
	HandleTaskCompleted(ctx context.Context, event TaskCompletedEvent) error
}



================================================
FILE: pkg/ifaces/event/types.go
================================================
package event

import "time"

// AlertTriggeredEvent represents an alert event that occurred in the system.
// zh: AlertTriggeredEvent 表示系統中觸發的告警事件。
type AlertTriggeredEvent struct {
	EventID    string    // zh: 事件唯一識別碼
	Timestamp  time.Time // zh: 事件發生時間
	AlertID    string    // zh: 告警事件 ID
	RuleName   string    // zh: 告警規則名稱
	Level      string    // zh: 告警嚴重程度，例如 critical、warning
	Instance   string    // zh: 實體名稱或識別，例如設備名稱
	Metric     string    // zh: 指標名稱
	Comparison string    // zh: 比較運算符，例如 >、<
	Value      float64   // zh: 實際值
	Threshold  float64   // zh: 閾值
	Message    string    // zh: 告警訊息內容
}

// TaskCompletedEvent represents a task completion event in the system.
// zh: TaskCompletedEvent 表示系統中某個任務完成的事件。
type TaskCompletedEvent struct {
	EventID   string    // zh: 事件唯一識別碼
	Timestamp time.Time // zh: 事件發生時間
	TaskID    string    // zh: 任務識別碼
	WorkerID  string    // zh: 執行任務的工作者 ID
	Status    string    // zh: 任務完成狀態，例如 success、failed
}

// HostDiscoveredEvent represents the discovery or registration of a host in the system.
// zh: HostDiscoveredEvent 表示系統中主機被發現或註冊的事件。
type HostDiscoveredEvent struct {
	EventID   string            // zh: 事件唯一識別碼
	Timestamp time.Time         // zh: 事件發生時間
	HostID    string            // zh: 主機識別碼
	Name      string            // zh: 主機名稱
	IP        string            // zh: 主機 IP 位址
	Source    string            // zh: 來源識別，例如由哪個掃描器或子系統發現
	Labels    map[string]string // zh: 附加標籤，例如 rack、zone 等
}

// MetricOverflowEvent represents an overflow condition in a monitored metric.
// zh: MetricOverflowEvent 表示某個監控指標超出預期範圍的事件。
type MetricOverflowEvent struct {
	EventID    string    // zh: 事件唯一識別碼
	Timestamp  time.Time // zh: 事件發生時間
	MetricName string    // zh: 指標名稱
	Value      float64   // zh: 實際值
	Threshold  float64   // zh: 閾值
	Source     string    // zh: 數據來源（原欄位，保留以維持相容）
	Instance   string    // zh: 實體名稱或識別，例如設備名稱
	Reason     string    // zh: 溢出的原因說明
}



================================================
FILE: pkg/ifaces/eventbus/eventbus.go
================================================
package eventbus

import (
	"context"

	"github.com/detectviz/detectviz/pkg/ifaces/event"
)

// EventDispatcher defines the interface for dispatching typed events between modules.
// zh: EventDispatcher 定義模組間明確事件分派的抽象介面。
type EventDispatcher interface {
	// -------------------------------------------------------------------------
	// AlertTriggeredEvent
	// -------------------------------------------------------------------------

	// DispatchAlertTriggered dispatches an AlertTriggeredEvent to all registered handlers.
	// zh: 將 AlertTriggeredEvent 分派給所有已註冊的處理器。
	DispatchAlertTriggered(ctx context.Context, e event.AlertTriggeredEvent) error

	// RegisterAlertHandler registers a new handler for AlertTriggeredEvent.
	// zh: 註冊一個處理 AlertTriggeredEvent 的事件處理器。
	RegisterAlertHandler(handler event.AlertEventHandler)

	// -------------------------------------------------------------------------
	// TaskCompletedEvent
	// -------------------------------------------------------------------------

	// DispatchTaskCompleted dispatches a TaskCompletedEvent to all registered handlers.
	// zh: 將 TaskCompletedEvent 分派給所有已註冊的處理器。
	DispatchTaskCompleted(ctx context.Context, e event.TaskCompletedEvent) error

	// RegisterTaskHandler registers a new handler for TaskCompletedEvent.
	// zh: 註冊一個處理 TaskCompletedEvent 的事件處理器。
	RegisterTaskHandler(handler event.TaskEventHandler)

	// -------------------------------------------------------------------------
	// HostDiscoveredEvent
	// -------------------------------------------------------------------------

	// DispatchHostDiscovered dispatches a HostDiscoveredEvent to all registered handlers.
	// zh: 將 HostDiscoveredEvent 分派給所有已註冊的處理器。
	DispatchHostDiscovered(ctx context.Context, e event.HostDiscoveredEvent) error

	// RegisterHostHandler registers a new handler for HostDiscoveredEvent.
	// zh: 註冊一個處理 HostDiscoveredEvent 的事件處理器。
	RegisterHostHandler(handler event.HostEventHandler)

	// -------------------------------------------------------------------------
	// MetricOverflowEvent
	// -------------------------------------------------------------------------

	// DispatchMetricOverflow dispatches a MetricOverflowEvent to all registered handlers.
	// zh: 將 MetricOverflowEvent 分派給所有已註冊的處理器。
	DispatchMetricOverflow(ctx context.Context, e event.MetricOverflowEvent) error

	// RegisterMetricHandler registers a new handler for MetricOverflowEvent.
	// zh: 註冊一個處理 MetricOverflowEvent 的事件處理器。
	RegisterMetricHandler(handler event.MetricEventHandler)
}

// AlertEventHandler defines a handler for AlertTriggeredEvent.
// zh: AlertEventHandler 定義處理 AlertTriggeredEvent 的事件處理器。
type AlertEventHandler interface {
	// HandleAlertTriggered processes the given alert event.
	// zh: 接收並處理 AlertTriggeredEvent。
	HandleAlertTriggered(ctx context.Context, event event.AlertTriggeredEvent) error
}

// TaskEventHandler defines a handler for TaskCompletedEvent.
// zh: TaskEventHandler 定義處理 TaskCompletedEvent 的事件處理器。
type TaskEventHandler interface {
	// HandleTaskCompleted processes the completed task event.
	// zh: 接收並處理 TaskCompletedEvent。
	HandleTaskCompleted(ctx context.Context, event event.TaskCompletedEvent) error
}

// HostEventHandler defines a handler for HostDiscoveredEvent.
// zh: HostEventHandler 定義處理 HostDiscoveredEvent 的事件處理器。
type HostEventHandler interface {
	// HandleHostDiscovered processes the discovered host event.
	// zh: 接收並處理 HostDiscoveredEvent。
	HandleHostDiscovered(ctx context.Context, event event.HostDiscoveredEvent) error
}

// MetricEventHandler defines a handler for MetricOverflowEvent.
// zh: MetricEventHandler 定義處理 MetricOverflowEvent 的事件處理器。
type MetricEventHandler interface {
	// HandleMetricOverflow processes the metric overflow event.
	// zh: 接收並處理 MetricOverflowEvent。
	HandleMetricOverflow(ctx context.Context, event event.MetricOverflowEvent) error
}

// RegisterPluginTaskHandler 是 plugin 註冊 TaskEventHandler 的統一介面。
// zh: 提供 plugin 註冊 TaskCompletedEvent 處理器。
type RegisterPluginTaskHandler interface {
	RegisterTaskHandler(handler TaskEventHandler)
}

// RegisterPluginHostHandler 是 plugin 註冊 HostEventHandler 的統一介面。
// zh: 提供 plugin 註冊 HostDiscoveredEvent 處理器。
type RegisterPluginHostHandler interface {
	RegisterHostHandler(handler HostEventHandler)
}

// RegisterPluginMetricHandler 是 plugin 註冊 MetricEventHandler 的統一介面。
// zh: 提供 plugin 註冊 MetricOverflowEvent 處理器。
type RegisterPluginMetricHandler interface {
	RegisterMetricHandler(handler MetricEventHandler)
}



================================================
FILE: pkg/ifaces/eventbus/provider.go
================================================
package eventbus

// DispatcherProvider 為 eventbus 後端實作註冊器
type DispatcherProvider interface {
	// Name 傳回 provider 名稱（如 in-memory, kafka）
	Name() string

	// Build 建構對應的 EventDispatcher
	Build() EventDispatcher
}



================================================
FILE: pkg/ifaces/logger/context.go
================================================
package logger

import "context"

// ctxKey 是用於 context 注入與擷取 logger 的專用 key。
// zh: ctxKey 用於避免與其他 context 欄位衝突。
type ctxKey struct{}

// WithContext 將 logger 實例注入至 context 中，供下游模組擷取使用。
// zh: 建議於 middleware 或 handler 中呼叫，讓後續模組皆可透過 FromContext 擷取同一 logger。
func WithContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext 從 context 中擷取 logger 實例，若未注入則回傳 NopLogger。
// zh: 此函式可保證不會回傳 nil，避免程式 panic。
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(ctxKey{}).(Logger); ok {
		return l
	}
	return NopLogger{} // fallback 預設 logger，靜默處理所有輸出
}



================================================
FILE: pkg/ifaces/logger/context_test.go
================================================
package logger_test

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/pkg/ifaces/logger"
)

type testLogger struct {
	called bool
}

func (t *testLogger) Debug(msg string, args ...any)                 { t.called = true }
func (t *testLogger) Info(msg string, args ...any)                  { t.called = true }
func (t *testLogger) Warn(msg string, args ...any)                  { t.called = true }
func (t *testLogger) Error(msg string, args ...any)                 { t.called = true }
func (t *testLogger) Sync() error                                   { return nil }
func (t *testLogger) WithFields(map[string]any) logger.Logger       { return t }
func (t *testLogger) WithContext(ctx context.Context) logger.Logger { return t }
func (t *testLogger) Named(string) logger.Logger                    { return t }

func TestFromContext_Default(t *testing.T) {
	ctx := context.Background()
	l := logger.FromContext(ctx)

	if l == nil {
		t.Fatal("expected non-nil fallback logger")
	}
}

func TestWithContext_InjectAndRetrieve(t *testing.T) {
	ctx := context.Background()
	mock := &testLogger{}
	ctx = logger.WithContext(ctx, mock)

	l := logger.FromContext(ctx)
	l.Info("test")

	if !mock.called {
		t.Error("expected injected logger to be called")
	}
}



================================================
FILE: pkg/ifaces/logger/logger.go
================================================
package logger

import "context"

// Logger defines the abstract interface for structured logging in Detectviz.
// This interface is designed to support OpenTelemetry-compatible logging,
// allowing correlation with trace and metric data.
// zh: Logger 定義 Detectviz 中結構化日誌的抽象介面。此介面設計支援 OpenTelemetry 相容的日誌功能，允許與追蹤與指標數據做關聯。
type Logger interface {
	// Info logs a message at the info level with optional structured fields.
	// zh: 記錄 info 級別的日誌，可附帶結構化欄位。
	Info(msg string, fields ...any)

	// Warn logs a message at the warning level.
	// zh: 記錄 warning 級別的日誌。
	Warn(msg string, fields ...any)

	// Error logs a message at the error level.
	// zh: 記錄 error 級別的日誌。
	Error(msg string, fields ...any)

	// Debug logs a message at the debug level.
	// zh: 記錄 debug 級別的日誌。
	Debug(msg string, fields ...any)

	// WithFields returns a logger with the provided structured fields included.
	// zh: 回傳一個包含指定結構化欄位的新 logger 實例。
	WithFields(fields map[string]any) Logger

	// WithContext returns a logger that extracts trace context from the given context.
	// zh: 從傳入的 context 中提取追蹤資訊（如 trace_id、span_id），並回傳帶有 context 的 logger。
	WithContext(ctx context.Context) Logger

	// Named returns a logger with an assigned name, typically per module or component.
	// zh: 回傳一個命名的 logger，常用於模組或元件名稱區分。
	Named(name string) Logger

	// Sync flushes any buffered log entries, if supported.
	// zh: 若 logger 支援緩衝區，則強制將其清空寫出。
	Sync() error
}



================================================
FILE: pkg/ifaces/logger/nop_logger.go
================================================
package logger

import "context"

// NopLogger implements Logger with no-op methods.
// zh: NopLogger 是 Logger 的空實作，不會輸出任何日誌。
type NopLogger struct{}

// Debug does nothing.
// zh: Debug 級別，不輸出內容。
func (NopLogger) Debug(msg string, args ...any) {}

// Sync does nothing and returns nil.
// zh: 不需 flush，直接回傳 nil。
func (NopLogger) Sync() error { return nil }

// WithFields returns NopLogger.
// zh: 回傳自身，不套用任何欄位。
func (NopLogger) WithFields(fields map[string]any) Logger {
	return NopLogger{}
}

// WithContext returns NopLogger.
// zh: 回傳自身，不注入 context。
func (NopLogger) WithContext(ctx context.Context) Logger {
	return NopLogger{}
}

// Named returns NopLogger.
// zh: 回傳自身，不套用 logger 名稱。
func (NopLogger) Named(name string) Logger {
	return NopLogger{}
}

// Error does nothing.
// zh: Error 級別，不輸出內容。
func (NopLogger) Error(msg string, args ...any) {}

// Info does nothing.
// zh: Info 級別，不輸出內容。
func (NopLogger) Info(msg string, args ...any) {}

// Warn does nothing.
// zh: Warn 級別，不輸出內容。
func (NopLogger) Warn(msg string, args ...any) {}



================================================
FILE: pkg/ifaces/metrics/metric.go
================================================
// Package metric provides shared metric-related interfaces for the detectviz project.
// zh: 提供 Detectviz 專案中與指標資料相關的共用介面。

package metric

import (
	"context"
)

// MetricWriter defines the interface for sending metric data to external systems.
// zh: MetricWriter 定義寫入指標資料至外部系統的介面（例如 InfluxDB、Pushgateway）。
type MetricWriter interface {
	// WritePoint writes a single metric point with measurement name, value, and labels.
	// zh: 寫入單筆指標資料，包含量測名稱、數值與標籤。
	WritePoint(ctx context.Context, measurement string, value float64, labels map[string]string) error
}

// MetricTransformer defines the interface for preprocessing metric data before evaluation or storage.
// zh: MetricTransformer 定義指標資料在評估或儲存前的預處理邏輯（例如單位轉換、標籤增補）。
type MetricTransformer interface {
	// Transform modifies the measurement name, value, and labels in-place before processing.
	// zh: 對指標名稱、數值與標籤進行轉換處理，會就地修改輸入值。
	Transform(measurement *string, value *float64, labels map[string]string) error
}

// MetricSeriesReader defines the interface for reading a time series of metric data.
// zh: MetricSeriesReader 定義讀取時間序列資料的介面，常用於報表、圖表或趨勢分析。
type MetricSeriesReader interface {
	// ReadSeries returns a list of timestamped values for a given expression and labels within a time range.
	// zh: 讀取指定表達式與標籤條件的時間序列資料，回傳時間戳與對應數值清單。
	ReadSeries(ctx context.Context, expr string, labels map[string]string, start, end int64) ([]TimePoint, error)
}



================================================
FILE: pkg/ifaces/metrics/query.go
================================================
package metric

import "context"

// MetricQueryAdapter defines the interface for querying metric values from various data sources.
// zh: MetricQueryAdapter 定義從各種資料來源查詢監控指標的介面。
type MetricQueryAdapter interface {
	// Query executes a query expression with optional label filters and returns a numeric result.
	// zh: 執行查詢語句並搭配可選的標籤篩選，回傳單一數值結果。
	Query(ctx context.Context, expr string, labels map[string]string) (float64, error)
}



================================================
FILE: pkg/ifaces/metrics/types.go
================================================
package metric

// TimePoint represents a single point in a time series.
// zh: TimePoint 表示時間序列中的一個時間點。
type TimePoint struct {
	Timestamp int64
	Value     float64
}

// TimeRange represents a start and end timestamp for querying time series data.
// zh: TimeRange 表示查詢時間序列資料時使用的起始與結束時間範圍。
type TimeRange struct {
	Start int64
	End   int64
}



================================================
FILE: pkg/ifaces/modules/modules.go
================================================
package modules

import "context"

// LifecycleModule defines a basic lifecycle interface for modules.
// zh: 定義模組基本生命週期介面。
type LifecycleModule interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// HealthCheckableModule adds health check capabilities to a lifecycle module.
// zh: 提供健康檢查能力的模組，擴充自 LifecycleModule。
type HealthCheckableModule interface {
	LifecycleModule
	Healthy() bool
}

// ModuleEngine coordinates registration and control of unnamed modules.
// zh: 控制匿名模組註冊與啟動的引擎。
type ModuleEngine interface {
	Register(m LifecycleModule)
	RunAll(ctx context.Context) error
	ShutdownAll(ctx context.Context) error
}

// ModuleRegistry manages named module registration and lookup.
// zh: 管理具名模組的註冊與查詢。
type ModuleRegistry interface {
	Register(name string, m LifecycleModule) error
	Get(name string) (LifecycleModule, bool)
	List() []string
}

// ModuleRunner controls the ordered startup and shutdown based on dependencies.
// zh: 根據依賴圖控制模組啟動與關閉順序。
type ModuleRunner interface {
	StartAll(ctx context.Context) error
	StopAll(ctx context.Context) error
}

// ModuleListener monitors health status of registered modules.
// zh: 監控模組健康狀態的監聽器。
type ModuleListener interface {
	Start(ctx context.Context)
	Stop()
}



================================================
FILE: pkg/ifaces/notifier/notifier.go
================================================
package notifier

import (
	"context"
	"time"
)

// Message represents a notification message to be sent.
// zh: Message 表示要傳送的通知訊息。
type Message struct {
	Title   string            // zh: 訊息標題
	Content string            // zh: 訊息內容
	Labels  map[string]string // zh: 附加標籤（例如等級、來源模組）
	Target  string            // zh: 接收對象或通道，例如 email、webhook URL
	Time    time.Time         // zh: 發送時間，預設為訊息產生時間
}

// Notifier defines an interface for sending notifications to various channels.
// zh: Notifier 定義通知傳送的介面，可擴充為多種通道實作（如 Email、Slack、Webhook）。
type Notifier interface {
	// Name returns the identifier of the notifier.
	// zh: 回傳 notifier 名稱。
	Name() string

	// Send delivers the message via this notifier channel.
	// zh: 傳送完整訊息結構至此通道，包含標籤與時間等欄位。
	Send(ctx context.Context, msg Message) error

	// Notify delivers a simple title-message pair.
	// zh: 傳送簡易通知（標題與訊息），通常用於預設用途或簡化介面。
	Notify(title, message string) error
}



================================================
FILE: pkg/ifaces/scheduler/mock_adapter_test.go
================================================
package scheduler_test

import (
	"context"
	"testing"

	"github.com/detectviz/detectviz/pkg/ifaces/scheduler"
	"github.com/stretchr/testify/assert"
)

// mockJob implements the scheduler.Job interface for testing.
// zh: mockJob 是模擬用的排程任務，實作 scheduler.Job 介面。
type mockJob struct {
	name  string // zh: 任務名稱
	spec  string // zh: 排程時間（cron 格式）
	calls int    // zh: 被執行次數
}

// Run increments the call counter.
// zh: Run 被呼叫時會遞增 calls 計數。
func (m *mockJob) Run(ctx context.Context) error {
	m.calls++
	return nil
}

// Name returns the job name.
// zh: 回傳任務名稱。
func (m *mockJob) Name() string {
	return m.name
}

// Spec returns the job's schedule.
// zh: 回傳排程規則。
func (m *mockJob) Spec() string {
	return m.spec
}

// mockScheduler is a lightweight mock implementation of Scheduler.
// zh: mockScheduler 是模擬用的 Scheduler 實作，不會實際執行任務。
type mockScheduler struct {
	Jobs []scheduler.Job // Registered jobs. zh: 已註冊的任務清單
}

// Register adds a job to the scheduler.
// zh: 註冊一個任務。
func (m *mockScheduler) Register(job scheduler.Job) {
	m.Jobs = append(m.Jobs, job)
}

// Start is a no-op.
// zh: 啟動排程器（模擬，無實際行為）。
func (m *mockScheduler) Start(ctx context.Context) error {
	return nil
}

// Stop is a no-op.
// zh: 停止排程器（模擬，無實際行為）。
func (m *mockScheduler) Stop(ctx context.Context) error {
	return nil
}

// TestMockSchedulerIntegration tests the registration and lifecycle flow of the mock scheduler.
// zh: 測試 mockScheduler 的任務註冊與啟停流程。
func TestMockSchedulerIntegration(t *testing.T) {
	job := &mockJob{name: "test-job", spec: "@every 1m"}
	mockSched := &mockScheduler{}

	mockSched.Register(job)

	assert.Len(t, mockSched.Jobs, 1)
	assert.Equal(t, "test-job", mockSched.Jobs[0].Name())

	err := mockSched.Start(context.Background())
	assert.NoError(t, err)

	err = mockSched.Stop(context.Background())
	assert.NoError(t, err)
}



================================================
FILE: pkg/ifaces/scheduler/scheduler.go
================================================
package scheduler

import "context"

// Job represents a task that can be scheduled and executed.
// zh: Job 代表一個可排程執行的任務。
type Job interface {
	// Run executes the job logic.
	// zh: 執行任務邏輯。
	Run(ctx context.Context) error

	// Name returns the job name.
	// zh: 回傳任務名稱（用於識別與日誌）。
	Name() string

	// Spec returns the cron-style schedule string for this job.
	// zh: 回傳此任務的排程時間字串（cron 格式）。
	Spec() string
}

// Scheduler defines the interface for a job scheduler.
// zh: Scheduler 定義排程器的操作介面。
type Scheduler interface {
	// Register registers a job with the scheduler.
	// zh: 註冊任務至排程器中。
	Register(job Job)

	// Start initiates the job scheduler.
	// zh: 啟動排程器。
	Start(ctx context.Context) error

	// Stop gracefully stops the scheduler and all running jobs.
	// zh: 停止排程器與所有正在執行的任務。
	Stop(ctx context.Context) error
}



================================================
FILE: pkg/ifaces/server/server.go
================================================
package server

import "context"

// Server defines the lifecycle of the main application server.
// zh: 定義伺服器核心生命週期的介面。
type Server interface {
	// Run starts the server and blocks until context is cancelled or an error occurs.
	// zh: 啟動伺服器，阻塞直到收到關閉訊號或發生錯誤。
	Run(ctx context.Context) error

	// Shutdown gracefully shuts down the server and its components.
	// zh: 優雅地關閉伺服器與所有相依元件。
	Shutdown(ctx context.Context) error
}


