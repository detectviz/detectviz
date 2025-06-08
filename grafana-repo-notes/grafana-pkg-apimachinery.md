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