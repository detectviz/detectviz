# WebUIPlugin 介面規格

WebUIPlugin 介面定義了 DetectViz Web UI 擴展插件的標準規範，允許插件註冊自訂路由、導覽節點和 UI 組件。

## 介面定義

```go
type WebUIPlugin interface {
    Plugin
    RegisterRoutes(router WebRouter) error
    RegisterNavNodes(navtree NavTreeBuilder) error
    RegisterComponents(registry ComponentRegistry) error
}
```

## 方法說明

### RegisterRoutes(router WebRouter) error

註冊 Web 路由，允許插件提供自訂的 HTTP 端點。

**用途：**
- 註冊頁面路由（GET 請求）
- 註冊 API 端點（GET、POST、PUT、DELETE）
- 建立路由群組

**範例：**
```go
func (p *MyPlugin) RegisterRoutes(router contracts.WebRouter) error {
    // 註冊主頁面
    router.GET("/my-plugin", p.handleMainPage)
    
    // 註冊 API 端點
    apiGroup := router.Group("/api/my-plugin")
    apiGroup.GET("/data", p.handleGetData)
    apiGroup.POST("/config", p.handlePostConfig)
    
    return nil
}
```

### RegisterNavNodes(navtree NavTreeBuilder) error

註冊導覽節點，將插件頁面添加到系統導覽選單中。

**用途：**
- 添加主選單項目
- 建立子選單結構
- 設定權限控制
- 添加徽章和圖示

**範例：**
```go
func (p *MyPlugin) RegisterNavNodes(navtree contracts.NavTreeBuilder) error {
    mainNode := contracts.NavNode{
        ID:         "my-plugin",
        Title:      "我的插件",
        Icon:       "fas fa-cog",
        URL:        "/my-plugin",
        Permission: "my-plugin.view",
        Order:      50,
        Visible:    true,
        Enabled:    true,
        Badge: &contracts.NavBadge{
            Text:  "New",
            Color: "info",
            Style: "pill",
        },
    }
    
    navtree.AddNode("my-plugin", mainNode)
    
    // 添加子節點
    subNode := contracts.NavNode{
        ID:         "my-plugin-config",
        Title:      "設定",
        Icon:       "fas fa-wrench",
        URL:        "/my-plugin/config",
        Permission: "my-plugin.config",
        Order:      1,
        Visible:    true,
        Enabled:    true,
    }
    
    navtree.AddChildNode("my-plugin", "my-plugin-config", subNode)
    
    return nil
}
```

### RegisterComponents(registry ComponentRegistry) error

註冊 UI 組件，包括 partials、widgets、themes 和 assets。

**用途：**
- 註冊可重複使用的模板片段（partials）
- 註冊互動式小工具（widgets）
- 註冊主題和樣式
- 註冊靜態資源

**範例：**
```go
func (p *MyPlugin) RegisterComponents(registry contracts.ComponentRegistry) error {
    // 註冊模板片段
    registry.RegisterPartial("my-card", "/templates/my-card.html")
    
    // 註冊 widget
    registry.RegisterWidget("status-indicator", p.statusWidget)
    
    // 註冊主題
    theme := contracts.Theme{
        Name:        "my-theme",
        DisplayName: "我的主題",
        Version:     "1.0.0",
        CSSFiles:    []string{"/assets/my-theme.css"},
        JSFiles:     []string{"/assets/my-theme.js"},
    }
    registry.RegisterTheme("my-theme", theme)
    
    // 註冊靜態資源
    asset := contracts.Asset{
        Name:    "my-plugin-css",
        Type:    "css",
        Path:    "/assets/my-plugin.css",
        Version: "1.0.0",
    }
    registry.RegisterAsset("my-plugin-css", asset)
    
    return nil
}
```

## Platform 自動載入與註冊流程

DetectViz 平台會自動掃描並載入實作 WebUIPlugin 介面的插件，以下是完整的註冊流程：

### 1. 插件掃描階段

**掃描流程**：
```go
// 平台啟動時的插件掃描流程
func (p *Platform) scanWebUIPlugins() error {
    // 1. 掃描所有已註冊的插件
    plugins := p.registry.GetAllPlugins()
    
    // 2. 檢查是否實作 WebUIPlugin 介面
    for _, plugin := range plugins {
        if webPlugin, ok := plugin.(contracts.WebUIPlugin); ok {
            // 3. 添加到 WebUI 插件列表
            p.webUIPlugins = append(p.webUIPlugins, webPlugin)
        }
    }
    
    return nil
}
```

**掃描路徑結構**：
```
plugins/
├── core/                    # 核心插件（自動掃描）
│   ├── auth/
│   └── middleware/
├── community/               # 社群插件（依 composition 載入）
│   └── integrations/
├── web/
│   ├── pages/          # 頁面插件
│   ├── components/     # 組件插件
│   ├── themes/         # 主題插件
│   └── widgets/        # 小工具插件
└── tools/              # 工具插件（開發時載入）
```

### 2. WebUI 註冊階段

**註冊順序**：
```go
// WebUI 插件註冊的執行順序
func (p *Platform) registerWebUIPlugins() error {
    for _, webPlugin := range p.webUIPlugins {
        // 1. 註冊路由
        if err := webPlugin.RegisterRoutes(p.webRouter); err != nil {
            return fmt.Errorf("failed to register routes for %s: %w", 
                webPlugin.Name(), err)
        }
        
        // 2. 註冊導覽節點
        if err := webPlugin.RegisterNavNodes(p.navTreeBuilder); err != nil {
            return fmt.Errorf("failed to register nav nodes for %s: %w", 
                webPlugin.Name(), err)
        }
        
        // 3. 註冊 UI 組件
        if err := webPlugin.RegisterComponents(p.componentRegistry); err != nil {
            return fmt.Errorf("failed to register components for %s: %w", 
                webPlugin.Name(), err)
        }
    }
    
    return nil
}
```

### 3. 路由掛載機制

**路由註冊行為**：
```go
// WebRouter 實作範例
type WebRouter struct {
    echo   *echo.Echo
    groups map[string]*echo.Group
}

func (wr *WebRouter) GET(path string, handler echo.HandlerFunc) {
    // 自動添加插件路由前綴
    fullPath := "/plugins" + path
    wr.echo.GET(fullPath, handler)
}

func (wr *WebRouter) Group(prefix string) contracts.RouterGroup {
    // 建立路由群組，支援中介層
    group := wr.echo.Group("/plugins" + prefix)
    return &RouterGroup{group: group}
}
```

**路由衝突處理**：
```go
// 路由衝突檢測與解決
func (wr *WebRouter) registerRoute(method, path string, handler echo.HandlerFunc) error {
    // 檢查路由是否已存在
    if wr.routeExists(method, path) {
        return fmt.Errorf("route conflict: %s %s already registered", method, path)
    }
    
    // 註冊路由
    wr.echo.Add(method, path, handler)
    return nil
}
```

### 4. 導覽樹建構

**NavTree 建構流程**：
```go
// NavTreeBuilder 實作
type NavTreeBuilder struct {
    nodes    map[string]*NavNode
    children map[string][]*NavNode
}

func (ntb *NavTreeBuilder) AddNode(id string, node contracts.NavNode) error {
    // 1. 驗證節點 ID 唯一性
    if _, exists := ntb.nodes[id]; exists {
        return fmt.Errorf("nav node %s already exists", id)
    }
    
    // 2. 設定預設值
    if node.Order == 0 {
        node.Order = 100 // 預設排序
    }
    
    // 3. 添加到節點映射
    ntb.nodes[id] = &node
    
    return nil
}

func (ntb *NavTreeBuilder) BuildTree() *NavTree {
    // 根據 Order 排序並建構樹狀結構
    return ntb.sortAndBuildTree()
}
```

### 5. 組件註冊系統

**ComponentRegistry 實作**：
```go
// 組件註冊表
type ComponentRegistry struct {
    partials  map[string]string           // 模板片段
    widgets   map[string]WidgetFunc       // 小工具
    themes    map[string]Theme            // 主題
    assets    map[string]Asset            // 靜態資源
}

func (cr *ComponentRegistry) RegisterPartial(name, templatePath string) error {
    // 1. 驗證模板檔案存在
    if !fileExists(templatePath) {
        return fmt.Errorf("template file not found: %s", templatePath)
    }
    
    // 2. 註冊模板片段
    cr.partials[name] = templatePath
    
    return nil
}
```

### 6. 生命週期整合

**與插件生命週期的整合**：
```go
// 插件生命週期中的 WebUI 註冊
func (p *Platform) initializePlugin(plugin contracts.Plugin) error {
    // 1. 標準插件初始化
    if err := plugin.Init(p.config); err != nil {
        return err
    }
    
    // 2. 檢查是否為 WebUI 插件
    if webPlugin, ok := plugin.(contracts.WebUIPlugin); ok {
        // 3. 執行 WebUI 註冊
        if err := p.registerWebUIPlugin(webPlugin); err != nil {
            return fmt.Errorf("failed to register WebUI plugin %s: %w", 
                plugin.Name(), err)
        }
    }
    
    return nil
}
```

## 支援的 Plugin 類型

WebUIPlugin 適用於以下插件分類：

### Web UI 擴展插件 (`plugins/web/`)
- **pages/**: 完整頁面插件
- **components/**: UI 組件插件
- **themes/**: 主題插件
- **widgets/**: 小工具插件

### 其他插件的 Web 擴展
任何實作 WebUIPlugin 介面的插件都可以提供 Web UI 功能：
- **core plugins**: 核心插件的管理界面
- **community plugins**: 社群插件的配置頁面
- **integration plugins**: 整合插件的監控面板

## 模板結構

### 頁面模板 (Page Templates)

完整頁面模板應包含：

```html
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}} - DetectViz</title>
    <meta name="description" content="{{.Description}}">
</head>
<body>
    <!-- Navigation -->
    {{template "navigation" .NavTree}}
    
    <!-- Breadcrumbs -->
    {{template "breadcrumbs" .Breadcrumbs}}
    
    <!-- Flash Messages -->
    {{range .Flash}}
        <div class="alert alert-{{.Type}}">{{.Message}}</div>
    {{end}}
    
    <!-- Main Content -->
    <main>
        {{template "content" .Data}}
    </main>
</body>
</html>
```

### 部分模板 (Partials)

可重複使用的模板片段：

```html
<!-- system-status-card.html -->
<div class="card">
    <div class="card-header">
        <h3 class="card-title">{{.title}}</h3>
    </div>
    <div class="card-body">
        <div class="row">
            <div class="col-md-3">
                <div class="info-box">
                    <span class="info-box-icon bg-info">
                        <i class="fas fa-server"></i>
                    </span>
                    <div class="info-box-content">
                        <span class="info-box-text">狀態</span>
                        <span class="info-box-number">{{.status}}</span>
                    </div>
                </div>
            </div>
            <!-- More columns... -->
        </div>
    </div>
</div>
```

### HTMX 整合

利用 HTMX 實現動態更新：

```html
<!-- 自動刷新的狀態卡片 -->
<div id="status-card" 
     hx-get="/api/system/status" 
     hx-trigger="every 30s"
     hx-target="#status-card"
     hx-swap="innerHTML">
    {{template "system-status-card" .}}
</div>

<!-- 表單提交 -->
<form hx-post="/api/my-plugin/config" 
      hx-target="#result"
      hx-swap="innerHTML">
    <input name="setting" type="text" value="{{.config.setting}}">
    <button type="submit">保存</button>
</form>
```

## 權限控制

### 導覽節點權限

```go
// 設定節點權限
navNode := contracts.NavNode{
    Permission: "my-plugin.view", // 需要的權限
}

// 在 handler 中檢查權限
func (p *MyPlugin) handlePage(ctx contracts.WebContext) error {
    if !ctx.HasPermission("view", "my-plugin") {
        return ctx.Redirect(403, "/forbidden")
    }
    // 處理請求...
}
```

### 路由層級權限

```go
// 在路由註冊時設定權限中介層
func (p *MyPlugin) RegisterRoutes(router contracts.WebRouter) error {
    // 受保護的路由群組
    adminGroup := router.Group("/my-plugin/admin")
    adminGroup.Use(PermissionMiddleware("my-plugin.admin"))
    adminGroup.GET("/settings", p.handleAdminSettings)
    
    return nil
}
```

## 最佳實踐

### 1. 命名規範
- Plugin ID 使用 kebab-case：`system-status`
- 路由路徑使用 kebab-case：`/system-status`
- CSS 類別使用 BEM 規範：`plugin-name__element--modifier`

### 2. 資源管理
- 使用版本化的靜態資源
- 提供 SRI (Subresource Integrity) 雜湊值
- 明確宣告依賴關係

### 3. 使用者體驗
- 提供適當的載入狀態
- 使用 HTMX 實現部分頁面更新
- 支援鍵盤導覽
- 提供適當的錯誤處理

### 4. 多語言支援
- 使用 i18n 鍵值而非硬編碼文字
- 支援從右到左 (RTL) 語言
- 提供文化相關的日期時間格式

## 實作範例

完整的 WebUI Plugin 實作範例請參考：
- `plugins/web/pages/system-status/plugin.go`
- `plugins/web/components/charts/plugin.go`
- `plugins/web/themes/dark-theme/plugin.go`

## 相關介面

- [Plugin](./plugin.md) - 基礎插件介面
- [Registry](./registry.md) - 插件註冊介面
- [Auth](./auth.md) - 認證與授權介面 