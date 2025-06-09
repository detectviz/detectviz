// Package modules provides global module registration and management.
// zh: modules 套件提供全域模組註冊與管理功能。
package modules

import (
	"fmt"
	"sync"
)

// Registry maintains a list of named modules for global management.
// zh: Registry 用於管理所有具名模組的註冊與查找。
type Registry struct {
	mu      sync.RWMutex
	modules map[string]Module
}

// NewRegistry creates a new module registry.
// zh: 建立一個模組註冊器。
func NewRegistry() *Registry {
	return &Registry{
		modules: make(map[string]Module),
	}
}

// Register registers a named module to the registry.
// Returns an error if the name already exists.
// zh: 註冊具名模組，若名稱已存在則回傳錯誤。
func (r *Registry) Register(name string, m Module) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.modules[name]; exists {
		return fmt.Errorf("module %q already registered", name)
	}

	r.modules[name] = m
	return nil
}

// Get retrieves a registered module by name.
// zh: 透過名稱查找已註冊模組。
func (r *Registry) Get(name string) (Module, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	m, ok := r.modules[name]
	return m, ok
}

// List returns all registered module names.
// zh: 回傳所有已註冊的模組名稱。
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.modules))
	for name := range r.modules {
		names = append(names, name)
	}
	return names
}
