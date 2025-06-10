## 目錄架構

## 1. `/internal/bootstrap` ：初始化與依賴注入
- 建立與初始化 registry、logger、config、scheduler 等元件
- 呼叫 registry.RegisterXxx(...) 來注入元件
- 不持有實作，只做 wiring

/internal/bootstrap/          
└── init.go     

## 2. `/internal/registry`：對應平台元件註冊中心（內部 handler）
- 不對外提供 HTTP API，用來「註冊與調用內部模組實作」，例如：

/internal/registry
├── notifier/
│   └── registry.go         # 註冊多個 notifier（webhook/slack/email）
├── scheduler/
│   └── registry.go         # 註冊 cron/workerpool 等
├── alert/
│   └── evaluator_registry.go
└── registry.go             # 統一導出

## 3. `/internal/modules` 模組啟動與生命週期控制：
- 對應 Grafana `/pkg/modules`
- 功能涵蓋：
	- RegisterModule(name, initFn)
	- RunAll() 模組依序啟動
	- 注入 logger, config, tracer 作為 dependency
	- 若某模組掛掉，自動關閉其他模組

/internal/modules/
├── engine.go             // 模組引擎與執行框架（對應 modules.Engine）
├── dependencies.go       // 定義模組名稱與依賴（如 Scheduler -> Logger）
├── registry.go           // 註冊模組（類似 RegisterModule）
├── listener.go           // 健康監控與錯誤處理
└── runner.go             // 啟動所有模組


## 4. `/internal/server`
- 對應 Grafana `/pkg/server`
- 預留很好，這一層通常對應主執行 entrypoint（如 main.go 或 server.New()）。
- 預計可整合 modules.RunAll() + HTTP Server 啟動 + graceful shutdown + PID 寫入 等邏輯。

/internal/server/
├── server.go
├── runner.go 
└── instrumentation.go


## 5. `/pkg/registry`：對應 Grafana-style 的 API Resource Registry
- 對應 Grafana `/pkg/registry`，支援：
	- GVK 註冊
	- YAML/CUE schema 驗證
	- RESTful CRUD + RBAC
	- 可搭配 OpenAPI 文件生成

```bash
/pkg/registry              
├── apis/
│   ├── host/
│   │   └── register.go
│   └── datasource/
│       └── register.go
├── schemas/                # GVK 與 CRD 對應 schema
├── kinds/                  # cue schema 解析與驗證
└── registry.go             # 資源註冊與 CRUD 綁定邏輯
```


## 6. `/internal/plugins`：對應 plugin 管理系統
- 支援 datasource plugin、notifier plugin、alert plugin 等功能模組的掛載與啟動。
- 對應 Grafana `/pkg/plugins`

/internal/plugins           (對應 Grafana /pkg/plugins)
├── manager/                # plugin 管理主體
│   ├── registry.go         # plugin 註冊中心（pluginRegistry）
│   ├── loader.go           # 掃描與加載 plugin
│   ├── process.go          # 啟動 plugin 後端
│   └── lifecycle.go        # 初始化、關閉 plugin
├── cdn/                    # CDN 插件資源載入（可選）
├── metadata/               # plugin.json、簽章驗證
└── plugin.go               # plugin interface 定義


# 參考 ========================================================


# pkg/modules:
`pkg/modules` 是 Grafana 中專門用來管理「模組化啟動流程」的元件，採用類似 `systemd` 或 `microservice graph` 的方式管理各個 background service 的初始化、依賴與生命週期。這個模組實作了「模組註冊 → 相依解析 → 啟動監控」的通用流程。

* * *

✅ 功能說明
------

| 功能 | 說明 |
| --- | --- |
| 模組註冊 | 支援 `RegisterModule()` / `RegisterInvisibleModule()`，註冊各個模組的啟動邏輯 |
| 依賴管理 | 使用 `dependencyMap` 建立模組依賴圖，保證順序與正確啟動 |
| 生命週期控制 | 管理所有模組的啟動、異常監控與停止流程（`StartAsync`, `AwaitStopped`, `Shutdown`） |
| 健康狀態監控 | 實作 `services.ManagerListener` 監控模組健康狀態，並在錯誤發生時全域停機 |

* * *

📦 實作重點（取自代碼）
-------------

- `pkg/modules/modules.go`: 核心 service 實作，實作了 `Engine`, `Manager` interfacegrafana-pkg-all-code
    
- `pkg/modules/dependencies.go`: 定義了模組名稱與相依圖，如 `StorageRing`, `FrontendServer`, `GrafanaAPIServer` 等grafana-pkg-all-code
    
- `pkg/modules/listener.go`: 定義模組出錯時的應對策略（例如 shutdown 全部）grafana-pkg-all-code
    

* * *

🔌 與其他模組的關聯
-----------

| 模組 | 角色 |
| --- | --- |
| `pkg/server` | 使用 `modules.New()` 啟動所有模組，作為 server 啟動的一部分 |
| `pkg/services` | 各服務（如 `authz`, `frontend`, `live`, `grpcserver`）會以「模組」身份註冊進來 |
| `github.com/grafana/dskit/modules` | Grafana 使用此套件來統一模組與 service 管理邏輯 |

* * *

🧠 設計動機與適用場景
------------

- 為了解耦 Grafana 核心服務（不必全綁死在 `main()` 啟動）
    
- 支援模組化部署與條件式啟用（僅啟動特定模組）
    
- 提高系統在背景服務 crash 時的可觀察性與容錯性
    

* * *

✅ 適合你參考的地方（detectviz 模組化啟動設計）
-----------------------------

你可以將這套架構視為：

- 比 `registry.RunAll()` 更有「生命週期控制」的通用模組啟動器
    
- 支援模組健康狀態監控與出錯處理
    
- 利用 `RegisterModule()` + `Run()` 構建 background pipeline
    

# pkg/server:
`pkg/server` 目錄是 Grafana 啟動與模組管理的核心邏輯之一，負責主程式啟動、模組註冊、HTTP/GRPC 伺服器初始化、與模組服務（如 distributor、storage、frontend）的協調與執行。以下是其主要功能說明：

* * *

📦 `pkg/server` 功能說明
--------------------

### 1. Server 啟動邏輯（`server.go`）

- 提供 `New()` 與 `Run()` 方法建立與執行核心 `Server` 實例。
    
- 管理背景服務的啟動與關閉。
    
- 負責初始化 log、設定檔、PID 寫入與 systemd 通知等grafana-pkg-all-code。
    

### 2. 模組伺服器（`module_server.go`）

- `ModuleServer` 是為了模組化 Grafana 運行的入口點設計。
    
- 支援 `init()` 啟動各個模組：如 `initDistributor()`、`initFrontend()`、`initInstrumentationServer()` 等。
    
- 用於建構 dskit-based 模組並管理模組依賴關係grafana-pkg-all-code。
    

### 3. 模組註冊與服務管理（與 `pkg/modules` 整合）

- 使用 `modules.Engine` 與 `services.BasicService` 管理模組生命週期。
    
- 各模組如 `Distributor`, `StorageServer` 等會用 `WithName()` 註冊名稱，方便統一管理與診斷grafana-pkg-all-code。
    

### 4. Instrumentation 伺服器（`instrumentation_service.go`）

- 啟動 Prometheus metrics HTTP server，對外提供 `/metrics` endpoint。
    
- 也可設定為健康檢查與 profiling 伺服器（透過 `gorilla/mux` 實作）grafana-pkg-all-code。
    

### 5. 模組型服務定義（如 `distributor.go`, `memberlist.go`）

- 每個模組會實作自己的 `initXXX()` 函式，返回一個符合 `services.Service` 介面的執行單元。
    
- e.g., `initDistributor()` 初始化一個 grpc-based 分發模組服務grafana-pkg-all-code。
    

### 6. Runner 與 ModuleRunner（`runner.go`, `module_runner.go`）

- `Runner` 是整體伺服器的執行者，注入 config、密鑰、user service 等元件。
    
- `ModuleRunner` 是簡化版，只載入模組需要的最小依賴（例如 feature toggles）grafana-pkg-all-code。
    

* * *

📂 總結
-----

| 子檔案 | 功能概要 |
| --- | --- |
| `server.go` | 核心伺服器建立與生命週期管理 |
| `module_server.go` | 模組化伺服器，管理多個 dskit 模組 |
| `instrumentation_service.go` | 提供 Prometheus metrics HTTP 伺服器 |
| `distributor.go` | 建立 grpc distributor 服務 |
| `memberlist.go` | 建立 kv-based 成員列表協調模組 |
| `runner.go` / `module_runner.go` | 提供 Grafana 啟動執行流程 |


## pkg/plugins 功能

此模組為 Grafana plugin 的生命週期管理核心，負責：
	•	掃描 plugins 來源（如：本地磁碟、CDN、設定檔）
	•	建立 plugin 實體（含 bootstrap、validation、init、terminate）
	•	啟動/停止 plugin backend process（如 data source backend plugin）
	•	簽章驗證（signature validation）
	•	支援 plugin class、plugin.json metadata、plugin 子元件（children）

具體模組子系統包含：
	•	manager/registry：plugin 寄存與查找服務
	•	manager/sources：plugin 來源（disk、CDN、配置）
	•	manager/pipeline：plugin load pipeline（discovery → bootstrap → validation → initialization → termination）
	•	manager/process：plugin backend 執行流程
	•	codegen：為 plugin 生成註冊檔案
	•	pluginscdn：CDN 插件資產路徑構建

⸻

## pkg/registry 功能

此模組為 Kubernetes 風格的資源控制與 API 物件註冊模組，定位在：
	•	定義 API 物件的存取邏輯（透過 k8s.io/apiserver/pkg/registry/rest）
	•	提供 CRUD 功能、OpenAPI 文件產出、RBAC 評估等
	•	對應 apiserver 的物件 REST endpoint，如 provisioning.Repository、userstorage、query

⸻

## 調用關係
	•	pkg/plugins 引用 pkg/registry
	•	pkg/plugins/codegen/jenny_plugin_registry.go 使用 pkg/registry/schemas 做為 plugin registry 檔案生成的儲存目錄 ￼。
	•	插件註冊時可能與 registry 管理的資源互動，例如：datasource 以 plugin ID 做為來源識別。
	•	pkg/registry 引用 pkg/plugins
	•	pkg/registry/apis/query/client/plugin.go 中，pluginRegistry 與 pluginClient 明確依賴 plugins.Plugin 與 pluginstore.Store ￼。
	•	這代表 query 服務是透過 registry 進行 plugin-based datasource 的 CRUD 與調用。

⸻


# pkg/registry:
`pkg/registry` 是 Grafana 架構中用來實作 Kubernetes-style API 資源註冊與管理系統 的核心模組。它是 Grafana 向平台化（Platform-as-a-Framework）邁進的重要基礎，讓不同模組（如 dashboard、datasource、secret）以資源（Resource）的形式註冊、暴露、授權與操作。

* * *

✅ 功能總覽
------

| 功能分類 | 說明 |
| --- | --- |
| 資源註冊 | 每個資源（如 `DashboardSnapshot`, `Datasource`, `IAM`）都透過 `register.go` 實作 `APIGroupBuilder` 並註冊到 API Servergrafana-pkg-registry-co… |
| 資源實作 | 每個資源會實作 `rest.Storage` / `rest.Connecter` 等介面，並定義 CRUD、validation、connect API 等行為 |
| 授權與驗證 | 整合 K8s 的 `authorizer.Attributes` 與 Grafana 身分管理來驗證命名空間與權限grafana-pkg-registry-co… |
| OpenAPI 文件 | 整合 `kube-openapi` 自動產生每個資源的 schema 文件與前端表單 |
| CUE schema 驗證 | 與 `pkg/kinds` 配合，支援 YAML / JSON 的格式驗證與轉換grafana-pkg-registry-co… |

* * *

架構與關聯模組
----------

| 模組 | 說明 |
| --- | --- |
| `pkg/registry/apis/` | 資源註冊實作目錄，每個子目錄對應一種 GVK（如 `iam`, `datasource`, `secret`） |
| `pkg/registry/apps/` | 將 app plugin 整合為 registry 註冊源（如 `advisor`, `playlist`） |
| `pkg/registry/usagestatssvcs/` | 整合統計上報模組，提供每個資源用量資料grafana-pkg-registry-co… |
| `pkg/registry/apis/{資源}/contracts/` | 抽象出每個資源的存取接口，如 `DecryptStorage`, `EncryptionManager`, `Repository`grafana-pkg-registry-co… |

* * *

📦 關鍵機制與介面
----------

- `APIGroupBuilder`：每個資源註冊單位，提供 GroupVersion、Authorizer、OpenAPI 定義等
    
- `rest.Storage` / `rest.Connecter`：定義 REST API 行為（如 Get/List/Post/Patch/Connect）
    
- `ExtraBuilder`：提供 mutate、OpenAPI 擴充、job 註冊等延伸邏輯grafana-pkg-registry-co…
    
- `GenericRegistryStore`：包裝資源資料表與策略處理邏輯，與 `apiserver` 整合
    

* * *

✅ 使用場景
------

| 情境 | 使用 `pkg/registry` 功能 |
| --- | --- |
| plugin 註冊自己的 REST 資源 | 定義 GVK，實作 `register.go`，呼叫 `RegisterAPI(builder)` |
| YAML / JSON 上傳並驗證格式 | 使用 `DecodeYAMLObject()` → `cue.Validate()` 檢查 |
| 自訂使用者資源管理（如 secret） | 實作 `rest.Storage` 與 `contracts.Database`，並註冊路由 |

* * *

🧠 結語：為何需要 `pkg/registry`
-------------------------

Grafana 過去以 dashboard 為主的單一資源中心，逐步走向多資源、多 API 群組的架構。`pkg/registry` 提供：

- 資源模組的統一註冊入口
    
- API group 的 自動掛載與生命週期管理
    
- 使用者自定資源（如 keeper, snapshot, serviceaccount）的一致性 CRUD 與驗證
    

* * *

若你要在 detectviz 支援「自訂資源註冊」與「YAML 驗證」，這個模組是非常值得借鏡的。

## 小結：兩者關係與建議

模組	功能定位	是否依賴對方	是否耦合
pkg/plugins	Plugin 掃描、驗證、啟動、管理	✅ 讀寫 pkg/registry/schemas	弱耦合
pkg/registry	K8s-like 物件 CRUD 與 OpenAPI	✅ 調用 plugins.Plugin 作為核心物件	中度耦合

若你計畫將 plugin 寄存與 CRUD 資源設計更平台化，建議：
	•	將 pkg/plugins/manager/registry 明確命名為 pluginstore 或 pluginregistry
	•	將 pkg/registry/apis/query 等模組抽象化支援非 plugin 型 data source
	•	若未來 plugin 資源設計類似 CRD，可將 plugin 註冊與 metadata 存儲納入 registry 統一管理

