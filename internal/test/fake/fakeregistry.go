package fake

import (
	"context"
	"fmt"
	"sync"

	"detectviz/pkg/platform/contracts"
)

// FakeRegistry is a mock implementation of plugin registry for testing.
// zh: FakeRegistry 是用於測試的插件註冊表模擬實作。
type FakeRegistry struct {
	mu       sync.RWMutex
	plugins  map[string]contracts.Plugin
	metadata map[string]*contracts.PluginMetadata
	configs  map[string]map[string]any
	started  map[string]bool
	errors   map[string]error // For simulating errors
}

// NewFakeRegistry creates a new fake registry instance.
// zh: NewFakeRegistry 建立新的假註冊表實例。
func NewFakeRegistry() *FakeRegistry {
	return &FakeRegistry{
		plugins:  make(map[string]contracts.Plugin),
		metadata: make(map[string]*contracts.PluginMetadata),
		configs:  make(map[string]map[string]any),
		started:  make(map[string]bool),
		errors:   make(map[string]error),
	}
}

// Register registers a plugin with the fake registry.
// zh: Register 在假註冊表中註冊插件。
func (fr *FakeRegistry) Register(plugin contracts.Plugin, metadata *contracts.PluginMetadata) error {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	name := plugin.Name()

	// Check for simulated errors
	if err, exists := fr.errors[name]; exists {
		return err
	}

	// Check if plugin already registered
	if _, exists := fr.plugins[name]; exists {
		return fmt.Errorf("plugin '%s' already registered", name)
	}

	fr.plugins[name] = plugin
	fr.metadata[name] = metadata
	fr.started[name] = false

	return nil
}

// Resolve resolves a plugin by name.
// zh: Resolve 根據名稱解析插件。
func (fr *FakeRegistry) Resolve(name string) (contracts.Plugin, error) {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	// Check for simulated errors
	if err, exists := fr.errors[name]; exists {
		return nil, err
	}

	plugin, exists := fr.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin '%s' not found", name)
	}

	return plugin, nil
}

// GetMetadata returns metadata for a plugin.
// zh: GetMetadata 返回插件的元數據。
func (fr *FakeRegistry) GetMetadata(name string) (*contracts.PluginMetadata, error) {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	metadata, exists := fr.metadata[name]
	if !exists {
		return nil, fmt.Errorf("metadata for plugin '%s' not found", name)
	}

	return metadata, nil
}

// ListPlugins returns all registered plugin names.
// zh: ListPlugins 返回所有已註冊的插件名稱。
func (fr *FakeRegistry) ListPlugins() []string {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	names := make([]string, 0, len(fr.plugins))
	for name := range fr.plugins {
		names = append(names, name)
	}
	return names
}

// InitializePlugin initializes a plugin with configuration.
// zh: InitializePlugin 使用配置初始化插件。
func (fr *FakeRegistry) InitializePlugin(name string, config map[string]any) error {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	// Check for simulated errors
	if err, exists := fr.errors[name]; exists {
		return err
	}

	plugin, exists := fr.plugins[name]
	if !exists {
		return fmt.Errorf("plugin '%s' not found", name)
	}

	// Store config
	fr.configs[name] = config

	// Initialize plugin
	return plugin.Init(config)
}

// StartPlugin starts a plugin.
// zh: StartPlugin 啟動插件。
func (fr *FakeRegistry) StartPlugin(name string) error {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	// Check for simulated errors
	if err, exists := fr.errors[name]; exists {
		return err
	}

	plugin, exists := fr.plugins[name]
	if !exists {
		return fmt.Errorf("plugin '%s' not found", name)
	}

	// Check if plugin implements LifecycleAware
	if lifecyclePlugin, ok := plugin.(contracts.LifecycleAware); ok {
		if err := lifecyclePlugin.OnStart(); err != nil {
			return fmt.Errorf("failed to start plugin '%s': %w", name, err)
		}
	}

	fr.started[name] = true
	return nil
}

// StopPlugin stops a plugin.
// zh: StopPlugin 停止插件。
func (fr *FakeRegistry) StopPlugin(name string) error {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	// Check for simulated errors
	if err, exists := fr.errors[name]; exists {
		return err
	}

	plugin, exists := fr.plugins[name]
	if !exists {
		return fmt.Errorf("plugin '%s' not found", name)
	}

	// Check if plugin implements LifecycleAware
	if lifecyclePlugin, ok := plugin.(contracts.LifecycleAware); ok {
		if err := lifecyclePlugin.OnStop(); err != nil {
			return fmt.Errorf("failed to stop plugin '%s': %w", name, err)
		}
	}

	fr.started[name] = false
	return nil
}

// IsPluginStarted checks if a plugin is started.
// zh: IsPluginStarted 檢查插件是否已啟動。
func (fr *FakeRegistry) IsPluginStarted(name string) bool {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	return fr.started[name]
}

// GetPluginConfig returns the configuration for a plugin.
// zh: GetPluginConfig 返回插件的配置。
func (fr *FakeRegistry) GetPluginConfig(name string) (map[string]any, bool) {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	config, exists := fr.configs[name]
	return config, exists
}

// SimulateError simulates an error for a specific plugin operation.
// zh: SimulateError 為特定插件操作模擬錯誤。
func (fr *FakeRegistry) SimulateError(pluginName string, err error) {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	fr.errors[pluginName] = err
}

// ClearError clears simulated error for a plugin.
// zh: ClearError 清除插件的模擬錯誤。
func (fr *FakeRegistry) ClearError(pluginName string) {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	delete(fr.errors, pluginName)
}

// ClearAllErrors clears all simulated errors.
// zh: ClearAllErrors 清除所有模擬錯誤。
func (fr *FakeRegistry) ClearAllErrors() {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	fr.errors = make(map[string]error)
}

// Reset resets the fake registry to initial state.
// zh: Reset 重置假註冊表到初始狀態。
func (fr *FakeRegistry) Reset() {
	fr.mu.Lock()
	defer fr.mu.Unlock()

	fr.plugins = make(map[string]contracts.Plugin)
	fr.metadata = make(map[string]*contracts.PluginMetadata)
	fr.configs = make(map[string]map[string]any)
	fr.started = make(map[string]bool)
	fr.errors = make(map[string]error)
}

// GetStats returns statistics about the fake registry.
// zh: GetStats 返回假註冊表的統計信息。
func (fr *FakeRegistry) GetStats() map[string]int {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	startedCount := 0
	for _, started := range fr.started {
		if started {
			startedCount++
		}
	}

	return map[string]int{
		"total_plugins":   len(fr.plugins),
		"started_plugins": startedCount,
		"stopped_plugins": len(fr.plugins) - startedCount,
		"error_count":     len(fr.errors),
	}
}

// FakePlugin is a simple plugin implementation for testing.
// zh: FakePlugin 是用於測試的簡單插件實作。
type FakePlugin struct {
	name        string
	version     string
	description string
	initialized bool
	config      map[string]any
	onRegister  func() error
	onStart     func() error
	onStop      func() error
	onShutdown  func() error
}

// NewFakePlugin creates a new fake plugin.
// zh: NewFakePlugin 建立新的假插件。
func NewFakePlugin(name, version, description string) *FakePlugin {
	return &FakePlugin{
		name:        name,
		version:     version,
		description: description,
		initialized: false,
	}
}

// Name returns the plugin name.
func (fp *FakePlugin) Name() string {
	return fp.name
}

// Version returns the plugin version.
func (fp *FakePlugin) Version() string {
	return fp.version
}

// Description returns the plugin description.
func (fp *FakePlugin) Description() string {
	return fp.description
}

// Init initializes the plugin.
func (fp *FakePlugin) Init(config any) error {
	if configMap, ok := config.(map[string]any); ok {
		fp.config = configMap
	} else {
		fp.config = make(map[string]any)
	}
	fp.initialized = true
	return nil
}

// Shutdown shuts down the plugin.
func (fp *FakePlugin) Shutdown() error {
	fp.initialized = false
	return nil
}

// IsInitialized returns whether the plugin is initialized.
func (fp *FakePlugin) IsInitialized() bool {
	return fp.initialized
}

// GetConfig returns the plugin configuration.
func (fp *FakePlugin) GetConfig() map[string]any {
	return fp.config
}

// SetLifecycleHandlers sets lifecycle event handlers.
// zh: SetLifecycleHandlers 設置生命週期事件處理器。
func (fp *FakePlugin) SetLifecycleHandlers(
	onRegister, onStart, onStop, onShutdown func() error,
) {
	fp.onRegister = onRegister
	fp.onStart = onStart
	fp.onStop = onStop
	fp.onShutdown = onShutdown
}

// OnRegister handles plugin registration.
func (fp *FakePlugin) OnRegister() error {
	if fp.onRegister != nil {
		return fp.onRegister()
	}
	return nil
}

// OnStart handles plugin start.
func (fp *FakePlugin) OnStart() error {
	if fp.onStart != nil {
		return fp.onStart()
	}
	return nil
}

// OnStop handles plugin stop.
func (fp *FakePlugin) OnStop() error {
	if fp.onStop != nil {
		return fp.onStop()
	}
	return nil
}

// OnShutdown handles plugin shutdown.
func (fp *FakePlugin) OnShutdown() error {
	if fp.onShutdown != nil {
		return fp.onShutdown()
	}
	return nil
}

// FakeHealthChecker implements health checking for fake plugins.
// zh: FakeHealthChecker 為假插件實作健康檢查。
type FakeHealthChecker struct {
	status  string
	message string
	metrics map[string]any
}

// NewFakeHealthChecker creates a new fake health checker.
// zh: NewFakeHealthChecker 建立新的假健康檢查器。
func NewFakeHealthChecker(status, message string) *FakeHealthChecker {
	return &FakeHealthChecker{
		status:  status,
		message: message,
		metrics: make(map[string]any),
	}
}

// CheckHealth returns fake health status.
func (fhc *FakeHealthChecker) CheckHealth(ctx context.Context) contracts.HealthStatus {
	return contracts.HealthStatus{
		Status:  fhc.status,
		Message: fhc.message,
		Details: fhc.metrics,
	}
}

// GetHealthMetrics returns fake health metrics.
func (fhc *FakeHealthChecker) GetHealthMetrics() map[string]any {
	return fhc.metrics
}

// SetHealthMetric sets a health metric.
// zh: SetHealthMetric 設置健康指標。
func (fhc *FakeHealthChecker) SetHealthMetric(key string, value any) {
	fhc.metrics[key] = value
}

// SetHealthStatus sets the health status.
// zh: SetHealthStatus 設置健康狀態。
func (fhc *FakeHealthChecker) SetHealthStatus(status, message string) {
	fhc.status = status
	fhc.message = message
}
