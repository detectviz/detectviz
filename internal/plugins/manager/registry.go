package manager

import (
	"fmt"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/plugins"
)

// ManagerRegistry 管理所有已載入的 plugins。
// zh: 負責註冊、查詢與統一管理所有插件實體。
type ManagerRegistry struct {
	mu      sync.RWMutex
	entries map[string]plugins.Plugin
}

// NewManagerRegistry 建立新的 plugin 註冊管理器。
// zh: 初始化 plugin 註冊資料表。
func NewManagerRegistry() *ManagerRegistry {
	return &ManagerRegistry{
		entries: make(map[string]plugins.Plugin),
	}
}

// Register 登錄插件至管理器。
// zh: 若名稱重複則回傳錯誤。
func (r *ManagerRegistry) Register(p plugins.Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.entries[p.Name()]; ok {
		return fmt.Errorf("plugin already exists: %s", p.Name())
	}
	r.entries[p.Name()] = p
	return nil
}

// Get 查詢插件。
// zh: 取得指定名稱的 plugin。
func (r *ManagerRegistry) Get(name string) (plugins.Plugin, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.entries[name]
	return p, ok
}

// List 回傳所有已註冊插件。
// zh: 回傳所有插件清單。
func (r *ManagerRegistry) List() []plugins.Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []plugins.Plugin
	for _, p := range r.entries {
		list = append(list, p)
	}
	return list
}
