# Detectviz

Detectviz 是一套基於 Clean Architecture 設計的模組化監控與告警平台，支援指標查詢、條件比對、事件發布與通知處理。透過 Plugin 機制整合各種數據來源（如 Prometheus, InfluxDB, Flux）與通知通道（如 Email, Slack, Webhook），並提供可維護、可擴充的事件處理架構。

---

## 專案目錄結構

```
detectviz/
├── apps/                     # 獨立 App 模組（如 alert-app、web-app 等）
├── cmd/                      # CLI 主程式進入點（可選）
├── internal/                 # 核心邏輯模組（僅供 apps 使用）
│   ├── adapters/             # 各模組實作（logger, notifier, scheduler 等）
│   ├── registry/             # 模組註冊中心（alert, notifier, scheduler 等）
│   └── test/                 # 整合測試、fakes、mocks、testutil 工具
├── pkg/                      # 共用抽象（interface、config、domain）
│   ├── config/               # 設定載入與注入模組
│   ├── ifaces/               # 模組介面定義（Logger, Scheduler, Notifier 等）
│   └── mocks/                # 使用 mockery 產出的 mock interface（自動生成）
├── plugins/                  # 插件模組（資料查詢來源、通知管道擴充等）
├── scripts/                  # 輔助腳本（備份、啟動、模擬工具）
├── deploy/                   # Docker 與環境部署相關設定
└── README.md
```

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
