package contracts

import (
	"net/http"
)

// WebUIPlugin defines the interface for web UI extensions.
// zh: WebUIPlugin 定義 Web UI 擴展介面。
type WebUIPlugin interface {
	Plugin
	RegisterRoutes(router WebRouter) error
	RegisterNavNodes(navtree NavTreeBuilder) error
	RegisterComponents(registry ComponentRegistry) error
}

// WebRouter defines the interface for web route registration.
// zh: WebRouter 定義 Web 路由註冊介面。
type WebRouter interface {
	GET(path string, handler WebHandler) error
	POST(path string, handler WebHandler) error
	PUT(path string, handler WebHandler) error
	DELETE(path string, handler WebHandler) error
	Group(prefix string) WebRouter
}

// WebHandler defines the interface for web request handlers.
// zh: WebHandler 定義 Web 請求處理器介面。
type WebHandler func(ctx WebContext) error

// WebContext defines the interface for web request context.
// zh: WebContext 定義 Web 請求上下文介面。
type WebContext interface {
	Request() *http.Request
	Response() http.ResponseWriter
	Param(name string) string
	Query(name string) string
	FormValue(name string) string
	JSON(code int, data any) error
	HTML(code int, template string, data any) error
	Redirect(code int, url string) error

	// Authentication & Authorization
	User() *UserInfo
	HasPermission(action, resource string) bool

	// Template rendering
	RenderPage(name string, data any) error
	RenderPartial(name string, data any) error
}

// NavTreeBuilder defines the interface for navigation tree construction.
// zh: NavTreeBuilder 定義導覽樹建構介面。
type NavTreeBuilder interface {
	AddNode(id string, node NavNode) error
	AddChildNode(parentID string, id string, node NavNode) error
	RemoveNode(id string) error
	GetNode(id string) (*NavNode, error)
	SetNodePermission(id string, permission string) error
}

// NavNode represents a navigation node in the web UI.
// zh: NavNode 代表 Web UI 中的導覽節點。
type NavNode struct {
	ID         string            `json:"id"`
	Title      string            `json:"title"`
	Icon       string            `json:"icon,omitempty"`
	URL        string            `json:"url,omitempty"`
	Permission string            `json:"permission,omitempty"`
	Order      int               `json:"order"`
	Target     string            `json:"target,omitempty"` // _blank, _self, etc.
	Badge      *NavBadge         `json:"badge,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
	Children   []NavNode         `json:"children,omitempty"`
	Visible    bool              `json:"visible"`
	Enabled    bool              `json:"enabled"`
}

// NavBadge represents a badge for navigation nodes.
// zh: NavBadge 代表導覽節點的徽章。
type NavBadge struct {
	Text  string `json:"text"`
	Color string `json:"color"` // success, warning, danger, info, etc.
	Style string `json:"style"` // pill, square, etc.
}

// ComponentRegistry defines the interface for UI component registration.
// zh: ComponentRegistry 定義 UI 組件註冊介面。
type ComponentRegistry interface {
	RegisterPartial(name string, templatePath string) error
	RegisterWidget(name string, handler WidgetHandler) error
	RegisterTheme(name string, theme Theme) error
	RegisterAsset(name string, asset Asset) error
	GetPartial(name string) (string, error)
	GetWidget(name string) (WidgetHandler, error)
	ListPartials() []string
	ListWidgets() []string
}

// WidgetHandler defines the interface for widget handlers.
// zh: WidgetHandler 定義小工具處理器介面。
type WidgetHandler func(ctx WebContext, params map[string]any) (any, error)

// Theme represents a UI theme configuration.
// zh: Theme 代表 UI 主題配置。
type Theme struct {
	Name        string            `json:"name"`
	DisplayName string            `json:"display_name"`
	Version     string            `json:"version"`
	Author      string            `json:"author"`
	Description string            `json:"description"`
	CSSFiles    []string          `json:"css_files"`
	JSFiles     []string          `json:"js_files"`
	Variables   map[string]string `json:"variables"`
	Preview     string            `json:"preview,omitempty"`
}

// Asset represents a static asset (CSS, JS, images, etc.).
// zh: Asset 代表靜態資源（CSS、JS、圖片等）。
type Asset struct {
	Name         string            `json:"name"`
	Type         string            `json:"type"` // css, js, image, font, etc.
	Path         string            `json:"path"`
	URL          string            `json:"url,omitempty"`
	Version      string            `json:"version"`
	Integrity    string            `json:"integrity,omitempty"` // SRI hash
	Attributes   map[string]string `json:"attributes,omitempty"`
	Dependencies []string          `json:"dependencies,omitempty"`
}

// WebUIRegistry defines the interface for web UI plugin registration.
// zh: WebUIRegistry 定義 Web UI 插件註冊介面。
type WebUIRegistry interface {
	RegisterWebUIPlugin(name string, factory func(config any) (WebUIPlugin, error)) error
	GetWebUIPlugin(name string) (WebUIPlugin, error)
	ListWebUIPlugins() []string
}

// TemplateData represents data passed to templates.
// zh: TemplateData 代表傳遞給模板的資料。
type TemplateData struct {
	Title       string         `json:"title"`
	Description string         `json:"description,omitempty"`
	Keywords    []string       `json:"keywords,omitempty"`
	User        *UserInfo      `json:"user,omitempty"`
	NavTree     []NavNode      `json:"nav_tree,omitempty"`
	Breadcrumbs []Breadcrumb   `json:"breadcrumbs,omitempty"`
	Flash       []FlashMessage `json:"flash,omitempty"`
	Data        map[string]any `json:"data"`
	Config      map[string]any `json:"config,omitempty"`
}

// Breadcrumb represents a breadcrumb navigation item.
// zh: Breadcrumb 代表麵包屑導覽項目。
type Breadcrumb struct {
	Title string `json:"title"`
	URL   string `json:"url,omitempty"`
	Icon  string `json:"icon,omitempty"`
}

// FlashMessage represents a flash message for user feedback.
// zh: FlashMessage 代表用於使用者回饋的快閃訊息。
type FlashMessage struct {
	Type    string `json:"type"` // success, error, warning, info
	Title   string `json:"title,omitempty"`
	Message string `json:"message"`
	Timeout int    `json:"timeout,omitempty"` // Auto-dismiss timeout in seconds
}
