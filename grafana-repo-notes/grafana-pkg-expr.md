# grafana/pkg/expr:
`grafana/pkg/expr` 是 Grafana 的「表達式引擎核心模組」，它負責將多種查詢來源（如 SQL、Prometheus、InfluxDB）與數學邏輯（如 reduce、math、threshold）**統一建模為一條 expression pipeline**，進行 alert 判斷或 dashboard 計算。

以下是重點解析：

* * *

✅ 主要功能
------

| 功能 | 說明 |
| --- | --- |
| 表達式建構 | 接收 alerting rule 中的 condition 與查詢資料，轉換為 `DataPipeline` 結構 |
| Pipeline 節點定義 | 每個 query/refId 被轉換為一個 `Node`，可能是 datasource、math、reduce、threshold |
| 表達式執行 | 將 pipeline 傳入 evaluator 執行，回傳結果與狀態 |
| Alert 判斷邏輯 | 與 `ngalert` 結合，在 `eval.go` 中呼叫 `expr.Service.BuildPipeline`grafana-pkg-services-ng… |

* * *

🧩 結構與模組子目錄（從其他模組引用觀察）
----------------------

1.  `expr.Service`
    
    *   提供 `BuildPipeline(req)` 方法：回傳 `DataPipeline`。
        
    *   實作 `expressionExecutor` interface：能夠執行 pipeline。
        
2.  子模組（根據引入行為）：
    
    *   `expr/classic`: 支援經典告警規則
        
    *   `expr/mathexp`: 提供 reducer / math / resample 計算邏輯
        
    *   `expr/node.go`: 定義各種 Node 類型與 pipeline 建構規則
        

* * *

🔌 Interface 定義（來自 `ngalert/eval.go` 中的依賴）
------------------------------------------

```go
type expressionExecutor interface {
    ExecutePipeline(ctx context.Context, now time.Time, pipeline expr.DataPipeline) (*backend.QueryDataResponse, error)
}

type expressionBuilder interface {
    expressionExecutor
    BuildPipeline(req *expr.Request) (expr.DataPipeline, error)
}
```

* * *

📦 常見 Node 類型（由 `expr.NodeTypeFromDatasourceUID()` 推導）
------------------------------------------------------

| Node 類型 | 描述 |
| --- | --- |
| `TypeDatasourceNode` | 一般查詢類型（Prometheus, MySQL 等） |
| `TypeMathNode` | 數學運算節點 |
| `TypeReduceNode` | 聚合函數（如 avg, max） |
| `TypeThresholdNode` | 評估閾值並回傳判斷狀態 |
| `TypeMLNode` | 機器學習運算節點 |
| `TypeCMDNode` | 指令型運算（特殊用途） |

* * *

🧠 在 Grafana 架構中的角色
-------------------

| 模組 | 用途 |
| --- | --- |
| `ngalert` | 主要使用者：建立與執行 alert rule 的條件判斷 |
| `dashboards` | 匯入儀表板時若有表達式，強制套用 `__expr__` 作為 pluginIdgrafana-pkg-services-al… |
| `eval_result` | 將 `classic.EvalMatch` 或 `mathexp.NumberValueCapture` 整合為前端呈現grafana-pkg-services-ng… |

* * *

如你有需要，我可以補上 `expr.Service` 的具體實作（如 `BuildPipeline()` 的行為流程），或幫你畫出 node pipeline 的資料流程圖。是否繼續？




# 遷移到 internal/alert/expr/ 後怎麼做

grafana/pkg/expr 是一個 表達式建構與運算引擎，功能高度集中在 告警條件判斷與數據轉換邏輯。若依照你正在建構的 Detectviz 架構分層原則，它應該：

⸻

✅ 遷移建議位置：internal/alert/expr/

⸻

📌 理由說明

項目	說明
📦 功能定位	屬於告警條件邏輯層，建構 DataPipeline、運行 reducer / math / threshold
🔗 使用模組	僅被 ngalert、alerting、dashboard 模組使用，高度專用於告警或計算應用邏輯
🔒 可測試性與邊界	可作為獨立單元測試，但無 plugin 化必要，不需開放 external 使用
❌ 非共用工具庫	雖具有 reusable 結構（如節點定義），但不具「跨應用共通性」
💡 與 alert/ 其他模組邏輯緊密	eval.go, rule.go 等皆與 pipeline 結構緊密結合，非獨立模組


⸻

🧱 建議遷移結構（Detectviz 架構）

internal/alert/
├── expr/
│   ├── service.go         # 原 expr.Service，含 BuildPipeline、Execute
│   ├── node.go            # Node 類型定義
│   ├── mathexp/           # math, reduce, resample
│   ├── classic/           # 經典 rule 兼容

若你未來要 plugin 化支援類似 ML Node 或特殊 DSL，也可以將：
	•	internal/alert/expr/plugin/：實作 MLNode, CMDNode
	•	pkg/iface/alert/expr.go：定義 PipelineBuilder, NodeExecutor

⸻

🚫 不建議放在：

路徑	原因
pkg/expr	表達式引擎無可重用性，且耦合具體 Alerting 結構
plugins/	不是 runtime plugin、無掃描註冊或獨立生命週期
apps/alert-app/	expression 是 Alert 的通用邏輯，應保留邏輯層中供多 app 使用


⸻

是否需要我幫你轉繪這個 expression pipeline 的 Mermaid 流程圖，或草擬 interface 對應（如 PipelineBuilder）？