# TODO 待補實作與修正項目清單（來自 detectviz-source-code.md 分析）


## 7. plugins/ 模組 Config 定義

- [ ] plugins/metric/flux/config.go：定義 Flux Plugin 設定結構與預設值。需設計 Flux 插件的設定 struct，並提供預設參數與驗證。
- [ ] plugins/metric/prom/config.go：定義 Prom Plugin 查詢參數。需實作 Prometheus 插件查詢設定 struct 與驗證邏輯。
- [ ] plugins/notifier/email/config.go：定義 Email 通知格式與驗證。需設計 email 通知的設定格式並驗證輸入。
- [ ] plugins/notifier/slack/config.go：Slack webhook 與格式設定。需設計 Slack webhook 設定 struct，支援訊息格式定義。
- [ ] plugins/notifier/webhook/config.go：Generic webhook 自定義結構。需設計通用 webhook 設定格式，支援自定義 header、body 等。
- [ ] docs/interfaces/plugin_config.md：各類 plugin config 結構與驗證方式說明。補充 flux/prom/email/slack/webhook 等格式範例與驗證邏輯。



## 8. 測試補齊與整合場景驗證

#### modules 模組
- [x] internal/test/fakes/fake_engine.go：Engine interface 測試用假件
- [x] internal/adapters/modules/engine_adapter_test.go：EngineAdapter 單元測試，驗證封裝與轉發
- [x] internal/test/modules/runner_test.go：模組整合注入測試，測試 Engine/Registry/Runner 全流程

#### importer 模組
- [x] internal/test/fakes/fake_importer.go：測試用 Importer 假件
- [x] internal/test/importer/importer_test.go：Importer 註冊與 CanImport/Import 行為測試

#### versioning 模組
- [x] internal/test/fakes/fake_versionstore.go：測試用 VersionStore 假件
- [x] internal/test/versioning/store_test.go：測試版本儲存與查詢邏輯
- [x] internal/adapters/versioning/store_adapter.go：封裝記憶體版或 DB 實作

#### libraryelements 模組
- [x] internal/test/fakes/fake_element_service.go：ElementService 假件
- [x] internal/test/libraryelements/service_test.go：測試 CRUD 與 Kind 註冊測試
- [x] internal/test/libraryelements/registry_test.go：測試元件渲染器註冊與查詢

#### plugin 模組
- [x] internal/test/registry/plugin_registry_test.go：PluginRegistry 註冊與查詢測試
- [x] internal/test/plugins/plugin_lifecycle_test.go：Plugin 啟動/關閉流程模擬

#### registry 模組
- [x] internal/test/registry/decoder_test.go：測試 DecodeAndValidate 對合法與不合法 schema 檔案的解析與驗證。

---

### 整合測試場景

#### plugin-importer-versioning 整合流程
- [ ] internal/test/integration/full_import_flow_test.go：
  - 模擬 plugin 掛載後註冊 importer
  - 模擬匯入 JSON/YAML → 轉換成 Element → 寫入 registry
  - 寫入 version store 並驗證版本記錄存在

#### bootstrap-wire 啟動流程
- [ ] internal/test/integration/bootstrap_test.go：
  - 測試完整 Init → BuildServer → Run 是否可完成啟動
  - 使用 fake_component 注入並驗證模組運作順序與日誌


## 9. /internal/adapters 與對應 interface 模組缺漏項目補充
interface: pkg/ifaces/{module}/*
adapter: internal/adapters/{module}/*

- [x] internal/adapters/eventbus/metric.go 測試缺失，需補對應 adapter 測試。
- [x] internal/adapters/logger/zap_adapter.go 未補對應單元測試。
- [x] internal/test/fakes/fake_logger.go 尚未建立，用於測試注入與日誌檢查。
- [x] internal/test/fakes/fake_notifier.go 尚未建立，用於測試多 notifier 整合。
- [x] internal/adapters/notifier/email_adapter_test.go 尚未補單元測試。
- [x] internal/adapters/eventbus/host.go 測試缺失，需補對應 adapter 測試。
- [x] internal/adapters/eventbus/task.go 測試缺失，需補對應 adapter 測試。
- [x] internal/adapters/eventbus/alert.go 測試缺失，需補 alert handler 對應事件處理測試。
- [x] internal/test/fakes/fake_eventbus.go 尚未建立，供 plugin 匯入測試使用。
- [x] internal/adapters/notifier/slack_adapter.go 未補單元測試。
- [x] internal/adapters/notifier/webhook_adapter.go 未補單元測試。
- [x] internal/adapters/logger/nop_adapter.go 未補單元測試。
 
## 10. 測試與樣板補齊: 確保前 4 步有可運行測試與測試替身
- [ ] pkg/ifaces/test/test.go：模組與插件測試 interface。可定義模組驗證器、模擬元件等輔助介面，簡化測試撰寫。
- [ ] internal/adapters/eventbus/alert.go 測試案例尚未建立。需設計單元測試覆蓋 alert eventbus handler 的各種情境。
- [ ] internal/plugins/eventbus/* 所有插件 handler 測試缺失。需針對 eventbus plugins 撰寫 handler測試案例。
- [ ] internal/test/plugins/：應建立 plugin 驗證測試樣板。需提供範例與基礎測試樣板，便於插件開發者驗證功能。
- [ ] pkg/ifaces/event/... 應補 mock 與泛用事件 payload 檢驗。需設計 mock event bus 與事件 payload 驗證工具。
- [ ] internal/test/testutil/assert_logger.go 功能可再擴充（包含結構化 log 比對）。需擴充 logger 測試工具，支援結構化日誌內容比對。
- [ ] docs/interfaces/test.md：測試樣板介面與模組說明文件。補充模組測試介面、mock 建立與 handler 驗證策略。
- [ ] internal/test/fakes/fake_server.go：提供 Server interface 假實作，用於測試注入與驗證呼叫。
- [ ] internal/test/server/server_test.go：測試 Server Run/Shutdown 邏輯與模組注入整合，驗證是否能正常啟動 HTTP 與模組流程。
- [x] internal/test/fakes/fake_config.go：提供 ConfigProvider 假實作，支援測試注入與設定模擬。
