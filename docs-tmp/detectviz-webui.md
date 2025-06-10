# 前端靜態資源模組

> **前端靜態資源、前端 API handler、或 HTMX / Web UI 的承載模組**。

> **不需要將 Web UI 獨立成一個前端 repo，完全可以和主後端（Detectviz）放在同一個 repo 內部，並直接由 app 管理它的畫面邏輯。**

### 前端技術棧

- 使用 Tabulator 呈現資料表
- 使用 HTMX 提供輕量互動式前端
- 使用 AdminLTE 提供管理介面
- 以 config 驅動頁面行為與 UI 顯現
- 遵循了「頁面即模組、元件可重用」的理念，每個頁面都是一個獨立模組，可直接載入。


建議方案（依照你目前架構現況）
-----------------

根據你先前描述：

- 你使用 HTMX + Web Component（非 SPA）
    
- 每個 app 都有自己的 route / 介面 / handler
    
- 強調模組獨立、可維護、避免耦合

# 前端組裝目錄結構 apps/

```bash
xxx-app/
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
├── conf/                    # 組態讀取（載入 pages.yaml、env.yaml）
├── models/                  # 結構定義（如 device, alert 等）
├── services/                # 可呼叫的邏輯單元（如 ExecService, LogService）
└── templates/               # 若有 templ 模板引擎（可用於 layout 組裝）
```

# 前端開發目錄結構

```bash
pkg/webui/                
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

✅ 原因：HTMX + Go SSR 是「後端內嵌式前端」
-----------------------------

### 特點如下：

| 特性 | 說明 |
| --- | --- |
| 🔧 單一構建流程 | 前端 HTML 就是 template，與 Go 一起編譯或部署 |
| 🧱 模板緊耦合業務邏輯 | 使用 `html/template` 或 `templ`，直接在 Go handler 中呼叫 |
| 🚫 無需 npm, webpack, vite | 不像 SPA 不需要 JS 編譯環境 |
| ✅ 可重用 layout / component | 可將共用 HTML component 抽象到 `pkg/webui/` 或 `apps/x/web/partial/` |
| ⛔ 不適合獨立開發 | 單靠 HTML + Go 無法做獨立前端開發流程 |

* * *

✅ 建議專案結構（HTMX + SSR 適用）
-----------------------

```bash
apps/
├── alert-app/
│   ├── main.go
│   ├── handler/
│   │   └── alert.go
│   └── web/
│       ├── layout/
│       ├── partials/
│       └── pages/
├── pdu-app/
│   └── web/
pkg/
└── webui/
    ├── component.go     ← alertCard, dataTable 等模板元件
    └── render.go        ← HTML renderer 抽象
```

- `web/` 是各 app 的 template
    
- `pkg/webui/` 是共用元件，例如 Alert Box、Tag 列、Table 分頁
    
- handler 直接 render 模板：
    

```go
webui.Render(ctx, "web/pages/alert_list.html", data)
```

* * *

✅ 開發 / 維運優勢
-----------

| 項目 | 說明 |
| --- | --- |
| 🧪 開發效率 | 修改畫面即時生效，不需 build / compile |
| 🧩 維護簡單 | 所有頁面綁定 app 本體，易於追蹤變更來源 |
| 🔁 可插拔性 | 可以讓某些 app 完全不含 web（純 API），其他有 UI |
| ⏱ 無前端 build 負擔 | CI/CD 無需跑 Vite / React build |


### 🧱 每個 app 擁有自己的 `web/`：

```bash
apps/
├── pdu-app/
│   ├── main.go
│   └── web/       ← templates, handlers, components
├── alert-app/
│   └── web/
```

### 🧱 抽象可重用元件至 `pkg/webui/`

```bash
pkg/
└── webui/
    ├── render.go      ← HTML render + layout handler
    └── component.go   ← 公用元件如 Tabs, Table, Dialog
```


