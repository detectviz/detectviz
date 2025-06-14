package composition

import (
	"context"
	"fmt"
	"sync"
	"time"

	"detectviz/pkg/platform/contracts"
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
	}

	if err != nil {
		event.Error = err.Error()
	}

	for _, listener := range lm.listeners {
		// Fire and forget - don't block lifecycle operations on listener errors
		go func(l contracts.LifecycleListener) {
			_ = l.OnLifecycleEvent(event)
		}(listener)
	}
}
