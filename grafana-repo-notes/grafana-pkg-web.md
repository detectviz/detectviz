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


