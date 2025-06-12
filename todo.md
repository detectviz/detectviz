# TODO 待補實作與修正項目清單（來自 detectviz-source-code.md 分析）

## 1. internal/api 路由與統一回應格式實作

參考文件：
- [docs/interfaces/api.md](./docs/interfaces/api.md)
- [docs/api-architecture.md](./docs/api-architecture.md)

### 起手實作項目

- [ ] internal/api/router.go：初始化 Echo，註冊 API 路由分組，掛載各模組 handler。
- [ ] internal/api/middleware.go：集中註冊 middleware 鏈（auth, logger, tracing, metrics 等）。
- [ ] internal/api/response.go：統一 JSON 回應與錯誤格式封裝方法。
- [ ] internal/api/errors/api_error.go：定義 API 專用錯誤型別與錯誤碼。
- [ ] internal/api/dtos/alert.go：定義 alert 模組用於 API 的 request/response DTO 結構。
- [ ] internal/api/dtos/rule.go：定義 rule 模組 DTO 結構。
- [ ] internal/api/dtos/notifier.go：定義 notifier 模組 DTO 結構。

### 測試與擴充

- [ ] internal/test/api/router_test.go：驗證 API router 註冊是否正確，涵蓋版本、路徑。
- [ ] internal/test/api/response_test.go：驗證回應格式是否一致，錯誤碼是否正常包裝。
- [ ] internal/test/api/dto_validation_test.go：測試 DTO binding 與 validator 整合行為。
- [ ] plugins/api/metrics.go：提供 plugin 擴充 API handler 範例（如 `/api/plugins/{mod}/metrics`）。

## 2. internal/middleware 掛載中介層、日誌、追蹤等

參考文件：
- [docs/interfaces/middleware.md](./docs/interfaces/middleware.md)
- [docs/middleware-architecture.md](./docs/middleware-architecture.md)

### 核心中介層

- [ ] internal/middleware/auth.go：呼叫 Authenticator 並注入 UserInfo
- [ ] internal/middleware/logger.go：記錄請求資訊（method, path, duration）
- [ ] internal/middleware/metrics.go：暴露 HTTP request metrics（status, duration）
- [ ] internal/middleware/tracing.go：插入 trace span 與 context propagation
- [ ] internal/middleware/recovery.go：捕捉 panic 回應 500 並記錄錯誤
- [ ] internal/middleware/gziper.go：自動壓縮回應內容（Content-Encoding）
- [ ] internal/middleware/cookies/cookies.go：封裝 cookie 設定與保護
- [ ] internal/middleware/csrf/csrf.go：產生與驗證 CSRF token
- [ ] internal/middleware/requestmeta/request_metadata.go：設定 X-Request-ID, User-Agent 等 header
- [ ] internal/middleware/testing.go：提供單元測試用中介層鏈結模擬器與輔助函式

### Plugin 擴充支援

- [ ] internal/registry/middleware/registry.go：支援 MiddlewarePlugin 註冊與使用
- [ ] plugins/middleware/cors.go：plugin 註冊 CORS middleware（範例）
- [ ] plugins/middleware/ratelimit.go：plugin 註冊限速 middleware（範例）


## 3. plugins/ 模組 Config 擴充設定標準化

- [ ] plugins/metric/flux/config.go：定義 Flux Plugin 設定結構與預設值。需設計 Flux 插件的設定 struct，並提供預設參數與驗證。
- [ ] plugins/metric/prom/config.go：定義 Prom Plugin 查詢參數。需實作 Prometheus 插件查詢設定 struct 與驗證邏輯。
- [ ] plugins/notifier/email/config.go：定義 Email 通知格式與驗證。需設計 email 通知的設定格式並驗證輸入。
- [ ] plugins/notifier/slack/config.go：Slack webhook 與格式設定。需設計 Slack webhook 設定 struct，支援訊息格式定義。
- [ ] plugins/notifier/webhook/config.go：Generic webhook 自定義結構。需設計通用 webhook 設定格式，支援自定義 header、body 等。
- [ ] docs/interfaces/plugin_config.md：各類 plugin config 結構與驗證方式說明。補充 flux/prom/email/slack/webhook 等格式範例與驗證邏輯。

## 4. internal/plugins lifecycle 控制與註冊

參考文件：
- [docs/architecture/plugins.md](./docs/architecture/plugins.md)
- [docs/interfaces/plugins.md](./docs/interfaces/plugins.md)

### Plugin 管理模組

- [x] internal/plugins/manager/lifecycle.go：控制 plugin 啟動與關閉流程（已實作）
- [x] internal/plugins/manager/loader.go：載入 plugins 之邏輯（已實作）
- [x] internal/plugins/manager/process.go：執行 Init/Close 流程控制（已實作）
- [x] internal/plugins/manager/registry.go：註冊並掛載 plugins（已實作）

### Plugin 註冊點與類型擴充

- [ ] internal/plugins/plugins/auth/keycloak/keycloak.go：實作 Keycloak 認證策略註冊範例
- [ ] internal/plugins/plugins/middleware/cors/cors.go：註冊 CORS MiddlewarePlugin
- [ ] internal/plugins/plugins/middleware/ratelimit/ratelimit.go：註冊 RateLimit MiddlewarePlugin
- [ ] internal/plugins/plugins/apihooks/alerts/alert_status.go：自訂 API route 插件範例
- [ ] internal/plugins/plugins/eventbus/alertlog/handler.go：事件擴充處理器

### Plugin 測試與驗證

- [ ] internal/test/plugins/plugin_lifecycle_test.go：測試註冊、啟動、關閉 plugin 的全流程
- [ ] internal/test/plugins/plugin_registry_test.go：驗證 plugin 註冊與反註冊的行為一致性
- [ ] internal/test/plugins/fake_plugins.go：定義 fake plugin 實作供測試注入與模擬
- [ ] internal/test/plugins/fake_context.go：模擬 plugin 執行上下文，驗證 Init/Close 行為
- [ ] internal/test/plugins/plugin_integration_test.go：驗證 plugin 註冊後與核心模組整合情境（ex: middleware 呼叫順序）

### 擴充建議

- [ ] 支援 plugins/ 目錄自動掃描與 lazy loading
- [ ] plugins 支援版本號與 metadata（如 name, description）
- [ ] CLI 工具自動產生 plugin scaffold（未來 roadmap）



## 5. internal/handlers 各模組 API 進入點實作

參考文件：
- [docs/interfaces/handlers.md](./docs/interfaces/handlers.md)
- [docs/handlers-architecture.md](./docs/handlers-architecture.md)

### 建立範本結構與核心組件

- [ ] internal/handlers/alert/v0alpha1/alert_handler.go：定義 Alert 模組 API handler，支援 GET /api/v0alpha1/alert/status。
- [ ] internal/handlers/rule/v0alpha1/rule_handler.go：定義 Rule 模組 API handler，支援基本 CRUD。
- [ ] internal/handlers/notifier/v0alpha1/notifier_handler.go：定義 Notifier 模組 handler。
- [ ] internal/handlers/common/response.go：定義 JSON success/error 包裝器。
- [ ] internal/handlers/common/errors.go：定義標準錯誤常數與結構。

### 測試與未來擴充

- [ ] internal/test/handlers/alert_handler_test.go：建立 Alert handler 單元測試。
- [ ] internal/test/handlers/rule_handler_test.go：建立 Rule handler 測試。
- [ ] internal/test/handlers/response_test.go：測試回應包裝行為。
- [ ] plugins/handlers/：提供 plugin handler 註冊範例，對應 plugin lifecycle。



## 6. internal/services 業務邏輯核心

參考文件：
- [docs/interfaces/services.md](./docs/interfaces/services.md)
- [docs/architecture/services.md](./docs/architecture/services.md)

### Service interface 與實作項目

- [ ] internal/services/interfaces.go：集中定義跨模組的 service 介面
- [ ] internal/services/rule/service.go：實作 RuleService，整合 rule store 與 notifier
- [ ] internal/services/notifier/service.go：實作 NotifierService，支援通知發送、註冊與篩選
- [ ] internal/services/logger/service.go：實作 LoggerService，支援結構化日誌寫入與查詢
- [ ] internal/services/metrics/service.go：實作 MetricsService，支援時序資料寫入與查詢
- [ ] internal/services/eventbus/service.go：實作 EventBusService，支援事件發送與訂閱觸發

### 每個 service 實作至少需支援以下邏輯：

- 注入對應模組 store（如 RuleStore、NotifierStore）
- 調用跨模組 service（如 RuleService 呼叫 NotifierService）
- 使用 eventbus 發送內部事件（如 rule 建立後發送通知）
- 回傳標準錯誤與結構封裝

### 測試與擴充建議

- [ ] internal/test/services/rule_service_test.go：測試 rule service CRUD 行為與事件觸發
- [ ] internal/test/services/notifier_service_test.go：測試通知篩選與發送整合流程
- [ ] internal/test/services/logger_service_test.go：測試日誌寫入與條件查詢
- [ ] internal/test/services/metrics_service_test.go：測試資料點寫入與範圍查詢
- [ ] internal/test/services/eventbus_service_test.go：測試多主題訂閱與事件觸發

### 擴充建議

- [ ] 支援 decorator：如 TracingService、AuditService 包裝原有 service
- [ ] 可拆分 UseCase 邏輯單元：ex. CreateRuleUseCase、TriggerAlertUseCase
- [ ] 支援 WithTenant(context.Context)：實作多租戶與權限管理流程

## 7. internal/store 資料儲存與讀寫後端

參考文件：
- [docs/interfaces/store.md](./docs/interfaces/store.md)
- [docs/architecture/store.md](./docs/architecture/store.md)

### 起手模組：RuleStore

- [ ] internal/store/rule/interfaces.go：定義 RuleStore 介面
- [ ] internal/store/rule/models/rule.go：定義 Rule 結構與查詢條件
- [ ] internal/store/rule/resolver.go：註冊與解析多 backend 實作
- [ ] internal/store/rule/memory/rule_store.go：提供測試用記憶體實作
- [ ] internal/store/rule/cache/rule_store.go：實作 Redis 快取邏輯（read-through）
- [ ] internal/store/rule/logfile/rule_store.go：檔案儲存實作（可 JSON Append-only）
- [ ] internal/store/rule/mysql/rule_store.go：MySQL-based 實作，使用 sqlx/GORM
- [ ] internal/store/rule/influxdb/rule_store.go：InfluxDB 實作（支援查詢與儲存）

### 起手模組補充：其他模組介面定義

- [ ] internal/store/notifier/interfaces.go：定義 NotifierStore 介面
- [ ] internal/store/logger/interfaces.go：定義 LoggerStore 介面
- [ ] internal/store/metrics/interfaces.go：定義 MetricsStore 介面
- [ ] internal/store/eventbus/interfaces.go：定義 EventBusStore 介面

### 其他模組實作項目

- [ ] internal/store/notifier/logfile/notifier_store.go：notifier 模組的檔案儲存實作
- [ ] internal/store/logger/logfile/logger_store.go：logger 模組的 logfile 實作
- [ ] internal/store/eventbus/memory/eventbus_store.go：eventbus 使用記憶體模擬儲存
- [ ] internal/store/metrics/influxdb/metrics_store.go：metrics 模組使用 InfluxDB 儲存

### Plugin 實作與註冊

- [ ] plugins/datasources/rulestore_influxdb/impl.go：plugin 實作 RuleStore 的 Influx 版本
- [ ] plugins/datasources/rulestore_mysql/impl.go：plugin 實作 RuleStore 的 MySQL 版本
- [ ] plugins/datasources/notifierstore_logfile/impl.go：plugin 實作 NotifierStore 的 Logfile 版本
- [ ] plugins/datasources/loggerstore_logfile/impl.go：plugin 實作 LoggerStore 的 Logfile 版本
- [ ] plugins/datasources/metricsstore_influxdb/impl.go：plugin 實作 MetricsStore 的 Influx 版本
- [ ] plugins/datasources/eventbusstore_memory/impl.go：plugin 實作 EventBusStore 的 Memory 版本
- [ ] internal/store/resolver.go：集中註冊所有 store 實作，支援多來源 fallback

### 測試與擴充

- [ ] internal/test/store/rule_store_test.go：針對 memory, mysql 等 RuleStore backend 撰寫測試
- [ ] internal/test/store/notifier_store_test.go：測試 notifier logfile 實作
- [ ] internal/test/store/logger_store_test.go：測試 logger store
- [ ] internal/test/store/metrics_store_test.go：測試 metrics store (influx)
- [ ] internal/test/store/cache_wrapper_test.go：測試快取包裝邏輯是否正確


## 8. internal/system 系統服務與整合層模組

參考文件：
- [docs/interfaces/system_components.md](./docs/interfaces/system_components.md)
- [docs/architecture/system.md](./docs/architecture/system.md)

### 預計目錄結構

```
internal/system/
├── http/              # 底層 HTTP 功能（proxy, cors, apiserver）
├── diagnostics/       # 健康診斷模組（stats, supportbundle）
├── integration/       # 外部整合介接（grpc, search, live）
├── platform/          # 快取、配額、生命週期等平台能力
```

---

### http/

- [ ] internal/system/http/apiserver.go：提供簡化 API server 包裝器（ex: Echo 初始化）
- [ ] internal/system/http/cors.go：提供 CORS middleware provider 實作
- [ ] internal/system/http/proxy.go：反向代理轉發工具

---

### diagnostics/

- [ ] internal/system/diagnostics/stats.go：統計資訊匯出介面（health, runtime, goroutine）
- [ ] internal/system/diagnostics/supportbundle.go：支援打包診斷資訊（log, config, env）

---

### integration/

- [ ] internal/system/integration/grpcserver.go：gRPC server 啟動器與生命週期管理
- [ ] internal/system/integration/search.go：統一 search 查詢入口整合
- [ ] internal/system/integration/live.go：支援 live update 或 websocket push（未來）

---

### platform/

- [ ] internal/system/platform/quota.go：定義與管理資源使用配額
- [ ] internal/system/platform/cache.go：全域快取 provider（如 redis 或 memory）
- [ ] internal/system/platform/hooks.go：模組初始化/終止時可註冊的全域 hook

---

### 測試與擴充建議

- [ ] internal/test/system/http_cors_test.go：測試 CORS middleware 行為與註冊
- [ ] internal/test/system/diagnostics_stats_test.go：測試健康診斷輸出格式與內容
- [ ] internal/test/system/platform_quota_test.go：測試配額限制是否正確判斷與拒絕
- [ ] internal/test/system/integration_grpcserver_test.go：測試 gRPC server 啟停與註冊服務是否正確




參考文件：
- [docs/interfaces/infra.md](./docs/interfaces/infra.md)
- [docs/architecture/infra.md](./docs/architecture/infra.md)

### 建議目錄結構

```
pkg/infra/
├── log/             # 封裝 zap logger
├── config/          # 組態載入與解析（env/yaml/flag）
├── trace/           # OpenTelemetry trace 工具
├── db/              # 資料庫驅動註冊與 sqlx 管理
├── cache/           # 快取 client（如 redis/memory）
├── metrics/         # Prometheus registry 與 middleware
├── signal/          # signal.NotifyContext 封裝
├── httpclient/      # 統一 HTTP 調用封裝（支援 middleware, retry, auth）
```

---

### 起手實作項目

- [ ] pkg/infra/log/zap_logger.go：建立 Logger interface 實作並導出 zap
- [ ] pkg/infra/config/provider.go：統一載入組態來源（env, yaml, flags）
- [ ] pkg/infra/trace/otel_tracer.go：整合 OpenTelemetry provider 與簡化 API
- [ ] pkg/infra/cache/memory.go：記憶體快取簡易實作（供測試與內嵌使用）
- [ ] pkg/infra/metrics/prometheus.go：初始化與註冊 Prometheus 指標
- [ ] pkg/infra/db/driver.go：註冊 mysql, sqlite driver，包裝連線 helper
- [ ] pkg/infra/signal/signalctx.go：封裝 context.WithSignal 與 cancel flow

---

- [ ] pkg/infra/httpclient/client.go：建立具 retry、timeout、auth header 的 HTTP client 抽象
- [ ] pkg/infra/httpclient/config.go：設定建構選項（如 timeout、base headers）
- [ ] pkg/infra/httpclient/middleware/trace.go：middleware 插入 OpenTelemetry trace context
- [ ] pkg/infra/httpclient/middleware/metrics.go：middleware 暴露 Prometheus 指標
- [ ] pkg/infra/httpclient/middleware/auth.go：middleware 插入 Authorization header

---

### 測試與擴充建議

- [ ] internal/test/infra/logger_test.go：測試 zap logger 是否符合 Logger interface
- [ ] internal/test/infra/config_provider_test.go：測試 env/yaml merge 行為與解析正確性
- [ ] internal/test/infra/cache_test.go：測試快取 get/set 行為與過期時間控制
- [ ] internal/test/infra/signalctx_test.go：測試 cancel 發生時 context 是否同步關閉

---

### 測試與擴充建議（httpclient）

- [ ] internal/test/infra/httpclient_client_test.go：測試基本 request/response、timeout 是否正常
- [ ] internal/test/infra/httpclient_retry_test.go：測試 retry 行為與失敗 fallback
- [ ] internal/test/infra/httpclient_auth_test.go：測試 Authorization middleware 是否正確加入 header
- [ ] internal/test/infra/httpclient_metrics_test.go：測試 Prometheus 指標是否正確記錄
- [ ] internal/test/infra/httpclient_trace_test.go：測試 trace context 是否正確注入與攜帶

---

### 擴充建議

- [ ] 支援 interface 替換與注入：Logger, ConfigProvider, Cache 等
- [ ] 支援 testing/fake 實作（如 fakeLogger, fakeCache）
- [ ] 可由 bootstrap 初始化階段建立並注入服務模組



## 9. pkg/security/encryption 加解密封裝模組

參考文件：
- [docs/interfaces/encryption.md](./docs/interfaces/encryption.md)

---

### 建議目錄結構

```
pkg/security/encryption/
├── interface.go         # 定義 Encryptor interface
├── aes.go               # AES 加密實作
├── registry.go          # 註冊與取得預設 Encryptor
├── error.go             # 加解密錯誤封裝
```

---

### Interface 與實作項目

- [ ] interface.go：定義 Encryptor interface（ID, Encrypt, Decrypt）
- [ ] aes.go：AES-CFB 模式的預設加密器，支援 key 設定與隨機 IV
- [ ] registry.go：加密 provider 註冊與預設注入邏輯（供 plugin 擴充）
- [ ] error.go：包裝錯誤型別與訊息，避免洩漏敏感資訊

---

### Plugin 實作（未來擴充）

- [ ] plugins/encryption/vault/impl.go：支援 HashiCorp Vault provider
- [ ] plugins/encryption/kms/impl.go：支援 AWS/GCP KMS provider
- [ ] plugins/encryption/plaintext/impl.go：純測試用 provider（無加密）

---

### 測試與覆蓋建議

- [ ] encryption/aes_test.go：測試加密/解密一致性、IV 正確性
- [ ] encryption/error_test.go：測試錯誤包裝與 unwrap 行為
- [ ] encryption/registry_test.go：測試註冊與替換 Encryptor provider

---

### 使用情境

- [ ] store 層寫入敏感資料時呼叫 Encryptor
- [ ] config 模組可提供預設金鑰載入
- [ ] service 層不得自行決定加密方式，應透過注入 provider 呼叫 Encrypt/Decrypt

---

## 10. pkg/validation 輸入驗證模組

參考文件：
- [docs/interfaces/validation.md](./docs/interfaces/validation.md)

---

### 建議目錄結構

```
pkg/validation/
├── basic.go         # 基本欄位驗證（required, min/max length）
├── regexp.go        # 格式驗證（email, IP, slug）
├── jsonschema.go    # JSON schema 驗證輔助（選用）
├── error.go         # 統一錯誤訊息結構
```

---

### 實作清單

- [ ] pkg/validation/basic.go：實作 Required、MaxLength、MinLength、NotEmptySlice 等基本驗證函式
- [ ] pkg/validation/basic_test.go：測試基本驗證函式的正確性與邊界情況
- [ ] pkg/validation/regexp.go：實作格式驗證函式，如 IsEmail、IsSlug、IsIPv4、IsUUID 等
- [ ] pkg/validation/regexp_test.go：測試格式驗證邏輯與特殊輸入案例
- [ ] pkg/validation/error.go：定義 FieldError 結構與錯誤包裝工具
- [ ] pkg/validation/error_test.go：測試錯誤包裝訊息與欄位回報正確性
- [ ] pkg/validation/jsonschema.go：提供 ValidateJSONSchema 方法，支援 schema 驗證（選用）
- [ ] pkg/validation/jsonschema_test.go：測試 schema 驗證錯誤與錯誤訊息包裝


## 11. pkg/cmd CLI 指令整合


參考文件：
- [docs/architecture/cmd.md](./docs/architecture/cmd.md)
- [docs/interfaces/cmd.md](./docs/interfaces/cmd.md)

---

### 建議目錄結構

```
pkg/cmd/
├── root.go             # detectviz 主命令，整合所有子指令
├── rule/
│   ├── apply.go        # detectviz rule apply
│   └── list.go         # detectviz rule list
├── plugin/
│   ├── list.go         # detectviz plugin list
│   └── enable.go       # detectviz plugin enable
├── config/
│   └── show.go         # detectviz config show
```

---

### CLI 指令項目

- [ ] root.go：初始化 rootCmd 與 global flags，匯出 Execute()
- [ ] rule/apply.go：實作匯入 YAML/JSON 規則並送出 API
- [ ] rule/list.go：呼叫 API 取得所有 rules 並輸出表格
- [ ] plugin/list.go：列出目前已載入的 plugin
- [ ] plugin/enable.go：啟用指定 plugin，支援 args 傳遞
- [ ] config/show.go：列出目前的設定來源與組態

---

### CLI 測試建議

- [ ] internal/test/cmd/root_test.go：測試空參數、預設幫助、flag 整合
- [ ] internal/test/cmd/rule_apply_test.go：模擬匯入檔案，驗證 service 呼叫
- [ ] internal/test/cmd/plugin_list_test.go：模擬 plugin registry 呼叫
- [ ] internal/test/cmd/config_show_test.go：測試組態輸出是否正確

---

### 實作整合點

- [ ] apps/cli/main.go：作為獨立 CLI 工具進入點，呼叫 cmd.Execute()
- [ ] 可與 internal/plugins/plugin.CLIPlugin 接口整合，支援 CLI 插件註冊

---

## 12. pkg/utils 輔助工具函式

參考文件：
- [docs/architecture/utils.md](./docs/architecture/utils.md)
- [docs/interfaces/utils.md](./docs/interfaces/utils.md)

---

### 建議模組結構

```
pkg/utils/
├── stringutil/
│   ├── string.go
│   └── string_test.go
├── timeutil/
│   ├── time.go
│   └── time_test.go
├── sliceutil/
│   ├── slice.go
│   └── slice_test.go
├── encodeutil/
│   ├── encode.go
│   └── encode_test.go
├── idutil/
│   ├── id.go
│   └── id_test.go
├── validateutil/
│   ├── validate.go
│   └── validate_test.go
├── retryer/
│   ├── retryer.go
│   └── retryer_test.go
├── pathutil/
│   ├── path.go
│   └── path_test.go
├── pointer/
│   ├── pointer.go
│   └── pointer_test.go
```

---

### 測試與驗證

- [ ] 每個子目錄皆需具備 `_test.go` 測試覆蓋
- [ ] 測試涵蓋正常、邊界與錯誤情境
- [ ] 禁止跨模組依賴（維持純工具屬性）

---

### 實作清單

- [ ] pkg/utils/stringutil/string.go：常用字串處理工具，如 Slugify、Elide 等
- [ ] pkg/utils/stringutil/string_test.go：覆蓋 stringutil 測試情境
- [ ] pkg/utils/timeutil/time.go：時間處理工具，如 ParseRFC3339、NowMillis 等
- [ ] pkg/utils/timeutil/time_test.go：覆蓋 timeutil 測試情境
- [ ] pkg/utils/sliceutil/slice.go：泛型 slice 操作，如 Dedup、Chunk、Filter
- [ ] pkg/utils/sliceutil/slice_test.go：測試泛型操作與邊界行為
- [ ] pkg/utils/encodeutil/encode.go：JSON、Base64 編碼工具
- [ ] pkg/utils/encodeutil/encode_test.go：測試 JSON 安全轉換與解碼行為
- [ ] pkg/utils/idutil/id.go：ID 產生器，如 UUIDv4、NanoID 等
- [ ] pkg/utils/idutil/id_test.go：測試 ID 格式與隨機性
- [ ] pkg/utils/validateutil/validate.go：基本驗證工具，如 IsEmail、IsUUID 等
- [ ] pkg/utils/validateutil/validate_test.go：驗證常見欄位與格式函式
- [ ] pkg/utils/retryer/retryer.go：封裝重試邏輯與 backoff 策略
- [ ] pkg/utils/retryer/retryer_test.go：測試重試成功/失敗與 RetryIf 判斷
- [ ] pkg/utils/pathutil/path.go：安全檔案與 URL 路徑處理
- [ ] pkg/utils/pathutil/path_test.go：測試路徑拼接、相對/絕對判定
- [ ] pkg/utils/pointer/pointer.go：指標輔助，如 StringPtr、BoolPtr 等
- [ ] pkg/utils/pointer/pointer_test.go：測試指標生成與指標值判斷

## 13. internal/rbac 權限控制模組與擴充支援

參考文件：
- [docs/architecture/rbac.md](./docs/architecture/rbac.md)
- [docs/interfaces/rbac.md](./docs/interfaces/rbac.md)

---

### 模組實作清單

- [ ] internal/rbac/accesscontrol/authorizer.go：主授權介面與 dispatch 控制
- [ ] internal/rbac/accesscontrol/scope.go：定義 scope 接口與 GlobalScope、OrgScope 等實作
- [ ] internal/rbac/accesscontrol/checker.go：快取與角色比對優化邏輯
- [ ] internal/rbac/accesscontrol/middleware.go：Echo middleware 授權驗證邏輯
- [ ] internal/rbac/org/org.go：組織邏輯與租戶切換實作
- [ ] internal/rbac/org/org_delete_svc.go：刪除組織與相關 metadata 的清理邏輯
- [ ] internal/rbac/org/model.go：組織模型與結構定義
- [ ] internal/rbac/team/team.go：支援群組管理（可選）
- [ ] internal/rbac/mock/fake.go：測試用假資料與假物件
- [ ] internal/rbac/mock/mock.go：mock interface 實作供 service 層注入

---

### plugin 擴充模組（可選實作）

- [ ] internal/plugins/serviceaccounts/serviceaccount_provider.go：服務帳號認證支援
- [ ] internal/plugins/tempuser/invite_token.go：臨時使用者建立與驗證

---

### 測試與驗證建議

- [ ] internal/test/rbac/authorizer_test.go：測試權限檢查邏輯與角色範疇判定
- [ ] internal/test/rbac/scope_test.go：測試 scope 組合與 Matches 函式行為
- [ ] internal/test/rbac/middleware_test.go：測試 Echo 路由 middleware 鍊授權處理
- [ ] internal/test/rbac/org_test.go：組織 CRUD 測試流程
- [ ] internal/test/rbac/mock_test.go：mock 對應行為正確性驗證