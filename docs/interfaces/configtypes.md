

# ConfigTypes 設定結構說明文件

> 本文件說明 Detectviz 專案中 `pkg/configtypes` 目錄內各設定結構（Config Struct）的用途、對應模組與設定來源。這些結構主要用於從靜態或動態設定載入模組所需參數，並透過 ConfigProvider 傳遞至模組內部。

---

## 設計目的（Design Purpose）

- 將每個功能模組所需的設定集中管理
- 可清楚對應設定來源（如 JSON / YAML / ENV）
- 搭配 ConfigProvider 實現動態注入或切換
- 提供 Registry 與 Bootstrap 使用的統一結構

---

## 設定結構總覽（Defined Config Structs）

| 結構名稱           | 檔案位置                           | 用途與對應模組                       |
|--------------------|------------------------------------|--------------------------------------|
| `LoggerConfig`     | `pkg/configtypes/logger_config.go` | 設定 logger 類型、Level、輸出格式等 |
| `SchedulerConfig`  | `pkg/configtypes/scheduler_config.go` | 設定排程器類型與參數                |
| `NotifierConfig`   | `pkg/configtypes/notifier_config.go` | 定義通知方式（email/slack/webhook） |
| `CacheConfig`      | `pkg/configtypes/cache_config.go`  | 定義快取實作方式與參數（memory/redis） |
| `AlertConfig`      | `pkg/configtypes/alert_config.go`  | 告警模組預設行為與閾值設定           |
| `BusConfig`        | `pkg/configtypes/bus_config.go`    | 指定使用哪一種 EventBus 實作         |
| `MetricsConfig`    | `pkg/configtypes/metrics_config.go` | 設定查詢資料來源（prom/flux）等     |

---

## 使用方式（Usage Pattern）

這些結構通常會出現在 `default.go` 或 `ConfigProvider` 的初始化中，例如：

```go
config := &configtypes.AppConfig{
    Logger:   configtypes.LoggerConfig{Level: "debug"},
    Scheduler: configtypes.SchedulerConfig{Type: "workerpool"},
}
```

或從 YAML / JSON 讀取時：

```yaml
logger:
  type: zap
  level: info

scheduler:
  type: cron
  spec: "0 * * * *"
```

---

## 測試建議（Testing Strategy）

- 可建立對應 `*_test.go` 測試配置反序列化行為
- 搭配 `testutil.LoadConfigFromYAML` 測試整體載入是否成功
- 支援 fallback 至預設值邏輯可單元測試驗證

---

## 擴充建議（Extensions）

- 若新增模組應同步建立對應 ConfigTypes 檔案（命名規則：`{module}_config.go`）
- 結構欄位建議皆標記 yaml/json tag
- 可未來對接 Config Schema 驗證或 UI 編輯介面

---