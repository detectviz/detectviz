package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"detectviz/pkg/platform/contracts"
)

// WebContext implements the WebContext interface for handling web requests.
// zh: WebContext 實作 WebContext 介面用於處理 Web 請求。
type WebContext struct {
	request   *http.Request
	response  http.ResponseWriter
	params    map[string]string
	user      *contracts.UserInfo
	templates *template.Template
	flashMsgs []contracts.FlashMessage
}

// NewWebContext creates a new WebContext instance.
// zh: NewWebContext 建立新的 WebContext 實例。
func NewWebContext(w http.ResponseWriter, r *http.Request) contracts.WebContext {
	return &WebContext{
		request:   r,
		response:  w,
		params:    make(map[string]string),
		flashMsgs: make([]contracts.FlashMessage, 0),
	}
}

// Request returns the HTTP request.
// zh: Request 回傳 HTTP 請求。
func (ctx *WebContext) Request() *http.Request {
	return ctx.request
}

// Response returns the HTTP response writer.
// zh: Response 回傳 HTTP 回應寫入器。
func (ctx *WebContext) Response() http.ResponseWriter {
	return ctx.response
}

// Param returns a URL parameter by name.
// zh: Param 根據名稱回傳 URL 參數。
func (ctx *WebContext) Param(name string) string {
	if value, exists := ctx.params[name]; exists {
		return value
	}

	// Extract from URL path (simple implementation)
	// In a real implementation, this would be handled by the router
	return ""
}

// Query returns a query parameter by name.
// zh: Query 根據名稱回傳查詢參數。
func (ctx *WebContext) Query(name string) string {
	return ctx.request.URL.Query().Get(name)
}

// FormValue returns a form value by name.
// zh: FormValue 根據名稱回傳表單值。
func (ctx *WebContext) FormValue(name string) string {
	return ctx.request.FormValue(name)
}

// JSON sends a JSON response.
// zh: JSON 發送 JSON 回應。
func (ctx *WebContext) JSON(code int, data any) error {
	ctx.response.Header().Set("Content-Type", "application/json")
	ctx.response.WriteHeader(code)

	encoder := json.NewEncoder(ctx.response)
	return encoder.Encode(data)
}

// HTML sends an HTML response.
// zh: HTML 發送 HTML 回應。
func (ctx *WebContext) HTML(code int, templateName string, data any) error {
	ctx.response.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.response.WriteHeader(code)

	if ctx.templates == nil {
		return fmt.Errorf("no templates configured")
	}

	return ctx.templates.ExecuteTemplate(ctx.response, templateName, data)
}

// Redirect sends a redirect response.
// zh: Redirect 發送重定向回應。
func (ctx *WebContext) Redirect(code int, url string) error {
	if code < 300 || code >= 400 {
		code = http.StatusFound
	}

	http.Redirect(ctx.response, ctx.request, url, code)
	return nil
}

// User returns the current user information.
// zh: User 回傳目前使用者資訊。
func (ctx *WebContext) User() *contracts.UserInfo {
	return ctx.user
}

// HasPermission checks if the current user has the specified permission.
// zh: HasPermission 檢查目前使用者是否具有指定權限。
func (ctx *WebContext) HasPermission(action, resource string) bool {
	if ctx.user == nil {
		return false
	}

	for _, permission := range ctx.user.Permissions {
		if permission.Action == action && permission.Resource == resource {
			// Check scope (simplified implementation)
			return len(permission.Scope) == 0 || contains(permission.Scope, "*")
		}
	}

	return false
}

// RenderPage renders a complete page template.
// zh: RenderPage 渲染完整的頁面模板。
func (ctx *WebContext) RenderPage(name string, data any) error {
	// Prepare template data
	templateData := ctx.prepareTemplateData(data)

	// Set content type
	ctx.response.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Load and execute template
	tmpl, err := ctx.loadTemplate(name + ".html")
	if err != nil {
		return fmt.Errorf("failed to load template %s: %w", name, err)
	}

	return tmpl.Execute(ctx.response, templateData)
}

// RenderPartial renders a partial template (for HTMX responses).
// zh: RenderPartial 渲染部分模板（用於 HTMX 回應）。
func (ctx *WebContext) RenderPartial(name string, data any) error {
	// Prepare template data
	templateData := ctx.prepareTemplateData(data)

	// Set content type
	ctx.response.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Load and execute partial template
	tmpl, err := ctx.loadTemplate("partials/" + name + ".html")
	if err != nil {
		return fmt.Errorf("failed to load partial template %s: %w", name, err)
	}

	return tmpl.Execute(ctx.response, templateData)
}

// Helper methods
// zh: 輔助方法

// SetUser sets the current user information.
// zh: SetUser 設定目前使用者資訊。
func (ctx *WebContext) SetUser(user *contracts.UserInfo) {
	ctx.user = user
}

// SetParam sets a URL parameter.
// zh: SetParam 設定 URL 參數。
func (ctx *WebContext) SetParam(name, value string) {
	ctx.params[name] = value
}

// AddFlash adds a flash message.
// zh: AddFlash 添加快閃訊息。
func (ctx *WebContext) AddFlash(msgType, title, message string) {
	flash := contracts.FlashMessage{
		Type:    msgType,
		Title:   title,
		Message: message,
		Timeout: 5, // Default 5 seconds
	}
	ctx.flashMsgs = append(ctx.flashMsgs, flash)
}

// GetFlashMessages returns all flash messages.
// zh: GetFlashMessages 回傳所有快閃訊息。
func (ctx *WebContext) GetFlashMessages() []contracts.FlashMessage {
	return ctx.flashMsgs
}

// ClearFlashMessages clears all flash messages.
// zh: ClearFlashMessages 清除所有快閃訊息。
func (ctx *WebContext) ClearFlashMessages() {
	ctx.flashMsgs = make([]contracts.FlashMessage, 0)
}

// prepareTemplateData prepares data for template rendering.
// zh: prepareTemplateData 為模板渲染準備資料。
func (ctx *WebContext) prepareTemplateData(data any) contracts.TemplateData {
	templateData := contracts.TemplateData{
		Title:       "DetectViz",
		Description: "Composable Observability Platform",
		User:        ctx.user,
		Flash:       ctx.flashMsgs,
		Data:        make(map[string]any),
	}

	// Handle different data types
	switch v := data.(type) {
	case contracts.TemplateData:
		return v
	case map[string]any:
		templateData.Data = v
		if title, ok := v["title"].(string); ok {
			templateData.Title = title
		}
		if desc, ok := v["description"].(string); ok {
			templateData.Description = desc
		}
	default:
		templateData.Data["content"] = data
	}

	return templateData
}

// loadTemplate loads a template by name.
// zh: loadTemplate 根據名稱載入模板。
func (ctx *WebContext) loadTemplate(name string) (*template.Template, error) {
	// This is a simplified implementation
	// In a real implementation, you would:
	// 1. Load templates from a template directory
	// 2. Cache compiled templates
	// 3. Support template inheritance/layouts
	// 4. Handle template hot-reloading in development

	templatePath := filepath.Join("internal/ports/web/templates", name)

	// Check if template file exists
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		// Fallback to embedded template or create a simple one
		return ctx.createFallbackTemplate(name)
	}

	return tmpl, nil
}

// createFallbackTemplate creates a fallback template for testing.
// zh: createFallbackTemplate 建立用於測試的備用模板。
func (ctx *WebContext) createFallbackTemplate(name string) (*template.Template, error) {
	// Simple fallback template
	var templateContent string

	if strings.Contains(name, "partial") {
		// Partial template
		templateContent = `<div class="partial-content">
	<h3>{{.Title}}</h3>
	<p>{{.Description}}</p>
	{{range $key, $value := .Data}}
		<div><strong>{{$key}}:</strong> {{$value}}</div>
	{{end}}
</div>`
	} else {
		// Full page template
		templateContent = `<!DOCTYPE html>
<html>
<head>
	<title>{{.Title}} - DetectViz</title>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<meta name="description" content="{{.Description}}">
</head>
<body>
	<div class="container">
		<header>
			<h1>{{.Title}}</h1>
			{{if .User}}
				<div class="user-info">Welcome, {{.User.DisplayName}}</div>
			{{end}}
		</header>
		
		{{range .Flash}}
			<div class="alert alert-{{.Type}}">
				{{if .Title}}<strong>{{.Title}}</strong>{{end}}
				{{.Message}}
			</div>
		{{end}}
		
		<main>
			{{range $key, $value := .Data}}
				<div class="data-item">
					<strong>{{$key}}:</strong> {{$value}}
				</div>
			{{end}}
		</main>
	</div>
</body>
</html>`
	}

	return template.New(name).Parse(templateContent)
}

// Utility functions
// zh: 工具函式

// contains checks if a slice contains a string.
// zh: contains 檢查切片是否包含字串。
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
