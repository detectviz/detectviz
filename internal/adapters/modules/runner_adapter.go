package modulesadapter

import (
	"context"

	core "github.com/detectviz/detectviz/internal/modules"
)

// RunnerAdapter wraps core.Runner to implement iface.ModuleRunner.
// zh: RunnerAdapter 包裝 core.Runner，使其實作 ModuleRunner 介面。
type RunnerAdapter struct {
	runner *core.Runner
}

// NewRunnerAdapter constructs a new RunnerAdapter instance.
// zh: 建立新的 RunnerAdapter 實例。
func NewRunnerAdapter(r *core.Runner) *RunnerAdapter {
	return &RunnerAdapter{
		runner: r,
	}
}

// StartAll starts all modules based on dependency order.
// zh: 依據依賴關係啟動所有模組。
func (a *RunnerAdapter) StartAll(ctx context.Context) error {
	return a.runner.Start(ctx)
}

// StopAll stops all modules in reverse order.
// zh: 依照啟動順序反向關閉所有模組。
func (a *RunnerAdapter) StopAll(ctx context.Context) error {
	return a.runner.Stop(ctx)
}
