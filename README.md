# Detectviz

Detectviz 是一套基於 Clean Architecture 設計的模組化監控與告警平台，支援指標查詢、條件比對、事件發布與通知處理。透過 Plugin 機制整合各種數據來源（如 Prometheus, InfluxDB, Flux）與通知通道（如 Email, Slack, Webhook），並提供可維護、可擴充的事件處理架構。

---

## 專案目錄結構

```bash
detectviz/
├── apps/                     # 每個 App 對應一套業務 API / UI
│   ├── {module}-app/
│   │   ├── main.go
│   │   ├── routes.go
│   │   ├── conf/
│   │   ├── web/              # HTMX 頁面（可含 layout, partials, pages）
│		│   └── handler/          # HTTP handler 層
│		│
├── internal/                 # 核心邏輯模組（僅供 apps 使用）
│   ├── {module}/
│		│   ├── alert.go          # Init, Enabled
│		│   ├── service.go        # 實作 interface
│		│   ├── interface.go      # 定義給 bootstrap 用的接口
│		│   ├── handler.go        # 若有 REST API
│		│   ├── cmd.go            # 若有 CLI
│		│   └── eventbus.go       # 若有事件訂閱
│		│
│   ├── adapters/             # 各模組抽象介面實作
│   ├── registry/             # 模組註冊中心
│   └── test/                 # 整合測試、fakes、mocks、testutil 工具
├── pkg/                      # 共用抽象（interface、config、domain）
│   ├── config/               # 設定載入與注入模組
│   ├── ifaces/               # 模組抽象介面定義
│   └── mocks/                # 使用 mockery 產出的 mock interface（自動生成）
├── plugins/                  # 可插拔模組：可獨立引用、註冊、替換
├── scripts/                  # 輔助腳本（備份、啟動、模擬工具）
├── deploy/                   # Docker 與環境部署相關設定
├── bulid/                    # 建置相關的工具和腳本，主要用於 CI/CD 和打包過程
├── docs/                     # 架構文件、介面規範、擴充開發指南
└── README.md
```

- `pkg`: 可重用模組、interface、工具（對外穩定）
- `internal`: 各業務邏輯模組（僅供 app 使用，不外部引用）

---

## 已實作模組

- **Logger**：支援 Zap 實作與 NopLogger。
- **ConfigProvider**：統一提供全域設定注入。
- **EventBus**：可註冊多種事件處理器（Host, Metric, Alert, Task）。
- **AlertEvaluator**：支援 Prometheus、Flux 查詢條件擴充。
- **Scheduler**：支援 Cron 與 WorkerPool 型任務排程。
- **Notifier**：支援 Email、Slack、Webhook 多種通道。

---

## 啟動方式

請搭配各 `apps/` 內主程式使用 `go run` 或 `make` 指令：

```bash
go run ./apps/alert-app/main.go
make run-scheduler
```

可參考 `scripts/` 或 `Makefile` 中的啟動流程與模擬指令。

---

## 文件導引

- [docs/interfaces/](docs/interfaces/)：介面定義與實作契約說明
- [internal/registry/](internal/registry/)：模組註冊流程（AlertEvaluator、Notifier、Scheduler 等）
- [internal/test/README.md](internal/test/README.md)：測試策略與實際目錄規劃
- [docs/develop-guide.md](docs/develop-guide.md)：設計原則與架構圖
- [docs/coding-style-guide.md](docs/coding-style-guide.md)：程式撰寫風格（命名規則、註解格式、golangci-lint 設定）
