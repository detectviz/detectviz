# Plugin Interface 基礎介面

> **檔案位置**: `pkg/platform/contracts/plugin.go`

## 概述

`Plugin` 是 DetectViz 平台中所有插件必須實作的基礎介面。它定義了插件的核心行為，包括基本資訊、初始化和關閉流程。

## 介面定義

```go
type Plugin interface {
    Name() string
    Version() string  
    Description() string
    Init(config any) error
    Shutdown() error
}
```

## 方法說明

### Name() string
- **用途**: 回傳插件的唯一識別名稱
- **回傳值**: 插件名稱字串，用於註冊和查找
- **範例**: `"jwt-authenticator"`, `"prometheus-importer"`

### Version() string  
- **用途**: 回傳插件版本資訊
- **回傳值**: 語意版本字串
- **範例**: `"1.0.0"`, `"2.1.3-beta"`

### Description() string
- **用途**: 回傳插件功能描述
- **回傳值**: 人類可讀的描述文字
- **範例**: `"JWT-based authentication plugin"`

### Init(config any) error
- **用途**: 初始化插件，解析配置並準備運行
- **參數**: `config` - 插件配置，可為 map[string]any 或結構體
- **回傳值**: 初始化失敗時回傳錯誤
- **注意**: 應支援 mapstructure 解碼配置

### Shutdown() error
- **用途**: 關閉插件，清理資源
- **回傳值**: 關閉失敗時回傳錯誤
- **注意**: 應確保安全關閉所有相關資源

## 相關介面

- **LifecycleAware**: 提供生命週期管理
- **HealthChecker**: 提供健康檢查能力
- **ConfigurablePlugin**: 提供配置驗證功能

## 實作範例

```go
type MyPlugin struct {
    name        string
    version     string
    description string
    config      *MyConfig
    initialized bool
}

func (p *MyPlugin) Name() string {
    return p.name
}

func (p *MyPlugin) Version() string {
    return p.version
}

func (p *MyPlugin) Description() string {
    return p.description
}

func (p *MyPlugin) Init(config any) error {
    if p.initialized {
        return nil
    }
    
    // 使用 mapstructure 解析配置
    if err := parsePluginConfig(config, p.config); err != nil {
        return fmt.Errorf("failed to parse config: %w", err)
    }
    
    p.initialized = true
    return nil
}

func (p *MyPlugin) Shutdown() error {
    p.initialized = false
    return nil
}
```

## 技術棧要求

1. **配置解析**: 使用 `mapstructure` 標籤解碼配置
2. **日誌記錄**: 使用 `otelzap` 或 `logrus`，避免 `fmt.Print`
3. **上下文傳遞**: 支援 `context.Context` 參數

## 相關文件

- [LifecycleAware Interface](./lifecycle.md)
- [HealthChecker Interface](./lifecycle.md#healthchecker)
- [Configuration Schema](../config/schema.md)
- [Plugin Development Guide](../develop-guide.md) 