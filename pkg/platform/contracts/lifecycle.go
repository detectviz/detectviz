package contracts

import (
	"context"
	"time"
)

// LifecycleManager defines the interface for managing plugin lifecycles.
// zh: LifecycleManager 定義管理插件生命週期的介面。
type LifecycleManager interface {
	Initialize(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Shutdown(ctx context.Context) error
	GetStatus() LifecycleStatus
}

// LifecycleStatus represents the current status of a component.
// zh: LifecycleStatus 代表組件的當前狀態。
type LifecycleStatus string

const (
	StatusUninitialized LifecycleStatus = "uninitialized"
	StatusInitializing  LifecycleStatus = "initializing"
	StatusInitialized   LifecycleStatus = "initialized"
	StatusStarting      LifecycleStatus = "starting"
	StatusRunning       LifecycleStatus = "running"
	StatusStopping      LifecycleStatus = "stopping"
	StatusStopped       LifecycleStatus = "stopped"
	StatusShuttingDown  LifecycleStatus = "shutting_down"
	StatusShutdown      LifecycleStatus = "shutdown"
	StatusError         LifecycleStatus = "error"
)

// LifecycleEvent represents a lifecycle event.
// zh: LifecycleEvent 代表生命週期事件。
type LifecycleEvent struct {
	Type      LifecycleEventType `json:"type"`
	Component string             `json:"component"`
	Timestamp time.Time          `json:"timestamp"`
	Status    LifecycleStatus    `json:"status"`
	Error     string             `json:"error,omitempty"`
	Metadata  map[string]any     `json:"metadata,omitempty"`
}

// LifecycleEventType represents the type of lifecycle event.
// zh: LifecycleEventType 代表生命週期事件類型。
type LifecycleEventType string

const (
	EventInitialize LifecycleEventType = "initialize"
	EventStart      LifecycleEventType = "start"
	EventStop       LifecycleEventType = "stop"
	EventShutdown   LifecycleEventType = "shutdown"
	EventError      LifecycleEventType = "error"
)

// LifecycleListener defines the interface for lifecycle event listeners.
// zh: LifecycleListener 定義生命週期事件監聽器介面。
type LifecycleListener interface {
	OnLifecycleEvent(event LifecycleEvent) error
}

// DependencyResolver defines the interface for resolving component dependencies.
// zh: DependencyResolver 定義解析組件依賴關係的介面。
type DependencyResolver interface {
	ResolveDependencies(components []ComponentInfo) ([]ComponentInfo, error)
	ValidateDependencies(components []ComponentInfo) error
	GetDependencyGraph() DependencyGraph
}

// ComponentInfo represents information about a component.
// zh: ComponentInfo 代表組件資訊。
type ComponentInfo struct {
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	Dependencies []string        `json:"dependencies"`
	Priority     int             `json:"priority"`
	Config       map[string]any  `json:"config"`
	Status       LifecycleStatus `json:"status"`
}

// DependencyGraph represents the dependency relationships between components.
// zh: DependencyGraph 代表組件間的依賴關係圖。
type DependencyGraph interface {
	AddNode(name string, info ComponentInfo) error
	AddEdge(from, to string) error
	GetTopologicalOrder() ([]string, error)
	HasCycle() bool
	GetDependents(name string) []string
	GetDependencies(name string) []string
}

// HealthChecker defines the interface for component health checking.
// zh: HealthChecker 定義組件健康檢查介面。
type HealthChecker interface {
	CheckHealth(ctx context.Context) HealthStatus
	GetHealthMetrics() map[string]any
}

// HealthStatus represents the health status of a component.
// zh: HealthStatus 代表組件的健康狀態。
type HealthStatus struct {
	Status    string         `json:"status"` // healthy, unhealthy, degraded
	Message   string         `json:"message"`
	Timestamp time.Time      `json:"timestamp"`
	Details   map[string]any `json:"details,omitempty"`
}
