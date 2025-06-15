package fake

import (
	"context"
	"fmt"
	"sync"
	"time"

	"detectviz/pkg/platform/contracts"
	"detectviz/pkg/shared/log"
)

// FakeComposition is a mock implementation of plugin composition for testing.
// zh: FakeComposition 是用於測試的插件組合模擬實作。
type FakeComposition struct {
	mu           sync.RWMutex
	name         string
	version      string
	plugins      map[string]*FakeCompositionPlugin
	dependencies map[string][]string // plugin -> dependencies
	loadOrder    []string
	started      bool
	errors       map[string]error
}

// FakeCompositionPlugin represents a plugin within a composition.
// zh: FakeCompositionPlugin 代表組合中的插件。
type FakeCompositionPlugin struct {
	Name         string         `yaml:"name" json:"name"`
	Type         string         `yaml:"type" json:"type"`
	Config       map[string]any `yaml:"config" json:"config"`
	Dependencies []string       `yaml:"dependencies" json:"dependencies"`
	Enabled      bool           `yaml:"enabled" json:"enabled"`
	Priority     int            `yaml:"priority" json:"priority"`
	Metadata     *contracts.PluginMetadata
	Instance     contracts.Plugin
	initialized  bool
	started      bool
}

// NewFakeComposition creates a new fake composition instance.
// zh: NewFakeComposition 建立新的假組合實例。
func NewFakeComposition(name, version string) *FakeComposition {
	return &FakeComposition{
		name:         name,
		version:      version,
		plugins:      make(map[string]*FakeCompositionPlugin),
		dependencies: make(map[string][]string),
		loadOrder:    make([]string, 0),
		started:      false,
		errors:       make(map[string]error),
	}
}

// AddPlugin adds a plugin to the composition.
// zh: AddPlugin 將插件添加到組合中。
func (fc *FakeComposition) AddPlugin(name, pluginType string, config map[string]any, dependencies []string) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if _, exists := fc.plugins[name]; exists {
		return fmt.Errorf("plugin '%s' already exists in composition", name)
	}

	plugin := &FakeCompositionPlugin{
		Name:         name,
		Type:         pluginType,
		Config:       config,
		Dependencies: dependencies,
		Enabled:      true,
		Priority:     0,
		Metadata: &contracts.PluginMetadata{
			Name:        name,
			Type:        pluginType,
			Category:    "test",
			Description: fmt.Sprintf("Test plugin %s", name),
			Version:     "1.0.0",
			Enabled:     true,
		},
		Instance:    NewFakePlugin(name, "1.0.0", fmt.Sprintf("Test plugin %s", name)),
		initialized: false,
		started:     false,
	}

	fc.plugins[name] = plugin
	fc.dependencies[name] = dependencies

	return nil
}

// RemovePlugin removes a plugin from the composition.
// zh: RemovePlugin 從組合中移除插件。
func (fc *FakeComposition) RemovePlugin(name string) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	if _, exists := fc.plugins[name]; !exists {
		return fmt.Errorf("plugin '%s' not found in composition", name)
	}

	delete(fc.plugins, name)
	delete(fc.dependencies, name)

	// Remove from load order
	for i, pluginName := range fc.loadOrder {
		if pluginName == name {
			fc.loadOrder = append(fc.loadOrder[:i], fc.loadOrder[i+1:]...)
			break
		}
	}

	return nil
}

// GetPlugin returns a plugin from the composition.
// zh: GetPlugin 從組合中返回插件。
func (fc *FakeComposition) GetPlugin(name string) (*FakeCompositionPlugin, error) {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	plugin, exists := fc.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin '%s' not found in composition", name)
	}

	return plugin, nil
}

// ListPlugins returns all plugin names in the composition.
// zh: ListPlugins 返回組合中所有插件名稱。
func (fc *FakeComposition) ListPlugins() []string {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	names := make([]string, 0, len(fc.plugins))
	for name := range fc.plugins {
		names = append(names, name)
	}
	return names
}

// CalculateLoadOrder calculates the plugin load order based on dependencies.
// zh: CalculateLoadOrder 根據依賴關係計算插件載入順序。
func (fc *FakeComposition) CalculateLoadOrder() error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	// Simple topological sort implementation
	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	order := make([]string, 0)

	var visit func(string) error
	visit = func(name string) error {
		if visiting[name] {
			return fmt.Errorf("circular dependency detected involving plugin '%s'", name)
		}
		if visited[name] {
			return nil
		}

		visiting[name] = true

		// Visit dependencies first
		for _, dep := range fc.dependencies[name] {
			if _, exists := fc.plugins[dep]; !exists {
				return fmt.Errorf("dependency '%s' not found for plugin '%s'", dep, name)
			}
			if err := visit(dep); err != nil {
				return err
			}
		}

		visiting[name] = false
		visited[name] = true
		order = append(order, name)

		return nil
	}

	// Visit all plugins
	for name := range fc.plugins {
		if err := visit(name); err != nil {
			return err
		}
	}

	fc.loadOrder = order
	return nil
}

// GetLoadOrder returns the calculated load order.
// zh: GetLoadOrder 返回計算出的載入順序。
func (fc *FakeComposition) GetLoadOrder() []string {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	order := make([]string, len(fc.loadOrder))
	copy(order, fc.loadOrder)
	return order
}

// InitializePlugins initializes all plugins in the composition.
// zh: InitializePlugins 初始化組合中的所有插件。
func (fc *FakeComposition) InitializePlugins(ctx context.Context) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	logger := log.L(ctx)

	// Initialize plugins
	for name, plugin := range fc.plugins {
		if !plugin.Enabled {
			logger.Info("Skipping disabled plugin", "plugin", name)
			continue
		}

		// Check for simulated errors
		if err, exists := fc.errors[name]; exists {
			return fmt.Errorf("simulated error for plugin '%s': %w", name, err)
		}

		logger.Info("Initializing plugin", "plugin", name)
		if err := plugin.Instance.Init(plugin.Config); err != nil {
			return fmt.Errorf("failed to initialize plugin '%s': %w", name, err)
		}

		plugin.initialized = true
		logger.Info("Plugin initialized successfully", "plugin", name)
	}

	return nil
}

// StartPlugins starts all initialized plugins in the composition.
// zh: StartPlugins 啟動組合中所有已初始化的插件。
func (fc *FakeComposition) StartPlugins(ctx context.Context) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	logger := log.L(ctx)

	for name, plugin := range fc.plugins {
		if !plugin.Enabled || !plugin.initialized {
			continue
		}

		// Check for simulated errors
		if err, exists := fc.errors[name]; exists {
			return fmt.Errorf("simulated error for plugin '%s': %w", name, err)
		}

		logger.Info("Starting plugin", "plugin", name)

		// Check if plugin implements LifecycleAware
		if lifecyclePlugin, ok := plugin.Instance.(contracts.LifecycleAware); ok {
			if err := lifecyclePlugin.OnStart(); err != nil {
				return fmt.Errorf("failed to start plugin '%s': %w", name, err)
			}
		}

		plugin.started = true
		logger.Info("Plugin started successfully", "plugin", name)
	}

	fc.started = true
	return nil
}

// StopPlugins stops all running plugins in the composition.
// zh: StopPlugins 停止組合中所有運行的插件。
func (fc *FakeComposition) StopPlugins(ctx context.Context) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	logger := log.L(ctx)

	// Stop plugins in reverse order
	for i := len(fc.loadOrder) - 1; i >= 0; i-- {
		name := fc.loadOrder[i]
		plugin := fc.plugins[name]
		if !plugin.started {
			continue
		}

		logger.Info("Stopping plugin", "plugin", name)

		// Check if plugin implements LifecycleAware
		if lifecyclePlugin, ok := plugin.Instance.(contracts.LifecycleAware); ok {
			if err := lifecyclePlugin.OnStop(); err != nil {
				logger.Error("Failed to stop plugin", "plugin", name, "error", err)
				// Continue stopping other plugins
			}
		}

		plugin.started = false
		logger.Info("Plugin stopped", "plugin", name)
	}

	fc.started = false
	return nil
}

// ShutdownPlugins shuts down all plugins in the composition.
// zh: ShutdownPlugins 關閉組合中的所有插件。
func (fc *FakeComposition) ShutdownPlugins(ctx context.Context) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	logger := log.L(ctx)

	// Shutdown plugins in reverse order
	for i := len(fc.loadOrder) - 1; i >= 0; i-- {
		name := fc.loadOrder[i]
		plugin := fc.plugins[name]
		if !plugin.initialized {
			continue
		}

		logger.Info("Shutting down plugin", "plugin", name)

		// Check if plugin implements LifecycleAware
		if lifecyclePlugin, ok := plugin.Instance.(contracts.LifecycleAware); ok {
			if err := lifecyclePlugin.OnShutdown(); err != nil {
				logger.Error("Failed to shutdown plugin", "plugin", name, "error", err)
				// Continue shutting down other plugins
			}
		}

		if err := plugin.Instance.Shutdown(); err != nil {
			logger.Error("Failed to shutdown plugin instance", "plugin", name, "error", err)
		}

		plugin.initialized = false
		plugin.started = false
		logger.Info("Plugin shut down", "plugin", name)
	}

	return nil
}

// GetCompositionInfo returns information about the composition.
// zh: GetCompositionInfo 返回組合的資訊。
func (fc *FakeComposition) GetCompositionInfo() map[string]any {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	initializedCount := 0
	startedCount := 0
	enabledCount := 0

	for _, plugin := range fc.plugins {
		if plugin.Enabled {
			enabledCount++
		}
		if plugin.initialized {
			initializedCount++
		}
		if plugin.started {
			startedCount++
		}
	}

	return map[string]any{
		"name":                fc.name,
		"version":             fc.version,
		"total_plugins":       len(fc.plugins),
		"enabled_plugins":     enabledCount,
		"initialized_plugins": initializedCount,
		"started_plugins":     startedCount,
		"composition_started": fc.started,
	}
}

// SimulateError simulates an error for a specific plugin.
// zh: SimulateError 為特定插件模擬錯誤。
func (fc *FakeComposition) SimulateError(pluginName string, err error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fc.errors[pluginName] = err
}

// ClearError clears simulated error for a plugin.
// zh: ClearError 清除插件的模擬錯誤。
func (fc *FakeComposition) ClearError(pluginName string) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	delete(fc.errors, pluginName)
}

// ClearAllErrors clears all simulated errors.
// zh: ClearAllErrors 清除所有模擬錯誤。
func (fc *FakeComposition) ClearAllErrors() {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fc.errors = make(map[string]error)
}

// Reset resets the composition to initial state.
// zh: Reset 重置組合到初始狀態。
func (fc *FakeComposition) Reset() {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	fc.plugins = make(map[string]*FakeCompositionPlugin)
	fc.dependencies = make(map[string][]string)
	fc.loadOrder = make([]string, 0)
	fc.started = false
	fc.errors = make(map[string]error)
}

// calculateLoadOrderUnsafe calculates load order without locking (internal use).
// zh: calculateLoadOrderUnsafe 計算載入順序但不加鎖（內部使用）。
func (fc *FakeComposition) calculateLoadOrderUnsafe() error {
	// Simple topological sort implementation
	visited := make(map[string]bool)
	visiting := make(map[string]bool)
	order := make([]string, 0)

	var visit func(string) error
	visit = func(name string) error {
		if visiting[name] {
			return fmt.Errorf("circular dependency detected involving plugin '%s'", name)
		}
		if visited[name] {
			return nil
		}

		visiting[name] = true

		// Visit dependencies first
		for _, dep := range fc.dependencies[name] {
			if _, exists := fc.plugins[dep]; !exists {
				return fmt.Errorf("dependency '%s' not found for plugin '%s'", dep, name)
			}
			if err := visit(dep); err != nil {
				return err
			}
		}

		visiting[name] = false
		visited[name] = true
		order = append(order, name)

		return nil
	}

	// Visit all plugins
	for name := range fc.plugins {
		if err := visit(name); err != nil {
			return err
		}
	}

	fc.loadOrder = order
	return nil
}

// FakeCompositionConfig represents a composition configuration for testing.
// zh: FakeCompositionConfig 代表用於測試的組合配置。
type FakeCompositionConfig struct {
	Name        string                            `yaml:"name" json:"name"`
	Version     string                            `yaml:"version" json:"version"`
	Description string                            `yaml:"description" json:"description"`
	Plugins     map[string]*FakeCompositionPlugin `yaml:"plugins" json:"plugins"`
	Settings    map[string]any                    `yaml:"settings" json:"settings"`
	CreatedAt   time.Time                         `yaml:"created_at" json:"created_at"`
	UpdatedAt   time.Time                         `yaml:"updated_at" json:"updated_at"`
}

// NewFakeCompositionConfig creates a new fake composition configuration.
// zh: NewFakeCompositionConfig 建立新的假組合配置。
func NewFakeCompositionConfig(name, version, description string) *FakeCompositionConfig {
	now := time.Now()
	return &FakeCompositionConfig{
		Name:        name,
		Version:     version,
		Description: description,
		Plugins:     make(map[string]*FakeCompositionPlugin),
		Settings:    make(map[string]any),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// AddPluginConfig adds a plugin configuration to the composition config.
// zh: AddPluginConfig 將插件配置添加到組合配置中。
func (fcc *FakeCompositionConfig) AddPluginConfig(name, pluginType string, config map[string]any, dependencies []string, enabled bool, priority int) {
	fcc.Plugins[name] = &FakeCompositionPlugin{
		Name:         name,
		Type:         pluginType,
		Config:       config,
		Dependencies: dependencies,
		Enabled:      enabled,
		Priority:     priority,
		Metadata: &contracts.PluginMetadata{
			Name:        name,
			Type:        pluginType,
			Category:    "test",
			Description: fmt.Sprintf("Test plugin %s", name),
			Version:     "1.0.0",
			Enabled:     enabled,
		},
	}
	fcc.UpdatedAt = time.Now()
}

// ToComposition converts the configuration to a FakeComposition.
// zh: ToComposition 將配置轉換為 FakeComposition。
func (fcc *FakeCompositionConfig) ToComposition() *FakeComposition {
	composition := NewFakeComposition(fcc.Name, fcc.Version)

	for name, pluginConfig := range fcc.Plugins {
		composition.AddPlugin(
			name,
			pluginConfig.Type,
			pluginConfig.Config,
			pluginConfig.Dependencies,
		)

		// Set enabled state and priority
		if plugin, err := composition.GetPlugin(name); err == nil {
			plugin.Enabled = pluginConfig.Enabled
			plugin.Priority = pluginConfig.Priority
		}
	}

	return composition
}
