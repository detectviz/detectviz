# ConfigProvider Interface 設計說明

> 本文件為 Detectviz 專案中 `ConfigProvider` 介面的設計原則、使用情境與實作擴充方式整理，並統一與其他 interface 文件格式。

- 統一管理所有模組設定結構於 `pkg/configtypes`，支援靜態與 plugin 配置

---

## 介面用途（What it does）

`ConfigProvider` 是 Detectviz 平台中用於讀取設定值的統一抽象介面，其目標如下：

- 解耦設定讀取來源與核心邏輯（YAML、ENV、Remote 等）
- 避免硬編碼與全域變數污染
- 支援 hot-reload（視實作而定）
- 提供良好測試性（可注入 Mock 設定來源）
- 支援擴充（如套用 Config Schema、版本控制）

---

## 使用情境（When and where it's used）

- 在 `internal/bootstrap/init.go` 中初始化後注入各核心模組
- 子模組可透過 `Get`、`GetInt`、`GetBool` 取得配置參數
- 搭配 logger、notifier、scheduler 等模組進行動態設定注入
- 測試時可替換為 Fake 或 Map-based 實作

---

## 方法定義（Interface Methods）

```go
type ConfigProvider interface {
    Get(key string) string
    GetInt(key string) int
    GetBool(key string) bool
    GetOrDefault(key string, defaultVal string) string
    GetCacheConfig() configtypes.CacheConfig
    GetNotifierConfigs() []configtypes.NotifierConfig
    Logger() logger.Logger
    Reload() error
}
```

- `Get`：傳回指定 key 的字串值，若不存在回空字串
- `GetInt`：傳回整數值，無法解析時預設為 0
- `GetBool`：傳回布林值，支援 "true"/"false" 字串轉換
- `GetOrDefault`：若指定 key 無值，則傳回提供的 default 值
- `GetCacheConfig`：回傳快取模組所需的結構設定
- `GetNotifierConfigs`：回傳通知通道模組的設定清單
- `Logger`：回傳 logger 實例，供模組共用
- `Reload`：重新載入設定來源，若支援 hot-reload 機制，否則為 no-op

---

## 模組設定結構存放位置

Detectviz 中各模組的設定結構（如 NotifierConfig、CacheConfig）應統一定義於：

```go
pkg/configtypes/
```

每個模組對應一個 `{name}_config.go`，便於：

- 保持 interface 與資料結構分離
- 自動對應 `config.UnmarshalKey()` 載入模組配置
- 提供良好 IDE 導引與文件生成支援

### 設定結構總覽（Defined Config Structs）

| 結構名稱             | 檔案位置                               | 用途與對應模組                                     |
|----------------------|----------------------------------------|----------------------------------------------------|
| `LoggerConfig`       | `pkg/configtypes/logger_config.go`     | 設定 logger 類型、Level、輸出格式等                |
| `SchedulerConfig`    | `pkg/configtypes/scheduler_config.go`  | 設定排程器類型與參數                              |
| `NotifierConfig`     | `pkg/configtypes/notifier_config.go`   | ✅ 已完成：定義通知方式（email/slack/webhook）     |
| `CacheConfig`        | `pkg/configtypes/cache_config.go`      | ✅ 已完成：定義快取實作方式與參數（memory/redis） |
| `AlertConfig`        | `pkg/configtypes/alert_config.go`      | 告警模組預設行為與閾值設定                         |
| `BusConfig`          | `pkg/configtypes/bus_config.go`        | 指定使用哪一種 EventBus 實作                       |
| `MetricsConfig`      | `pkg/configtypes/metrics_config.go`    | 設定查詢資料來源（prometheus/flux）等             |
| `WebConfig`          | `pkg/configtypes/web_config.go`        | 🟡 預計實作：設定 HTMX 前端與靜態資源               |
| `InfraConfig`        | `pkg/configtypes/infra_config.go`      | 🟡 預計實作：底層 Infra（如 HTTP Client, Retry）    |
| `EncryptionConfig`   | `pkg/configtypes/encryption_config.go` | 🟡 預計實作：加密金鑰、演算法與 key rotation 設定  |
| `ValidationConfig`   | `pkg/configtypes/validation_config.go` | 🟡 預計實作：輸入資料驗證控制與規則                |
| `PluginConfig[T any]`| `pkg/configtypes/plugin_config.go`     | 提供 plugin 專屬結構包裝與載入                     |

### PluginConfig[T] 用法

若需為 plugin（如 `auth.keycloak`, `store.redis`）定義專屬設定結構，可透過以下方式實作：

```go
type KeycloakConfig struct {
  URL string `mapstructure:"url" yaml:"url"`
}

type PluginConfig[T any] struct {
  Config T
}
```

搭配：

```go
config.UnmarshalKey("plugin.auth.keycloak", &PluginConfig[KeycloakConfig]{})
```

### 測試建議（Testing Strategy）

- 建立對應 `*_test.go` 檔案測試配置反序列化與 fallback
- 使用 `testutil.LoadConfigFromYAML` 驗證整體載入行為
- PluginConfig 可獨立以泛型方式驗證任意結構

### 擴充建議（Extensions）

- 所有 struct 欄位應標註 `yaml/json/mapstructure` tag，建議格式：

```go
type Example struct {
  Enabled bool `yaml:"enabled" json:"enabled" mapstructure:"enabled"`
}
```

- 每新增模組應同步建立對應 `{module}_config.go`
- 可進一步對接設定 Schema 驗證與 UI 編輯工具

---

## 預期實作（Expected Implementations）

| 類型     | 路徑位置                                   | 描述                            |
|----------|--------------------------------------------|---------------------------------|
| 預設     | `pkg/config/default.go`                    | 使用 map + ENV 實作設定讀取器   |
| YAML     | `internal/adapters/config/yaml.go`         | 從指定 YAML 檔案讀取設定         |
| Remote   | `internal/adapters/config/remote.go`       | 支援 HTTP / gRPC 動態設定服務   |
| Nop/Fake | `internal/adapters/config/mock_adapter.go`<br>`internal/test/fakes/fake_config.go` | 提供測試用途的空實作或假資料     |

---

## 測試建議（Testing Strategy）

- 可透過注入 `FakeConfigProvider` 模擬錯誤或邊界值
- 建議單元測試涵蓋 `Reload` 行為與 fallback 邏輯
- 可加入 `Set()` 方法以便測試程式中直接設定參數（建議僅於測試實作中使用）
- 結合整合測試驗證模組是否正確依據 config 行為切換

---

## 擴充與整合建議（Extensions & Integration）

- 可整合 Config Schema，支援欄位驗證與版本轉換
- 搭配熱更新模組（如 fsnotify）實現自動 reload
- 與遠端設定中心（如 Consul、etcd、Spring Cloud Config）整合
- 規劃支援 json-schema-version，可提升未來 JSON 設定相容性管理

---

---
（本文已整併 configtypes.md 內容）