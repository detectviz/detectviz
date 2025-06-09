# Plugin Interface 說明

`Plugin` 為 detectviz 所有外掛模組的標準化通用介面，定義於 `pkg/ifaces/plugins/plugin.go`。所有可掛載模組皆應實作該介面以便於生命週期管理與版本控管。

## Plugin Interface 定義

```go
type Plugin interface {
  Name() string
  Version() string
  Init() error
  Close() error
}
```

### 方法說明：

| 方法     | 回傳型別 | 說明                                         |
|----------|----------|----------------------------------------------|
| `Name()` | `string` | 插件唯一識別名稱，用於註冊與日誌識別等用途     |
| `Version()` | `string` | 插件版本（字串格式，可與 Versioning 比較） |
| `Init()` | `error`  | 啟用插件邏輯，進行初始化與依賴注入             |
| `Close()` | `error` | 關閉插件時釋放資源與終止任務                   |

---

## 使用情境

- 每個插件皆需明確實作 `Plugin` 介面，支援註冊、掛載與卸載操作
- 可透過 Registry 檢查版本差異、註冊至 Plugin Manager 統一管理
- 可搭配 eventbus, metrics, notifier 等模組建立 plugin 擴充架構

---

## 延伸設計

Plugin Interface 可作為以下模組的基礎：

- `AlertPlugin`, `DataSourcePlugin`, `RendererPlugin` 等具類型擴充的插件接口
- 支援動態掃描 `/internal/plugins/` 或 `/plugins/` 目錄註冊 Plugin 實例
- 可與版本控管模組整合，檢查相容性與載入優先順序

---

## 測試支援與假件

為方便測試 Plugin 註冊與生命週期控制，系統建議實作 `plugins.Plugin` 介面的簡易 mock plugin，搭配 PluginLifecycleManager、ManagerRegistry 等模組進行單元測試與整合驗證。

### 假件範例：

```go
type mockPlugin struct {
  name    string
  version string
  inited  bool
  closed  bool
}

func (p *mockPlugin) Name() string    { return p.name }
func (p *mockPlugin) Version() string { return p.version }
func (p *mockPlugin) Init() error     { p.inited = true; return nil }
func (p *mockPlugin) Close() error    { p.closed = true; return nil }
```

可用於以下模組測試：

- `internal/plugins/manager/registry_test.go`
- `internal/plugins/manager/lifecycle_test.go`
- 動態註冊測試、Init/Close 執行驗證等

此測試模式支援 plugin 模組邏輯抽換與邊界行為驗證，便於未來新增 plugin 類型時複用。