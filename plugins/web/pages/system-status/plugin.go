package systemstatus

import (
	"fmt"
	"runtime"
	"time"

	"detectviz/pkg/platform/contracts"
)

// SystemStatusPlugin provides system status information web UI.
// zh: SystemStatusPlugin 提供系統狀態資訊的 Web UI。
type SystemStatusPlugin struct {
	name        string
	version     string
	description string
	config      *Config
	initialized bool
}

// Config contains configuration for the system status plugin.
// zh: Config 包含系統狀態插件的配置。
type Config struct {
	Title       string `yaml:"title" json:"title"`
	RefreshRate int    `yaml:"refresh_rate" json:"refresh_rate"` // seconds
	ShowMemory  bool   `yaml:"show_memory" json:"show_memory"`
	ShowCPU     bool   `yaml:"show_cpu" json:"show_cpu"`
	ShowPlugins bool   `yaml:"show_plugins" json:"show_plugins"`
}

// NewSystemStatusPlugin creates a new system status plugin instance.
// zh: NewSystemStatusPlugin 建立新的系統狀態插件實例。
func NewSystemStatusPlugin(config any) (contracts.Plugin, error) {
	var cfg Config

	// Set default configuration
	cfg = Config{
		Title:       "系統狀態",
		RefreshRate: 30,
		ShowMemory:  true,
		ShowCPU:     true,
		ShowPlugins: true,
	}

	// Parse configuration if provided
	if config != nil {
		if configMap, ok := config.(map[string]any); ok {
			if title, ok := configMap["title"].(string); ok {
				cfg.Title = title
			}
			if refreshRate, ok := configMap["refresh_rate"].(int); ok {
				cfg.RefreshRate = refreshRate
			}
			if showMemory, ok := configMap["show_memory"].(bool); ok {
				cfg.ShowMemory = showMemory
			}
			if showCPU, ok := configMap["show_cpu"].(bool); ok {
				cfg.ShowCPU = showCPU
			}
			if showPlugins, ok := configMap["show_plugins"].(bool); ok {
				cfg.ShowPlugins = showPlugins
			}
		}
	}

	return &SystemStatusPlugin{
		name:        "system-status",
		version:     "1.0.0",
		description: "System status monitoring web UI plugin",
		config:      &cfg,
		initialized: false,
	}, nil
}

// Plugin interface implementation
// zh: Plugin 介面實作

// Name returns the plugin name.
// zh: Name 回傳插件名稱。
func (p *SystemStatusPlugin) Name() string {
	return p.name
}

// Version returns the plugin version.
// zh: Version 回傳插件版本。
func (p *SystemStatusPlugin) Version() string {
	return p.version
}

// Description returns the plugin description.
// zh: Description 回傳插件描述。
func (p *SystemStatusPlugin) Description() string {
	return p.description
}

// Init initializes the plugin.
// zh: Init 初始化插件。
func (p *SystemStatusPlugin) Init(config any) error {
	if p.initialized {
		return nil
	}

	p.initialized = true
	return nil
}

// Shutdown shuts down the plugin.
// zh: Shutdown 關閉插件。
func (p *SystemStatusPlugin) Shutdown() error {
	p.initialized = false
	return nil
}

// WebUIPlugin interface implementation
// zh: WebUIPlugin 介面實作

// RegisterRoutes registers web routes for the system status plugin.
// zh: RegisterRoutes 為系統狀態插件註冊 Web 路由。
func (p *SystemStatusPlugin) RegisterRoutes(router contracts.WebRouter) error {
	// Register main system status page
	if err := router.GET("/system/status", p.handleStatusPage); err != nil {
		return fmt.Errorf("failed to register status page route: %w", err)
	}

	// Register API endpoint for status data
	if err := router.GET("/api/system/status", p.handleStatusAPI); err != nil {
		return fmt.Errorf("failed to register status API route: %w", err)
	}

	// Register refresh endpoint
	if err := router.GET("/system/status/refresh", p.handleRefresh); err != nil {
		return fmt.Errorf("failed to register refresh route: %w", err)
	}

	return nil
}

// RegisterNavNodes registers navigation nodes for the system status plugin.
// zh: RegisterNavNodes 為系統狀態插件註冊導覽節點。
func (p *SystemStatusPlugin) RegisterNavNodes(navtree contracts.NavTreeBuilder) error {
	statusNode := contracts.NavNode{
		ID:         "system-status",
		Title:      p.config.Title,
		Icon:       "fas fa-heartbeat",
		URL:        "/system/status",
		Permission: "system.status.view",
		Order:      100,
		Visible:    true,
		Enabled:    true,
		Badge: &contracts.NavBadge{
			Text:  "Live",
			Color: "success",
			Style: "pill",
		},
	}

	if err := navtree.AddNode("system-status", statusNode); err != nil {
		return fmt.Errorf("failed to add system status nav node: %w", err)
	}

	return nil
}

// RegisterComponents registers UI components for the system status plugin.
// zh: RegisterComponents 為系統狀態插件註冊 UI 組件。
func (p *SystemStatusPlugin) RegisterComponents(registry contracts.ComponentRegistry) error {
	// Register status card widget
	if err := registry.RegisterWidget("system-status-card", p.statusCardWidget); err != nil {
		return fmt.Errorf("failed to register status card widget: %w", err)
	}

	// Register memory usage widget
	if err := registry.RegisterWidget("memory-usage", p.memoryUsageWidget); err != nil {
		return fmt.Errorf("failed to register memory usage widget: %w", err)
	}

	// Register plugin list widget
	if err := registry.RegisterWidget("plugin-list", p.pluginListWidget); err != nil {
		return fmt.Errorf("failed to register plugin list widget: %w", err)
	}

	return nil
}

// Web handlers
// zh: Web 處理器

// handleStatusPage handles the main system status page.
// zh: handleStatusPage 處理主要系統狀態頁面。
func (p *SystemStatusPlugin) handleStatusPage(ctx contracts.WebContext) error {
	data := contracts.TemplateData{
		Title:       p.config.Title,
		Description: "DetectViz 系統狀態監控",
		Data: map[string]any{
			"config":      p.config,
			"system_info": p.getSystemInfo(),
			"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	return ctx.RenderPage("system-status", data)
}

// handleStatusAPI handles the system status API endpoint.
// zh: handleStatusAPI 處理系統狀態 API 端點。
func (p *SystemStatusPlugin) handleStatusAPI(ctx contracts.WebContext) error {
	statusData := map[string]any{
		"system_info": p.getSystemInfo(),
		"timestamp":   time.Now().Unix(),
	}

	return ctx.JSON(200, statusData)
}

// handleRefresh handles the status refresh endpoint.
// zh: handleRefresh 處理狀態刷新端點。
func (p *SystemStatusPlugin) handleRefresh(ctx contracts.WebContext) error {
	// Return updated system status as HTML partial
	data := map[string]any{
		"system_info": p.getSystemInfo(),
		"timestamp":   time.Now().Format("15:04:05"),
	}

	return ctx.RenderPartial("system-status-content", data)
}

// Widget handlers
// zh: Widget 處理器

// statusCardWidget handles the system status card widget.
// zh: statusCardWidget 處理系統狀態卡片 widget。
func (p *SystemStatusPlugin) statusCardWidget(ctx contracts.WebContext, params map[string]any) (any, error) {
	return map[string]any{
		"title":       "系統狀態",
		"status":      "running",
		"uptime":      p.getUptime(),
		"memory_used": p.getMemoryUsage(),
		"cpu_count":   runtime.NumCPU(),
	}, nil
}

// memoryUsageWidget handles the memory usage widget.
// zh: memoryUsageWidget 處理記憶體使用量 widget。
func (p *SystemStatusPlugin) memoryUsageWidget(ctx contracts.WebContext, params map[string]any) (any, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]any{
		"allocated": memStats.Alloc,
		"total":     memStats.TotalAlloc,
		"sys":       memStats.Sys,
		"gc_count":  memStats.NumGC,
	}, nil
}

// pluginListWidget handles the plugin list widget.
// zh: pluginListWidget 處理插件清單 widget。
func (p *SystemStatusPlugin) pluginListWidget(ctx contracts.WebContext, params map[string]any) (any, error) {
	// Mock plugin data - in real implementation, this would come from registry
	plugins := []map[string]any{
		{
			"name":     "jwt-authenticator",
			"version":  "1.0.0",
			"status":   "active",
			"category": "core",
		},
		{
			"name":     "prometheus-importer",
			"version":  "1.0.0",
			"status":   "active",
			"category": "community",
		},
		{
			"name":     "system-status",
			"version":  "1.0.0",
			"status":   "active",
			"category": "web",
		},
	}

	return map[string]any{
		"plugins": plugins,
		"total":   len(plugins),
	}, nil
}

// Helper methods
// zh: 輔助方法

// getSystemInfo returns system information.
// zh: getSystemInfo 回傳系統資訊。
func (p *SystemStatusPlugin) getSystemInfo() map[string]any {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]any{
		"go_version":   runtime.Version(),
		"go_os":        runtime.GOOS,
		"go_arch":      runtime.GOARCH,
		"cpu_count":    runtime.NumCPU(),
		"goroutines":   runtime.NumGoroutine(),
		"memory_alloc": memStats.Alloc,
		"memory_sys":   memStats.Sys,
		"gc_count":     memStats.NumGC,
	}
}

// getUptime returns system uptime (placeholder).
// zh: getUptime 回傳系統運行時間（佔位符）。
func (p *SystemStatusPlugin) getUptime() string {
	// In real implementation, this would track actual uptime
	return "24h 35m 12s"
}

// getMemoryUsage returns memory usage percentage (placeholder).
// zh: getMemoryUsage 回傳記憶體使用百分比（佔位符）。
func (p *SystemStatusPlugin) getMemoryUsage() string {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Simple calculation - in real implementation, this would be more accurate
	usagePercent := float64(memStats.Alloc) / float64(memStats.Sys) * 100
	return fmt.Sprintf("%.1f%%", usagePercent)
}

// Register registers the system status plugin.
// zh: Register 註冊系統狀態插件。
func Register(registry contracts.Registry) error {
	return registry.RegisterPlugin("system-status", NewSystemStatusPlugin)
}

// RegisterWebUIPlugin registers the system status web UI plugin.
// zh: RegisterWebUIPlugin 註冊系統狀態 Web UI 插件。
func RegisterWebUIPlugin(registry contracts.WebUIRegistry) error {
	factory := func(config any) (contracts.WebUIPlugin, error) {
		plugin, err := NewSystemStatusPlugin(config)
		if err != nil {
			return nil, err
		}

		webUIPlugin, ok := plugin.(contracts.WebUIPlugin)
		if !ok {
			return nil, fmt.Errorf("plugin does not implement WebUIPlugin interface")
		}

		return webUIPlugin, nil
	}

	return registry.RegisterWebUIPlugin("system-status", factory)
}
