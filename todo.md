# scaffold: 初始目錄與範本建立

> 備註：原本 detectviz 舊版目錄搬移至 docs/detectviz-deprecated

- [Cursor Scaffold 指引補充](docs/README.md)

## scaffold: 完成元件實作

- [x] 建立基礎目錄結構
  - `apps/` ✅
  - `pkg/` ✅
  - `internal/` ✅
  - `plugins/` ✅
  - `compositions/` ✅
  - `scripts/` ✅
  - `tools/` ✅
  - `docs/` ✅

- [x] 建立 plugins scaffold
  - `plugins/core/auth/jwt/plugin.go` ✅
  - `plugins/core/middleware/logging/plugin.go` ⏳
  - `plugins/community/importers/prometheus/plugin.go` ✅
  - `plugins/community/exporters/influxdb/plugin.go` ⏳
  - `plugins/community/integrations/security/keycloak/plugin.go` ⏳
  - `plugins/tools/supportbundles/plugin.go` ⏳

- [x] 建立 contracts interface 定義
  - `pkg/platform/contracts/plugin.go` ✅
  - `pkg/platform/contracts/importers.go` ✅
  - `pkg/platform/contracts/exporters.go` ✅
  - `pkg/platform/contracts/auth.go` ✅
  - `pkg/platform/contracts/lifecycle.go` ✅

- [x] 建立組合與註冊框架
  - `internal/platform/registry/registry.go` ✅
  - `internal/platform/composition/lifecycle.go` ✅
  - `internal/platform/composition/resolver.go` ⏳
  - `pkg/config/loader/config_loader.go` ⏳

- [x] 建立最小平台組合檔
  - `compositions/minimal-platform/composition.yaml` ✅
  - `compositions/alloy-devkit/alloy-config.yaml` ✅

- [x] 建立 scaffold 註解 README
  - `plugins/core/auth/jwt/README.md` ✅
  - `plugins/community/importers/prometheus/README.md` ✅
  - 每個 plugin 目錄包含 `README.md` 說明該 plugin 類型與功能範本說明

## scaffold: 未完成元件實作補強（詳細）

- [ ] 實作 LifecycleManager 核心啟動流程
  - 檔案：`internal/platform/composition/lifecycle.go`
  - 區段：lines 228–247
  - 內容：
    - 補齊 `StartAll()`, `ShutdownAll()`, `HealthCheck()` 函式
    - 對已註冊的 plugin 執行 lifecycle hook：`OnStart()`, `OnShutdown()`
    - 檢查是否為 `LifecycleAware` 再呼叫對應函式

- [ ] 新增組合依賴解析器
  - 檔案：`internal/platform/composition/resolver.go`
  - 內容：
    - 定義 plugin 組合結構與 `depends_on` 欄位
    - 建立拓撲排序邏輯（Topological Sort）
    - 回傳排序後的啟動順序供 `LifecycleManager` 使用

- [ ] 實作 YAML 設定載入器
  - 檔案：`pkg/config/loader/config_loader.go`
  - 內容：
    - 讀取 `composition.yaml`
    - 將每個 plugin 的 config 區塊對應傳遞給該 plugin 實例
    - 支援 `map[string]any`、`yaml.Unmarshal` 或 `viper` 自動解碼

- [ ] 建立最小 Server 入口範例
  - 檔案：`apps/server/main.go`
  - 內容：
    - 建立 registry、載入組合、啟動 plugin lifecycle
    - 印出啟動成功、錯誤日誌
    - 最小平台組合使用：`compositions/minimal-platform/composition.yaml`

- [ ] 修正 plugin 忽略設定參數問題
  - 檔案：`plugins/core/auth/jwt/plugin.go` → `NewJWTAuthenticator(cfg any)`
  - 檔案：`plugins/community/importers/prometheus/plugin.go` → `NewPrometheusImporter(cfg any)`
  - 內容：
    - 應將 `cfg any` decode 成結構體，解析 `config["issuer"]`, `config["targets"]` 等欄位
    - 建議使用 `mapstructure.Decode()` 或自定 `ParsePluginConfig()`
    - 記得補上預設值 fallback 與格式錯誤處理

## scaffold: 功能擴充與測試驗證

- [ ] 實作 plugin config schema 驗證器
  - 檔案：`pkg/config/schema/validator.go`
  - 內容：
    - 接收 `map[string]any` 或結構體，執行欄位格式與預設值驗證
    - 可使用 `mapstructure`, `go-playground/validator`, `yaml` 等方式

- [ ] 增加 plugin 健康狀態查詢函式
  - 檔案：`internal/platform/composition/lifecycle.go`
  - 內容：
    - `HealthCheck()` 支援 plugin 回報健康狀態（若實作 `HealthAware` interface）
    - 彙整狀態供外部查詢

- [ ] 增加 enable/disable 機制到 composition.yaml 支援
  - 檔案：`composition.yaml`, `registry.go`
  - 內容：
    - 若 `enabled: false` 則跳過註冊與啟動該 plugin
    - Registry 與 LifecycleManager 均須處理

- [ ] 建立 plugin discovery 機制
  - 檔案：`internal/platform/registry/discovery.go`
  - 內容：
    - 掃描 `plugins/**/plugin.go`
    - 自動呼叫其 `Register()`，支援熱插拔擴充用

- [ ] 增加 scaffold integration test
  - 檔案：`internal/test/integration/scaffold_test.go`
  - 內容：
    - 載入 minimal composition → 組合 → 啟動 plugins → 驗證運作順序與 lifecycle 正確