package composition

import (
	"context"
	"fmt"
	"sync"
	"time"

	"detectviz/pkg/platform/contracts"
	"detectviz/pkg/shared/log"
	"detectviz/plugins/core/logging/otelzap"
	// Import core logging plugin for auto-registration
)

// LifecycleManager implements the core lifecycle management functionality.
// zh: LifecycleManager 實作核心生命週期管理功能。
type LifecycleManager struct {
	components map[string]contracts.ComponentInfo
	status     contracts.LifecycleStatus
	listeners  []contracts.LifecycleListener
	resolver   contracts.DependencyResolver
	mutex      sync.RWMutex
}

// NewLifecycleManager creates a new lifecycle manager.
// zh: NewLifecycleManager 建立新的生命週期管理器。
func NewLifecycleManager(resolver contracts.DependencyResolver) *LifecycleManager {
	return &LifecycleManager{
		components: make(map[string]contracts.ComponentInfo),
		status:     contracts.StatusUninitialized,
		listeners:  make([]contracts.LifecycleListener, 0),
		resolver:   resolver,
	}
}

// Initialize initializes all registered components in dependency order.
// zh: Initialize 按依賴順序初始化所有已註冊的組件。
func (lm *LifecycleManager) Initialize(ctx context.Context) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	if lm.status != contracts.StatusUninitialized {
		return fmt.Errorf("lifecycle manager already initialized")
	}

	lm.status = contracts.StatusInitializing
	lm.publishEvent(contracts.EventInitialize, "lifecycle-manager", lm.status, nil)

	// Get components in dependency order
	components := make([]contracts.ComponentInfo, 0, len(lm.components))
	for _, comp := range lm.components {
		components = append(components, comp)
	}

	orderedComponents, err := lm.resolver.ResolveDependencies(components)
	if err != nil {
		lm.status = contracts.StatusError
		lm.publishEvent(contracts.EventError, "lifecycle-manager", lm.status, err)
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	// Initialize components in order
	for _, comp := range orderedComponents {
		if err := lm.initializeComponent(ctx, comp); err != nil {
			lm.status = contracts.StatusError
			lm.publishEvent(contracts.EventError, comp.Name, contracts.StatusError, err)
			return fmt.Errorf("failed to initialize component %s: %w", comp.Name, err)
		}
		comp.Status = contracts.StatusInitialized
		lm.components[comp.Name] = comp
	}

	lm.status = contracts.StatusInitialized
	lm.publishEvent(contracts.EventInitialize, "lifecycle-manager", lm.status, nil)
	return nil
}

// Start starts all initialized components in dependency order.
// zh: Start 按依賴順序啟動所有已初始化的組件。
func (lm *LifecycleManager) Start(ctx context.Context) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	if lm.status != contracts.StatusInitialized {
		return fmt.Errorf("lifecycle manager not initialized")
	}

	lm.status = contracts.StatusStarting
	lm.publishEvent(contracts.EventStart, "lifecycle-manager", lm.status, nil)

	// Start components in dependency order
	components := make([]contracts.ComponentInfo, 0, len(lm.components))
	for _, comp := range lm.components {
		components = append(components, comp)
	}

	orderedComponents, err := lm.resolver.ResolveDependencies(components)
	if err != nil {
		lm.status = contracts.StatusError
		lm.publishEvent(contracts.EventError, "lifecycle-manager", lm.status, err)
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	for _, comp := range orderedComponents {
		if err := lm.startComponent(ctx, comp); err != nil {
			lm.status = contracts.StatusError
			lm.publishEvent(contracts.EventError, comp.Name, contracts.StatusError, err)
			return fmt.Errorf("failed to start component %s: %w", comp.Name, err)
		}
		comp.Status = contracts.StatusRunning
		lm.components[comp.Name] = comp
	}

	lm.status = contracts.StatusRunning
	lm.publishEvent(contracts.EventStart, "lifecycle-manager", lm.status, nil)
	return nil
}

// Stop stops all running components in reverse dependency order.
// zh: Stop 按反向依賴順序停止所有運行中的組件。
func (lm *LifecycleManager) Stop(ctx context.Context) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	if lm.status != contracts.StatusRunning {
		return fmt.Errorf("lifecycle manager not running")
	}

	lm.status = contracts.StatusStopping
	lm.publishEvent(contracts.EventStop, "lifecycle-manager", lm.status, nil)

	// Stop components in reverse dependency order
	components := make([]contracts.ComponentInfo, 0, len(lm.components))
	for _, comp := range lm.components {
		components = append(components, comp)
	}

	orderedComponents, err := lm.resolver.ResolveDependencies(components)
	if err != nil {
		lm.status = contracts.StatusError
		lm.publishEvent(contracts.EventError, "lifecycle-manager", lm.status, err)
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	// Reverse the order for stopping
	for i := len(orderedComponents) - 1; i >= 0; i-- {
		comp := orderedComponents[i]
		if err := lm.stopComponent(ctx, comp); err != nil {
			lm.status = contracts.StatusError
			lm.publishEvent(contracts.EventError, comp.Name, contracts.StatusError, err)
			return fmt.Errorf("failed to stop component %s: %w", comp.Name, err)
		}
		comp.Status = contracts.StatusStopped
		lm.components[comp.Name] = comp
	}

	lm.status = contracts.StatusStopped
	lm.publishEvent(contracts.EventStop, "lifecycle-manager", lm.status, nil)
	return nil
}

// Shutdown shuts down all components and cleans up resources.
// zh: Shutdown 關閉所有組件並清理資源。
func (lm *LifecycleManager) Shutdown(ctx context.Context) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	lm.status = contracts.StatusShuttingDown
	lm.publishEvent(contracts.EventShutdown, "lifecycle-manager", lm.status, nil)

	// Shutdown components in reverse dependency order
	components := make([]contracts.ComponentInfo, 0, len(lm.components))
	for _, comp := range lm.components {
		components = append(components, comp)
	}

	orderedComponents, err := lm.resolver.ResolveDependencies(components)
	if err != nil {
		lm.status = contracts.StatusError
		lm.publishEvent(contracts.EventError, "lifecycle-manager", lm.status, err)
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	// Reverse the order for shutdown
	for i := len(orderedComponents) - 1; i >= 0; i-- {
		comp := orderedComponents[i]
		if err := lm.shutdownComponent(ctx, comp); err != nil {
			lm.status = contracts.StatusError
			lm.publishEvent(contracts.EventError, comp.Name, contracts.StatusError, err)
			return fmt.Errorf("failed to shutdown component %s: %w", comp.Name, err)
		}
		comp.Status = contracts.StatusShutdown
		lm.components[comp.Name] = comp
	}

	lm.status = contracts.StatusShutdown
	lm.publishEvent(contracts.EventShutdown, "lifecycle-manager", lm.status, nil)
	return nil
}

// GetStatus returns the current lifecycle status.
// zh: GetStatus 回傳當前生命週期狀態。
func (lm *LifecycleManager) GetStatus() contracts.LifecycleStatus {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	return lm.status
}

// RegisterComponent registers a component for lifecycle management.
// zh: RegisterComponent 註冊組件進行生命週期管理。
func (lm *LifecycleManager) RegisterComponent(info contracts.ComponentInfo) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	if _, exists := lm.components[info.Name]; exists {
		return fmt.Errorf("component %s already registered", info.Name)
	}

	lm.components[info.Name] = info
	return nil
}

// AddLifecycleListener adds a lifecycle event listener.
// zh: AddLifecycleListener 新增生命週期事件監聽器。
func (lm *LifecycleManager) AddLifecycleListener(listener contracts.LifecycleListener) {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	lm.listeners = append(lm.listeners, listener)
}

// Helper methods for component lifecycle operations
func (lm *LifecycleManager) initializeComponent(ctx context.Context, comp contracts.ComponentInfo) error {
	// TODO: Implement actual component initialization logic
	// This would involve calling the component's Init method if it implements LifecycleAware
	return nil
}

func (lm *LifecycleManager) startComponent(ctx context.Context, comp contracts.ComponentInfo) error {
	// TODO: Implement actual component start logic
	return nil
}

func (lm *LifecycleManager) stopComponent(ctx context.Context, comp contracts.ComponentInfo) error {
	// TODO: Implement actual component stop logic
	return nil
}

func (lm *LifecycleManager) shutdownComponent(ctx context.Context, comp contracts.ComponentInfo) error {
	// TODO: Implement actual component shutdown logic
	return nil
}

func (lm *LifecycleManager) publishEvent(eventType contracts.LifecycleEventType, component string, status contracts.LifecycleStatus, err error) {
	event := contracts.LifecycleEvent{
		Type:      eventType,
		Component: component,
		Timestamp: time.Now(),
		Status:    status,
		Metadata:  make(map[string]any),
	}

	if err != nil {
		event.Error = err.Error()
	}

	// Notify all listeners
	for _, listener := range lm.listeners {
		if listenerErr := listener.OnLifecycleEvent(event); listenerErr != nil {
			// Log the error but don't fail the lifecycle operation
			ctx := context.Background()
			log.L(ctx).Warn("Lifecycle listener error", "error", listenerErr, "component", component, "event", eventType)
		}
	}
}

// RegisterLoggerPlugin registers and initializes the logging plugin if available
// zh: RegisterLoggerPlugin 註冊並初始化日誌插件（如果可用）
func (lm *LifecycleManager) RegisterLoggerPlugin(ctx context.Context, registry contracts.Registry) error {
	// Use the global registry from otelzap package for plugin access
	globalRegistry := otelzap.GetGlobalRegistry()

	// Check if otelzap logging plugin is available
	pluginNames := globalRegistry.ListPlugins()
	for _, name := range pluginNames {
		if name == "otelzap" {
			plugin, err := globalRegistry.GetPlugin(name)
			if err != nil {
				log.L(ctx).Warn("Failed to get otelzap plugin", "error", err)
				return nil // Don't fail the entire startup if logger plugin fails
			}

			// Check if plugin implements LoggerProvider interface
			if loggerProvider, ok := plugin.(contracts.LoggerProvider); ok {
				// Initialize the logger plugin first
				if lifecyclePlugin, ok := plugin.(contracts.LifecycleAware); ok {
					if err := lifecyclePlugin.OnRegister(); err != nil {
						log.L(ctx).Warn("Failed to register otelzap plugin", "error", err)
						return nil
					}
				}

				// Set as global logger
				globalLogger := loggerProvider.Logger()
				log.SetGlobalLogger(globalLogger)

				log.L(ctx).Info("Successfully registered otelzap logging plugin as global logger")

				// Register as component in our local registry too if needed
				if err := registry.RegisterPlugin(name, otelzap.NewOtelZapPlugin); err != nil {
					log.L(ctx).Warn("Failed to register otelzap plugin in local registry", "error", err)
				}

				// Register as component
				comp := contracts.ComponentInfo{
					Name:   name,
					Type:   "logger-plugin",
					Status: contracts.StatusInitialized,
				}
				lm.components[name] = comp
				lm.publishEvent(contracts.EventInitialize, name, contracts.StatusInitialized, nil)
			}
			break
		}
	}

	return nil
}

// StartAll starts all registered plugins with lifecycle management
// zh: StartAll 啟動所有已註冊的插件並進行生命週期管理
func (lm *LifecycleManager) StartAll(ctx context.Context, registry contracts.Registry) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	if lm.status != contracts.StatusInitialized && lm.status != contracts.StatusStopped {
		return fmt.Errorf("lifecycle manager is not in a state to start (current: %s)", lm.status)
	}

	lm.status = contracts.StatusStarting
	lm.publishEvent(contracts.EventStart, "lifecycle-manager", lm.status, nil)

	// First, register logger plugin if available
	if err := lm.RegisterLoggerPlugin(ctx, registry); err != nil {
		log.L(ctx).Warn("Failed to register logger plugin", "error", err)
		// Continue with other plugins even if logger plugin fails
	}

	// Get all plugins from registry
	pluginNames := registry.ListPlugins()

	// Start each plugin
	for _, name := range pluginNames {
		plugin, err := registry.GetPlugin(name)
		if err != nil {
			lm.status = contracts.StatusError
			lm.publishEvent(contracts.EventError, name, contracts.StatusError, err)
			return fmt.Errorf("failed to get plugin %s: %w", name, err)
		}

		// Check if plugin is LifecycleAware and call OnStart
		if lifecyclePlugin, ok := plugin.(contracts.LifecycleAware); ok {
			if err := lifecyclePlugin.OnStart(); err != nil {
				lm.status = contracts.StatusError
				lm.publishEvent(contracts.EventError, name, contracts.StatusError, err)
				return fmt.Errorf("failed to start plugin %s: %w", name, err)
			}
		}

		// Update component status
		comp := contracts.ComponentInfo{
			Name:   name,
			Type:   "plugin",
			Status: contracts.StatusRunning,
		}
		lm.components[name] = comp
		lm.publishEvent(contracts.EventStart, name, contracts.StatusRunning, nil)
	}

	lm.status = contracts.StatusRunning
	lm.publishEvent(contracts.EventStart, "lifecycle-manager", lm.status, nil)
	return nil
}

// ShutdownAll shuts down all plugins in reverse order
// zh: ShutdownAll 按反向順序關閉所有插件
func (lm *LifecycleManager) ShutdownAll(ctx context.Context, registry contracts.Registry) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()

	lm.status = contracts.StatusShuttingDown
	lm.publishEvent(contracts.EventShutdown, "lifecycle-manager", lm.status, nil)

	// Get all plugin names in reverse order
	pluginNames := registry.ListPlugins()

	// Shutdown each plugin in reverse order
	for i := len(pluginNames) - 1; i >= 0; i-- {
		name := pluginNames[i]
		plugin, err := registry.GetPlugin(name)
		if err != nil {
			lm.publishEvent(contracts.EventError, name, contracts.StatusError, err)
			// Continue with other plugins even if one fails
			continue
		}

		// Check if plugin is LifecycleAware and call OnShutdown
		if lifecyclePlugin, ok := plugin.(contracts.LifecycleAware); ok {
			if err := lifecyclePlugin.OnShutdown(); err != nil {
				lm.publishEvent(contracts.EventError, name, contracts.StatusError, err)
				// Continue with other plugins even if one fails
				continue
			}
		}

		// Call plugin's Shutdown method
		if err := plugin.Shutdown(); err != nil {
			lm.publishEvent(contracts.EventError, name, contracts.StatusError, err)
			// Continue with other plugins even if one fails
			continue
		}

		// Update component status
		if comp, exists := lm.components[name]; exists {
			comp.Status = contracts.StatusShutdown
			lm.components[name] = comp
		}
		lm.publishEvent(contracts.EventShutdown, name, contracts.StatusShutdown, nil)
	}

	lm.status = contracts.StatusShutdown
	lm.publishEvent(contracts.EventShutdown, "lifecycle-manager", lm.status, nil)
	return nil
}

// HealthCheck performs health checks on all registered plugins
// zh: HealthCheck 對所有已註冊的插件執行健康檢查
func (lm *LifecycleManager) HealthCheck(ctx context.Context, registry contracts.Registry) map[string]contracts.HealthStatus {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	healthStatus := make(map[string]contracts.HealthStatus)
	pluginNames := registry.ListPlugins()

	for _, name := range pluginNames {
		plugin, err := registry.GetPlugin(name)
		if err != nil {
			healthStatus[name] = contracts.HealthStatus{
				Status:    "unhealthy",
				Message:   fmt.Sprintf("Failed to get plugin: %v", err),
				Timestamp: time.Now(),
			}
			continue
		}

		// Check if plugin implements HealthChecker interface
		if healthChecker, ok := plugin.(contracts.HealthChecker); ok {
			status := healthChecker.CheckHealth(ctx)
			healthStatus[name] = status
		} else {
			// Default health status for plugins that don't implement HealthChecker
			healthStatus[name] = contracts.HealthStatus{
				Status:    "healthy",
				Message:   "Plugin is running (no health check implemented)",
				Timestamp: time.Now(),
			}
		}
	}

	return healthStatus
}

// GetComponentStatus returns the status of a specific component
// zh: GetComponentStatus 回傳特定組件的狀態
func (lm *LifecycleManager) GetComponentStatus(name string) (contracts.ComponentInfo, bool) {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	comp, exists := lm.components[name]
	return comp, exists
}

// ListComponents returns all registered components
// zh: ListComponents 回傳所有已註冊的組件
func (lm *LifecycleManager) ListComponents() []contracts.ComponentInfo {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()

	components := make([]contracts.ComponentInfo, 0, len(lm.components))
	for _, comp := range lm.components {
		components = append(components, comp)
	}
	return components
}
