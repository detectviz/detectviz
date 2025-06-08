
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



