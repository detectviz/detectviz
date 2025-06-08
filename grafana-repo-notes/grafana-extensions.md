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

