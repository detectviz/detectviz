package host

import (
	"context"
	"fmt"
	"sync"

	"github.com/detectviz/detectviz/pkg/ifaces/registry"
)

var (
	gvk = registry.GVK{
		Group:   "core",
		Version: "v1",
		Kind:    "Host",
	}
)

// hostResource represents a simple host with a name field.
// zh: hostResource 是用來測試的基本資源，僅包含名稱。
type hostResource struct {
	Name string
}

func (h *hostResource) GetName() string {
	return h.Name
}

// memoryHostHandler is an in-memory implementation of ResourceHandler for Host.
// zh: 用 map 儲存 Host 資源的 CRUD 操作實作。
type memoryHostHandler struct {
	mu    sync.RWMutex
	store map[string]registry.Resource
}

func newMemoryHostHandler() *memoryHostHandler {
	return &memoryHostHandler{
		store: make(map[string]registry.Resource),
	}
}

func (m *memoryHostHandler) Get(ctx context.Context, name string) (registry.Resource, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res, ok := m.store[name]
	if !ok {
		return nil, fmt.Errorf("host not found: %s", name)
	}
	return res, nil
}

func (m *memoryHostHandler) List(ctx context.Context) ([]registry.Resource, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]registry.Resource, 0, len(m.store))
	for _, res := range m.store {
		list = append(list, res)
	}
	return list, nil
}

func (m *memoryHostHandler) Create(ctx context.Context, res registry.Resource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := res.GetName()
	if _, exists := m.store[name]; exists {
		return fmt.Errorf("host already exists: %s", name)
	}
	m.store[name] = res
	return nil
}

func (m *memoryHostHandler) Update(ctx context.Context, res registry.Resource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := res.GetName()
	if _, exists := m.store[name]; !exists {
		return fmt.Errorf("host does not exist: %s", name)
	}
	m.store[name] = res
	return nil
}

func (m *memoryHostHandler) Delete(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.store[name]; !exists {
		return fmt.Errorf("host does not exist: %s", name)
	}
	delete(m.store, name)
	return nil
}

// RegisterHost registers the Host resource handler into the given registry.
// zh: 註冊 Host 資源處理器到指定的 Registry。
func RegisterHost(r registry.Registry) error {
	return r.Register(gvk, newMemoryHostHandler())
}
