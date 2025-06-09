package registry

import (
	"fmt"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/registry"
)

// Engine 提供資源註冊的中心控制器。
// zh: Engine 是整個 GVK registry 的主控單元，統一管理註冊與查詢。
type Engine struct {
	mu       sync.RWMutex
	handlers map[registry.GVK]registry.ResourceHandler
	schemas  map[registry.GVK]string // optional: schema file path
}

// NewEngine 建立一個新的資源註冊 Engine。
// zh: 初始化 Engine 並準備註冊表。
func NewEngine() *Engine {
	return &Engine{
		handlers: make(map[registry.GVK]registry.ResourceHandler),
		schemas:  make(map[registry.GVK]string),
	}
}

// RegisterHandler 註冊 GVK 與對應資源處理器。
// zh: 登記特定 GVK 對應的 CRUD 實作。
func (e *Engine) RegisterHandler(gvk registry.GVK, h registry.ResourceHandler) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.handlers[gvk]; exists {
		return fmt.Errorf("GVK already registered: %v", gvk)
	}
	e.handlers[gvk] = h
	return nil
}

// RegisterSchema 註冊 GVK 與對應的 schema 檔案路徑。
// zh: 登記特定 GVK 對應的 schema 定義檔。
func (e *Engine) RegisterSchema(gvk registry.GVK, path string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.schemas[gvk] = path
}

// Handler 取得指定 GVK 的處理器。
// zh: 查詢指定 GVK 的 CRUD handler。
func (e *Engine) Handler(gvk registry.GVK) (registry.ResourceHandler, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	h, ok := e.handlers[gvk]
	if !ok {
		return nil, fmt.Errorf("handler not found for GVK: %v", gvk)
	}
	return h, nil
}

// SchemaPath 回傳指定 GVK 的 schema 路徑。
// zh: 查詢指定 GVK 的 schema 檔案。
func (e *Engine) SchemaPath(gvk registry.GVK) (string, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	path, ok := e.schemas[gvk]
	return path, ok
}
