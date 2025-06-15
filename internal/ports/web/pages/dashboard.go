package pages

import (
	"context"
	"fmt"
	"html/template"
	"runtime"
	"time"

	"detectviz/pkg/platform/contracts"
	"detectviz/pkg/shared/log"
)

// DashboardPlugin implements a basic dashboard WebUI plugin.
// zh: DashboardPlugin 實作基本的儀表板 WebUI 插件。
type DashboardPlugin struct {
	name        string
	version     string
	description string
	initialized bool
}

// DashboardData represents data for the dashboard page.
// zh: DashboardData 代表儀表板頁面的資料。
type DashboardData struct {
	Title       string                 `json:"title"`
	Timestamp   time.Time              `json:"timestamp"`
	SystemInfo  map[string]interface{} `json:"system_info"`
	Metrics     []DashboardMetric      `json:"metrics"`
	Alerts      []DashboardAlert       `json:"alerts"`
	PluginCount int                    `json:"plugin_count"`
}

// DashboardMetric represents a metric displayed on the dashboard.
// zh: DashboardMetric 代表儀表板上顯示的指標。
type DashboardMetric struct {
	Name        string  `json:"name"`
	Value       float64 `json:"value"`
	Unit        string  `json:"unit"`
	Description string  `json:"description"`
	Status      string  `json:"status"` // healthy, warning, critical
}

// DashboardAlert represents an alert displayed on the dashboard.
// zh: DashboardAlert 代表儀表板上顯示的警報。
type DashboardAlert struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Message      string    `json:"message"`
	Severity     string    `json:"severity"` // info, warning, error, critical
	Timestamp    time.Time `json:"timestamp"`
	Source       string    `json:"source"`
	Acknowledged bool      `json:"acknowledged"`
}

// NewDashboardPlugin creates a new dashboard plugin instance.
// zh: NewDashboardPlugin 建立新的儀表板插件實例。
func NewDashboardPlugin(config any) (contracts.Plugin, error) {
	return &DashboardPlugin{
		name:        "dashboard-webui",
		version:     "1.0.0",
		description: "DetectViz Dashboard WebUI Plugin",
		initialized: false,
	}, nil
}

// Plugin interface implementation
// zh: Plugin 介面實作

// Name returns the plugin name.
// zh: Name 回傳插件名稱。
func (d *DashboardPlugin) Name() string {
	return d.name
}

// Version returns the plugin version.
// zh: Version 回傳插件版本。
func (d *DashboardPlugin) Version() string {
	return d.version
}

// Description returns the plugin description.
// zh: Description 回傳插件描述。
func (d *DashboardPlugin) Description() string {
	return d.description
}

// Init initializes the dashboard plugin.
// zh: Init 初始化儀表板插件。
func (d *DashboardPlugin) Init(config any) error {
	d.initialized = true

	ctx := context.Background()
	log.L(ctx).Info("Dashboard WebUI plugin initialized")

	return nil
}

// Shutdown shuts down the dashboard plugin.
// zh: Shutdown 關閉儀表板插件。
func (d *DashboardPlugin) Shutdown() error {
	d.initialized = false

	ctx := context.Background()
	log.L(ctx).Info("Dashboard WebUI plugin shutdown")

	return nil
}

// WebUIPlugin interface implementation
// zh: WebUIPlugin 介面實作

// RegisterRoutes registers HTTP routes for the dashboard.
// zh: RegisterRoutes 為儀表板註冊 HTTP 路由。
func (d *DashboardPlugin) RegisterRoutes(router contracts.WebRouter) error {
	// Register dashboard page route
	if err := router.GET("/dashboard", d.handleDashboard); err != nil {
		return fmt.Errorf("failed to register dashboard route: %w", err)
	}

	// Register API routes
	if err := router.GET("/api/dashboard/data", d.handleDashboardAPI); err != nil {
		return fmt.Errorf("failed to register dashboard data API route: %w", err)
	}

	if err := router.GET("/api/dashboard/metrics", d.handleMetricsAPI); err != nil {
		return fmt.Errorf("failed to register metrics API route: %w", err)
	}

	if err := router.GET("/api/dashboard/alerts", d.handleAlertsAPI); err != nil {
		return fmt.Errorf("failed to register alerts API route: %w", err)
	}

	ctx := context.Background()
	log.L(ctx).Info("Dashboard routes registered", "routes", []string{"/dashboard", "/api/dashboard/data", "/api/dashboard/metrics", "/api/dashboard/alerts"})

	return nil
}

// RegisterNavNodes registers navigation nodes for the dashboard.
// zh: RegisterNavNodes 為儀表板註冊導覽節點。
func (d *DashboardPlugin) RegisterNavNodes(navtree contracts.NavTreeBuilder) error {
	// Add dashboard to main navigation
	dashboardNode := contracts.NavNode{
		ID:      "dashboard",
		Title:   "Dashboard",
		Icon:    "fas fa-tachometer-alt",
		URL:     "/dashboard",
		Order:   1,
		Enabled: true,
		Visible: true,
		Children: []contracts.NavNode{
			{
				ID:      "dashboard-overview",
				Title:   "Overview",
				URL:     "/dashboard",
				Order:   1,
				Enabled: true,
				Visible: true,
			},
			{
				ID:      "dashboard-metrics",
				Title:   "Metrics",
				URL:     "/dashboard/metrics",
				Order:   2,
				Enabled: true,
				Visible: true,
			},
			{
				ID:      "dashboard-alerts",
				Title:   "Alerts",
				URL:     "/dashboard/alerts",
				Order:   3,
				Enabled: true,
				Visible: true,
			},
		},
	}

	err := navtree.AddNode("dashboard", dashboardNode)
	if err != nil {
		return fmt.Errorf("failed to register dashboard navigation: %w", err)
	}

	ctx := context.Background()
	log.L(ctx).Info("Dashboard navigation registered")

	return nil
}

// RegisterComponents registers UI components for the dashboard.
// zh: RegisterComponents 為儀表板註冊 UI 組件。
func (d *DashboardPlugin) RegisterComponents(registry contracts.ComponentRegistry) error {
	// Register dashboard widgets
	widgets := []struct {
		name    string
		handler contracts.WidgetHandler
	}{
		{
			name:    "system-status",
			handler: d.systemStatusWidget,
		},
		{
			name:    "metrics-chart",
			handler: d.metricsChartWidget,
		},
		{
			name:    "alerts-panel",
			handler: d.alertsPanelWidget,
		},
		{
			name:    "plugin-status",
			handler: d.pluginStatusWidget,
		},
	}

	for _, widget := range widgets {
		if err := registry.RegisterWidget(widget.name, widget.handler); err != nil {
			return fmt.Errorf("failed to register widget %s: %w", widget.name, err)
		}
	}

	ctx := context.Background()
	log.L(ctx).Info("Dashboard widgets registered", "count", len(widgets))

	return nil
}

// Web handlers (WebHandler signature: func(ctx WebContext) error)
// zh: Web 處理器（WebHandler 簽名：func(ctx WebContext) error）

// handleDashboard handles the main dashboard page.
// zh: handleDashboard 處理主要的儀表板頁面。
func (d *DashboardPlugin) handleDashboard(ctx contracts.WebContext) error {
	// Get dashboard data
	data := d.getDashboardData(ctx.Request().Context())

	// Render dashboard template
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}} - DetectViz</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 0; padding: 20px; background-color: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; }
        .header { background: white; padding: 20px; border-radius: 8px; margin-bottom: 20px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
        .card { background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .metric { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
        .metric-value { font-size: 24px; font-weight: bold; }
        .status-healthy { color: #28a745; }
        .status-warning { color: #ffc107; }
        .status-critical { color: #dc3545; }
        .alert { padding: 10px; margin-bottom: 10px; border-radius: 4px; }
        .alert-info { background-color: #d1ecf1; border-left: 4px solid #17a2b8; }
        .alert-warning { background-color: #fff3cd; border-left: 4px solid #ffc107; }
        .alert-error { background-color: #f8d7da; border-left: 4px solid #dc3545; }
        .timestamp { color: #666; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.Title}}</h1>
            <p class="timestamp">Last updated: {{.Timestamp.Format "2006-01-02 15:04:05"}}</p>
        </div>
        
        <div class="grid">
            <div class="card">
                <h3>System Metrics</h3>
                {{range .Metrics}}
                <div class="metric">
                    <span>{{.Name}}</span>
                    <span class="metric-value status-{{.Status}}">{{.Value}} {{.Unit}}</span>
                </div>
                {{end}}
            </div>
            
            <div class="card">
                <h3>Recent Alerts</h3>
                {{range .Alerts}}
                <div class="alert alert-{{.Severity}}">
                    <strong>{{.Title}}</strong><br>
                    {{.Message}}<br>
                    <small class="timestamp">{{.Timestamp.Format "15:04:05"}} - {{.Source}}</small>
                </div>
                {{end}}
            </div>
            
            <div class="card">
                <h3>System Information</h3>
                <div class="metric">
                    <span>Active Plugins</span>
                    <span class="metric-value">{{.PluginCount}}</span>
                </div>
                {{range $key, $value := .SystemInfo}}
                <div class="metric">
                    <span>{{$key}}</span>
                    <span>{{$value}}</span>
                </div>
                {{end}}
            </div>
        </div>
    </div>
    
    <script>
        // Auto-refresh dashboard every 30 seconds
        setTimeout(function() {
            window.location.reload();
        }, 30000);
    </script>
</body>
</html>
`

	_, err := template.New("dashboard").Parse(tmpl)
	if err != nil {
		log.L(ctx.Request().Context()).Error("Dashboard template error", "error", err)
		return ctx.HTML(500, "error", map[string]any{"error": "Template error"})
	}

	return ctx.HTML(200, tmpl, data)
}

// handleDashboardAPI handles the dashboard data API.
// zh: handleDashboardAPI 處理儀表板資料 API。
func (d *DashboardPlugin) handleDashboardAPI(ctx contracts.WebContext) error {
	data := d.getDashboardData(ctx.Request().Context())

	response := map[string]any{
		"title":        data.Title,
		"timestamp":    data.Timestamp.Format(time.RFC3339),
		"plugin_count": data.PluginCount,
		"status":       "healthy",
	}

	return ctx.JSON(200, response)
}

// handleMetricsAPI handles the metrics API.
// zh: handleMetricsAPI 處理指標 API。
func (d *DashboardPlugin) handleMetricsAPI(ctx contracts.WebContext) error {
	data := d.getDashboardData(ctx.Request().Context())

	response := map[string]any{
		"metrics": []map[string]any{
			{"name": "CPU Usage", "value": 45.2, "unit": "%", "status": "healthy"},
			{"name": "Memory Usage", "value": 67.8, "unit": "%", "status": "warning"},
			{"name": "Disk Usage", "value": 23.1, "unit": "%", "status": "healthy"},
			{"name": "Active Connections", "value": len(data.Metrics), "unit": "", "status": "healthy"},
		},
	}

	return ctx.JSON(200, response)
}

// handleAlertsAPI handles the alerts API.
// zh: handleAlertsAPI 處理警報 API。
func (d *DashboardPlugin) handleAlertsAPI(ctx contracts.WebContext) error {
	data := d.getDashboardData(ctx.Request().Context())

	response := map[string]any{
		"alerts": []map[string]any{
			{
				"id":           "alert-001",
				"title":        "High Memory Usage",
				"message":      "Memory usage is above 80%",
				"severity":     "warning",
				"timestamp":    time.Now().Format(time.RFC3339),
				"source":       "system-monitor",
				"acknowledged": false,
			},
		},
		"total": len(data.Alerts),
	}

	return ctx.JSON(200, response)
}

// Widget handlers (WidgetHandler signature: func(ctx WebContext, params map[string]any) (any, error))
// zh: Widget 處理器（WidgetHandler 簽名：func(ctx WebContext, params map[string]any) (any, error)）

// systemStatusWidget handles the system status widget.
// zh: systemStatusWidget 處理系統狀態 widget。
func (d *DashboardPlugin) systemStatusWidget(ctx contracts.WebContext, params map[string]any) (any, error) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]any{
		"title":       "系統狀態",
		"status":      "running",
		"uptime":      d.getUptime(),
		"memory_used": memStats.Alloc,
		"cpu_count":   runtime.NumCPU(),
	}, nil
}

// metricsChartWidget handles the metrics chart widget.
// zh: metricsChartWidget 處理指標圖表 widget。
func (d *DashboardPlugin) metricsChartWidget(ctx contracts.WebContext, params map[string]any) (any, error) {
	return map[string]any{
		"chart_type": "line",
		"data": []map[string]any{
			{"time": time.Now().Add(-5 * time.Minute).Unix(), "cpu": 42.1, "memory": 65.3},
			{"time": time.Now().Add(-4 * time.Minute).Unix(), "cpu": 45.2, "memory": 67.8},
			{"time": time.Now().Add(-3 * time.Minute).Unix(), "cpu": 43.8, "memory": 66.1},
			{"time": time.Now().Add(-2 * time.Minute).Unix(), "cpu": 47.3, "memory": 69.2},
			{"time": time.Now().Add(-1 * time.Minute).Unix(), "cpu": 44.6, "memory": 68.5},
		},
	}, nil
}

// alertsPanelWidget handles the alerts panel widget.
// zh: alertsPanelWidget 處理警報面板 widget。
func (d *DashboardPlugin) alertsPanelWidget(ctx contracts.WebContext, params map[string]any) (any, error) {
	return map[string]any{
		"alerts": []map[string]any{
			{
				"id":        "alert-001",
				"title":     "High Memory Usage",
				"severity":  "warning",
				"timestamp": time.Now().Add(-10 * time.Minute).Unix(),
			},
			{
				"id":        "alert-002",
				"title":     "Disk Space Low",
				"severity":  "info",
				"timestamp": time.Now().Add(-25 * time.Minute).Unix(),
			},
		},
		"total": 2,
	}, nil
}

// pluginStatusWidget handles the plugin status widget.
// zh: pluginStatusWidget 處理插件狀態 widget。
func (d *DashboardPlugin) pluginStatusWidget(ctx contracts.WebContext, params map[string]any) (any, error) {
	return map[string]any{
		"plugins": []map[string]any{
			{"name": "prometheus-importer", "status": "healthy", "version": "1.0.0"},
			{"name": "influxdb-exporter", "status": "healthy", "version": "1.0.0"},
			{"name": "jwt-authenticator", "status": "healthy", "version": "1.0.0"},
		},
		"total": 3,
	}, nil
}

// Helper methods
// zh: 輔助方法

// getDashboardData generates dashboard data.
// zh: getDashboardData 產生儀表板資料。
func (d *DashboardPlugin) getDashboardData(ctx context.Context) *DashboardData {
	return &DashboardData{
		Title:     "DetectViz Dashboard",
		Timestamp: time.Now(),
		SystemInfo: map[string]interface{}{
			"Version":    "1.0.0",
			"Uptime":     d.getUptime(),
			"Go Version": runtime.Version(),
			"Platform":   runtime.GOOS + "/" + runtime.GOARCH,
		},
		Metrics: []DashboardMetric{
			{Name: "CPU Usage", Value: 45.2, Unit: "%", Description: "Current CPU utilization", Status: "healthy"},
			{Name: "Memory Usage", Value: 67.8, Unit: "%", Description: "Current memory utilization", Status: "warning"},
			{Name: "Disk Usage", Value: 23.1, Unit: "%", Description: "Current disk utilization", Status: "healthy"},
			{Name: "Active Connections", Value: 12, Unit: "", Description: "Number of active connections", Status: "healthy"},
		},
		Alerts: []DashboardAlert{
			{
				ID:           "alert-001",
				Title:        "High Memory Usage",
				Message:      "Memory usage is above 80%",
				Severity:     "warning",
				Timestamp:    time.Now().Add(-10 * time.Minute),
				Source:       "system-monitor",
				Acknowledged: false,
			},
			{
				ID:           "alert-002",
				Title:        "Disk Space Low",
				Message:      "Available disk space is below 20%",
				Severity:     "info",
				Timestamp:    time.Now().Add(-25 * time.Minute),
				Source:       "disk-monitor",
				Acknowledged: true,
			},
		},
		PluginCount: 3,
	}
}

// getUptime returns a mock uptime string.
// zh: getUptime 回傳模擬的運行時間字串。
func (d *DashboardPlugin) getUptime() string {
	// This is a mock implementation
	// In a real implementation, you would track the actual start time
	return "2h 15m"
}

// Register registers the dashboard plugin with the registry.
// zh: Register 向註冊表註冊儀表板插件。
func Register(registry contracts.Registry) error {
	return registry.RegisterPlugin("dashboard-webui", NewDashboardPlugin)
}
