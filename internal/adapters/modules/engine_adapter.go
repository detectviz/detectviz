package modulesadapter

import (
	"context"

	core "github.com/detectviz/detectviz/internal/modules"
	iface "github.com/detectviz/detectviz/pkg/ifaces/modules"
)

// EngineAdapter wraps core.Engine to implement iface.ModuleEngine.
// zh: EngineAdapter 包裝 core.Engine，使其實作 ModuleEngine 介面。
type EngineAdapter struct {
	engine core.Engine
}

// NewEngineAdapter constructs a new EngineAdapter instance.
// zh: 建立新的 EngineAdapter 實例。
func NewEngineAdapter(e core.Engine) *EngineAdapter {
	return &EngineAdapter{
		engine: e,
	}
}

// Register adds a module to the core engine.
// zh: 註冊模組至底層 Engine 實例。
func (a *EngineAdapter) Register(m iface.LifecycleModule) {
	a.engine.Register(m)
}

// RunAll starts all registered modules via the core engine.
// zh: 啟動所有已註冊的模組。
func (a *EngineAdapter) RunAll(ctx context.Context) error {
	return a.engine.RunAll(ctx)
}

// ShutdownAll stops all modules via the core engine.
// zh: 關閉所有已註冊的模組。
func (a *EngineAdapter) ShutdownAll(ctx context.Context) error {
	return a.engine.ShutdownAll(ctx)
}
