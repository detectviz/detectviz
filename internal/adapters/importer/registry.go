package importeradapter

import (
	"fmt"
	"sync"

	"github.com/detectviz/detectviz/pkg/importer"
)

// Registry 提供 Importer 註冊機制。
// zh: 記錄所有已註冊的 Importer，供系統查詢與調用。
type Registry struct {
	mu        sync.RWMutex
	importers map[string]importer.Importer
}

// NewRegistry 建立 Importer Registry 實例。
func NewRegistry() *Registry {
	return &Registry{
		importers: make(map[string]importer.Importer),
	}
}

// Register 註冊一個 Importer。
func (r *Registry) Register(i importer.Importer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := i.Name()
	if _, exists := r.importers[name]; exists {
		return fmt.Errorf("importer already registered: %s", name)
	}
	r.importers[name] = i
	return nil
}

// Get 根據名稱取得 Importer。
func (r *Registry) Get(name string) (importer.Importer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	i, ok := r.importers[name]
	return i, ok
}

// List 回傳所有註冊的 Importer。
func (r *Registry) List() []importer.Importer {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]importer.Importer, 0, len(r.importers))
	for _, i := range r.importers {
		result = append(result, i)
	}
	return result
}
