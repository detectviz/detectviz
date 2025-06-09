package registry

import (
	"context"
	"fmt"
	"sync"

	ifacereg "github.com/detectviz/detectviz/pkg/ifaces/registry"
)

// memoryRegistry provides an in-memory implementation of the Registry interface.
// zh: memoryRegistry 是 Registry 介面的記憶體實作。
type memoryRegistry struct {
	mu       sync.RWMutex
	handlers map[ifacereg.GVK]ifacereg.ResourceHandler
}

// NewMemoryRegistry creates a new in-memory registry instance.
// zh: 建立一個新的記憶體註冊表實例。
func NewMemoryRegistry() ifacereg.Registry {
	return &memoryRegistry{
		handlers: make(map[ifacereg.GVK]ifacereg.ResourceHandler),
	}
}

// Register registers a resource handler for a given GVK.
// zh: 註冊對應 GVK 的資源操作處理器。
func (r *memoryRegistry) Register(gvk ifacereg.GVK, handler ifacereg.ResourceHandler) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.handlers[gvk]; exists {
		return fmt.Errorf("GVK already registered: %v", gvk)
	}

	r.handlers[gvk] = handler
	return nil
}

// Get retrieves a resource by GVK and name.
// zh: 根據 GVK 與名稱取得資源。
func (r *memoryRegistry) Get(ctx context.Context, gvk ifacereg.GVK, name string) (ifacereg.Resource, error) {
	handler, err := r.getHandler(gvk)
	if err != nil {
		return nil, err
	}
	return handler.Get(ctx, name)
}

// List lists all resources under a GVK.
// zh: 列出指定 GVK 下所有資源。
func (r *memoryRegistry) List(ctx context.Context, gvk ifacereg.GVK) ([]ifacereg.Resource, error) {
	handler, err := r.getHandler(gvk)
	if err != nil {
		return nil, err
	}
	return handler.List(ctx)
}

// Create adds a new resource under the given GVK.
// zh: 在指定 GVK 下新增資源。
func (r *memoryRegistry) Create(ctx context.Context, gvk ifacereg.GVK, res ifacereg.Resource) error {
	handler, err := r.getHandler(gvk)
	if err != nil {
		return err
	}
	return handler.Create(ctx, res)
}

// Update updates an existing resource.
// zh: 更新指定資源內容。
func (r *memoryRegistry) Update(ctx context.Context, gvk ifacereg.GVK, res ifacereg.Resource) error {
	handler, err := r.getHandler(gvk)
	if err != nil {
		return err
	}
	return handler.Update(ctx, res)
}

// Delete removes a resource by name.
// zh: 根據名稱刪除指定資源。
func (r *memoryRegistry) Delete(ctx context.Context, gvk ifacereg.GVK, name string) error {
	handler, err := r.getHandler(gvk)
	if err != nil {
		return err
	}
	return handler.Delete(ctx, name)
}

// getHandler looks up a registered handler for the given GVK.
// zh: 根據 GVK 尋找對應的資源處理器。
func (r *memoryRegistry) getHandler(gvk ifacereg.GVK) (ifacereg.ResourceHandler, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handler, exists := r.handlers[gvk]
	if !exists {
		return nil, fmt.Errorf("GVK not registered: %v", gvk)
	}
	return handler, nil
}
