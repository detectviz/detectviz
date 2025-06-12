# DetectViz 組合式架構目錄命名標準

## 架構特性分析

DetectViz 採用 **組合式架構 (Composable Architecture)** 設計，具備以下特點：
1. **Clean Architecture 分層**：明確的職責分離和依賴方向
2. **Plugin-First 設計**：核心功能通過插件系統擴展
3. **Platform-as-Code**：架構本身作為平台，支援快速組裝
4. **模組化組合**：不同應用可選擇性組合所需模組

## 當前問題重新分析

基於組合式架構的需求，目錄結構問題如下：
1. **插件邊界不清**：Plugin 相關目錄分散在不同層級
2. **平台層缺失**：缺乏明確的平台抽象層
3. **組合邏輯混亂**：模組組合和註冊邏輯分散
4. **可組合性差**：難以快速組裝新的應用組合

## 組合式架構重構建議

### 核心設計原則
1. **平台優先**：Platform 層提供核心組合能力
2. **插件驅動**：功能通過 Plugin 系統擴展
3. **組合透明**：模組組合邏輯清晰可見
4. **契約明確**：介面契約統一定義

### 重構後目錄結構

```bash
detectviz/
├── apps/                    # 🎯 應用組合層（組裝不同的應用）
│   ├── server/              # Web API 伺服器應用
│   ├── cli/                 # CLI 工具應用  
│   ├── agent/               # 監控代理應用
│   └── testkit/             # 測試工具應用
│
├── pkg/                     # 📋 公共契約與平台介面
│   ├── platform/            # 🏗️ 平台核心抽象
│   │   ├── contracts/       # 跨模組契約定義（原 ifaces）
│   │   │   ├── alerting/
│   │   │   ├── monitoring/
│   │   │   ├── notification/
│   │   │   ├── storage/
│   │   │   └── lifecycle/   # 插件生命週期
│   │   ├── registry/        # 註冊機制抽象
│   │   │   ├── interface.go
│   │   │   ├── discovery.go
│   │   │   └── composer.go  # 組合邏輯
│   │   └── composition/     # 組合模式定義
│   │       ├── app.go       # 應用組合介面
│   │       ├── module.go    # 模組組合介面
│   │       └── plugin.go    # 插件組合介面
│   │
│   ├── domain/              # 🎯 領域模型（業務實體）
│   │   ├── alert/
│   │   ├── metric/
│   │   ├── rule/
│   │   └── notification/
│   │
│   ├── config/              # ⚙️ 統一配置管理
│   │   ├── types/           # 配置結構定義（原 configtypes）
│   │   ├── schema/          # 配置模式驗證
│   │   ├── composition/     # 組合配置
│   │   └── loader/          # 配置載入器
│   │
│   └── shared/              # 🔧 共用工具與常數
│       ├── errors/
│       ├── utils/
│       ├── constants/
│       └── types/           # 基礎類型定義
│
├── internal/                # 🏭 實作層（Clean Architecture 分層）
│   ├── platform/            # 🏗️ 平台實作層
│   │   ├── registry/        # 註冊機制實作
│   │   │   ├── manager.go   # 註冊管理器
│   │   │   ├── discovery.go # 自動發現
│   │   │   └── resolver.go  # 依賴解析
│   │   ├── composition/     # 組合引擎實作
│   │   │   ├── builder.go   # 組合建構器
│   │   │   ├── lifecycle.go # 生命週期管理
│   │   │   └── injector.go  # 依賴注入
│   │   └── runtime/         # 運行時管理
│   │       ├── bootstrap.go
│   │       ├── shutdown.go
│   │       └── health.go
│   │
│   ├── adapters/            # 🔌 適配器層（外部系統介接）
│   │   ├── datasources/     # 資料源適配器
│   │   │   ├── prometheus/
│   │   │   ├── influxdb/
│   │   │   └── mysql/
│   │   ├── notifications/   # 通知適配器
│   │   │   ├── email/
│   │   │   ├── slack/
│   │   │   └── webhook/
│   │   └── storage/         # 儲存適配器
│   │       ├── redis/
│   │       ├── memory/
│   │       └── filesystem/
│   │
│   ├── services/            # 💼 服務層（業務邏輯）
│   │   ├── alerting/        # 告警服務
│   │   ├── monitoring/      # 監控服務
│   │   ├── notification/    # 通知服務
│   │   ├── reporting/       # 報告服務
│   │   └── composition/     # 組合服務
│   │
│   ├── repositories/        # 🗄️ 資料存取層
│   │   ├── alert/
│   │   ├── rule/
│   │   ├── metric/
│   │   ├── user/
│   │   └── plugin/          # 插件元資料儲存
│   │
│   ├── infrastructure/      # 🏗️ 基礎設施層
│   │   ├── eventbus/        # 事件匯流排
│   │   ├── cache/           # 快取系統
│   │   ├── metrics/         # 指標收集
│   │   ├── logging/         # 日誌系統
│   │   └── tracing/         # 追蹤系統
│   │
│   └── ports/               # 🚪 輸入輸出埠
│       ├── http/            # HTTP API 埠
│       │   ├── handlers/    # API 處理器
│       │   ├── middleware/  # HTTP 中介層
│       │   └── routes/      # 路由定義
│       ├── grpc/            # gRPC 埠
│       ├── cli/             # CLI 埠
│       └── web/             # Web UI 埠
│
├── plugins/                 # 🔌 插件生態系統
│   ├── core/                # 核心插件（內建）
│   │   ├── auth/            # 認證插件
│   │   │   ├── basic/
│   │   │   ├── jwt/
│   │   │   └── ldap/
│   │   ├── datasources/     # 資料源插件
│   │   │   ├── prometheus/
│   │   │   ├── influxdb/
│   │   │   └── elasticsearch/
│   │   ├── notifiers/       # 通知插件
│   │   │   ├── email/
│   │   │   ├── slack/
│   │   │   └── teams/
│   │   └── middleware/      # 中介層插件
│   │       ├── cors/
│   │       ├── ratelimit/
│   │       └── logging/
│   │
│   ├── community/           # 社群插件
│   │   ├── exporters/       # 資料匯出器
│   │   ├── importers/       # 資料匯入器
│   │   └── integrations/    # 第三方整合
│   │
│   └── custom/              # 自訂插件
│       └── example/         # 插件範例
│
├── compositions/            # 🎼 組合定義（可選）
│   ├── monitoring-stack/    # 監控堆疊組合
│   ├── alerting-platform/   # 告警平台組合
│   └── full-platform/       # 完整平台組合
│
└── examples/                # 📝 範例與教學
    ├── quick-start/
    ├── custom-plugin/
    └── composition-guide/
```

## 組合式架構的關鍵概念

### 1. 平台抽象層 (pkg/platform/)
```go
// 組合介面定義
type ApplicationComposer interface {
    Compose(config *CompositionConfig) (*Application, error)
    AddModule(module Module) error
    AddPlugin(plugin Plugin) error
}

// 模組介面定義
type Module interface {
    ID() string
    Dependencies() []string
    Initialize(ctx Context) error
    Shutdown(ctx Context) error
}

// 插件介面定義  
type Plugin interface {
    Metadata() PluginMetadata
    Register(registry Registry) error
    Configure(config PluginConfig) error
}
```

### 2. 組合配置 (pkg/config/composition/)
```yaml
# 應用組合範例
composition:
  name: "monitoring-platform"
  modules:
    - alerting
    - monitoring  
    - notification
  plugins:
    - core/auth/jwt
    - core/datasources/prometheus
    - core/notifiers/slack
  configuration:
    alerting:
      evaluation_interval: "30s"
    monitoring:
      scrape_interval: "15s"
```

### 3. 註冊與發現 (internal/platform/registry/)
```go
// 自動發現插件
type PluginDiscovery interface {
    Discover(paths []string) ([]Plugin, error)
    Watch(callback DiscoveryCallback) error
}

// 依賴解析
type DependencyResolver interface {
    Resolve(dependencies []Dependency) ([]Module, error)
    ValidateComposition(composition *Composition) error
}
```

## 插件系統重新設計

### 插件分類與命名
```bash
plugins/
├── core/                    # 內建核心插件
│   ├── auth/
│   │   ├── basic/
│   │   │   ├── plugin.go
│   │   │   ├── strategy.go
│   │   │   └── config.go
│   │   └── jwt/
│   ├── datasources/
│   └── notifiers/
├── community/               # 社群貢獻插件
└── custom/                  # 使用者自訂插件
```

### 插件元資料標準
```go
type PluginMetadata struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Version     string            `json:"version"`
    Type        PluginType        `json:"type"`
    Category    string            `json:"category"`
    Dependencies []string         `json:"dependencies"`
    Capabilities []Capability     `json:"capabilities"`
    Schema      *ConfigSchema     `json:"schema,omitempty"`
}
```

## 組合式開發流程

### 1. 模組開發
```go
// 1. 定義契約 (pkg/platform/contracts/)
type AlertingService interface {
    EvaluateRules(ctx context.Context) error
    ProcessAlert(alert *domain.Alert) error
}

// 2. 實作服務 (internal/services/alerting/)
type alertingService struct {
    evaluator contracts.AlertEvaluator
    notifier  contracts.Notifier
}

// 3. 註冊模組
func init() {
    platform.RegisterModule(&AlertingModule{})
}
```

### 2. 插件開發
```go
// 1. 實作插件介面
type PrometheusPlugin struct {
    config *PrometheusConfig
}

func (p *PrometheusPlugin) Metadata() PluginMetadata {
    return PluginMetadata{
        ID:       "prometheus-datasource",
        Type:     "datasource",
        Category: "monitoring",
    }
}

// 2. 註冊能力
func (p *PrometheusPlugin) Register(registry Registry) error {
    return registry.RegisterDatasource("prometheus", p.newDatasource)
}
```

### 3. 應用組合
```go
// apps/server/main.go
func main() {
    composer := platform.NewComposer()
    
    // 載入組合配置
    config := config.LoadComposition("monitoring-platform")
    
    // 建構應用
    app, err := composer.Compose(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // 啟動應用
    app.Run()
}
```

## 實作階段規劃

### 第一階段：平台基礎 (Week 1-3)
1. 建立 `pkg/platform/` 基礎抽象
2. 實作基本的註冊與組合機制
3. 定義插件介面和元資料標準
4. 重構現有的註冊邏輯

### 第二階段：插件系統 (Week 4-6)  
1. 重構現有插件到新的目錄結構
2. 實作插件自動發現和載入
3. 建立插件依賴解析機制
4. 實作組合引擎

### 第三階段：組合式配置 (Week 7-8)
1. 設計組合配置格式
2. 實作配置驗證和載入
3. 建立應用組合範例
4. 完善文檔和教學

### 第四階段：生態系統 (Week 9-12)
1. 遷移現有模組到新架構
2. 建立插件開發工具
3. 完善測試和CI/CD
4. 建立社群貢獻流程

## 向後相容策略

### 漸進式遷移
```go
// 使用 alias 保持相容性
package ifaces

import "detectviz/pkg/platform/contracts/alerting"

// Deprecated: Use pkg/platform/contracts/alerting.Evaluator
type AlertEvaluator = alerting.Evaluator
```

### 組合模式相容
```go
// 舊的直接初始化方式仍然支援
func NewLegacyApplication() *Application {
    return &Application{
        // 傳統初始化邏輯
    }
}

// 新的組合式初始化
func NewComposedApplication(config CompositionConfig) *Application {
    composer := platform.NewComposer()
    return composer.Compose(config)
}
```

---

這個重構將 DetectViz 轉變為真正的組合式平台，支援靈活的模組組合和插件生態系統，同時保持 Clean Architecture 的設計原則。 

---

