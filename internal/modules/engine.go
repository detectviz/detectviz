package modules

import (
	"context"
	"fmt"
	"sync"
)

// Module defines the standard lifecycle interface for all modules.
// zh: 定義所有模組的統一生命週期介面。
type Module interface {
	// Run starts the module and blocks until it is stopped or encounters an error.
	// zh: 啟動模組，阻塞直到停止或發生錯誤。
	Run(ctx context.Context) error

	// Shutdown gracefully stops the module.
	// zh: 優雅地關閉模組。
	Shutdown(ctx context.Context) error
}

// Engine defines the interface for module lifecycle management.
// zh: 模組生命週期管理介面。
type Engine interface {
	Register(m Module)
	RunAll(ctx context.Context) error
	ShutdownAll(ctx context.Context) error
}

// engineImpl 是 Engine 介面的實作結構。
// zh: Engine 的內部實作。
type engineImpl struct {
	modules []Module
	mu      sync.Mutex
}

// NewEngine creates a new Engine instance to manage modules.
// zh: 建立新的模組控制引擎。
func NewEngine() Engine {
	return &engineImpl{
		modules: make([]Module, 0),
	}
}

// Register adds a module to the Engine for lifecycle management.
// zh: 註冊一個模組到引擎中，由引擎管理其生命週期。
func (e *engineImpl) Register(m Module) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.modules = append(e.modules, m)
}

// RunAll starts all registered modules sequentially.
// It returns an error immediately if any module fails to start.
// zh: 順序啟動所有已註冊模組，若任一模組啟動失敗則立刻回傳錯誤。
func (e *engineImpl) RunAll(ctx context.Context) error {
	for _, m := range e.modules {
		if err := m.Run(ctx); err != nil {
			return err
		}
	}
	return nil
}

// ShutdownAll gracefully stops all registered modules.
// If multiple modules fail to shutdown, their errors are combined.
// zh: 優雅地關閉所有模組，若有多個錯誤會合併回傳。
func (e *engineImpl) ShutdownAll(ctx context.Context) error {
	var shutdownErr error

	for _, m := range e.modules {
		if err := m.Shutdown(ctx); err != nil {
			if shutdownErr == nil {
				shutdownErr = err
			} else {
				shutdownErr = fmt.Errorf("%w; %v", shutdownErr, err)
			}
		}
	}

	return shutdownErr
}
