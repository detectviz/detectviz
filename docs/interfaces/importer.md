

# Importer Interface 說明

`pkg/importer/interface.go` 定義 Importer 介面，作為 detectviz 系統中的資源匯入模組擴充入口。每個 Importer 可負責處理特定類型的檔案，並轉換為符合 `registry.Resource` 的結構供後續註冊使用。

## Interface 定義

```go
type Importer interface {
  Name() string
  GVK() registry.GVK
  Load(ctx context.Context, filePath string) ([]registry.Resource, error)
}
```

### 方法說明

| 方法        | 回傳類型                     | 說明                                |
|-------------|------------------------------|-------------------------------------|
| `Name()`    | `string`                     | 回傳此匯入器的唯一識別名稱         |
| `GVK()`     | `registry.GVK`               | 匯入器對應的資源類型               |
| `Load()`    | `[]registry.Resource, error` | 載入指定檔案，解析為資源結構清單   |

## 使用情境

- 從 JSON/YAML 檔案載入一批 `Host` 資源
- 轉換第三方監控系統匯出的格式為內部資源
- 載入 CSV 並轉換為 `Datasource` 結構

## 設計原則

- 每個 Importer 實作應僅處理一組 GVK
- 應明確分離資料來源與資料格式解析責任
- 可由 plugin 或系統模組透過統一介面載入匯入器模組

## 延伸說明

Importer 模組將可與 plugin loader、版本控制（versioning）、系統配置結合，實現模組化可擴充的資源初始化與載入機制。