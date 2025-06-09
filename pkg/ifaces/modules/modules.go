package modules

import "context"

// LifecycleModule defines a basic lifecycle interface for modules.
// zh: 定義模組基本生命週期介面。
type LifecycleModule interface {
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// HealthCheckableModule adds health check capabilities to a lifecycle module.
// zh: 提供健康檢查能力的模組，擴充自 LifecycleModule。
type HealthCheckableModule interface {
	LifecycleModule
	Healthy() bool
}

// ModuleEngine coordinates registration and control of unnamed modules.
// zh: 控制匿名模組註冊與啟動的引擎。
type ModuleEngine interface {
	Register(m LifecycleModule)
	RunAll(ctx context.Context) error
	ShutdownAll(ctx context.Context) error
}

// ModuleRegistry manages named module registration and lookup.
// zh: 管理具名模組的註冊與查詢。
type ModuleRegistry interface {
	Register(name string, m LifecycleModule) error
	Get(name string) (LifecycleModule, bool)
	List() []string
}

// ModuleRunner controls the ordered startup and shutdown based on dependencies.
// zh: 根據依賴圖控制模組啟動與關閉順序。
type ModuleRunner interface {
	StartAll(ctx context.Context) error
	StopAll(ctx context.Context) error
}

// ModuleListener monitors health status of registered modules.
// zh: 監控模組健康狀態的監聽器。
type ModuleListener interface {
	Start(ctx context.Context)
	Stop()
}

// Engine defines the interface for module lifecycle management.
// zh: 模組生命週期管理介面。
type Engine interface {
	// Register registers a lifecycle module for management.
	// zh: 註冊模組以便生命週期管理。
	Register(m LifecycleModule)

	// RunAll starts all registered modules.
	// zh: 啟動所有已註冊的模組。
	RunAll(ctx context.Context) error

	// ShutdownAll gracefully shuts down all registered modules.
	// zh: 優雅地關閉所有已註冊的模組。
	ShutdownAll(ctx context.Context) error
}
