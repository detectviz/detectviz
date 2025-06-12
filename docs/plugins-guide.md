

# Plugins Guide

本文件說明 detectviz 插件系統的設計原則、註冊方式與擴充範例，幫助開發者撰寫與整合自定義插件模組。

---

## 插件設計目標

- 支援外部模組獨立註冊、載入與初始化
- 可針對以下類型進行 plugin 擴充：
  - Auth 驗證策略（如 keycloak, ldap）
  - Middleware 中介層（如限速器、觀察者）
  - Store 後端儲存擴充（如 influxdb、logfile）
  - CLI 指令（如 rule apply, plugin enable）
  - API Hook（擴充額外 handler）

---

## 插件管理目錄結構

```
internal/plugins/
├── auth/               # 擴充 AuthStrategy 註冊點
├── middleware/         # 擴充 Middleware plugin 註冊點
├── store/              # 擴充 Store plugin 註冊點（支援多 source）
├── apihooks/           # 註冊自定 API handler
├── eventbus/           # 擴充事件監控模組（如 alertlog）
├── manager/            # 核心 lifecycle 管理器
│   ├── loader.go       # 掃描 plugin 資料夾與初始化
│   ├── registry.go     # 全域註冊器
│   └── process.go      # 掛載插件到運行時
└── plugin.go           # Plugin interface 與註冊入口
```

---

## Plugin 註冊方式說明

```go
// 定義 plugin
type Plugin interface {
    ID() string
    Init() error
}

// 透過 manager 註冊
plugins.Register("auth.keycloak", &KeycloakStrategy{})
plugins.Register("store.redis", &RedisStore{})
plugins.Register("middleware.logger", NewLoggerMiddleware())
```

---

## 常見插件類型實作範例

### Auth Plugin

```go
type KeycloakStrategy struct{}

func (s *KeycloakStrategy) ID() string { return "keycloak" }
func (s *KeycloakStrategy) Init() error {
    auth.RegisterStrategy("keycloak", s)
    return nil
}
```

### Store Plugin

```go
type RedisStore struct{}

func (s *RedisStore) ID() string { return "redis" }
func (s *RedisStore) Init() error {
    store.Register("redis", s)
    return nil
}
```

---

## 自定義 Plugin 實作建議

- 插件應實作 `Plugin` interface 並在 `Init()` 方法中自行註冊至對應模組
- 插件不得擁有全域副作用（如環境變數或未控制的 logger）
- 插件初始化順序由 `manager/loader.go` 控制

---

## 插件掃描與掛載流程

```mermaid
graph LR
    A[plugins.Register] --> B[manager.registry]
    B --> C[loader.LoadAll]
    C --> D[plugin.Init()]
    D --> E[plugin-specific Register]
```

---

## 插件載入與設定建議

- 可透過 `pkg/config` 提供 plugin 個別設定
- 設定命名建議：`plugin.{type}.{id}.{field}`

```yaml
plugin.auth.keycloak.url: https://auth.example.com
plugin.store.redis.addr: localhost:6379
```

---

## 測試建議

- plugin 應提供 `*_test.go` 單元測試
- 可搭配 `internal/test/fakes` 建立模擬註冊流程
- plugin 實作應可重複 Init 並具備可觀測性

---