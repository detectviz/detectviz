package libraryelements

import (
	"fmt"
	"sync"
)

// Registry 提供元件類型與渲染器註冊功能。
// zh: 管理 ElementRenderer 註冊與查詢。
type Registry struct {
	mu        sync.RWMutex
	renderers map[string]ElementRenderer
}

// NewRegistry 建立空白元件註冊表。
// zh: 初始化內部映射結構。
func NewRegistry() *Registry {
	return &Registry{
		renderers: make(map[string]ElementRenderer),
	}
}

// RegisterRenderer 註冊指定元件類型的渲染器。
// zh: 為指定 Type 註冊轉譯處理器。
func (r *Registry) RegisterRenderer(kind string, renderer ElementRenderer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.renderers[kind]; exists {
		return fmt.Errorf("renderer already registered for type: %s", kind)
	}
	r.renderers[kind] = renderer
	return nil
}

// GetRenderer 取得指定元件類型的渲染器。
// zh: 根據元件類型查詢對應轉譯器。
func (r *Registry) GetRenderer(kind string) (ElementRenderer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	renderer, ok := r.renderers[kind]
	return renderer, ok
}
