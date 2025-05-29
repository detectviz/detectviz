# Detectviz

Detectviz 是一套專為工控與異質場域設計的 **智能監控控制平台**，建構於 [Telegraf](https://github.com/influxdata/telegraf) 的穩定 agent 架構之上。

Detectviz 不重新實作收值引擎，而是以 Telegraf 作為資料收集與執行層，專注提供：

- **設備自動掃描與模組調度**
- **多資料來源 plugin 編排管理**
- **異常偵測、告警判斷與自動化回應**
- **AutoML + LLM 報表洞察**
- **多節點調度與可視化中控**

Detectviz 透過 Config 定義節點模組組合、plugin 組裝、資料轉向與行為邏輯，使工控監控從收值到告警行動形成統一、自動化、可管理的閉環。

聚焦在「**效能監控 + 告警判斷 + 分層通知 + 可視化 + 回溯分析**」，模組化、資料驅動、低耦合、高可重用。
> **「將繁雜異質的資料與告警來源，抽象為模組化架構，建立一個低門檻、可維護、具備洞察力的可視化與告警平台。」**
產品定義與設計理念 
你的系統架構（尤其是前端的 `web/` 與後端的 `viot-web/`）**完全對齊「以 config 驅動場景、模組化構建平台」的核心思想**
---

## 1. 產品定位與核心理念

Detectviz 為一平台型產品，具有以下核心能力：

- **插件可編排**：基於 Telegraf plugin 架構，收值 collector 皆以 autodeploy 模組管理，輸出可編排至 Telegraf 應用目錄。
- **節點可分散**：支援多場域、多節點部署，具備「monitor / analyzer / control」分工模式。
- **告警可判斷**：支援閾值、統計與 AI 模型混合檢測，具備告警分級與多通道通報。
- **行為可自動化**：任意告警可對應自動修復腳本、Webhook 或任務清單執行，並可透過 AutoML 模型生成告警摘要。
- **部署可標準化**：透過 `detectviz.conf` 組裝後端模組配置、`pages.yaml` 組裝前端頁面配置，適用多場域多樣化場景複用。
- **平台可觀測**：透過 OTEL 定義格式，配合 Tempo / Loki / Prometheus / Mimir 實現統一觀測。

## 平台產品應具備的標準：

1. 模組化設計
2. 低耦合高可重用
3. 以 interface 驅動業務邏輯
4. 統一配置、熱重載、集中管理
5. 支援 plugin、微服務、動態注入
6. 與 Telegraf/Grafana/Loki/LLM 的整合性強
7. 各模組皆符合 **SRP（單一職責原則）**，未出現過度集中邏輯。
8. 符合事件驅動系統的基本原則
9. 使用 Redis Stream 作為事件匯流排，以支援微服務化與可插拔設計。
10. 可支援 registry 熱註冊與 mock 測試，滿足未來 extensibility 要求
11. 符合 Clean Architecture 原則（分層、依賴倒置）
12. 模組獨立、可插拔、熱部署
    

> 所有 ML / LLM 模型皆封裝為 Python 獨立服務，放置於 `services/` 目錄下，與 Go 核心模組解耦。Detectviz 後端模組以 HTTP 呼叫這些服務，以實現模型模組化、部署彈性與跨語言維護。

## 2. 技術規範

### 依賴服務

1. Keycloak for authentication and authorization
2. Telegraf for collect prometheus metrics
3. InfluxDB for time-series database
4. MySQL for relational database
5. Redis for distributed cache and event bus
6. Tempo for traces
7. Loki for logs
8. Grafana for visualization
9. Ollama for LLM processing

### 分層

| 層級 | 說明 |
| --- | --- |
| `apps/viot-web/` | 應用進入點、統整 API 與前端載入 |
| `internal/` | 核心商業邏輯模組，清楚切分職責 |
| `pkg/` | 通用可重用套件（logger、telemetry、config） |
| `plugins/` | 可擴充動作 plugin，符合 plugin registry 設計模式 |
| `web/` | 前端模組化設計，頁面即模組，元件可重用 |

## 分層原則

- `internal/`: 專屬 Detectviz 的應用邏輯與事件流模組，不可被其他 repo 匯入
- `pkg/`: 可被內部或其他專案重用的通用工具、客戶端、抽象層（無業務耦合）
- `apps/viot-web/`: 驅動整個架構，串接一切模組，提供 REST / HTMX 前端


### 觀測全平台架構圖

```bash
[Go / Python 模組]
   └─ 使用 OpenTelemetry SDK
        ├── Trace 資料
        ├── Log 資料
        └── Metric 資料
            ↓ OTLP (gRPC or HTTP)
[OpenTelemetry Collector]
    ├── Exporter: Tempo (Trace)
    ├── Exporter: Loki  (Log)
    └── Exporter: Mimir (Metric)
```

### 建議導入組合

參考[OpenTelemetry GitHub](https://github.com/orgs/open-telemetry) 的 repo 清單，導入適合的組合。

| 組件 | 建議 Repo / 套件 |
| --- | --- |
| Go SDK | `opentelemetry-go` + `opentelemetry-go-contrib` |
| Python SDK | `opentelemetry-python` + `python-contrib` |
| Exporter/轉發 | `opentelemetry-collector` + `collector-contrib` |
| 整合 Grafana | 搭配 Tempo, Loki, Prometheus 使用 OTLP |

### 後端技術棧

後端統一為 Go 的編排框架，將 ML/LLM 模型或數學引擎封裝為獨立 Python API 微服務，並透過 gRPC 或 HTTP 與主框架通訊。
- 採用 Go Clean Architecture 架構
- 使用 Echo 提供 RESTful API 與 HTMX 前端
- 使用 Viper 讀取與 hot reload 設定檔
- 使用 Gorm / golang-migrate 操作資料庫 MySQL
- 使用 go-redis 操作 Redis 作為事件流
- 使用 Prometheus 收集指標
- 使用 Optlzap 作為日誌管理

### 前端技術棧

- 使用 Tabulator 呈現資料表
- 使用 HTMX 提供輕量互動式前端
- 使用 AdminLTE 提供管理介面
- 以 config 驅動頁面行為與 UI 顯現
- 遵循了「頁面即模組、元件可重用」的理念，每個頁面都是一個獨立模組，可直接載入。

### 模組架構 (internal/ 目錄)

- 每個模組都用 interface 抽象暴露給 orchestrator / router / handler 層，並透過 DI 注入實例。
- pkg 中的模組：都應暴露 interface，集中定義於 `pkg/iface`，例如 EventBus, RedisStore, LLMClient
- internal 專用介面 `internal/<module>/interface.go`，例如 notifier.Service, scheduler.Job

### 模組健康檢查 (internal/healthcheck/ 目錄)

- 每個模組實作 `/api/health`，返回狀態 JSON
- 由 Telegraf inputs.http_response 定時打所有模組 /api/health，回傳 HTTP code
- 使用者輸入的 URL / method / 狀態碼 → 寫入 healthcheck_targets.csv 或 config
- plugins/generator 將這些轉為 Telegraf inputs.http_response 配置
- 自動更新 /etc/telegraf/telegraf.d/health.conf
- 重啟 telegraf 或熱重載

#### 1.	Healthcheck 模組
- healthcheck 模組負責所有模組提供 /api/health
- Telegraf 定時打 API，產生 health_status=1|0 等 metrics

#### 2.	Alert 模組
- alert 模組負責根據 health metrics + service status 訂定規則
- 觸發告警等級（Info / Warning / Critical）
- 可與 notification 通知分流

#### 3.	Automation 模組
- automation 模組應該作為「執行器（executor）」，而不是直接實作所有行動邏輯。
- 所有實際行動應抽象為可擴展的 Plugins，每個 action 各自獨立封裝。
- 接收告警事件（或從 Redis Stream 消費）
- 執行：修復腳本、建立工單、呼叫外部 Webhook、與 LLM 整合進行 log 分析或自動回應建議

#### 3.	Autodeploy 模組
| 模組 | 定位與責任 |
| --- | --- |
| `autodeploy` | 控流程，定義 scan → validate → deploy → snapshot 的主流程 |
| `plugins/scanner` | 定義如何發現設備（SNMP, Modbus, others） |
| `plugins/validator` | 定義如何驗證設備通訊與欄位正確性 |
| `plugins/deployer` | 定義如何生成 Telegraf 設定檔並寫入路徑 |
| `plugins/snapshot` | 定義如何對設定檔做版本備份與比較 |

### 告警事件的「編排者」與「中心樞紐」模組 (internal/alert/ 目錄)

| 模組 | 定位與責任 |
| --- | --- |
| `alert/` | 同時觸發：`notifier.Notify(alert)` + `automation.Trigger(alert)` |
| `notifier/` | 保留：轉為統一通知協調器，呼叫 `plugins/notify/*` |
| `plugins/notify/` | Email / Line / Slack / Webhook 實作 |
| `automation/` | 不直接通知人，而是執行任務（可間接觸發通知任務） |

#### 最佳實踐架構流向

```
[alert]
   │
   ├──▶ [notifier] ──▶ [plugins/notify/...]
   │                      ├── email
   │                      ├── slack
   │                      ├── line
   │                      └── webhook
   │
   └──▶ [automation] ──▶ [plugins/...]
                          ├── shell       # 自動修復腳本
                          ├── webhook     # 呼叫第三方 API
                          ├── ticket      # 建立工單
                          └── llm         # 呼叫 LLM 分析 log / 產生摘要
```

#### 說明

- alert 模組 是流程編排者，觸發 analyzer、reporter。
- analyzer 模組 不實作模型，而是根據規則呼叫 anomaly-python 提供的多種 API。
- reporter 模組 可接收標記後的結果，呼叫 llm-python 生成摘要，再整合為報表。
- anomaly-python 是多種異常偵測模型的 REST API 服務。
- llm-python 是單一 LLM API（摘要、分析、重寫）服務。


### 事件流

- 所有事件皆以 Redis Stream 驅動，具備可追蹤、可重放、可靠交付特性

### OpenTelemetry 整合

- 比照 Grafana 官方預設的整合與 drilldown 鏈接邏輯，在 Grafana 實現「告警帶出觀測數據」，讓使用者能「點選一筆告警 → 自動帶出 trace / log / metrics」
- 可以在 Grafana 點選：
  - 查看 trace（Tempo）
  - 查看 log（Loki）
  - 查看 metrics（Mimir）

### 串接流程總覽圖

```
[使用者登入]
    │
    ▼
Detectviz (串 Keycloak)
    │
    ▼
後端 Cookie / Token 保存 → 植入 iframe URL
    │
    ▼
[iframe: Grafana URL]
    │
    ▼
Nginx 反向代理處理：
- 自動帶上 JWT / Cookie
- 注入 X-Scope-OrgID header 給 Mimir
```

### 整合清單

1. Grafana 設定 OAuth 登入: 使用 Keycloak client
2. iframe 開啟 `allow_embedding`: Grafana + Nginx 配置
3. Nginx 注入 X-Scope-OrgID: 若 Mimir 多租戶有啟用
4. Detectviz 傳遞 org_id 給 iframe: 依據登入者或 token claim 動態帶入
5. iframe 搭配變數切換: 使用 var-xxx URL 參數

## 系統架構

### 目錄結構

```bash
detectviz/
└── apps/
    └── viot-web/                    # 節點應用程式
        ├── main.go                  # 啟動伺服器，初始化 router、middleware、config
        ├── router/
        │   └── router.go            # 定義 Echo 路由（/pages/:name, /components/:name, /api/...）
        │
        ├── handler/                 # 處理請求的邏輯
        │   ├── page_handler.go      # 載入 web/pages/* 組合為 HTML
        │   ├── component_handler.go # 載入元件片段（含 partial 渲染）
        │   └── api_handler.go       # 提供資料來源（讀取 CSV / JSON / log / exec）
        │
        ├── middleware/              # 權限、log、session 等中介層
        │   └── telemetry.go         # otelecho 中介層自動注入 trace_id / span_id
        ├── config/                  # 組態讀取（載入 pages.yaml、env.yaml）
        ├── utils/                   # 公共工具函式
        ├── models/                  # 結構定義（如 device, alert 等）
        ├── services/                # 可呼叫的邏輯單元（如 ExecService, LogService）
        └── templates/               # 若有 templ 模板引擎（可用於 layout 組裝）

├── data/               # 資料儲存目錄
│   ├── exports/        # 報表輸出
│   └── embeddings/     # 向量資料儲存

├── conf/               # 設定檔目錄
│   ├── .env            # 環境變數設定
│   ├── detectviz.conf  # 系統核心設定
│   ├── pages.yaml      # 驅動前端頁面行為與欄位結構
│   ├── otel-collector-config.yaml # 串接 Tempo / Loki / Mimir 的設定檔
│   └── provisioning/   # 配置
│       ├── schemas/    # 資料結構定義
│       ├── migrations/ # 資料庫遷移腳本
│       └── templates/  # 各模組的模板
│           ├── dashboard/     # Grafana 儀表板
│           ├── reporter/      # 報表模板
│           ├── llm/           # 提示詞模板
│           └── telegraf/      # Telegraf 設定檔模板
│               ├── http_response.conf # 檢查設備在線設定檔
│               ├── net_response.conf # 檢查應用模組健康設定檔
│               ├── telegraf.conf    # Telegraf 設定檔
│               └── snmp.conf        # SNMP 設定檔

├── pkg/                # 平台共用程式庫
│   ├── iface/          # 介面定義集中於此，可被內部或其他專案重用
│   ├── env/            # 無狀態的環境變數處理
│   ├── config/         # 讀取 conf/ 搭配 viper 支援 hot reload
│   ├── logger/         # 日誌管理 optlzap 注入 trace_id
│   │   └── otel.go     # 整合 otel 的 logger
│   ├── telemetry/      # otel/metric 指標收集
│   │   ├── otel.go     # 基礎 OpenTelemetry SDK 設定
│   │   └── wrapper.go  # 各模組監控包裝器

│   ├── scheduler/      # 提供背景任務註冊與觸發，支援內部模組使用
│   ├── eventbus/       # 事件匯流排
│   │   ├── interface.go 
│   ├── auth/           # 認證模組，sso 可替換為其他認證模組
│   ├── registry/       # 初始化 Redis、MySQL、Keycloak 等 client
│   │   ├── keycloak/
│   │   ├── redis/
│   │   ├── influxdb/
│   │   └── mysql/
├── services/               # 各種可獨立部署的微服務或外部依賴程式
│   ├── anomaly-python/     # 異常偵測模型 REST API（Python）
│   │   └── main.py
│   └── llm-python/         # LLM 模型 API（Python）
│       └── main.py
│   
├── internal/           # 平台核心功能
│   ├── orchestrator/   # 僅負責模組協調與服務註冊/上下線
│   ├── healthcheck/    # 模組健康檢查
│   ├── rule/           # 規則 CRUD 管理
│   ├── contacts/       # 聯絡人 CRUD 管理
│   ├── notifier/       # 通知模組
│   ├── alert/          # 告警事件的「編排者」與「中心樞紐」模組
│   ├── autodeploy/     # 自動部署器任務協調層，根據掃描偵測對象觸發自動行為
│   ├── automation/     # 自動化任務協調層，根據事件觸發自動行為
│   │   ├── executor.go
│   │   ├── trigger.go
│   │   └── registry.go
│   │
│   ├── analyzer/           # 掛接 API，負責決策與調度
│   │   ├── rule_checker.go # 規則檢查器
│   │   ├── service.go      # 透過 interface 調用 plugins/anomaly 或 plugins/llm
│   │   ├── detector.go     # 整合異常結果、標記異常數據
│   │   └── summarizer.go   # Prompt 上下文編排邏輯
│   │
│   ├── reporter/            # 報表模板與產出邏輯
│   │   ├── generator.go
│   │   └── exporter.go
│   │
│   ├── plugins/            # 實際動作 plugin（可動態增加）
│   │   ├── iface.go        # 明確定義所有 Plugin Interface，支援熱插拔
│   │   ├── registry.go     # 註冊機制（可動態載入 plugin）
│   │   ├── generator.go    # 生成設定檔
│   │   ├── deployer.go     # 部署器
│   │   ├── reload.go       # 重載 Telegraf 服務
│   │   ├── validator.go    # 驗證設定檔
│   │   ├── snapshot.go     # 快照管理
│   │   ├── scanner/        # 掃描偵測對象
│   │   │   ├── snmp.go
│   │   │   └── modbus.go
│   │   ├── shell/          # 執行 shell script
│   │   ├── ticket/         # 建立工單
│   │   ├── webhook/        # 執行 webhook(ITSM / Slack / ServiceNow)
│   │   ├── notify/         # 發送簡訊或 Line
│   │   │   ├── email.go    # 發送 email
│   │   │   ├── slack.go    # 發送 slack
│   │   │   ├── line.go     # 發送 line
│   │   │   └── webhook.go  # 發送 webhook
│   │   ├── anomaly/        # 呼叫 anomaly-python API 進行 SPC, Threshold, ML 分析
│   │   └── llm/            # 呼叫 llm-python API 產生自然語言摘要、log 解釋
│   
├── scripts/            # 工具腳本
├── test/               # 測試程式
└── web/                # 前端介面
    ├── pages/                    # 每個頁面子資料夾為獨立模組，支援 config.yaml 組態
    │   ├── login/                # 登入頁
    │   ├── grid-status/          # 顯示裝置狀態
    │   ├── record-status/        # 告警 / 任務 / 稽核記錄
    │   ├── table-config/         # Tabulator 編輯器
    │   ├── csv-editor/           # CSV 表格 CRUD 編輯器
    │   ├── log-viewer/           # 系統日誌
    │   ├── file-editor/          # 設定檔 CRUD (yaml、json、txt)
    │   └── task-runner/          # 執行任務（Exec 模組）
    │
    ├── components/               # 可重用元件，抽離出小模組
    │
    ├── layouts/                  # 共用佈局（如 sidebar、navbar、footer）
    │   └── base.html             # 預設 layout，包含 sidebar、navbar、footer
    │
    ├── static/                   # 靜態資源：圖片、icon、JS libs
    │   ├── img/
    │   ├── libs/
    │   │   ├── tabulator.min.js
    │   │   ├── htmx.min.js
    │   │   └── bootstrap.css
    │   └── fonts/
    │
    └── index.html               # 可選 demo preview 或 redirect 頁

```


# pkg 共用程式庫的 interface

所有定義於 `pkg/iface/` 的介面，目的是讓各個 internal 模組可以透過依賴注入方式使用共用邏輯元件，例如：Redis client、Logger、設定載入器、OpenTelemetry 等，達到模組解耦與單元測試替換的目的。

根據系統架構，這些 interface 主要會在以下幾個目錄中實作：

1. **核心基礎設施 (pkg/)**
   - `env.go` → `pkg/env/`
   - `config.go` → `pkg/config/`
   - `logger.go` → `pkg/logger/`
   - `telemetry.go` → `pkg/telemetry/`
   - `error.go` → `pkg/error/`

2. **資料存取層 (pkg/registry/)**
   - `redis.go` → `pkg/registry/redis/`
   - `mysql.go` → `pkg/registry/mysql/`
   - `influxdb.go` → `pkg/registry/influxdb/`
   - `keycloak.go` → `pkg/registry/keycloak/`

3. **事件與檔案系統 (pkg/)**
   - `eventbus.go` → `pkg/eventbus/`
   - `filestore.go` → `pkg/filestore/`
   - `datasource.go` → `pkg/datasource/`

4. **安全與快取 (pkg/)**
   - `auth.go` → `pkg/auth/`
   - `cache.go` → `pkg/cache/`
   - `secret.go` → `pkg/secret/`

5. **排程系統 (pkg/)**
   - `scheduler.go` → `pkg/scheduler/`

### 詳細的對應表：

| Interface | 實作目錄 | 主要依賴 | 使用場景 |
|-----------|----------|----------|----------|
| `EnvReader` | `pkg/env/` | 無 | 環境變數讀取 |
| `ConfigManager` | `pkg/config/` | `EnvReader` | 設定檔管理 |
| `Logger` | `pkg/logger/` | `Telemetry` | 日誌記錄 |
| `Telemetry` | `pkg/telemetry/` | 無 | 追蹤與指標 |
| `RedisClient` | `pkg/registry/redis/` | 無 | Redis 操作 |
| `SQLClient` | `pkg/registry/mysql/` | 無 | MySQL 操作 |
| `InfluxWriter` | `pkg/registry/influxdb/` | 無 | InfluxDB 寫入 |
| `AuthProvider` | `pkg/registry/keycloak/` | 無 | 認證授權 |
| `EventBus` | `pkg/eventbus/` | `RedisClient` | 事件處理 |
| `FileStore` | `pkg/filestore/` | 無 | 檔案操作 |
| `DataSource` | `pkg/datasource/` | `SQLClient` | 資料存取 |
| `CacheManager` | `pkg/cache/` | `RedisClient` | 快取管理 |
| `JobScheduler` | `pkg/scheduler/` | `Logger` | 任務排程 |
| `SecretManager` | `pkg/secret/` | `FileStore` | 密鑰管理 |
| `ErrorHandler` | `pkg/error/` | `Logger` | 錯誤處理 |

實作建議：

1. **目錄結構**：
   - 每個 interface 都應該有自己的目錄
   - 目錄下應包含 `interface.go`、`impl.go` 和 `mock.go`
   - 可以加入 `options.go` 用於配置選項

2. **依賴注入**：
   - 所有實作都應該支援依賴注入
   - 使用 factory 模式創建實例
   - 提供預設實作和自定義選項

3. **測試支援**：
   - 每個目錄都應該包含測試檔案
   - 提供 mock 實作用於測試
   - 包含整合測試和單元測試

4. **文檔要求**：
   - 每個目錄都應該有 README.md
   - 包含使用範例和配置說明
   - 說明依賴關係和使用限制

5. **錯誤處理**：
   - 統一使用 `ErrorHandler`
   - 提供詳細的錯誤資訊
   - 支援錯誤追蹤和重試

這樣的組織方式可以確保：
- 清晰的職責分離
- 良好的可測試性
- 靈活的擴展性
- 統一的錯誤處理
- 完整的文檔支援

需要我為任何特定 interface 提供更詳細的實作建議嗎？


## 介面檔案清單與說明

### 1. `env.go`
```go
// EnvReader：環境變數讀取器介面，統一管理 .env 或系統環境變數
type EnvReader interface {
	Get(key string) string
	GetOrDefault(key, fallback string) string
}
```

### 2. `config.go`
```go
// ConfigManager：整合環境變數與設定檔管理
// 功能：
// - 統一管理環境變數和設定檔
// - 支援設定檔熱重載
// - 提供型別安全的配置讀取
// - 支援配置驗證
type ConfigManager interface {
	// 設定檔操作
	// Load：載入指定路徑的設定檔
	// GetConfig：取得指定 key 的設定值
	// Watch：註冊設定變更回調函數
	Load(path string) error
	GetConfig(key string) any
	Watch(callback func(key string, value any))
	
	// 驗證與轉換
	// Validate：驗證所有設定值是否符合規則
	// GetString/GetInt/GetBool：型別安全的配置讀取
	Validate() error
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
}
```

### 3. `logger.go`
```go
// Logger：通用日誌介面，支援多級別與上下文
// 功能：
// - 支援多級別日誌（Debug/Info/Warn/Error）
// - 支援結構化日誌
// - 支援上下文追蹤
// - 支援錯誤鏈追蹤
type Logger interface {
	// 基本日誌方法
	// Debug：用於開發時期的詳細資訊
	// Info：用於一般資訊記錄
	// Warn：用於警告訊息
	// Error：用於錯誤訊息
	Debug(msg string, fields ...any)
	Info(msg string, fields ...any)
	Warn(msg string, fields ...any)
	Error(msg string, fields ...any)
	
	// 上下文擴充
	// With：新增固定欄位到日誌
	// WithContext：從 context 提取追蹤資訊
	// WithError：記錄錯誤資訊
	With(fields ...any) Logger
	WithContext(ctx context.Context) Logger
	WithError(err error) Logger
}
```

### 4. `telemetry.go`
```go
// Telemetry：整合追蹤與指標收集
// 功能：
// - 分散式追蹤
// - 指標收集
// - 取樣率控制
// - 自定義標籤
type Telemetry interface {
	// 追蹤相關
	// StartSpan：開始新的追蹤區間
	// GenerateTraceID：生成追蹤 ID
	// SetSamplingRate：設定取樣率
	StartSpan(ctx context.Context, name string) (context.Context, Span)
	GenerateTraceID() string
	SetSamplingRate(rate float64)
	
	// 指標相關
	// RecordCounter：記錄計數器
	// RecordGauge：記錄儀表值
	// RecordHistogram：記錄直方圖
	RecordCounter(name string, value int64, labels map[string]string)
	RecordGauge(name string, value float64, labels map[string]string)
	RecordHistogram(name string, value float64, labels map[string]string)
}

// Span：抽象化的 OTel span
// 功能：
// - 追蹤區間管理
// - 事件記錄
// - 屬性設定
type Span interface {
	// End：結束追蹤區間
	// AddEvent：新增事件
	// SetAttributes：設定屬性
	End()
	AddEvent(name string)
	SetAttributes(attrs map[string]any)
}
```

### 5. `redis.go`
```go
// RedisClient：Redis 操作介面
// 功能：
// - 鍵值操作
// - 發布訂閱
// - 集合操作
// - 事務支援
type RedisClient interface {
	// 基本操作
	// Set：設定鍵值，支援 TTL
	// Get：取得鍵值
	// Delete：刪除鍵值
	Set(key string, value any, ttl time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
	
	// 發布訂閱
	// Publish：發布訊息
	// Subscribe：訂閱訊息
	Publish(channel string, payload any) error
	Subscribe(channel string, handler func(payload string)) error
	
	// 集合操作
	// SAdd：新增集合成員
	// SMembers：取得集合成員
	SAdd(key string, members ...string) error
	SMembers(key string) ([]string, error)
	
	// 事務支援
	// Begin：開始事務
	Begin() (Transaction, error)
}

// Transaction：Redis 事務介面
// 功能：
// - 事務提交
// - 事務回滾
type Transaction interface {
	// Exec：提交事務
	// Rollback：回滾事務
	Exec() error
	Rollback() error
}
```

### 6. `influxdb.go`
```go
// InfluxWriter：InfluxDB 寫入介面
type InfluxWriter interface {
	WritePoint(measurement string, tags map[string]string, fields map[string]any, ts time.Time) error
}
```

### 7. `keycloak.go`
```go
// AuthProvider：認證提供者介面（如：Keycloak）
type AuthProvider interface {
	Login(username, password string) (token string, err error)
	ValidateToken(token string) (bool, error)
}
```


### 8. `mysql.go`
```go
// SQLClient：MySQL 操作封裝介面
// 功能：
// - SQL 查詢執行
// - 事務管理
// - 預處理語句
// - 連接池管理
type SQLClient interface {
	// 基本操作
	// Query：執行查詢
	// Exec：執行更新
	Query(query string, args ...any) ([]map[string]any, error)
	Exec(query string, args ...any) error
	
	// 事務支援
	// Begin：開始事務
	Begin() (Transaction, error)
	
	// 預處理語句
	// Prepare：準備預處理語句
	Prepare(query string) (Stmt, error)
	
	// 連接池管理
	// SetMaxOpenConns：設定最大連接數
	// SetMaxIdleConns：設定最大閒置連接數
	// SetConnMaxLifetime：設定連接生命週期
	SetMaxOpenConns(n int)
	SetMaxIdleConns(n int)
	SetConnMaxLifetime(d time.Duration)
}

// Transaction：SQL 事務介面
// 功能：
// - 事務提交
// - 事務回滾
type Transaction interface {
	// Commit：提交事務
	// Rollback：回滾事務
	Commit() error
	Rollback() error
}

// Stmt：預處理語句介面
// 功能：
// - 參數化查詢
// - 參數化更新
// - 資源釋放
type Stmt interface {
	// Query：執行參數化查詢
	// Exec：執行參數化更新
	// Close：釋放資源
	Query(args ...any) ([]map[string]any, error)
	Exec(args ...any) error
	Close() error
}
```

---

### 9. `eventbus.go`
```go
// EventBus：事件匯流排介面
// 功能：
// - 事件發布訂閱
// - 重試機制
// - 死信佇列
// - 訊息過期
type EventBus interface {
	// 基本操作
	// Publish：發布事件
	// Subscribe：訂閱事件
	Publish(stream string, payload map[string]any) error
	Subscribe(stream string, handler func(entry map[string]any)) error
	
	// 進階功能
	// PublishWithRetry：支援重試的發布
	// SubscribeWithDeadLetter：支援死信佇列的訂閱
	// SetMessageTTL：設定訊息過期時間
	PublishWithRetry(stream string, payload map[string]any, maxRetries int) error
	SubscribeWithDeadLetter(stream string, handler func(entry map[string]any), dlq string) error
	SetMessageTTL(stream string, ttl time.Duration) error
}
```

---

### 10. `filestore.go`
```go
// FileStore：檔案操作介面
// 功能：
// - 檔案讀寫
// - 檔案鎖定
// - 版本控制
// - 元資料管理
type FileStore interface {
	// 基本操作
	// Read：讀取檔案
	// Write：寫入檔案
	// Delete：刪除檔案
	Read(path string) ([]byte, error)
	Write(path string, content []byte) error
	Delete(path string) error
	
	// 進階功能
	// Lock：檔案鎖定
	// Unlock：檔案解鎖
	// GetMetadata：取得檔案元資料
	// ListVersions：列出檔案版本
	Lock(path string) error
	Unlock(path string) error
	GetMetadata(path string) (FileMetadata, error)
	ListVersions(path string) ([]FileVersion, error)
}

// FileMetadata：檔案元資料
// 功能：記錄檔案的基本資訊
type FileMetadata struct {
	Size    int64     // 檔案大小
	ModTime time.Time // 修改時間
	Hash    string    // 檔案雜湊值
}

// FileVersion：檔案版本
// 功能：記錄檔案的版本資訊
type FileVersion struct {
	Version   string    // 版本號
	Timestamp time.Time // 版本時間
	Hash      string    // 版本雜湊值
}
```

---

### 11. `datasource.go`
```go
// DataSource：通用資料來源介面
// 功能：
// - CRUD 操作
// - 分頁查詢
// - 排序功能
// - 批量操作
// - 資料驗證
type DataSource interface {
	// 基本操作
	// List：列出資料
	// Create：新增資料
	// Update：更新資料
	// Delete：刪除資料
	List(filter map[string]any) ([]map[string]any, error)
	Create(record map[string]any) error
	Update(id string, record map[string]any) error
	Delete(id string) error
	
	// 進階功能
	// ListWithPagination：分頁查詢
	// ListWithSort：排序查詢
	// BatchCreate：批量新增
	// BatchUpdate：批量更新
	// Validate：資料驗證
	ListWithPagination(filter map[string]any, page, size int) ([]map[string]any, int64, error)
	ListWithSort(filter map[string]any, sortBy string, ascending bool) ([]map[string]any, error)
	BatchCreate(records []map[string]any) error
	BatchUpdate(records []map[string]any) error
	Validate(record map[string]any) error
}
```

### 12. `auth.go`
```go
// AuthProvider：認證提供者介面
// 功能：
// - 使用者認證
// - 令牌管理
// - 權限檢查
// - 角色管理
type AuthProvider interface {
	// 基本操作
	// Login：使用者登入
	// ValidateToken：驗證令牌
	// Logout：使用者登出
	Login(username, password string) (token string, err error)
	ValidateToken(token string) (bool, error)
	Logout(token string) error
	
	// 進階功能
	// RefreshToken：更新令牌
	// GetUserInfo：取得使用者資訊
	// CheckPermission：檢查權限
	// ListRoles：列出角色
	RefreshToken(token string) (string, error)
	GetUserInfo(token string) (UserInfo, error)
	CheckPermission(token string, resource string, action string) (bool, error)
	ListRoles() ([]Role, error)
}

// UserInfo：使用者資訊
// 功能：記錄使用者的基本資訊
type UserInfo struct {
	ID       string   // 使用者 ID
	Username string   // 使用者名稱
	Email    string   // 電子郵件
	Roles    []string // 角色列表
}

// Role：角色資訊
// 功能：記錄角色的權限資訊
type Role struct {
	ID          string       // 角色 ID
	Name        string       // 角色名稱
	Permissions []Permission // 權限列表
}
```

### 13. `cache.go`
```go
// CacheManager：快取管理介面
// 功能：
// - 快取存取
// - TTL 管理
// - 批量操作
// - 快取回填
type CacheManager interface {
	// 基本操作
	// Get：取得快取
	// Set：設定快取
	// Delete：刪除快取
	// Clear：清除所有快取
	Get(key string) (any, error)
	Set(key string, value any, ttl time.Duration) error
	Delete(key string) error
	Clear() error
	
	// 進階功能
	// GetOrSet：取得或設定快取
	// GetMulti：批量取得快取
	// SetMulti：批量設定快取
	GetOrSet(key string, ttl time.Duration, getter func() (any, error)) (any, error)
	GetMulti(keys []string) (map[string]any, error)
	SetMulti(items map[string]any, ttl time.Duration) error
}
```

### 14. `scheduler.go`
```go
// JobScheduler：任務排程介面
// 功能：
// - 任務排程
// - 任務管理
// - 執行記錄
// - 重試機制
type JobScheduler interface {
	// 基本操作
	// Schedule：排程任務
	// Cancel：取消任務
	// ListJobs：列出任務
	// GetJobStatus：取得任務狀態
	Schedule(job Job) error
	Cancel(jobID string) error
	ListJobs() ([]Job, error)
	GetJobStatus(jobID string) (JobStatus, error)
	
	// 進階功能
	// PauseJob：暫停任務
	// ResumeJob：恢復任務
	// UpdateJob：更新任務
	// GetJobHistory：取得任務歷史
	PauseJob(jobID string) error
	ResumeJob(jobID string) error
	UpdateJob(job Job) error
	GetJobHistory(jobID string) ([]JobExecution, error)
}

// Job：任務定義
// 功能：定義任務的基本資訊
type Job struct {
	ID          string       // 任務 ID
	Name        string       // 任務名稱
	Cron        string       // Cron 表達式
	Handler     func() error // 處理函數
	RetryPolicy RetryPolicy  // 重試策略
}

// JobStatus：任務狀態
// 功能：記錄任務的執行狀態
type JobStatus struct {
	Status     string    // 狀態
	LastRun    time.Time // 上次執行時間
	NextRun    time.Time // 下次執行時間
	Error      error     // 錯誤資訊
	RetryCount int       // 重試次數
}

// JobExecution：任務執行記錄
// 功能：記錄任務的執行歷史
type JobExecution struct {
	StartTime time.Time // 開始時間
	EndTime   time.Time // 結束時間
	Status    string    // 狀態
	Error     error     // 錯誤資訊
	Output    string    // 輸出資訊
}
```

### 15. `secret.go`
```go
// SecretManager：密鑰管理介面
// 功能：
// - 密鑰存取
// - 版本控制
// - 密鑰輪換
// - 元資料管理
type SecretManager interface {
	// 基本操作
	// GetSecret：取得密鑰
	// SetSecret：設定密鑰
	// DeleteSecret：刪除密鑰
	// ListSecrets：列出密鑰
	GetSecret(key string) (string, error)
	SetSecret(key string, value string) error
	DeleteSecret(key string) error
	ListSecrets() ([]string, error)
	
	// 進階功能
	// GetSecretWithVersion：取得指定版本的密鑰
	// RotateSecret：輪換密鑰
	// GetSecretMetadata：取得密鑰元資料
	GetSecretWithVersion(key string, version string) (string, error)
	RotateSecret(key string) error
	GetSecretMetadata(key string) (SecretMetadata, error)
}

// SecretMetadata：密鑰元資料
// 功能：記錄密鑰的基本資訊
type SecretMetadata struct {
	CreatedAt  time.Time // 建立時間
	UpdatedAt  time.Time // 更新時間
	Version    string    // 版本號
	ExpiresAt  time.Time // 過期時間
}
```

### 16. `error.go`
```go
// ErrorHandler：錯誤處理介面
// 功能：
// - 錯誤處理
// - 錯誤分類
// - 錯誤追蹤
// - 錯誤上下文
type ErrorHandler interface {
	// 基本操作
	// HandleError：處理錯誤
	// IsRetryable：判斷是否可重試
	// GetErrorCode：取得錯誤代碼
	// GetErrorDetails：取得錯誤詳情
	HandleError(err error) error
	IsRetryable(err error) bool
	GetErrorCode(err error) string
	GetErrorDetails(err error) map[string]any
	
	// 進階功能
	// WrapError：包裝錯誤
	// UnwrapError：解包錯誤
	// GetErrorStack：取得錯誤堆疊
	// GetErrorContext：取得錯誤上下文
	WrapError(err error, message string) error
	UnwrapError(err error) error
	GetErrorStack(err error) []string
	GetErrorContext(err error) map[string]any
}
```

---

## 使用建議

1. **依賴注入**：所有模組應透過依賴注入使用這些介面，避免直接依賴具體實作
2. **錯誤處理**：使用 `ErrorHandler` 統一處理錯誤，確保錯誤資訊完整
3. **配置管理**：使用 `ConfigManager` 統一管理環境變數和設定檔
4. **日誌追蹤**：使用 `Logger` 和 `Telemetry` 確保系統可觀測性
5. **資料存取**：使用 `DataSource` 統一資料存取介面，支援多種資料來源
6. **事件處理**：使用 `EventBus` 處理模組間通訊，確保事件可靠傳遞
7. **任務排程**：使用 `JobScheduler` 管理背景任務，支援重試和監控
8. **密鑰管理**：使用 `SecretManager` 安全地管理敏感資訊
9. **快取管理**：使用 `CacheManager` 提升系統效能
10. **認證授權**：使用 `AuthProvider` 統一管理使用者認證和權限

# internal 模組的 interface

所有定義於 `internal/` 的介面，目的是讓各個模組之間可以透過依賴注入方式進行解耦，並支援單元測試替換。這些介面主要用於模組間的內部通訊，與 `pkg/iface` 中的共用介面不同。

## 介面檔案清單與說明

### 1. `analyzer/interface.go`
```go
// AnalyzerPlugin：異常分析插件介面
// 功能：
// - 異常檢測
// - 模型訓練
// - 模型資訊管理
// - 配置更新
type AnalyzerPlugin interface {
	// 基本操作
	// Detect：執行異常檢測
	Detect(metricSet any) (AnomalyResult, error)
	
	// 進階功能
	// Train：訓練模型
	// GetModelInfo：取得模型資訊
	// UpdateConfig：更新配置
	Train(data []any) error
	GetModelInfo() ModelInfo
	UpdateConfig(config map[string]any) error
}

// ModelInfo：模型資訊
// 功能：記錄模型的基本資訊和效能指標
type ModelInfo struct {
	Name        string             // 模型名稱
	Version     string             // 模型版本
	LastUpdated time.Time          // 最後更新時間
	Metrics     map[string]float64 // 效能指標
}

// AnomalyResult：異常分析結果
// 功能：記錄異常檢測的結果和詳情
type AnomalyResult struct {
	IsAnomaly bool              // 是否為異常
	Score     float64           // 異常分數
	Reason    string            // 異常原因
	Details   map[string]any    // 詳細資訊
}
```

### 2. `llm/interface.go`
```go
// LLMPlugin：自然語言處理插件介面
// 功能：
// - 文字摘要
// - 文本分析
// - 文本生成
// - 向量嵌入
type LLMPlugin interface {
	// 基本操作
	// Summarize：生成文字摘要
	Summarize(input LLMInput) (string, error)
	
	// 進階功能
	// Analyze：分析文本內容
	// Generate：生成文本內容
	// GetEmbedding：取得文本向量
	Analyze(input LLMInput) (AnalysisResult, error)
	Generate(input LLMInput) (string, error)
	GetEmbedding(text string) ([]float64, error)
}

// LLMInput：LLM 輸入資料結構
// 功能：定義 LLM 處理所需的輸入參數
type LLMInput struct {
	Context    string         // 上下文資訊
	Prompt     string         // 提示詞
	Parameters map[string]any // 模型參數
}

// AnalysisResult：分析結果
// 功能：記錄文本分析的結果
type AnalysisResult struct {
	Summary    string   // 摘要
	Keywords   []string // 關鍵詞
	Sentiment  float64  // 情感分數
	Categories []string // 分類結果
}
```

### 3. `notifier/interface.go`
```go
// NotifierPlugin：通知插件介面
// 功能：
// - 事件通知
// - 模板通知
// - 重試機制
// - 狀態追蹤
type NotifierPlugin interface {
	// 基本操作
	// Notify：發送通知
	Notify(event AlertEvent) error
	
	// 進階功能
	// NotifyWithTemplate：使用模板發送通知
	// NotifyWithRetry：支援重試的通知
	// GetNotificationStatus：取得通知狀態
	NotifyWithTemplate(event AlertEvent, template string) error
	NotifyWithRetry(event AlertEvent, maxRetries int) error
	GetNotificationStatus(eventID string) (NotificationStatus, error)
}

// NotificationStatus：通知狀態
// 功能：記錄通知的發送狀態
type NotificationStatus struct {
	EventID    string    // 事件 ID
	Status     string    // 通知狀態
	SentAt     time.Time // 發送時間
	RetryCount int       // 重試次數
	Error      error     // 錯誤資訊
}
```

### 4. `automation/interface.go`
```go
// AutomationPlugin：自動化執行動作的介面
// 功能：
// - 動作執行
// - 超時控制
// - 重試機制
// - 執行記錄
type AutomationPlugin interface {
	// 基本操作
	// Execute：執行自動化動作
	Execute(action AutomationAction) (AutomationResult, error)
	
	// 進階功能
	// ExecuteWithTimeout：支援超時控制的執行
	// ExecuteWithRetry：支援重試的執行
	// GetExecutionHistory：取得執行歷史
	ExecuteWithTimeout(action AutomationAction, timeout time.Duration) (AutomationResult, error)
	ExecuteWithRetry(action AutomationAction, maxRetries int) (AutomationResult, error)
	GetExecutionHistory(actionID string) ([]ExecutionRecord, error)
}

// AutomationAction：自動化動作定義
// 功能：定義自動化動作的參數
type AutomationAction struct {
	ID          string            // 動作 ID
	Name        string            // 動作名稱
	Type        string            // 動作類型
	Parameters  map[string]string // 動作參數
	RetryPolicy RetryPolicy       // 重試策略
}

// AutomationResult：自動化執行結果
// 功能：記錄自動化執行的結果
type AutomationResult struct {
	ActionID  string    // 動作 ID
	Success   bool      // 是否成功
	Output    string    // 輸出結果
	StartTime time.Time // 開始時間
	EndTime   time.Time // 結束時間
	Error     error     // 錯誤資訊
}

// ExecutionRecord：執行記錄
// 功能：記錄執行的詳細資訊
type ExecutionRecord struct {
	ActionID  string    // 動作 ID
	Status    string    // 執行狀態
	StartTime time.Time // 開始時間
	EndTime   time.Time // 結束時間
	Output    string    // 輸出結果
	Error     error     // 錯誤資訊
}
```

### 5. `scheduler/interface.go`
```go
// SchedulerService：排程服務介面
// 功能：
// - 任務排程
// - 任務管理
// - 狀態追蹤
// - 歷史記錄
type SchedulerService interface {
	// 基本操作
	// ScheduleJob：排程任務
	// CancelJob：取消任務
	// GetJobStatus：取得任務狀態
	ScheduleJob(job JobDefinition) error
	CancelJob(jobID string) error
	GetJobStatus(jobID string) (JobStatus, error)
	
	// 進階功能
	// PauseJob：暫停任務
	// ResumeJob：恢復任務
	// UpdateJobSchedule：更新任務排程
	// GetJobHistory：取得任務歷史
	PauseJob(jobID string) error
	ResumeJob(jobID string) error
	UpdateJobSchedule(jobID string, schedule string) error
	GetJobHistory(jobID string) ([]JobExecution, error)
}

// JobDefinition：任務定義
// 功能：定義任務的基本資訊
type JobDefinition struct {
	ID          string        // 任務 ID
	Name        string        // 任務名稱
	Schedule    string        // 排程設定
	Handler     func() error  // 處理函數
	RetryPolicy RetryPolicy   // 重試策略
	Timeout     time.Duration // 超時設定
}

// JobStatus：任務狀態
// 功能：記錄任務的執行狀態
type JobStatus struct {
	ID         string    // 任務 ID
	Status     string    // 執行狀態
	LastRun    time.Time // 上次執行時間
	NextRun    time.Time // 下次執行時間
	Error      error     // 錯誤資訊
	RetryCount int       // 重試次數
}
```

### 6. `healthcheck/interface.go`
```go
// HealthChecker：健康檢查介面
// 功能：
// - 健康狀態檢查
// - 檢查項目管理
// - 詳細狀態查詢
type HealthChecker interface {
	// 基本操作
	// CheckHealth：執行健康檢查
	CheckHealth() (HealthStatus, error)
	
	// 進階功能
	// RegisterCheck：註冊檢查項目
	// UnregisterCheck：取消註冊檢查項目
	// GetDetailedStatus：取得詳細狀態
	RegisterCheck(name string, check func() error)
	UnregisterCheck(name string)
	GetDetailedStatus() map[string]ComponentStatus
}

// HealthStatus：健康狀態
// 功能：記錄系統的健康狀態
type HealthStatus struct {
	Status    string                    // 整體狀態
	Timestamp time.Time                 // 檢查時間
	Details   map[string]ComponentStatus // 元件狀態
}

// ComponentStatus：元件狀態
// 功能：記錄個別元件的狀態
type ComponentStatus struct {
	Status    string    // 元件狀態
	Message   string    // 狀態訊息
	LastCheck time.Time // 最後檢查時間
}
```

### 7. `rule/interface.go`
```go
// RuleManager：規則管理介面
// 功能：
// - 規則 CRUD
// - 規則啟用/停用
// - 規則驗證
// - 規則過濾
type RuleManager interface {
	// 基本操作
	// CreateRule：建立規則
	// UpdateRule：更新規則
	// DeleteRule：刪除規則
	// GetRule：取得規則
	CreateRule(rule Rule) error
	UpdateRule(rule Rule) error
	DeleteRule(ruleID string) error
	GetRule(ruleID string) (Rule, error)
	
	// 進階功能
	// ListRules：列出規則
	// EnableRule：啟用規則
	// DisableRule：停用規則
	// ValidateRule：驗證規則
	ListRules(filter RuleFilter) ([]Rule, error)
	EnableRule(ruleID string) error
	DisableRule(ruleID string) error
	ValidateRule(rule Rule) error
}

// Rule：規則定義
// 功能：定義規則的基本資訊
type Rule struct {
	ID          string   // 規則 ID
	Name        string   // 規則名稱
	Description string   // 規則描述
	Condition   string   // 條件表達式
	Actions     []Action // 動作列表
	Enabled     bool     // 是否啟用
	Priority    int      // 優先級
}

// RuleFilter：規則過濾條件
// 功能：定義規則查詢的過濾條件
type RuleFilter struct {
	Enabled  *bool    // 啟用狀態
	Priority *int     // 優先級
	Tags     []string // 標籤列表
}
```

### 8. `contacts/interface.go`
```go
// ContactManager：聯絡人管理介面
// 功能：
// - 聯絡人 CRUD
// - 群組管理
// - 通知管道管理
type ContactManager interface {
	// 基本操作
	// CreateContact：建立聯絡人
	// UpdateContact：更新聯絡人
	// DeleteContact：刪除聯絡人
	// GetContact：取得聯絡人
	CreateContact(contact Contact) error
	UpdateContact(contact Contact) error
	DeleteContact(contactID string) error
	GetContact(contactID string) (Contact, error)
	
	// 進階功能
	// ListContacts：列出聯絡人
	// AddToGroup：加入群組
	// RemoveFromGroup：從群組移除
	ListContacts(filter ContactFilter) ([]Contact, error)
	AddToGroup(contactID string, groupID string) error
	RemoveFromGroup(contactID string, groupID string) error
}

// Contact：聯絡人定義
// 功能：定義聯絡人的基本資訊
type Contact struct {
	ID        string              // 聯絡人 ID
	Name      string              // 聯絡人名稱
	Email     string              // 電子郵件
	Phone     string              // 電話號碼
	Groups    []string            // 群組列表
	Channels  []NotificationChannel // 通知管道
}

// ContactFilter：聯絡人過濾條件
// 功能：定義聯絡人查詢的過濾條件
type ContactFilter struct {
	Groups    []string // 群組列表
	Channels  []string // 通知管道
	Active    *bool    // 是否啟用
}
```

---

## 使用建議

1. **模組解耦**：使用介面定義模組間的依賴關係，避免直接依賴具體實作
2. **單元測試**：透過介面可以輕鬆替換實作，方便進行單元測試
3. **錯誤處理**：統一使用 error 回傳錯誤，並提供詳細的錯誤資訊
4. **狀態追蹤**：所有操作都應提供狀態查詢功能，方便監控和除錯
5. **擴展性**：介面設計應考慮未來擴展需求，預留擴充點
6. **文檔完整**：每個介面都應有完整的註解說明，包含使用範例
7. **版本控制**：介面變更應考慮向後相容性，必要時提供遷移方案
8. **效能考量**：介面設計應考慮效能影響，避免不必要的抽象層
9. **安全性**：涉及敏感操作的介面應提供適當的安全控制
10. **可觀測性**：所有操作都應提供足夠的日誌和指標資訊

