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