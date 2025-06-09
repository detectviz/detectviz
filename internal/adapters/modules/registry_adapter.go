package modulesadapter

import (
	core "github.com/detectviz/detectviz/internal/modules"
	iface "github.com/detectviz/detectviz/pkg/ifaces/modules"
)

// RegistryAdapter wraps core.Registry to implement iface.ModuleRegistry.
// zh: RegistryAdapter 包裝 core.Registry，使其實作 ModuleRegistry 介面。
type RegistryAdapter struct {
	registry *core.Registry
}

// NewRegistryAdapter constructs a new RegistryAdapter instance.
// zh: 建立新的 RegistryAdapter 實例。
func NewRegistryAdapter(r *core.Registry) *RegistryAdapter {
	return &RegistryAdapter{
		registry: r,
	}
}

// Register adds a named module to the registry.
// zh: 註冊具名模組至底層 registry。
func (a *RegistryAdapter) Register(name string, m iface.LifecycleModule) error {
	return a.registry.Register(name, m)
}

// Get retrieves a registered module by name.
// zh: 根據名稱查詢模組。
func (a *RegistryAdapter) Get(name string) (iface.LifecycleModule, bool) {
	return a.registry.Get(name)
}

// List returns all registered module names.
// zh: 回傳所有已註冊模組名稱。
func (a *RegistryAdapter) List() []string {
	return a.registry.List()
}
