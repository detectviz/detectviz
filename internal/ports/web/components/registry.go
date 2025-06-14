package components

import (
	"fmt"
	"sync"

	"detectviz/pkg/platform/contracts"
)

// Registry implements ComponentRegistry interface for managing UI components.
// zh: Registry 實作 ComponentRegistry 介面用於管理 UI 組件。
type Registry struct {
	partials map[string]string                  // partial name -> template path
	widgets  map[string]contracts.WidgetHandler // widget name -> handler
	themes   map[string]contracts.Theme         // theme name -> theme config
	assets   map[string]contracts.Asset         // asset name -> asset config
	mutex    sync.RWMutex
}

// NewRegistry creates a new component registry.
// zh: NewRegistry 建立新的組件註冊表。
func NewRegistry() *Registry {
	return &Registry{
		partials: make(map[string]string),
		widgets:  make(map[string]contracts.WidgetHandler),
		themes:   make(map[string]contracts.Theme),
		assets:   make(map[string]contracts.Asset),
	}
}

// RegisterPartial registers a partial template.
// zh: RegisterPartial 註冊部分模板。
func (r *Registry) RegisterPartial(name string, templatePath string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.partials[name]; exists {
		return fmt.Errorf("partial %s already registered", name)
	}

	r.partials[name] = templatePath
	return nil
}

// RegisterWidget registers a widget handler.
// zh: RegisterWidget 註冊小工具處理器。
func (r *Registry) RegisterWidget(name string, handler contracts.WidgetHandler) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.widgets[name]; exists {
		return fmt.Errorf("widget %s already registered", name)
	}

	r.widgets[name] = handler
	return nil
}

// RegisterTheme registers a theme configuration.
// zh: RegisterTheme 註冊主題配置。
func (r *Registry) RegisterTheme(name string, theme contracts.Theme) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.themes[name]; exists {
		return fmt.Errorf("theme %s already registered", name)
	}

	// Validate theme
	if theme.Name == "" {
		theme.Name = name
	}
	if theme.DisplayName == "" {
		theme.DisplayName = name
	}

	r.themes[name] = theme
	return nil
}

// RegisterAsset registers a static asset.
// zh: RegisterAsset 註冊靜態資源。
func (r *Registry) RegisterAsset(name string, asset contracts.Asset) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.assets[name]; exists {
		return fmt.Errorf("asset %s already registered", name)
	}

	// Validate asset
	if asset.Name == "" {
		asset.Name = name
	}
	if asset.Type == "" {
		// Infer type from path extension
		asset.Type = r.inferAssetType(asset.Path)
	}

	r.assets[name] = asset
	return nil
}

// GetPartial retrieves a partial template path by name.
// zh: GetPartial 根據名稱取得部分模板路徑。
func (r *Registry) GetPartial(name string) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	templatePath, exists := r.partials[name]
	if !exists {
		return "", fmt.Errorf("partial %s not found", name)
	}

	return templatePath, nil
}

// GetWidget retrieves a widget handler by name.
// zh: GetWidget 根據名稱取得小工具處理器。
func (r *Registry) GetWidget(name string) (contracts.WidgetHandler, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	handler, exists := r.widgets[name]
	if !exists {
		return nil, fmt.Errorf("widget %s not found", name)
	}

	return handler, nil
}

// ListPartials returns a list of all registered partial names.
// zh: ListPartials 回傳所有已註冊部分模板名稱的清單。
func (r *Registry) ListPartials() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.partials))
	for name := range r.partials {
		names = append(names, name)
	}
	return names
}

// ListWidgets returns a list of all registered widget names.
// zh: ListWidgets 回傳所有已註冊小工具名稱的清單。
func (r *Registry) ListWidgets() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.widgets))
	for name := range r.widgets {
		names = append(names, name)
	}
	return names
}

// Additional methods for enhanced functionality
// zh: 增強功能的額外方法

// GetTheme retrieves a theme configuration by name.
// zh: GetTheme 根據名稱取得主題配置。
func (r *Registry) GetTheme(name string) (contracts.Theme, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	theme, exists := r.themes[name]
	if !exists {
		return contracts.Theme{}, fmt.Errorf("theme %s not found", name)
	}

	return theme, nil
}

// ListThemes returns a list of all registered theme names.
// zh: ListThemes 回傳所有已註冊主題名稱的清單。
func (r *Registry) ListThemes() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.themes))
	for name := range r.themes {
		names = append(names, name)
	}
	return names
}

// GetAsset retrieves an asset configuration by name.
// zh: GetAsset 根據名稱取得資源配置。
func (r *Registry) GetAsset(name string) (contracts.Asset, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	asset, exists := r.assets[name]
	if !exists {
		return contracts.Asset{}, fmt.Errorf("asset %s not found", name)
	}

	return asset, nil
}

// ListAssets returns a list of all registered asset names.
// zh: ListAssets 回傳所有已註冊資源名稱的清單。
func (r *Registry) ListAssets() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	names := make([]string, 0, len(r.assets))
	for name := range r.assets {
		names = append(names, name)
	}
	return names
}

// GetAssetsByType returns assets filtered by type.
// zh: GetAssetsByType 回傳按類型篩選的資源。
func (r *Registry) GetAssetsByType(assetType string) []contracts.Asset {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var assets []contracts.Asset
	for _, asset := range r.assets {
		if asset.Type == assetType {
			assets = append(assets, asset)
		}
	}
	return assets
}

// ExecuteWidget executes a widget and returns its output.
// zh: ExecuteWidget 執行小工具並回傳其輸出。
func (r *Registry) ExecuteWidget(name string, ctx contracts.WebContext, params map[string]any) (any, error) {
	handler, err := r.GetWidget(name)
	if err != nil {
		return nil, err
	}

	return handler(ctx, params)
}

// UnregisterPartial removes a partial template from the registry.
// zh: UnregisterPartial 從註冊表中移除部分模板。
func (r *Registry) UnregisterPartial(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.partials[name]; !exists {
		return fmt.Errorf("partial %s not found", name)
	}

	delete(r.partials, name)
	return nil
}

// UnregisterWidget removes a widget from the registry.
// zh: UnregisterWidget 從註冊表中移除小工具。
func (r *Registry) UnregisterWidget(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.widgets[name]; !exists {
		return fmt.Errorf("widget %s not found", name)
	}

	delete(r.widgets, name)
	return nil
}

// UnregisterTheme removes a theme from the registry.
// zh: UnregisterTheme 從註冊表中移除主題。
func (r *Registry) UnregisterTheme(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.themes[name]; !exists {
		return fmt.Errorf("theme %s not found", name)
	}

	delete(r.themes, name)
	return nil
}

// UnregisterAsset removes an asset from the registry.
// zh: UnregisterAsset 從註冊表中移除資源。
func (r *Registry) UnregisterAsset(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.assets[name]; !exists {
		return fmt.Errorf("asset %s not found", name)
	}

	delete(r.assets, name)
	return nil
}

// GetStats returns statistics about the component registry.
// zh: GetStats 回傳組件註冊表的統計資訊。
func (r *Registry) GetStats() map[string]any {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	stats := map[string]any{
		"total_partials": len(r.partials),
		"total_widgets":  len(r.widgets),
		"total_themes":   len(r.themes),
		"total_assets":   len(r.assets),
		"assets_by_type": r.getAssetsByTypeStats(),
	}

	return stats
}

// Helper methods
// zh: 輔助方法

// inferAssetType infers asset type from file path.
// zh: inferAssetType 從檔案路徑推斷資源類型。
func (r *Registry) inferAssetType(path string) string {
	// Simple file extension-based type inference
	if len(path) < 4 {
		return "unknown"
	}

	ext := path[len(path)-4:]
	switch ext {
	case ".css":
		return "css"
	case ".js":
		return "js"
	case ".png", ".jpg", ".gif", ".svg":
		return "image"
	case ".ttf", ".otf", ".woff", ".woff2":
		return "font"
	default:
		// Check for longer extensions
		if len(path) >= 5 {
			ext5 := path[len(path)-5:]
			switch ext5 {
			case ".jpeg":
				return "image"
			case ".woff", ".woff2":
				return "font"
			}
		}
		return "unknown"
	}
}

// getAssetsByTypeStats returns asset count grouped by type.
// zh: getAssetsByTypeStats 回傳按類型分組的資源數量。
func (r *Registry) getAssetsByTypeStats() map[string]int {
	typeCount := make(map[string]int)

	for _, asset := range r.assets {
		typeCount[asset.Type]++
	}

	return typeCount
}

// ValidateTheme validates a theme configuration.
// zh: ValidateTheme 驗證主題配置。
func (r *Registry) ValidateTheme(theme contracts.Theme) error {
	if theme.Name == "" {
		return fmt.Errorf("theme name is required")
	}

	if theme.DisplayName == "" {
		return fmt.Errorf("theme display name is required")
	}

	if theme.Version == "" {
		return fmt.Errorf("theme version is required")
	}

	// Validate CSS and JS files exist (could be enhanced with actual file checks)
	for _, cssFile := range theme.CSSFiles {
		if cssFile == "" {
			return fmt.Errorf("empty CSS file path in theme %s", theme.Name)
		}
	}

	for _, jsFile := range theme.JSFiles {
		if jsFile == "" {
			return fmt.Errorf("empty JS file path in theme %s", theme.Name)
		}
	}

	return nil
}

// ValidateAsset validates an asset configuration.
// zh: ValidateAsset 驗證資源配置。
func (r *Registry) ValidateAsset(asset contracts.Asset) error {
	if asset.Name == "" {
		return fmt.Errorf("asset name is required")
	}

	if asset.Path == "" && asset.URL == "" {
		return fmt.Errorf("asset path or URL is required")
	}

	if asset.Type == "" {
		return fmt.Errorf("asset type is required")
	}

	return nil
}

// ClearAll clears all registered components (for testing purposes).
// zh: ClearAll 清除所有已註冊的組件（用於測試）。
func (r *Registry) ClearAll() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.partials = make(map[string]string)
	r.widgets = make(map[string]contracts.WidgetHandler)
	r.themes = make(map[string]contracts.Theme)
	r.assets = make(map[string]contracts.Asset)
}
