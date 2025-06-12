# DetectViz 組合式架構遷移計劃

## 當前架構 vs 建議架構對照表

### 📊 整體對照分析

| 層級 | 當前結構 | 建議結構 | 狀態 | 行動 |
|------|----------|----------|------|------|
| **應用層** | `apps/` ✅ | `apps/` | 🟢 保持 | 無需變更 |
| **公共契約** | `pkg/ifaces/` | `pkg/platform/contracts/` | 🟡 重構 | 遷移+重組 |
| **配置管理** | `pkg/configtypes/` | `pkg/config/types/` | 🟡 重構 | 遷移+統一 |
| **註冊機制** | `pkg/registry/` + `internal/registry/` | `pkg/platform/registry/` + `internal/platform/registry/` | 🔴 重構 | 統一+分層 |
| **插件系統** | `internal/plugins/` + `plugins/` | `plugins/core/` + `internal/platform/` | 🟡 重構 | 重組+擴展 |
| **適配器層** | `internal/adapters/` ✅ | `internal/adapters/` | 🟢 保持 | 輕微調整 |
| **服務層** | 分散在各模組 | `internal/services/` | 🔴 新增 | 提取+統一 |
| **基礎設施** | 分散實作 | `internal/infrastructure/` | 🔴 新增 | 重組+統一 |

## 詳細遷移對照

### 🏗️ **平台抽象層 (新增)**

#### 需要新增的目錄結構：
```bash
pkg/platform/                    # 🆕 全新建立
├── contracts/                   # 🔄 從 pkg/ifaces/ 遷移
│   ├── alerting/               # 🔄 重組 pkg/ifaces/alert/
│   ├── monitoring/             # 🔄 重組 pkg/ifaces/metrics/
│   ├── notification/           # 🔄 重組 pkg/ifaces/notifier/
│   ├── storage/                # 🔄 重組 pkg/ifaces/cachestore/
│   ├── scheduling/             # 🔄 重組 pkg/ifaces/scheduler/
│   └── lifecycle/              # 🆕 新增插件生命週期契約
├── registry/                   # 🆕 統一註冊抽象
│   ├── interface.go            # 🆕 註冊介面定義
│   ├── discovery.go            # 🆕 自動發現機制
│   └── composer.go             # 🆕 組合邏輯
└── composition/                # 🆕 組合模式定義
    ├── app.go                  # 🆕 應用組合介面
    ├── module.go               # 🆕 模組組合介面
    └── plugin.go               # 🆕 插件組合介面
```

#### 遷移動作：
1. **創建新目錄結構**
2. **遷移現有介面**：
   - `pkg/ifaces/alert/` → `pkg/platform/contracts/alerting/`
   - `pkg/ifaces/metrics/` → `pkg/platform/contracts/monitoring/`
   - `pkg/ifaces/notifier/` → `pkg/platform/contracts/notification/`
   - `pkg/ifaces/scheduler/` → `pkg/platform/contracts/scheduling/`
   - `pkg/ifaces/cachestore/` → `pkg/platform/contracts/storage/`
3. **新增組合相關介面**

### 📋 **領域模型層 (新增)**

#### 需要新增：
```bash
pkg/domain/                     # 🆕 全新建立
├── alert/                      # 🆕 告警領域模型
│   ├── alert.go               # 🆕 Alert 實體
│   ├── rule.go                # 🆕 Rule 值物件
│   └── condition.go           # 🆕 Condition 值物件
├── metric/                     # 🆕 指標領域模型
│   ├── metric.go              # 🆕 Metric 實體
│   └── series.go              # 🆕 TimeSeries 值物件
├── rule/                       # 🆕 規則領域模型
└── notification/               # 🆕 通知領域模型
    ├── channel.go             # 🆕 Channel 實體
    └── template.go            # 🆕 Template 值物件
```

#### 遷移動作：
1. **從現有代碼提取領域實體**
2. **建立值物件和聚合根**
3. **定義領域服務介面**

### ⚙️ **配置管理統一化**

#### 當前狀態：
```bash
pkg/configtypes/                # 🔄 需要遷移
├── cache_config.go
└── notifier_config.go

pkg/config/                     # 🔄 需要擴展
├── default.go
└── README.md
```

#### 建議結構：
```bash
pkg/config/                     # 🔄 擴展現有
├── types/                      # 🔄 從 pkg/configtypes/ 遷移
│   ├── server.go              # 🆕 新增
│   ├── database.go            # 🆕 新增
│   ├── cache.go               # 🔄 從 cache_config.go 遷移
│   ├── alert.go               # 🆕 新增
│   └── notification.go        # 🔄 從 notifier_config.go 遷移
├── schema/                     # 🆕 新增配置模式驗證
├── composition/                # 🆕 新增組合配置
└── loader/                     # 🔄 擴展現有
```

#### 遷移動作：
1. **移動配置類型**：`pkg/configtypes/` → `pkg/config/types/`
2. **新增配置驗證和組合邏輯**
3. **更新所有 import 路徑**

### 🏭 **內部實作層重構**

#### 需要新增的平台實作：
```bash
internal/platform/              # 🆕 全新建立
├── registry/                   # 🔄 整合現有註冊邏輯
│   ├── manager.go             # 🔄 從 internal/plugins/manager/ 整合
│   ├── discovery.go           # 🆕 新增自動發現
│   └── resolver.go            # 🆕 新增依賴解析
├── composition/                # 🆕 新增組合引擎
│   ├── builder.go             # 🆕 組合建構器
│   ├── lifecycle.go           # 🔄 從 internal/plugins/manager/lifecycle.go 整合
│   └── injector.go            # 🆕 依賴注入
└── runtime/                    # 🆕 新增運行時管理
    ├── bootstrap.go           # 🔄 從 internal/bootstrap/ 整合
    ├── shutdown.go            # 🆕 新增
    └── health.go              # 🆕 新增
```

#### 需要新增的服務層：
```bash
internal/services/              # 🆕 全新建立
├── alerting/                   # 🆕 提取告警業務邏輯
├── monitoring/                 # 🆕 提取監控業務邏輯
├── notification/               # 🆕 提取通知業務邏輯
├── reporting/                  # 🆕 新增報告服務
└── composition/                # 🆕 組合服務
```

#### 需要新增的資料存取層：
```bash
internal/repositories/          # 🆕 全新建立
├── alert/                      # 🆕 告警資料存取
├── rule/                       # 🆕 規則資料存取
├── metric/                     # 🆕 指標資料存取
├── user/                       # 🆕 使用者資料存取
└── plugin/                     # 🆕 插件元資料儲存
```

#### 需要新增的基礎設施層：
```bash
internal/infrastructure/        # 🆕 全新建立
├── eventbus/                   # 🔄 從 internal/adapters/eventbus/ 遷移
├── cache/                      # 🔄 從 internal/adapters/cachestore/ 遷移
├── metrics/                    # 🔄 從 internal/adapters/metrics/ 遷移
├── logging/                    # 🔄 從 internal/adapters/logger/ 遷移
└── tracing/                    # 🆕 新增追蹤系統
```

#### 需要新增的輸入輸出埠：
```bash
internal/ports/                 # 🆕 全新建立
├── http/                       # 🆕 HTTP API 埠
│   ├── handlers/              # 🆕 API 處理器
│   ├── middleware/            # 🆕 HTTP 中介層
│   └── routes/                # 🆕 路由定義
├── grpc/                       # 🆕 gRPC 埠
├── cli/                        # 🆕 CLI 埠
└── web/                        # 🆕 Web UI 埠
```

### 🔌 **插件系統重組**

#### 當前狀態：
```bash
internal/plugins/               # 🔄 需要重構
├── manager/                    # 🔄 遷移到 internal/platform/
├── eventbus/                   # 🔄 保持
└── plugin.go                   # 🔄 重構

plugins/                        # 🔄 需要重組
├── auth/                       # 🔄 遷移到 plugins/core/auth/
├── datasources/                # 🔄 遷移到 plugins/core/datasources/
├── exporter/                   # 🔄 遷移到 plugins/community/exporters/
├── tools/                      # 🔄 遷移到 plugins/community/tools/
└── visuals/                    # 🔄 遷移到 plugins/community/visuals/
```

#### 建議結構：
```bash
plugins/                        # 🔄 重組現有
├── core/                       # 🆕 內建核心插件
│   ├── auth/                   # 🔄 從 plugins/auth/ 遷移
│   │   ├── basic/             # 🆕 基礎認證
│   │   ├── jwt/               # 🆕 JWT 認證
│   │   └── ldap/              # 🆕 LDAP 認證
│   ├── datasources/            # 🔄 從 plugins/datasources/ 遷移
│   │   ├── prometheus/        # 🆕 Prometheus 插件
│   │   ├── influxdb/          # 🆕 InfluxDB 插件
│   │   └── elasticsearch/     # 🆕 Elasticsearch 插件
│   ├── notifiers/              # 🆕 通知插件
│   │   ├── email/             # 🆕 郵件通知
│   │   ├── slack/             # 🆕 Slack 通知
│   │   └── teams/             # 🆕 Teams 通知
│   └── middleware/             # 🆕 中介層插件
│       ├── cors/              # 🆕 CORS 中介層
│       ├── ratelimit/         # 🆕 限速中介層
│       └── logging/           # 🆕 日誌中介層
├── community/                  # 🆕 社群插件
│   ├── exporters/              # 🔄 從 plugins/exporter/ 遷移
│   ├── importers/              # 🆕 資料匯入器
│   ├── integrations/           # 🆕 第三方整合
│   ├── tools/                  # 🔄 從 plugins/tools/ 遷移
│   └── visuals/                # 🔄 從 plugins/visuals/ 遷移
└── custom/                     # 🆕 自訂插件
    └── example/                # 🆕 插件範例
```

### 🎼 **組合定義 (新增)**

#### 需要新增：
```bash
compositions/                   # 🆕 全新建立
├── monitoring-stack/           # 🆕 監控堆疊組合
│   ├── composition.yaml       # 🆕 組合配置
│   └── README.md              # 🆕 說明文檔
├── alerting-platform/          # 🆕 告警平台組合
│   ├── composition.yaml       # 🆕 組合配置
│   └── README.md              # 🆕 說明文檔
└── full-platform/              # 🆕 完整平台組合
    ├── composition.yaml        # 🆕 組合配置
    └── README.md               # 🆕 說明文檔
```

### 📝 **範例與教學 (新增)**

#### 需要新增：
```bash
examples/                       # 🆕 全新建立
├── quick-start/                # 🆕 快速開始範例
├── custom-plugin/              # 🆕 自訂插件範例
└── composition-guide/          # 🆕 組合指南範例
```

## 🚀 實作階段規劃

### **第一階段：平台基礎建設 (Week 1-3)**

#### Week 1: 平台抽象層
- [ ] 創建 `pkg/platform/` 目錄結構
- [ ] 遷移 `pkg/ifaces/` → `pkg/platform/contracts/`
- [ ] 定義組合相關介面
- [ ] 更新所有 import 路徑

#### Week 2: 配置管理統一
- [ ] 遷移 `pkg/configtypes/` → `pkg/config/types/`
- [ ] 新增配置驗證和組合邏輯
- [ ] 建立組合配置格式
- [ ] 更新配置載入邏輯

#### Week 3: 註冊機制統一
- [ ] 整合 `pkg/registry/` 和 `internal/registry/`
- [ ] 建立統一註冊介面
- [ ] 實作自動發現機制
- [ ] 實作依賴解析邏輯

### **第二階段：插件系統重構 (Week 4-6)**

#### Week 4: 插件目錄重組
- [ ] 重組 `plugins/` 目錄結構
- [ ] 遷移現有插件到新分類
- [ ] 定義插件元資料標準
- [ ] 更新插件載入邏輯

#### Week 5: 插件管理增強
- [ ] 整合 `internal/plugins/manager/` 到 `internal/platform/`
- [ ] 實作插件自動發現
- [ ] 實作插件依賴解析
- [ ] 實作插件生命週期管理

#### Week 6: 組合引擎實作
- [ ] 實作組合建構器
- [ ] 實作依賴注入機制
- [ ] 實作運行時管理
- [ ] 建立組合配置範例

### **第三階段：架構分層完善 (Week 7-9)**

#### Week 7: 服務層建立
- [ ] 創建 `internal/services/` 目錄
- [ ] 提取業務邏輯到服務層
- [ ] 建立服務間通信機制
- [ ] 實作服務註冊和發現

#### Week 8: 資料存取層
- [ ] 創建 `internal/repositories/` 目錄
- [ ] 實作資料存取介面
- [ ] 遷移現有資料存取邏輯
- [ ] 建立資料存取抽象

#### Week 9: 基礎設施層
- [ ] 創建 `internal/infrastructure/` 目錄
- [ ] 重組基礎設施元件
- [ ] 實作統一的基礎設施介面
- [ ] 建立基礎設施配置

### **第四階段：輸入輸出埠與整合 (Week 10-12)**

#### Week 10: 輸入輸出埠
- [ ] 創建 `internal/ports/` 目錄
- [ ] 實作 HTTP API 埠
- [ ] 實作 CLI 埠
- [ ] 實作 Web UI 埠

#### Week 11: 領域模型
- [ ] 創建 `pkg/domain/` 目錄
- [ ] 提取領域實體和值物件
- [ ] 建立領域服務
- [ ] 實作領域事件

#### Week 12: 整合測試與文檔
- [ ] 建立組合範例
- [ ] 撰寫遷移指南
- [ ] 完善測試覆蓋
- [ ] 更新所有文檔

## 🔄 向後相容策略

### 漸進式遷移
```go
// 在舊包中使用 alias 保持相容性
package ifaces

import "detectviz/pkg/platform/contracts/alerting"

// Deprecated: Use pkg/platform/contracts/alerting.Evaluator
type AlertEvaluator = alerting.Evaluator
```

### 雙模式支援
```go
// 同時支援舊的直接初始化和新的組合式初始化
func NewApplication(config Config) *Application {
    if config.UseComposition {
        return NewComposedApplication(config.Composition)
    }
    return NewLegacyApplication(config)
}
```

## 📊 風險評估與緩解

### 高風險項目
1. **大量 import 路徑變更** - 使用自動化腳本批量處理
2. **插件相容性問題** - 保持舊介面 alias，逐步遷移
3. **配置格式變更** - 支援舊格式自動轉換

### 緩解策略
1. **分階段實施** - 每個階段獨立驗證
2. **自動化測試** - 確保每次變更不破壞現有功能
3. **回滾計劃** - 每個階段都有明確的回滾步驟

## ✅ 驗收標準

### 功能驗收
- [ ] 所有現有功能正常運作
- [ ] 新的組合式架構可以正常組裝應用
- [ ] 插件系統支援動態載入和卸載
- [ ] 配置系統支援組合式配置

### 效能驗收
- [ ] 啟動時間不超過原來的 120%
- [ ] 記憶體使用不超過原來的 110%
- [ ] API 回應時間保持在相同水準

### 可維護性驗收
- [ ] 新增插件的開發時間縮短 50%
- [ ] 新增應用組合的時間縮短 70%
- [ ] 程式碼覆蓋率維持在 80% 以上

---

這個遷移計劃將 DetectViz 從傳統架構轉變為真正的組合式平台，支援靈活的模組組合和插件生態系統。 

---

# GPT

## 📁 目錄結構對照分析表

| Cursor 建議目錄 | DetectViz 現況 | 建議調整 | 備註說明 |
| --- | --- | --- | --- |
| `apps/server/`  
`apps/cli/`  
`apps/agent/` | ✅ 已有 `/apps/server/`, `/apps/cli/` | ✅ 採納 | 可對應應用組裝層 |
| `pkg/platform/contracts/` | `pkg/ifaces/` | ✅ **採納並改名** | 改為 `pkg/platform/contracts/`，語意更明確、符合平台分層 |
| `pkg/platform/registry/` | `pkg/registry/` | ✅ 採納 | 定義 interface 註冊行為，對齊 plugin-first |
| `pkg/platform/composition/` | ❌ 無 | ✅ 新增 | 定義 app/module/plugin 組合邏輯與抽象 |
| `pkg/domain/` | ✅ 有 `pkg/domain/{mod}` | ⚠️ 微調建議 | 長期可分 context（e.g. domain/alerting） |
| `pkg/config/types/` | `pkg/configtypes/` | ✅ 改名統一 | 合併為 `pkg/config/types/` 統一層級與風格 |
| `pkg/config/loader/` | ✅ 已有 | ✅ 不動 |  |
| `pkg/shared/utils/`, `pkg/shared/errors/` | `pkg/utils/`, `pkg/errors/` | ⚠️ **微調建議：改為 `pkg/common/`** | `shared/` 命名偏少見，`common/` 更廣泛接受 |
| `internal/services/{mod}/` | ✅ 有 | ✅ 不動 | 正確放置 Service 實作邏輯 |
| `internal/adapters/{mod}/` | ✅ 有 | ✅ 不動（但命名統一 `_adapter.go`） | 建議保留 adapter 後綴命名一致性 |
| `internal/repositories/{mod}/` | 部分存在 | ✅ 採納 | 取代過去 `internal/store/` 混用邏輯 |
| `internal/infrastructure/cache/` 等 | ✅ 有 `internal/system`、`infra` | ✅ 採納並重整為 `infrastructure/` | `metrics`, `logging`, `eventbus` 可歸入 |
| `internal/ports/http/handlers/`  
`.../middleware/`  
`.../routes/` | ✅ 有 `/internal/web`, `/internal/api` | ✅ 採納 → 合併為 `/ports/http` 下分層 | 對齊 DDD port 層次 |
| `plugins/core/`  
`plugins/community/`  
`plugins/custom/` | ✅ 有 | ✅ 不動 | 非常清晰的分層，已貼近 Grafana/Thanos 模式 |
| `compositions/` | ❌ 無 | ✅ 新增 | 支援組合式堆疊定義，對應 AppComposer |
| `examples/` | ✅ 有 `/docs/examples/` | ✅ 可保留 | 教學與腳本存放處 |

* * *

## 📌 建議採納情況總結

| 分類 | 採納 | 微調建議 | 無需調整 |
| --- | --- | --- | --- |
| ✅ 採納 | apps/, platform/, config/, plugins/, ports/, services/, compositions/ | domain/, adapters/, shared/ → common/, configtypes/ | loader/, registry/, examples/ |

* * *

## 🧭 補充建議：detectviz 特有補充模組

| detectviz 特有模組 | Cursor 沒有 | 建議 |
| --- | --- | --- |
| `internal/system/` | ❌ | ✅ 可併入 `internal/infrastructure/` |
| `internal/modules/` | ❌ | ✅ 應遷移為 `pkg/platform/composition/` + plugin 註冊 |
| `pkg/config/composition/` | ❌ | ✅ 需建立，支援組合配置與驗證 |
| `internal/plugins/manager` | ❌ | ✅ 移至 `internal/platform/registry/` 為 runtime 註冊入口 |

* * *

## ✅ 總結：你是否應接納這份 Cursor 架構？

**非常建議採納**。你目前的 detectviz 架構與 Cursor 所推薦的組合式架構幾乎完全對齊。只需微調命名與整併部分目錄，便能：

* 提高一致性與可維護性
    
* 加速 plugin 與組合式模組開發流程
    
* 支援 future-proof 的平台化開發模式
    
