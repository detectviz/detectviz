
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