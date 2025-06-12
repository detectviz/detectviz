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
│		│                       
│   ├── middleware/           # 所有 HTTP middleware 實作
│		│		├── auth.go
│		│		├── logger.go
│		│		├── tracing.go
│		│		├── metrics.go
│		│		├── recovery.go
│		│		├── csrf/csrf.go
│		│		├── cookies/cookies.go
│		│		├── requestmeta/request_metadata.go
│		│		└── testing.go
│		│
│   ├── api/                   # API server 啟動與路由註冊邏輯
│		│		├── errors/
│		│		│   └── api_error.go   # 定義 ErrorCode, Message, ToJSON()
│		│		├── dtos/              # 用於 DTO / Response 共用模型
│		│		│   ├── alert.go       # CreateAlertRequest, AlertResponse
│		│		│   ├── rule.go
│		│		│   └── common.go
│		│		├── response/
│		│		│   └── json_response.go # JSON(), ErrorJSON(), WithStatus()
│		│   ├── router.go          # 統一註冊模組 API route
│		│   ├── middleware.go      # 組合中介層鏈結
│		│   └── server.go          # HTTP Server 啟動與註冊控制
│		│
│   ├── handlers/              # 功能模組 API handler 與 controller 實作
│		│   ├── alert/v0alpha1/    # Alert 模組 REST handler
│		│   ├── rule/v0alpha1/     # Rule 模組 handler
│		│   ├── report/v0alpha1/   # 報表模組 handler
│		│   └── common/            # 回應格式、錯誤、驗證工具
│		│
│   ├── auth/                  # Authenticator 策略與登入邏輯（Plugin 可擴充）
│		│   ├── login/             # 傳統帳號密碼登入
│		│   ├── oauthtoken/        # 儲存 Token / Session
│		│   ├── ssosettings/       # 設定各種 SSO 登入參數
│		│   ├── strategies/        # 每種登入方式的策略模組（可 plugin 注入）
│		│   ├── context.go         # UserInfo 注入與提取
│		│   ├── identity.go        # Requester 介面（參考 Grafana）
│		│   └── registry.go        # 動態註冊多組 authenticator
│		│
│   ├── services/
│		│   ├── alert/
│		│		│   ├── alert.go           # Init, Enabled
│		│		│   ├── service.go         # 實作 interface
│		│		│   ├── interface.go       # 定義給 bootstrap 用的接口
│		│		│   ├── handler.go         # 若有 REST API
│		│		│   ├── cmd.go             # 若有 CLI
│		│		│   └── eventbus.go        # 若有事件訂閱
│		│		│
│   ├── adapters/              # 各模組抽象介面實作
│   │   ├── alert
│   │   │   ├── evaluator.go
│   │   │   ├── flux
│   │   │   ├── mock_adapter_test.go
│   │   │   ├── mock_adapter.go
│   │   │   └── prom
│   │   ├── cachestore
│   │   │   ├── memory
│   │   │   ├── redis
│   │   │   ├── registry_test.go
│   │   │   └── registry.go
│   │   ├── eventbus
│   │   │   ├── alert_test.go
│   │   │   ├── alert.go
│   │   │   ├── host_test.go
│   │   │   ├── host.go
│   │   │   ├── inmemory.go
│   │   │   ├── metric_test.go
│   │   │   ├── metric.go
│   │   │   ├── task_test.go
│   │   │   └── task.go
│   │   ├── importer
│   │   │   ├── registry_test.go
│   │   │   └── registry.go
│   │   ├── libraryelements
│   │   │   ├── service_adapter_test.go
│   │   │   └── service_adapter.go
│   │   ├── logger
│   │   │   ├── logger_test.go
│   │   │   ├── nop_adapter.go
│   │   │   └── zap_adapter.go
│   │   ├── metrics
│   │   │   ├── aggregator.go
│   │   │   ├── query_adapter.go
│   │   │   ├── series_reader_adapter.go
│   │   │   ├── transformer_adapter.go
│   │   │   └── writer_adapter.go
│   │   ├── modules
│   │   │   ├── engine_adapter_test.go
│   │   │   ├── engine_adapter.go
│   │   │   ├── listener_adapter.go
│   │   │   ├── registry_adapter.go
│   │   │   └── runner_adapter.go
│   │   ├── notifier
│   │   │   ├── email_adapter_test.go
│   │   │   ├── email_adapter.go
│   │   │   ├── mock_adapter.go
│   │   │   ├── multi.go
│   │   │   ├── nop.go
│   │   │   ├── slack_adapter_test.go
│   │   │   ├── slack_adapter.go
│   │   │   ├── webhook_adapter_test.go
│   │   │   └── webhook_adapter.go
│   │   ├── scheduler
│   │   │   ├── cron_adapter_test.go
│   │   │   ├── cron_adapter.go
│   │   │   ├── mock_adapter.go
│   │   │   ├── workerpool_adapter_test.go
│   │   │   └── workerpool_adapter.go
│   │   ├── server
│   │   │   └── server_adapter.go
│   │   └── versioning
│   │       ├── store_adapter_test.go
│   │       └── store_adapter.go
│   ├── registry/                # 模組註冊中心
│   │   ├── alert
│   │   │   └── registry.go
│   │   ├── cachestore
│   │   │   └── registry.go
│   │   ├── config
│   │   │   └── registry.go
│   │   ├── eventbus
│   │   │   ├── plugins.go
│   │   │   ├── providers.go
│   │   │   ├── registry_inmemory.go
│   │   │   └── registry.go
│   │   ├── logger
│   │   │   └── registry.go
│   │   ├── notifier
│   │   │   ├── registry_test.go
│   │   │   └── registry.go
│   │   ├── scheduler
│   │   │   └── registry.go
│   │   ├── decoder.go
│   │   ├── engine.go
│   │   ├── loader.go
│   │   └── registry.go
│   ├── store/             # 只依賴 interface，不直接操作底層存取

│   ├── plugins/               # 可插拔模組擴充（可註冊 middleware, auth 策略等）
│		│   ├── auth/              # 額外擴充的登入策略
│		│   ├── middleware/        # 其他中介層插件（如 CORS、限速器）
│   │   ├── apihooks/          # 提供平台 API 擴充註冊點
│   │   ├── eventbus/
│   │   │   └── alertlog
│   │   ├── manager/
│   │   │   ├── lifecycle_test.go
│   │   │   ├── lifecycle.go
│   │   │   ├── loader_test.go
│   │   │   ├── loader.go
│   │   │   ├── process_test.go
│   │   │   ├── process.go
│   │   │   ├── registry_test.go
│   │   │   └── registry.go
│   │   └── plugin.go
│   ├── rbac/
│		│   ├── accesscontrol     # 權限控管與角色資源策略
│		│   ├── org               # 組織管理與切換
│		│   ├── team              # 使用者群組功能
│		│   └── user              # 使用者 CRUD 與偏好設定
│   ├── system/
│		│		├── apiserver/        # REST 接口建構器
│		│		├── grpcserver/       # gRPC 接口與注入點
│		│		├── datasourceproxy/  # 多數據源後端轉發器
│		│		├── caching/          # 快取框架與策略
│		│		├── quota/            # 資源使用限制機制
│		│		├── supportbundles/   # 問題診斷壓縮包產生器
│		│		├── stats/            # 平台統計收集
│		│		├── hooks/            # 模組內事件 hook 機制
│		│		├── live/             # Live 推播或事件橋接器
│		│  	└── search/           # 資料或資源統一查詢服務
│   ├── server/
│   │   ├── instrumentation.go
│   │   ├── runner.go
│   │   └── server.go
│   ├── bootstrap/
│   │   ├── config_loader.go
│   │   ├── elements_loader.go
│   │   ├── versioning_loader.go
│   │   ├── init.go
│   │   └── wire.go
│   ├── modules/
│   │   ├── dependencies.go
│   │   ├── engine.go
│   │   ├── listener.go
│   │   ├── registry.go
│   │   └── runner.go
│   └── test/                 # 整合測試、fakes、mocks、testutil 工具
│
├── pkg/                      # 共用抽象（interface、config）            
│   ├── config/								# 設定載入與注入模組
│   │   ├── default.go
│   │   └── README.md
│   ├── configtypes/
│   │   ├── cache_config.go
│   │   └── notifier_config.go
│   ├── ifaces/               # 模組抽象介面定義
│   │   ├── alert/
│   │   │   ├── evaluate_test.go
│   │   │   ├── evaluate.go
│   │   │   └── evaluator.go
│   │   ├── bus/
│   │   │   ├── alert.go
│   │   │   ├── host.go
│   │   │   ├── metric.go
│   │   │   ├── task.go
│   │   │   └── types.go
│   │   ├── cachestore/
│   │   │   └── cachestore.go
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── event/
│   │   │   ├── alert.go
│   │   │   ├── host.go
│   │   │   ├── metric.go
│   │   │   ├── task.go
│   │   │   └── types.go
│   │   ├── eventbus/
│   │   │   ├── eventbus.go
│   │   │   └── provider.go
│   │   ├── logger/
│   │   │   ├── context_test.go
│   │   │   ├── context.go
│   │   │   ├── logger.go
│   │   │   └── nop_logger.go
│   │   ├── metrics/
│   │   │   ├── metric.go
│   │   │   ├── query.go
│   │   │   └── types.go
│   │   ├── modules/
│   │   │   └── modules.go
│   │   ├── notifier/
│   │   │   └── notifier.go
│   │   ├── plugins/
│   │   │   └── plugin.go
│   │   ├── registry/
│   │   │   └── registry.go
│   │   ├── scheduler/
│   │   │   ├── mock_adapter_test.go
│   │   │   └── scheduler.go
│   │   ├── server/
│   │   │   └── server.go
│   │   └── web/
│   │				├── context.go         # 自訂 context 包含 request/user/logger
│   │				├── router.go          # 註冊與匹配邏輯
│   │				├── binding.go         # JSON bind 與驗證
│   │				├── response_writer.go # 攔截與回應控制
│   │				├── web.go             # 主入口：定義 router、middleware
│   │ 			└── webtest/           # 單元測試與 chain 模擬
│   ├── registry/
│   │   ├── apis/
│   │   │   ├── datasource/
│   │   │   ├── host/
│   │   │   └── plugin/
│   │   ├── kinds/
│   │   │   ├── testdata
│   │   │   └── validator.go
│   │   ├── registry.go
│   │   └── schemas/
│   │       ├── datasource.schema.yaml
│   │       ├── host.schema.yaml
│   │       └── index.yaml
│   ├── importer/
│   │   └── interface.go
│   ├── libraryelements/
│   │   ├── interface.go
│   │   ├── registry.go
│   │   ├── store_memory.go
│   │   └── types.go
│   ├── validations           # 表單驗證、參數邏輯聚焦
│   ├── infra/
│   ├── utils/                # 各類通用工具、輔助函式
│   └── mocks/                # 使用 mockery 產出的 mock interface（自動生成）
├── plugins/                  # 可插拔模組：可獨立引用、註冊、替換
│   ├── auth
│   ├── datasources
│   ├── exporter
│   ├── tools
│   └── visuals
├── scripts/                  # 輔助腳本（備份、啟動、模擬工具）
├── deploy/                   # Docker 與環境部署相關設定
├── build/                    # 建置相關的工具和腳本，主要用於 CI/CD 和打包過程
├── docs/                     # 架構文件、介面規範、擴充開發指南
└── README.md
```

### 補充說明:
- `pkg`: 可重用模組、interface、工具（對外穩定）
- `internal`: 各業務邏輯模組（僅供 app 使用，不外部引用）
- `internal/api`:  API server 啟動與路由註冊邏輯
- `internal/apis`: 各模組功能 API handler 與版本切分
- `internal/auth`: 支援多種登入策略，可註冊外部 plugin
- `internal/middleware`: 提供通用中介層，可由 plugins 擴充
- 各 app 的 API route 將由各模組自行註冊並統一導入 router
- API 路由與模組 API 採用 plugin 式注入

---

## 已實作模組

- **Logger**：支援 Zap 實作與 NopLogger。
- **ConfigProvider**：統一提供全域設定注入。
- **EventBus**：可註冊多種事件處理器（Host, Metric, Alert, Task）。
- **AlertEvaluator**：支援 Prometheus、Flux 查詢條件擴充。
- **Scheduler**：支援 Cron 與 WorkerPool 型任務排程。
- **Notifier**：支援 Email、Slack、Webhook 多種通道。

## 尚未實作的模組（但已設計初版）

- `internal/web/`：負責 HTMX Web 組件、模板與畫面渲染
- `internal/store/`：提供統一 CRUD 接口，可由 plugin datasource backend 注入實作
- `internal/plugins/datasources/`：各種資料來源實作（influxdb, loki, file 等）
- `internal/services/`：封裝業務邏輯，不直接操作 handler 或 adapter
- `pkg/infra/metrics/`：Prometheus exporter 模組
- `pkg/infra/tracing/`：OpenTelemetry 追蹤邏輯
- `pkg/infra/httpclient/`：統一 http 呼叫邏輯與 middleware
- `pkg/security/encryption/`：AES 封裝與 provider 模組
- `pkg/utils/`：各類通用函式（pointer, string, retryer, uri 等）

### 未實作:
apps/alert-app
internal/middleware
internal/api
internal/apis
internal/auth
internal/alert
internal/rabc
internal/system/*
internal/plugins/auth
internal/plugins/middleware
internal/plugins/apihooks
pkg/infra/*
pkg/utils/*
plugins/*
internal/web



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

---


