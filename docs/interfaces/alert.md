# AlertEvaluator Interface 說明文件

> 本文件為 Detectviz 專案中 `AlertEvaluator` 告警評估介面的設計說明與擴充建議，依據 Clean Architecture 原則，資料查詢行為由 `MetricQueryAdapter` 注入，保持告警邏輯與資料來源解耦。

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
- 生產環境可共用一套評估器，但對接不同資料來源（Prometheus、InfluxDB 等）

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

| 檔案                          | 功能描述                               |
|-------------------------------|----------------------------------------|
| `internal/adapters/alert/static.go` | 靜態值或固定閾值比對邏輯（單元測試常用） |
| `internal/adapters/alert/nop.go`    | 不進行任何評估的 Noop 實作               |

註：告警模組不需針對每種資料來源製作獨立 Evaluator，皆透過注入 `MetricQueryAdapter` 統一處理資料查詢，邏輯不變。

---

## 擴充資料來源時應實作哪些部分（How to add a new data source）

### 1. 實作 MetricQueryAdapter（參見 `docs/interfaces/metric.md`）

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

| 欄位     | 說明                 |
|----------|----------------------|
| ID       | 條件識別碼           |
| Expr     | 查詢語法（如 Flux）  |
| Threshold | 閾值數值            |
| Labels   | 過濾條件的標籤組     |

### AlertResult

| 欄位     | 說明                   |
|----------|------------------------|
| Firing   | 是否觸發告警           |
| Value    | 查詢到的實際值         |
| Message  | 描述文字，可供記錄與通知 |

---

## 關聯模組與擴充性（Related & extensibility）

- 評估器由 `AlertScheduler` 控制執行時機
- 結果可拋給 `EventDispatcher` 通知對應管道
- 下層資料查詢可使用統一介面（`MetricQueryAdapter`）自由替換與擴充