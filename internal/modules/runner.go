package modules

import (
	"context"
	"fmt"
)

// Runner coordinates the startup and shutdown of all modules based on dependencies.
// zh: Runner 根據模組依賴關係協調所有模組的啟動與停止。
type Runner struct {
	engine   Engine
	registry *Registry
	graph    *DependencyGraph
	started  []string
}

// NewRunner creates a new module runner.
// zh: 建立模組啟動與關閉的協調器。
func NewRunner(engine Engine, registry *Registry, graph *DependencyGraph) *Runner {
	return &Runner{
		engine:   engine,
		registry: registry,
		graph:    graph,
	}
}

// Start starts all registered modules in topological order.
// zh: 依拓撲排序啟動所有模組。
func (r *Runner) Start(ctx context.Context) error {
	order, err := r.graph.TopologicalSort()
	if err != nil {
		return fmt.Errorf("failed to resolve dependencies: %w", err)
	}

	for _, name := range order {
		module, ok := r.registry.Get(name)
		if !ok {
			return fmt.Errorf("module %q not found in registry", name)
		}
		if err := module.Run(ctx); err != nil {
			return fmt.Errorf("failed to start module %q: %w", name, err)
		}
		r.started = append(r.started, name)
	}

	return nil
}

// Stop stops all started modules in reverse order.
// zh: 依反向順序關閉所有模組。
func (r *Runner) Stop(ctx context.Context) error {
	for i := len(r.started) - 1; i >= 0; i-- {
		name := r.started[i]
		module, ok := r.registry.Get(name)
		if ok {
			if err := module.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed to stop module %q: %w", name, err)
			}
		}
	}
	return nil
}
