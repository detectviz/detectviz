package datasource

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
		Kind:    "Datasource",
	}
)

// datasourceResource represents a simplified datasource with only a name.
// zh: datasourceResource 是測試用簡化資源，僅含名稱。
type datasourceResource struct {
	Name string
}

func (d *datasourceResource) GetName() string {
	return d.Name
}

// memoryDatasourceHandler provides in-memory CRUD for datasourceResource.
// zh: datasource 記憶體處理器，實作 ResourceHandler。
type memoryDatasourceHandler struct {
	mu    sync.RWMutex
	store map[string]registry.Resource
}

func newMemoryDatasourceHandler() *memoryDatasourceHandler {
	return &memoryDatasourceHandler{
		store: make(map[string]registry.Resource),
	}
}

func (m *memoryDatasourceHandler) Get(ctx context.Context, name string) (registry.Resource, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res, ok := m.store[name]
	if !ok {
		return nil, fmt.Errorf("datasource not found: %s", name)
	}
	return res, nil
}

func (m *memoryDatasourceHandler) List(ctx context.Context) ([]registry.Resource, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := make([]registry.Resource, 0, len(m.store))
	for _, res := range m.store {
		list = append(list, res)
	}
	return list, nil
}

func (m *memoryDatasourceHandler) Create(ctx context.Context, res registry.Resource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := res.GetName()
	if _, exists := m.store[name]; exists {
		return fmt.Errorf("datasource already exists: %s", name)
	}
	m.store[name] = res
	return nil
}

func (m *memoryDatasourceHandler) Update(ctx context.Context, res registry.Resource) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	name := res.GetName()
	if _, exists := m.store[name]; !exists {
		return fmt.Errorf("datasource does not exist: %s", name)
	}
	m.store[name] = res
	return nil
}

func (m *memoryDatasourceHandler) Delete(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.store[name]; !exists {
		return fmt.Errorf("datasource does not exist: %s", name)
	}
	delete(m.store, name)
	return nil
}

// RegisterDatasource registers the datasource handler in the given Registry.
// zh: 註冊 Datasource GVK 與處理器至註冊表。
func RegisterDatasource(r registry.Registry) error {
	return r.Register(gvk, newMemoryDatasourceHandler())
}
