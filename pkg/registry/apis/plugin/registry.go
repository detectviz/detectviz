package plugin

import (
	"fmt"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/plugins"
)

// Registry 提供插件註冊表。
// zh: 管理所有 Plugin 實體的集中註冊系統。
type Registry struct {
	mu      sync.RWMutex
	entries map[string]plugins.Plugin
}

// New 建立新的插件註冊表。
// zh: 初始化 Registry 實體。
func New() *Registry {
	return &Registry{
		entries: make(map[string]plugins.Plugin),
	}
}

// Register 註冊插件。
// zh: 若名稱已存在則回傳錯誤。
func (r *Registry) Register(p plugins.Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	name := p.Name()
	if _, exists := r.entries[name]; exists {
		return fmt.Errorf("plugin already registered: %s", name)
	}
	r.entries[name] = p
	return nil
}

// Get 查詢指定名稱的插件。
// zh: 若找不到則回傳 false。
func (r *Registry) Get(name string) (plugins.Plugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.entries[name]
	return p, ok
}

// List 回傳所有已註冊插件。
// zh: 回傳所有已登錄的 Plugin 清單。
func (r *Registry) List() []plugins.Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]plugins.Plugin, 0, len(r.entries))
	for _, p := range r.entries {
		list = append(list, p)
	}
	return list
}
