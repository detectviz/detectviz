# ConfigProvider Interface 設計說明

> 本文件為 Detectviz 專案中 `ConfigProvider` 介面的設計原則、使用情境與實作擴充方式整理，並統一與其他 interface 文件格式。

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