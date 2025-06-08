# ConfigProvider Interface 說明文件

> 本文件為 Detectviz 專案中 `ConfigProvider` 介面的設計說明與使用情境整理。

## 介面用途（What it does）

ConfigProvider 是 Detectviz 後端平台中統一存取設定值的標準介面。其設計目標為：

- 解耦設定來源與業務邏輯
- 避免使用全域變數與硬編碼
- 提供 hot-reload 支援與擴充性（YAML, DB, 遠端等）
- 易於 mock、測試與插件化

## 使用情境（When and where it's used）

- 在 `internal/bootstrap/init.go` 中初始化並注入各模組
- 各子模組透過 `Get`, `GetBool` 等方式查詢設定值
- 支援後續單元測試與動態切換設定來源

## 方法說明（Methods）

- `Get(key string) string`：回傳指定 key 對應的字串設定值
- `GetInt(key string) int`：回傳整數型別設定值
- `GetBool(key string) bool`：回傳布林值設定值
- `GetOrDefault(key, defaultVal string) string`：回傳 key 設定值，若無則回傳預設值
- `Reload() error`：重新載入設定來源（若實作支援），可用於 hot-reload

## 預期實作（Expected implementations）

- `pkg/config/default.go`：預設記憶體實作，支援環境變數與 map
- `internal/adapters/config/yaml.go`（可選擴充）：支援 YAML 檔案讀取
- `internal/adapters/config/remote.go`（可選擴充）：支援 HTTP 或 config service

## 關聯模組與擴充性（Related & extensibility）

- 可搭配 `watcher` 類型模組監控設定變更（future）
- 可與遠端設定管理服務（如 Consul、etcd）整合
- interface 抽象化允許單元測試時注入 mock 設定來源