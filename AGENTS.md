# DetectViz Agent 指南

> 此文件是為 Codex 和其他 AI 助理提供的專案說明書，幫助理解 DetectViz 的架構、開發規範和協作模式。

## 專案架構概述

DetectViz 是一個可組合的可觀測性平台，採用插件化架構，支援靈活的功能擴展和組合：

### 核心分層結構

```
DetectViz Platform
├── Apps Layer (apps/)             # 應用組合層
│   ├── server/                   # HTTP/gRPC 服務應用
│   ├── cli/                      # 命令列管理工具
│   ├── agent/                    # 分散式收集代理
│   └── testkit/                  # 測試工具組合
│
├── Platform Core (pkg/platform/) # 平台核心抽象
│   ├── contracts/                # 跨模組契約定義
│   ├── registry/                 # 插件註冊機制
│   └── composition/              # 組合引擎
│
├── Internal Platform (internal/) # 平台實作層
│   ├── platform/                 # 核心平台實作
│   ├── ports/                    # 輸入輸出介面
│   └── services/                 # 業務服務實作
│
├── Plugins Ecosystem (plugins/)  # 插件生態系統
│   ├── core/                     # 核心插件
│   │   ├── auth/                 # 認證插件
│   │   ├── storage/              # 存儲插件
│   │   └── alerting/             # 告警插件
│   │
│   ├── community/                # 社群插件
│   │   ├── importers/            # 資料匯入器
│   │   ├── exporters/            # 資料匯出器
│   │   ├── integrations/         # 第三方整合
│   │   └── tools/                # 工具類插件
│   │
│   └── extensions/               # 擴展插件
│       ├── analytics/            # 分析功能
│       └── web/                  # Web UI 擴展
│
└── Compositions (compositions/)   # 組合方案
    ├── minimal-platform/         # 最小平台組合
    ├── monitoring-stack/         # 監控堆疊組合
    └── alloy-devkit/             # Alloy 開發套件
```

### 核心概念

1. **插件 (Plugin)**: 實作特定功能的獨立模組，支援動態載入和配置
2. **契約 (Contract)**: 定義插件間互動的介面規範
3. **組合 (Composition)**: 將多個插件組合成完整的應用方案
4. **註冊 (Registry)**: 管理插件的註冊、發現和依賴解析
5. **生命週期 (Lifecycle)**: 控制插件的初始化、啟動、停止和關閉

## 插件開發指南

### 插件類型與介面

DetectViz 支援多種類型的插件，每種類型實作對應的介面：

1. **資料匯入器 (Importer)**
   - 介面定義: `pkg/platform/contracts/importers.go`
   - 實作範例: `plugins/community/importers/prometheus/`
   - 用途: 從外部系統收集資料

2. **資料匯出器 (Exporter)**
   - 介面定義: `pkg/platform/contracts/exporters.go`
   - 實作範例: `plugins/community/exporters/influxdb/`
   - 用途: 將資料輸出到外部系統

3. **認證器 (Authenticator)**
   - 介面定義: `pkg/platform/contracts/auth.go`
   - 實作範例: `plugins/core/auth/jwt/`, `plugins/community/integrations/security/keycloak/`
   - 用途: 提供身份認證和授權功能

4. **Web UI 插件 (WebUIPlugin)**
   - 介面定義: `pkg/platform/contracts/webplugin.go`
   - 實作範例: `plugins/extensions/web/systemstatus/`
   - 用途: 擴展 Web 介面功能

5. **中介層 (Middleware)**
   - 實作範例: `plugins/tools/middleware/requestmeta/`
   - 用途: 提供橫切面功能（如日誌、監控、權限檢查）

### 插件開發模式

#### 基本插件結構

```go
package myplugin

import (
    "context"
    "detectviz/pkg/platform/contracts"
)

type MyPlugin struct {
    name        string
    version     string
    description string
    config      *Config
    initialized bool
    started     bool
}

type Config struct {
    Setting1 string `yaml:"setting1" json:"setting1"`
    Setting2 int    `yaml:"setting2" json:"setting2"`
}

func NewMyPlugin(config any) (contracts.Plugin, error) {
    // 解析配置
    // 設定預設值
    // 驗證必要參數
    // 回傳插件實例
}

// 實作基本 Plugin 介面
func (p *MyPlugin) Name() string { return p.name }
func (p *MyPlugin) Version() string { return p.version }
func (p *MyPlugin) Description() string { return p.description }
func (p *MyPlugin) Init(config any) error { /* 初始化邏輯 */ }
func (p *MyPlugin) Shutdown() error { /* 清理邏輯 */ }

// 實作 LifecycleAware 介面（如需要）
func (p *MyPlugin) OnRegister() error { /* 註冊時邏輯 */ }
func (p *MyPlugin) OnStart() error { /* 啟動時邏輯 */ }
func (p *MyPlugin) OnStop() error { /* 停止時邏輯 */ }
func (p *MyPlugin) OnShutdown() error { /* 關閉時邏輯 */ }

// 實作 HealthChecker 介面（如需要）
func (p *MyPlugin) CheckHealth(ctx context.Context) contracts.HealthStatus {
    // 檢查插件健康狀態
}

// 實作特定功能介面（如 Importer, Exporter 等）

// 註冊函式
func Register(registry contracts.Registry) error {
    return registry.RegisterPlugin("my-plugin", NewMyPlugin)
}
```

#### 配置解析模式

由於專案不使用外部依賴（如 mapstructure），採用手動解析配置：

```go
func parseConfig(config any, target *Config) error {
    if config == nil {
        return nil
    }

    configMap, ok := config.(map[string]any)
    if !ok {
        return fmt.Errorf("config must be a map[string]any")
    }

    if setting1, exists := configMap["setting1"]; exists {
        if str, ok := setting1.(string); ok {
            target.Setting1 = str
        }
    }
    
    if setting2, exists := configMap["setting2"]; exists {
        if intVal, ok := setting2.(int); ok {
            target.Setting2 = intVal
        }
    }

    return nil
}
```

### 插件放置規則

- **核心插件**: `plugins/core/{category}/{name}/plugin.go`
  - 例: `plugins/core/auth/jwt/plugin.go`
  
- **社群插件**: `plugins/community/{category}/{subcategory}/{name}/plugin.go`
  - 例: `plugins/community/importers/prometheus/plugin.go`
  - 例: `plugins/community/exporters/influxdb/plugin.go`
  - 例: `plugins/community/integrations/security/keycloak/plugin.go`

- **擴展插件**: `plugins/extensions/{category}/{name}/plugin.go`
  - 例: `plugins/extensions/web/systemstatus/plugin.go`

## 測試策略與最佳實務

### 測試分層

1. **單元測試**: 測試插件核心邏輯
   - 檔案命名: `plugin_test.go`
   - 測試配置解析、業務邏輯、錯誤處理

2. **整合測試**: 測試插件間互動
   - 位置: `internal/test/integration/`
   - 測試範例: `scaffold_test.go`, `webui_test.go`, `config_validation_test.go`

3. **系統測試**: 測試完整組合
   - 使用 `compositions/` 中的配置進行端到端測試

### 測試原則

1. **隔離性**: 每個測試獨立運行，不依賴其他測試
2. **可重複性**: 測試結果穩定，不受環境影響
3. **覆蓋性**: 涵蓋正常路徑、異常路徑和邊界條件
4. **真實性**: 使用真實的 HTTP 服務器和資料庫連接進行測試

### 測試工具與模式

```go
// 使用 httptest.NewServer 測試 HTTP 互動
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // 模擬外部服務回應
}))
defer server.Close()

// 使用 context 控制測試超時
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// 測試健康檢查
healthStatus := plugin.CheckHealth(ctx)
if healthStatus.Status != "healthy" {
    t.Errorf("Expected healthy status, got %s", healthStatus.Status)
}
```

## Web UI 開發指南

DetectViz 採用 HTMX + Echo SSR + Templ 技術棧：

### Web UI 插件開發

1. **註冊路由**: 透過 `RegisterRoutes` 方法註冊 HTTP 路由
2. **註冊導覽**: 透過 `RegisterNavTreeNodes` 方法添加導覽項目
3. **註冊組件**: 透過 `RegisterComponents` 方法添加 UI 組件

```go
func (p *MyWebPlugin) RegisterRoutes(router contracts.WebRouter) error {
    router.GET("/my-plugin/status", p.handleStatus)
    router.POST("/api/my-plugin/action", p.handleAction)
    return nil
}

func (p *MyWebPlugin) RegisterNavTreeNodes() []contracts.NavTreeNode {
    return []contracts.NavTreeNode{
        {
            ID:    "my-plugin",
            Title: "My Plugin",
            URL:   "/my-plugin/status",
            Icon:  "fas fa-cog",
        },
    }
}
```

### 前端技術棧

- **HTMX**: 實作動態更新，無需複雜 JavaScript
- **AdminLTE**: UI 框架，提供一致的視覺風格
- **Templ**: 型別安全的 HTML 範本生成
- **Tabulator.js**: 表格視覺化組件

## 可觀測性整合

### OpenTelemetry 整合

DetectViz 內建 OpenTelemetry 支援，插件應使用標準的追蹤和指標收集：

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func (p *MyPlugin) ProcessData(ctx context.Context, data []byte) error {
    ctx, span := otel.Tracer("my-plugin").Start(ctx, "process-data")
    defer span.End()
    
    // 業務邏輯處理
    
    return nil
}
```

### Grafana Alloy 整合

平台整合 Grafana Alloy 作為統一的可觀測性代理：

- **配置管理**: 透過 `alloy-config.river` 管理監控配置
- **OTLP 支援**: 完整支援 OpenTelemetry Protocol
- **多後端支援**: 支援 Tempo (追蹤)、Loki (日誌)、Mimir (指標)

## 配置與組合

### 插件配置

所有插件配置透過 `composition.yaml` 統一管理：

```yaml
plugins:
  - name: my-plugin
    type: importer
    enabled: true
    config:
      setting1: "value1"
      setting2: 42
      timeout: "30s"
```

### 配置驗證

插件應支援配置模式驗證：

```go
// pkg/config/schema/validator.go 中註冊模式
validator.RegisterSchema("my-plugin", &schema.PluginSchema{
    Name: "my-plugin",
    Fields: map[string]*schema.FieldSchema{
        "setting1": {
            Type:        "string",
            Required:    true,
            Description: "Setting 1 description",
        },
        "setting2": {
            Type:        "int",
            Required:    false,
            Default:     42,
            Description: "Setting 2 description",
        },
    },
})
```

## 錯誤處理與日誌

### 錯誤處理模式

```go
// 使用有意義的錯誤訊息
if err != nil {
    return fmt.Errorf("failed to process data: %w", err)
}

// 區分可恢復和不可恢復錯誤
type RecoverableError struct {
    Err error
}

func (e RecoverableError) Error() string {
    return e.Err.Error()
}

func (e RecoverableError) Unwrap() error {
    return e.Err
}
```

### 日誌記錄

```go
// 使用結構化日誌
fmt.Printf("Plugin %s started successfully on %s\n", p.name, endpoint)
fmt.Printf("Plugin %s failed to start: %v\n", p.name, err)

// 包含相關上下文資訊
fmt.Printf("Processing batch of %d items for plugin %s\n", len(batch), p.name)
```

## 效能與資源管理

### 連接管理

```go
// 使用 HTTP 客戶端池
httpClient := &http.Client{
    Timeout: time.Duration(config.Timeout) * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

// 正確關閉資源
defer resp.Body.Close()
```

### 批次處理

```go
// 支援批次處理以提升效能
func (p *MyPlugin) BatchProcess(ctx context.Context, batch []interface{}) error {
    if len(batch) == 0 {
        return nil
    }
    
    // 批次處理邏輯
    for _, item := range batch {
        if err := p.processItem(ctx, item); err != nil {
            return fmt.Errorf("failed to process item: %w", err)
        }
    }
    
    return nil
}
```

## 開發工作流程

### 新插件開發流程

1. **確認需求**: 閱讀 `todo.md` 了解待完成項目
2. **選擇介面**: 根據功能需求選擇實作的 contracts 介面
3. **創建檔案**: 按照目錄規則創建插件檔案
4. **實作功能**: 按照介面定義實作核心功能
5. **撰寫測試**: 創建對應的測試檔案
6. **配置驗證**: 添加配置模式驗證（如需要）
7. **文件更新**: 更新相關文件

### 程式碼品質檢查

開發完成後，請執行以下檢查：

1. **語法檢查**: `go build ./...`
2. **測試執行**: `go test ./...`
3. **格式化**: `go fmt ./...`
4. **配置驗證**: 確保插件配置能正確解析和驗證

### 除錯與故障排除

1. **健康檢查**: 實作 `HealthChecker` 介面提供狀態資訊
2. **指標收集**: 實作 `GetHealthMetrics` 方法提供監控資料
3. **錯誤日誌**: 記錄詳細的錯誤資訊和上下文
4. **連接測試**: 在初始化時測試外部服務連接

## 協作指南

### 與 Codex 協作

當使用 Codex 進行開發時：

1. **明確任務**: 提供具體的功能需求和實作要求
2. **參考現有**: 指向相似的現有插件作為參考
3. **測試要求**: 要求生成對應的測試檔案
4. **配置驗證**: 要求添加配置模式驗證支援

### 程式碼審查要點

1. **介面遵循**: 確保正確實作所需的 contracts 介面
2. **錯誤處理**: 檢查錯誤處理是否完整和有意義
3. **資源管理**: 確認連接和資源正確釋放
4. **測試覆蓋**: 驗證測試是否覆蓋主要功能路徑
5. **文件一致**: 確保程式碼與文件描述一致

## 常見問題與解決方案

### Q: 如何實作一個新的資料匯入器？

A: 參考 `plugins/community/importers/prometheus/plugin.go`，實作 `contracts.Importer` 介面，包含 `Import` 方法和生命週期管理。

### Q: 如何添加 Web UI 功能？

A: 實作 `contracts.WebUIPlugin` 介面，註冊路由、導覽節點和 UI 組件。參考 `plugins/extensions/web/systemstatus/plugin.go`。

### Q: 如何處理插件配置？

A: 使用手動解析配置的模式，並在 `pkg/config/schema/validator.go` 中註冊配置模式進行驗證。

### Q: 如何進行插件測試？

A: 參考 `internal/test/integration/` 中的測試範例，使用 `httptest.NewServer` 進行 HTTP 測試，使用 context 控制超時。

---

此文件將持續更新，以反映 DetectViz 架構的最新變化和最佳實務。如有疑問，請參考相關的介面文件和現有實作範例。