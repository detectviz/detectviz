Directory structure:
grafana-reference-note/
├── grafana-extensions.md
├── grafana-pkg-aggregator.md
├── grafana-pkg-apimachinery.md
├── grafana-pkg-auth.md
├── grafana-pkg-build.md
├── grafana-pkg-db.md
├── grafana-pkg-event.md
├── grafana-pkg-expr.md
├── grafana-pkg-ifaces-gcsifaces.md
├── grafana-pkg-login.md
├── grafana-pkg-modules.md
├── grafana-pkg-other.md
├── grafana-pkg-plugins.md
├── grafana-pkg-registry.md
├── grafana-pkg-server.md
├── grafana-pkg-setting.md
├── grafana-pkg-tests.md
├── grafana-pkg-util.md
├── grafana-pkg-web.md
└── services-featuremgmt.md

================================================
FILE: grafana-repo-notes/grafana-extensions.md
================================================
# extensions 概覽

## Response:
`pkg/extensions` 目錄（或其相關結構，如 `ExtensionsV2`）在 Grafana 中的功能為：

* * *

✅ 功能定位：**UI 插件可插拔能力的描述與配置支援**
-----------------------------

這個模組**不是用來提供後端 runtime 插件邏輯**，而是針對 plugin.json 中 `extensions` 欄位提供資料結構、解析與轉換邏輯。其目的為讓前端 UI 可以從 plugin 中**插入元件、連結、功能點等擴展點（extension points）**。

* * *

🔧 定義的核心結構
----------

來自 `ExtensionsV2` 結構grafana-pkg-plugins-code：

```go
type ExtensionsV2 struct {
  AddedLinks        []AddedLink
  AddedComponents   []AddedComponent
  ExposedComponents []ExposedComponent
  ExtensionPoints   []ExtensionPoint
  AddedFunctions    []AddedFunction
}
```

每一項代表：

| 欄位 | 說明 |
| --- | --- |
| `AddedLinks` | 插入 UI 的導航連結（如左側選單、設定連結） |
| `AddedComponents` | 插入 UI 的可視元件 |
| `ExposedComponents` | Plugin 暴露給其他 plugin 可呼叫的前端元件 |
| `ExtensionPoints` | 插件所提供的可擴展位置（可被其他插件填入） |
| `AddedFunctions` | 插件提供的可呼叫函式（可能是 JS 層 Function） |

* * *

🧩 使用場景
-------

*   Plugin 在 `plugin.json` 中描述：
    
    ```json
    "extensions": {
      "addedLinks": [...],
      "addedComponents": [...],
      "extensionPoints": [...]
    }
    ```
    
*   Grafana 啟動時解析 plugin.json → 轉換成 `ExtensionsV2` 結構
    
*   前端根據 extensionPoint 的 ID 決定要將哪個 plugin 的元件渲染在哪裡
    

* * *

❌ 與 plugin 執行邏輯無關
-----------------

`pkg/extensions` 或其對應資料結構 **並不涉及 plugin 的後端實作、資料查詢、監控功能等**，純粹作為：

*   plugin metadata 描述的 struct 定義
    
*   提供 backward compatibility（支援 V1 / V2 格式）
    

* * *

✅ 與 plugin.json 的關係
-------------------

`Extensions` 結構通常出現在：

```go
type JSONData struct {
  ...
  Extensions Extensions `json:"extensions"`
}
```

並由 `UnmarshalJSON()` 做版本兼容轉換（如 V1 → V2）grafana-pkg-plugins-code。

* * *

🧠 結論：這是前端擴展點的**資料描述模組**
------------------------

| 對象 | 用途 |
| --- | --- |
| `pkg/extensions` | 定義與解析 plugin.json 的 UI 插件擴展格式 |
| `pkg/plugins` | Plugin 掃描、載入、執行邏輯核心 |
| `pkg/components` | 共用邏輯模組，如快取、轉換、狀態機 |
| `pkg/services/...` | 具體的後端 service 或 plugin 註冊實作 |

* * *

需要的話我可以幫你畫出一張「Extensions 機制的 Plugin 與 UI 對應架構圖」，或針對 detectviz 設計你自己的 plugin extensionPoint 格式。要幫你建議嗎？

## Response:
`pkg/extensions` 模組的產生，是為了解決 Grafana 在插件生態系中出現的一個**關鍵需求：讓外部 plugin 能夠「非侵入式地擴展 UI 功能與互動」**。

* * *

✅ 問題背景：傳統 plugin 模式的限制
----------------------

早期 Grafana 的 plugin 架構專注於：

*   提供資料來源（Datasource plugin）
    
*   視覺化面板（Panel plugin）
    
*   告警、設定頁等單點功能
    

然而出現以下限制：

1.  插件無法自定義「新增 UI 區塊」或「掛載額外的功能頁」
    
2.  插件無法「動態對接其他 plugin 的畫面或邏輯」
    
3.  缺乏類似 `extension point` 的機制，導致擴展時只能 fork Grafana 原始碼
    

* * *

🧩 為了解決這些擴展困境，Grafana 引入 `extensions` 概念
----------------------------------------

類似於：

| 概念 | 類似對應架構 |
| --- | --- |
| `extensionPoints` | Eclipse 的 plugin point |
| `addedComponents` | React 插槽 / 插件注入 |
| `exposedComponents` | JS module export for inter-plugin calls |

這種設計方式讓 plugin：

*   可以聲明「我要把元件掛在某個指定位置」
    
*   可以提供可被其他 plugin 呼叫的 UI 元件或函式
    
*   可以在 UI 菜單中插入自己的設定連結或頁面
    

* * *

🧠 為什麼封裝成 `pkg/extensions` 模組？
------------------------------

Grafana 必須要：

*   **將 plugin.json 中的 `extensions` 轉換為嚴格定義的 struct（`ExtensionsV2`）**
    
*   **處理不同版本格式（V1/V2）的轉換與相容性處理**
    
*   **提供工具與 interface 給前端與 runtime 做 UI 組合與渲染**
    

這些需求屬於資料描述與解析邏輯，因此獨立為 `pkg/extensions`，而非混入 plugin runtime。

* * *

✅ 整體意圖：**平台化生態系擴展**
-------------------

Grafana 正在從「單一儀表板工具」轉型為「平台型生態」，而 `extensions` 就是這個方向的基礎：

*   前端：提供 plugin 注入點與視覺化組合
    
*   後端：可搭配 `apiserver` + `registry` 自訂資源物件
    
*   Plugin 作者：能做完整 UI + API 的擴展，而不 fork Grafana
    

* * *

如果你想在 detectviz 中也支援這類 extension 模式，我可以幫你設計一個 Plugin Interface 規格，例如：

```json
{
  "extensionPoints": ["alert-summary", "custom-dashboard-toolbar"],
  "addedComponents": [{ "type": "button", "target": "alert-toolbar", "id": "csv-exporter" }]
}
```




================================================
FILE: grafana-repo-notes/grafana-pkg-aggregator.md
================================================
# grafana/pkg/aggregator:
`grafana/pkg/aggregator` 目錄主要是提供 **Grafana Aggregator API Server** 的實作與設定，這個 aggregator 是用來整合統一的 RESTful API 對外暴露，例如：資料存取、多租戶資料彙整、查詢服務等。它與你所關注的 `expr`、`ngalert` 中的數據聚合（如 reduce、sum）**並沒有直接關聯**。

* * *

✅ `pkg/aggregator` 功能說明
-----------------------

| 區塊 | 說明 |
| --- | --- |
| `apiserver/` | 使用 K8s `genericapiserver` 建立 aggregator 專用的 API Server。類似於 K8s Aggregated API 架構。 |
| `apis/aggregation/v0alpha1/` | 定義 Aggregator 專用的 API Group 與資源格式，例如：聚合查詢、彙整規則。 |
| `options/` | 提供 `GrafanaAggregatorOptions`，供主系統啟動時注入 aggregator 配置與 OpenAPI 文件grafana-pkg-services-al…。 |
| `scheme/` | 提供 aggregator API 的 runtime scheme 註冊與 OpenAPI 定義整合。 |

* * *

🔌 Interface：`AggregatorRunner`
-------------------------------

定義在 `pkg/services/all-code` 中：

```go
type AggregatorRunner interface {
    Configure(opts *Options, config *RecommendedConfig, delegateAPIServer DelegationTarget, scheme *runtime.Scheme, builders []APIGroupBuilder) (*GenericAPIServer, error)
    Run(ctx context.Context, transport *RoundTripperFunc, stoppedCh chan error) (*GenericAPIServer, error)
}
```

用於初始化與啟動 aggregator server，實作可支援：

*   可擴充的 API group
    
*   自訂 OpenAPI schema
    
*   熱注入 API 實體與服務
    

* * *

與 `expr` / `ngalert` 差異比較
-------------------------

| 模組 | 用途 | 聚合功能來源 |
| --- | --- | --- |
| `expr` | 表達式運算（數據來源、運算、閾值判斷） | 由 `reduce` node / `math` node 實作 |
| `ngalert` | Alert 評估與狀態管理 | 依賴 `expr` pipeline 結果 |
| `aggregator` | 提供獨立 API server 對外整合資料服務 | 與數據計算無直接關係，偏向統一入口與 OpenAPI 定義 |

* * *

如你要了解的是「數學運算與聚合方法（sum, avg, count）」的程式碼位置，那麼應查看的是：

*   `pkg/expr/mathexp/`：包含 `reducer.go`, `resample.go`
    
*   `pkg/services/ngalert/eval/`：在 alert pipeline 評估中調用 reducer 與 upsample logicgrafana-pkg-services-ng…
    

需要我幫你對比這些模組如何共用 `reducer` enum 與邏輯嗎？還是你想聚焦在 `aggregator` 如何整合成獨立 server？

## Response:
`pkg/aggregator` 的啟動是條件式觸發的，**只有在某些 Feature Toggle（功能開關）啟用時才會初始化與執行**。以下是整理的觸發時機與呼叫流程：

* * *

✅ 啟動條件
------

Aggregator server 的啟動由以下兩個 Feature Flag 控制：

| Feature Flag | 說明 |
| --- | --- |
| `kubernetesAggregator` | 控制是否啟動 Kubernetes Aggregator |
| `dataplaneAggregator` | 控制是否啟動 Dataplane Aggregator |

這兩者由 `s.features.IsEnabledGlobally(...)` 判斷是否啟動grafana-pkg-services-al…。

* * *

🚀 啟動流程（摘要）
-----------

1.  **Grafana server 啟動進入點**：
    
    *   `main.go` → `commands.ServerCommand()` → `RunServer(...)`
        
    *   呼叫到 `apiserver.service.Run()`，此處控制所有 apiserver 啟動邏輯。
        
2.  **進入 Aggregator 啟動判斷**：
    
    ```go
    if isKubernetesAggregatorEnabled {
        aggregatorServer, err := s.aggregatorRunner.Configure(...) // 初始化
        ...
        runningServer, err = s.aggregatorRunner.Run(...)           // 啟動
    }
    ```
    
    若啟用且配置成功，就會呼叫 `Run()` 啟動整個 aggregator servergrafana-pkg-services-al…。
    
3.  **額外分支：Dataplane Aggregator** 也可能會被 `startDataplaneAggregator()` 啟動（通常與 plugin runtime 整合）。
    

* * *

🧩 `Configure()` 和 `Run()` 的具體內容
--------------------------------

兩個方法都實作自 `AggregatorRunner` interface：

```go
type AggregatorRunner interface {
    Configure(opts, config, delegate, scheme, builders)
    Run(ctx, transport, stoppedCh)
}
```

*   `Configure()` 會建立 server config、初始化 OpenAPI 設定與 etcd 儲存設定grafana-pkg-services-al…。
    
*   `Run()` 會：
    
    *   建立 `GenericAPIServer`
        
    *   設定 loopback handler
        
    *   呼叫 `PrepareRun()` 並以 goroutine 執行 server
        

* * *

⚙️ 什麼情況下不會啟動？
-------------

如果：

*   Feature toggle `kubernetesAggregator` 沒有打開（預設通常是關閉）
    
*   或是 Grafana 不是 enterprise 版本（部分實作會 fallback 成 noop）
    
*   或是 Grafana 沒有註冊任何 aggregator API group
    

就會完全跳過 aggregator 的 `Configure` / `Run()`，轉而執行 core servergrafana-pkg-services-al…。

* * *

是否需要我幫你畫出一張包含 `Grafana Server → apiserver → aggregator` 的 UML 或啟動流程圖？或者你要自己擴充 aggregator 模組？

## Response:
正確，**Grafana OSS 版本預設是不會啟動 `pkg/aggregator` 的**。原因如下：

* * *

✅ 為什麼 OSS 版本不會使用 `aggregator`
-----------------------------

1.  **功能旗標（Feature Flag）預設關閉**：
    
    *   `kubernetesAggregator` 和 `dataplaneAggregator` 這兩個開關在 OSS 中預設是 **關閉的**。
        
    *   這些 flag 屬於 Enterprise feature toggle 管理的範圍grafana-pkg-services-al…。
        
2.  **依賴 Enterprise 服務與模組**：
    
    *   `aggregatorRunner` 的實作通常會註冊於 Enterprise build 的 registry 中。
        
    *   OSS 版本中會 fallback 成 `noopAggregatorRunner`，其 `Run()` 及 `Configure()` 都是空實作grafana-pkg-all-code。
        
3.  **模組用途多半為 Enterprise 整合與多租戶 API 擴展**：
    
    *   例如統一整合各種 Plugin Runtime API（如 ML plugin, Dataplane plugin）、Kubernetes CRD 管理。
        
    *   在 OSS 環境下無此需求。
        

* * *

🔍 你可以這樣確認是否有啟動
---------------

在啟動日誌中：

*   若未出現與 aggregator server (`starting generic API server`, `registered API groups`) 相關訊息，則表示未啟動。
    
*   可檢查 feature flags 中是否有：`kubernetesAggregator`, `dataplaneAggregator`
    

* * *

🧪 若你想在 OSS 嘗試啟用
----------------

這**理論上可行**，但實務上會遇到以下限制：

*   需要將 `aggregatorRunner` 的實作注入 `registry`（這在 Enterprise 內部完成）
    
*   缺乏必要的 API group handler（例如：`pkg/aggregator/apis/aggregation/v0alpha1` 的實際用途未暴露）
    
*   可能還要 patch 一些 `build.go` 或 `main.go` 注入流程
    

* * *

若你有「自行擴充 Aggregator API group」的需求，或想將其變成 OSS 版本的一部分，我可以協助你**抽出可獨立運行的版本**，或者幫你對照 OSS 與 Enterprise 的差異設計。是否要我補這部分？

## Response:
非常有參考性。即使你不打算啟用 `pkg/aggregator`，它在 Grafana 中的「條件式觸發啟動流程」具有高度參考價值，尤其對你目前架構有幾個重要啟發：

* * *

✅ 可參考的設計模式
----------

### 1. **Feature Toggle 驅動模組啟動**

```go
if s.features.IsEnabledGlobally("kubernetesAggregator") {
	s.aggregatorRunner.Configure(...)
	s.aggregatorRunner.Run(...)
}
```

*   ✅ **優點**：允許功能模組隨時被打開 / 關閉，且不影響其他系統模組。
    
*   🧩 **應用場景**：你在 `detectviz` 中若要啟用：
    
    *   `alert` 模組
        
    *   `automation` 模組
        
    *   `plugin/llm` 模組  
        可採類似設計，讓這些模組在 config 或環境變數開啟後才進入啟動流程。
        

* * *

### 2. **interface + noop fallback 實作**

```go
type AggregatorRunner interface {
	Configure(...) error
	Run(...) error
}

type noopAggregatorRunner struct{}

func (n *noopAggregatorRunner) Configure(...) error { return nil }
func (n *noopAggregatorRunner) Run(...) error { return nil }
```

*   ✅ **優點**：
    
    *   在主流程中可以不做判斷，統一呼叫。
        
    *   測試、OSS/EE 拆分、模組選擇都變得非常簡潔。
        
*   🧩 **應用建議**：你可以為未啟用的模組提供 `noop.Service`，例如：
    

```go
type SchedulerService interface {
	Start() error
	Stop() error
}

type noopScheduler struct{}

func (n *noopScheduler) Start() error { return nil }
func (n *noopScheduler) Stop() error { return nil }
```

* * *

### 3. **延遲初始化、避免全域耦合**

Aggregator server 並不是一開機就執行，而是延遲在 `apiserver.Run()` 中根據條件註冊並啟動，這種設計：

*   ✅ **避免初始化時就加載所有模組**
    
*   ✅ **讓模組切換、熱部署、Mock 測試更簡單**
    

你也可以在 `detectviz` 的 orchestrator 中參考這種延遲註冊設計。

* * *

📦 範例應用：Detectviz 模組延遲啟動
------------------------

```go
if config.Features["automation"] {
	registry.Register("automation", automation.NewService(...))
}
```

或用 DI + interface + noop：

```go
type Automation interface {
	Run() error
}

func NewAutomation(cfg Config) Automation {
	if cfg.Enabled {
		return &realAutomation{}
	}
	return &noopAutomation{}
}
```

* * *

如你有意針對 detectviz 的模組設計一套「feature toggle + 延遲注入」架構，我可以幫你列出標準 interface、主流程注入範本與組織策略，要幫你規劃嗎？



================================================
FILE: grafana-repo-notes/grafana-pkg-apimachinery.md
================================================
# grafana/pkg/apimachinery:
`grafana/pkg/apimachinery` 是 **Grafana 自定義的 apimachinery 抽象層**，靈感來自 Kubernetes 的 `k8s.io/apimachinery`，但它並不是直接取用 Kubernetes 的，而是將某些機制內化並整合為 Grafana 自身的「統一資源管理基礎層」。

* * *

✅ 功能總覽
------

| 功能類別 | 說明 |
| --- | --- |
| **Identity 身分識別** | 定義 `identity.Requester`，用來標準化存取上下文中的身分，例如登入者、service account 等 |
| **MetaAccessor** | 抽象物件的 metadata 存取工具（annotation、label、resourceVersion 等） |
| **Error 轉換** | 將內部錯誤轉換為 `k8s.io/apimachinery/api/errors` 相容格式，供 apiserver 回應使用 |
| **OpenAPI 整合** | 提供自訂資源的 OpenAPI 定義與轉換工具（如 `GetOpenAPIDefinitions`） |
| **Namespace 與權限控制** | 定義類似 k8s 的 namespace 與 RBAC 群組概念，供 apiserver 授權判斷 |

* * *

🔌 Interface / 模組範例
-------------------

*   `identity.GetRequester(ctx)`：從 context 取得呼叫者身分
    
*   `MetaAccessor(obj)`：取得物件中的 metadata 欄位（支援泛型）
    
*   `utils.AnnoKeyCreatedBy` 等：annotation 的標準鍵名
    
*   `ErrorLister`：提供 `field.ErrorList`，與 k8s validation 系統兼容
    

* * *

🧩 與 apiserver/aggregator 的整合
-----------------------------

| 模組 | 功能 |
| --- | --- |
| `apiserver/auth/authenticator.go` | 使用 `identityAuthenticator` 來整合 context 中的 Requester 作為登入判斷依據grafana-pkg-services-al… |
| `apiserver/options/grafana-aggregator.go` | 結合 `commonv0alpha1.GetOpenAPIDefinitions()`，來生成 aggregator 的 OpenAPI 文件grafana-pkg-services-al… |
| `registry/apis/*/subresource.go` | 利用 `utils.MetaAccessor()` 處理 annotation、resource version、名稱等欄位grafana-pkg-registry-co… |

* * *

🚀 為何這個模組重要？
------------

因為 Grafana 正在逐步往「類似 Kubernetes 的資源抽象架構」前進，例如：

*   `dashboard`、`user`、`secret` 都有 GVK（GroupVersionKind）表示法
    
*   `registry/apis/` 中所有資源都使用 `runtime.Object` 介面包裝
    
*   實作類似 k8s 的 REST Storage 與 Webhook 設計
    

而 `pkg/apimachinery` 就是這一切的基礎，讓 Grafana 得以重用既有生態如 `client-go`、`watch.Interface`、`cache.Indexer`。

* * *

如果你正在設計 `detectviz` 的資源模型（像是 alert, rule, report），這個模組的架構可以**高度參考**。是否需要我幫你拆解 `identity` 模組的介面設計或 `MetaAccessor` 的泛型寫法？

## Response:
是「**仿照 Kubernetes 架構設計，但不限於 Kubernetes 模式才會用到**」。

* * *

✅ 更明確的說法是：
----------

> `pkg/apimachinery` 是 Grafana 內部設計的 **通用資源管理框架**，靈感來自 Kubernetes 的 `apimachinery`，但它並**不需要 Grafana 在 K8s 上執行，也不是僅限於 K8s 模式**。而是為了讓 Grafana 本身擁有類似的資源抽象與操作方式。

* * *

📦 使用場景不侷限於 Kubernetes
----------------------

### 在非 K8s 模式也會用到：

*   ✅ **Aggregator Server 模式**（如 plugin runtime 的擴展 API）
    
*   ✅ **Registry 架構下的資源 CRUD**（如 dashboards、secrets、teams 等）
    
*   ✅ **OpenAPI Schema 自動產生**
    
*   ✅ **Context 中身份管理（Requester）**
    
*   ✅ **Namespace、Annotation、Label 操作**
    

這些全部來自 `pkg/apimachinery` 的 model/metadata/identity/util/tools。

* * *

🧠 為何模仿 Kubernetes？
-------------------

因為 Grafana 正在轉型為一個「可管理、可擴展、多租戶」的**平台級產品**，仿照 K8s 的資源定義與 API 結構可以帶來幾個優勢：

| 優勢 | 說明 |
| --- | --- |
| 資源統一語意 | 每個物件都有 GVK、metadata、spec/status |
| 支援動態擴展 | 可像 K8s 一樣註冊 CRD 或 plugin API |
| 授權控制明確 | RBAC, namespace, ownerReference 模式可用 |
| Schema 驗證自動化 | 可產生 OpenAPI + JSON Schema |
| 客戶端協定一致 | 將來支援 REST/gRPC/Watch 就像 K8s API Server 一樣 |

* * *

🧪 總結：你可以這樣看
------------

| 模式 | 是否會用到 `pkg/apimachinery` |
| --- | --- |
| 本地部署 OSS | ✅ 用於 registry 與 plugin 資源模型 |
| 雲端 SaaS | ✅ 用於多租戶管理與動態 API 擴展 |
| K8s Plugin | ✅ 與 Aggregator API 整合 |
| 非 K8s 模式 | ✅ Identity、Annotation、OpenAPI 皆會用到 |

* * *

是否需要我幫你畫一張「Grafana 的 internal API 資源架構圖」或「`apimachinery` 應用示意圖」？這對你定義 `detectviz` 自訂資源格式會非常有幫助。

## Response:
以下是 Grafana 原始碼中三個目錄的功能差異說明：

* * *

✅ `pkg/api`
-----------

### 🔧 功能：**傳統 HTTP 路由與 UI 對應**

*   提供早期 Grafana OSS 的 HTTP API 與 UI 路由（非 Kubernetes-style）。
    
*   實作傳統 RESTful handler，例如 `/api/health`、`/api/dashboards/db/...`。
    
*   採用自訂 middleware（如 `csrf`, `requestmeta`, `loggermw`）包裝。
    

### 📦 範例元件：

*   `http_server.go`: 核心 HTTPServer 實作與 TLS 啟動邏輯。
    
*   `routing/`: 定義各條路由與 handler。
    
*   `static/`: 提供前端資源（JS、CSS）靜態路由。
    

* * *

✅ `pkg/apis`
------------

### 🔧 功能：**以 Kubernetes CRD 風格定義 Grafana 內部資源 API**

*   每個子目錄對應一個 GVK（Group-Version-Kind），類似 CRD 定義。
    
*   提供強類型資源，如 `provisioning/v0alpha1`, `query/v0alpha1`, `secret/v0alpha1`。
    
*   與 `pkg/registry/apis` 搭配，提供完整的 REST Storage 與 Resource 註冊。
    

### 📦 用途：

*   定義各資源的 schema、OpenAPI metadata。
    
*   結合 `pkg/apimachinery` 供 `apiserver` 與 `aggregator` 使用。
    

* * *

✅ `pkg/apiserver`
-----------------

### 🔧 功能：**建構 Kubernetes 風格的 API Server（模組化、多租戶、OpenAPI）**

*   包含 Server 啟動器、路由鏈組裝器、認證與授權器、OpenAPI 支援、admission webhook。
    
*   可與 OSS 模式共存，也可在 aggregator 模式下條件式啟動。
    

### 📦 關鍵元件：

*   `service.go`: 啟動主流程，初始化 `genericapiserver.Config`。
    
*   `builder/`: 提供 `APIGroupBuilder`, `Admission`, `Metrics`, `OpenAPI` 等抽象。
    
*   `auth/`: 包含 `authenticator` 與 `authorizer`，支援 namespace-based RBAC。
    
*   `endpoints/`: 提供 custom response writer、context filters。
    

* * *

🧩 三者差異對照表
----------

| 目錄 | 角色 | 主要用途 | 接口風格 | 使用情境 |
| --- | --- | --- | --- | --- |
| `pkg/api` | 傳統 HTTP handler | 傳統 API 與 UI | RESTful, 路由為主 | OSS 基礎 |
| `pkg/apis` | 資源定義 | schema, CRD 結構, GVK | CRD style | `apiserver` 所需 |
| `pkg/apiserver` | API server runtime | 啟動、註冊、OpenAPI、auth | genericapiserver | Aggregator/Enterprise 模式 |

* * *

是否要我幫你畫出一張「Grafana API server 分層與對應元件圖」？這將有助於你為 detectviz 設計 plugin 或自定資源 API 結構。


================================================
FILE: grafana-repo-notes/grafana-pkg-auth.md
================================================
# /pkg/services/auth 總覽:
根據你提供的 Grafana 專案結構與原始碼（含 `grafana-pkg-services-all-code.md` 等多個檔案），我針對 `/pkg/services/auth` 目錄下的 interface 做了整理，以下是目前該模組中明確定義的介面（interface）與用途概覽：

* * *

`/pkg/services/auth` 中的主要介面（Interfaces）
---------------------------------------

### 1. `Authenticator` interface

**來源位置**：`pkg/services/auth/authenticator.go`

```go
type Authenticator interface {
    Login(ctx *models.ReqContext) (*models.User, error)
    SyncUser(user *models.User) error
}
```

**說明**：

*   此為登入與同步使用者資訊的核心介面。
    
*   實作範例如：LDAP、OAuth、Grafana native login 等，都會實作此介面。
    

* * *

### 2. `AuthTokenService` interface

**來源位置**：`pkg/services/auth/tokens.go`

```go
type AuthTokenService interface {
    CreateToken(userID int64, clientIP string) (*models.AuthToken, error)
    LookupToken(unhashedToken string) (*models.AuthToken, error)
    RevokeToken(unhashedToken string) error
}
```

**說明**：

*   用於處理 token 的建立、查詢與撤銷。
    
*   常用於實作 API Token 或 Session 管理。
    

* * *

### 3. `SignedTokenService` interface

**來源位置**：`pkg/services/auth/signedtoken.go`

```go
type SignedTokenService interface {
    CreateSignedToken(data string, expires time.Duration) (string, error)
    ValidateSignedToken(token string) (string, error)
}
```

**說明**：

*   用於產生具有簽章的安全 token，常用於 password reset、email verification。
    
*   透過簽名方式防止被竄改。
    

* * *

### 4. `LoginService` interface

**來源位置**：`pkg/services/auth/login.go`

```go
type LoginService interface {
    LoginUserWithUser(user *models.User, ctx *models.ReqContext) error
}
```

**說明**：

*   提供直接登入已知使用者的方法。
    
*   常見於開發或測試用途，例如 admin impersonation。
    

* * *

這些 interface 通常會在 `auth.go`, `tokens.go`, `signedtoken.go`, `login.go` 中以 struct 實作對應邏輯，並註冊到依賴注入的容器（通常是 `registry` 模組）。

* * *

需要我幫你列出這些介面在哪些實作中具體被用到，或在 `main.go` / `registry` 中如何註冊的嗎？

## Grafana Auth 原始碼結構:
根據 Grafana 原始碼結構與代碼分析，以下是你關注的五個目錄的 **功能與定義的 interface** 整理：

* * *

1. `pkg/services/auth`
-----------------------

### 🔧 功能：

統一定義與管理授權（Authentication）相關介面與抽象，包括 Token 管理、ID 簽署、外部 session 存取等。

### 🔌 Interface：

*   `UserTokenService`：管理 user auth token，例如建立、查詢、撤銷等。
    
*   `ExternalSessionStore`：處理外部登入會話（如 OAuth）的 CRUD 操作grafana-pkg-services-al…。
    
*   `IDService`：為 identity 請求者簽發 ID Token。
    
*   `IDSigner`：以 signing key 實作簽發 JWT。
    
*   `IDClaims`：JWT 內容的結構別名。
    

* * *

2. `pkg/services/authimpl`
---------------------------

### 🔧 功能：

此為 `auth` interface 的具體實作（implementation），例如 user token 實作、session 儲存、token 清除等。

### 🔌 Interface 實作：

*   `UserAuthTokenService`：實作 `UserTokenService`，提供 token lifecycle 管理、revoked tokens 查詢、rotation 機制等。
    
*   `store`：實作 `ExternalSessionStore`，負責 DB 交易與密文儲存解密grafana-pkg-services-al…。
    

* * *

3. `pkg/services/gcomsso`
--------------------------

### 🔧 功能：

Grafana Cloud 的單一登出（SSO）hook 實作，提供整合外部系統的 logout 機制。

### 🔌 Interface：

未定義獨立的公開 interface，主要提供 `gcom_logout_hook.go` 中的邏輯，作為 `authn` 登出 hook 的實作模組。

* * *

4. `pkg/services/idimpl`
-------------------------

### 🔧 功能：

提供 ID Token 的本地簽署實作，內部使用 JWT + JOSE + signingkeys，支援組織角色、快取、claim 建構等。

### 🔌 Interface 實作：

*   `LocalSigner`：實作 `IDSigner`，透過 `signingkeys.Service` 提供簽署服務grafana-pkg-services-al…。
    

* * *

5. `pkg/services/jwt`
----------------------

### 🔧 功能：

專責 JWT 的驗證邏輯，包括 key set 檢查、claim 驗證、token 解碼等。

### 🔌 Interface：

*   `JWTService`：
    
    ```go
    type JWTService interface {
        Verify(ctx context.Context, strToken string) (map[string]any, error)
    }
    ```
    
*   `FakeJWTService`：測試用 mock 實作grafana-pkg-services-al…。
    

* * *

如果你有需要，我也可以幫你畫出這些 interface 的 UML 關聯圖，或補上 registry 的註冊方式與注入流程。是否要繼續？



================================================
FILE: grafana-repo-notes/grafana-pkg-build.md
================================================

## pkg/build:
`pkg/build` 是 Grafana 用來在編譯期間注入「版本與建置資訊」的模組，屬於 **build metadata 管理元件**，搭配 Go build flags 使用。其功能與 `pkg/versions` 互補，但用途更接近「執行檔等級的 build metadata」。

* * *

✅ 主要功能
------

| 功能類別 | 說明 |
| --- | --- |
| 定義變數 | 定義變數如 `Version`, `Commit`, `BuildStamp`, `BuildEdition` 等 |
| 可由 `-ldflags` 動態注入 | 支援從 `go build` 傳入版本與 commit hash，自動寫入該模組 |
| 提供執行檔版本顯示 | 用於 CLI `grafana-server --version`、Web UI 中顯示版本資訊 |
| CLI 輸出用 | 搭配 `pkg/cmd/grafana-server` 使用，輸出建置訊息與環境資訊 |

* * *

📦 變數範例
-------

```go
package build

var (
  Version     = "dev"
  Commit      = "unknown"
  BuildStamp  = ""
  BuildEdition = "oss"
)
```

在 build 階段使用：

```bash
go build -ldflags "-X github.com/grafana/grafana/pkg/build.Version=10.3.1 \
                   -X github.com/grafana/grafana/pkg/build.Commit=abcdef123"
```

* * *

🧩 與 `pkg/versions` 差異比較
------------------------

| 模組 | 目的 | 使用場景 | 資料來源 |
| --- | --- | --- | --- |
| `pkg/build` | 記錄執行檔編譯資訊 | CLI 顯示、後台版本比對 | build flags 注入 |
| `pkg/versions` | 管理版本邏輯與 schema 版本 | 各模組內部比對用 | 程式內固定定義 |

兩者會互相參考，但 `build` 更偏向「Build pipeline 的輸入點」。

* * *

✅ 使用場景
------

*   `grafana-server --version` 輸出
    
*   Grafana 前端 `/public/views/index.html` 加入 `BuildEdition`
    
*   REST API `/api/frontend/settings` 回傳版本資訊
    
*   Prometheus metrics export 標示版本欄位
    

* * *

detectviz 借鏡建議
--------------

若你希望：

*   顯示 `--version` / `--commit` 等 CLI 資訊
    
*   前端或 API 回傳版本號
    

你可以仿照 Grafana 設計一個 `pkg/build/build.go` 並用以下方式設定：

```bash
go build -ldflags "-X github.com/detectviz/pkg/build.Version=$(VERSION) \
                   -X github.com/detectviz/pkg/build.Commit=$(GIT_COMMIT)"
```

# 遷移到 Detectviz /build 後怎麼做

pkg/build 屬於 執行檔等級的建置資訊模組（Build Metadata），其內容（如 Version, Commit, BuildStamp）不屬於應用邏輯或共用工具，也不應被其他模組 import 作為核心依賴。因此，建議如下：

⸻

✅ 最佳遷移位置：internal/buildmeta/ 或根目錄下的 build/

選項	說明	適用情境
internal/buildmeta/	表示此為內部框架層 metadata 專用，非模組介面一部分	✅ 遵循 Clean Architecture 分層，嚴格封裝
build/（專案根目錄）	與 main.go 同層，清楚表示此為執行檔層的 metadata 設定	✅ 適合簡單專案或希望集中控制版本輸出


⸻

🧱 具體建議結構（Detectviz 建議）

internal/buildmeta/
├── version.go     // 定義 Version, Commit, BuildStamp 等變數

或更簡化為：

build/
└── meta.go        // 根目錄下控制整體版號與 commit

並在 go build 時使用：

go build -ldflags "-X github.com/detectviz/internal/buildmeta.Version=1.0.0"


⸻

🔧 延伸應用建議

用途	說明
CLI 顯示	--version command 印出 buildmeta.Version
REST API	GET /api/meta 回傳版本、build edition 等
Prometheus	build_info{version=...} metrics 可加入 buildmeta 提供的欄位


⸻

✅ 總結

模組	遷移建議
pkg/build	⏩ 遷移為 internal/buildmeta/ 或根目錄 build/
變數命名	Version、Commit、BuildStamp 保留不變
建置參數	維持 -ldflags 注入策略


⸻

是否需要我協助重構路徑並補上最小版的 meta.go 範例？或協助你產生 --version CLI 輸出邏輯？


================================================
FILE: grafana-repo-notes/grafana-pkg-db.md
================================================

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


================================================
FILE: grafana-repo-notes/grafana-pkg-event.md
================================================
# pkg/events:
以下是 `grafana/pkg/bus` 與 `grafana/pkg/events` 兩個目錄的功能與差異說明：

* * *

✅ `pkg/bus`：**In-Process Message Bus**
--------------------------------------

### 📌 功能

*   提供 **同步發送與處理事件（Command / Message）的總線機制**，類似 CQRS 架構中的 Command Bus。
    
*   所有模組可透過 `bus.Publish(ctx, msg)` 發送訊息，由對應 handler 同步處理。
    

### 🔌 定義的介面與結構

```go
type Msg any
type HandlerFunc any

type Bus interface {
  Publish(ctx context.Context, msg Msg) error
  AddEventListener(handler HandlerFunc)
}
```

### 🧠 核心元件

*   `InProcBus`：Grafana 內部唯一實作。
    
*   `AddEventListener(fn)`：註冊處理特定訊息類型的 handler。
    
*   `Publish(ctx, msg)`：會找出該訊息類型的所有 handler 並呼叫grafana-pkg-all-code。
    

### 📦 使用場景

*   發送內部事件，例如：`FolderRenamed{FolderID, NewName}`。
    
*   減少模組之間直接依賴，達到模組解耦。
    

* * *

✅ `pkg/events`：**外部通知與系統事件格式定義**
--------------------------------

### 📌 功能

*   定義 **可傳送至外部系統（如 AMQP, Kafka, webhook）使用的事件資料結構**。
    
*   所有事件都為資料轉移物件（DTO），設計為**序列化安全且向後相容**。
    

### 🧾 事件格式範例：

```go
type OrgCreated struct {
	Timestamp time.Time
	Id        int64
	Name      string
}

type UserCreated struct {
	Timestamp time.Time
	Id        int64
	Name      string
	Login     string
	Email     string
}
```

### 📦 使用場景

*   事件可由內部觸發後寫入外部資料通道。
    
*   例如：新使用者註冊後觸發 `UserCreated` 事件，傳送到通知系統或審計系統grafana-pkg-all-code。
    

* * *

🔍 差異與關聯對照
----------

| 面向 | `pkg/bus` | `pkg/events` |
| --- | --- | --- |
| 定位 | 執行階段事件調度 | 定義靜態事件資料結構（DTO） |
| 資料流向 | 同步呼叫（in-process only） | 外部序列化使用（Kafka, webhook 等） |
| 用法 | `bus.Publish(ctx, &SomeCommand{})` | `emit(&UserCreated{...})` |
| handler 型態 | 函數註冊（透過 reflect 類型比對） | 不處理邏輯，只是 struct 定義 |
| 解耦程度 | 用於模組內部解耦 | 用於模組與外部世界解耦 |

* * *

### 🧩 組合應用建議（你可套用於 detectviz）

*   將 `bus` 作為模組間調度中心（例如 alert → reporter）
    
*   將 `events` 作為告警紀錄 / webhook / 通知格式標準層
    
*   分離處理邏輯（Handler）與資料定義（DTO），保持模組責任單一（符合 SRP）
    

是否需要我提供 bus 使用範例或在 detectviz 中建立事件發送範本？



================================================
FILE: grafana-repo-notes/grafana-pkg-expr.md
================================================
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


================================================
FILE: grafana-repo-notes/grafana-pkg-ifaces-gcsifaces.md
================================================



================================================
FILE: grafana-repo-notes/grafana-pkg-login.md
================================================
# grafana/pkg/login/ 概覽
`pkg/login` 目錄在 Grafana 中的主要功能是：**提供登入系統的抽象與實作，特別是整合各類 OAuth 社交登入（social login）與其邏輯處理機制**。

* * *

✅ 功能說明
------

| 子模組 | 說明 |
| --- | --- |
| `social/` | 定義各種 Social 登入 provider（GitHub、GitLab、Google、Okta、AzureAD 等）的共通介面與欄位結構 |
| `social/connectors/` | 每個 provider 對應一個 `*_oauth.go` 檔，包含 OAuth2 流程實作、API 解析、權限比對、帳號同步 |
| `socialimpl/` | 將設定（grafana.ini 或 dynamic settings）轉換為 provider 實例，註冊並提供存取方法 |
| `socialtest/` | 提供 `FakeSocialService`, `MockSocialConnector` 等測試替身與驗證工具 |

* * *

🔌 定義的核心介面
----------

```go
type SocialConnector interface {
  UserInfo(ctx, client, token) (*BasicUserInfo, error)
  IsEmailAllowed(email string) bool
  IsSignupAllowed() bool
  GetOAuthInfo() *OAuthInfo
}
```

```go
type Service interface {
  GetOAuthProviders() map[string]bool
  GetConnector(string) (SocialConnector, error)
}
```

這些介面讓其他模組（如 `authn`, `middleware`, `apiserver`）能夠統一調用各種登入來源，實現身份認證與授權。

* * *

🧩 與其他模組的關聯
-----------

| 模組 | 用途 |
| --- | --- |
| `pkg/services/authn/` | 將 `login.SocialConnector` 注入為 `authn.Client`，支援 token 驗證、自動註冊 |
| `pkg/services/ssosettings` | 允許使用者透過 API 設定 OAuth provider |
| `pkg/middleware/auth.go` | 在 API 呼叫過程中處理未登入者導向社交登入入口grafana-pkg-all-code |

* * *

🧠 設計意圖與可延伸價值
-------------

Grafana 將 `login/social` 拆為獨立模組，而非混入 `auth`、`authn`，目的是：

*   模組解耦（符合 SRP）
    
*   支援多種登入策略並行
    
*   可替換或禁用單一 provider
    
*   適合 SaaS 多租戶或動態設定場景
    

* * *

🚀 對 detectviz 的啟發
------------------

若你將來想支援：

*   OAuth2 API token 交換
    
*   單一入口整合 GitHub / Google 登入
    
*   動態切換登入策略或身份驗證模組
    

你可以參考這個設計：

*   interface 抽象放 `pkg/login`
    
*   provider 建構器放 `pkg/login/connectors`
    
*   統一注入與管理邏輯放 `pkg/login/service.go`
    
*   模擬測試放 `pkg/login/test/`
    

是否需要我幫你定義 detectviz 的 login interface 與 plugin-based 登入機制草稿？

# 遷移到 Detectviz internal/auth/ 後怎麼做

✅ 是的，你可以將 pkg/login/ 遷移至 internal/auth/，而且這麼做在 Detectviz 的分層原則下是合理且建議的。以下是詳細說明：

⸻

✅ 為什麼建議遷移至 internal/auth/

🎯 1. 屬於應用邏輯（Authentication Usecase）
	•	pkg/login/ 實際負責的是登入行為與 provider 驗證邏輯（OAuth 流程、帳號綁定…）
	•	它屬於 Detectviz 內部的身份驗證流程實作，而非獨立共用模組
	•	不太可能在其他專案或模組中重用（不像 redis client、http client 等）

⸻

📦 2. interface 抽象應該再整理進 pkg/iface/auth.go
	•	若你想保留 plugin 化的可能，可將 SocialConnector interface 抽象定義移至 pkg/iface/auth.go
	•	將所有實作與 provider 流程（如 google_oauth.go）搬至 internal/auth/social/

⸻

✅ 遷移建議結構

internal/auth/
├── social/              # provider 實作（google, github, ...）
├── socialimpl/          # provider 建構與註冊
├── socialtest/          # mock 與替身
├── service.go           # 登入流程邏輯
└── handler.go           # 登入 API 或表單頁

pkg/iface/auth.go        # 定義 SocialConnector 等抽象


⸻

🧭 延伸規劃（選擇性）
	•	若將來登入行為也支援 plugin 化，可讓 provider 實作放到 plugins/login-provider/xxx
	•	但在目前階段，「Detectviz 預設支援多種 provider」這件事仍屬於 app 內部邏輯，因此維持在 internal/ 更合理




================================================
FILE: grafana-repo-notes/grafana-pkg-modules.md
================================================
# grafana/pkg/modules 模組化啟動流程:
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

*   `pkg/modules/modules.go`: 核心 service 實作，實作了 `Engine`, `Manager` interfacegrafana-pkg-all-code
    
*   `pkg/modules/dependencies.go`: 定義了模組名稱與相依圖，如 `StorageRing`, `FrontendServer`, `GrafanaAPIServer` 等grafana-pkg-all-code
    
*   `pkg/modules/listener.go`: 定義模組出錯時的應對策略（例如 shutdown 全部）grafana-pkg-all-code
    

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

*   為了解耦 Grafana 核心服務（不必全綁死在 `main()` 啟動）
    
*   支援模組化部署與條件式啟用（僅啟動特定模組）
    
*   提高系統在背景服務 crash 時的可觀察性與容錯性
    

* * *

✅ 適合你參考的地方（detectviz 模組化啟動設計）
-----------------------------

你可以將這套架構視為：

*   比 `registry.RunAll()` 更有「生命週期控制」的通用模組啟動器
    
*   支援模組健康狀態監控與出錯處理
    
*   利用 `RegisterModule()` + `Run()` 構建 background pipeline
    
# 遷移到 Detectviz internal/auth/ 後怎麼做

依照你目前的 Detectviz 架構與分層原則，pkg/modules/ 這種「模組生命週期管理器」性質的元件，建議 不應放在 pkg/，而是遷移至：

⸻

✅ 遷移建議路徑：internal/infra/modules/

📌 理由：
	1.	屬於應用啟動框架的一部分，不具重用性
	•	雖然是模組化設計，但用途專屬於 Detectviz 的啟動組裝流程（如 DI 或 plugins 組裝）
	•	並不適合開放給外部套件或 app 使用，因此不應放在 pkg/
	2.	定位上屬於基礎設施（infra）
	•	與 logger, redis, httpserver 等類似，都是「支撐系統本體運作的服務控制器」
	3.	與 pkg/server, pkg/services 的耦合已知
	•	遷移後可以讓 app 組裝者更明確知道：這是 Detectviz 系統初始化流程的一部分，而非可以外部重用的模組

⸻

🧱 最終建議目錄：

internal/infra/
├── modules/           # 原 pkg/modules 遷移內容
│   ├── modules.go
│   ├── dependencies.go
│   └── listener.go
├── logger/
├── redis/
├── httpserver/


⸻

✅ 後續對應動作
	•	將 pkg/modules 全部搬至 internal/infra/modules
	•	interface 若需抽象化，可於 pkg/iface/module.go 定義 ModuleEngine 等（若 plugin 要擴充）
	•	原 pkg/server 或 app 組裝流程中使用 modules.New() 的地方調整 import

⸻

✅ 是的，將 modules 遷移至 internal/bootstrap/modules/ 是完全合理且更語義明確的做法，甚至比放在 internal/infra/ 更好，原因如下：

⸻

✅ 為什麼建議放 internal/bootstrap/modules/

理由	說明
🧠 功能屬於「系統啟動流程控制」	模組的註冊、依賴解析、啟動順序 → 明確屬於 bootstrap layer 的一部分
🔗 與 server, registry, services 密切協作	所有模組化的啟動、生命週期監控都在 bootstrap 階段進行
🚫 不應放在 infra/	infra 是提供通用技術元件，如 redis/logger，不處理主流程控制邏輯
📦 更清楚劃分責任	bootstrap/ 是組裝（composition）區，代表 Detectviz 本體啟動器，與 apps/ 對應，邏輯清晰


⸻

🧱 建議遷移結構

internal/bootstrap/
├── modules/
│   ├── modules.go            # Engine, Manager 實作
│   ├── dependencies.go       # 建構模組依賴圖
│   ├── listener.go           # 異常監控與 shutdown 策略
├── server/                   # 組裝模組與 API
├── runtime.go                # NewServer() / DI Container 啟動點

後續你若加上 bootstrap/events/、bootstrap/loader/，也能自然擴充這層功能

⸻

🧩 Interface 抽象（可選）

如果你打算讓 plugin 也可註冊模組，可：
	•	在 pkg/iface/module.go 中定義：

type Module interface {
    Name() string
    DependsOn() []string
    Run(ctx context.Context) error
}


	•	或定義 ModuleLifecycle、ModuleBuilder 介面供擴充

⸻

✅ 結論

遷移路徑	評價
internal/bootstrap/modules/	✅ 最佳選擇，語義清楚，擴充彈性高
internal/infra/modules/	⛔ 太偏向工具層，會誤導為通用模組
pkg/modules/	⛔ 高耦合應用邏輯，不該暴露在 pkg 層


⸻

是否需要我協助你進行完整搬遷（含調整 import）、interface 抽出、或自動化 module registry？我可以立即幫你處理。


================================================
FILE: grafana-repo-notes/grafana-pkg-other.md
================================================

# pkg/components:
`pkg/components` 是 Grafana 架構中專門用來存放「**小型通用邏輯元件與工具模組**」的目錄，這些模組雖不屬於核心層（如 `services`、`api`、`registry`），但會被多處業務重複使用，具有下列特性：

* * *

✅ 功能定位與設計原則
-----------

| 面向 | 說明 |
| --- | --- |
| 📦 封裝通用元件 | 提供簡單的 data structure 處理、JSON 工具、快取、metrics helper 等 |
| 🧩 可被其他模組引入 | 被 `services`, `registry`, `dashboard`, `publicdashboards` 等直接依賴 |
| 🔌 非 plugin 模組 | 與 plugins/ 不同，components 不涉及可載入的外部元件邏輯，而是內部共用邏輯 |

* * *

🧱 常見子模組功能範例
------------

根據實際代碼結構（如 `pkg/services/publicdashboards` 中的 import），以下是部分實例：

| 子模組 | 說明 |
| --- | --- |
| `components/simplejson` | 提供兼容性高的 JSON 讀寫工具（比官方 encoding/json 寬容） |
| `components/statemachine` | 定義狀態機處理邏輯，用於 alert / dashboard workflows |
| `components/cache` | 簡化的快取模組（LRU、Map-based） |
| `components/conversion` | 資料類型轉換（如 dto → model） |
| `components/crypto` | 加解密、雜湊等工具（非 infra/ 內的 TLS） |

* * *

🧩 與 plugins 的差異
----------------

| 比較項目 | `pkg/components` | `pkg/plugins` |
| --- | --- | --- |
| 類型 | 內部工具模組 | 外掛管理系統與 plugin runtime |
| 用途 | 被程式碼呼叫 | 被 runtime 載入 |
| 掛載方式 | 透過 import 使用 | 透過 plugin.json 掃描與註冊 |
| 關聯性 | 服務層共用邏輯 | 用戶擴展使用者功能與視覺化 |

* * *

💡 可參考用途（你在 detectviz 可套用的場景）
-----------------------------

*   將你未來的 `AlertLevelClassifier`、`EventBusAdapter`、`TimeRangeParser` 等邏輯模組統一放在 `pkg/components/` 中，保持可重用性又不與核心 service 模組耦合。
    
*   若某些模組需要簡單的快取、驗證、資料轉換，也可封裝於此。
    

* * *

是否需要我列出目前 Grafana 中完整的 `components` 子模組列表與用途？或幫你為 `detectviz/components/` 草擬標準目錄結構？


# pkg/extensions:
`pkg/extensions` 目錄（或其相關結構，如 `ExtensionsV2`）在 Grafana 中的功能為：

* * *

✅ 功能定位：**UI 插件可插拔能力的描述與配置支援**
-----------------------------

這個模組**不是用來提供後端 runtime 插件邏輯**，而是針對 plugin.json 中 `extensions` 欄位提供資料結構、解析與轉換邏輯。其目的為讓前端 UI 可以從 plugin 中**插入元件、連結、功能點等擴展點（extension points）**。

* * *

🔧 定義的核心結構
----------

來自 `ExtensionsV2` 結構grafana-pkg-plugins-code：

```go
type ExtensionsV2 struct {
  AddedLinks        []AddedLink
  AddedComponents   []AddedComponent
  ExposedComponents []ExposedComponent
  ExtensionPoints   []ExtensionPoint
  AddedFunctions    []AddedFunction
}
```

每一項代表：

| 欄位 | 說明 |
| --- | --- |
| `AddedLinks` | 插入 UI 的導航連結（如左側選單、設定連結） |
| `AddedComponents` | 插入 UI 的可視元件 |
| `ExposedComponents` | Plugin 暴露給其他 plugin 可呼叫的前端元件 |
| `ExtensionPoints` | 插件所提供的可擴展位置（可被其他插件填入） |
| `AddedFunctions` | 插件提供的可呼叫函式（可能是 JS 層 Function） |

* * *

🧩 使用場景
-------

*   Plugin 在 `plugin.json` 中描述：
    
    ```json
    "extensions": {
      "addedLinks": [...],
      "addedComponents": [...],
      "extensionPoints": [...]
    }
    ```
    
*   Grafana 啟動時解析 plugin.json → 轉換成 `ExtensionsV2` 結構
    
*   前端根據 extensionPoint 的 ID 決定要將哪個 plugin 的元件渲染在哪裡
    

* * *

❌ 與 plugin 執行邏輯無關
-----------------

`pkg/extensions` 或其對應資料結構 **並不涉及 plugin 的後端實作、資料查詢、監控功能等**，純粹作為：

*   plugin metadata 描述的 struct 定義
    
*   提供 backward compatibility（支援 V1 / V2 格式）
    

* * *

✅ 與 plugin.json 的關係
-------------------

`Extensions` 結構通常出現在：

```go
type JSONData struct {
  ...
  Extensions Extensions `json:"extensions"`
}
```

並由 `UnmarshalJSON()` 做版本兼容轉換（如 V1 → V2）grafana-pkg-plugins-code。

* * *

🧠 結論：這是前端擴展點的**資料描述模組**
------------------------

| 對象 | 用途 |
| --- | --- |
| `pkg/extensions` | 定義與解析 plugin.json 的 UI 插件擴展格式 |
| `pkg/plugins` | Plugin 掃描、載入、執行邏輯核心 |
| `pkg/components` | 共用邏輯模組，如快取、轉換、狀態機 |
| `pkg/services/...` | 具體的後端 service 或 plugin 註冊實作 |

* * *


`pkg/extensions` 模組的產生，是為了解決 Grafana 在插件生態系中出現的一個**關鍵需求：讓外部 plugin 能夠「非侵入式地擴展 UI 功能與互動」**。

* * *

✅ 問題背景：傳統 plugin 模式的限制
----------------------

早期 Grafana 的 plugin 架構專注於：

*   提供資料來源（Datasource plugin）
    
*   視覺化面板（Panel plugin）
    
*   告警、設定頁等單點功能
    

然而出現以下限制：

1.  插件無法自定義「新增 UI 區塊」或「掛載額外的功能頁」
    
2.  插件無法「動態對接其他 plugin 的畫面或邏輯」
    
3.  缺乏類似 `extension point` 的機制，導致擴展時只能 fork Grafana 原始碼
    

* * *

🧩 為了解決這些擴展困境，Grafana 引入 `extensions` 概念
----------------------------------------

類似於：

| 概念 | 類似對應架構 |
| --- | --- |
| `extensionPoints` | Eclipse 的 plugin point |
| `addedComponents` | React 插槽 / 插件注入 |
| `exposedComponents` | JS module export for inter-plugin calls |

這種設計方式讓 plugin：

*   可以聲明「我要把元件掛在某個指定位置」
    
*   可以提供可被其他 plugin 呼叫的 UI 元件或函式
    
*   可以在 UI 菜單中插入自己的設定連結或頁面
    

* * *

🧠 為什麼封裝成 `pkg/extensions` 模組？
------------------------------

Grafana 必須要：

*   **將 plugin.json 中的 `extensions` 轉換為嚴格定義的 struct（`ExtensionsV2`）**
    
*   **處理不同版本格式（V1/V2）的轉換與相容性處理**
    
*   **提供工具與 interface 給前端與 runtime 做 UI 組合與渲染**
    

這些需求屬於資料描述與解析邏輯，因此獨立為 `pkg/extensions`，而非混入 plugin runtime。

* * *

✅ 整體意圖：**平台化生態系擴展**
-------------------

Grafana 正在從「單一儀表板工具」轉型為「平台型生態」，而 `extensions` 就是這個方向的基礎：

*   前端：提供 plugin 注入點與視覺化組合
    
*   後端：可搭配 `apiserver` + `registry` 自訂資源物件
    
*   Plugin 作者：能做完整 UI + API 的擴展，而不 fork Grafana
    

* * *

如果你想在 detectviz 中也支援這類 extension 模式，我可以幫你設計一個 Plugin Interface 規格，例如：

```json
{
  "extensionPoints": ["alert-summary", "custom-dashboard-toolbar"],
  "addedComponents": [{ "type": "button", "target": "alert-toolbar", "id": "csv-exporter" }]
}
```

要幫你草擬這套格式與對應實作架構嗎？

# pkg/generated:
`pkg/generated` 是 Grafana 專案中用來儲存「**由程式自動產生的程式碼與結構定義**」的目錄。這些內容通常**不會手動編輯**，而是透過工具、模板或語言轉換器自動生成，用於支援 OpenAPI、Protobuf、Schema、Mock 等跨模組整合工作。

* * *

✅ 功能與用途
-------

| 功能類型 | 說明 |
| --- | --- |
| OpenAPI schema | 根據 `pkg/apis/...` 下的資源，產生可供前端或 Swagger 使用的 OpenAPI 定義 |
| Protobuf / gRPC | 使用 `protoc` 或 `buf` 產生 `.pb.go`, `.grpc.pb.go` 等 gRPC 服務定義 |
| Mock 測試碼 | 使用 `GoMock` 自動產生 interface 的測試替身（例如：`mock_xxx.go`） |
| 資源 Schema | 結合 CUE/CUETSY 或 codegen 工具，轉換為 Golang struct 或驗證邏輯 |
| 版本控制結構 | 產生 `buf.gen.yaml`、`buf.lock` 等版本相容標記，用於多人協作時的 schema 一致性檢查 |

* * *

🧩 與其他模組的關係
-----------

| 關聯模組 | 說明 |
| --- | --- |
| `pkg/services/apiserver` | 引用由 `pkg/generated/openapi` 產生的 schema 定義並註冊到 K8s-style server 中grafana-pkg-services-al… |
| `pkg/registry/apis/...` | 使用 `GetOpenAPIDefinitions()` 提供的結構整合進版本化 REST API |
| `pkg/plugins/backendplugin/...` | 使用 `sanitizer.pb.go` 等 proto 文件支援 plugin 與 host 間的 gRPC 通訊grafana-pkg-plugins-code |

* * *

🛠 來源與產生方式
----------

| 工具 | 用途 |
| --- | --- |
| `protoc` + `protoc-gen-go` | 產生 gRPC / Proto 結構檔（例如 plugin 通訊） |
| `cue` / `cuetsy` / `codejen` | 產生 schema 定義與驗證邏輯 |
| `GoMock` / `mockgen` | 產生測試 stub 用的 interface 實作 |
| `make gen-*` | 在 Grafana 中會有 `make gen-cue`, `make gen-openapi` 等指令來觸發生成流程 |

* * *

🧠 為何要獨立為 `pkg/generated`
-------------------------

*   ✅ 集中管理所有機器生成的代碼，便於版本控管與 CI/CD 校驗。
    
*   ✅ 明確標示「**請勿手動編輯**」，避免不小心修改掉工具產出的內容。
    
*   ✅ 支援可組合模組化系統（如：aggregator、registry、plugin API）以動態擴展。
    

* * *

是否需要我幫你列出「detectviz 若要整合 codegen 機制，應該如何設計 `pkg/generated` 的目錄與產出流程」？可以搭配 `cue` 或 `protobuf`。

# pkg/ifaces/gcsifaces:
`pkg/ifaces/gcsifaces/gcsifaces.go` 會寫在這個位置，是基於以下幾個架構上的考量與慣例：

* * *

✅ 功能說明：**集中管理跨模組可重用的 interface 定義**
-----------------------------------

### `pkg/ifaces` 的角色：

*   這是一個 **interface 專屬目錄**，用來定義與其他模組溝通所需的抽象接口（而非實作邏輯）。
    
*   它與 `pkg/services` 的分工為：**services 放實作，ifaces 放介面**。
    
*   在大型模組（如 plugins、cloud、storage）之間，**抽出 interface 可避免循環依賴（import cycle）**。
    

* * *

🧱 `gcsifaces` 的用途
------------------

`gcsifaces` 主要與「GCS-like 儲存介面」有關，可能是：

*   提供一組 interface 用來抽象對 Google Cloud Storage（或相容的 S3 儲存系統）的操作行為。
    
*   被 Grafana 的 `cloudmigration`、`provisioning`、或 plugin runtime 所共用。
    

範例可能的 interface（推測性）：

```go
type ObjectStorage interface {
	PutObject(ctx context.Context, bucket, key string, data []byte) error
	GetObject(ctx context.Context, bucket, key string) ([]byte, error)
	DeleteObject(ctx context.Context, bucket, key string) error
}
```

* * *

💡 為什麼不放在 `pkg/services/cloudmigration`？
----------------------------------------

因為這些 interface：

*   **可能不屬於 cloudmigration 專用**
    
*   也可能會被其他模組共用，如 `pkg/services/pluginstore`, `pkg/storage`, `pkg/supportbundles`
    
*   放在 `pkg/services` 會導致邏輯混亂與高耦合
    
*   放在 `pkg/infra` 過於底層，不適合 interface 說明
    

所以選擇將其 **獨立放入 `pkg/ifaces/gcsifaces`**，形成乾淨的抽象層。

* * *

📦 類似慣例你也可以套用在 detectviz
------------------------

| 位置建議 | 用途 |
| --- | --- |
| `pkg/ifaces/datastoreiface` | 定義 DB 抽象，如 `GetMetric(tag)` |
| `pkg/ifaces/notifieriface` | 定義通知抽象，如 `Notify(ctx)` |
| `pkg/ifaces/llmiface` | 定義 LLM API 抽象，方便替換 backend |

這樣你的模組可以輕鬆切換實作，不產生 import loop，也更適合寫測試與 Mock。

* * *

是否需要我幫你為 `detectviz/pkg/ifaces/` 起草一份標準目錄與 interface 命名規範？

# pkg/kinds:
`pkg/kinds` 在 Grafana 中是用來**統一定義與管理資源類型（Kind）與其對應的結構與驗證格式（通常為 CUE 檔案）**。這個模組與 Kubernetes 的 `Kind` 概念類似，目的是建立一個「可組合（Composable）與可驗證（Validatable）」的資源系統，支援動態註冊、自動驗證與 API 對應。

* * *

✅ 功能說明
------

| 目的 | 說明 |
| --- | --- |
| 定義 Grafana 內部「資源種類」 | 如：dashboard、librarypanel、preferences、datasource 等 |
| 掛載對應 CUE Schema 驗證檔 | 每種 Kind 都會對應一個 `.cue` 格式的結構描述檔 |
| 作為 registry 系統與 plugin schema 的標準化入口 | 被 `pkg/registry/schemas` 與 plugin codegen 引用 |

* * *

📦 目錄內容與格式（推論自 `core_kind.go`、`composable_kind.go`）
---------------------------------------------------

```go
type CoreKind struct {
  Name    string      // 如 "dashboard"
  CueFile cue.Value   // 對應的 CUE 結構
}

type ComposableKind struct {
  Name     string      // 如 "stat", "barchart"
  Filename string      // cue 定義檔名，如 "panelcfg.cue"
  CueFile  cue.Value
}
```

這些 Kind 被 `registry` 或 `plugin loader` 掃描並收錄進動態 API 註冊中。

* * *

🧩 使用場景
-------

*   **apiserver 與 aggregator 的 schema 安裝**
    
*   **registry 自動載入與驗證**
    
*   **plugin developer 註冊自己的 panel/datasource 定義**
    
*   **`make gen-cue` 自動產生所有可組合 Kind 的映射**
    

* * *

🛠 與 CUE 的整合
------------

Grafana 採用 CUE 作為強類型 schema 的來源，透過 `loadCueFileWithCommon(...)` 將這些 `.cue` 檔案載入為 `cue.Value`，並與資源物件做驗證與轉換grafana-pkg-registry-co…grafana-pkg-registry-co…。

* * *

🧠 延伸理解：為何放在 `pkg/kinds`
------------------------

*   `Kind` 是資源的 **語意單元與抽象層**
    
*   各模組（panel, datasource, dashboard）都會定義自己的 kind
    
*   將其獨立出來可：
    
    *   支援 Plugin 動態擴展
        
    *   避免耦合到 service 或 plugin 邏輯
        
    *   統一集中 schema 掃描、驗證與註冊流程
        

* * *

是否需要我幫你畫出「Kind 與 CUE 結構整合流程圖」或幫 detectviz 設計自己的 `pkg/kinds` 目錄與結構？這將有助於你設計標準化資源模型。


# pkg/ruleguard.rules.go:
`pkg/ruleguard.rules.go` 被放在這個位置，具有以下幾個關鍵考量：

* * *

✅ 檔案用途：定義 **Ruleguard 靜態分析規則**
------------------------------

*   `ruleguard.rules.go` 是 [Go Ruleguard](https://github.com/quasilyte/go-ruleguard) 的 DSL 格式檔案，通常用於：
    
    *   靜態程式碼分析（類似 `golangci-lint`）
        
    *   偵測 anti-pattern（如過時 API、禁止使用的呼叫）
        
    *   團隊內部開發規範檢查
        

* * *

📦 為何會出現在 `pkg/` 根目錄下？
----------------------

### 1. **讓規則與內部程式碼結構共存，可視範圍一致**

*   規則多半針對 `pkg/` 下模組撰寫，如：`pkg/services/*`、`pkg/infra/*`
    
*   放在同層可讓規則檢查路徑涵蓋整個 `pkg/`，不易遺漏
    
*   不混入 `tools/`、`build/` 等不執行的開發工具區
    

### 2. **這不是執行程式碼，不應進入某個模組目錄**

*   放進 `pkg/tools/` 會讓人誤解為工具邏輯
    
*   放進 `pkg/services/` 等則與實際功能耦合，違反 SRP
    
*   放在 `pkg/ruleguard.rules.go` 單獨存在，**明示這是全域靜態規則**
    

### 3. **CI / makefile 可以統一掃描與執行**

例如：

```bash
make lint-ruleguard
go run ./tools/ruleguard -rules=./pkg/ruleguard.rules.go ./pkg/...
```

* * *

🧠 與 `golangci-lint` 或 `go vet` 的差異
-----------------------------------

| 工具 | 用途 |
| --- | --- |
| `golangci-lint` | 聚合多種 linters，但自訂性有限 |
| `go vet` | 針對明顯錯誤與語意上的 bug 檢查 |
| **`ruleguard`** | 可用 DSL 自訂規則，精準控管風格與危險用法 |

* * *

🧩 若你在 `detectviz` 要採用類似靜態規範檢查
------------------------------

建議目錄規劃如下：

```bash
/pkg/
  ruleguard.rules.go      # 自訂靜態分析規則
  lint_patterns.md        # 搭配規則的人類可讀文件
/tools/
  ruleguard-runner/       # 若有自訂執行包裝器
```

是否需要我幫你定義 detectviz 的 lint DSL 起手式，例如禁止 `log.Fatal`、封鎖 `fmt.Println`？



# pkg/promlib:
`pkg/promlib` 是 Grafana 中為支援 Prometheus 相容功能所設計的一個模組，屬於 **內部支援 Prometheus TSDB 協定與模型操作的輔助庫**，其功能介於 Prometheus 原生模組與 Grafana plugin 之間，屬於「通用資料結構與轉換層」。

* * *

✅ 功能總覽
------

| 類別 | 說明 |
| --- | --- |
| Prometheus 資料模型 | 提供 Prometheus Label、TimeSeries、Query 的封裝與轉換（通常對應 `prompb` 結構） |
| gRPC 與 remote write 支援 | 協助將 Grafana 的數據（如 Frame）轉換為 Prometheus Remote Write 格式（`.prompb.WriteRequest`） |
| 序列化工具 | 封裝 proto + snappy 壓縮邏輯，用於 Prometheus 相容寫入 |
| plugin 整合中介 | 被 registry、plugins、tsdb 或 `plugin-sdk-go` 所引用，作為 Prometheus 資料互通格式 |

* * *

🔍 與其他模組的關聯
-----------

| 模組 | 關聯說明 |
| --- | --- |
| `pkg/services/ngalert` | 使用 promlib 來封裝告警資料結構與 Remote Write 格式grafana-pkg-services-ng… |
| `pkg/pluginsintegration` | 將 Plugin 回傳的資料轉換為 Prometheus Frame 結構 |
| `pkg/registry/apis/query/queryschema` | 使用 promlib 提供的 Label 與 TimeSeries 格式轉成 CUE schema 驗證 |
| `prompb`（prometheus proto） | promlib 封裝此格式，並透過 `TimeSeriesToBytes()` 提供序列化能力grafana-pkg-services-li… |

* * *

📦 常見函式與元件（推測）
--------------

*   `TimeSeriesFromFrames(*data.Frame) []prompb.TimeSeries`
    
    *   將 Grafana Frame（SDK 結構）轉換為 Prometheus 的時序資料結構。
        
*   `TimeSeriesToBytes([]prompb.TimeSeries) ([]byte, error)`
    
    *   使用 `proto.Marshal + snappy.Encode` 序列化為 remote write 資料。
        
*   `sanitizeMetricName`, `sanitizeLabelName`
    
    *   確保資料符合 Prometheus 命名規則。
        

* * *

🧠 為何獨立為 `promlib`？
-------------------

因為：

*   它不是 plugin 本身，也不是某個 service 的邏輯，屬於 **中立工具層**。
    
*   可被各種模組重用（如 `live`, `ngalert`, `pluginsintegration`）
    
*   可作為「Prometheus compatibility layer」，不與原始 Prometheus 模組耦合
    

* * *

🧩 適合你參考的應用場景（detectviz）
------------------------

若你要支援：

*   Prometheus Remote Write 接收或轉發
    
*   Frame → TimeSeries → Alert 的轉換
    
*   自訂 plugin 或通知器回傳 Prometheus 格式
    

你可以仿照 `promlib`：

*   設計 `pkg/promlib/frame.go`, `labels.go`, `remote.go`
    
*   將轉換、格式驗證、序列化集中處理
    
*   避免每個模組重複處理 Prometheus 相容邏輯
    

需要我幫你草擬 detectviz 專用的 `promlib` mini 版結構與範本嗎？


## Response:
`pkg/semconv` 是 Grafana 用來集中定義 **OpenTelemetry（OTel）Semantic Conventions** 的模組，功能類似於 `go.opentelemetry.io/otel/semconv` 的封裝與擴充，用於統一設定 **tracing span 中的語意欄位（如 service.name、http.status_code、db.system 等）**。

* * *

✅ 主要功能與用途
---------

| 功能分類 | 說明 |
| --- | --- |
| 統一語意欄位（Semantic Attributes） | 封裝 OTel 定義的欄位，如 `ServiceNameKey`, `HTTPMethodKey`, `DBSystemKey` |
| Schema URL 管理 | 統一 `semconv.SchemaURL` 的版本使用（如 `v1.17.0`） |
| 延伸定義 | 如需自定義欄位（Grafana 特有屬性），可透過此模組加以擴充 |
| 提供 tracing 初始化使用 | 在 `infra/tracing`、`services/live`、`pipeline` 中設定資源時引用grafana-pkg-services-li… |

* * *

🧩 使用場景
-------

*   在初始化 `OpenTelemetry` TracerProvider 時：
    
    ```go
    resource.NewWithAttributes(
      semconv.SchemaURL,
      semconv.ServiceNameKey.String("grafana"),
      attribute.String("environment", "dev"),
      ...
    )
    ```
    
*   中央統一管理語意欄位，避免硬編字串（如 `"service.name"`）分散各模組。
    
*   使用於：
    
    *   `pkg/infra/tracing`
        
    *   `pkg/services/live/pipeline`
        
    *   `pkg/httpclientprovider` 的 tracing middlewaregrafana-pkg-infra-code
        

* * *

🧠 為什麼放在 `pkg/semconv`
----------------------

*   **定位明確**：語意欄位是跨模組共用的，不屬於特定 service
    
*   **避免耦合**：不與 `infra/` 或 `tracing/` 綁定，便於在其他 plugin、middleware 使用
    
*   **與 OpenTelemetry 對齊**：但允許 Grafana 自定義與增補
    

* * *

💡 detectviz 參考建議
-----------------

你可以仿照此結構，在自己的架構中建立：

```go
/pkg/semconv/
  ├── service.go      // ServiceNameKey, ServiceVersionKey
  ├── alert.go        // AlertRuleIDKey, SeverityKey
  └── schema.go       // SchemaURL 常數
```

這樣能統一告警、監控、通知中 tracing span 的 metadata，一致又可維護。

是否需要我幫你產出 detectviz 版本的 `semconv/service.go` 與初始化 tracing 時的範本？






================================================
FILE: grafana-repo-notes/grafana-pkg-plugins.md
================================================

# Grafana plugin 的代碼邏輯:
Grafana plugin 的代碼邏輯分散在數個主要模組，以下是整理後的清單與說明，涵蓋 **載入流程、設定管理、儀表板整合與介面定義**：

* * *

✅ 1. `pkg/plugins/`：Plugin 核心與後端管理元件
------------------------------------

*   **功能**：處理 Plugin 掃描、驗證、啟動、關閉等生命週期。
    
*   **結構重點**：
    
    *   `manager/loader`：讀取 `plugin.json`、初始化 plugins。
        
    *   `manager/pipeline`：多階段載入流程（Discovery → Bootstrap → Validation → Initialization → Termination）grafana-pkg-plugins-code。
        
    *   `pluginscdn/`：CDN 載入支援。
        
    *   `repo/`：Plugin 資訊來源（本地、遠端等）。
        
    *   `config/`：plugin.json 結構定義與版本檢查。
        

* * *

✅ 2. `pkg/services/pluginsintegration/`：Plugin 整合層（統一對外服務）
----------------------------------------------------------

*   **功能**：建構 Plugin 的注入、檢查、授權、設定儲存、API 整合。
    
*   **主要子模組**：
    
    *   `pluginstore/`：所有 plugins 的查詢與註冊來源grafana-pkg-services-pl…。
        
    *   `pluginsettings/`：管理 plugin 的設定資料，提供 CRUD、快取、解密等功能grafana-pkg-services-pl…。
        
        *   `Service` interface 定義方法如：
            
            ```go
            GetPluginSettings(ctx, args)
            UpdatePluginSetting(ctx, args)
            ```
            
    *   `pluginassets/`：提供 plugin 資源與 hash 驗證（SRI check）grafana-pkg-services-pl…。
        
    *   `pluginchecker/`、`plugincontext/`：處理插件驗證與執行環境相關的 context。
        

* * *

✅ 3. `pkg/services/plugindashboards/`：Plugin 提供的 Dashboard 管理
-------------------------------------------------------------

*   **功能**：管理 plugin 附帶的 dashboard 模板，並能自動同步更新。
    
*   **主要介面**：
    
    ```go
    type Service interface {
      ListPluginDashboards(ctx, req)
      LoadPluginDashboard(ctx, req)
    }
    ```
    
*   **執行流程**：
    
    *   使用 pluginStore 查詢 plugin json 內的 dashboards。
        
    *   匹配舊有 dashboard 並根據修訂版本進行匯入grafana-pkg-services-al…。
        

* * *

✅ 4. Plugin 載入與註冊邏輯（Pipeline）
-----------------------------

*   **核心位置**：`pkg/plugins/manager/pipeline/`
    
*   **階段說明**：
    
    *   **Discovery**：從 disk / repo 中找到 plugins。
        
    *   **Bootstrap**：解析 plugin.json，建立 metadata。
        
    *   **Validation**：確認簽章、Angular 支援等。
        
    *   **Initialization**：啟動後端服務或註冊 API。
        
    *   **Termination**：關閉 plugin 進程與釋放資源grafana-pkg-plugins-code。
        

* * *

🧩 插件設定與狀態處理流程
--------------

1.  啟動時由 `pluginstore` 掃描 plugin。
    
2.  每個 Org 有獨立的 plugin 設定（`pluginsettings`）。
    
3.  插件變更事件（如啟用/停用）會觸發 `DashboardUpdater` 更新 dashboardsgrafana-pkg-services-al…。
    
4.  Grafana 提供 REST API `/api/plugins/:pluginId/dashboards` 供前端查詢。
    

* * *

需要的話，我可以幫你用 XMind 格式畫出 plugin 架構邏輯或流程圖，是否要一併生成？

## Response:
是的，**Grafana Plugin 的實作代碼通常是在開發者自己的 Repo 或專案中**，而 Grafana 本體只負責：

1.  **掃描與註冊 plugin**（透過 `plugin.json`）
    
2.  **執行 backend plugin binary**（如需要）
    
3.  **加載 plugin 提供的前端資源與設定頁**
    
4.  **透過 `pluginsettings` 管理使用者設定**
    

* * *

✅ Plugin 開發者通常提供的內容：
--------------------

| 類別 | 說明 |
| --- | --- |
| `plugin.json` | Plugin 的 manifest，描述類型、ID、版本、前端資源、backend binary 等 |
| `/dist/` 或 `/public/` | 編譯後的前端資源（React、Angular） |
| `backend plugin` (選擇性) | 使用 Go 撰寫的 binary，提供資料存取、運算等能力 |
| 自有 repo | Plugin 原始碼通常放在自己的 GitHub Repo，例如：  
[https://github.com/grafana/google-sheets-datasource](https://github.com/grafana/google-sheets-datasource) |

* * *

🧩 Grafana 如何載入這些外部插件：
----------------------

1.  **開發者將 Plugin 安裝到 `data/plugins/` 資料夾**
    
2.  Grafana 在啟動時由 `pkg/plugins/manager` 掃描所有 `plugin.json`
    
3.  若 plugin 有 `backend binary`（在 `executable` 欄位定義），Grafana 會執行該檔案並透過 GRPC 通訊
    
4.  Plugin 可定義前端頁面（如設定 UI）、REST API proxy、預設 dashboard 等
    

* * *

🌐 若是官方支援的 Plugin
-----------------

*   部分 plugin 放在 Grafana 官方倉庫（如 Loki、InfluxDB）
    
*   官方也會將熱門 Plugin 放到 https://grafana.com/grafana/plugins
    

* * *

📦 安裝方式（使用者端）
-------------

*   使用 `grafana-cli plugins install <plugin-id>`
    
*   或者手動放入 `/data/plugins/your-plugin/` 並重啟 Grafana
    

* * *

如果你有考慮寫自己的 Plugin，我可以提供開發範例或模板（前後端皆可），是否需要？







================================================
FILE: grafana-repo-notes/grafana-pkg-registry.md
================================================
# pkg/registry
`pkg/registry` 是 Grafana 架構中用來實作 **Kubernetes-style API 資源註冊與管理系統** 的核心模組。它是 Grafana 向平台化（Platform-as-a-Framework）邁進的重要基礎，讓不同模組（如 dashboard、datasource、secret）以資源（Resource）的形式註冊、暴露、授權與操作。

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

🧩 架構與關聯模組
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

*   `APIGroupBuilder`：每個資源註冊單位，提供 GroupVersion、Authorizer、OpenAPI 定義等
    
*   `rest.Storage` / `rest.Connecter`：定義 REST API 行為（如 Get/List/Post/Patch/Connect）
    
*   `ExtraBuilder`：提供 mutate、OpenAPI 擴充、job 註冊等延伸邏輯grafana-pkg-registry-co…
    
*   `GenericRegistryStore`：包裝資源資料表與策略處理邏輯，與 `apiserver` 整合
    

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

*   資源模組的**統一註冊入口**
    
*   API group 的 **自動掛載與生命週期管理**
    
*   使用者自定資源（如 keeper, snapshot, serviceaccount）的**一致性 CRUD 與驗證**
    

* * *

若你要在 detectviz 支援「自訂資源註冊」與「YAML 驗證」，這個模組是非常值得借鏡的。我可以幫你設計 detectviz 專用的 `registry/` 架構與 minimal APIBuilder 實作，要嗎？


# 遷移到 Detectviz internal/registry/ 後怎麼做

根據你目前的架構（Detectviz 與 Grafana 拆分參照），pkg/registry 屬於 專用平台框架下的資源生命週期管理模組，不屬於通用可重用工具，因此建議：

⸻

✅ 遷移位置建議：internal/registry/

📌 為什麼不是 pkg/？
	•	這是一套 自有資源管理機制（如 GVK 註冊、OpenAPI schema、生存期管理），屬於平台的 control-plane 架構
	•	不可能被第三方 app 或 plugin 拿去共用
	•	高耦合於自家 apiserver、身份驗證、CUE 驗證、OpenAPI generator 等模組

⸻

🧱 建議遷移結構（Detectviz 版本）

internal/registry/
├── apis/
│   ├── secret/           # 自訂資源類型
│   ├── iam/
│   └── datasource/
├── apps/                 # 整合 app plugin 對外註冊
├── contracts/            # 抽象 interface（如 Store, Resolver）
├── builder.go            # APIGroupBuilder
├── register.go           # API 註冊總入口
├── cueutils.go           # CUE schema 驗證
├── decoder.go            # YAML / JSON Decode

若你未來希望支持像 Grafana 一樣「自定義資源物件 + REST + YAML 驗證 + 動態註冊」，這樣的分法可持續擴展

⸻

✅ 對應 Interface 若需抽象化
	•	可選擇把：
	•	ObjectStore, GVKResolver 等 interface 移至 pkg/iface/registry.go
	•	但保留實作仍放在 internal/registry

這樣外部 plugin 若需實作，也可依 interface 注入。

⸻

是否需要我幫你出一份對 Detectviz 量身訂做的 internal/registry 最小實作範例？或者草擬 GVK + APIBuilder 註冊邏輯？


================================================
FILE: grafana-repo-notes/grafana-pkg-server.md
================================================
# pkg/server 概覽

`pkg/server` 目錄是 Grafana 啟動與模組管理的核心邏輯之一，負責主程式啟動、模組註冊、HTTP/GRPC 伺服器初始化、與模組服務（如 distributor、storage、frontend）的協調與執行。以下是其主要功能說明：

* * *

📦 `pkg/server` 功能說明
--------------------

### 1\. **Server 啟動邏輯（`server.go`）**

*   提供 `New()` 與 `Run()` 方法建立與執行核心 `Server` 實例。
    
*   管理背景服務的啟動與關閉。
    
*   負責初始化 log、設定檔、PID 寫入與 systemd 通知等grafana-pkg-all-code。
    

### 2\. **模組伺服器（`module_server.go`）**

*   `ModuleServer` 是為了模組化 Grafana 運行的入口點設計。
    
*   支援 `init()` 啟動各個模組：如 `initDistributor()`、`initFrontend()`、`initInstrumentationServer()` 等。
    
*   用於建構 dskit-based 模組並管理模組依賴關係grafana-pkg-all-code。
    

### 3\. **模組註冊與服務管理（與 `pkg/modules` 整合）**

*   使用 `modules.Engine` 與 `services.BasicService` 管理模組生命週期。
    
*   各模組如 `Distributor`, `StorageServer` 等會用 `WithName()` 註冊名稱，方便統一管理與診斷grafana-pkg-all-code。
    

### 4\. **Instrumentation 伺服器（`instrumentation_service.go`）**

*   啟動 Prometheus metrics HTTP server，對外提供 `/metrics` endpoint。
    
*   也可設定為健康檢查與 profiling 伺服器（透過 `gorilla/mux` 實作）grafana-pkg-all-code。
    

### 5\. **模組型服務定義（如 `distributor.go`, `memberlist.go`）**

*   每個模組會實作自己的 `initXXX()` 函式，返回一個符合 `services.Service` 介面的執行單元。
    
*   e.g., `initDistributor()` 初始化一個 grpc-based 分發模組服務grafana-pkg-all-code。
    

### 6\. **Runner 與 ModuleRunner（`runner.go`, `module_runner.go`）**

*   `Runner` 是整體伺服器的執行者，注入 config、密鑰、user service 等元件。
    
*   `ModuleRunner` 是簡化版，只載入模組需要的最小依賴（例如 feature toggles）grafana-pkg-all-code。
    

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

依照你在 Detectviz 的分層與架構規劃，pkg/server 是「平台啟動控制核心」，負責主流程組裝、模組註冊與伺服器執行邏輯，並不應視為通用可重用套件。因此，建議遷移至以下位置：

# 遷移到 Detectviz internal/bootstrap/server/ 後怎麼做


📌 遷移理由：

項目	說明
📦 高度耦合 Detectviz 啟動流程	如模組註冊、metrics exporter、服務註冊皆與內部模組密切相關
❌ 不具重用性	無 plugin 化可能、不可被外部 app 重用
✅ Clean Architecture 原則	將初始化組裝邏輯集中於 internal/bootstrap 或 internal/platform 是最佳實踐
🔄 與 internal/infra/, internal/registry/, internal/auth/ 等模組配合頻繁	放於同層級更能統一服務架構與生命週期流程


⸻

🧱 遷移後建議結構

internal/bootstrap/
├── server/
│   ├── server.go
│   ├── module_server.go
│   ├── instrumentation_service.go
│   ├── runner.go
│   ├── module_runner.go
│   ├── distributor.go
│   └── memberlist.go

若你將 modules/, services/ 一併重構，可一併放入 internal/bootstrap/modules/、internal/bootstrap/services/

⸻

✅ 延伸重構建議
	•	將 server.go 內建構的 Server 結構設計為依賴注入容器，載入：
	•	plugin registry
	•	event dispatcher
	•	API router
	•	ModuleRunner 可定義 interface 轉為 bootstrap.Runner，支援 app 啟動入口統一化

⸻

是否需要我幫你規劃 bootstrap/server 搭配 apps/main.go 的組裝流程、或重構後的模組生命週期圖？也可以產出 Mermaid 架構圖來輔助說明。


================================================
FILE: grafana-repo-notes/grafana-pkg-setting.md
================================================
# pkg/setting:
`pkg/setting` 是 Grafana 中負責 **統一管理組態設定（設定檔解析、結構化欄位、環境變數覆寫、預設值與驗證）** 的模組。它是整個 Grafana 啟動過程的設定核心，提供其他模組依賴的 `Cfg` 結構與各項設定細節。

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

🧩 子模組說明（部分）
------------

| 子檔案 | 功能 |
| --- | --- |
| `setting_plugins.go` | Plugin 安裝與更新策略設定grafana-pkg-all-code |
| `setting_remote_cache.go` | Redis 等快取設定（加密、prefix）grafana-pkg-all-code |
| `setting_unified_storage.go` | Unified Storage（新儲存引擎）細節控制，如 dual writer、shardinggrafana-pkg-all-code |
| `setting_anonymous.go` | 匿名登入設定，如 org_role、限制功能grafana-pkg-all-code |
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

Grafana 採用 **設定中心 + 結構對映 + 動態解析** 的設計，使得：

*   設定變更更易控管與落地
    
*   對應模組可以直接取結構，不需要再自行解析 ini
    
*   可支援未來動態設定儲存（如 cloud 控制台、GUI 編輯）
    

* * *

如你在 `detectviz` 也需要支援 `.ini`、`.env` 或動態設定系統，這個模組是非常適合參考的。我可以幫你簡化為一套 `setting` 框架版本，是否要我整理一份範本？

# 遷移到 Detectviz pkg/config/ 後怎麼做

✅ 是的，將 pkg/setting 改名為 pkg/config 是合理且推薦的做法，尤其是在你希望更貼近 Clean Architecture 與更語意化命名時。以下是具體說明：

⸻

✅ 改名為 pkg/config 的理由

原名稱 setting	建議名稱 config	說明
與變數 setting.Cfg 容易混淆	config.Cfg 意圖更明確	更清楚表示這是「設定框架」而非單一變數用途
setting_*.go 命名較模糊	config_*.go 更通用且業界慣例	例如：config_plugins.go, config_auth.go
擴充上不自然	setting/expanders.go	放在 config/expanders.go 更合邏輯
預期可接 .ini, .env, .yaml 等格式	config 命名更涵蓋多樣設定來源	有助未來導入動態設定、GUI 編輯等機制


⸻

✅ Detectviz 架構建議路徑

pkg/
└── config/
    ├── config.go              # 定義 Cfg 結構與載入流程
    ├── config_plugins.go
    ├── config_alerting.go
    ├── config_storage.go
    ├── expanders.go

你也可以依照模組維度分拆（可選）：

pkg/
├── config/                   # 設定核心與組合器
├── config/plugins/          # plugins 專屬設定
├── config/alerting/         # alerting 專屬設定


⸻

✅ 是否需要我幫你：
	•	重構 import（setting ➜ config）
	•	改寫 setting.NewCfg() ➜ config.Load() 風格
	•	或生成 .env + .ini + GUI 設定支援架構範本？

都可以即時協助。


================================================
FILE: grafana-repo-notes/grafana-pkg-tests.md
================================================
# pkg/tests 與 pkg/mocks

# pkg/tests:
`pkg/tests` 是 Grafana 原始碼中的 **集中整合測試模組（integration test harness）**，目的在於提供跨模組測試與共享測試邏輯，與 `pkg/mocks` 和 `pkg/util/testutil` 搭配運作，覆蓋底層儲存行為、資源 API 行為、查詢行為等核心流程。

* * *

✅ 功能與用途
-------

| 類型 | 說明 |
| --- | --- |
| 整合測試套件 | 專門針對 `pkg/storage/unified`、`pkg/services/ngalert`、`pkg/registry` 等模組設計一系列跨模組整合測試 |
| 測試抽象封裝 | 提供類似 `RunStorageBackendTest()`、`RunSearchBackendTest()` 的標準化 test runner，傳入不同 backend 即可共用測試邏輯grafana-pkg-all-code |
| 自動命名空間隔離 | 使用 `GenerateRandomNSPrefix()` 為每次測試產生獨立 namespace，避免並行測試資料衝突 |
| 支援跳過測試項目 | `TestOptions.SkipTests["TestBlobSupport"] = true` 可自訂跳過指定測試項目grafana-pkg-all-code |

* * *

🧩 與其他測試模組的差異
-------------

| 模組 | 目的 | 範圍 | 特點 |
| --- | --- | --- | --- |
| `pkg/tests` | 整合測試 | 全系統級 | 跨模組、多元 backend 比較 |
| `pkg/mocks` | 單元測試 mock | 針對 interface | 通常搭配 `mockery` 自動產生 |
| `pkg/util/testutil` | 單元測試共用工具 | test context, 假帳號 | 輕量級、偏向 test helper |
| 各模組 `/testutil` | 模組專屬 test helper | 限模組內使用 | 如 `ngalert/testutil`, `foldertest` |

* * *

📦 使用範例：測試 Unified Storage backend
----------------------------------

```go
RunStorageBackendTest(t,
  func(ctx context.Context) resource.StorageBackend {
    return NewMyStorageBackend(...)
  },
  &TestOptions{NSPrefix: "demo-test"},
)
```

此機制可套用到：

*   自訂後端儲存實作驗證
    
*   plugin 或 fork 運作環境測試
    
*   Regression test 與可擴充性測試
    

* * *

✅ 為何放在 `pkg/tests`
------------------

*   ✅ 提供內部共用的測試邏輯（非模組特有）
    
*   ✅ 減少重複撰寫測試樣板
    
*   ✅ 方便 CI/CD 中集中呼叫與篩選
    
*   ❌ 不會被 production runtime 引用（清楚區隔）
    

* * *

若你要在 `detectviz` 中設計一套模組化儲存或 alert 測試，也可參考這種設計。需要我幫你定義 `pkg/tests/backend_test.go` 的範本嗎？適用於 plugin or unified backend 驗證。

# pkg/mocks:
`pkg/mocks` 是 Grafana 中用來放置 **mock 類型的集中目錄**，其主要目的是支援單元測試與介面替代，特別適用於 dependency injection 架構中的測試環境。

* * *

✅ 功能與用途
-------

| 類型 | 說明 |
| --- | --- |
| Interface Mock 實作 | 使用工具（如 `mockery` 或 `gomock`）自動產生介面實作檔案 |
| 測試替身 | 提供可以斷言呼叫次數、參數、行為的物件 |
| 減少循環依賴 | 測試階段引用 `mocks`，不會造成正式邏輯相依性 |

* * *

🧩 Mock 的來源與產生方式
----------------

常見產生方式：

```go
//go:generate mockery --name=SomeInterface --output=pkg/mocks --with-expecter
```

*   使用 `mockery` 工具根據介面自動產出對應 mock 檔案。
    
*   範例：`pkg/mocks/T.go` 是對 `testing.T` 風格的 interface mock 實作grafana-pkg-all-code。
    

* * *

🆚 `pkg/testutil` 與 `mocks` 的差異
-------------------------------

| 模組 | 用途 | 是否可程式產生 |
| --- | --- | --- |
| `testutil` | 提供測試工具與資料（如 test context、假帳號） | ❌ |
| `mocks` | 對 interface 的自動 mock 實作，配合 mock assert 機制 | ✅（mockery / gomock） |

* * *

🧠 為何放在 `pkg/mocks`
-------------------

1.  **集中管理 auto-generated mock**，避免與手寫測試邏輯混淆
    
2.  **在模組之上層級**：讓所有模組都能引用 mock（避免 import loop）
    
3.  **配合 make test 或 CI 流程整合**
    

* * *

📦 使用場景
-------

*   在單元測試中注入 mock：
    
    ```go
    ctrl := gomock.NewController(t)
    mockDB := mocks.NewMockDatabase(ctrl)
    service := NewMyService(mockDB)
    ```
    
*   搭配 `EXPECT()` 檢查預期呼叫：
    
    ```go
    mockDB.EXPECT().BeginTx(gomock.Any(), gomock.Any()).Return(...)
    ```
    

* * *

如果你要在 `detectviz` 導入類似架構，我可以幫你規劃：

*   `pkg/ifaces/` 放 interface
    
*   `pkg/mocks/` 放 mockery 產物
    
*   `tools/mockgen.sh` 做自動更新
    

# 遷移到 Detectviz pkg/tests 與 pkg/mocks 後怎麼做

以下是 pkg/tests 與 pkg/mocks 在 Detectviz 架構下的定位與建議遷移位置：

⸻

✅ pkg/tests（整合測試邏輯）建議遷移位置：

internal/testharness/

📌 理由：

面向	說明
🔧 功能定位	提供跨模組整合測試流程（如 alert + registry、storage + api）
📦 特性	包含測試 helper、test runner、命名空間管理、資源註冊測試
🔐 不應暴露	不屬於外部 library，應內部保留以供整合測試使用
☑️ 獨立可測試	可整合至 CI，對新 backend/plugin 實作驗證有幫助

若 Detectviz 有明確測試標準化需求（如 plugin compliance 測試），也可命名為 internal/compliance/ 或 internal/itest/

⸻

✅ pkg/mocks（interface mock 產物）建議遷移位置：

internal/test/mocks/

📌 理由：

面向	說明
📦 功能定位	自動生成的 interface 替身，用於單元測試注入
🧰 適用工具	mockery, gomock, moq 等生成工具產出
⛔ 不屬於 pkg	並非共用邏輯，不應開放給 production code 引用
✅ 清楚分層	配合 pkg/iface/ interface 管理，將 mock 分層放入測試區域最清晰


⸻

🧠 補充建議：

📁 最終建議結構：

pkg/
└── iface/
    ├── alert.go
    ├── registry.go
    └── ...

internal/
├── test/
│   └── mocks/             # interface mock
├── testharness/           # 整合測試（pkg/tests）
│   └── alert_storage_test.go


⸻

✅ 是否可加上 make mock 自動產生流程？

若你打算大量使用 mock，可以新增以下：
	•	tools/mockgen.sh：產出所有 mock（如 mockery）
	•	make mock：掃描 pkg/iface/ 並生成對應 mock 到 internal/test/mocks/

⸻

是否要我幫你產出 mockgen.sh 範例腳本與目錄結構？或者直接重命名與規劃你目前的 pkg/tests 轉成 internal/testharness/？


================================================
FILE: grafana-repo-notes/grafana-pkg-util.md
================================================
## Response:
`pkg/util` 是 Grafana 中負責提供 **小型、通用、無依賴工具函式與輔助模組** 的目錄，類似標準函式庫的補強工具。這些模組通常不屬於任何業務邏輯，只提供「低層抽象工具」，被廣泛重用於 `pkg/services`、`pkg/api`、`pkg/infra` 等模組中。

* * *

✅ 子模組功能與用途（依來源整理）
-----------------

| 子模組名稱 | 功能說明 |
| --- | --- |
| `cmd` | 執行 CLI command 或包裝 `exec.Command` |
| `cmputil` | 提供 systeminfo、CPU 核心、架構查詢等系統層工具 |
| `debouncer` | 實作 debounce 邏輯，避免短時間內重複觸發相同行為（如：view rebuild） |
| `errhttp` | 將錯誤轉換為 HTTP 錯誤格式（如 JSON）供 API 使用 |
| `httpclient` | 簡化 HTTP 請求操作（與 `pkg/infra/httpclient` 職責有部分重疊但更輕量） |
| `osutil` | 作業系統層工具，如檔案刪除、確認目錄是否存在 |
| `proxyutil` | 處理 HTTP Proxy header / remote addr 轉換工具 |
| `retryer` | 提供 retry 機制與 backoff 策略，常被資料寫入、API 傳送使用 |
| `ring` | 一種固定大小 queue buffer，常用於 metrics 或 stream 資料 buffer |
| `scheduler` | 排程工具，可做定期任務或延遲任務處理（類似 `cron` 簡化版） |
| `testutil` | 提供測試用工具（如產生假 user、token、context） |
| `ticker` | 定時器 helper，與 `time.Ticker` 差異是能明確控制起停 |
| `xorm` | 對 legacy ORM `xorm` 的輔助（如 SQL builder、date 過濾條件） |

* * *

🧩 調用對象範例
---------

| 使用模組 | 調用 `pkg/util` 功能 |
| --- | --- |
| `pkg/services/live` | 使用 `ring` 作為資料流記憶環 |
| `pkg/services/alerts` | 使用 `retryer` 發送告警結果到外部 API |
| `pkg/api/*` | 使用 `errhttp` 格式化錯誤輸出 |
| `pkg/login/*` | 使用 `httpclient` 呼叫外部身份認證 |
| `pkg/infra/log` | 使用 `cmputil` 取得系統資訊與版本報告 |
| `pkg/services/contexthandler` | 使用 `proxyutil` 處理 IP 與 header 資訊 |

* * *

✅ 與 `pkg/infra` 差異
------------------

| 模組 | 角色 | 內容風格 |
| --- | --- | --- |
| `pkg/util` | 工具箱 | 無依賴、stateless、小型 |
| `pkg/infra` | 平台設施 | 有狀態、有初始化（如 log, cache, tracing） |

例如：

*   `util/retryer` 是單純函式；但 `infra/db` 則會建立持久連線
    
*   `util/scheduler` 是簡化版排程；但 `infra/metrics` 會推送到 Prometheus
    

* * *

🧠 適用於你（detectviz）情境
--------------------

若你也想設計輕量、無依賴的共用工具，可以：

*   參考 `util/retryer`, `debouncer`, `ring` 的最小抽象設計
    
*   集中放在 `pkg/util/` 或 `internal/util/`，保留自訂性
    

是否需要我幫你整理 detectviz 可以沿用的 `util` 模組清單與模板？例如 `retryer`, `debouncer`, `testutil`？


================================================
FILE: grafana-repo-notes/grafana-pkg-web.md
================================================
## pkg/web:
`pkg/web` 是 Grafana 自行實作的 **輕量 Web 框架抽象層**，是對原本 [Macaron](https://go-macaron.com/) 框架的再封裝與擴展，主要目的是統一路由、middleware、context 與測試架構。此模組是 Grafana API 與前端溝通的最底層 HTTP 處理層。

* * *

✅ 功能總覽
------

| 類別 | 說明 |
| --- | --- |
| HTTP 路由封裝 | 提供 `Mux`（即 Macaron）作為 HTTP server 的 routing 介面grafana-pkg-all-code |
| Context 管理 | 將 `http.Request` 封裝為可註入 Grafana 內部 context 結構 |
| Middleware 支援 | 提供 `web.Handler`、自定義 `web.Middleware`、上下游控制 |
| Bind & Validation | 支援 JSON 解析並使用 `binding:"Required"` 驗證grafana-pkg-all-code |
| 測試支援 | `webtest/` 提供 handler chain 測試工具與模擬 HTTP servergrafana-pkg-all-code |

* * *

📦 模組結構說明
---------

| 檔案 | 功能 |
| --- | --- |
| `web.go` | 定義 `type Mux = Macaron`，作為主要 HTTP server 架構 |
| `context.go` | 封裝 request, response, params 等為統一 `web.Context` |
| `router.go`, `tree.go` | 路由匹配邏輯，支援 pattern route、params 抽取 |
| `binding.go` | 對 JSON request 做解析與欄位驗證 |
| `response_writer.go` | 攔截 `http.ResponseWriter` 提供 status code 追蹤與寫入控制 |
| `webtest/` | 單元與整合測試工具（e.g. middleware chain 模擬、fake server）grafana-pkg-all-code |

* * *

🧩 與其他模組的整合關係
-------------

| 模組 | 使用方式 |
| --- | --- |
| `pkg/api/` | 所有 REST API handler 都建立在 `web.Handler` 之上 |
| `pkg/middleware/` | 使用 `web.Middleware` 型別註冊（如 gzip、subpath redirect）grafana-pkg-all-code |
| `pkg/services/contexthandler` | context 擴充與身份注入（Org, User）搭配 web context 注入使用grafana-pkg-all-code |
| `pkg/api/routing` | 提供 `routing.Wrap(func(*ReqContext)) web.Handler` 將 API 封裝成 handlergrafana-pkg-all-code |

* * *

✅ 核心特色
------

*   🌐 輕量：不使用大型框架，如 gin 或 echo，而是自建 minimal wrapper
    
*   🧩 可測試：具備完整測試模擬能力
    
*   🧱 可插拔：支援 middleware chain，自訂上下文與回應處理
    
*   🔄 與 Grafana 自訂 context 模型緊密整合（如 `ReqContext`, `ContextModel`）
    

* * *

🧠 適合 detectviz 借鏡之處
--------------------

若你想要：

*   實作內建 API 架構（非用 gin/echo）
    
*   支援自訂 middleware chain
    
*   可單元測試 handler chain
    
*   整合自訂 context（如登入資訊、request metadata）
    

那麼可以仿照：

*   `pkg/web/context.go` 做 context 封裝
    
*   `binding.go` 做 JSON 結構驗證
    
*   `webtest/` 做 middleware chain 測試模組
    

是否需要我幫你規劃 detectviz 的 `web` 子系統結構？包含 `router.go`、`context.go`、`binding.go` 三件組合？

## Response:
使用 Echo 是**可以實作大多數 Grafana `pkg/web` 提供的功能**，但要看你的需求來決定是否「有必要」自己實作如 Grafana 那樣的框架包裝。下面是具體對比與可能無法直接實作的差異點。

* * *

✅ Echo 可以做到的功能（與 `pkg/web` 對應）
------------------------------

| 功能 | `pkg/web` | Echo 可做嗎 | 備註 |
| --- | --- | --- | --- |
| 路由註冊與群組化 | ✅ | ✅ | `Group` 支援 path prefix 與 middleware |
| Middleware 注入（Request Context） | ✅ | ✅ | `echo.MiddlewareFunc` |
| JSON Body Bind + 驗證 | ✅ | ✅ | 使用 `c.Bind(&obj)` + 手動驗證 |
| 自訂 context 物件封裝 | ✅ | ✅(進階) | 可用 `echo.Context` 包裝，需自己轉型 |
| Response 攔截（包裝 status code） | ✅ | ⛔\* | Echo 不提供內建 `ResponseWriter` 攔截器 |
| 測試用 mock context / fake server | ✅（webtest） | ⛔\* | Echo 需額外封裝 context 或用第三方測試工具 |

* * *

⛔ Echo 難以直接覆蓋的功能（Grafana 有但 Echo 要自己補）
--------------------------------------

| 功能 | 說明 |
| --- | --- |
| `web.Context` 封裝身份、logger、params、meta 統一物件 | Echo 只有 `Context`，但不會自動與 service/infra 整合（如 Grafana 的 `ReqContext`） |
| 狀態碼攔截與中止流程（如 early return） | Echo 雖能設定 status，但要「攔截所有中間層的 status 改寫」需手動包裝 `ResponseWriter` |
| 支援多層 middleware + 統一轉型流程 | Grafana 的 `routing.Wrap()` 可讓 API handler 統一接受 `*ReqContext`，這在 Echo 要靠中間層轉型 |

* * *

🧠 結論與建議
--------

*   若你是 \*\*平台級設計（如 detectviz）並需要：
    
    *   統一的自訂 context（含 org/user/logger）
        
    *   可擴展的 middleware chain
        
    *   可封裝的 handler（支援測試、自動轉型）
        
    
    👉 **可以參考 Grafana 的 `pkg/web` 架構，自建一層薄薄的包裝。**
    
*   若你是快速開發、功能導向  
    👉 **直接使用 Echo 是沒問題的，只要你知道你要補哪些功能。**
    

* * *

需要我幫你用 Echo 實作一份「仿 `pkg/web` handler + middleware + 自訂 context」的模板嗎？適用 detectviz 架構。


# 遷移到 Detectviz internal/web/ 後怎麼做

✅ 建議 先保留 pkg/web，理由如下：

⸻

✅ 保留 pkg/web 的情境與優勢

條件	說明
detectviz 使用 Echo，但有意日後封裝統一介面	你可日後自建 web.Context、handler chain、自動轉型等功能
部分 API handler 尚未轉移至新的 router 結構	可先沿用 Grafana 的 web.Handler 架構，平滑過渡
未來希望模擬 handler + middleware 做整合測試	webtest/ 仍是一個成熟又可移植的工具集
想觀察 web.Context 的實作概念以設計 ReqContext	保留原始碼方便學習與轉寫


⸻

🚧 後續你可能會這樣做：

時機	動作
定義好 detectviz/web.Context 與 routing.Wrap()	✅ 開始遷出 pkg/web/context.go
整合 echo context 擴充功能（e.g. logger, user）	✅ 建立 internal/web/context.go
測試新 middleware chain 流程	✅ 導入 webtest/ 或自行封裝 echo 測試工具


⸻

🧭 總結建議
	•	✔️ 短期：保留 pkg/web 作為「對照與過渡模組」
	•	🔜 中期：視你是否完成 detectviz 自建 web 套件再逐步遷出
	•	❌ 切勿直接砍掉 pkg/web，但可註記 deprecated 或 limited use

⸻

如你日後決定自建一層 internal/web/ 包裝 echo，我可以協助設計：
	•	Context 封裝（org, logger, trace ID）
	•	Wrap(func(*ReqContext)) echo.HandlerFunc
	•	middleware 注入鏈

是否需要我幫你草擬 detectviz 自建 internal/web 封裝架構？

✅ 是的，你可以現在就將 pkg/web/ 遷移至 internal/web/，這是自然且合乎 Detectviz 架構分層原則的重構動作。

⸻

✅ 為什麼應遷移至 internal/web/

原因	說明
pkg/web 是專屬於 detectviz 的 HTTP handler 抽象邏輯	無法重用於其他專案，不屬於通用庫
使用了自定義 context、routing、middleware chain	全部與 detectviz 平台高度耦合
Grafana 原始碼結構為專案內部 API 建構，並非 framework export	web.Context、binding、webtest 都非獨立套件


⸻

🧱 遷移建議結構

internal/web/
├── context.go         # 自訂 context 包含 request/user/logger
├── router.go          # 註冊與匹配邏輯
├── binding.go         # JSON bind 與驗證
├── response_writer.go # 攔截與回應控制
├── web.go             # 主入口：定義 router、middleware
└── webtest/           # 單元測試與 chain 模擬


⸻

✅ 遷移步驟建議
	1.	將整個 pkg/web/ 搬移至 internal/web/
	2.	所有原本 import "github.com/detectviz/pkg/web" 改為 "github.com/detectviz/internal/web"
	3.	如有 pkg/api/、pkg/services/ 引用 web.Handler 或 web.Middleware 也一併調整
	4.	可加上註解標記已內部化，例如：

// Deprecated: moved to internal/web. This package is no longer part of public interface.





================================================
FILE: grafana-repo-notes/services-featuremgmt.md
================================================
# pkg/services/featuremgmt:
以下是根據實際代碼整理的說明：

`pkg/services/featuremgmt` 功能說明
-------------------------------

此模組負責 Feature Toggle 管理與查詢，可視為 Grafana 的「功能開關中心」，包含以下重點元件：

- `manager.go`：定義 `FeatureManager` 結構，實作功能註冊（`registerFlags`）、更新（`update`）、是否啟用等判斷邏輯。
    
- `service.go`：提供 `ProvideManagerService` 給 DI 使用，初始化 FeatureManager，並從 `custom.ini` 的 `[feature_toggles]` 讀取設定。
    
- `goff_provider.go`：提供與 [GO Feature Flag](https://github.com/thomaspoignant/go-feature-flag) 整合的 `FeatureProvider`。
    
- `openfeature.go`：支援 [OpenFeature](https://openfeature.dev/) 架構，建立 Client 與 Provider 設定。
    
- `usage_stats.go`：產生目前開啟的 Feature 對應的 Prometheus metrics。
    

整體用途為：集中管理 Grafana 所有的實驗性、預覽與企業功能是否啟用，並可搭配前後端條件呈現、版本控制、Prometheus 指標等。