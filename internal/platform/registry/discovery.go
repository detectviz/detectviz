package registry

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"plugin"
	"runtime"
	"strings"

	"detectviz/pkg/platform/contracts"
)

// PluginDiscovery handles automatic discovery and registration of plugins.
// zh: PluginDiscovery 處理插件的自動發現和註冊。
type PluginDiscovery struct {
	registry    *Manager
	basePath    string
	pluginPaths []string
	discovered  map[string]*DiscoveredPlugin
	loadedSOs   map[string]*plugin.Plugin
}

// DiscoveredPlugin represents a discovered plugin.
// zh: DiscoveredPlugin 代表已發現的插件。
type DiscoveredPlugin struct {
	Name         string `json:"name"`
	Path         string `json:"path"`
	Type         string `json:"type"`
	Category     string `json:"category"`
	RegisterFunc string `json:"register_func"`
	PackageName  string `json:"package_name"`
	SourceFile   string `json:"source_file"`
	Loaded       bool   `json:"loaded"`
	Error        string `json:"error,omitempty"`
}

// DiscoveryConfig contains configuration for plugin discovery.
// zh: DiscoveryConfig 包含插件發現的配置。
type DiscoveryConfig struct {
	BasePath       string   `json:"base_path"`
	ScanPaths      []string `json:"scan_paths"`
	ExcludePaths   []string `json:"exclude_paths"`
	AutoRegister   bool     `json:"auto_register"`
	LoadSharedLibs bool     `json:"load_shared_libs"`
	ScanGoSources  bool     `json:"scan_go_sources"`
}

// NewPluginDiscovery creates a new plugin discovery instance.
// zh: NewPluginDiscovery 建立新的插件發現實例。
func NewPluginDiscovery(registry *Manager, config *DiscoveryConfig) *PluginDiscovery {
	if config == nil {
		config = &DiscoveryConfig{
			BasePath:       ".",
			ScanPaths:      []string{"plugins"},
			AutoRegister:   true,
			LoadSharedLibs: false,
			ScanGoSources:  true,
		}
	}

	return &PluginDiscovery{
		registry:    registry,
		basePath:    config.BasePath,
		pluginPaths: config.ScanPaths,
		discovered:  make(map[string]*DiscoveredPlugin),
		loadedSOs:   make(map[string]*plugin.Plugin),
	}
}

// DiscoverPlugins scans for plugins in the configured paths.
// zh: DiscoverPlugins 在配置的路徑中掃描插件。
func (pd *PluginDiscovery) DiscoverPlugins() ([]*DiscoveredPlugin, error) {
	var allPlugins []*DiscoveredPlugin

	for _, scanPath := range pd.pluginPaths {
		fullPath := filepath.Join(pd.basePath, scanPath)

		plugins, err := pd.scanDirectory(fullPath)
		if err != nil {
			return nil, fmt.Errorf("failed to scan directory %s: %w", fullPath, err)
		}

		allPlugins = append(allPlugins, plugins...)
	}

	// Store discovered plugins
	for _, plugin := range allPlugins {
		pd.discovered[plugin.Name] = plugin
	}

	return allPlugins, nil
}

// scanDirectory recursively scans a directory for plugin files.
// zh: scanDirectory 遞歸掃描目錄中的插件檔案。
func (pd *PluginDiscovery) scanDirectory(dirPath string) ([]*DiscoveredPlugin, error) {
	var plugins []*DiscoveredPlugin

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Look for plugin.go files
		if strings.HasSuffix(path, "plugin.go") {
			plugin, err := pd.analyzePluginFile(path)
			if err != nil {
				// Log error but continue scanning
				fmt.Printf("Error analyzing plugin file %s: %v\n", path, err)
				return nil
			}

			if plugin != nil {
				plugins = append(plugins, plugin)
			}
		}

		// Look for .so files if enabled
		if strings.HasSuffix(path, ".so") {
			plugin, err := pd.analyzeSharedLibrary(path)
			if err != nil {
				fmt.Printf("Error analyzing shared library %s: %v\n", path, err)
				return nil
			}

			if plugin != nil {
				plugins = append(plugins, plugin)
			}
		}

		return nil
	})

	return plugins, err
}

// analyzePluginFile analyzes a Go plugin source file.
// zh: analyzePluginFile 分析 Go 插件原始檔。
func (pd *PluginDiscovery) analyzePluginFile(filePath string) (*DiscoveredPlugin, error) {
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}

	discovered := &DiscoveredPlugin{
		Path:        filePath,
		SourceFile:  filePath,
		PackageName: node.Name.Name,
		Loaded:      false,
	}

	// Determine plugin category from path
	discovered.Category = pd.determineCategory(filePath)

	// Determine plugin type from path
	discovered.Type = pd.determineType(filePath)

	// Look for Register function
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Name.Name == "Register" {
				discovered.RegisterFunc = "Register"
				discovered.Name = pd.extractPluginName(filePath, node)
				break
			}
		}
	}

	// If no Register function found, skip this file
	if discovered.RegisterFunc == "" {
		return nil, nil
	}

	return discovered, nil
}

// analyzeSharedLibrary analyzes a shared library (.so) file.
// zh: analyzeSharedLibrary 分析共享函式庫 (.so) 檔案。
func (pd *PluginDiscovery) analyzeSharedLibrary(filePath string) (*DiscoveredPlugin, error) {
	// Only supported on Linux/Unix systems
	if runtime.GOOS == "windows" {
		return nil, nil
	}

	discovered := &DiscoveredPlugin{
		Path:       filePath,
		SourceFile: filePath,
		Category:   pd.determineCategory(filePath),
		Type:       pd.determineType(filePath),
		Name:       pd.extractNameFromPath(filePath),
		Loaded:     false,
	}

	return discovered, nil
}

// determineCategory determines plugin category from file path.
// zh: determineCategory 從檔案路徑判斷插件類別。
func (pd *PluginDiscovery) determineCategory(filePath string) string {
	if strings.Contains(filePath, "plugins/core") {
		return "core"
	}
	if strings.Contains(filePath, "plugins/community") {
		return "community"
	}
	if strings.Contains(filePath, "plugins/custom") {
		return "custom"
	}
	return "unknown"
}

// determineType determines plugin type from file path.
// zh: determineType 從檔案路徑判斷插件類型。
func (pd *PluginDiscovery) determineType(filePath string) string {
	if strings.Contains(filePath, "/auth/") {
		return "auth"
	}
	if strings.Contains(filePath, "/importers/") {
		return "importer"
	}
	if strings.Contains(filePath, "/exporters/") {
		return "exporter"
	}
	if strings.Contains(filePath, "/middleware/") {
		return "middleware"
	}
	if strings.Contains(filePath, "/integrations/") {
		return "integration"
	}
	if strings.Contains(filePath, "/tools/") {
		return "tool"
	}
	return "unknown"
}

// extractPluginName extracts plugin name from source code or path.
// zh: extractPluginName 從原始碼或路徑提取插件名稱。
func (pd *PluginDiscovery) extractPluginName(filePath string, node *ast.File) string {
	// Look for string literals in RegisterPlugin calls
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			if funcDecl.Body != nil {
				for _, stmt := range funcDecl.Body.List {
					if name := pd.extractNameFromCallExpr(stmt); name != "" {
						return name
					}
				}
			}
		}
	}

	// Fallback to extracting from path
	return pd.extractNameFromPath(filePath)
}

// extractNameFromCallExpr extracts plugin name from RegisterPlugin call expressions.
// zh: extractNameFromCallExpr 從 RegisterPlugin 呼叫表達式提取插件名稱。
func (pd *PluginDiscovery) extractNameFromCallExpr(stmt ast.Stmt) string {
	var name string

	// Handle expression statements (direct calls)
	if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
		name = pd.extractNameFromCall(exprStmt.X)
	}

	// Handle return statements (return RegisterPlugin(...))
	if retStmt, ok := stmt.(*ast.ReturnStmt); ok {
		for _, result := range retStmt.Results {
			if foundName := pd.extractNameFromCall(result); foundName != "" {
				name = foundName
				break
			}
		}
	}

	return name
}

// extractNameFromCall extracts plugin name from any call expression
// zh: extractNameFromCall 從任何呼叫表達式提取插件名稱
func (pd *PluginDiscovery) extractNameFromCall(expr ast.Expr) string {
	if callExpr, ok := expr.(*ast.CallExpr); ok {
		if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if selectorExpr.Sel.Name == "RegisterPlugin" && len(callExpr.Args) > 0 {
				if basicLit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
					// Remove quotes from string literal
					return strings.Trim(basicLit.Value, "\"")
				}
			}
		}
	}
	return ""
}

// extractNameFromPath extracts plugin name from file path.
// zh: extractNameFromPath 從檔案路徑提取插件名稱。
func (pd *PluginDiscovery) extractNameFromPath(filePath string) string {
	// Extract directory name containing the plugin
	dir := filepath.Dir(filePath)
	parts := strings.Split(dir, string(filepath.Separator))

	if len(parts) > 0 {
		// Get the last directory name
		name := parts[len(parts)-1]

		// Add type suffix if not already present
		pluginType := pd.determineType(filePath)
		if pluginType != "unknown" && !strings.Contains(name, pluginType) {
			return name + "-" + pluginType
		}

		return name
	}

	return "unknown-plugin"
}

// RegisterDiscoveredPlugins registers all discovered plugins.
// zh: RegisterDiscoveredPlugins 註冊所有已發現的插件。
func (pd *PluginDiscovery) RegisterDiscoveredPlugins() error {
	for _, discoveredPlugin := range pd.discovered {
		if err := pd.registerPlugin(discoveredPlugin); err != nil {
			discoveredPlugin.Error = err.Error()
			fmt.Printf("Failed to register plugin %s: %v\n", discoveredPlugin.Name, err)
			continue
		}
		discoveredPlugin.Loaded = true
	}
	return nil
}

// registerPlugin registers a single discovered plugin.
// zh: registerPlugin 註冊單個已發現的插件。
func (pd *PluginDiscovery) registerPlugin(discoveredPlugin *DiscoveredPlugin) error {
	if strings.HasSuffix(discoveredPlugin.Path, ".so") {
		return pd.registerSharedLibraryPlugin(discoveredPlugin)
	}

	// For Go source plugins, we need to use reflection or code generation
	// Since we can't directly load Go source files at runtime,
	// we'll assume they've been compiled into the binary
	return pd.registerSourcePlugin(discoveredPlugin)
}

// registerSharedLibraryPlugin registers a plugin from a shared library.
// zh: registerSharedLibraryPlugin 從共享函式庫註冊插件。
func (pd *PluginDiscovery) registerSharedLibraryPlugin(discoveredPlugin *DiscoveredPlugin) error {
	if runtime.GOOS == "windows" {
		return fmt.Errorf("shared library plugins not supported on Windows")
	}

	plug, err := plugin.Open(discoveredPlugin.Path)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %w", discoveredPlugin.Path, err)
	}

	pd.loadedSOs[discoveredPlugin.Name] = plug

	// Look for Register symbol
	registerSymbol, err := plug.Lookup("Register")
	if err != nil {
		return fmt.Errorf("Register function not found in plugin %s: %w", discoveredPlugin.Path, err)
	}

	// Call Register function
	registerFunc, ok := registerSymbol.(func(contracts.Registry) error)
	if !ok {
		return fmt.Errorf("Register function has wrong signature in plugin %s", discoveredPlugin.Path)
	}

	return registerFunc(pd.registry)
}

// registerSourcePlugin registers a plugin compiled into the binary.
// zh: registerSourcePlugin 註冊編譯到二進位檔案中的插件。
func (pd *PluginDiscovery) registerSourcePlugin(discoveredPlugin *DiscoveredPlugin) error {
	// This is a placeholder for source-based plugin registration
	// In a real implementation, you would use build tags or code generation
	// to conditionally compile and register plugins

	fmt.Printf("Source plugin %s found but not automatically registered (requires manual registration)\n",
		discoveredPlugin.Name)

	return nil
}

// GetDiscoveredPlugins returns all discovered plugins.
// zh: GetDiscoveredPlugins 回傳所有已發現的插件。
func (pd *PluginDiscovery) GetDiscoveredPlugins() map[string]*DiscoveredPlugin {
	result := make(map[string]*DiscoveredPlugin)
	for k, v := range pd.discovered {
		result[k] = v
	}
	return result
}

// GetPluginInfo returns information about a specific discovered plugin.
// zh: GetPluginInfo 回傳特定已發現插件的資訊。
func (pd *PluginDiscovery) GetPluginInfo(pluginName string) (*DiscoveredPlugin, bool) {
	plugin, exists := pd.discovered[pluginName]
	return plugin, exists
}

// RefreshDiscovery re-scans for plugins and updates the discovery cache.
// zh: RefreshDiscovery 重新掃描插件並更新發現快取。
func (pd *PluginDiscovery) RefreshDiscovery() error {
	// Clear existing discovered plugins
	pd.discovered = make(map[string]*DiscoveredPlugin)

	// Re-discover plugins
	_, err := pd.DiscoverPlugins()
	return err
}

// UnloadPlugin unloads a previously loaded plugin (for shared libraries).
// zh: UnloadPlugin 卸載先前載入的插件（適用於共享函式庫）。
func (pd *PluginDiscovery) UnloadPlugin(pluginName string) error {
	if plug, exists := pd.loadedSOs[pluginName]; exists {
		// Note: Go's plugin package doesn't support unloading
		// This is a limitation of the Go plugin system
		_ = plug
		delete(pd.loadedSOs, pluginName)

		if discoveredPlugin, exists := pd.discovered[pluginName]; exists {
			discoveredPlugin.Loaded = false
		}
	}

	return nil
}
