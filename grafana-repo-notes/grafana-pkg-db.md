
# pkg/storage 跟 pkg/tsdb
`pkg/storage` 與 `pkg/tsdb` 是 Grafana 中處理「資料存取」的兩個核心模組，但功能與定位**完全不同**：

* * *

✅ `pkg/storage`：**Grafana Unified Storage 架構**
----------------------------------------------

### 📌 主要用途

*   作為 **Grafana 平台型資源管理儲存層**，支援多種資源物件（如 dashboard、secret、library panel）的讀寫與查詢
    
*   提供統一介面支援 legacy SQL、Unified SQL、Parquet、Blob、gRPC storage 等
    

### 📦 子模組結構

*   `unified/`: 次世代儲存核心（支援 gRPC、分片、blob）
    
*   `legacysql/`: 舊有 SQL-based 儲存層（採用 xorm）
    
*   `resourcepb/`: 儲存通訊格式（gRPC proto）
    
*   `resource/`: 核心 CRUD / watcher 實作
    
*   `parquet/`: 支援大資料匯出（批次/壓縮/高效格式）
    

### 📌 功能特性

*   多種儲存格式支援（text、JSON、blob、parquet）
    
*   支援 K8s-style GVK API Resource 存取
    
*   被 apiserver / registry / provisioning 模組大量依賴
    
*   支援「dualwrite」（雙寫 legacy 與 unified）機制grafana-pkg-all-code
    

* * *

✅ `pkg/tsdb`：**Grafana 面板數據查詢（Query Engine）**
---------------------------------------------

### 📌 主要用途

*   用於處理 Dashboard Panel 的 **資料查詢**，例如 Prometheus、InfluxDB、Loki、Graphite 等
    
*   負責轉換 panel 查詢（如 dashboard → query model → datasource plugin）
    

### 📦 子模組結構

*   `tsdb/grafanads/`: Grafana 資料源標準 interface 定義
    
*   `tsdb/prometheus/`, `tsdb/graphite/`, `tsdb/elasticsearch/`：每種資料來源的查詢適配器
    
*   `tsdb/backendplugin/`: 與遠端 plugin 溝通查詢資料的後端協定
    

### 📌 功能特性

*   用於資料查詢流程：Query → Eval → Transform → Frame
    
*   實作 Panel Plugin 所需的資料來源統一查詢介面
    
*   不直接與 Unified Storage 結構整合
    

* * *

🆚 差異總結
-------

| 面向 | `pkg/storage` | `pkg/tsdb` |
| --- | --- | --- |
| 目的 | 儲存 API-style 資源 | 查詢時序/指標資料 |
| 接收來源 | REST API / Registry / Plugin | Dashboard Panel 查詢引擎 |
| 資料型態 | Dashboard、Blob、Secret、KV 物件等 | Prometheus, InfluxDB 等 query result |
| 通訊協定 | gRPC, SQL, Parquet, Blob | Plugin backend, JSON RPC |
| 擴展形式 | 可 plug storage backend（如 unified, legacy） | 可 plug datasource query adapter |

* * *

🧠 適合應用參考
---------

*   若你想實作「模組化資料儲存」架構（如 detectviz support alert config/blob export），**請參考 `pkg/storage`**
    
*   若你開發「資料查詢引擎/面板模組」，如轉換成 `frame` 結構再呈現，**請參考 `pkg/tsdb`**
    

# 遷移到 Detectviz internal/storage/ internal/tsquery/  後怎麼做

根據 Detectviz 的分層原則與你的平台定位，pkg/storage 與 pkg/tsdb 應該如下分類與遷移：

⸻

✅ pkg/storage → 遷移到 internal/storage/

📌 為什麼？
	•	專屬於 Detectviz 自有資源的儲存邏輯（如 alert rule、dataset blob、export parquet）
	•	與 app、registry、event 等模組強耦合，並不通用
	•	屬於應用層的「儲存策略與生命週期」管理

✅ 遷移建議路徑：

internal/storage/
├── unified/       # 類似 grafana unified storage（gRPC blob / parquet）
├── legacysql/     # SQL 模式（如 GORM, xorm）
├── resource/      # 主要 CRUD 與 watcher
├── resourcepb/    # protobuf 定義（若有用 gRPC）

若你希望提供外部模組使用，可將 interface 抽出至 pkg/iface/storage.go

⸻

✅ pkg/tsdb → 建議遷移至 internal/queryengine/ 或 internal/tsquery/

📌 為什麼？
	•	它是「查詢引擎」，專門為面板、報表、alert 引擎提供資料查詢邏輯
	•	通常與 dashboard panel 或 alert rule 模組直接互動
	•	不建議放在 pkg/，因為你不會把 query adapter 提供給其他 repo 使用

✅ 遷移建議路徑：

internal/queryengine/
├── prometheus/
├── influxdb/
├── graphite/
├── adapter.go
├── interface.go     # QueryAdapter interface

或簡化為：

internal/tsquery/

如果你日後要 plugin 化 datasources，也可將 adapter 拆到 plugins/querysource/，搭配 interface 放在 pkg/iface/query.go

⸻

🧭 總結表

模組	功能定位	建議路徑	是否抽象 interface
pkg/storage	App 資源儲存層	internal/storage/	✅ pkg/iface/storage.go（可選）
pkg/tsdb	面板/告警的查詢引擎	internal/queryengine/ 或 internal/tsquery/	✅ pkg/iface/query.go（可選）


⸻

是否需要我幫你生成搬遷後的樹狀結構，或建立對應的 interface 定義草稿？