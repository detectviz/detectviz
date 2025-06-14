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
- 所有 interface 定義應集中於 `pkg/platform/contracts/`，並維持穩定與可測試性。
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

- 為避免與 `pkg/platform/contracts/metric` 等核心介面衝突，建議 adapter 的實作包使用 `package xxxadapter` 命名，例如 `metricsadapter`, `loggeradapter`
- 這有助於區分：
  - 核心抽象層（如 `metric.Writer`、`logger.Logger`）
  - 對應的實作包（如 `internal/adapters/metrics`）
- 匯入時可維持一致，例如：

```go
import metricsadapter "detectviz/internal/adapters/metrics"
```

- 適用情境：
  - 該模組有清楚對應 interface（如 `pkg/platform/contracts/metric.Writer`）
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